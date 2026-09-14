package dto

import "time"

type CreateOrderKTLRequest struct {
	CustomerName string  `json:"customer_name" form:"customer_name" validate:"required"`
	PhoneNumber  string  `json:"phone_number" form:"phone_number" validate:"required"`
	VehicleType  string  `json:"vehicle_type" form:"vehicle_type" validate:"required"`
	Region       string  `json:"region" form:"region" validate:"required"`
	Fee          float64 `json:"fee" form:"fee" validate:"required,gt=0"`
	MerchantRef  string  `json:"merchant_ref,omitempty" form:"merchant_ref"`
}

type OrderKTLResponse struct {
	ID                  string    `json:"id"`
	CustomerName        string    `json:"customer_name"`
	PhoneNumber         string    `json:"phone_number"`
	VehicleType         string    `json:"vehicle_type"`
	Region              string    `json:"region"`
	Fee                 float64   `json:"fee"`
	MerchantRef         *string   `json:"merchant_ref,omitempty"`
	Status              string    `json:"status"`
	BuktiKtlPath        *string   `json:"bukti_ktl_path,omitempty"`
	BuktiPembayaranPath *string   `json:"bukti_pembayaran_path,omitempty"`
	OrderDate           time.Time `json:"order_date"`
}
