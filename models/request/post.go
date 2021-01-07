package request

type PostURIRequest struct {
	PostUID string `uri:"post_id" binding:"required,validPostUID"`
}

type CommentPost struct {
	PostUID string `json:"post_uid" binding:"required,validPostUID"`
	Comment string `json:"comment" binding:"required"`
}

type PublishRequest struct {
	DratID string `json:"draft_id" binding:"required"`
}
