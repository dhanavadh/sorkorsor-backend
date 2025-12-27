package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dhanavadh/sorkorsor-backend/config"
)

type R2Storage struct {
	client     *s3.Client
	bucketName string
	accountID  string
}

func NewR2Storage(cfg *config.Config) *R2Storage {
	client := s3.New(s3.Options{
		Region: "auto",
		BaseEndpoint: aws.String(
			fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID),
		),
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.R2AccessKeyID,
			cfg.R2SecretAccessKey,
			"",
		),
	})
	return &R2Storage{
		client:     client,
		bucketName: cfg.R2BucketName,
		accountID:  cfg.R2AccountID,
	}
}

func (r *R2Storage) UploadFile(ctx context.Context, file multipart.File, fileName string, contentType string) (string, error) {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	url := fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s", r.accountID, fileName)
	return url, nil
}
