package dto

import (
	"time"

	"bj/internal/domain/service/entity"
)

type CreateServiceRequest struct {
	ServiceName      string  `json:"service_name" validate:"required"`
	Category         string  `json:"category" validate:"required"`
	ServiceFee       float64 `json:"service_fee" validate:"gte=0"`
	ExpressFee       float64 `json:"express_fee" validate:"gte=0"`
	AdminFee         float64 `json:"admin_fee" validate:"gte=0"`
	PhysicalCheckFee float64 `json:"physical_check_fee" validate:"gte=0"`
	SWDKLLJFee       float64 `json:"swdkllj_fee" validate:"gte=0"`
	Description      string  `json:"description"`
}

func (req *CreateServiceRequest) ToEntity(id string) *entity.Service {
	return &entity.Service{
		ID:               id,
		ServiceName:      req.ServiceName,
		Category:         req.Category,
		ServiceFee:       req.ServiceFee,
		ExpressFee:       req.ExpressFee,
		AdminFee:         req.AdminFee,
		PhysicalCheckFee: req.PhysicalCheckFee,
		SWDKLLJFee:       req.SWDKLLJFee,
		Description:      req.Description,
		CreatedAt:        time.Now(),
	}
}

type UpdateServiceRequest struct {
	ServiceName      string  `json:"service_name"`
	Category         string  `json:"category"`
	ServiceFee       float64 `json:"service_fee" validate:"omitempty,gte=0"`
	ExpressFee       float64 `json:"express_fee" validate:"omitempty,gte=0"`
	AdminFee         float64 `json:"admin_fee" validate:"omitempty,gte=0"`
	PhysicalCheckFee float64 `json:"physical_check_fee" validate:"omitempty,gte=0"`
	SWDKLLJFee       float64 `json:"swdkllj_fee" validate:"omitempty,gte=0"`
	Description      string  `json:"description"`
}

func (req *UpdateServiceRequest) ApplyToEntity(s *entity.Service) {
	if req.ServiceName != "" {
		s.ServiceName = req.ServiceName
	}
	if req.Category != "" {
		s.Category = req.Category
	}
	s.ServiceFee = req.ServiceFee
	s.ExpressFee = req.ExpressFee
	s.AdminFee = req.AdminFee
	s.PhysicalCheckFee = req.PhysicalCheckFee
	s.SWDKLLJFee = req.SWDKLLJFee
	if req.Description != "" {
		s.Description = req.Description
	}
}

type ServiceResponse struct {
	ID               string    `json:"id"`
	ServiceName      string    `json:"service_name"`
	Category         string    `json:"category"`
	ServiceFee       float64   `json:"service_fee"`
	ExpressFee       float64   `json:"express_fee"`
	AdminFee         float64   `json:"admin_fee"`
	PhysicalCheckFee float64   `json:"physical_check_fee"`
	SWDKLLJFee       float64   `json:"swdkllj_fee"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
}

func ToServiceResponse(s entity.Service) ServiceResponse {
	return ServiceResponse{
		ID:               s.ID,
		ServiceName:      s.ServiceName,
		Category:         s.Category,
		ServiceFee:       s.ServiceFee,
		ExpressFee:       s.ExpressFee,
		AdminFee:         s.AdminFee,
		PhysicalCheckFee: s.PhysicalCheckFee,
		SWDKLLJFee:       s.SWDKLLJFee,
		Description:      s.Description,
		CreatedAt:        s.CreatedAt,
	}
}

func ToServiceResponses(services []entity.Service) []ServiceResponse {
	responses := make([]ServiceResponse, 0, len(services))
	for _, s := range services {
		responses = append(responses, ToServiceResponse(s))
	}
	return responses
}
