package dto

type CreateModelRequest struct {
	BrandID     string `json:"brand_id" validate:"required"`
	ModelName   string `json:"model_name" validate:"required"`
	VehicleType string `json:"vehicle_type" validate:"required,oneof=motor mobil"`
}

type UpdateModelRequest struct {
	BrandID     string `json:"brand_id"`
	ModelName   string `json:"model_name"`
	VehicleType string `json:"vehicle_type" validate:"omitempty,oneof=motor mobil"`
}

type ModelResponse struct {
	ID          string `json:"id"`
	BrandID     string `json:"brand_id"`
	BrandName   string `json:"brand_name"`
	ModelName   string `json:"model_name"`
	VehicleType string `json:"vehicle_type"`
}
