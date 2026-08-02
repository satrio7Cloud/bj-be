package entity

type Vehicle struct {
	ID          string `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()"`
	CustomerID  string `gorm:"type:uniqueidentifier;not null"`
	ModelID     string `gorm:"type:uniqueidentifier;not null"`
	PlateNumber string `gorm:"type:varchar(20);unique;not null"`
	Year        int    `gorm:"type:int;not null"`
}
