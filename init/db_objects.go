package init

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/inclusi-blog/gola-utils/tracing"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"post-api/configuration"
	"post-api/dbhelper"
)

func Db(configData *configuration.ConfigData, awsSession *session.Session) *sqlx.DB {
	logger := logging.GetLogger(context.Background())
	connectionString := dbhelper.BuildConnectionString()
	if configData.Environment != "local" {
		var err error
		connectionString, err = dbhelper.BuildConnectionStringCloud(awsSession)
		if err != nil {
			logger.Panic("unable to get value from aws")
		}
	}
	db, err := tracing.InitPostgresDBWithInstrumentationAndConnectionConfig("postgres", connectionString, configData.GetDBConnectionPoolConfig())
	if err != nil {
		//logger.Panic("Could not connect to POST DB", err)
	}

	return db
}
