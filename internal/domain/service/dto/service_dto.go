package dto

import "time"

type CreateServiceRequest struct {
	ServiceName      string  `json:"service_name" validate:"required"`
	Category         string  `json:"category" validate:"required"`
	ServiceFee       float64 `json:"service_fee" validate:"gte=0"`
	ExpressFee       float64 `json:"express_fee" validate:"gte=0"`
	AdminFee         float64 `json:"admin_fee" validate:"gte=0"`
	PhysicalCheckFee float64 `json:"physical_check_fee" validate:"gte=0"`
	Description      string  `json:"description"`
}

type UpdateServiceRequest struct {
	ServiceName      string  `json:"service_name"`
	Category         string  `json:"category"`
	ServiceFee       float64 `json:"service_fee" validate:"omitempty,gte=0"`
	ExpressFee       float64 `json:"express_fee" validate:"omitempty,gte=0"`
	AdminFee         float64 `json:"admin_fee" validate:"omitempty,gte=0"`
	PhysicalCheckFee float64 `json:"physical_check_fee" validate:"omitempty,gte=0"`
	Description      string  `json:"description"`
}

type ServiceResponse struct {
	ID               string    `json:"id"`
	ServiceName      string    `json:"service_name"`
	Category         string    `json:"category"`
	ServiceFee       float64   `json:"service_fee"`
	ExpressFee       float64   `json:"express_fee"`
	AdminFee         float64   `json:"admin_fee"`
	PhysicalCheckFee float64   `json:"physical_check_fee"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
}
