package request

import (
	"github.com/google/uuid"
)

type TaglineSaveRequest struct {
	UserID  uuid.UUID `json:"-"`
	DraftID uuid.UUID `json:"-"`
	Tagline string    `json:"tagline" binding:"required" db:"tagline"`
}

type InterestsSaveRequest struct {
	UserID    uuid.UUID `json:"-"`
	DraftID   uuid.UUID `json:"-"`
	Interests []string  `json:"interests" binding:"required" db:"interest" `
}

type PreviewImageSaveRequest struct {
	UserID   uuid.UUID `json:"-"`
	DraftID  uuid.UUID `json:"-"`
	UploadID string    `json:"upload_id" binding:"required" db:"preview_image"`
}

type DraftURIRequest struct {
	DraftID string `uri:"draft_id" binding:"required,validPostUID"`
}
