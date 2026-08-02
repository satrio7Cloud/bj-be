package dto

type CreateVehicleRequest struct {
	CustomerID  string `json:"customer_id" validate:"required"`
	ModelID     string `json:"model_id" validate:"required"`
	PlateNumber string `json:"plate_number" validate:"required"`
	Year        int    `json:"year" validate:"required,gt=1900"`
}

type UpdateVehicleRequest struct {
	ModelID     string `json:"model_id"`
	PlateNumber string `json:"plate_number"`
	Year        int    `json:"year" validate:"omitempty,gt=1900"`
}

type VehicleResponse struct {
	ID          string `json:"id"`
	CustomerID  string `json:"customer_id"`
	ModelID     string `json:"model_id"`
	PlateNumber string `json:"plate_number"`
	Year        int    `json:"year"`
}
