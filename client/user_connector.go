package main

import (
	"api"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type UserConnector interface {
	CallServeUser(ctx context.Context, request *api.UserRequest) (*api.UserResponse, error)
}

type UserConnect struct {
	ServerDomain string
}

func NewUserConnect(params ServeConnectParams) *UserConnect {
	return &UserConnect{ServerDomain: params.ServerDomain}
}

func (u *UserConnect) CallServeUser(ctx context.Context, request *api.UserRequest) (*api.UserResponse, error) {
	conn, err := grpc.Dial(u.ServerDomain, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	c := api.NewUserServerClient(conn)

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	r, err := c.Register(cancelCtx, request)
	return r, err
}
