package entity

import "time"

type Order struct {
	ID            string  `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()"`
	CustomerID    string  `gorm:"type:uniqueidentifier;not null"`
	VehicleID     string  `gorm:"type:uniqueidentifier;not null"`
	ServiceID     string  `gorm:"type:uniqueidentifier;not null"`
	ProcessedByID *string `gorm:"type:uniqueidentifier"` // Boleh null jika belum diproses
	Status        string  `gorm:"type:varchar(50);not null;default:'Pending'"`

	// --- Delivery Options ---
	PickupMethod           string  `gorm:"type:varchar(50)"`
	ReturnMethod           string  `gorm:"type:varchar(50)"`
	CustomerTrackingNumber *string `gorm:"type:varchar(100)"`
	BiroJasaTrackingNumber *string `gorm:"type:varchar(100)"`
	PickupAddress          *string `gorm:"type:varchar(max)"`

	// --- Fees & Pricing ---
	TaxAmount        float64 `gorm:"type:decimal(18,2);not null;default:0"`
	ServiceFee       float64 `gorm:"type:decimal(18,2);not null;default:0"`
	AdminFee         float64 `gorm:"type:decimal(18,2);not null;default:0"`
	PhysicalCheckFee float64 `gorm:"type:decimal(18,2);not null;default:0"`
	DeliveryFee      float64 `gorm:"type:decimal(18,2);not null;default:0"`
	ExpressFee       float64 `gorm:"type:decimal(18,2);not null;default:0"`
	TotalCost        float64 `gorm:"type:decimal(18,2);not null;default:0"`

	// --- Service Level & Flags ---
	ServiceLevel           string `gorm:"type:varchar(50);not null;default:'STANDARD'"`
	IsNameTransferRequired bool   `gorm:"type:bit;not null;default:0"`

	SamsatOrigin      *string `gorm:"type:varchar(100)"`
	SamsatDestination *string `gorm:"type:varchar(100)"`

	Notes     string    `gorm:"type:varchar(max)"`
	OrderDate time.Time `gorm:"autoCreateTime"`

	TrackingNumber      *string   `gorm:"type:varchar(100)"`
	CourierName         *string   `gorm:"type:varchar(50)"`
	CourierService      *string   `gorm:"type:varchar(100)"`
	ShippingCost        float64   `gorm:"type:decimal(18,2);default:0"`
}
