package db

type SaveUserDetails struct {
	UUID     string `db:"UUID"`
	Username string `db:"USERNAME"`
	Email    string `db:"EMAIL"`
	Password string `db:"PASSWD"`
	IsActive bool   `db:"IS_ACTIVE"`
}
