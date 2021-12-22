package request

import "github.com/google/uuid"

type Report struct {
	PostID uuid.UUID
	UserID uuid.UUID
	Reason string `json:"reason,omitempty"`
}
