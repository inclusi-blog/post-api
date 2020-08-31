package dbhelper

import (
	"fmt"
	"os"
)

func BuildConnectionString() string {
	dbConnectionString := fmt.Sprintf("%s/%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SERVICE_NAME"),"parseTime","true")
	return dbConnectionString
}
