package init

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"os"
	"post-api/configuration"
	"post-api/dbhelper"
)

func Db(configData *configuration.ConfigData) neo4j.Session {
	logger := logging.GetLogger(context.Background())
	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }
	driver, err := neo4j.NewDriver(dbhelper.BuildConnectionString(), neo4j.BasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), ""), configForNeo4j40)
	if err != nil {
		logger.Fatal("error occurred while connecting neo4j database")
	}

	if driver == nil {
		logger.Fatal("driver is nil")
	}

	logger.Info("GOT ")
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: os.Getenv("DB_SERVICE_NAME")}
	session, err := driver.NewSession(sessionConfig)
	if err != nil {
		logger.Fatalf("error occurred while creating session neo4j database %v", err)
	}

	return session
}
