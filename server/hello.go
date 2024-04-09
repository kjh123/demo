package main

import (
	"api"
	"context"
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrUnknown = errors.New("unknown")
)

type HelloService struct {
	api.UnimplementedHelloServerServer
}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (h *HelloService) Hello(_ context.Context, request *api.HelloRequest) (*api.HelloResponse, error) {
	if request.Name == "" {
		return nil, errors.Wrap(ErrUnknown, "form")
	}

	return &api.HelloResponse{
		Message: fmt.Sprintf("hello: %s", request.Name),
	}, nil
}
