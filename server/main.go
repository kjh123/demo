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
	Addr  string `default:":8000"`
	Mysql string

	KafkaAddr  string `name:"kafka-addr" default:"http://127.0.0.1:9092"`
	KafkaTopic string `name:"kafka-topic"`

	InfluxEnable bool   `name:"enable-influx" default:"false"`
	InfluxHost   string `name:"influx-host" default:"http://127.0.0.1:8086"`
	InfluxToken  string `name:"influx-token"`
	InfluxOrg    string `name:"influx-org" default:"org"`
	InfluxBucket string `name:"influx-bucket" default:"bucket"`

	ClickhouseEnable bool   `name:"enable-clickhouse" default:"false"`
	ClickHouseHost   string `name:"clickhouse-host" default:"http://127.0.0.1:9000"`
	ClickHouseUser   string `name:"clickhouse-user"`
	ClickHousePass   string `name:"clickhouse-pass"`
	ClickHouseDB     string `name:"clickhouse-db"`
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

			fx.Annotate(ServerConf.InfluxEnable, fx.ResultTags(`name:"enable_influx"`)),
			fx.Annotate(ServerConf.InfluxHost, fx.ResultTags(`name:"influx_host"`)),
			fx.Annotate(ServerConf.InfluxToken, fx.ResultTags(`name:"influx_token"`)),
			fx.Annotate(ServerConf.InfluxOrg, fx.ResultTags(`name:"influx_org"`)),
			fx.Annotate(ServerConf.InfluxBucket, fx.ResultTags(`name:"influx_bucket"`)),

			fx.Annotate(ServerConf.ClickhouseEnable, fx.ResultTags(`name:"enable_clickhouse"`)),
			fx.Annotate(ServerConf.ClickHouseHost, fx.ResultTags(`name:"clickhouse_host"`)),
			fx.Annotate(ServerConf.ClickHouseUser, fx.ResultTags(`name:"clickhouse_user"`)),
			fx.Annotate(ServerConf.ClickHousePass, fx.ResultTags(`name:"clickhouse_pass"`)),
			fx.Annotate(ServerConf.ClickHouseDB, fx.ResultTags(`name:"clickhouse_db"`)),
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
