package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/visor-tax/firemodel"
)

var showCmd = &cobra.Command{
	Use:   "show-languages",
	Short: "Show all available languages.",
	Run: func(cmd *cobra.Command, args []string) {
		for _, language := range firemodel.AllModelers() {
			fmt.Println(language)
		}
	},
}
