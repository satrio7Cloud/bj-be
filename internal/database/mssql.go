package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"bj/internal/config"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectSQLServer(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&connection+timeout=3",
		cfg.SQLServer.User,
		cfg.SQLServer.Password,
		cfg.SQLServer.Host,
		cfg.SQLServer.Port,
		cfg.SQLServer.Database,
	)

	// Setup GORM Logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Toleransi query lambat
			LogLevel:                  logger.Info, // Tampilkan semua query (ubah ke Warn di production)
			IgnoreRecordNotFoundError: true,        // Abaikan error NotFound
			Colorful:                  true,        // Teks berwarna
		},
	)

	var db *gorm.DB
	var err error
	maxRetries := 5

	// Retry loop untuk mengantisipasi database booting atau mati sementara
	for i := 1; i <= maxRetries; i++ {
		log.Printf("Mencoba terhubung ke database (Percobaan %d/%d)...\n", i, maxRetries)
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err == nil {
			sqlDB, errPool := db.DB()
			if errPool == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					sqlDB.SetMaxIdleConns(10)
					sqlDB.SetMaxOpenConns(100)
					sqlDB.SetConnMaxLifetime(time.Hour)
					log.Println("Berhasil terhubung ke database", cfg.SQLServer.Database)
					return db
				}
			}
		}

		log.Printf("Gagal terhubung ke database (%v). Mencoba kembali dalam 3 detik...\n", err)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("Gagal terhubung ke database setelah %d percobaan: %v", maxRetries, err)
	return nil
}
