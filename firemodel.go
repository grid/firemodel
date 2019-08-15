package firemodel

import (
	"io"
	"strings"

	"github.com/go-errors/errors"
)

type SourceCoder interface {
	NewFile(filename string) (io.WriteCloser, error)
	Flush() error
}

type GenOptions map[string]string

func (options GenOptions) Get(s string) string {
	return options[s]
}

type Modeler interface {
	Model(schema *Schema, options GenOptions, sourceCoder SourceCoder) error
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

			output := language.Output
			options := map[string]string{}
			if strings.Contains(language.Output, ":") {
				n := strings.SplitN(language.Output, ":", 1)
				output = n[1]
				optionPairs := strings.Split(n[0], "=")
				if len(optionPairs)%2 != 0 {
					panic(errors.Errorf("Invalid output options: %s", language.Output))
				}
				for idx := 0; idx < len(optionPairs); idx += 2 {
					options[optionPairs[0]] = optionPairs[1]
				}
			}

			sourceCoder := config.SourceCoderProvider(output)
			if err := modeler.Model(schema, options, sourceCoder); err != nil {
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
