package request

type PostLikeRequest struct {
	PostUID string `uri:"post_id" binding:"required,validPostUID"`
}
