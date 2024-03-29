package main

import (
	"api"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestHelloService_Hello(t *testing.T) {
	for _, tt := range []struct {
		headline string
		req      *api.HelloRequest
		err      error
	}{
		{"success", &api.HelloRequest{Name: "done"}, nil},
		{"unknown", &api.HelloRequest{}, ErrUnknown},
	} {
		t.Run(tt.headline, func(t *testing.T) {
			s := NewHelloService()
			_, err := s.Hello(context.Background(), tt.req)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
