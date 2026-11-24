package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"bj/internal/domain/deliverymethod/dto"
	"bj/internal/domain/deliverymethod/repository"
	"bj/internal/shared/cache"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	deliveryMethodsAllCacheKey    = "master:delivery_methods:all"
	deliveryMethodsActiveCacheKey = "master:delivery_methods:active:"
	deliveryMethodsCacheTTL       = 30 * time.Minute
)

type DeliveryMethodUsecase interface {
	Create(req dto.CreateDeliveryMethodRequest) (*dto.DeliveryMethodResponse, error)
	GetAll() ([]dto.DeliveryMethodResponse, error)
	GetActive(direction string) ([]dto.DeliveryMethodResponse, error)
	GetByID(id string) (*dto.DeliveryMethodResponse, error)
	Update(id string, req dto.UpdateDeliveryMethodRequest) error
	ToggleActive(id string) error
	Delete(id string) error
}

type deliveryMethodUsecase struct {
	repo        repository.DeliveryMethodRepository
	redisClient *redis.Client
}

func NewDeliveryMethodUsecase(repo repository.DeliveryMethodRepository, redisClient *redis.Client) DeliveryMethodUsecase {
	return &deliveryMethodUsecase{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (u *deliveryMethodUsecase) invalidateCache() {
	ctx := context.Background()
	cache.Invalidate(ctx, u.redisClient, deliveryMethodsAllCacheKey)
	cache.InvalidateByPattern(ctx, u.redisClient, deliveryMethodsActiveCacheKey+"*")
}

func (u *deliveryMethodUsecase) Create(req dto.CreateDeliveryMethodRequest) (*dto.DeliveryMethodResponse, error) {
	code := strings.ToUpper(strings.TrimSpace(req.MethodCode))

	existing, err := u.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("method_code sudah digunakan")
	}

	m := req.ToEntity(uuid.New().String())
	if err := u.repo.Create(m); err != nil {
		return nil, err
	}

	u.invalidateCache()

	resp := dto.ToDeliveryMethodResponse(*m)
	return &resp, nil
}

func (u *deliveryMethodUsecase) GetAll() ([]dto.DeliveryMethodResponse, error) {
	return cache.GetOrSet(context.Background(), u.redisClient, deliveryMethodsAllCacheKey, deliveryMethodsCacheTTL, func() ([]dto.DeliveryMethodResponse, error) {
		list, err := u.repo.FindAll()
		if err != nil {
			return nil, err
		}
		return dto.ToDeliveryMethodResponses(list), nil
	})
}

func (u *deliveryMethodUsecase) GetActive(direction string) ([]dto.DeliveryMethodResponse, error) {
	cleanDir := strings.ToUpper(strings.TrimSpace(direction))
	cacheKey := deliveryMethodsActiveCacheKey + cleanDir

	return cache.GetOrSet(context.Background(), u.redisClient, cacheKey, deliveryMethodsCacheTTL, func() ([]dto.DeliveryMethodResponse, error) {
		list, err := u.repo.FindActive(cleanDir)
		if err != nil {
			return nil, err
		}
		return dto.ToDeliveryMethodResponses(list), nil
	})
}

func (u *deliveryMethodUsecase) GetByID(id string) (*dto.DeliveryMethodResponse, error) {
	m, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("metode pengantaran tidak ditemukan")
	}

	resp := dto.ToDeliveryMethodResponse(*m)
	return &resp, nil
}

func (u *deliveryMethodUsecase) Update(id string, req dto.UpdateDeliveryMethodRequest) error {
	existing, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("metode pengantaran tidak ditemukan")
	}

	req.ApplyToEntity(existing)

	err = u.repo.Update(existing)
	if err == nil {
		u.invalidateCache()
	}
	return err
}

func (u *deliveryMethodUsecase) ToggleActive(id string) error {
	existing, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("metode pengantaran tidak ditemukan")
	}

	err = u.repo.ToggleActive(id)
	if err == nil {
		u.invalidateCache()
	}
	return err
}

func (u *deliveryMethodUsecase) Delete(id string) error {
	existing, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("metode pengantaran tidak ditemukan")
	}

	err = u.repo.Delete(id)
	if err == nil {
		u.invalidateCache()
	}
	return err
}
