package cmd

import (
	"gitbrush/internal/ui"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "faux-commit-painter",
	Short: "Paint your Github contribution graph!",
	Run:  func(cmd *cobra.Command, args []string) {
		config, err := ui.AskConfig()
		



	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
