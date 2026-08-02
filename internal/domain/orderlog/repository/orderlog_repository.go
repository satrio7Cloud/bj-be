package repository

import (
	"bj/internal/domain/orderlog/entity"

	"gorm.io/gorm"
)

type OrderLogRepository interface {
	Create(log *entity.OrderLog) error
	FindByOrderID(orderID string) ([]entity.OrderLog, error)
}

type orderLogRepository struct {
	db *gorm.DB
}

func NewOrderLogRepository(db *gorm.DB) OrderLogRepository {
	return &orderLogRepository{db}
}

func (r *orderLogRepository) Create(log *entity.OrderLog) error {
	query := `INSERT INTO dbo.OrderLogs (ID, OrderID, AdminID, OldStatus, NewStatus, Notes, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, log.ID, log.OrderID, log.AdminID, log.OldStatus, log.NewStatus, log.Notes, log.CreatedAt).Error
}

func (r *orderLogRepository) FindByOrderID(orderID string) ([]entity.OrderLog, error) {
	var logs []entity.OrderLog
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, CAST(OrderID AS VARCHAR(36)) AS OrderID, CAST(AdminID AS VARCHAR(36)) AS AdminID, OldStatus, NewStatus, Notes, CreatedAt FROM dbo.OrderLogs WHERE OrderID = ? ORDER BY CreatedAt ASC`
	err := r.db.Raw(query, orderID).Scan(&logs).Error
	return logs, err
}
