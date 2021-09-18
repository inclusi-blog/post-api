package db

import (
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/google/uuid"
	"post-api/story/models"
	"strings"
	"time"
)

type Draft struct {
	DraftID      uuid.UUID         `json:"draft_id" db:"id"`
	UserID       uuid.UUID         `json:"user_id" db:"user_id"`
	Data         models.JSONString `json:"data" db:"data"`
	PreviewImage *string           `json:"preview_image" db:"preview_image"`
	Tagline      *string           `json:"tagline" db:"tagline"`
	Interests    *string           `json:"-" db:"interests"`
	CreatedAt    *time.Time        `json:"created_at" db:"created_at"`
	InterestTags []Interests       `json:"interests"`
}

type DraftPreview struct {
	DraftID   uuid.UUID         `json:"id"`
	UserID    uuid.UUID         `json:"user_id"`
	Data      models.JSONString `json:"data"`
	Title     string            `json:"title"`
	Tagline   string            `json:"tagline"`
	Interests []Interests       `json:"interests"`
	CreatedAt *time.Time        `json:"created_at"`
}

func (draft *Draft) ConvertInterests(onSuccess func(interests []string) *golaerror.Error) *golaerror.Error {
	if draft.Interests != nil {
		length := len(*draft.Interests)
		sliced := (*draft.Interests)[:length-1]
		overAll := sliced[1:]
		interests := strings.Split(overAll, ",")
		return onSuccess(interests)
	}
	return nil
}
