package firemodel

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type Schema struct {
	Models  []*SchemaModel
	Enums   []*SchemaEnum
	Structs []*SchemaStruct
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

type SchemaStruct struct {
	Name    string
	Comment string
	Fields  []*SchemaField
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

// GetFirestorePath returns the templatized Firestore path where this model is located in Firestore.
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
	if len(components)%2 != 0 {
		err = errors.Errorf(`firemodel: invalid path option (must be even number of components) "%s"`, pathTemplate)
		return
	}

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

func (options SchemaModelOptions) GetFirestoreModelName() (modelName string, err error) {
	modelName, ok := options.Get("firestore")["model_name"]
	if !ok {
		return
	}

	if len(modelName) == 0 {
		err = errors.Errorf(`firemodel: invalid model name "%s"`, modelName)
		return
	}

	return modelName, nil
}

func (options SchemaModelOptions) GetAutoTimestamp() bool {
	if autoTimestamp, ok := options.Get("firestore")["autotimestamp"]; ok && autoTimestamp != "false" {

		return true
	}
	return false
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

type SchemaNestedCollection struct {
	Name    string
	Comment string
	Type    *SchemaModel
}
