package mocks

import (
	"github.com/stretchr/testify/mock"

	"ashwitha/workspace/demo/utils"
)

// MockFetcher implements a mock Fetcher.
type MockFetcher struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: req.
func (m *MockFetcher) Fetch(req utils.Request) (responseBody []byte, err error) {
	args := m.Called(req)

	if args.Get(0) != nil {
		responseBody = args.Get(0).([]byte)
	}

	err = args.Error(1)

	return responseBody, err
}

