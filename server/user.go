package main

import (
	"api"
	"context"
	"data"
	"github.com/pkg/errors"
	"time"
)

var ErrEmpty = errors.New("empty user")

type UserService struct {
	api.UnimplementedUserServerServer
	userCommander data.UserCommander
}

func NewUserService(userCommander data.UserCommander) *UserService {
	return &UserService{userCommander: userCommander}
}

func (s *UserService) Register(ctx context.Context, request *api.UserRequest) (*api.UserResponse, error) {
	if request.GetName() == "" {
		return nil, ErrEmpty
	}

	err := s.userCommander.Create(ctx, data.User{
		Name:      request.GetName(),
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &api.UserResponse{Message: "success"}, nil
}
