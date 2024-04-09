package main

import (
	"api"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (u *UserConnect) CallServeUser(ctx context.Context, request *api.UserRequest) (response *api.UserResponse, err error) {
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(u.ServerDomain, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer func() {
		err = conn.Close()
	}()

	if err != nil {
		return nil, err
	}

	c := api.NewUserServerClient(conn)

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	response, err = c.Register(cancelCtx, request)
	return
}
