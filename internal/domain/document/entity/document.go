package entity

type Document struct {
	ID           string `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()"`
	OrderID      string `gorm:"type:uniqueidentifier;not null"`
	DocumentType string `gorm:"type:varchar(50);not null"`
	FilePath     string `gorm:"type:varchar(255);not null"`
	IsVerified   bool   `gorm:"type:bit;default:0"`
}
