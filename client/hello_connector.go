package main

import (
	"api"
	"context"
	"time"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HelloConnector interface {
	CallServer(context.Context, *api.HelloRequest) (*api.HelloResponse, error)
}

type HelloConnect struct {
	ServerDomain string
}

type ServeConnectParams struct {
	fx.In

	ServerDomain string `name:"server_domain"`
}

func NewHelloConnect(params ServeConnectParams) *HelloConnect {
	return &HelloConnect{ServerDomain: params.ServerDomain}
}

func (h *HelloConnect) CallServer(ctx context.Context, request *api.HelloRequest) (resp *api.HelloResponse, err error) {
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(h.ServerDomain, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer func() {
		err = conn.Close()
	}()
	if err != nil {
		return nil, err
	}

	c := api.NewHelloServerClient(conn)

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	r, err := c.Hello(cancelCtx, request)
	return r, err
}
