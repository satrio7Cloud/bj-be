package entity

import "time"

// ShippingSettings adalah konfigurasi pengaturan kurir oleh Admin (singleton record).
// Menyimpan preferensi teknis paket (berat, asuransi) dan kurir rekanan aktif.
type ShippingSettings struct {
	ID                string    `gorm:"type:varchar(36);primaryKey"`
	EnabledCouriers   string    `gorm:"type:nvarchar(max)"`
	DefaultItemWeight int       `gorm:"type:int;not null;default:200"`
	DefaultItemValue  int64     `gorm:"type:bigint;not null;default:100000"`
	UpdatedBy         string    `gorm:"type:varchar(36)"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}
