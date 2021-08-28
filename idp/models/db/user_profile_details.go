package db

type UserProfile struct {
	UserID   string `json:"user_id" db:"uuid"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	IsActive bool   `json:"isActive" db:"is_active"`
}
