package response

type InterestCountDetails struct {
	FollowersCount int64 `json:"followers_count" db:"followers_count"`
	IsFollowed     bool  `json:"is_followed" db:"is_followed"`
}
