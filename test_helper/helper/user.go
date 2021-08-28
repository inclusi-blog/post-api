package helper

type CreateUserRequest struct {
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
	Username string `json:"username"`
}
