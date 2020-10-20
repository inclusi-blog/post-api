package request

type TaglineSaveRequest struct {
	UserID  string `json:"user_id" db:"USER_ID"`
	DraftID string `json:"draft_id" db:"DRAFT_ID"`
	Tagline string `json:"tagline" db:"TAGLINE"`
}
