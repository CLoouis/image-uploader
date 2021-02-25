package s3

import (
	"fmt"
	"time"

	"github.com/CLoouis/image-uploader/pkg/utl/uploader"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type (
	S3Uploader struct {
		awsService *s3.S3
		bucket     string
	}
)

func NewS3Uploader(id, region, key, bucket, token string) uploader.Uploader {
	creds := credentials.NewStaticCredentials(id, key, token)
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
	}

	cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)

	s3Session, _ := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})

	awsService := s3.New(s3Session, cfg)
	return &S3Uploader{awsService: awsService, bucket: bucket}
}

func (s S3Uploader) GetPresignUploadUrl(fileName string) (string, error) {
	request, _ := s.awsService.PutObjectRequest(&s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(fileName),
		ContentLength: aws.Int64(5242880),
	})

	url, err := request.Presign(3 * time.Minute)

	if err != nil {
		return "", err
	}

	return url, nil
}

func (s S3Uploader) GetPresignFetchUrl(fileName string) (string, error) {
	request, _ := s.awsService.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	})

	url, err := request.Presign(1 * time.Minute)

	if err != nil {
		return "", err
	}

	return url, nil
}
