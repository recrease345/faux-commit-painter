package painter

import (
	"fmt"
	"gitbrush/internal/ui"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/schollz/progressbar/v3"
)

func PaintGraph(cfg *ui.Config, repoURL string) (int, error) {
	tmpDir, err := os.MkdirTemp("", "commit_painter-*")
	if err != nil {
		return 0, err
	}

	defer os.RemoveAll(tmpDir)

	auth := &http.BasicAuth{
		Username: cfg.Username,
		Password: cfg.Token,
	}

	cloneBar := progressbar.Default(-1, "Cloning repository...")
	repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:  repoURL,
		Auth: auth,
	})
	if err != nil {
		cloneBar.Clear()
		return 0, fmt.Errorf("clone failed: %v", err)
	}
	cloneBar.Finish()

	w, err := repo.Worktree()
	if err != nil {
		return 0, err
	}

	startDate := time.Now().AddDate(0, 0, -364)
	offset := int(startDate.Weekday())
	startDate = startDate.AddDate(0, 0, -offset)

	var ascii []string
	if cfg.Mode == "text" {
		ascii = strings.Split(strings.TrimSpace(cfg.ASCII), "\n")
	}

	days := 364

	bar := progressbar.Default(int64(days), "Painting graph...")
	totalComitsCounter := 0

	for i := 0; i < days; i++ {
		currentData := startDate.AddDate(0, 0, i)

		shouldCommit := false
		commitsPerDay := 1

		switch cfg.Mode {
		case "full":
			shouldCommit = true
			commitsPerDay = 10
		case "chaotic":
			shouldCommit = rand.IntN(3) == 1
			commitsPerDay = rand.IntN(7) + 1
		case "text":
			week := i / 7
			day := i % 7

			if day < len(ascii) && week < len(ascii[day]) {
				char := string(ascii[day][week])

				if char != " " && char != "." && char != "_" {
					shouldCommit = true
					commitsPerDay = 15
				}
			}
		}

		if shouldCommit {
			for commit := 0; commit < commitsPerDay; commit++ {
				fileName := "faux-commit-painter-activity.txt"
				filePath := filepath.Join(tmpDir, fileName)

				content := fmt.Sprintf("Commit %d-%d\n", i, commit)
				err := os.WriteFile(filePath, []byte(content), 0644)
				if err != nil {
					return 0, err
				}

				_, err = w.Add(fileName)
				if err != nil {
					return 0, err
				}

				commitTime := currentData.Add(time.Duration(commit) * time.Hour)

				_, err = w.Commit(fmt.Sprintf("faux-commit-painter %d-%d", i, commit), &git.CommitOptions{
					Author: &object.Signature{
						Name:  cfg.Username,
						Email: cfg.Email,
						When:  commitTime,
					},
				})
				if err != nil {
					return 0, fmt.Errorf("failed to commit: %w", err)
				}

				totalComitsCounter++
			}
		}

		bar.Add(1)
	}

	pushBar := progressbar.Default(-1, "Pushing commit to GitHub...")
	err = repo.Push(&git.PushOptions{
		Auth: auth,
	})
	if err != nil {
		pushBar.Clear()
		return 0, fmt.Errorf("failed to push: %w", err)
	}
	pushBar.Finish()

	return totalComitsCounter, nil
}
