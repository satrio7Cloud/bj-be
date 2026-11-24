package dto

import (
	"time"

	"bj/internal/domain/order/entity"
	"bj/internal/utils"
)

type CreateOrderKTLRequest struct {
	CustomerName string  `json:"customer_name" form:"customer_name" validate:"required"`
	PhoneNumber  string  `json:"phone_number" form:"phone_number" validate:"required"`
	VehicleType  string  `json:"vehicle_type" form:"vehicle_type" validate:"required"`
	Region       string  `json:"region" form:"region" validate:"required"`
	Fee          float64 `json:"fee" form:"fee" validate:"required,gt=0"`
}

func (req *CreateOrderKTLRequest) ToEntity(id string, buktiKtlPath *string) *entity.OrderKTL {
	return &entity.OrderKTL{
		ID:                  id,
		CustomerName:        req.CustomerName,
		PhoneNumber:         req.PhoneNumber,
		VehicleType:         req.VehicleType,
		Region:              req.Region,
		Fee:                 req.Fee,
		MerchantRef:         nil,
		Status:              "Menunggu Pembayaran",
		BuktiKtlPath:        buktiKtlPath,
		BuktiPembayaranPath: nil,
		OrderDate:           time.Now(),
	}
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
	SnapToken           string    `json:"snap_token,omitempty"`
	OrderDate           time.Time `json:"order_date"`
}

// MidtransWebhookPayload is the expected payload from Midtrans
type MidtransWebhookPayload struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

func ToOrderKTLResponse(o entity.OrderKTL, snapToken ...string) OrderKTLResponse {
	token := ""
	if len(snapToken) > 0 {
		token = snapToken[0]
	}
	return OrderKTLResponse{
		ID:                  o.ID,
		CustomerName:        o.CustomerName,
		PhoneNumber:         o.PhoneNumber,
		VehicleType:         o.VehicleType,
		Region:              o.Region,
		Fee:                 o.Fee,
		MerchantRef:         o.MerchantRef,
		Status:              o.Status,
		BuktiKtlPath:        o.BuktiKtlPath,
		BuktiPembayaranPath: o.BuktiPembayaranPath,
		SnapToken:           token,
		OrderDate:           utils.ToWIB(o.OrderDate),
	}
}

func ToOrderKTLResponses(orders []entity.OrderKTL) []OrderKTLResponse {
	responses := make([]OrderKTLResponse, 0, len(orders))
	for _, o := range orders {
		responses = append(responses, ToOrderKTLResponse(o))
	}
	return responses
}
