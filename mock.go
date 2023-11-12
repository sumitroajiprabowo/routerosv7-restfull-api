// mocks.go

package routerosv7_restfull_api

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockRouterOSAPI struct {
	mock.Mock
}

func (m *MockRouterOSAPI) Print(ctx context.Context, host, username, password, command string) (interface{}, error) {
	args := m.Called(ctx, host, username, password, command)
	return args.Get(0), args.Error(1)
}

func (m *MockRouterOSAPI) Add(
	ctx context.Context, host, username, password, command string, payload []byte,
) (interface{}, error) {
	args := m.Called(ctx, host, username, password, command, payload)
	return args.Get(0), args.Error(1)
}

func (m *MockRouterOSAPI) Set(
	ctx context.Context, host, username, password, command string, payload []byte,
) (interface{}, error) {
	args := m.Called(ctx, host, username, password, command, payload)
	return args.Get(0), args.Error(1)
}

func (m *MockRouterOSAPI) Remove(ctx context.Context, host, username, password, command string) (interface{}, error) {
	args := m.Called(ctx, host, username, password, command)
	return args.Get(0), args.Error(1)
}

func (m *MockRouterOSAPI) Run(
	ctx context.Context, host, username, password, command string, payload []byte,
) (interface{}, error) {
	args := m.Called(ctx, host, username, password, command, payload)
	return args.Get(0), args.Error(1)
}
