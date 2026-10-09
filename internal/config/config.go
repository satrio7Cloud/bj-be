package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type SQLServerConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type MinIOConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type QrisConfig struct {
	StaticPayload string
}

type MidtransConfig struct {
	MerchantID string
	ClientKey  string
	ServerKey  string
	Env        string
}

type BiteshipConfig struct {
	APIKey           string
	BaseURL          string
	DefaultCouriers  string
	OriginPostalCode string
}

type Config struct {
	AppHost          string
	AppPort          string
	NotificationCron string
	SQLServer        SQLServerConfig
	MinIO            MinIOConfig
	Redis            RedisConfig
	Qris             QrisConfig
	Midtrans         MidtransConfig
	Biteship         BiteshipConfig
}

func LoadConfig() *Config {
	// Mencoba memuat file .env dari beberapa kemungkinan lokasi (root atau dalam folder cmd)
	err := godotenv.Overload(".env")
	if err != nil {
		err = godotenv.Overload("../../.env")
	}
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan nilai bawaan OS")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "localhost" {
		dbHost = "127.0.0.1"
	}

	return &Config{
		AppHost:          os.Getenv("APP_HOST"),
		AppPort:          os.Getenv("APP_PORT"),
		NotificationCron: os.Getenv("CRON_SCHEDULE_NOTIFICATION"),
		SQLServer: SQLServerConfig{
			Host:     dbHost,
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
		},
		MinIO: MinIOConfig{
			Endpoint:   os.Getenv("MINIO_ENDPOINT"),
			AccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
			SecretKey:  os.Getenv("MINIO_SECRET_KEY"),
			BucketName: os.Getenv("MINIO_BUCKET_NAME"),
			UseSSL:     os.Getenv("MINIO_USE_SSL") == "true",
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		Qris: QrisConfig{
			StaticPayload: os.Getenv("QRIS_STATIC_PAYLOAD"),
		},
		Midtrans: MidtransConfig{
			MerchantID: os.Getenv("MIDTRANS_MERCHANT_ID"),
			ClientKey:  os.Getenv("MIDTRANS_CLIENT_KEY"),
			ServerKey:  os.Getenv("MIDTRANS_SERVER_KEY"),
			Env:        os.Getenv("MIDTRANS_ENV"),
		},
		Biteship: BiteshipConfig{
			APIKey:           os.Getenv("BITESHIP_API_KEY"),
			BaseURL:          os.Getenv("BITESHIP_BASE_URL"),
			DefaultCouriers:  os.Getenv("BITESHIP_DEFAULT_COURIERS"),
			OriginPostalCode: os.Getenv("BITESHIP_ORIGIN_POSTAL_CODE"),
		},
	}
}
