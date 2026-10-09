package entity

import "time"

type OrderKTL struct {
	ID           string    `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()"`
	CustomerName string    `gorm:"type:varchar(100);not null"`
	PhoneNumber  string    `gorm:"type:varchar(20);not null"`
	VehicleType  string    `gorm:"type:varchar(20);not null;"` // 'Motor' or 'Mobil'
	Region       string    `gorm:"type:varchar(50);not null"`  // 'Jakarta', 'Tangerang', 'Depok'
	Fee                 float64   `gorm:"type:decimal(18,2);not null"`
	MerchantRef         *string   `gorm:"type:varchar(100)"` // Catatan atau referensi merchant
	Status              string    `gorm:"type:varchar(50);not null;default:'Pending'"`
	BuktiKtlPath        *string   `gorm:"type:varchar(255)"`
	BuktiPembayaranPath *string   `gorm:"type:varchar(255)"`
	OrderDate           time.Time `gorm:"autoCreateTime"`
	
	TrackingNumber      *string   `gorm:"type:varchar(100)"`
	CourierName         *string   `gorm:"type:varchar(50)"`
	CourierService      *string   `gorm:"type:varchar(100)"`
	ShippingCost        float64   `gorm:"type:decimal(18,2);default:0"`
}
