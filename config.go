package firemodel

import "fmt"

type Schema struct {
	Models     []*SchemaModel
	Enums      []*SchemaEnum
	Structs    []*SchemaStruct
	Interfaces []*SchemaInterface
}

type Config struct {
	Languages           []Language
	SourceCoderProvider func(prefix string) SourceCoder
}

type SchemaModel struct {
	Name          string
	FirestorePath SchemaModelPathTemplate
	Implements    []*SchemaInterface
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

type SchemaInterface struct {
	Name    string
	Comment string
	Fields  []*SchemaField
}

type SchemaStruct struct {
	Name       string
	Comment    string
	Implements []*SchemaInterface
	Fields     []*SchemaField
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
}

type SchemaFieldType interface {
	isSchemaTypeName()
}

func equalSchemaFieldTypes(a, b SchemaFieldType) bool {
	switch a := a.(type) {
	case *Boolean:
		_, ok := b.(*Boolean)
		return ok
	case *Integer:
		_, ok := b.(*Integer)
		return ok
	case *Double:
		_, ok := b.(*Double)
		return ok
	case *GeoPoint:
		_, ok := b.(*GeoPoint)
		return ok
	case *Timestamp:
		_, ok := b.(*Timestamp)
		return ok
	case *String:
		_, ok := b.(*String)
		return ok
	case *Bytes:
		_, ok := b.(*Bytes)
		return ok
	case *Reference:
		b, ok := b.(*Reference)
		return ok && a.T.Name == b.T.Name
	case *Array:
		b, ok := b.(*Array)
		return ok && equalSchemaFieldTypes(a.T, b.T)
	case *Map:
		b, ok := b.(*Map)
		return ok && equalSchemaFieldTypes(a.T, b.T)
	case *Struct:
		b, ok := b.(*Struct)
		return ok && a.T.Name == b.T.Name
	case *Enum:
		b, ok := b.(*Enum)
		return ok && a.T.Name == b.T.Name
	case *URL:
		_, ok := b.(*URL)
		return ok
	case *File:
		_, ok := b.(*File)
		return ok
	}
	return false
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

func (t *Boolean) String() string   { return "boolean" }
func (t *Integer) String() string   { return "integer" }
func (t *Double) String() string    { return "double" }
func (t *GeoPoint) String() string  { return "geopoint" }
func (t *Timestamp) String() string { return "timestamp" }
func (t *String) String() string    { return "string" }
func (t *Bytes) String() string     { return "bytes" }
func (t *Reference) String() string { return fmt.Sprintf("reference<%s>", t.T.Name) }
func (t *Array) String() string     { return fmt.Sprintf("array<%s>", t.T) }
func (t *Map) String() string       { return fmt.Sprintf("map<%s>", t.T) }
func (t *Struct) String() string    { return t.T.Name }
func (t *Enum) String() string      { return t.T.Name }
func (t *URL) String() string       { return "URL" }
func (t *File) String() string      { return "File" }
