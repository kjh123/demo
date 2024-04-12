package logs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/pkg/errors"

	"go.uber.org/fx"
)

const logTableName = "logs"

type BehaviorLog struct {
	UID int64 `json:"uid"`
	// UUID string `json:"uuid"`
	IP   string   `json:"ip"`
	Tags []string `json:"tags"`
	UA   string   `json:"ua"`
	// AppVersion  string   `json:"app_version"`
	// DeviceID    string   `json:"device_id"`
	// DeviceType  string   `json:"device_type"`
	// DeviceName  string   `json:"device_name"`
	// DeviceOS    string   `json:"device_os"`
	// DeviceModel string   `json:"device_model"`
	// DeviceLoc   string   `json:"device_loc"`
	// DeviceLang  string   `json:"device_lang"`
	// DeviceIP    string   `json:"device_ip"`
}

type LogWriter interface {
	Writer(ctx context.Context, log BehaviorLog) error
}

var Module = fx.Module(
	"data",
	fx.Provide(
		setupInflux,
		setupClickHouse,

		// fx.Annotate(NewInfluxRepository, fx.As(new(LogWriter))),
		fx.Annotate(NewClickHouseRepository, fx.As(new(LogWriter))),
	),

	fx.Invoke(func(client *InfluxClient, clickHouseClient *ClickHouseClient, lifecycle fx.Lifecycle) {
		lifecycle.Append(fx.Hook{OnStop: func(_ context.Context) error {
			client.Close()
			return clickHouseClient.Close()
		}})
	}),
)

type Params struct {
	fx.In

	EnableInflux bool   `name:"enable_influx"`
	InfluxHost   string `name:"influx_host"`
	InfluxToken  string `name:"influx_token"`
	InfluxOrg    string `name:"influx_org"`
	InfluxBucket string `name:"influx_bucket"`

	EnableClickHouse bool     `name:"enable_clickhouse"`
	ClickHouseHost   []string `name:"clickhouse_host"`
	ClickHouseUser   string   `name:"clickhouse_user"`
	ClickHousePass   string   `name:"clickhouse_pass"`
	ClickHouseDB     string   `name:"clickhouse_db"`
}

type InfluxClient struct {
	Client      influxdb2.Client
	WriteAPI    api.WriteAPIBlocking
	QueryAPI    api.QueryAPI
	Measurement string
}

func (i *InfluxClient) Close() {
	if i == nil {
		return
	}

	i.Client.Close()
}

func setupInflux(params Params) *InfluxClient {
	if !params.EnableInflux {
		return nil
	}

	client := influxdb2.NewClient(params.InfluxHost, params.InfluxToken)
	if ok, err := client.Ping(context.Background()); !ok || err != nil {
		panic(err)
	}

	return &InfluxClient{
		Client:      client,
		WriteAPI:    client.WriteAPIBlocking(params.InfluxOrg, params.InfluxBucket),
		QueryAPI:    client.QueryAPI(params.InfluxOrg),
		Measurement: logTableName,
	}
}

type ClickHouseClient struct {
	Conn  driver.Conn
	DB    string
	Table string
}

func (c *ClickHouseClient) Close() error {
	if c == nil {
		return nil
	}

	return c.Conn.Close()
}

const (
	ClickHouseDBSQL    = `CREATE DATABASE IF NOT EXISTS %s`
	ClickHouseTableSQL = `CREATE TABLE IF NOT EXISTS %s.%s 
(
	date Date DEFAULT toDate(timestamp),
	uid UInt64,
    ip String,
    tags Array(String),
	ua String,
    timestamp DateTime64(3, 'Asia/Shanghai') DEFAULT now()
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(date) -- 按月分区
ORDER BY (date, timestamp) -- 按日期和时间戳排序
SETTINGS index_granularity = 8192;
`
) // MergeTree 系列的引擎被设计用于插入极大量的数据到一张表当中。数据可以以数据片段的形式一个接着一个的快速写入，数据片段在后台按照一定的规则进行合并。相比在插入时不断修改（重写）已存储的数据，这种策略会高效很多。

func setupClickHouse(params Params) *ClickHouseClient {
	if !params.EnableClickHouse {
		return nil
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: params.ClickHouseHost,
		Auth: clickhouse.Auth{
			Database: params.ClickHouseDB,
			Username: params.ClickHouseUser,
			Password: params.ClickHousePass,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout: time.Duration(10) * time.Second,
		Debug:       true,
	})
	if err != nil {
		panic(errors.Wrap(err, "connect clickhouse hose"))
	}

	v, err := conn.ServerVersion()
	if err != nil {
		panic(errors.Wrap(err, "clickhouse server version error"))
	}

	if err = conn.Exec(context.Background(), fmt.Sprintf(ClickHouseDBSQL, params.ClickHouseDB)); err != nil {
		panic(errors.Wrap(err, "create clickhouse db"))
	}

	if err = conn.Exec(context.Background(), fmt.Sprintf(ClickHouseTableSQL, params.ClickHouseDB, logTableName)); err != nil {
		panic(errors.Wrap(err, "create clickhouse table"))
	}

	slog.Info("clickhouse version", "version", v.Version.String())
	return &ClickHouseClient{Conn: conn, DB: params.ClickHouseDB, Table: logTableName}
}
