package firemodel

import (
	"regexp"
)

type Schema struct {
	Models  []*SchemaModel
	Enums   []*SchemaEnum
	Structs []*SchemaStruct
}

type Config struct {
	Languages           []Language
	SourceCoderProvider func(prefix string) SourceCoder
}

type SchemaModel struct {
	Name          string
	FirestorePath SchemaModelPathTemplate
	Comment       string
	Fields        []*SchemaField
}

type SchemaModelPathTemplate struct {
	Pattern         string
	CollectionParts []SchemaModelPathTemplatePart
}

type SchemaModelPathTemplatePart struct {
	CollectionName      string
	DocumentPlaceholder string
}

func (pt SchemaModelPathTemplate) String() string {
	return pt.Pattern
}

type SchemaStruct struct {
	Name    string
	Comment string
	Fields  []*SchemaField
}

var (
	firestorePathVariablePattern = regexp.MustCompile("^{([a-zA-Z0-9_-]+)}$")
	firestorePathConstantPattern = regexp.MustCompile("^([a-zA-Z0-9_-]+)$")
)

type SchemaEnum struct {
	Name    string
	Comment string
	Values  []*SchemaEnumValue
}

type SchemaField struct {
	Name    string
	Comment string
	Type    SchemaFieldType
}

type SchemaFieldType interface {
	isSchemaTypeName()
}

type SchemaEnumValue struct {
	Name            string
	Comment         string
	AssociatedValue *Struct
}

type Boolean struct{}
type Integer struct{}
type Double struct{}
type GeoPoint struct{}
type Timestamp struct{}
type String struct{}
type Bytes struct{}
type Reference struct{ T *SchemaModel }
type Array struct{ T SchemaFieldType }
type Map struct{ T SchemaFieldType }
type Struct struct{ T *SchemaStruct }
type Enum struct{ T *SchemaEnum }
type URL struct{}
type File struct{}

func (t *Boolean) isSchemaTypeName()   {}
func (t *Integer) isSchemaTypeName()   {}
func (t *Double) isSchemaTypeName()    {}
func (t *GeoPoint) isSchemaTypeName()  {}
func (t *Timestamp) isSchemaTypeName() {}
func (t *String) isSchemaTypeName()    {}
func (t *Bytes) isSchemaTypeName()     {}
func (t *Reference) isSchemaTypeName() {}
func (t *Array) isSchemaTypeName()     {}
func (t *Map) isSchemaTypeName()       {}
func (t *Struct) isSchemaTypeName()    {}
func (t *Enum) isSchemaTypeName()      {}
func (t *URL) isSchemaTypeName()       {}
func (t *File) isSchemaTypeName()      {}
