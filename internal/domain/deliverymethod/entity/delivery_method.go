package entity

import "time"

// DeliveryMethod merepresentasikan tabel dbo.MasterShippingMethods.
// Mengelola metode penyerahan (PICKUP) dan pengembalian (RETURN) berkas biro jasa.
type DeliveryMethod struct {
	ID              string    `gorm:"type:varchar(36);primaryKey"`
	MethodCode      string    `gorm:"type:varchar(50);unique;not null"`
	MethodName      string    `gorm:"type:varchar(100);not null"`
	Direction       string    `gorm:"type:varchar(20);not null"` // PICKUP | RETURN | BOTH
	ProviderType    string    `gorm:"type:varchar(30);not null"` // SELF | COURIER | INTERNAL
	Description     string    `gorm:"type:nvarchar(255)"`
	FlatFee         int64     `gorm:"type:bigint;not null;default:0"`
	RequiresAddress bool      `gorm:"type:bit;not null;default:0"`
	RequiresCourier bool      `gorm:"type:bit;not null;default:0"`
	IsActive        bool      `gorm:"type:bit;not null;default:1"`
	DisplayOrder    int       `gorm:"type:int;not null;default:1"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

