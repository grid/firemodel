package firemodel

import (
	"github.com/mickeyreiss/firemodel/internal/ast"
	"io"
	"github.com/go-errors/errors"
)

const (
	Boolean   = SchemaFieldType(ast.Boolean)
	Integer   = SchemaFieldType(ast.Integer)
	Double    = SchemaFieldType(ast.Double)
	Timestamp = SchemaFieldType(ast.Timestamp)
	String    = SchemaFieldType(ast.String)
	Bytes     = SchemaFieldType(ast.Bytes)
	Reference = SchemaFieldType(ast.Reference)
	GeoPoint  = SchemaFieldType(ast.GeoPoint)
	Array     = SchemaFieldType(ast.Array)
	Map       = SchemaFieldType(ast.Map)
)

type SourceCoder interface {
	NewFile(filename string) (io.WriteCloser, error)
	Flush() error
}

type Modeler interface {
	Model(schema *Schema, sourceCoder SourceCoder) error
}

func Run(
	schema *Schema,
	config *Config,
) error {
	if config == nil {
		err := errors.New("firemodel: config not set")
		return err
	}

	for _, language := range config.Languages {
		if err := func(language *Language) error {
			modeler := language.Modeler()
			sourceCoder := config.SourceCoderProvider(language.Output)
			if err := modeler.Model(schema, sourceCoder); err != nil {
				return err
			}
			if err := sourceCoder.Flush(); err != nil {
				return err
			}
			return nil
		}(&language); err != nil {
			return err
		}
	}
	return nil
}
