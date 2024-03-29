package data

import (
	"database/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	client   *sqlx.DB
	fixtures *testfixtures.Loader
)

func TestMain(m *testing.M) {
	if mysqlDSN := os.Getenv("MYSQL"); mysqlDSN != "" {
		dsn := parseMySQLDSN(mysqlDSN)
		client = setupData(Params{Mysql: dsn.FormatDSN()})
		defer client.Close()

		db, err := sql.Open("mysql", dsn.FormatDSN())
		defer db.Close()
		if err != nil {
			panic(err)
		}

		fixtures, err = testfixtures.New(
			testfixtures.DangerousSkipTestDatabaseCheck(),
			testfixtures.Database(db),
			testfixtures.Dialect("mysql"),
			testfixtures.Directory("fixtures"),
		)
		if err != nil {
			panic(err)
		}
	}

	os.Exit(m.Run())
}

func LoadFixtures(t *testing.T) {
	if fixtures == nil {
		t.Skip("fixtures not init.")
	}

	assert.NoError(t, fixtures.Load())
}
