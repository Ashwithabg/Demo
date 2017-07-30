package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"ashwitha/workspace/demo/api"
)

type GitHubAPI struct {
	mock.Mock
}

func (ga GitHubAPI) GetRepo(ctx context.Context) (repo []api.Repository, err error) {
	args := ga.Called(ctx)
	if args.Get(0) != nil {
		repo = args.Get(0).([]api.Repository)
	}
	err = args.Error(1)
	return repo, err
}
