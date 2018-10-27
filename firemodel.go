package firemodel

import (
	"io"

	"github.com/go-errors/errors"
)

type SourceCoder interface {
	NewFile(filename string) (io.WriteCloser, error)
	Flush() error
}

type Modeler interface {
	Model(schema *Schema, sourceCoder SourceCoder) error
}

type Client interface {
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
