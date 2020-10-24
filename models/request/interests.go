package request

type SearchInterests struct {
	SearchKeyword string   `json:"searchKeyword,omitempty"`
	SelectedTags  []string `json:"selectedTags,omitempty" binding:"required"`
}
