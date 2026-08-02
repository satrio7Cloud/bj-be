package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"bj/internal/domain/service/dto"
	"bj/internal/domain/service/entity"
	"bj/internal/domain/service/repository"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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

// Mapper functions
func mapServiceToResponse(s entity.Service) dto.ServiceResponse {
	return dto.ServiceResponse{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Category:    s.Category,
		ServiceFee:  s.ServiceFee,
		ExpressFee:  s.ExpressFee,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
	}
}

func (u *serviceUsecase) CreateService(req dto.CreateServiceRequest) error {
	service := &entity.Service{
		ID:          uuid.New().String(),
		ServiceName: req.ServiceName,
		Category:    req.Category,
		ServiceFee:  req.ServiceFee,
		ExpressFee:  req.ExpressFee,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}
	
	
	// Invalidate Redis cache on create
	ctx := context.Background()
	u.redisClient.Del(ctx, "services:all")
	
	return u.serviceRepo.Create(service)
}

func (u *serviceUsecase) GetAllServices() ([]dto.ServiceResponse, error) {
	ctx := context.Background()
	cacheKey := "services:all"

	// 1. Check Redis cache first
	val, err := u.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit
		var cachedResponses []dto.ServiceResponse
		if err := json.Unmarshal([]byte(val), &cachedResponses); err == nil {
			return cachedResponses, nil
		}
	}

	// 2. Cache miss, fetch from DB
	services, err := u.serviceRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.ServiceResponse
	for _, s := range services {
		responses = append(responses, mapServiceToResponse(s))
	}

	// 3. Update Redis cache
	if cacheBytes, err := json.Marshal(responses); err == nil {
		u.redisClient.Set(ctx, cacheKey, cacheBytes, 5*time.Minute)
	}

	return responses, nil
}

func (u *serviceUsecase) GetServiceByID(id string) (*dto.ServiceResponse, error) {
	service, err := u.serviceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, errors.New("layanan tidak ditemukan")
	}

	resp := mapServiceToResponse(*service)
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

	if req.ServiceName != "" {
		service.ServiceName = req.ServiceName
	}
	if req.Category != "" {
		service.Category = req.Category
	}
	// Note: 0 is valid for ServiceFee if free, but normally we check if it's passed or use pointers
	service.ServiceFee = req.ServiceFee
	service.ExpressFee = req.ExpressFee
	if req.Description != "" {
		service.Description = req.Description
	}

	err = u.serviceRepo.Update(service)
	
	// Invalidate Redis cache on update
	if err == nil {
		ctx := context.Background()
		u.redisClient.Del(ctx, "services:all")
	}
	
	return err
}

func (u *serviceUsecase) DeleteService(id string) error {
	err := u.serviceRepo.Delete(id)
	
	// Invalidate Redis cache on delete
	if err == nil {
		ctx := context.Background()
		u.redisClient.Del(ctx, "services:all")
	}
	
	return err
}
