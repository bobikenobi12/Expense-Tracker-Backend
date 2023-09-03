package config

import (
	"ExpenseTracker/models"
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
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" && os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		return errors.New("no AWS credentials found. Please set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY")
	}
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

func UploadToS3Bucket(file *multipart.File, filename string, fileType string) (*models.S3Object, error) {
	bucket := os.Getenv("S3_BUCKET")

	if bucket == "" {
		return nil, errors.New("no S3 bucket found. Please set S3_BUCKET")
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("profile_pics/" + filename),
		Body:   *file,
		ContentType: aws.String(
			fileType,
		),
	})

	if err != nil {
		return nil, err
	}

	return &models.S3Object{
		ETag:      result.ETag,
		VersionId: result.VersionID,
		Location:  result.Location,
		Key:       "profile_pics/" + filename,
	}, nil
}

func GetObjectFromS3Bucket(s3Object *models.S3Object) (*s3.GetObjectOutput, error) {
	bucket := os.Getenv("S3_BUCKET")

	if bucket == "" {
		return nil, errors.New("no S3 bucket found. Please set S3_BUCKET")
	}

	result, err := uploader.S3.GetObject(&s3.GetObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(s3Object.Key),
		VersionId: s3Object.VersionId,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
