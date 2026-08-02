package entity

import "time"

type OrderLog struct {
	ID        string    `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()"`
	OrderID   string    `gorm:"type:uniqueidentifier;not null"`
	AdminID   *string   `gorm:"type:uniqueidentifier;null"`
	OldStatus string    `gorm:"type:varchar(50)"`
	NewStatus string    `gorm:"type:varchar(50);not null"`
	Notes     string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
