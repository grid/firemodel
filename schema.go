package firemodel

import (
	"io"

	"github.com/iancoleman/strcase"
	"github.com/mickeyreiss/firemodel/internal/ast"
	"github.com/pkg/errors"
)

func ParseSchema(r io.Reader) (*Schema, error) {
	tree, err := ast.ParseSchema(r)
	if err != nil {
		return nil, err
	}
	compiler := &configSchemaCompiler{ast: tree}
	config, err := compiler.compileConfig()
	if err != nil {
		return nil, errors.Wrap(err, "firemodel/schema")
	}
	return config, nil
}

type configSchemaCompiler struct {
	models  []*SchemaModel
	structs []*SchemaStruct
	enums   []*SchemaEnum

	ast *ast.AST
}

func (c *configSchemaCompiler) compileConfig() (*Schema, error) {
	if err := c.precompileEnumTypes(); err != nil {
		return nil, err
	}
	if err := c.precompileModelTypes(); err != nil {
		return nil, err
	}
	if err := c.precompileStructTypes(); err != nil {
		return nil, err
	}

	return &Schema{
		Models:  c.compileModels(),
		Enums:   c.compileEnums(),
		Structs: c.compileStructs(),
		Options: c.compileLanguageOptions(),
	}, nil
}

func (c *configSchemaCompiler) precompileEnumTypes() error {
	c.enums = make([]*SchemaEnum, 0)
	for _, v := range c.ast.Types {
		if v.Enum == nil {
			continue
		}

		if v.Enum.Identifier.IsReserved() {
			err := errors.Errorf("firemodel/schema: can't name enum %s, %s is a reserved word.", v.Enum.Identifier, v.Enum.Identifier)
			return err
		}

		c.enums = append(c.enums, &SchemaEnum{
			Name: strcase.ToCamel(string(v.Enum.Identifier)),
		})
	}
	return nil
}

func (c *configSchemaCompiler) precompileModelTypes() error {
	c.models = make([]*SchemaModel, 0)
	for _, v := range c.ast.Types {
		if v.Model == nil {
			continue
		}

		if v.Model.Identifier.IsReserved() {
			err := errors.Errorf("firemodel/schema: can't name model %s, %s is a reserved word.", v.Model.Identifier, v.Model.Identifier)
			return err
		}

		c.models = append(c.models, &SchemaModel{
			Name: strcase.ToCamel(string(v.Model.Identifier)),
		})
	}
	return nil
}

func (c *configSchemaCompiler) precompileStructTypes() error {
	c.structs = make([]*SchemaStruct, 0)
	for _, v := range c.ast.Types {
		if v.Struct == nil {
			continue
		}

		if v.Struct.Identifier.IsReserved() {
			err := errors.Errorf("firemodel/schema: can't name struct %s, %s is a reserved word.", v.Struct.Identifier, v.Struct.Identifier)
			return err
		}

		c.structs = append(c.structs, &SchemaStruct{
			Name: strcase.ToCamel(string(v.Struct.Identifier)),
		})
	}
	return nil
}

func (c *configSchemaCompiler) compileModels() (out []*SchemaModel) {
	for _, v := range c.ast.Types {
		if v.Model == nil {
			continue
		}

		if v.Model.Identifier.IsReserved() {
			err := errors.Errorf("firemodel/schema: can't name model %s, %s is a reserved word.", v.Model.Identifier, v.Model.Identifier)
			panic(err)
		}

		out = append(out, &SchemaModel{
			Name:        strcase.ToCamel(string(v.Model.Identifier)),
			Comment:     v.Comment,
			Fields:      c.compileModelFields(v.Model.Elements),
			Collections: c.compileCollections(v.Model.Elements),
			Options:     c.compileModelOptions(v.Model.Elements),
		})
	}
	return
}

func (c *configSchemaCompiler) compileStructs() (out []*SchemaStruct) {
	for _, v := range c.ast.Types {
		if v.Struct == nil {
			continue
		}

		if v.Struct.Identifier.IsReserved() {
			err := errors.Errorf("firemodel/schema: can't name struct %s, %s is a reserved word.", v.Struct.Identifier, v.Struct.Identifier)
			panic(err)
		}

		out = append(out, &SchemaStruct{
			Name:    strcase.ToCamel(string(v.Struct.Identifier)),
			Comment: v.Comment,
			Fields:  c.compileStructFields(v.Struct.Elements),
		})
	}
	return
}

func (c *configSchemaCompiler) compileEnums() (out []*SchemaEnum) {
	for _, v := range c.ast.Types {
		if v.Enum == nil {
			continue
		}
		out = append(out, &SchemaEnum{
			Name:    strcase.ToCamel(string(v.Enum.Identifier)),
			Comment: v.Comment,
			Values:  c.enumValuesToConfig(v.Enum.Values),
		})
	}
	return
}

func (c *configSchemaCompiler) compileLanguageOptions() (out SchemaOptions) {
	out = SchemaOptions{}
	for _, v := range c.ast.Types {
		opt := v.Option
		if opt == nil {
			continue
		}
		if out[opt.Language] == nil {
			out[opt.Language] = map[string]string{}
		}
		if opt.Key.IsReserved() {
			err := errors.Errorf("firemodel/schema: can't use option key %s, %s is a reserved word.", opt.Key, opt.Key)
			panic(err)
		}
		out[opt.Language][string(opt.Key)] = opt.Value
	}
	return
}

func (c *configSchemaCompiler) enumValuesToConfig(values []*ast.ASTEnumValue) (out []*SchemaEnumValue) {
	for _, enumValue := range values {
		out = append(out, &SchemaEnumValue{
			Name:    strcase.ToSnake(enumValue.Name),
			Comment: enumValue.Comment,
		})
	}
	return
}

func (c *configSchemaCompiler) compileModelFields(elements []*ast.ASTModelElement) (out []*SchemaField) {
	for _, element := range elements {
		field := element.Field
		if field == nil {
			continue
		}
		if field.Type.Base.IsCollection() {
			continue // handled in compileCollections
		}

		out = append(out, &SchemaField{
			Name:    strcase.ToSnake(field.Name),
			Comment: field.Comment,
			Type:    c.compileFieldType(field.Type),
		})
	}
	return
}

func (c *configSchemaCompiler) compileStructFields(elements []*ast.ASTStructElement) (out []*SchemaField) {
	for _, element := range elements {
		field := element.Field
		if field == nil {
			continue
		}
		if field.Type.Base.IsCollection() {
			continue // handled in compileCollections
		}

		out = append(out, &SchemaField{
			Name:    strcase.ToSnake(field.Name),
			Comment: field.Comment,
			Type:    c.compileFieldType(field.Type),
		})
	}
	return
}

func (c *configSchemaCompiler) compileCollections(elements []*ast.ASTModelElement) (out []*SchemaNestedCollection) {
	for _, element := range elements {
		field := element.Field
		if field == nil {
			continue // element is not a Field
		}
		if !field.Type.Base.IsCollection() {
			continue // handled in compileFields
		}
		modelType, ok := c.assertModelType(field.Type.Generic)
		if !ok {
			err := errors.Errorf("invalid collection type: %s", field.Type)
			panic(err)
		}
		out = append(out, &SchemaNestedCollection{
			Name:    field.Name,
			Comment: field.Comment,
			Type:    modelType,
		})
	}
	return
}

func (c *configSchemaCompiler) compileFieldType(astFieldType *ast.ASTFieldType) SchemaFieldType {
	if c.enums == nil {
		panic("bug: enum types not yet registered")
	}
	if enum, ok := c.assertEnumType(astFieldType); ok {
		return &Enum{T: enum}
	}
	if _, ok := c.assertModelType(astFieldType); ok {
		err := errors.Errorf("firemodel/schema: can't use models as field types (got %s); please use reference, collection or switch model to struct instead", astFieldType)
		panic(err)
	}
	if structT, ok := c.assertStructType(astFieldType); ok {
		return &Struct{T: structT}
	}
	switch astFieldType.Base {
	case ast.Boolean:
		return &Boolean{}
	case ast.Integer:
		return &Integer{}
	case ast.Double:
		return &Double{}
	case ast.Timestamp:
		return &Timestamp{}
	case ast.String:
		return &String{}
	case ast.Bytes:
		return &Bytes{}
	case ast.GeoPoint:
		return &GeoPoint{}
	case ast.File:
		return &File{}
	case ast.URL:
		return &URL{}
	case ast.Map:
		if generic := astFieldType.Generic; generic != nil {
			return &Map{T: c.compileFieldType(generic)}
		}
		return &Map{}
	case ast.Array:
		if generic := astFieldType.Generic; generic != nil {
			return &Array{T: c.compileFieldType(generic)}
		}
		return &Array{}
	case ast.Reference:
		if astFieldType.Generic == nil {
			return &Reference{}
		} else if modelType, ok := c.assertModelType(astFieldType.Generic); ok {
			return &Reference{T: modelType}
		} else {
			err := errors.Errorf("firemodel: invalid generic type %s in %s<%s> (must be a model type)", astFieldType.Generic, astFieldType.Base, astFieldType.Generic)
			panic(err)
		}
	}

	err := errors.Errorf("invalid type: %s", astFieldType.Base)
	panic(err)
}

func (c *configSchemaCompiler) assertModelType(astFieldType *ast.ASTFieldType) (*SchemaModel, bool) {
	if c.models == nil {
		panic("bug: model types not yet registered")
	}
	if astFieldType == nil {
		return nil, false
	}
	astType := astFieldType.Base
	for _, model := range c.models {
		if model.Name == strcase.ToCamel(string(astType)) {
			return model, true
		}
	}
	return nil, false
}

func (c *configSchemaCompiler) assertStructType(astFieldType *ast.ASTFieldType) (*SchemaStruct, bool) {
	if c.structs == nil {
		panic("bug: model types not yet registered")
	}
	if astFieldType == nil {
		return nil, false
	}
	astType := astFieldType.Base
	for _, schemaStruct := range c.structs {
		if schemaStruct.Name == strcase.ToCamel(string(astType)) {
			return schemaStruct, true
		}
	}
	return nil, false
}

func (c *configSchemaCompiler) assertEnumType(astType *ast.ASTFieldType) (*SchemaEnum, bool) {
	if astType == nil {
		return nil, false
	}
	for _, enum := range c.enums {
		if enum.Name == strcase.ToCamel(string(astType.Base)) {
			if astType.Generic != nil {
				panic(errors.Errorf("generic enums are not supported: %s %v", astType.Base, astType.Generic))
			}
			return enum, true
		}

	}
	return nil, false
}

func (c *configSchemaCompiler) compileModelOptions(elements []*ast.ASTModelElement) SchemaModelOptions {
	out := SchemaModelOptions{}
	for _, element := range elements {
		option := element.Option
		if option == nil {
			continue
		}
		if out[option.Language] == nil {
			out[option.Language] = map[string]string{}
		}
		if option.Key.IsReserved() {
			err := errors.Errorf("firemodel/schema: can't use option key %s, %s is a reserved word.", option.Key, option.Key)
			panic(err)
		}
		out[option.Language][string(option.Key)] = option.Value
	}
	return out
}
