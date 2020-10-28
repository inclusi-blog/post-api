package db

type PreviewPost struct {
	Model
	PostID       int64  `json:"post_id" db:"post_id"`
	Title        string `json:"title" db:"title"`
	Tagline      string `json:"tagline" db:"tagline"`
	PreviewImage string `json:"preview_image" db:"preview_image"`
	LikeCount    int64  `json:"like_count" db:"like_count"`
	CommentCount int64  `json:"comment_count" db:"comment_count"`
	ViewTime     int64  `json:"view_time" db:"view_time"`
}
