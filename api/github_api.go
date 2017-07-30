package api

import (
	"context"
	"fmt"
	"encoding/json"
	"bytes"

	"ashwitha/workspace/demo/utils"
)

type githubAPI struct {
	baseURL    string
	httpClient utils.Fetcher
}

type GithubAPI interface {
	GetRepo(ctx context.Context) ([]Repository, error)
}

func NewGithubAPI(baseURL string, client utils.Fetcher) GithubAPI {
	return &githubAPI{
		baseURL:    baseURL,
		httpClient: client,
	}
}

func (ga githubAPI) GetRepo(ctx context.Context) ([]Repository, error) {
	req := utils.Request{
		BaseURL: ga.baseURL,
		Route: "/repos",
		Method: "GET",
	}

	body, err := ga.httpClient.Fetch(req)
	if err != nil {
		fmt.Printf("Fetch: unable to get gitHub repos: %s", err)
		return nil, fmt.Errorf("Fetch: unable to get gitHub repos: %s", err)
	}

	repos := make([]Repository,0)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&repos)
	if err != nil {
		return nil, fmt.Errorf("getRepos: unable to unmarshal response: %s", err)
	}

	return repos, nil
}