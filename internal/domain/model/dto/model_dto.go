package dto

import "bj/internal/domain/model/entity"

type CreateModelRequest struct {
	BrandID     string `json:"brand_id" validate:"required"`
	ModelName   string `json:"model_name" validate:"required"`
	VehicleType string `json:"vehicle_type" validate:"required,oneof=motor mobil"`
}

func (req *CreateModelRequest) ToEntity(id string) *entity.Model {
	return &entity.Model{
		ID:          id,
		BrandID:     req.BrandID,
		ModelName:   req.ModelName,
		VehicleType: req.VehicleType,
	}
}

type UpdateModelRequest struct {
	BrandID     string `json:"brand_id"`
	ModelName   string `json:"model_name"`
	VehicleType string `json:"vehicle_type" validate:"omitempty,oneof=motor mobil"`
}

func (req *UpdateModelRequest) ApplyToEntity(m *entity.Model) {
	if req.BrandID != "" {
		m.BrandID = req.BrandID
	}
	if req.ModelName != "" {
		m.ModelName = req.ModelName
	}
	if req.VehicleType != "" {
		m.VehicleType = req.VehicleType
	}
}

type ModelResponse struct {
	ID          string `json:"id"`
	BrandID     string `json:"brand_id"`
	BrandName   string `json:"brand_name"`
	ModelName   string `json:"model_name"`
	VehicleType string `json:"vehicle_type"`
}

func ToModelResponse(m entity.Model) ModelResponse {
	return ModelResponse{
		ID:          m.ID,
		BrandID:     m.BrandID,
		BrandName:   m.BrandName,
		ModelName:   m.ModelName,
		VehicleType: m.VehicleType,
	}
}

func ToModelResponses(models []entity.Model) []ModelResponse {
	responses := make([]ModelResponse, 0, len(models))
	for _, m := range models {
		responses = append(responses, ToModelResponse(m))
	}
	return responses
}
