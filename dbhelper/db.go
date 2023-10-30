package dbhelper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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

func BuildConnectionStringCloud(svc *session.Session) (string, error) {
	manager := secretsmanager.New(svc, nil)
	value, err := manager.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String("dev/postgres/password"),
	})
	if err != nil {
		log.Println("unable to get db password from aws secrets ", err)
		return "", err
	}
	secretString := []byte(*value.SecretString)
	var data map[string]string
	err = json.Unmarshal(secretString, &data)
	if err != nil {
		log.Println("unable to unmarshal db password from aws secrets ", err)
		return "", err
	}
	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		data["STORY_PASSWORD"],
		os.Getenv("DB_SERVICE_NAME"), "sslmode=disable")
	return dbConnectionString, nil
}
