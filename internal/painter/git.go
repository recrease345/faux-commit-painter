package painter

import (
	"errors"
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

	now := time.Now().UTC()
	startDate := now.AddDate(0, 0, -364)

	offset := int(startDate.Weekday())
	if offset != 0 {
		startDate = startDate.AddDate(0, 0, -offset)
	}
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 12, 0, 0, 0, time.UTC)

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	var ascii []string
	if cfg.Mode == "text" {
		ascii = strings.Split(strings.TrimSpace(cfg.ASCII), "\n")
	}

	days := 371 // 53 weeks

	bar := progressbar.Default(int64(days), "Painting graph...")
	totalCommitsCounter := 0
	pushed := 0

	syncToGitHub := func(label string) error {
		if totalCommitsCounter == pushed {
			return nil
		}
		chunkBar := progressbar.Default(-1, label)
		err := repo.Push(&git.PushOptions{
			Auth: auth,
		})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			chunkBar.Clear()
			return fmt.Errorf("failed to push chunk: %w", err)
		}
		chunkBar.Finish()
		pushed = totalCommitsCounter
		return nil
	}

	for i := 0; i < days; i++ {
		currentDate := startDate.AddDate(0, 0, i)

		if currentDate.Format("2006-01-02") > today.Format("2006-01-02") {
			break
		}

		shouldCommit := false
		commitsPerDay := 1

		switch cfg.Mode {
		case "full":
			shouldCommit = true
			commitsPerDay = 5
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

				commitTime := currentDate.Add(time.Duration(commit) * time.Minute)

				_, err = w.Commit(fmt.Sprintf("faux-commit-painter %d-%d", i, commit), &git.CommitOptions{
					Author: &object.Signature{
						Name:  cfg.Username,
						Email: cfg.Email,
						When:  commitTime,
					},
					Committer: &object.Signature{
						Name:  cfg.Username,
						Email: cfg.Email,
						When:  commitTime,
					},
				})
				if err != nil {
					return 0, fmt.Errorf("failed to commit: %w", err)
				}

				totalCommitsCounter++
			}
		}

		if i > 0 && i%28 == 0 {
			if err := syncToGitHub(fmt.Sprintf("Syncing to GitHub (day %d/%d)...", i, days)); err != nil {
				return 0, err
			}
		}

		bar.Add(1)
	}

	if err := syncToGitHub("Final sync to GitHub..."); err != nil {
		return 0, err
	}
	bar.Finish()

	return totalCommitsCounter, nil
}
