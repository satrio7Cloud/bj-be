package usecase

import (
	"errors"

	"bj/internal/domain/model/dto"
	"bj/internal/domain/model/entity"
	"bj/internal/domain/model/repository"

	"github.com/google/uuid"
)

type ModelUsecase interface {
	CreateModel(req dto.CreateModelRequest) error
	GetAllModels() ([]dto.ModelResponse, error)
	GetModelsByBrand(brandID string) ([]dto.ModelResponse, error)
	GetModelByID(id string) (*dto.ModelResponse, error)
	UpdateModel(id string, req dto.UpdateModelRequest) error
	DeleteModel(id string) error
}

type modelUsecase struct {
	modelRepo repository.ModelRepository
}

func NewModelUsecase(repo repository.ModelRepository) ModelUsecase {
	return &modelUsecase{repo}
}

// Mapper functions
func mapModelToResponse(m entity.Model) dto.ModelResponse {
	return dto.ModelResponse{
		ID:          m.ID,
		BrandID:     m.BrandID,
		BrandName:   m.BrandName,
		ModelName:   m.ModelName,
		VehicleType: m.VehicleType,
	}
}

func (u *modelUsecase) CreateModel(req dto.CreateModelRequest) error {
	model := &entity.Model{
		ID:          uuid.New().String(),
		BrandID:     req.BrandID,
		ModelName:   req.ModelName,
		VehicleType: req.VehicleType,
	}
	return u.modelRepo.Create(model)
}

func (u *modelUsecase) GetAllModels() ([]dto.ModelResponse, error) {
	models, err := u.modelRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var responses []dto.ModelResponse
	for _, m := range models {
		responses = append(responses, mapModelToResponse(m))
	}
	return responses, nil
}

func (u *modelUsecase) GetModelsByBrand(brandID string) ([]dto.ModelResponse, error) {
	models, err := u.modelRepo.FindByBrandID(brandID)
	if err != nil {
		return nil, err
	}
	var responses []dto.ModelResponse
	for _, m := range models {
		responses = append(responses, mapModelToResponse(m))
	}
	return responses, nil
}

func (u *modelUsecase) GetModelByID(id string) (*dto.ModelResponse, error) {
	model, err := u.modelRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, errors.New("model tidak ditemukan")
	}

	resp := mapModelToResponse(*model)
	return &resp, nil
}

func (u *modelUsecase) UpdateModel(id string, req dto.UpdateModelRequest) error {
	model, err := u.modelRepo.FindByID(id)
	if err != nil {
		return err
	}
	if model == nil {
		return errors.New("model tidak ditemukan")
	}

	if req.BrandID != "" {
		model.BrandID = req.BrandID
	}
	if req.ModelName != "" {
		model.ModelName = req.ModelName
	}
	if req.VehicleType != "" {
		model.VehicleType = req.VehicleType
	}

	return u.modelRepo.Update(model)
}

func (u *modelUsecase) DeleteModel(id string) error {
	return u.modelRepo.Delete(id)
}
