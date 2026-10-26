package usecase

import (
	"bj/internal/domain/company/dto"
	"bj/internal/domain/company/entity"
	"bj/internal/domain/company/repository"
	"time"

	"github.com/google/uuid"
)

type CompanyProfileUsecase interface {
	GetProfile() (*dto.CompanyProfileResponse, error)
	CreateProfile(req dto.CreateCompanyProfileRequest) (*dto.CompanyProfileResponse, error)
	UpdateProfile(req dto.UpdateCompanyProfileRequest) (*dto.CompanyProfileResponse, error)
}

type companyProfileUsecase struct {
	repo repository.CompanyProfileRepository
}

func NewCompanyProfileUsecase(repo repository.CompanyProfileRepository) CompanyProfileUsecase {
	return &companyProfileUsecase{repo: repo}
}

func (u *companyProfileUsecase) GetProfile() (*dto.CompanyProfileResponse, error) {
	profile, err := u.repo.Get()
	if err != nil {
		return nil, err
	}

	return &dto.CompanyProfileResponse{
		ID:               profile.ID,
		OfficeName:       profile.OfficeName,
		FullAddress:      profile.FullAddress,
		PostCode:         profile.PostCode,
		PhoneNumber:      profile.PhoneNumber,
		OperationalHours: profile.OperationalHours,
		CreatedAt:        profile.CreatedAt,
		UpdatedAt:        profile.UpdatedAt,
	}, nil
}

func (u *companyProfileUsecase) CreateProfile(req dto.CreateCompanyProfileRequest) (*dto.CompanyProfileResponse, error) {
	now := time.Now()
	profile := &entity.CompanyProfile{
		ID:               uuid.New().String(),
		OfficeName:       req.OfficeName,
		FullAddress:      req.FullAddress,
		PostCode:         req.PostCode,
		PhoneNumber:      req.PhoneNumber,
		OperationalHours: req.OperationalHours,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := u.repo.Create(profile); err != nil {
		return nil, err
	}

	return &dto.CompanyProfileResponse{
		ID:               profile.ID,
		OfficeName:       profile.OfficeName,
		FullAddress:      profile.FullAddress,
		PostCode:         profile.PostCode,
		PhoneNumber:      profile.PhoneNumber,
		OperationalHours: profile.OperationalHours,
		CreatedAt:        profile.CreatedAt,
		UpdatedAt:        profile.UpdatedAt,
	}, nil
}

func (u *companyProfileUsecase) UpdateProfile(req dto.UpdateCompanyProfileRequest) (*dto.CompanyProfileResponse, error) {
	profile, err := u.repo.Get()
	if err != nil {
		return nil, err
	}

	profile.OfficeName = req.OfficeName
	profile.FullAddress = req.FullAddress
	profile.PostCode = req.PostCode
	profile.PhoneNumber = req.PhoneNumber
	profile.OperationalHours = req.OperationalHours
	profile.UpdatedAt = time.Now()

	if err := u.repo.Update(profile); err != nil {
		return nil, err
	}

	return &dto.CompanyProfileResponse{
		ID:               profile.ID,
		OfficeName:       profile.OfficeName,
		FullAddress:      profile.FullAddress,
		PostCode:         profile.PostCode,
		PhoneNumber:      profile.PhoneNumber,
		OperationalHours: profile.OperationalHours,
		CreatedAt:        profile.CreatedAt,
		UpdatedAt:        profile.UpdatedAt,
	}, nil
}
