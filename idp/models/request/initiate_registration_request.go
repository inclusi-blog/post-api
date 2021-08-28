package request

type InitiateRegistrationRequest struct {
	Email    string `json:"email" binding:"required" example:"email"`
	Password string `json:"password" binding:"required" example:"encrypted-password"`
	Username string `json:"username" binding:"required" example:"user123"`
	UUID     string `json:"uuid" binding:"required"`
}

type NewRegistrationRequest struct {
	Email    string `json:"email" binding:"required" example:"email"`
	Password string `json:"password" binding:"required" example:"encrypted-password"`
	Username string `json:"username" binding:"required" example:"user123"`
}

type EmailAvailabilityRequest struct {
	Email string `json:"email" binding:"required" example:"email"`
}

type UsernameAvailabilityRequest struct {
	Username string `json:"username" binding:"required" example:"username"`
}
