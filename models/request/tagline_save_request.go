package request

type TaglineSaveRequest struct {
	UserID  string `json:"user_id" binding:"required" db:"user_id"`
	DraftID string `json:"draft_id" binding:"required" db:"draft_id"`
	Tagline string `json:"tagline" binding:"required" db:"tagline"`
}

type InterestsSaveRequest struct {
	UserID   string `json:"user_id" binding:"required" db:"user_id" `
	DraftID  string `json:"draft_id" binding:"required" db:"draft_id" `
	Interest string `json:"interest" binding:"required" db:"interest" `
}

type PreviewImageSaveRequest struct {
	UserID          string `json:"user_id" binding:"required" db:"user_id"`
	DraftID         string `json:"draft_id" binding:"required" db:"draft_id"`
	PreviewImageUrl string `json:"preview_image" binding:"required" db:"preview_image"`
}

type DraftURIRequest struct {
	DraftID string `uri:"draft_id" binding:"required,validPostUID"`
}
