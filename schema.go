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
	models []*SchemaModel
	enums  []*SchemaEnum

	ast *ast.AST
}

func (c *configSchemaCompiler) compileConfig() (*Schema, error) {
	if err := c.precompileEnumTypes(); err != nil {
		return nil, err
	}
	if err := c.precompileModelTypes(); err != nil {
		return nil, err
	}

	return &Schema{
		Models:  c.compileModels(),
		Enums:   c.compileEnums(),
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
			Fields:      c.compileFields(v.Model.Elements),
			Collections: c.compileCollections(v.Model.Elements),
			Options:     c.compileModelOptions(v.Model.Elements),
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

func (c *configSchemaCompiler) compileFields(elements []*ast.ASTModelElement) (out []*SchemaField) {
	for _, element := range elements {
		field := element.Field
		if field == nil {
			continue // element is not a Field
		}
		if field.Type.Base.IsCollection() {
			continue // handled in compileCollections
		}

		out = append(out, &SchemaField{
			Name:    strcase.ToSnake(field.Name),
			Comment: field.Comment,
			Type:    c.compileFieldType(field.Type),
			Extras:  c.compileExtras(field.Type),
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
	if astFieldType.Base.IsPrimitive() {
		return SchemaFieldType(astFieldType.Base)
	}
	if c.enums == nil {
		panic("bug: enum types not yet registered")
	}
	if astFieldType.Base == ast.File {
		return Map
	}
	if astFieldType.Base == ast.URL {
		return String
	}
	_, ok := c.assertEnumType(astFieldType)
	if ok {
		return String
	}
	_, ok = c.assertModelType(astFieldType.Base)
	if ok {
		return Map
	}
	err := errors.Errorf("invalid type: %s", astFieldType.Base)
	panic(err)
}

func (c *configSchemaCompiler) compileModelType(astType *ast.ASTFieldType) string {
	if astType.Generic != "" {
		err := errors.Errorf("models cannot have generics: %s<%s>", astType.Base, astType.Generic)
		panic(err)
	}

	if modelType, ok := c.assertModelType(astType.Base); ok {
		return modelType
	}

	err := errors.Errorf("invalid type: %s", astType.Base)
	panic(err)
}

func (c *configSchemaCompiler) assertModelType(astType ast.ASTType) (string, bool) {
	if c.models == nil {
		panic("bug: model types not yet registered")
	}
	for _, model := range c.models {
		if model.Name == strcase.ToCamel(string(astType)) {
			return model.Name, true
		}
	}
	return "", false
}

func (c *configSchemaCompiler) assertEnumType(astType *ast.ASTFieldType) (ast.ASTType, bool) {
	for _, enum := range c.enums {
		if enum.Name == strcase.ToCamel(string(astType.Base)) {
			return astType.Base, true
		}

	}
	return "", false
}

func (c *configSchemaCompiler) compileExtras(astType *ast.ASTFieldType) *SchemaFieldExtras {
	out := &SchemaFieldExtras{}

	if enumType, ok := c.assertEnumType(astType); ok {
		out.EnumType = string(enumType)
	}
	if modelType, ok := c.assertModelType(astType.Base); ok {
		out.MapTo = string(modelType)
	}

	switch astType.Base {
	case ast.URL:
		out.URL = true
	case ast.File:
		out.File = true
	}

	if astType.Generic != "" {
		switch astType.Base {
		case ast.Map:
			if astType.Generic.IsPrimitive() {
				out.MapToPrimitive = SchemaFieldType(astType.Generic)
			} else {
				out.MapTo = string(astType.Generic)
			}
		case ast.Array:
			if astType.Generic.IsPrimitive() {
				out.ArrayOfPrimitive = SchemaFieldType(astType.Generic)
			} else {
				out.ArrayOf = string(astType.Generic)
			}
		case ast.Reference:
			if astType.Generic.IsPrimitive() {
				err := errors.Errorf("firemodel: invalid generic type %s in %s<%s> (must be a model type)", astType.Generic, astType.Base, astType.Generic)
				panic(err)
			} else {
				out.ReferenceTo = string(astType.Generic)
			}
		default:
			err := errors.Errorf("firemodel: invalid generic type on %s", astType.Base)
			panic(err)
		}
	}
	return out
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
