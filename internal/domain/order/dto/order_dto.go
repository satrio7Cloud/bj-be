package dto

import (
	"time"

	"bj/internal/constants"
	custEntity "bj/internal/domain/customer/entity"
	"bj/internal/domain/order/entity"
	vehEntity "bj/internal/domain/vehicle/entity"
	"bj/internal/utils"
)

type CreateOrderAdminRequest struct {
	CustomerID string `json:"customer_id" validate:"required,uuid"`
	VehicleID  string `json:"vehicle_id" validate:"required,uuid"`
	ServiceID  string `json:"service_id" validate:"required,uuid"`

	// --- Delivery Options ---
	PickupMethod           string  `json:"pickup_method"`
	ReturnMethod           string  `json:"return_method"`
	CustomerTrackingNumber *string `json:"customer_tracking_number,omitempty"`
	PickupAddress          *string `json:"pickup_address,omitempty"`

	// --- Fees & Pricing ---
	TaxAmount        float64 `json:"tax_amount"`
	ServiceFee       float64 `json:"service_fee"`
	AdminFee         float64 `json:"admin_fee"`
	PhysicalCheckFee float64 `json:"physical_check_fee"`
	DeliveryFee      float64 `json:"delivery_fee"`
	ExpressFee       float64 `json:"express_fee"`

	// --- Service Level ---
	ServiceLevel           string `json:"service_level"`
	IsNameTransferRequired bool   `json:"is_name_transfer_required"`

	SamsatOrigin      *string `json:"samsat_origin,omitempty"`
	SamsatDestination *string `json:"samsat_destination,omitempty"`

	Notes string `json:"notes"`
}

func (req *CreateOrderAdminRequest) CalculateTotalCost() float64 {
	return req.TaxAmount + req.ServiceFee + req.AdminFee + req.PhysicalCheckFee + req.DeliveryFee + req.ExpressFee
}

func (req *CreateOrderAdminRequest) ToEntity(id string, adminID string) *entity.Order {
	return &entity.Order{
		ID:         id,
		CustomerID: req.CustomerID,
		VehicleID:  req.VehicleID,
		ServiceID:  req.ServiceID,

		PickupMethod:           req.PickupMethod,
		ReturnMethod:           req.ReturnMethod,
		CustomerTrackingNumber: req.CustomerTrackingNumber,
		PickupAddress:          req.PickupAddress,

		TaxAmount:        req.TaxAmount,
		ServiceFee:       req.ServiceFee,
		AdminFee:         req.AdminFee,
		PhysicalCheckFee: req.PhysicalCheckFee,
		DeliveryFee:      req.DeliveryFee,
		ExpressFee:       req.ExpressFee,
		TotalCost:        req.CalculateTotalCost(),

		ServiceLevel:           req.ServiceLevel,
		IsNameTransferRequired: req.IsNameTransferRequired,

		SamsatOrigin:      req.SamsatOrigin,
		SamsatDestination: req.SamsatDestination,

		Notes:         req.Notes,
		Status:        constants.OrderStatusPending,
		ProcessedByID: &adminID,
		OrderDate:     time.Now(),
	}
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
	ServiceID string `json:"service_id" validate:"required,uuid"`

	PickupMethod           string  `json:"pickup_method"`
	ReturnMethod           string  `json:"return_method"`
	CustomerTrackingNumber *string `json:"customer_tracking_number,omitempty"`
	PickupAddress          *string `json:"pickup_address,omitempty"`

	DeliveryFee           float64 `json:"delivery_fee"`
	CourierName           string  `json:"courier_name,omitempty"`
	CourierService        string  `json:"courier_service,omitempty"`
	PostalCode            string  `json:"postal_code,omitempty"`
	DestinationPostalCode string  `json:"destination_postal_code,omitempty"`

	ServiceLevel           string `json:"service_level"`
	IsNameTransferRequired bool   `json:"is_name_transfer_required"`

	SamsatOrigin      *string `json:"samsat_origin,omitempty"`
	SamsatDestination *string `json:"samsat_destination,omitempty"`

	Notes string `json:"notes"`
}

func (req *CreateOrderCustomerRequest) ToCustomerEntity(id string) *custEntity.Customer {
	return &custEntity.Customer{
		ID:          id,
		Fullname:    req.Fullname,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		NIK:         req.NIK,
		Email:       req.Email,
		CreatedAt:   time.Now(),
	}
}

func (req *CreateOrderCustomerRequest) ToVehicleEntity(id string, customerID string) *vehEntity.Vehicle {
	year := req.Year
	if year == 0 {
		year = time.Now().Year()
	}
	return &vehEntity.Vehicle{
		ID:          id,
		CustomerID:  customerID,
		ModelID:     req.ModelID,
		PlateNumber: req.PlateNumber,
		Year:        year,
	}
}

// OrderFees menyimpan rincian biaya yang dikalkulasikan untuk pesanan customer.
type OrderFees struct {
	ServiceFee       float64
	AdminFee         float64
	PhysicalCheckFee float64
	DeliveryFee      float64
	ExpressFee       float64
	TotalCost        float64
}

func (req *CreateOrderCustomerRequest) ToOrderEntity(id string, customerID, vehicleID string, fees OrderFees) *entity.Order {
	var courierName *string
	if req.CourierName != "" {
		courierName = &req.CourierName
	}
	var courierService *string
	if req.CourierService != "" {
		courierService = &req.CourierService
	}

	pickupAddress := req.PickupAddress
	if (pickupAddress == nil || *pickupAddress == "") && req.Address != "" {
		pickupAddress = &req.Address
	}

	return &entity.Order{
		ID:         id,
		CustomerID: customerID,
		VehicleID:  vehicleID,
		ServiceID:  req.ServiceID,

		PickupMethod:           req.PickupMethod,
		ReturnMethod:           req.ReturnMethod,
		CustomerTrackingNumber: req.CustomerTrackingNumber,
		PickupAddress:          pickupAddress,

		TaxAmount:        0,
		ServiceFee:       fees.ServiceFee,
		AdminFee:         fees.AdminFee,
		PhysicalCheckFee: fees.PhysicalCheckFee,
		DeliveryFee:      fees.DeliveryFee,
		ExpressFee:       fees.ExpressFee,
		TotalCost:        fees.TotalCost,

		CourierName:    courierName,
		CourierService: courierService,
		ShippingCost:   fees.DeliveryFee,

		ServiceLevel:           req.ServiceLevel,
		IsNameTransferRequired: req.IsNameTransferRequired,

		SamsatOrigin:      req.SamsatOrigin,
		SamsatDestination: req.SamsatDestination,

		Notes:         req.Notes,
		Status:        constants.OrderStatusPending,
		ProcessedByID: nil,
		OrderDate:     time.Now(),
	}
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
	ID            string  `json:"id"`
	CustomerID    string  `json:"customer_id"`
	VehicleID     string  `json:"vehicle_id"`
	ServiceID     string  `json:"service_id"`
	ProcessedByID *string `json:"processed_by_id"`
	Status        string  `json:"status"`

	PickupMethod           string  `json:"pickup_method"`
	ReturnMethod           string  `json:"return_method"`
	CustomerTrackingNumber *string `json:"customer_tracking_number,omitempty"`
	BiroJasaTrackingNumber *string `json:"biro_jasa_tracking_number,omitempty"`
	PickupAddress          *string `json:"pickup_address,omitempty"`

	TaxAmount        float64 `json:"tax_amount"`
	ServiceFee       float64 `json:"service_fee"`
	AdminFee         float64 `json:"admin_fee"`
	PhysicalCheckFee float64 `json:"physical_check_fee"`
	DeliveryFee      float64 `json:"delivery_fee"`
	ExpressFee       float64 `json:"express_fee"`
	TotalCost        float64 `json:"total_cost"`

	ServiceLevel           string `json:"service_level"`
	IsNameTransferRequired bool   `json:"is_name_transfer_required"`

	SamsatOrigin      *string `json:"samsat_origin,omitempty"`
	SamsatDestination *string `json:"samsat_destination,omitempty"`

	Notes     string    `json:"notes"`
	OrderDate time.Time `json:"order_date"`

	TrackingNumber *string `json:"tracking_number,omitempty"`
	CourierName    *string `json:"courier_name,omitempty"`
	CourierService *string `json:"courier_service,omitempty"`
	ShippingCost   float64 `json:"shipping_cost"`

	History []OrderStatusHistoryResponse `json:"history,omitempty"`
}

func ToOrderResponse(o entity.Order) OrderResponse {
	return OrderResponse{
		ID:            o.ID,
		CustomerID:    o.CustomerID,
		VehicleID:     o.VehicleID,
		ServiceID:     o.ServiceID,
		ProcessedByID: o.ProcessedByID,
		Status:        o.Status,

		PickupMethod:           o.PickupMethod,
		ReturnMethod:           o.ReturnMethod,
		CustomerTrackingNumber: o.CustomerTrackingNumber,
		BiroJasaTrackingNumber: o.BiroJasaTrackingNumber,
		PickupAddress:          o.PickupAddress,

		TaxAmount:        o.TaxAmount,
		ServiceFee:       o.ServiceFee,
		AdminFee:         o.AdminFee,
		PhysicalCheckFee: o.PhysicalCheckFee,
		DeliveryFee:      o.DeliveryFee,
		ExpressFee:       o.ExpressFee,
		TotalCost:        o.TotalCost,

		ServiceLevel:           o.ServiceLevel,
		IsNameTransferRequired: o.IsNameTransferRequired,

		SamsatOrigin:      o.SamsatOrigin,
		SamsatDestination: o.SamsatDestination,

		Notes:     o.Notes,
		OrderDate: utils.ToWIB(o.OrderDate),

		TrackingNumber: o.TrackingNumber,
		CourierName:    o.CourierName,
		CourierService: o.CourierService,
		ShippingCost:   o.ShippingCost,
	}
}

func ToOrderResponses(orders []entity.Order) []OrderResponse {
	responses := make([]OrderResponse, 0, len(orders))
	for _, o := range orders {
		responses = append(responses, ToOrderResponse(o))
	}
	return responses
}
