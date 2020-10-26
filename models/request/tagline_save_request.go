package request

import "post-api/models"

type TaglineSaveRequest struct {
	UserID  string `json:"user_id" binding:"required" db:"USER_ID"`
	DraftID string `json:"draft_id" binding:"required" db:"DRAFT_ID"`
	Tagline string `json:"tagline" binding:"required" db:"TAGLINE"`
}

type InterestsSaveRequest struct {
	UserID    string            `json:"user_id" binding:"required" db:"USER_ID"`
	DraftID   string            `json:"draft_id" binding:"required" db:"DRAFT_ID"`
	Interests models.JSONString `json:"interests" binding:"required" db:"INTEREST"`
}
