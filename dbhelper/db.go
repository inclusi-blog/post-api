package dbhelper

import (
	"fmt"
	"os"
)

func BuildConnectionString() string {
	dbConnectionString := fmt.Sprintf("%s://%s:%s",
		os.Getenv("CONNECTION_TYPE"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))
	return dbConnectionString
}
