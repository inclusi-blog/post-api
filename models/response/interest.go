package response

type CategoryAndInterest struct {
	Category  string             `json:"category"`
	Interests []InterestWithIcon `json:"interests"`
}

type InterestWithIcon struct {
	Name             string `json:"name"`
	Image            string `json:"image"`
	IsFollowedByUser bool   `json:"isFollowedByUser"`
}
