package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"time"

	"bj/internal/config"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type minioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(client *minio.Client, cfg *config.Config) StorageService {
	return &minioStorage{
		client:     client,
		bucketName: cfg.MinIO.BucketName,
	}
}

func (s *minioStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, folderName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %w", err)
	}
	defer src.Close()

	// Generate nama file unik dengan mempertahankan nama aslinya
	// Format: folder/uuid-namaasli.ext
	originalName := filepath.Base(file.Filename)
	objectName := fmt.Sprintf("%s/%s-%s", folderName, uuid.New().String(), originalName)

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = s.client.PutObject(ctx, s.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("gagal mengunggah file ke MinIO: %w", err)
	}

	return objectName, nil
}

func (s *minioStorage) GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, expires, reqParams)
	if err != nil {
		return "", fmt.Errorf("gagal membuat URL file: %w", err)
	}
	return presignedURL.String(), nil
}

func (s *minioStorage) DeleteFile(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("gagal menghapus file dari MinIO: %w", err)
	}
	return nil
}
