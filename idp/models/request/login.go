package request

type UserLoginRequest struct {
	Email          string `json:"email" binding:"required" example:"someuser@gmail.com"`
	Password       string `json:"password" binding:"required" example:"encrypted password"`
	LoginChallenge string `json:"login_challenge" binding:"required" example:"login challenge from hydra"`
}
