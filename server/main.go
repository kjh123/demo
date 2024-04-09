package main

import (
	"context"
	"data"
	"event"
	"log/slog"
	"net"

	"github.com/alecthomas/kong"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var ServerConf struct {
	Addr       string `default:":8000"`
	Mysql      string
	KafkaAddr  string `name:"kafka-addr"`
	KafkaTopic string `name:"kafka-topic"`
}

func main() {
	kong.Parse(&ServerConf)
	app := fx.New(
		event.Module,
		data.Module,
		fx.Provide(
			NewServer,
		),

		fx.Supply(
			fx.Annotate(ServerConf.Mysql, fx.ResultTags(`name:"mysql_dsn"`)),
			fx.Annotate(ServerConf.KafkaAddr, fx.ResultTags(`name:"kafka_addr"`)),
			fx.Annotate(ServerConf.KafkaTopic, fx.ResultTags(`name:"kafka_topic"`)),
		),

		fx.Invoke(func(s *grpc.Server, lifecycle fx.Lifecycle) {
			lifecycle.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					lis, err := net.Listen("tcp", ServerConf.Addr)
					if err != nil {
						return err
					}

					go func(lis net.Listener) {
						slog.Info("server start and listening at:", "addr", lis.Addr())
						if err := s.Serve(lis); err != nil {
							panic(err)
						}
					}(lis)
					return nil
				},
				OnStop: func(_ context.Context) error {
					s.Stop()
					return nil
				},
			})
		}),

		fx.Invoke(setupServer),
		fx.Invoke(Receiver),
	)
	app.Run()
}
