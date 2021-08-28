package response

type EmailAvailabilityResponse struct {
	IsAvailable bool `json:"isEmailAvailable"`
}

type UsernameAvailabilityResponse struct {
	IsAvailable bool `json:"isUsernameAvailable"`
}
