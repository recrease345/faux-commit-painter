package painter

import (
	"fmt"
	"gitbrush/internal/ui"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
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

	for i := 0; i < days; i++ {
		currentData := startDate.AddDate(0, 0, i)
	}
}
