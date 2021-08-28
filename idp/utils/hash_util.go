package util

import (
	"golang.org/x/crypto/bcrypt"
)

type HashUtil interface {
	GenerateBcryptHash(text string) (string, error)
	MatchBcryptHash(hashedText, plainText string) error
}

type hashUtil struct {
}

func NewHashUtil() HashUtil {
	return hashUtil{}
}

func (hashUtil hashUtil) GenerateBcryptHash(text string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (hashUtil hashUtil) MatchBcryptHash(hashedText, plainText string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(plainText))
}
