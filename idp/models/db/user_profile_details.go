package db

import "github.com/google/uuid"

type UserProfile struct {
	Id       uuid.UUID `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	IsActive bool   `json:"isActive" db:"is_active"`
}
