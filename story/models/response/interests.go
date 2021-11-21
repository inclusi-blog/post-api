package response

import "github.com/google/uuid"

type InterestCountDetails struct {
	FollowersCount int64     `json:"followers_count" db:"followers_count"`
	IsFollowed     bool      `json:"is_followed" db:"is_followed"`
	InterestID     uuid.UUID `json:"interest_id" db:"interest_id"`
	Name           string    `json:"name" db:"name"`
}
