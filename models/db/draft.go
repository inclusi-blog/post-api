package db

import (
	"post-api/models"
)

type Draft struct {
	DraftID      string            `json:"draft_id" db:"draft_id"`
	UserID       string            `json:"user_id" db:"user_id"`
	PostData     models.JSONString `json:"post_data" db:"post_data"`
	PreviewImage *string           `json:"preview_image" db:"preview_image"`
	Tagline      *string           `json:"tagline" db:"tagline"`
	Interest     []string          `json:"interest" db:"interest"`
}

type DraftDB struct {
	DraftID      string            `json:"draft_id" db:"draftId"`
	UserID       string            `json:"user_id" db:"userId"`
	PostData     models.JSONString `json:"post_data" db:"postData"`
	PreviewImage string            `json:"preview_image" db:"previewImage"`
	Tagline      string            `json:"tagline" db:"tagline"`
	Interest     []string          `json:"interest" db:"interests"`
	IsPublished  bool              `json:"is_published"`
}

type AllDraft struct {
	DraftID   string            `json:"draft_id" db:"draft_id"`
	UserID    string            `json:"user_id" db:"user_id"`
	PostData  models.JSONString `json:"post_data" db:"post_data"`
	TitleData string            `json:"title_data" db:"title_data"`
	Tagline   *string           `json:"tagline" db:"tagline"`
	Interest  []string          `json:"interest" db:"interest"`
}
