package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v60/github"
)

func CreateRepo(token, repoName, username string) (string, error) {
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	_, resp, err := client.Repositories.Get(ctx, username, repoName)
	if err == nil && resp.StatusCode == 200 {
		return fmt.Sprintf("https://github.com/%s/%s.git", username, repoName), nil
	}

	repo := &github.Repository{
		Name: github.String(repoName),
		Private: github.Bool(true),
		AutoInit: github.Bool(true),
	}

	_, _, err = client.Repositories.Create(ctx, "", repo)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://github.com/%s/%s.git")
}
