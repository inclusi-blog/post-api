package db

type Interest struct {
	Name string `json:"name" db:"name"`
}

type CategoryAndInterest struct {
	Category  string             `json:"category" db:"category"`
	Interests []InterestWithIcon `json:"interests" db:"interests"`
}

type InterestWithIcon struct {
	Name  string `json:"name" db:"name"`
	Image string `json:"image" db:"image"`
}
