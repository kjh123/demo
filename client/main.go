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

	Addr  string
	Mysql string

	KafkaAddr  string `name:"kafka-addr" default:"http://127.0.0.1:9092"`
	KafkaTopic string `name:"kafka-topic"`

	InfluxEnable bool   `name:"enable-influx" default:"false"`
	InfluxHost   string `name:"influx-host" default:"http://127.0.0.1:8086"`
	InfluxToken  string `name:"influx-token"`
	InfluxOrg    string `name:"influx-org" default:"org"`
	InfluxBucket string `name:"influx-bucket" default:"bucket"`

	ClickhouseEnable bool     `name:"enable-clickhouse" default:"true"`
	ClickHouseHost   []string `name:"clickhouse-host" default:"127.0.0.1:9000"`
	ClickHouseUser   string   `name:"clickhouse-user" default:"default"`
	ClickHousePass   string   `name:"clickhouse-pass"`
	ClickHouseDB     string   `name:"clickhouse-db"`
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
			fx.Annotate(NewLogService, fx.As(new(MountController)), fx.ResultTags(`group:"controller"`)),
		),

		fx.Supply(
			fx.Annotate(ClientConf.Addr, fx.ResultTags(`name:"addr"`)),
			fx.Annotate(ClientConf.Mysql, fx.ResultTags(`name:"mysql_dsn"`)),
			fx.Annotate(ClientConf.ServerDomain, fx.ResultTags(`name:"server_domain"`)),

			fx.Annotate(ClientConf.KafkaAddr, fx.ResultTags(`name:"kafka_addr"`)),
			fx.Annotate(ClientConf.KafkaTopic, fx.ResultTags(`name:"kafka_topic"`)),

			fx.Annotate(ClientConf.InfluxEnable, fx.ResultTags(`name:"enable_influx"`)),
			fx.Annotate(ClientConf.InfluxHost, fx.ResultTags(`name:"influx_host"`)),
			fx.Annotate(ClientConf.InfluxToken, fx.ResultTags(`name:"influx_token"`)),
			fx.Annotate(ClientConf.InfluxOrg, fx.ResultTags(`name:"influx_org"`)),
			fx.Annotate(ClientConf.InfluxBucket, fx.ResultTags(`name:"influx_bucket"`)),

			fx.Annotate(ClientConf.ClickhouseEnable, fx.ResultTags(`name:"enable_clickhouse"`)),
			fx.Annotate(ClientConf.ClickHouseHost, fx.ResultTags(`name:"clickhouse_host"`)),
			fx.Annotate(ClientConf.ClickHouseUser, fx.ResultTags(`name:"clickhouse_user"`)),
			fx.Annotate(ClientConf.ClickHousePass, fx.ResultTags(`name:"clickhouse_pass"`)),
			fx.Annotate(ClientConf.ClickHouseDB, fx.ResultTags(`name:"clickhouse_db"`)),
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
