package init

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/tracing"
	"github.com/jmoiron/sqlx"
	"post-api/configuration"
	"post-api/dbhelper"
)

func Db(configData *configuration.ConfigData) *sqlx.DB {
	logger := logging.GetLogger(context.Background())
	connectionString := dbhelper.BuildConnectionString()
	db, err := tracing.InitSqlxOracleDBWithInstrumentationAndConnectionConfig("postgres", connectionString, configData.GetDBConnectionPoolConfig())
	if err != nil {
		logger.Panic("Could not connect to POST DB", err)
	}

	return db
}
