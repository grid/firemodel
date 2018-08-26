package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/mickeyreiss/firemodel"
)

var showCmd = &cobra.Command{
	Use:   "show-languages",
	Short: "Show all available languages for --out",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available languages:")
		for _, language := range firemodel.AllModelers() {
			fmt.Println("-", language)
		}
	},
}
