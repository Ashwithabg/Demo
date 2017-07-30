package api

import (
	"testing"
	"context"
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"

	"ashwitha/workspace/demo/mocks/utils"
	"ashwitha/workspace/demo/utils"
)

var testResponse = `
[
    {
        "name": "Assignments",
        "description": null,
        "language": "Java"
    },
    {
        "name": "Demo",
        "description": null,
        "language": null
    }
]`

func TestGetRepo(t *testing.T) {
	ctx := context.Background()

	req := utils.Request{
		BaseURL:"/github/api/repo",
		Route: "/repos",
		Method: "GET",
	}
	fetcher := &mocks.MockFetcher{}
	fetcher.On("Fetch", req).Return([]byte(testResponse), nil)

	githubAPI := NewGithubAPI("/github/api/repo", fetcher)
	repos, err := githubAPI.GetRepo(ctx)

	expectedResp := []Repository{
		{
			Name: "Assignments",
			Description: "",
			Language: "Java",
		},
		{
			Name: "Demo",
			Description: "",
			Language: "",
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedResp, repos)
	mock.AssertExpectationsForObjects(t, fetcher)
}

func TestGetRepoFails(t *testing.T) {
	ctx := context.Background()

	req := utils.Request{
		BaseURL:"/github/api/repo",
		Route: "/repos",
		Method: "GET",
	}
	fetcher := &mocks.MockFetcher{}
	fetcher.On("Fetch", req).Return(nil, errors.New("connection lost"))

	githubAPI := NewGithubAPI("/github/api/repo", fetcher)
	_, err := githubAPI.GetRepo(ctx)

	assert.EqualError(t, err, "Fetch: unable to get gitHub repos: connection lost")
	mock.AssertExpectationsForObjects(t, fetcher)
}

func TestGetRepoInvalidResponse(t *testing.T) {
	ctx := context.Background()

	req := utils.Request{
		BaseURL:"/github/api/repo",
		Route: "/repos",
		Method: "GET",
	}
	invalidResp := " Invalid Response"
	fetcher := &mocks.MockFetcher{}
	fetcher.On("Fetch", req).Return([]byte(invalidResp), nil)

	githubAPI := NewGithubAPI("/github/api/repo", fetcher)
	_, err := githubAPI.GetRepo(ctx)

	assert.Contains(t, err.Error(), "getRepos: unable to unmarshal response: ")
	mock.AssertExpectationsForObjects(t, fetcher)
}
