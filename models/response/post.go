package response

import "post-api/models"

type LikedByCount struct {
	LikeCount int64 `json:"like_count"`
}

type Post struct {
	PostID                 string            `json:"post_id"`
	PostData               models.JSONString `json:"post_data" db:"data"`
	LikeCount              int64             `json:"like_count" db:"likeCount"`
	CommentCount           int64             `json:"comment_count" db:"commentCount"`
	Interests              []string          `json:"interests" db:"interests"`
	AuthorID               string            `json:"author_id" db:"authorID"`
	PreviewImage           string            `json:"preview_image" db:"previewImage"`
	PublishedAt            int64             `json:"published_at" db:"publishedAt"`
	IsViewerLiked          bool              `json:"is_viewer_liked" db:"isViewerLiked"`
	IsViewIsAuthor         bool              `json:"is_view_is_author" db:"isAuthorViewing"`
	IsViewerFollowedAuthor bool              `json:"is_viewer_followed_author"`
}

type DBPost struct {
	PostID                 string   `json:"post_id"`
	PostData               string   `json:"post_data" db:"data"`
	LikeCount              int64    `json:"like_count" db:"likeCount"`
	CommentCount           int64    `json:"comment_count" db:"commentCount"`
	Interests              []string `json:"interests" db:"interests"`
	AuthorID               string   `json:"author_id" db:"authorID"`
	PreviewImage           string   `json:"preview_image" db:"previewImage"`
	PublishedAt            int64    `json:"published_at" db:"publishedAt"`
	IsViewerLiked          bool     `json:"is_viewer_liked" db:"isViewerLiked"`
	IsViewIsAuthor         bool     `json:"is_view_is_author" db:"isAuthorViewing"`
	IsViewerFollowedAuthor bool     `json:"is_viewer_followed_author"`
}
