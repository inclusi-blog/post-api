package models

import "post-api/story/models"

type ExploreInterest struct {
	Category  string            `json:"category" db:"category"`
	Interests models.JSONString `json:"interests" db:"interests"`
}
