package models

import "github.com/google/uuid"

type UpsertDraft struct {
	DraftID uuid.UUID
	UserID  uuid.UUID
	Data    JSONString `json:"data" db:"data"`
}

type GetAllDraftRequest struct {
	UserID     uuid.UUID
	StartValue int `json:"start_value"`
	Limit      int `json:"limit"`
}

type CreateDraft struct {
	Data   JSONString `json:"data"`
	UserID uuid.UUID
}
