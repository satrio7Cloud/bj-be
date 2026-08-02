package entity

type Brand struct {
	ID        string `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()"`
	BrandName string `gorm:"type:varchar(50);not null"`
}
