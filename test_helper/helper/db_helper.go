package helper

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/inclusi-blog/gola-utils/logging"
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
	var tableInOrderDeletion = []string{
		"abstract_post",
		"post_x_interests",
		"comments",
		"likes",
		"posts",
		"abstract_post",
		"drafts",
		"users",
	}

	tx := dbHelper.db.MustBegin()
	for _, tableName := range tableInOrderDeletion {
		if tableName == "users" {
			sql := fmt.Sprintf("delete from %s where role_id != (select id from roles where name = 'Admin')", tableName)
			err := dbHelper.retryDelete(tx, sql, tableName)
			if err != nil {
				return err
			}
			continue
		}
		sql := fmt.Sprintf("DELETE FROM %s", tableName)
		err := dbHelper.retryDelete(tx, sql, tableName)
		if err != nil {
			return err
		}
	}

	_ = tx.Commit()
	return nil
}

func (dbHelper DbHelper) retryDelete(tx *sqlx.Tx, sql string, tableName string) error {
	logger := logging.GetLogger(context.Background())
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
	return nil
}

func (dbHelper DbHelper) GetInterestsIDs(interestNames []string) ([]uuid.UUID, error) {
	var interestsIDs []uuid.UUID
	query, args, err := sqlx.In("SELECT id from interests where name in (?);", interestNames)

	query = dbHelper.db.Rebind(query)
	rows, err := dbHelper.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		interestsIDs = append(interestsIDs, id)
	}
	if err != nil {
		return nil, err
	}
	return interestsIDs, nil
}
