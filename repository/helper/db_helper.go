package helper

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type DbHelper struct {
	db neo4j.Session
}

func NewDbHelper(db neo4j.Session) DbHelper {
	return DbHelper{
		db: db,
	}
}

func (dbHelper DbHelper) ClearAll() error {
	logger := logging.GetLogger(context.Background())
	_, err := dbHelper.db.Run("MATCH (node) detach delete node", map[string]interface{}{})
	if err != nil {
		logger.Error("unable to clear data")
		return err
	}
	return nil
}
