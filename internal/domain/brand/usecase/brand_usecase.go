package usecase

import (
	"errors"

	"bj/internal/domain/brand/dto"
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

func (u *brandUsecase) CreateBrand(req dto.CreateBrandRequest) error {
	brand := req.ToEntity(uuid.New().String())
	return u.brandRepo.Create(brand)
}

func (u *brandUsecase) GetAllBrands() ([]dto.BrandResponse, error) {
	brands, err := u.brandRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.ToBrandResponses(brands), nil
}

func (u *brandUsecase) GetBrandByID(id string) (*dto.BrandResponse, error) {
	brand, err := u.brandRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, errors.New("merek tidak ditemukan")
	}

	resp := dto.ToBrandResponse(*brand)
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

	req.ApplyToEntity(brand)
	return u.brandRepo.Update(brand)
}

func (u *brandUsecase) DeleteBrand(id string) error {
	return u.brandRepo.Delete(id)
}
