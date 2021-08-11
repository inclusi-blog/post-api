package db

import (
	"github.com/google/uuid"
	"post-api/models"
)

type PublishPost struct {
	UserID   uuid.UUID         `json:"author_id" db:"author_id"`
	PostData models.JSONString `json:"data" db:"data"`
	DraftID  uuid.UUID         `json:"draft_id" db:"draft_id"`
}

type LikedByRes struct {
	LikedByID string `json:"id" db:"id"`
}
