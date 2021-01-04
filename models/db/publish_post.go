package db

import "post-api/models"

type PublishPost struct {
	PUID         string            `json:"puid" db:"puid"`
	UserID       string            `json:"user_id" db:"user_id"`
	PostData     models.JSONString `json:"post_data" db:"post_data"`
	ReadTime     int               `json:"read_time" db:"read_time"`
	Interest     []string          `json:"interest"`
	Title        string            `json:"title"`
	Tagline      string            `json:"tagline"`
	PreviewImage string            `json:"preview_image"`
	PostUrl      string            `json:"post_url"`
}

type LikedByRes struct {
	LikedByID string `json:"id" db:"id"`
}
