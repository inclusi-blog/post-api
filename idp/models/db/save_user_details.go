package db

import "github.com/google/uuid"

type SaveUserDetails struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	IsActive bool      `db:"is_active"`
}
