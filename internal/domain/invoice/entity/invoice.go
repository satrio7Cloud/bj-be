package entity

import "time"

type Invoice struct {
	ID            string     `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()"`
	OrderID       string     `gorm:"type:uniqueidentifier;not null;uniqueIndex"`
	InvoiceNumber string     `gorm:"type:varchar(50);not null;uniqueIndex"`
	IssueDate     time.Time  `gorm:"autoCreateTime"`
	DueDate       time.Time  `gorm:"not null"`
	PaymentStatus string     `gorm:"type:varchar(20);not null;default:'UNPAID'"`
	PaymentMethod *string    `gorm:"type:varchar(50)"`
	PaidAt        *time.Time `gorm:"type:datetime"`
	QrisPayload   *string    `gorm:"type:varchar(500)"`
	QrisAmount    *float64   `gorm:"type:decimal(18,2)"`
	PaymentToken  *string    `gorm:"type:varchar(255)"`
	PaymentUrl    *string    `gorm:"type:varchar(255)"`
}

type InvoiceWithDetails struct {
	Invoice
	CustomerName  string  `gorm:"column:CustomerName"`
	CustomerEmail string  `gorm:"column:CustomerEmail"`
	CustomerPhone string  `gorm:"column:CustomerPhone"`
	VehiclePlate  string  `gorm:"column:VehiclePlate"`
	ServiceName   string  `gorm:"column:ServiceName"`
	TotalCost     float64 `gorm:"column:TotalCost"`
}
