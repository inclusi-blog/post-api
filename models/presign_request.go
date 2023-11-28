package models

type CoverPreSign struct {
	Extension string `form:"extension" binding:"oneof=png jpg jpeg webp svg"`
}

type UploadImage struct {
	UploadID string `json:"upload_id" validate:"required"`
}

type Pagination struct {
	Start int `form:"start"`
	Limit int `form:"limit"`
}
