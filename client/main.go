package main

import (
	"context"
	"data"
	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"net/http"
)

var ClientConf struct {
	ServerDomain string `default:"127.0.0.1:8000"`

	Addr  string
	Mysql string
}

func main() {
	kong.Parse(&ClientConf)
	app := fx.New(
		data.Module,
		fx.Provide(
			setupRouter,
			setupServer,
		),

		fx.Provide(
			fx.Annotate(NewHelloConnect, fx.As(new(HelloConnector))),
			fx.Annotate(NewHelloService, fx.As(new(MountController)), fx.ResultTags(`group:"controller"`)),

			fx.Annotate(NewUserConnect, fx.As(new(UserConnector))),
			fx.Annotate(NewUserService, fx.As(new(MountController)), fx.ResultTags(`group:"controller"`)),
		),

		fx.Supply(
			fx.Annotate(ClientConf.Mysql, fx.ResultTags(`name:"mysql_dsn"`)),
			fx.Annotate(ClientConf.Addr, fx.ResultTags(`name:"addr"`)),
			fx.Annotate(ClientConf.ServerDomain, fx.ResultTags(`name:"server_domain"`)),
		),

		fx.Invoke(func(server *Server, lifecycle fx.Lifecycle) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						if err := server.Start(ctx); err != nil && !errors.Is(http.ErrServerClosed, err) {
							panic(err)
						}
					}()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					return server.Stop(ctx)
				},
			})
		}),
	)
	app.Run()
}
