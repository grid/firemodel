package cmd

import (
	// Modeler registrations:
	"fmt"

	"github.com/spf13/cobra"
	_ "github.com/visor-tax/firemodel/langs/go"
	_ "github.com/visor-tax/firemodel/langs/ios"
	_ "github.com/visor-tax/firemodel/langs/ts"
	"github.com/visor-tax/firemodel/version"
)

var (
	req struct {
		schemas []string
	}
)

var rootCmd = &cobra.Command{
	Use:     "firemodel",
	Short:   "Type-safe, cross-platform models for Firestore",
	Version: version.Version,
}

func init() {
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(compileCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
