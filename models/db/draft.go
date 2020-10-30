package db

import (
	"database/sql"
	"post-api/models"
)

type Draft struct {
	DraftID      string            `json:"draft_id" db:"draft_id"`
	UserID       string            `json:"user_id" db:"user_id"`
	PostData     models.JSONString `json:"post_data" db:"post_data"`
	PreviewImage string            `json:"preview_image" db:"preview_image"`
	Tagline      string            `json:"tagline" db:"tagline"`
	Interest     models.JSONString `json:"interest" db:"interest"`
}

type DraftDB struct {
	DraftID      string            `json:"draft_id" db:"draft_id"`
	UserID       string            `json:"user_id" db:"user_id"`
	PostData     models.JSONString `json:"post_data" db:"post_data"`
	PreviewImage sql.NullString    `json:"preview_image" db:"preview_image"`
	Tagline      sql.NullString    `json:"tagline" db:"tagline"`
	Interest     models.JSONString `json:"interest" db:"interest"`
}

type AllDraft struct {
	DraftID   string            `json:"draft_id" db:"draft_id"`
	UserID    string            `json:"user_id" db:"user_id"`
	PostData  models.JSONString `json:"post_data" db:"post_data"`
	TitleData string            `json:"title_data" db:"title_data"`
	Tagline   string            `json:"tagline" db:"tagline"`
	Interest  models.JSONString `json:"interest" db:"interest"`
}
