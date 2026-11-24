package dto

import (
	"strings"
	"time"

	"bj/internal/domain/deliverymethod/entity"
)

// CreateDeliveryMethodRequest adalah payload untuk membuat master metode pengantaran baru.
type CreateDeliveryMethodRequest struct {
	MethodCode      string `json:"method_code" validate:"required,max=50"`
	MethodName      string `json:"method_name" validate:"required,max=100"`
	Direction       string `json:"direction" validate:"required,oneof=PICKUP RETURN BOTH"`
	ProviderType    string `json:"provider_type" validate:"required,oneof=SELF COURIER INTERNAL"`
	Description     string `json:"description" validate:"max=255"`
	FlatFee         int64  `json:"flat_fee" validate:"gte=0"`
	RequiresAddress bool   `json:"requires_address"`
	RequiresCourier bool   `json:"requires_courier"`
	IsActive        *bool  `json:"is_active"`
	DisplayOrder    int    `json:"display_order" validate:"gte=1"`
}

func (req *CreateDeliveryMethodRequest) ToEntity(id string) *entity.DeliveryMethod {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	displayOrder := req.DisplayOrder
	if displayOrder <= 0 {
		displayOrder = 1
	}

	return &entity.DeliveryMethod{
		ID:              id,
		MethodCode:      strings.ToUpper(strings.TrimSpace(req.MethodCode)),
		MethodName:      strings.TrimSpace(req.MethodName),
		Direction:       strings.ToUpper(strings.TrimSpace(req.Direction)),
		ProviderType:    strings.ToUpper(strings.TrimSpace(req.ProviderType)),
		Description:     strings.TrimSpace(req.Description),
		FlatFee:         req.FlatFee,
		RequiresAddress: req.RequiresAddress,
		RequiresCourier: req.RequiresCourier,
		IsActive:        isActive,
		DisplayOrder:    displayOrder,
	}
}

// UpdateDeliveryMethodRequest adalah payload untuk mengedit master metode pengantaran.
type UpdateDeliveryMethodRequest struct {
	MethodName      string `json:"method_name" validate:"required,max=100"`
	Direction       string `json:"direction" validate:"required,oneof=PICKUP RETURN BOTH"`
	ProviderType    string `json:"provider_type" validate:"required,oneof=SELF COURIER INTERNAL"`
	Description     string `json:"description" validate:"max=255"`
	FlatFee         int64  `json:"flat_fee" validate:"gte=0"`
	RequiresAddress bool   `json:"requires_address"`
	RequiresCourier bool   `json:"requires_courier"`
	IsActive        bool   `json:"is_active"`
	DisplayOrder    int    `json:"display_order" validate:"gte=1"`
}

func (req *UpdateDeliveryMethodRequest) ApplyToEntity(m *entity.DeliveryMethod) {
	m.MethodName = strings.TrimSpace(req.MethodName)
	m.Direction = strings.ToUpper(strings.TrimSpace(req.Direction))
	m.ProviderType = strings.ToUpper(strings.TrimSpace(req.ProviderType))
	m.Description = strings.TrimSpace(req.Description)
	m.FlatFee = req.FlatFee
	m.RequiresAddress = req.RequiresAddress
	m.RequiresCourier = req.RequiresCourier
	m.IsActive = req.IsActive
	m.DisplayOrder = req.DisplayOrder
}

// DeliveryMethodResponse merepresentasikan data respon satu metode pengantaran.
type DeliveryMethodResponse struct {
	ID              string    `json:"id"`
	MethodCode      string    `json:"method_code"`
	MethodName      string    `json:"method_name"`
	Direction       string    `json:"direction"`
	ProviderType    string    `json:"provider_type"`
	Description     string    `json:"description"`
	FlatFee         int64     `json:"flat_fee"`
	RequiresAddress bool      `json:"requires_address"`
	RequiresCourier bool      `json:"requires_courier"`
	IsActive        bool      `json:"is_active"`
	DisplayOrder    int       `json:"display_order"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func ToDeliveryMethodResponse(m entity.DeliveryMethod) DeliveryMethodResponse {
	return DeliveryMethodResponse{
		ID:              m.ID,
		MethodCode:      m.MethodCode,
		MethodName:      m.MethodName,
		Direction:       m.Direction,
		ProviderType:    m.ProviderType,
		Description:     m.Description,
		FlatFee:         m.FlatFee,
		RequiresAddress: m.RequiresAddress,
		RequiresCourier: m.RequiresCourier,
		IsActive:        m.IsActive,
		DisplayOrder:    m.DisplayOrder,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func ToDeliveryMethodResponses(methods []entity.DeliveryMethod) []DeliveryMethodResponse {
	responses := make([]DeliveryMethodResponse, 0, len(methods))
	for _, m := range methods {
		responses = append(responses, ToDeliveryMethodResponse(m))
	}
	return responses
}
