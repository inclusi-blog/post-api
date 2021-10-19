package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"post-api/configuration"
	"time"
)

type awsServices struct {
	session *session.Session
	config  *configuration.ConfigData
}

type AwsServices interface {
	GetObjectInS3(key string, expiryTime time.Duration) (string, error)
	PutObjectInS3(key string) (string, error)
	CheckS3Object(key string) (bool, error)
}

func (service awsServices) GetObjectInS3(key string, expiryTime time.Duration) (string, error) {
	svc := s3.New(service.session)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(service.config.AwsBucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(expiryTime)

	if err != nil {
		return "", nil
	}

	return urlStr, err
}

func (service awsServices) PutObjectInS3(key string) (string, error) {
	svc := s3.New(service.session)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(service.config.AwsBucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(30 * time.Minute)

	if err != nil {
		return "", err
	}

	return urlStr, nil
}

func (service awsServices) CheckS3Object(key string) (bool, error) {
	svc := s3.New(service.session)
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(service.config.AwsBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}

func NewAwsServices(session *session.Session, data *configuration.ConfigData) AwsServices {
	return awsServices{
		session: session,
		config:  data,
	}
}
