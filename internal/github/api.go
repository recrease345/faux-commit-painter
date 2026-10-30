package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v60/github"
)

func CreateRepo(username, token, repoName string) (string, error) {
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	_, resp, err := client.Repositories.Get(ctx, username, repoName)
	if err == nil && resp != nil && resp.StatusCode == 200 {
		return fmt.Sprintf("https://github.com/%s/%s.git", username, repoName), nil
	}

	repo := &github.Repository{
		Name:     github.String(repoName),
		Private:  github.Bool(false),
		AutoInit: github.Bool(true),
	}

	_, _, err = client.Repositories.Create(ctx, "", repo)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://github.com/%s/%s.git", username, repoName), nil
}

func LatestReleaseVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := github.NewClient(nil)
	rel, resp, err := client.Repositories.GetLatestRelease(ctx, "recrease345", "faux-commit-painter")
	if err != nil {
		if resp != nil && resp.StatusCode == 403 {
			return "", nil
		}
		return "", err
	}
	return rel.GetTagName(), nil
}
