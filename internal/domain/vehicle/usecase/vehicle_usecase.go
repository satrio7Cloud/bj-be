package usecase

import (
	"errors"

	"bj/internal/domain/vehicle/dto"
	"bj/internal/domain/vehicle/entity"
	"bj/internal/domain/vehicle/repository"

	"github.com/google/uuid"
)

type VehicleUsecase interface {
	CreateVehicle(req dto.CreateVehicleRequest) (string, error)
	GetAllVehicles() ([]dto.VehicleResponse, error)
	GetVehiclesByCustomer(customerID string) ([]dto.VehicleResponse, error)
	GetVehicleByID(id string) (*dto.VehicleResponse, error)
	UpdateVehicle(id string, req dto.UpdateVehicleRequest) error
	DeleteVehicle(id string) error
}

type vehicleUsecase struct {
	vehicleRepo repository.VehicleRepository
}

func NewVehicleUsecase(repo repository.VehicleRepository) VehicleUsecase {
	return &vehicleUsecase{repo}
}

// Mapper functions
func mapVehicleToResponse(v entity.Vehicle) dto.VehicleResponse {
	return dto.VehicleResponse{
		ID:          v.ID,
		CustomerID:  v.CustomerID,
		ModelID:     v.ModelID,
		PlateNumber: v.PlateNumber,
		Year:        v.Year,
	}
}

func (u *vehicleUsecase) CreateVehicle(req dto.CreateVehicleRequest) (string, error) {
	newID := uuid.New().String()
	vehicle := &entity.Vehicle{
		ID:          newID,
		CustomerID:  req.CustomerID,
		ModelID:     req.ModelID,
		PlateNumber: req.PlateNumber,
		Year:        req.Year,
	}
	err := u.vehicleRepo.Create(vehicle)
	return newID, err
}

func (u *vehicleUsecase) GetAllVehicles() ([]dto.VehicleResponse, error) {
	vehicles, err := u.vehicleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var responses []dto.VehicleResponse
	for _, v := range vehicles {
		responses = append(responses, mapVehicleToResponse(v))
	}
	return responses, nil
}

func (u *vehicleUsecase) GetVehiclesByCustomer(customerID string) ([]dto.VehicleResponse, error) {
	vehicles, err := u.vehicleRepo.FindByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	var responses []dto.VehicleResponse
	for _, v := range vehicles {
		responses = append(responses, mapVehicleToResponse(v))
	}
	return responses, nil
}

func (u *vehicleUsecase) GetVehicleByID(id string) (*dto.VehicleResponse, error) {
	vehicle, err := u.vehicleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("kendaraan tidak ditemukan")
	}

	resp := mapVehicleToResponse(*vehicle)
	return &resp, nil
}

func (u *vehicleUsecase) UpdateVehicle(id string, req dto.UpdateVehicleRequest) error {
	vehicle, err := u.vehicleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if vehicle == nil {
		return errors.New("kendaraan tidak ditemukan")
	}

	if req.ModelID != "" {
		vehicle.ModelID = req.ModelID
	}
	if req.PlateNumber != "" {
		vehicle.PlateNumber = req.PlateNumber
	}
	if req.Year > 0 {
		vehicle.Year = req.Year
	}

	return u.vehicleRepo.Update(vehicle)
}

func (u *vehicleUsecase) DeleteVehicle(id string) error {
	return u.vehicleRepo.Delete(id)
}
