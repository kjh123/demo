package main

import (
	"api"
	"context"
	"data"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockUserCommander struct {
	mock.Mock
}

func (m *mockUserCommander) Create(ctx context.Context, user data.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestUserService_Register(t *testing.T) {
	for _, tt := range []struct {
		headline string
		req      *api.UserRequest
		err      error
	}{
		{"success", &api.UserRequest{Name: "lt"}, nil},
		{"empty error", &api.UserRequest{}, ErrEmpty},
		{"internal error", &api.UserRequest{Name: "n"}, errors.New("internal error")},
	} {
		t.Run(tt.headline, func(t *testing.T) {
			userCommander := &mockUserCommander{}
			service := NewUserService(userCommander)

			userCommander.On("Create", mock.Anything, mock.AnythingOfType("data.User")).Return(tt.err)
			_, err := service.Register(context.Background(), tt.req)
			assert.ErrorIs(t, err, err)
		})
	}
}
