package entity

import "time"

type Customer struct {
	ID          string    `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()"`
	Fullname    string    `gorm:"type:varchar(100);not null"`
	PhoneNumber string    `gorm:"type:varchar(20);not null"`
	Address     string    `gorm:"type:text"`
	NIK         string    `gorm:"type:varchar(16);unique;not null"`
	Email       string    `gorm:"type:varchar(100)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
