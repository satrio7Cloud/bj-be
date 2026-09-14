package repository

import (
	"fmt"
	"bj/internal/domain/notification/entity"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *entity.Notification) error
	GetUnreadNotifs(limit int) ([]entity.Notification, error)
	MarkAsRead(id string) error

	// Scanner Methods for Cronjob
	NotificationExists(refID string, title string) bool
	FindUnverifiedDocuments() ([]map[string]interface{}, error)
	FindStuckOrders() ([]map[string]interface{}, error)
	FindStuckKTLs() ([]map[string]interface{}, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *entity.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) GetUnreadNotifs(limit int) ([]entity.Notification, error) {
	var notifs []entity.Notification
	
	query := `
		SELECT 
			CAST(ID AS VARCHAR(36)) AS ID, 
			Title, 
			Message, 
			Type, 
			CAST(ReferenceID AS VARCHAR(36)) AS ReferenceID,
			IsRead, 
			CreatedAt 
		FROM dbo.Notifications 
		WHERE IsRead = 0 
		ORDER BY CreatedAt DESC
	`
	
	if limit > 0 {
		query += fmt.Sprintf(" OFFSET 0 ROWS FETCH NEXT %d ROWS ONLY", limit)
	}

	err := r.db.Raw(query).Scan(&notifs).Error
	return notifs, err
}

func (r *notificationRepository) MarkAsRead(id string) error {
	return r.db.Model(&entity.Notification{}).Where("ID = ?", id).Update("IsRead", 1).Error
}

func (r *notificationRepository) NotificationExists(refID string, title string) bool {
	var count int64
	r.db.Model(&entity.Notification{}).Where("ReferenceID = ? AND Title = ?", refID, title).Count(&count)
	return count > 0
}

func (r *notificationRepository) FindUnverifiedDocuments() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `
		SELECT CAST(D.ID AS VARCHAR(36)) AS ID, C.FullName AS CustomerName 
		FROM dbo.OrderDocuments D
		JOIN dbo.Orders O ON D.OrderID = O.ID
		JOIN dbo.Customers C ON O.CustomerID = C.ID
		WHERE D.IsVerified = 0 AND D.DeletedAt IS NULL
	`
	err := r.db.Raw(query).Scan(&results).Error
	return results, err
}

func (r *notificationRepository) FindStuckOrders() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `
		SELECT CAST(O.ID AS VARCHAR(36)) AS ID, C.FullName AS CustomerName 
		FROM dbo.Orders O
		JOIN dbo.Customers C ON O.CustomerID = C.ID
		WHERE O.Status = 'Processing' AND DATEDIFF(day, O.OrderDate, GETDATE()) >= 3
	`
	err := r.db.Raw(query).Scan(&results).Error
	return results, err
}

func (r *notificationRepository) FindStuckKTLs() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `
		SELECT CAST(ID AS VARCHAR(36)) AS ID, CustomerName 
		FROM dbo.OrderKTLs
		WHERE Status = 'Processing' AND DATEDIFF(day, OrderDate, GETDATE()) >= 3
	`
	err := r.db.Raw(query).Scan(&results).Error
	return results, err
}

