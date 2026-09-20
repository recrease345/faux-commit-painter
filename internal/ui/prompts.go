package ui

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/common-nighthawk/go-figure"
)

type Config struct {
	Username string
	Token    string
	Email    string
	RepoName string
	Mode     string
	ASCII    string
}

func AskConfig() (*Config, error) {
	cfg := &Config{}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your Github username").Placeholder("octocat").Value(&cfg.Username),

			huh.NewInput().
				Title("Enter your Github email").Placeholder("octocat@example.com").Value(&cfg.Email),

			huh.NewInput().
				Title("Enter new repository name(it will be private)").Placeholder("my-paint-repo").Value(&cfg.RepoName),

			huh.NewNote().
				Title("How to get your GitHub Personal Access Token").
				Description(`1. Go to: GitHub -> Settings -> Developer settings
2. Select: personal access tokens -> tokens (classic)
3. Click: generate new token (classic)
4. Scopes: check the box next to "repo" (control of private repositories)
5. Click: generate token
6. Copy your token (it starts "ghp_")

It is needed to automatically create 
a repository for commits. Your tokens are not transferred or stored!`),

			huh.NewInput().
				Title("Paste your Github token(input is hidden)").Placeholder("ghp_...").EchoMode(huh.EchoModePassword).Value(&cfg.Token),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose painting mode").
				Options(
					huh.NewOption("Text to ASCII(type text, we draw it)", "text"),
					huh.NewOption("Chaotic (random commits)", "chaotic"),
					huh.NewOption("Full (all green)", "full"),
					// huh.NewOption("Custom ASCII Art(soon...)", "ascii"),
				).Value(&cfg.Mode),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	if cfg.Mode == "text" {
		var userInput string
		textForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Enter text to draw(hi, go, 2026...)").Placeholder("2026").Value(&userInput),
			),
		)

		if err := textForm.Run(); err != nil {
			return nil, err
		}

		fig := figure.NewFigure(strings.ToUpper(userInput), "banner3", true)
		asciiStr := fig.String()

		lines := strings.Split(strings.TrimRight(asciiStr, "\n"), "\n")
		var builder strings.Builder

		for i, line := range lines {
			for _, ch := range line {
				if ch == ' ' {
					builder.WriteString(".")
				} else {
					builder.WriteString(string(ch))
				}
			}

			if i < len(lines)-1 {
				builder.WriteString("\n")
			}
		}

		cfg.ASCII = builder.String()

		previewForm := huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("Preview your art").
					Description(cfg.ASCII + "\n\npress Enter to start..."),
			),
		)

		previewForm.Run()
	}

	/*if cfg.Mode == "ascii" {
		asciiForm := huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("ASCII Generator").
			)
		)
	}*/

	return cfg, nil
}
