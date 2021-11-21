package request

import "github.com/google/uuid"

type InterestURIRequest struct {
	InterestUID string `uri:"interest_id" binding:"required,validPostUID"`
}

type InterestRequest struct {
	InterestUID uuid.UUID
	Start       int `form:"start"`
	Limit       int `form:"limit" binding:"required"`
}

type InterestNameRequest struct {
	Name string `json:"name" binding:"required"`
}
