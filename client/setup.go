package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

type MountController interface {
	Mount(router gin.IRouter)
}

type RouteParams struct {
	fx.In

	Controllers []MountController `group:"controller"`
}

func setupRouter(params RouteParams) *gin.Engine {
	r := gin.Default()
	anonymous := r.Group("_")
	for _, controller := range params.Controllers {
		controller.Mount(anonymous)
	}
	return r
}

type serverParams struct {
	fx.In

	Gin   *gin.Engine
	Addr  string `name:"addr"`
	MySQL string `name:"mysql_dsn"`
}

func setupServer(params serverParams) *Server {
	server := &http.Server{Addr: params.Addr, Handler: params.Gin.Handler()}
	return &Server{
		server: server,
	}
}

type Server struct {
	server *http.Server
}

func (s *Server) Start(_ context.Context) error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
