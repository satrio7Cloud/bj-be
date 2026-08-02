package minio

import (
	"context"
	"log"

	"bj/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewMinioClient menginisialisasi koneksi ke server MinIO
func NewMinioClient(cfg *config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Cek apakah bucket sudah ada
	ctx := context.Background()
	bucketName := cfg.MinIO.BucketName
	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	if errBucketExists == nil && exists {
		log.Printf("MinIO: Bucket '%s' sudah tersedia", bucketName)
	} else {
		// Jika belum ada, buat bucket baru
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("MinIO: Gagal membuat bucket '%s': %v", bucketName, err)
			return nil, err
		}
		log.Printf("MinIO: Bucket '%s' berhasil dibuat", bucketName)
	}

	return minioClient, nil
}
