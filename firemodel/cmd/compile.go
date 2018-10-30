package cmd

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/visor-tax/firemodel"
	"github.com/visor-tax/firemodel/internal/tempwriter"
	"github.com/spf13/cobra"
	"os"
)

var compileReq struct {
	wipe        bool
	langOutDirs map[string]*string
}

func init() {
	compileCmd.PersistentFlags().StringVar(&req.schema, "schema", "schema.firemodel", "Path to firemodel schema.")
	compileCmd.PersistentFlags().BoolVarP(&compileReq.wipe, "wipe", "f", false, "Confirms it is ok to rm -rf the output directories. (This is generally something you want, but defaults off for safety.)")

	compileReq.langOutDirs = make(map[string]*string)
	for _, modeler := range firemodel.AllModelers() {
		compileReq.langOutDirs[modeler] = new(string)
		compileCmd.PersistentFlags().StringVar(compileReq.langOutDirs[modeler], modeler+"_out", "", fmt.Sprintf("%s output directory", modeler))
	}
}

var compileCmd = &cobra.Command{
	Use:     "compile",
	Aliases: []string{"c"},
	Short:   "Type-safe, cross-platform models for Firestore",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		for k := range compileReq.langOutDirs {
			if !cmd.Flag(k + "_out").Changed {
				delete(compileReq.langOutDirs, k)
			}
		}
		if len(compileReq.langOutDirs) == 0 {
			return errors.New("no languages requested")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		r, err := os.Open(req.schema)
		if err != nil {
			panic(err)
		}

		schema, err := firemodel.ParseSchema(r)
		if err != nil {
			panic(err)
		}

		config := &firemodel.Config{
			SourceCoderProvider: func(prefix string) firemodel.SourceCoder {
				return tempwriter.New(prefix, compileReq.wipe)
			},
		}

		for language, outDir := range compileReq.langOutDirs {
			if outDir == nil {
				continue
			}
			config.Languages = append(config.Languages, firemodel.Language{
				Language: language,
				Output:   *outDir,
			})
		}

		if err := firemodel.Run(schema, config); err != nil {
			panic(err)
		}

	},
}
