package main

import (
	"api"
	"context"
	"data"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"testing"
)

type mockUserConnector struct {
	mock.Mock
}

func (m *mockUserConnector) CallServeUser(ctx context.Context, request *api.UserRequest) (*api.UserResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*api.UserResponse), nil
}

type mockUserQueryer struct {
	mock.Mock
}

func (m *mockUserQueryer) Find(ctx context.Context, id int) (*data.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*data.User), nil
}

func (m *mockUserQueryer) Select(ctx context.Context) ([]data.User, error) {
	args := m.Called(ctx)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	if args.Get(0) == nil {
		return nil, nil
	}

	return args.Get(0).([]data.User), nil
}

func TestUserService_Create(t *testing.T) {
	for _, tt := range []struct {
		headline   string
		body       url.Values
		response   *api.UserResponse
		err        error
		statusCode int
	}{
		{"success", url.Values{"name": {"t"}}, &api.UserResponse{Message: "success"}, nil, http.StatusCreated},
		{"bad request error", url.Values{}, &api.UserResponse{Message: "success"}, nil, http.StatusBadRequest},
		{"internal error", url.Values{"name": {"t"}}, nil, errors.New("internal"), http.StatusInternalServerError},
	} {
		t.Run(tt.headline, func(t *testing.T) {
			userConnector := &mockUserConnector{}
			userQueryer := &mockUserQueryer{}
			service := NewUserService(userConnector, userQueryer)

			ctx := &gin.Context{
				Writer:  &mockWriter{header: make(http.Header)},
				Request: &http.Request{Form: tt.body},
			}

			userConnector.On("CallServeUser", mock.Anything, mock.Anything).Return(tt.response, tt.err)

			service.Create(ctx)
			assert.NoError(t, ctx.Err())
			assert.Equal(t, tt.statusCode, ctx.Writer.Status())
		})
	}
}

func TestUserService_List(t *testing.T) {
	for _, tt := range []struct {
		headline   string
		err        error
		statusCode int
		users      []data.User
	}{
		{"success", nil, 200, []data.User{{ID: 1, Name: "t1"}, {ID: 2, Name: "t2"}}},
		{"internal error", errors.New("internal error"), 500, nil},
	} {
		t.Run(tt.headline, func(t *testing.T) {
			userConnector := &mockUserConnector{}
			userQueryer := &mockUserQueryer{}
			service := NewUserService(userConnector, userQueryer)

			ctx := &gin.Context{
				Writer:  &mockWriter{header: make(http.Header)},
				Request: &http.Request{},
			}
			userQueryer.On("Select", mock.Anything).Return(tt.users, tt.err)
			service.List(ctx)
			assert.NoError(t, ctx.Err())
			assert.Equal(t, tt.statusCode, ctx.Writer.Status())
		})
	}
}
