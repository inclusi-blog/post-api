package request

type UserDetailsUpdate struct {
	Name        string `json:"name" binding:"omitempty,min=0,max=100"`
	About       string `json:"about" binding:"omitempty,min=0,max=108"`
	Username    string `json:"username" binding:"omitempty,min=0,max=100"`
	FacebookURL string `json:"facebook_url" binding:"omitempty"`
	LinkedInURL string `json:"linked_in_url" binding:"omitempty"`
	TwitterURL  string `json:"twitter_url" binding:"omitempty"`
}

type CoverPreSign struct {
	Extension string `form:"extension" binding:"oneof=png jpg jpeg webp svg"`
}

type UploadImage struct {
	UploadID string `json:"upload_id" validate:"required"`
}
