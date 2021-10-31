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

type InterestNameURIRequest struct {
	Name string `uri:"name" binding:"required"`
}
