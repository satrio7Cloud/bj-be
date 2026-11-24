package dto

import (
	"time"

	custEntity "bj/internal/domain/customer/entity"
	invEntity "bj/internal/domain/invoice/entity"
	orderEntity "bj/internal/domain/order/entity"
	srvEntity "bj/internal/domain/service/entity"
	vehEntity "bj/internal/domain/vehicle/entity"
	"bj/internal/utils"
)

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

// ToInvoiceDetailResponse merakit entity Invoice, Order, Customer, Vehicle, dan Service menjadi response DTO
func ToInvoiceDetailResponse(
	invoice *invEntity.Invoice,
	order *orderEntity.Order,
	cust *custEntity.Customer,
	veh *vehEntity.Vehicle,
	srv *srvEntity.Service,
) *InvoiceDetailResponse {
	if invoice == nil {
		return nil
	}

	resp := &InvoiceDetailResponse{
		InvoiceID:     invoice.ID,
		OrderID:       invoice.OrderID,
		InvoiceNumber: invoice.InvoiceNumber,
		IssueDate:     utils.ToWIB(invoice.IssueDate),
		DueDate:       utils.ToWIB(invoice.DueDate),
		PaymentStatus: invoice.PaymentStatus,
		PaymentMethod: invoice.PaymentMethod,
		PaidAt:        utils.ToWIBPtr(invoice.PaidAt),
		PaymentToken:  invoice.PaymentToken,
		PaymentUrl:    invoice.PaymentUrl,

		QrisPayload: invoice.QrisPayload,
		QrisAmount:  invoice.QrisAmount,
	}

	if order != nil {
		resp.TaxAmount = order.TaxAmount
		resp.ServiceFee = order.ServiceFee
		resp.AdminFee = order.AdminFee
		resp.PhysicalCheckFee = order.PhysicalCheckFee
		resp.DeliveryFee = order.DeliveryFee
		resp.ExpressFee = order.ExpressFee
		resp.TotalCost = order.TotalCost
	}

	if cust != nil {
		resp.CustomerName = cust.Fullname
		resp.CustomerPhone = cust.PhoneNumber
		resp.CustomerEmail = cust.Email
	}

	if veh != nil {
		resp.VehiclePlate = veh.PlateNumber
	}

	if srv != nil {
		resp.ServiceName = srv.ServiceName
	}

	return resp
}

// FromInvoiceWithDetails merakit entity InvoiceWithDetails menjadi response DTO
func FromInvoiceWithDetails(
	inv *invEntity.InvoiceWithDetails,
	paymentStatus string,
	paidAt *time.Time,
	paymentMethod *string,
) InvoiceDetailResponse {
	return InvoiceDetailResponse{
		InvoiceID:     inv.ID,
		OrderID:       inv.OrderID,
		InvoiceNumber: inv.InvoiceNumber,
		IssueDate:     utils.ToWIB(inv.IssueDate),
		DueDate:       utils.ToWIB(inv.DueDate),
		PaymentStatus: paymentStatus,
		PaymentMethod: paymentMethod,
		PaidAt:        utils.ToWIBPtr(paidAt),
		QrisPayload:   inv.QrisPayload,
		QrisAmount:    inv.QrisAmount,
		CustomerName:  inv.CustomerName,
		CustomerEmail: inv.CustomerEmail,
		CustomerPhone: inv.CustomerPhone,
		VehiclePlate:  inv.VehiclePlate,
		ServiceName:   inv.ServiceName,
		TotalCost:     inv.TotalCost,
	}
}
