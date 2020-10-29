package db

import "post-api/models"

type PublishPost struct {
	PUID      string            `json:"puid" db:"puid"`
	UserID    string            `json:"user_id" db:"user_id"`
	PostData  models.JSONString `json:"post_data" db:"post_data"`
	ReadTime  int               `json:"read_time" db:"read_time"`
	ViewCount int               `json:"view_count" db:"view_count"`
}
