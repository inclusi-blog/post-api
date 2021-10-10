package request

import "github.com/google/uuid"

type Comment struct {
	Data        string `json:"data" binding:"required" db:"data"`
	PostID      uuid.UUID
	CommentedBy uuid.UUID `db:"commented_by"`
}

type FetchComments struct {
	PostID uuid.UUID
	Start  int `form:"start"`
	Limit  int `form:"limit" binding:"required"`
}
