package request

type ForgetPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Password string `json:"password"`
}
