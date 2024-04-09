package main

import (
	"api"
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var _ HelloConnector = (*mockHelloConnector)(nil)

type mockHelloConnector struct {
	mock.Mock
}

func (m *mockHelloConnector) CallServer(ctx context.Context, request *api.HelloRequest) (*api.HelloResponse, error) {
	args := m.Mock.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*api.HelloResponse), nil
}

type mockWriter struct {
	gin.ResponseWriter
	header http.Header
	status int
}

func (i *mockWriter) WriteHeader(code int) {
	i.status = code
}

func (i *mockWriter) Header() http.Header {
	return i.header
}

func (i *mockWriter) Write(_ []byte) (int, error) {
	return 0, nil
}

func (i *mockWriter) Status() int {
	return i.status
}

func TestHelloService_Hello(t *testing.T) {
	for _, item := range []struct {
		headline     string
		connectorErr error
		statusCode   int
	}{
		{"success", nil, 200},
		{"connector error", errors.New("connector"), 500},
	} {
		t.Run(item.headline, func(t *testing.T) {
			connector := &mockHelloConnector{}
			service := NewHelloService(connector)
			if item.connectorErr == nil {
				connector.On("CallServer", mock.Anything, mock.AnythingOfType("*api.HelloRequest")).Return(&api.HelloResponse{Message: "xxx"}, nil)
			} else {
				connector.On("CallServer", mock.Anything, mock.AnythingOfType("*api.HelloRequest")).Return(nil, item.connectorErr)
			}

			ctx := &gin.Context{
				Writer:  &mockWriter{header: make(http.Header)},
				Request: &http.Request{URL: &url.URL{RawQuery: "name=t&to=t1&message=m"}},
			}
			service.Hello(ctx)
			assert.NoError(t, ctx.Err())
			assert.Equal(t, item.statusCode, ctx.Writer.Status())
		})
	}
}
