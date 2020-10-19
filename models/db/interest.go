package db

type Interest struct {
	Name string `json:"name" db:"NAME"`
	ID   string `json:"id" db:"INTEREST_ID"`
}
