package request

import "github.com/google/uuid"

type InitiateRegistrationRequest struct {
	Email    string    `json:"email" binding:"required" example:"email"`
	Password string    `json:"password" binding:"required" example:"encrypted-password"`
	Id       uuid.UUID `json:"id" binding:"required"`
}

type NewRegistrationRequest struct {
	Email    string `json:"email" binding:"required" example:"email"`
	Password string `json:"password" binding:"required" example:"encrypted-password"`
}

type EmailAvailabilityRequest struct {
	Email string `json:"email" binding:"required" example:"email"`
}

type UsernameAvailabilityRequest struct {
	Username string `json:"username" binding:"required" example:"username"`
}
