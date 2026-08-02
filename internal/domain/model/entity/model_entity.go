package entity

type Model struct {
	ID          string `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()"`
	BrandID     string `gorm:"type:uniqueidentifier;not null"`
	BrandName   string `gorm:"column:BrandName"` // Added from JOIN
	ModelName   string `gorm:"type:varchar(100);not null"`
	VehicleType string `gorm:"type:varchar(50);not null"`
}
