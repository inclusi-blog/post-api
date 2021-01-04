package db

import (
	"context"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/constants"
	"post-api/models"
)

type Draft struct {
	DraftID      string            `json:"draft_id" db:"draft_id"`
	PostData     models.JSONString `json:"post_data" db:"post_data"`
	PreviewImage *string           `json:"preview_image" db:"preview_image"`
	Tagline      *string           `json:"tagline" db:"tagline"`
	Interest     []string          `json:"interest" db:"interest"`
}

func (draft Draft) IsValidInterest(ctx context.Context, config map[string]int, readTime int, minimumReadTime int) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "PostValidator").WithField("method", "IsValidInterest")
	if !draft.isValidInterestCount() {
		return &constants.MinimumInterestCountNotMatchErr
	}
	for _, interest := range draft.Interest {
		if interest == "" {
			return &constants.DraftInterestParseError
		}
		configReadTime := config[interest]
		if configReadTime != 0 {
			if readTime < configReadTime {
				logger.Errorf("post interest doesn't meet required read time %v .%v", draft.DraftID, readTime)
				return &constants.InterestReadTimeDoesNotMeetErr
			}
			continue
		}
		if readTime < minimumReadTime {
			logger.Errorf("post doesn't meet minimum read time %v .%v", draft.DraftID, readTime)
			return &constants.ReadTimeNotMeetError
		}
	}
	return nil
}

func (draft Draft) isValidInterestCount() bool {
	return len(draft.Interest) >= 3 && len(draft.Interest) <= 5
}

type DraftDB struct {
	DraftID      string            `json:"draft_id" db:"draftId"`
	UserID       string            `json:"user_id" db:"userId"`
	PostData     models.JSONString `json:"post_data" db:"postData"`
	PreviewImage string            `json:"preview_image" db:"previewImage"`
	Tagline      string            `json:"tagline" db:"tagline"`
	Interest     []string          `json:"interest" db:"interests"`
	IsPublished  bool              `json:"is_published"`
	CreatedAt    int64             `json:"created_at"`
}

type AllDraft struct {
	DraftID    string   `json:"draft_id" db:"draft_id"`
	TitleData  string   `json:"title_data" db:"title_data"`
	Tagline    *string  `json:"tagline" db:"tagline"`
	Interest   []string `json:"interest" db:"interest"`
	CreatedAt  int64    `json:"created_at"`
	AuthorName string   `json:"author_name"`
}
