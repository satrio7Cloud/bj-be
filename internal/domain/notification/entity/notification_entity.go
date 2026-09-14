package entity

import "time"

type Notification struct {
	ID          string    `gorm:"column:ID;type:uniqueidentifier;primaryKey;default:NEWID()"`
	Title       string    `gorm:"column:Title;type:varchar(100);not null"`
	Message     string    `gorm:"column:Message;type:varchar(255);not null"`
	Type        string    `gorm:"column:Type;type:varchar(50);not null"` // 'warning', 'info', 'danger'
	ReferenceID *string   `gorm:"column:ReferenceID;type:varchar(36)"`          // Opsional UUID
	IsRead      bool      `gorm:"column:IsRead;type:bit;default:0"`
	CreatedAt   time.Time `gorm:"column:CreatedAt;autoCreateTime"`
}
