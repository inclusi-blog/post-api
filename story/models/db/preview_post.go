package db

import (
	"github.com/google/uuid"
	"time"
)

type AbstractPost struct {
	Model
	PostID       uuid.UUID `json:"post_id" db:"post_id"`
	Title        string    `json:"title" db:"title"`
	Tagline      string    `json:"tagline" db:"tagline"`
	PreviewImage string    `json:"preview_image" db:"preview_image"`
	ViewTime     int64     `json:"view_time" db:"view_time"`
	URL          string    `json:"url" db:"url"`
}

type HomeFeedPost struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Tagline       string    `json:"tagline" db:"tagline"`
	ViewTime      int64     `json:"view_time" db:"view_time"`
	PublishedDate time.Time `json:"published_date" db:"published_date"`
	InterestNames []string  `json:"interest_names" db:"interest_names"`
	AuthorName    string    `json:"author_name" db:"author_name"`
	LikeCount     int64     `json:"like_count" db:"like_count"`
	UserLiked     bool      `json:"user_liked" db:"user_liked"`
	PreviewImage  string    `json:"preview_image" db:"preview_image"`
}
