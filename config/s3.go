package config

import (
	"errors"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var uploader *s3manager.Uploader

func S3() error {
	credentials := aws.NewConfig().Credentials
	region := os.Getenv("S3_REGION")

	if region == "" {
		return errors.New("no AWS region found. Please set S3_REGION")
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials,
	}))

	s3Svc := s3.New(sess)

	uploader = s3manager.NewUploaderWithClient(s3Svc, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // 5MB per part
		u.Concurrency = 10
	})

	return nil
}

func UploadToS3Bucket(file *multipart.File) (*s3manager.UploadOutput, error) {
	bucket := os.Getenv("S3_BUCKET")
	key := os.Getenv("S3_KEY")

	if bucket == "" {
		return nil, errors.New("no S3 bucket found. Please set S3_BUCKET")
	}

	if key == "" {
		return nil, errors.New("no S3 key found. Please set S3_KEY")
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   *file,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
