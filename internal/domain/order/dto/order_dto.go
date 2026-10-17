package dto

import "time"

type CreateOrderAdminRequest struct {
	CustomerID string `json:"customer_id" validate:"required,uuid"`
	VehicleID  string `json:"vehicle_id" validate:"required,uuid"`
	ServiceID              string  `json:"service_id" validate:"required,uuid"`
	
	// --- Delivery Options ---
	PickupMethod           string  `json:"pickup_method"`
	ReturnMethod           string  `json:"return_method"`
	CustomerTrackingNumber *string `json:"customer_tracking_number,omitempty"`
	PickupAddress          *string `json:"pickup_address,omitempty"`
	
	// --- Fees & Pricing ---
	TaxAmount              float64 `json:"tax_amount"`
	ServiceFee             float64 `json:"service_fee"`
	AdminFee               float64 `json:"admin_fee"`
	PhysicalCheckFee       float64 `json:"physical_check_fee"`
	DeliveryFee            float64 `json:"delivery_fee"`
	ExpressFee             float64 `json:"express_fee"`
	
	// --- Service Level ---
	ServiceLevel           string  `json:"service_level"`
	IsNameTransferRequired bool    `json:"is_name_transfer_required"`
	
	SamsatOrigin           *string `json:"samsat_origin,omitempty"`
	SamsatDestination      *string `json:"samsat_destination,omitempty"`

	Notes                  string  `json:"notes"`
}

type CreateOrderCustomerRequest struct {
	// Customer Details (Bisa kirim ID atau data langsung)
	CustomerID  string `json:"customer_id,omitempty"`
	Fullname    string `json:"fullname,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	NIK         string `json:"nik,omitempty"`
	Email       string `json:"email,omitempty"`
	Address     string `json:"address,omitempty"`

	// Vehicle Details (Bisa kirim ID atau data langsung)
	VehicleID   string `json:"vehicle_id,omitempty"`
	PlateNumber string `json:"plate_number,omitempty"`
	ModelID     string `json:"model_id,omitempty"`
	Year        int    `json:"year,omitempty"`

	// Service & Logistics
	ServiceID   string `json:"service_id" validate:"required,uuid"`
	
	PickupMethod           string  `json:"pickup_method"`
	ReturnMethod           string  `json:"return_method"`
	CustomerTrackingNumber *string `json:"customer_tracking_number,omitempty"`
	PickupAddress          *string `json:"pickup_address,omitempty"`
	
	DeliveryFee            float64 `json:"delivery_fee"`
	CourierName            string  `json:"courier_name,omitempty"`
	CourierService         string  `json:"courier_service,omitempty"`
	PostalCode             string  `json:"postal_code,omitempty"`
	DestinationPostalCode  string  `json:"destination_postal_code,omitempty"`

	ServiceLevel           string  `json:"service_level"`
	IsNameTransferRequired bool    `json:"is_name_transfer_required"`
	
	SamsatOrigin           *string `json:"samsat_origin,omitempty"`
	SamsatDestination      *string `json:"samsat_destination,omitempty"`

	Notes                  string  `json:"notes"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=Pending Processing Ready Completed Cancelled Rejected"`
	Notes  string `json:"notes"`
}

type OrderStatusHistoryResponse struct {
	Status    string    `json:"status"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderResponse struct {
	ID            string    `json:"id"`
	CustomerID    string    `json:"customer_id"`
	VehicleID     string    `json:"vehicle_id"`
	ServiceID     string    `json:"service_id"`
	ProcessedByID          *string   `json:"processed_by_id"`
	Status                 string    `json:"status"`
	
	PickupMethod           string    `json:"pickup_method"`
	ReturnMethod           string    `json:"return_method"`
	CustomerTrackingNumber *string   `json:"customer_tracking_number,omitempty"`
	BiroJasaTrackingNumber *string   `json:"biro_jasa_tracking_number,omitempty"`
	PickupAddress          *string   `json:"pickup_address,omitempty"`
	
	TaxAmount              float64   `json:"tax_amount"`
	ServiceFee             float64   `json:"service_fee"`
	AdminFee               float64   `json:"admin_fee"`
	PhysicalCheckFee       float64   `json:"physical_check_fee"`
	DeliveryFee            float64   `json:"delivery_fee"`
	ExpressFee             float64   `json:"express_fee"`
	TotalCost              float64   `json:"total_cost"`
	
	ServiceLevel           string    `json:"service_level"`
	IsNameTransferRequired bool      `json:"is_name_transfer_required"`
	
	SamsatOrigin           *string   `json:"samsat_origin,omitempty"`
	SamsatDestination      *string   `json:"samsat_destination,omitempty"`

	Notes                  string    `json:"notes"`
	OrderDate              time.Time `json:"order_date"`
	
	TrackingNumber         *string   `json:"tracking_number,omitempty"`
	CourierName            *string   `json:"courier_name,omitempty"`
	CourierService         *string   `json:"courier_service,omitempty"`
	ShippingCost           float64   `json:"shipping_cost"`

	History []OrderStatusHistoryResponse `json:"history,omitempty"`
}
