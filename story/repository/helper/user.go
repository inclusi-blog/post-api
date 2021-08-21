package helper

type CreateUserRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}
