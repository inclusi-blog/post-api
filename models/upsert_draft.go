package models

type UpsertDraft struct {
	DraftID  string     `json:"draft_id" binding:"required" db:"DRAFT_ID"`
	UserID   string     `json:"user_id" binding:"required" db:"USER_ID"`
	PostData JSONString `json:"post_data" db:"POST_DATA"`
}

type GetAllDraftRequest struct {
	UserID     string `json:"user_id" binding:"required" db:"USER_ID"`
	StartValue int    `json:"start_value" binding:"required" `
	Limit      int    `json:"limit" binding:"required" `
}
