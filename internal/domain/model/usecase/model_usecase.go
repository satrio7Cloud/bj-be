package usecase

import (
	"context"
	"errors"
	"time"

	"bj/internal/domain/model/dto"
	"bj/internal/domain/model/repository"
	"bj/internal/shared/cache"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	modelsAllCacheKey   = "master:models:all"
	modelsBrandCacheKey = "master:models:brand:"
	modelsCacheTTL      = 1 * time.Hour
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
	modelRepo   repository.ModelRepository
	redisClient *redis.Client
}

func NewModelUsecase(repo repository.ModelRepository, redisClient *redis.Client) ModelUsecase {
	return &modelUsecase{
		modelRepo:   repo,
		redisClient: redisClient,
	}
}

func (u *modelUsecase) invalidateCache() {
	ctx := context.Background()
	cache.Invalidate(ctx, u.redisClient, modelsAllCacheKey)
	cache.InvalidateByPattern(ctx, u.redisClient, modelsBrandCacheKey+"*")
}

func (u *modelUsecase) CreateModel(req dto.CreateModelRequest) error {
	model := req.ToEntity(uuid.New().String())
	err := u.modelRepo.Create(model)
	if err == nil {
		u.invalidateCache()
	}
	return err
}

func (u *modelUsecase) GetAllModels() ([]dto.ModelResponse, error) {
	return cache.GetOrSet(context.Background(), u.redisClient, modelsAllCacheKey, modelsCacheTTL, func() ([]dto.ModelResponse, error) {
		models, err := u.modelRepo.FindAll()
		if err != nil {
			return nil, err
		}
		return dto.ToModelResponses(models), nil
	})
}

func (u *modelUsecase) GetModelsByBrand(brandID string) ([]dto.ModelResponse, error) {
	cacheKey := modelsBrandCacheKey + brandID

	return cache.GetOrSet(context.Background(), u.redisClient, cacheKey, modelsCacheTTL, func() ([]dto.ModelResponse, error) {
		models, err := u.modelRepo.FindByBrandID(brandID)
		if err != nil {
			return nil, err
		}
		return dto.ToModelResponses(models), nil
	})
}

func (u *modelUsecase) GetModelByID(id string) (*dto.ModelResponse, error) {
	model, err := u.modelRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, errors.New("model tidak ditemukan")
	}

	resp := dto.ToModelResponse(*model)
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

	req.ApplyToEntity(model)

	err = u.modelRepo.Update(model)
	if err == nil {
		u.invalidateCache()
	}
	return err
}

func (u *modelUsecase) DeleteModel(id string) error {
	err := u.modelRepo.Delete(id)
	if err == nil {
		u.invalidateCache()
	}
	return err
}
