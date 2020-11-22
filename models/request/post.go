package request

type LikedByCount struct {
	LikeCount string `json:"like_count"`
}

type PostLikeRequest struct {
	PostUID string `uri:"post_id" binding:"required,validPostUID"`
}
