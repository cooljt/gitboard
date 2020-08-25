package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

type gitHubClient struct {
	client *github.Client
	ctx    context.Context
}

func (c gitHubClient) searchGithubRepoByStar(num int, searchWord string, language string) []*github.Repository {
	client := c.client
	ctx := c.ctx

	queryString := fmt.Sprintf("%s language:%s", searchWord, language)
	queryOptions := &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
	}

	repos, _, err := client.Search.Repositories(
		ctx,
		queryString,
		queryOptions,
	)

	if err != nil {
		return nil
	}

	total := repos.Total
	return repos.Repositories[:min(num, *total)]

}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
