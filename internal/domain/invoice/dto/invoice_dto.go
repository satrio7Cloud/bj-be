package dto

import "time"

type GenerateInvoiceRequest struct {
	OrderID string `json:"order_id" validate:"required"`
}

type PayInvoiceRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required"`
}

type InvoiceDetailResponse struct {
	InvoiceID        string     `json:"invoice_id"`
	OrderID          string     `json:"order_id"`
	InvoiceNumber    string     `json:"invoice_number"`
	IssueDate        time.Time  `json:"issue_date"`
	DueDate          time.Time  `json:"due_date"`
	PaymentStatus    string     `json:"payment_status"`
	PaymentMethod    *string    `json:"payment_method,omitempty"`
	PaidAt           *time.Time `json:"paid_at,omitempty"`
	QrisPayload      *string    `json:"qris_payload,omitempty"`
	QrisAmount       *float64   `json:"qris_amount,omitempty"`
	PaymentToken     *string    `json:"payment_token,omitempty"`
	PaymentUrl       *string    `json:"payment_url,omitempty"`
	
	TaxAmount        float64    `json:"tax_amount"`
	ServiceFee       float64    `json:"service_fee"`
	AdminFee         float64    `json:"admin_fee"`
	PhysicalCheckFee float64    `json:"physical_check_fee"`
	DeliveryFee      float64    `json:"delivery_fee"`
	ExpressFee       float64    `json:"express_fee"`
	TotalCost        float64    `json:"total_cost"`
	
	CustomerName     string     `json:"customer_name,omitempty"`
	CustomerEmail    string     `json:"customer_email,omitempty"`
	CustomerPhone    string     `json:"customer_phone,omitempty"`
	VehiclePlate     string     `json:"vehicle_plate,omitempty"`
	ServiceName      string     `json:"service_name,omitempty"`
}
