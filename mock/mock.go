// mocks.go

package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type RouterOSAPI struct {
	mock.Mock
}

func (m *RouterOSAPI) Auth(ctx context.Context, host, username, password string) (interface{}, error) {
	args := m.Called(ctx, host, username, password)
	return args.Get(0), args.Error(1)
}

func (m *RouterOSAPI) Print(ctx context.Context, host, username, password, command string) (interface{}, error) {
	args := m.Called(ctx, host, username, password, command)
	return args.Get(0), args.Error(1)
}

func (m *RouterOSAPI) Add(
	ctx context.Context, host, username, password, command string, payload []byte,
) (interface{}, error) {
	args := m.Called(ctx, host, username, password, command, payload)
	return args.Get(0), args.Error(1)
}

func (m *RouterOSAPI) Set(
	ctx context.Context, host, username, password, command string, payload []byte,
) (interface{}, error) {
	args := m.Called(ctx, host, username, password, command, payload)
	return args.Get(0), args.Error(1)
}

func (m *RouterOSAPI) Remove(ctx context.Context, host, username, password, command string) (interface{}, error) {
	args := m.Called(ctx, host, username, password, command)
	return args.Get(0), args.Error(1)
}

func (m *RouterOSAPI) Run(
	ctx context.Context, host, username, password, command string, payload []byte,
) (interface{}, error) {
	args := m.Called(ctx, host, username, password, command, payload)
	return args.Get(0), args.Error(1)
}
