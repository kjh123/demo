package data

import (
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"data",
	fx.Provide(
		setupData,

		fx.Annotate(NewUserRepository, fx.As(new(UserQueryer))),
		fx.Annotate(NewUserRepository, fx.As(new(UserCommander))),
	),
)

type Params struct {
	fx.In

	Mysql string `name:"mysql_dsn"`
}

var schema = `
CREATE TABLE IF NOT EXISTS users(
	id int unsigned auto_increment PRIMARY KEY,
	name VARCHAR(20) not null DEFAULT '',
	created_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP
);`

func setupData(params Params) *sqlx.DB {
	dsn := parseMySQLDSN(params.Mysql)
	db := sqlx.MustConnect("mysql", dsn.FormatDSN())
	db.MustExec(schema)
	return db
}

func parseMySQLDSN(mysqlDSN string) *mysql.Config {
	dsn, err := mysql.ParseDSN(mysqlDSN)
	if err != nil {
		panic(err)
	}
	dsn.ParseTime = true
	return dsn
}
