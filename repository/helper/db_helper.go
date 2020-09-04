package helper

import (
	"context"
	"fmt"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
)

type DbHelper struct {
	db *sqlx.DB
}

func NewDbHelper(db *sqlx.DB) DbHelper {
	return DbHelper{
		db: db,
	}
}

func (dbHelper DbHelper) ClearAll() error {
	var tableInOrderDeletion = []string {
		"DRAFTS",
	}

	tx := dbHelper.db.MustBegin()
	logger := logging.GetLogger(context.Background())
	for _, tableName := range tableInOrderDeletion {
		sql := fmt.Sprintf("DELETE FROM %s", tableName)
		_, e := tx.Exec(sql)
		if e != nil {
			logger.Error("Delete failed for", tableName, " : ", e.Error())
			logger.Error("Retry attempted")
			_, e := tx.Exec(sql)
			if e != nil {
				logger.Error("Delete failed again ", tableName, " : ", e.Error())
				_ = tx.Rollback()
				return e
			}
		}
	}

	_ = tx.Commit()
	return nil
}
