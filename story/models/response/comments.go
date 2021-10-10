package response

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Data        string    `json:"data" db:"data"`
	PostID      uuid.UUID `json:"post_id" db:"post_id"`
	Username    string    `json:"username" db:"username"`
	CommentedAt time.Time `json:"commented_at" db:"created_at"`
}
