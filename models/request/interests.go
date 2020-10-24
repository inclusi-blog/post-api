package request

type SearchInterests struct {
	SearchKeyword string `json:"search_keyword,omitempty" binding:"required"`
}
