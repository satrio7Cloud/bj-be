package usecase

import (
	"context"
	"errors"
	"time"

	"bj/internal/domain/service/dto"
	"bj/internal/domain/service/repository"
	"bj/internal/shared/cache"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	servicesAllCacheKey = "services:all"
	servicesCacheTTL    = 5 * time.Minute
)

type ServiceUsecase interface {
	CreateService(req dto.CreateServiceRequest) error
	GetAllServices() ([]dto.ServiceResponse, error)
	GetServiceByID(id string) (*dto.ServiceResponse, error)
	UpdateService(id string, req dto.UpdateServiceRequest) error
	DeleteService(id string) error
}

type serviceUsecase struct {
	serviceRepo repository.ServiceRepository
	redisClient *redis.Client
}

func NewServiceUsecase(repo repository.ServiceRepository, redisClient *redis.Client) ServiceUsecase {
	return &serviceUsecase{
		serviceRepo: repo,
		redisClient: redisClient,
	}
}

func (u *serviceUsecase) invalidateCache() {
	cache.Invalidate(context.Background(), u.redisClient, servicesAllCacheKey)
}

func (u *serviceUsecase) CreateService(req dto.CreateServiceRequest) error {
	service := req.ToEntity(uuid.New().String())
	err := u.serviceRepo.Create(service)
	if err == nil {
		u.invalidateCache()
	}
	return err
}

func (u *serviceUsecase) GetAllServices() ([]dto.ServiceResponse, error) {
	return cache.GetOrSet(context.Background(), u.redisClient, servicesAllCacheKey, servicesCacheTTL, func() ([]dto.ServiceResponse, error) {
		services, err := u.serviceRepo.FindAll()
		if err != nil {
			return nil, err
		}
		return dto.ToServiceResponses(services), nil
	})
}

func (u *serviceUsecase) GetServiceByID(id string) (*dto.ServiceResponse, error) {
	service, err := u.serviceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, errors.New("layanan tidak ditemukan")
	}

	resp := dto.ToServiceResponse(*service)
	return &resp, nil
}

func (u *serviceUsecase) UpdateService(id string, req dto.UpdateServiceRequest) error {
	service, err := u.serviceRepo.FindByID(id)
	if err != nil {
		return err
	}
	if service == nil {
		return errors.New("layanan tidak ditemukan")
	}

	req.ApplyToEntity(service)

	err = u.serviceRepo.Update(service)
	if err == nil {
		u.invalidateCache()
	}
	return err
}

func (u *serviceUsecase) DeleteService(id string) error {
	err := u.serviceRepo.Delete(id)
	if err == nil {
		u.invalidateCache()
	}
	return err
}
