package data

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/gommon/log"
)

type S3 struct {
	client *s3.Client
}

func NewS3(client *s3.Client) *S3 {
	return &S3{client: client}
}

func (s *S3) UploadFile(ctx context.Context, bucket, key string, file *os.File) error {
	putObjectOutput, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return err
	}

	log.Info(putObjectOutput.ResultMetadata)

	return nil
}

func (s *S3) DownloadFile(ctx context.Context, bucket, key, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3) DeleteFile(ctx context.Context, bucket, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}
