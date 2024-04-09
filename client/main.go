package main

import (
	"context"
	"data"
	"event"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var ClientConf struct {
	ServerDomain string `default:"127.0.0.1:8000"`

	Addr       string
	Mysql      string
	KafkaAddr  string `name:"kafka-addr"`
	KafkaTopic string `name:"kafka-topic"`
}

func main() {
	kong.Parse(&ClientConf)
	app := fx.New(
		event.Module,
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

			fx.Annotate(NewPusher, fx.As(new(MountController)), fx.ResultTags(`group:"controller"`)),
		),

		fx.Supply(
			fx.Annotate(ClientConf.Mysql, fx.ResultTags(`name:"mysql_dsn"`)),
			fx.Annotate(ClientConf.Addr, fx.ResultTags(`name:"addr"`)),
			fx.Annotate(ClientConf.ServerDomain, fx.ResultTags(`name:"server_domain"`)),
			fx.Annotate(ClientConf.KafkaAddr, fx.ResultTags(`name:"kafka_addr"`)),
			fx.Annotate(ClientConf.KafkaTopic, fx.ResultTags(`name:"kafka_topic"`)),
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
