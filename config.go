package firemodel

import (
	"strings"
	"github.com/pkg/errors"
	"regexp"
)

type Schema struct {
	Models  []*SchemaModel
	Enums   []*SchemaEnum
	Options SchemaOptions
}

type Config struct {
	Languages           []Language
	SourceCoderProvider func(prefix string) SourceCoder
}

type SchemaModel struct {
	Name        string
	Comment     string
	Fields      []*SchemaField
	Collections []*SchemaNestedCollection
	Options     SchemaModelOptions
}

type SchemaOptions map[string]map[string]string

func (options SchemaOptions) Get(key string) map[string]string {
	if res, ok := options[key]; ok {
		return res
	}
	return map[string]string{}
}

type SchemaModelOptions SchemaOptions

func (options SchemaModelOptions) Get(key string) map[string]string {
	if res, ok := options[key]; ok {
		return res
	}
	return map[string]string{}
}

var (
	firestorePathVariablePattern = regexp.MustCompile("^{([a-zA-Z0-9_-]+)}$")
	firestorePathConstantPattern = regexp.MustCompile("^([a-zA-Z0-9_-]+)$")
)

// GetFirestorePath returns the templetized Firestore path where this model is located in Firestore.
//
// This method requires that the model includes an option called firestore.path.
//
// The path may include variables, wrapped in curly brackets: e.g. `users/{user_id}`. Variables are
// replaced with %s, so that they can be interpolated by printf functions. vars provides the names
// of these interpolation variables.
func (options SchemaModelOptions) GetFirestorePath() (format string, args []string, err error) {
	pathTemplate, ok := options.Get("firestore")["path"]
	if !ok {
		return
	}
	if len(pathTemplate) == 0 {
		err = errors.Errorf(`firemodel: invalid path option "%s"`, pathTemplate)
		return
	}

	components := strings.Split(pathTemplate, "/")
	for idx, component := range components {
		if firestorePathConstantPattern.MatchString(component) {
			continue
		} else if variableComponents := firestorePathVariablePattern.FindStringSubmatch(component); variableComponents != nil {
			args = append(args, variableComponents[1])
			components[idx] = "%s"
		} else {
			err = errors.Errorf(`firemodel: invalid path option "%s (component=%s)"`, pathTemplate, component)
			return
		}
	}
	format = strings.Join(components, "/")
	return
}

type SchemaEnum struct {
	Name    string
	Comment string
	Values  []*SchemaEnumValue
}

type SchemaField struct {
	Name    string
	Comment string
	Type    SchemaFieldType
	Extras  *SchemaFieldExtras
}

type SchemaFieldType string

type SchemaEnumValue struct {
	Name    string
	Comment string
}

type SchemaFieldExtras struct {
	ReferenceTo      string
	ArrayOfPrimitive SchemaFieldType
	ArrayOf          string
	MapToPrimitive   SchemaFieldType
	MapTo            string
	EnumType         string
	URL              bool
	File             bool
}

type SchemaNestedCollection struct {
	Name    string
	Comment string
	Type    string
}
