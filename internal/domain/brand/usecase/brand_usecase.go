package usecase

import (
	"errors"

	"bj/internal/domain/brand/dto"
	"bj/internal/domain/brand/entity"
	"bj/internal/domain/brand/repository"

	"github.com/google/uuid"
)

type BrandUsecase interface {
	CreateBrand(req dto.CreateBrandRequest) error
	GetAllBrands() ([]dto.BrandResponse, error)
	GetBrandByID(id string) (*dto.BrandResponse, error)
	UpdateBrand(id string, req dto.UpdateBrandRequest) error
	DeleteBrand(id string) error
}

type brandUsecase struct {
	brandRepo repository.BrandRepository
}

func NewBrandUsecase(repo repository.BrandRepository) BrandUsecase {
	return &brandUsecase{repo}
}

// Mapper functions
func mapBrandToResponse(b entity.Brand) dto.BrandResponse {
	return dto.BrandResponse{
		ID:        b.ID,
		BrandName: b.BrandName,
	}
}

func (u *brandUsecase) CreateBrand(req dto.CreateBrandRequest) error {
	brand := &entity.Brand{
		ID:        uuid.New().String(),
		BrandName: req.BrandName,
	}
	return u.brandRepo.Create(brand)
}

func (u *brandUsecase) GetAllBrands() ([]dto.BrandResponse, error) {
	brands, err := u.brandRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.BrandResponse
	for _, b := range brands {
		responses = append(responses, mapBrandToResponse(b))
	}
	return responses, nil
}

func (u *brandUsecase) GetBrandByID(id string) (*dto.BrandResponse, error) {
	brand, err := u.brandRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, errors.New("merek tidak ditemukan")
	}

	resp := mapBrandToResponse(*brand)
	return &resp, nil
}

func (u *brandUsecase) UpdateBrand(id string, req dto.UpdateBrandRequest) error {
	brand, err := u.brandRepo.FindByID(id)
	if err != nil {
		return err
	}
	if brand == nil {
		return errors.New("merek tidak ditemukan")
	}

	brand.BrandName = req.BrandName
	return u.brandRepo.Update(brand)
}

func (u *brandUsecase) DeleteBrand(id string) error {
	return u.brandRepo.Delete(id)
}
