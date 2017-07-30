package handlers

import (
	"testing"
	"context"
	"net/http/httptest"
	"bytes"
	"io/ioutil"
	"errors"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ashwitha/workspace/demo/api"
	"ashwitha/workspace/demo/mocks/api"
)

func TestGetRepo(t *testing.T) {
	ctx := context.Background()
	repo := []api.Repository{
		{
			Name: "HotelReservation",
			Description:"reserve hotel",
			Language:"java",
		},
	}

	githubAPI := &mocks.GitHubAPI{}
	githubAPI.On("GetRepo", ctx).Return(repo, nil)

	gh := GetReposHandler(githubAPI)

	r := httptest.NewRequest("POST", "/repo", bytes.NewBufferString(""))
	w := httptest.NewRecorder()

	gh(w, r)
	body, err := ioutil.ReadAll(w.Body)
	require.Nil(t, err)

	expectedRes := "[{\"name\":\"HotelReservation\",\"description\":\"reserve hotel\",\"language\":\"java\"}]\n"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t,  expectedRes, string(body))

	mock.AssertExpectationsForObjects(t, githubAPI)
}

func TestGetRepoFails(t *testing.T) {
	ctx := context.Background()

	githubAPI := &mocks.GitHubAPI{}
	githubAPI.On("GetRepo", ctx).Return(nil, errors.New("invalid url"))

	gh := GetReposHandler(githubAPI)

	r := httptest.NewRequest("POST", "/repo", bytes.NewBufferString(""))
	w := httptest.NewRecorder()

	gh(w, r)
	body, err := ioutil.ReadAll(w.Body)
	require.Nil(t, err)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "Failed to get github repositories\n", string(body))

	mock.AssertExpectationsForObjects(t, githubAPI)
}
