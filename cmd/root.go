package cmd

import (
	"fmt"
	"gitbrush/internal/github"
	"gitbrush/internal/painter"
	"gitbrush/internal/ui"
	"os"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "faux-commit-painter",
	Short: "Paint your Github contribution graph!",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("faux-commit-painter %s\n\n", ui.AppVersion)

		if v, err := github.LatestReleaseVersion(); err == nil && v != "" && v != ui.AppVersion {
			fmt.Printf("Update available: %s (you have %s)\n", v, ui.AppVersion)
			fmt.Println("  Reinstall with one command:")
			fmt.Println("  curl -fsSL https://raw.githubusercontent.com/recrease345/faux-commit-painter/main/install.sh | sh\n")
		}

		config, err := ui.AskConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		repoBar := progressbar.Default(-1, "Creating repository on Github...")
		repoUrl, err := github.CreateRepo(config.Username, config.Token, config.RepoName)
		if err != nil {
			repoBar.Clear()
			fmt.Printf("\n Failed to create repo: %v\n", err)
			os.Exit(1)
		}
		repoBar.Finish()

		commits, err := painter.PaintGraph(config, repoUrl)
		if err != nil {
			fmt.Printf("Failed to paint graph: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Done! Succesfule pushed %d commits\n", commits)
		fmt.Println("ℹ️ It may take GitHub up to a few hours to fully update your contribution graph after completion")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
