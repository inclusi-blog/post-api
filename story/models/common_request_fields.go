package models

import "github.com/jmoiron/sqlx/types"

type JSONString struct {
	types.JSONText
}
