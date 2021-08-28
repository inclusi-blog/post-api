package db

type RegistrationHash struct {
	ActivationHash string `db:"ACTIVATION_HASH"`
	UserID         string `db:"USER_ID"`
}
