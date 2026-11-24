package dto

import "bj/internal/domain/vehicle/entity"

type CreateVehicleRequest struct {
	CustomerID  string `json:"customer_id" validate:"required"`
	ModelID     string `json:"model_id" validate:"required"`
	PlateNumber string `json:"plate_number" validate:"required"`
	Year        int    `json:"year" validate:"required,gt=1900"`
}

func (req *CreateVehicleRequest) ToEntity(id string) *entity.Vehicle {
	return &entity.Vehicle{
		ID:          id,
		CustomerID:  req.CustomerID,
		ModelID:     req.ModelID,
		PlateNumber: req.PlateNumber,
		Year:        req.Year,
	}
}

type UpdateVehicleRequest struct {
	ModelID     string `json:"model_id"`
	PlateNumber string `json:"plate_number"`
	Year        int    `json:"year" validate:"omitempty,gt=1900"`
}

func (req *UpdateVehicleRequest) ApplyToEntity(v *entity.Vehicle) {
	if req.ModelID != "" {
		v.ModelID = req.ModelID
	}
	if req.PlateNumber != "" {
		v.PlateNumber = req.PlateNumber
	}
	if req.Year > 0 {
		v.Year = req.Year
	}
}

type VehicleResponse struct {
	ID          string `json:"id"`
	CustomerID  string `json:"customer_id"`
	ModelID     string `json:"model_id"`
	PlateNumber string `json:"plate_number"`
	Year        int    `json:"year"`
}

func ToVehicleResponse(v entity.Vehicle) VehicleResponse {
	return VehicleResponse{
		ID:          v.ID,
		CustomerID:  v.CustomerID,
		ModelID:     v.ModelID,
		PlateNumber: v.PlateNumber,
		Year:        v.Year,
	}
}

func ToVehicleResponses(vehicles []entity.Vehicle) []VehicleResponse {
	responses := make([]VehicleResponse, 0, len(vehicles))
	for _, v := range vehicles {
		responses = append(responses, ToVehicleResponse(v))
	}
	return responses
}
