package main

import (
	"context"
	"data"
	"github.com/alecthomas/kong"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

var ServerConf struct {
	Addr  string `default:":8000"`
	Mysql string
}

func main() {
	kong.Parse(&ServerConf)
	app := fx.New(
		data.Module,
		fx.Provide(
			NewServer,
		),

		fx.Supply(fx.Annotate(ServerConf.Mysql, fx.ResultTags(`name:"mysql_dsn"`))),

		fx.Invoke(func(s *grpc.Server, lifecycle fx.Lifecycle) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
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
				OnStop: func(ctx context.Context) error {
					s.Stop()
					return nil
				},
			})
		}),

		fx.Invoke(setupServer),
	)
	app.Run()
}
