package response

import (
	"post-api/story/models"
	"time"
)

type Post struct {
	PostID                 string            `json:"post_id" db:"id"`
	PostData               models.JSONString `json:"post_data" db:"data"`
	LikeCount              int64             `json:"like_count" db:"likes_count"`
	CommentCount           int64             `json:"comment_count" db:"comments_count"`
	Interests              models.JSONString `json:"interests" db:"interests"`
	AuthorID               string            `json:"author_id" db:"author_id"`
	AuthorName             string            `json:"author_name" db:"author_name"`
	PreviewImage           string            `json:"preview_image" db:"preview_image"`
	PublishedAt            time.Time         `json:"published_at" db:"published_at"`
	IsViewerLiked          bool              `json:"is_viewer_liked" db:"is_viewer_liked"`
	IsViewerIsAuthor       bool              `json:"is_viewer_is_author" db:"is_viewer_is_author"`
	IsViewerFollowedAuthor bool              `json:"is_viewer_followed_author"`
}

type PublishedPost struct {
	LikesCount   int64             `json:"likes_count" db:"likes_count"`
	ID           string            `json:"id" db:"id"`
	Title        string            `json:"title" db:"title"`
	Tagline      string            `json:"tagline" db:"tagline"`
	PreviewImage string            `json:"preview_image" db:"preview_image"`
	CreatedAt    time.Time         `json:"published_at" db:"created_at"`
	Interests    models.JSONString `json:"interests" db:"interests"`
}
