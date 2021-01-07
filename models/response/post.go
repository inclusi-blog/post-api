package response

import "post-api/models"

type LikedByCount struct {
	LikeCount int64 `json:"like_count"`
}

type Post struct {
	PostID                 string            `json:"post_id" db:"postId"`
	PostData               models.JSONString `json:"post_data" db:"data"`
	LikeCount              int64             `json:"like_count" db:"likeCount"`
	CommentCount           int64             `json:"comment_count" db:"commentCount"`
	Interests              []string          `json:"interests" db:"interests"`
	AuthorID               string            `json:"author_id" db:"authorID"`
	AuthorName             string            `json:"author_name" db:"authorName"`
	PreviewImage           string            `json:"preview_image" db:"previewImage"`
	PublishedAt            int64             `json:"published_at" db:"publishedAt"`
	IsViewerLiked          bool              `json:"is_viewer_liked" db:"isViewerLiked"`
	IsViewerIsAuthor       bool              `json:"is_viewer_is_author" db:"isAuthorViewing"`
	IsViewerFollowedAuthor bool              `json:"is_viewer_followed_author"`
}

type DBPost struct {
	PostID                 string   `json:"post_id" db:"postId"`
	PostData               string   `json:"post_data" db:"data"`
	LikeCount              int64    `json:"like_count" db:"likeCount"`
	CommentCount           int64    `json:"comment_count" db:"commentCount"`
	Interests              []string `json:"interests" db:"interests"`
	AuthorID               string   `json:"author_id" db:"authorID"`
	AuthorName             string   `json:"author_name" db:"authorName"`
	PreviewImage           string   `json:"preview_image" db:"previewImage"`
	PublishedAt            int64    `json:"published_at" db:"publishedAt"`
	IsViewerLiked          bool     `json:"is_viewer_liked" db:"isViewerLiked"`
	IsViewerIsAuthor       bool     `json:"is_viewer_is_author" db:"isAuthorViewing"`
	IsViewerFollowedAuthor bool     `json:"is_viewer_followed_author"`
}
