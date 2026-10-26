package entity

import "time"

type CompanyProfile struct {
	ID               string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	OfficeName       string    `gorm:"type:varchar(100);not null" json:"office_name"`
	FullAddress      string    `gorm:"type:text;not null" json:"full_address"`
	PostCode         string    `gorm:"type:varchar(10);not null" json:"postal_code"`
	PhoneNumber      string    `gorm:"type:varchar(20);not null" json:"phone_number"`
	OperationalHours string    `gorm:"type:varchar(100);not null" json:"operational_hours"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
