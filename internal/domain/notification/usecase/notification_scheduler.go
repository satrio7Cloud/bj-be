package usecase

import (
	"fmt"
	"log"

	"bj/internal/domain/notification/entity"

	"github.com/robfig/cron/v3"
)

func (u *notificationUsecase) StartNotificationScheduler() {
	// Buat scheduler cron baru
	c := cron.New()

	// Jalankan setiap jam 8 pagi ("0 8 * * *")
	// _, err := c.AddFunc("0 8 * * *", func() {
	_, err := c.AddFunc("*/1 * * * *", func() {
		log.Println("[CRON] Menjalankan pengecekan notifikasi harian...")
		u.scanAndNotify()
	})

	if err != nil {
		log.Printf("[CRON ERROR] Gagal mendaftarkan cron: %v\n", err)
		return
	}

	c.Start()
	log.Println("[CRON] Scheduler notifikasi harian telah berjalan.")
}

func (u *notificationUsecase) scanAndNotify() {
	// 1. Unverified Documents
	unverifiedDocs, err := u.repo.FindUnverifiedDocuments()
	if err == nil {
		for _, doc := range unverifiedDocs {
			id := doc["ID"].(string)
			name := doc["CustomerName"].(string)
			title := "Dokumen Belum Diverifikasi"

			if !u.repo.NotificationExists(id, title) {
				_ = u.repo.Create(&entity.Notification{
					Title:       title,
					Message:     fmt.Sprintf("Dokumen untuk order atas nama %s butuh verifikasi.", name),
					Type:        "warning",
					ReferenceID: &id,
				})
			}
		}
	}

	// 2. Stuck Orders
	stuckOrders, err := u.repo.FindStuckOrders()
	if err == nil {
		for _, order := range stuckOrders {
			id := order["ID"].(string)
			name := order["CustomerName"].(string)
			title := "Order Reguler Nyangkut"

			if !u.repo.NotificationExists(id, title) {
				_ = u.repo.Create(&entity.Notification{
					Title:       title,
					Message:     fmt.Sprintf("Order atas nama %s sudah tertahan > 3 hari.", name),
					Type:        "danger",
					ReferenceID: &id,
				})
			}
		}
	}

	// 3. Stuck KTLs
	stuckKTLs, err := u.repo.FindStuckKTLs()
	if err == nil {
		for _, ktl := range stuckKTLs {
			id := ktl["ID"].(string)
			name := ktl["CustomerName"].(string)
			title := "Order KTL Nyangkut"

			if !u.repo.NotificationExists(id, title) {
				_ = u.repo.Create(&entity.Notification{
					Title:       title,
					Message:     fmt.Sprintf("Order KTL atas nama %s sudah tertahan > 3 hari.", name),
					Type:        "danger",
					ReferenceID: &id,
				})
			}
		}
	}
}
