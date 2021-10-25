package request

import "github.com/google/uuid"

type PostURIRequest struct {
	PostUID string `uri:"post_id" binding:"required,validPostUID"`
}

type PostLikeRequest struct {
	PostUID string `uri:"post_id" binding:"required,validPostUID"`
}

type GetPublishedPostRequest struct {
	UserID     uuid.UUID
	StartValue int `json:"start_value"`
	Limit      int `json:"limit"`
}

type PostRequest struct {
	UserID uuid.UUID
	Start  int `form:"start"`
	Limit  int `form:"limit" binding:"required"`
}
