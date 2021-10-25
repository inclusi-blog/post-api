package request

type InterestURIRequest struct {
	InterestUID string `uri:"interest_id" binding:"required,validPostUID"`
}
