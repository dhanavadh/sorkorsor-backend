package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dhanavadh/sorkorsor-backend/config"
)

type R2Storage struct {
	client     *s3.Client
	presign    *s3.PresignClient
	bucketName string
	publicURL  string
}

type PresignedURL struct {
	UploadURL string `json:"upload_url"`
	Key       string `json:"key"`
	PublicURL string `json:"public_url"`
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
		presign:    s3.NewPresignClient(client),
		bucketName: cfg.R2BucketName,
		publicURL:  cfg.R2PublicURL,
	}
}

func (r *R2Storage) GeneratePresignedURL(ctx context.Context, contentType string, ext string) (*PresignedURL, error) {
	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	key := fmt.Sprintf("images/%d_%s%s", time.Now().UnixNano(), hex.EncodeToString(randBytes), ext)

	presignResult, err := r.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return &PresignedURL{
		UploadURL: presignResult.URL,
		Key:       key,
		PublicURL: fmt.Sprintf("%s/%s", r.publicURL, key),
	}, nil
}

func (r *R2Storage) GenerateMultiplePresignedURLs(ctx context.Context, files []FileInfo) ([]PresignedURL, error) {
	results := make([]PresignedURL, 0, len(files))
	for _, file := range files {
		result, err := r.GeneratePresignedURL(ctx, file.ContentType, file.Extension)
		if err != nil {
			return nil, err
		}
		results = append(results, *result)
	}
	return results, nil
}

type FileInfo struct {
	ContentType string `json:"content_type"`
	Extension   string `json:"extension"`
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
	return fmt.Sprintf("%s/%s", r.publicURL, fileName), nil
}

func (r *R2Storage) UploadMultiple(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	urls := make([]string, 0, len(files))
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}

		ext := filepath.Ext(fileHeader.Filename)
		key := fmt.Sprintf("images/%d%s", time.Now().UnixNano(), ext)

		url, err := r.UploadFile(ctx, file, key, fileHeader.Header.Get("Content-Type"))
		file.Close()
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}
