package request

type PostURIRequest struct {
	PostUID string `uri:"post_id" binding:"required,validPostUID"`
}

type PostLikeRequest struct {
	PostUID string `uri:"post_id" binding:"required,validPostUID"`
}
