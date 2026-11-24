package usecase

import (
	"errors"

	"bj/internal/domain/vehicle/dto"
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

func (u *vehicleUsecase) CreateVehicle(req dto.CreateVehicleRequest) (string, error) {
	newID := uuid.New().String()
	vehicle := req.ToEntity(newID)
	err := u.vehicleRepo.Create(vehicle)
	return newID, err
}

func (u *vehicleUsecase) GetAllVehicles() ([]dto.VehicleResponse, error) {
	vehicles, err := u.vehicleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.ToVehicleResponses(vehicles), nil
}

func (u *vehicleUsecase) GetVehiclesByCustomer(customerID string) ([]dto.VehicleResponse, error) {
	vehicles, err := u.vehicleRepo.FindByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	return dto.ToVehicleResponses(vehicles), nil
}

func (u *vehicleUsecase) GetVehicleByID(id string) (*dto.VehicleResponse, error) {
	vehicle, err := u.vehicleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("kendaraan tidak ditemukan")
	}

	resp := dto.ToVehicleResponse(*vehicle)
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

	req.ApplyToEntity(vehicle)
	return u.vehicleRepo.Update(vehicle)
}

func (u *vehicleUsecase) DeleteVehicle(id string) error {
	return u.vehicleRepo.Delete(id)
}
