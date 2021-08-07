package db

type Interest struct {
	Name string `json:"name" db:"name"`
	ID   string `json:"id" db:"id"`
}
