package cmd

import (
	// Modeler registrations:
	_ "github.com/mickeyreiss/firemodel/langs/go"
	_ "github.com/mickeyreiss/firemodel/langs/ios"
	_ "github.com/mickeyreiss/firemodel/langs/ts"
)

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	req struct {
		schema string
	}
)

var rootCmd = &cobra.Command{
	Use:   "firemodel",
	Short: "Type-safe, cross-platform models for Firestore",
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
