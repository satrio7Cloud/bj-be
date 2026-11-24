package entity

import "time"

// Courier merepresentasikan master kurir & layanan di tabel dbo.ShippingCouriers.
// Berisi hasil sinkronisasi dari Biteship dan bertindak sebagai cache lokal di database.
type Courier struct {
	ID                 string    `gorm:"type:varchar(36);primaryKey"`
	CourierCode        string    `gorm:"type:varchar(50);not null"`
	CourierName        string    `gorm:"type:varchar(100);not null"`
	CourierServiceCode string    `gorm:"type:varchar(50);not null"`
	CourierServiceName string    `gorm:"type:varchar(100);not null"`
	Tier               string    `gorm:"type:varchar(50)"`
	ServiceType        string    `gorm:"type:varchar(50)"`
	ShippingType       string    `gorm:"type:varchar(50)"`
	DurationRange      string    `gorm:"type:varchar(50)"`
	DurationUnit       string    `gorm:"type:varchar(20)"`
	IsAvailable        bool      `gorm:"type:bit;not null;default:1"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

