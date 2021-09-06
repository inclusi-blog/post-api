package db

import "github.com/google/uuid"

type Interests struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
