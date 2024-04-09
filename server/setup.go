package main

import (
	"api"
	"data"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Server struct {
	fx.Out

	Server *grpc.Server
}

func NewServer() Server {
	s := grpc.NewServer()
	return Server{Server: s}
}

type BootParams struct {
	fx.In

	Server        *grpc.Server
	UserCommander data.UserCommander
}

func setupServer(param BootParams) {
	api.RegisterHelloServerServer(param.Server, NewHelloService())
	api.RegisterUserServerServer(param.Server, NewUserService(param.UserCommander))
}
