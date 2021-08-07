package dbhelper

import (
	"fmt"
	"os"
)

func BuildConnectionString() string {
	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SERVICE_NAME"), "sslmode=disable")
	return dbConnectionString
}
