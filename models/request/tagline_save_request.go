package request

type TaglineSaveRequest struct {
	UserID  string `json:"user_id" binding:"required" db:"USER_ID"`
	DraftID string `json:"draft_id" binding:"required" db:"DRAFT_ID"`
	Tagline string `json:"tagline" binding:"required" db:"TAGLINE"`
}
