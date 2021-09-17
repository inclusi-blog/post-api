package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/jmoiron/sqlx/types"
)

type JSONString struct {
	types.JSONText
}

func (j JSONString) Value() (driver.Value, error) {
	var m json.RawMessage
	var err = j.JSONText.Unmarshal(&m)
	if err != nil {
		return "", err
	}
	return string(j.JSONText), nil
}
