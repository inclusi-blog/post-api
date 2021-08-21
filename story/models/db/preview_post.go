package db

import "github.com/google/uuid"

type AbstractPost struct {
	Model
	PostID       uuid.UUID `json:"post_id" db:"post_id"`
	Title        string    `json:"title" db:"title"`
	Tagline      string    `json:"tagline" db:"tagline"`
	PreviewImage string    `json:"preview_image" db:"preview_image"`
	ViewTime     int64     `json:"view_time" db:"view_time"`
}
