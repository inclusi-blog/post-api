package db

import "post-api/models"

type Draft struct {
	DraftID   string            `json:"draft_id" db:"draft_id"`
	UserID    string            `json:"user_id" db:"user_id"`
	PostData  models.JSONString `json:"post_data" db:"post_data"`
	TitleData models.JSONString `json:"title_data" db:"title_data"`
	Tagline   string            `json:"tagline" db:"tagline"`
	Interest  models.JSONString `json:"interest" db:"interest"`
}

type AllDraft struct {
	DraftID   string            `json:"draft_id" db:"draft_id"`
	UserID    string            `json:"user_id" db:"user_id"`
	PostData  models.JSONString `json:"post_data" db:"post_data"`
	TitleData string            `json:"title_data" db:"title_data"`
	Tagline   string            `json:"tagline" db:"tagline"`
	Interest  models.JSONString `json:"interest" db:"interest"`
}
