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

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	log.Println("Berhasil terhubung ke database", cfg.SQLServer.Database)
	return db
}
