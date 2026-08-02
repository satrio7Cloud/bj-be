package entity

import "time"

type Service struct {
	ID          string    `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()"`
	ServiceName string    `gorm:"type:varchar(100);not null"`
	Category    string    `gorm:"type:varchar(50);not null"`
	ServiceFee  float64   `gorm:"type:decimal(18,2);not null;default:0"`
	ExpressFee  float64   `gorm:"type:decimal(18,2);not null;default:0"`
	Description string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
