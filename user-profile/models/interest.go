package models

import (
	"github.com/google/uuid"
)

type Interest struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	CoverPic string    `json:"cover_pic"`
}

type ExploreInterests struct {
	Category  string     `json:"category"`
	Interests []Interest `json:"interests"`
}
