package response

import "github.com/google/uuid"

type PreviewDraft struct {
	DraftID      uuid.UUID `json:"draft_id" binding:"required"`
	Title        string    `json:"title" binding:"required"`
	Tagline      string    `json:"tagline" binding:"required"`
	Interest     []string  `json:"interest" binding:"required"`
	PreviewImage string    `json:"preview_image" binding:"required"`
	AuthorName   string    `json:"author_name" binding:"required"`
}
