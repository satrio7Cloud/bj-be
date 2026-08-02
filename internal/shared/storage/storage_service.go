package storage

import (
	"context"
	"mime/multipart"
	"time"
)

// StorageService adalah antarmuka untuk menyimpan dan mengambil file
type StorageService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folderName string) (string, error)
	GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error)
	DeleteFile(ctx context.Context, objectName string) error
}
