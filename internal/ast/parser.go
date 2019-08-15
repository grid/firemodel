package ast

import (
	"fmt"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/scanner"
)

var commentPattern = regexp.MustCompile(`//\s*(?P<Comment>.*)$|/\*(?P<BlockComment>(?s:.*))\*/`)

func ParseSchema(r io.Reader) (*AST, error) {
	parser := participle.MustBuild(
		&AST{},
		participle.Lexer(&lexerDefinition{}),
		participle.Map(func(token lexer.Token) (lexer.Token, error) {
			if os.Getenv("DEBUG") != "" {
				log.Print(scanner.TokenString(rune(token.Type)), " -> ", token.Value)
			}
			if token.Type == scanner.Comment {
				submatch := commentPattern.FindStringSubmatch(token.Value)
				if lineComment := submatch[1]; lineComment != "" {
					token.Value = lineComment
				} else if blockComment := submatch[2]; blockComment != "" {
					lines := strings.Split(blockComment, "\n")
					for lineIdx, line := range lines {
						lines[lineIdx] = strings.TrimLeft(strings.TrimSpace(line), "*")
					}
					token.Value = strings.Join(lines, "\n")
				}
			}
			return token, nil
		}),
	)

	s := &AST{}
	if err := parser.Parse(r, s); err != nil {
		return nil, errors.Wrap(err, "firemodel syntax error")
	}
	return s, nil
}

type lexerDefinition struct{}

func (d *lexerDefinition) Lex(r io.Reader) (lexer.Lexer, error) {
	s := &scanner.Scanner{}
	l := lexer.LexWithScanner(r, s)
	s.Mode = scanner.GoTokens ^ scanner.SkipComments
	return l, nil
}

func (d *lexerDefinition) Symbols() map[string]rune {
	return map[string]rune{
		"EOF":       scanner.EOF,
		"Ident":     scanner.Ident,
		"String":    scanner.String,
		"RawString": scanner.RawString,
		"Int":       scanner.Int,
		"Comment":   scanner.Comment,
	}
}

// AST represents the structure of a parsed firemodel schema file.
//
// Read about the magical annotations here: https://github.com/alecthomas/participle/.
type AST struct {
	Types []*ASTElement `parser:"( @@ | Comment )*"`
}

type ASTElement struct {
	Comment   string        `parser:"( @Comment )?"`
	Model     *ASTModel     `parser:"( 'model' @@"`
	Interface *ASTInterface `parser:"| 'interface' @@"`
	Enum      *ASTEnum      `parser:"| 'enum' @@"`
	Option    *ASTOption    `parser:"| 'option' @@"`
	Struct    *ASTStruct    `parser:"| 'struct' @@ )"`
}

type ASTModel struct {
	Identifier   ASTIdentifier      `parser:"@Ident ':'"`
	PathTemplate *ASTPathTemplate   `parser:"@@"`
	Implements   []ASTType          `parser:"( 'implements' ( @Ident ',' )* @Ident )?"`
	Elements     []*ASTModelElement `parser:"'{' ( @@  | Comment  )* '}'"`
}

type ASTPathTemplatePart struct {
	CollectionName      string
	DocumentPlaceholder string
}

type ASTPathTemplate struct {
	Pattern         string
	CollectionParts []ASTPathTemplatePart
}

func (pt *ASTPathTemplate) Parse(lex lexer.PeekingLexer) error {
	tok, err := lex.Next()
	if err != nil {
		return err
	}
	if tok.Type != scanner.String {
		return participle.NextMatch
	}
	if len(tok.Value) == 0 {
		return errors.New("missing path pattern")
	}
	parts := strings.Split(tok.Value, "/")
	if parts[0] != "" {
		return errors.Errorf("path pattern should start with /: %s", tok)
	}
	if len(parts)%2 != 1 {
		return errors.Errorf("bad path pattern: %s", tok)
	}
	for k, v := range parts[1:] { // start after blank for leading /
		switch k % 2 {
		case 0: //collections
			if !collectionPattern.MatchString(v) {
				return errors.Errorf("bad collection %s in pattern %s", v, tok.Value)
			}
		case 1: // documents
			if !documentTemplatePattern.MatchString(v) {
				return errors.Errorf("bad document template %s in pattern %s", v, tok.Value)
			}
			pt.CollectionParts = append(pt.CollectionParts, ASTPathTemplatePart{
				CollectionName:      parts[k],
				DocumentPlaceholder: parts[k+1],
			})
		}
	}
	pt.Pattern = tok.Value
	return nil
}

var (
	collectionPattern       = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	documentTemplatePattern = regexp.MustCompile(`^{[a-zA-Z0-9_-]+}$`)
)

type ASTInterface struct {
	Identifier ASTIdentifier       `parser:"@Ident"`
	Elements   []*ASTStructElement `parser:"'{' ( @@ )* '}'"`
}

type ASTStruct struct {
	Identifier ASTIdentifier       `parser:"@Ident"`
	Implements []ASTType           `parser:"( 'implements' ( @Ident ',' )* @Ident )?"`
	Elements   []*ASTStructElement `parser:"'{' ( @@  | Comment  )* '}'"`
}

type ASTIdentifier string

var (
	reservedIdentifiers = []string{
		// Primitive types.
		"boolean", "integer", "double", "timestamp", "string",
		"bytes", "reference", "geopoint", "array", "map", "url",
		"collection",
		// Keywords.
		"model", "option", "enum", "implements", "struct", "at",
	}
)

func init() {
	sort.Strings(reservedIdentifiers)
}

func (id ASTIdentifier) IsReserved() bool {
	needle := strings.ToLower(string(id))
	idx := sort.SearchStrings(reservedIdentifiers, needle)
	if idx >= len(reservedIdentifiers) {
		return false
	}
	if reservedIdentifiers[idx] != needle {
		return false
	}
	return true
}

type ASTStructElement struct {
	Field *ASTField `parser:"@@"`
}

type ASTModelElement struct {
	Option *ASTOption `parser:"  'option' @@"`
	Field  *ASTField  `parser:"| @@"`
}

type ASTEnum struct {
	Identifier ASTIdentifier   `parser:"@Ident '{'"`
	Values     []*ASTEnumValue `parser:"( ( @@ ',' ) | Comment )* '}'"`
}

type ASTOption struct {
	Language string        `parser:"@Ident '.'"`
	Key      ASTIdentifier `parser:"@Ident '='"`
	Value    string        `parser:"@('true' | 'false' | 'null' | String | Int) ';'"`
}

type ASTEnumValue struct {
	Comment         string        `parser:"( @Comment )?"`
	Name            string        `parser:"@Ident"`
	AssociatedValue *ASTFieldType `parser:"( '(' @@ ')')?"`
}

type ASTField struct {
	Comment string        `parser:"( @Comment )?"`
	Type    *ASTFieldType `parser:"@@"`
	Name    string        `parser:"@Ident ';'"`
}

type ASTFieldType struct {
	Base    ASTType       `parser:"@Ident"`
	Generic *ASTFieldType `parser:"( '<' @@ '>' )?"`
}

func (ft *ASTFieldType) String() string {
	if ft.Generic != nil {
		return fmt.Sprintf("%s<%s>", ft.Base, ft.Generic)
	} else {
		return fmt.Sprint(ft.Base)
	}
}

func (ft *ASTFieldType) IsPrimitive() bool {
	switch ft.Base {
	case String,
		Integer,
		Bytes,
		Double,
		Timestamp,
		Boolean,
		Reference,
		GeoPoint,
		Array,
		Map:
		return true
	default:
		return false
	}
}

type ASTType string

const (
	Boolean   ASTType = "boolean"
	Integer   ASTType = "integer"
	Double    ASTType = "double"
	Timestamp ASTType = "timestamp"
	String    ASTType = "string"
	Bytes     ASTType = "bytes"
	Reference ASTType = "reference"
	GeoPoint  ASTType = "geopoint"
	Array     ASTType = "array"
	Map       ASTType = "map"
	// Fake types.
	URL ASTType = "URL"
)
