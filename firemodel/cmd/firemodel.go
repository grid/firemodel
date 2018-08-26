package cmd

import (
	"os"

	_ "github.com/mickeyreiss/firemodel/langs/go"
	_ "github.com/mickeyreiss/firemodel/langs/ios"
	_ "github.com/mickeyreiss/firemodel/langs/ts"
	"path"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/go-errors/errors"
	"github.com/mickeyreiss/firemodel/internal/tempwriter"
	"github.com/mickeyreiss/firemodel"
)

var (
	req struct {
		schema string
		outDir string
		wipe   bool
	}
)

var rootCmd = &cobra.Command{
	Use:   "firemodel",
	Short: "Type-safe, cross-platform models for Firestore",
}

var compileCmd = &cobra.Command{
	Use:     "compile",
	Aliases: []string{"c"},
	Short:   "Type-safe, cross-platform models for Firestore",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.Errorf("Please specify at least one language. %s %v.", rootCmd.Use, firemodel.AllModelers())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
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
				return tempwriter.New(prefix, req.wipe)
			},
		}
		for _, language := range args {
			config.Languages = append(config.Languages, firemodel.Language{
				Language: language,
				Output:   path.Join(req.outDir, language),
			})
		}

		if err := firemodel.Run(schema, config); err != nil {
			panic(err)
		}

	},
}

var cleanCmd = &cobra.Command{
	Use: "clean",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !req.wipe {
			return errors.New("Refusing to proceed without --wipe or -f.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if req.wipe {
			if err := os.RemoveAll(req.outDir); err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	compileCmd.PersistentFlags().StringVar(&req.schema, "schema", "schema.firemodel", "Path to firemodel schema.")
	rootCmd.PersistentFlags().BoolVarP(&req.wipe, "wipe", "f", false, "Whether it is ok to unlink the output directory.")
	compileCmd.PersistentFlags().StringVarP(&req.outDir, "output-dir", "o", "./.build/firemodel", "Directory path for firemodel output.")
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(compileCmd)
	rootCmd.AddCommand(cleanCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
