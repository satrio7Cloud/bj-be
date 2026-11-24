package usecase

import (
	"context"
	"time"

	"bj/internal/domain/company/dto"
	"bj/internal/domain/company/repository"
	"bj/internal/shared/cache"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	companyProfileCacheKey = "master:company_profile"
	companyProfileCacheTTL = 1 * time.Hour
)

type CompanyProfileUsecase interface {
	GetProfile() (*dto.CompanyProfileResponse, error)
	CreateProfile(req dto.CreateCompanyProfileRequest) (*dto.CompanyProfileResponse, error)
	UpdateProfile(req dto.UpdateCompanyProfileRequest) (*dto.CompanyProfileResponse, error)
}

type companyProfileUsecase struct {
	repo        repository.CompanyProfileRepository
	redisClient *redis.Client
}

func NewCompanyProfileUsecase(repo repository.CompanyProfileRepository, redisClient *redis.Client) CompanyProfileUsecase {
	return &companyProfileUsecase{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (u *companyProfileUsecase) invalidateCache() {
	cache.Invalidate(context.Background(), u.redisClient, companyProfileCacheKey)
}

func (u *companyProfileUsecase) GetProfile() (*dto.CompanyProfileResponse, error) {
	resp, err := cache.GetOrSet(context.Background(), u.redisClient, companyProfileCacheKey, companyProfileCacheTTL, func() (dto.CompanyProfileResponse, error) {
		profile, err := u.repo.Get()
		if err != nil {
			return dto.CompanyProfileResponse{}, err
		}
		return dto.ToCompanyProfileResponse(*profile), nil
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (u *companyProfileUsecase) CreateProfile(req dto.CreateCompanyProfileRequest) (*dto.CompanyProfileResponse, error) {
	profile := req.ToEntity(uuid.New().String())
	if err := u.repo.Create(profile); err != nil {
		return nil, err
	}

	u.invalidateCache()

	resp := dto.ToCompanyProfileResponse(*profile)
	return &resp, nil
}

func (u *companyProfileUsecase) UpdateProfile(req dto.UpdateCompanyProfileRequest) (*dto.CompanyProfileResponse, error) {
	profile, err := u.repo.Get()
	if err != nil {
		return nil, err
	}

	req.ApplyToEntity(profile)

	if err := u.repo.Update(profile); err != nil {
		return nil, err
	}

	u.invalidateCache()

	resp := dto.ToCompanyProfileResponse(*profile)
	return &resp, nil
}
