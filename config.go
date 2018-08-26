package firemodel

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
	Options     SchemaOptions
}

type SchemaOptions map[string]map[string]string

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
