package usecase

import (
	"errors"
	"time"

	biteshipSvc "bj/internal/domain/biteship/service"
	"bj/internal/domain/courier/dto"
	"bj/internal/domain/courier/entity"
	"bj/internal/domain/courier/repository"

	"github.com/google/uuid"
)

type CourierUsecase interface {
	SyncCouriers() (*dto.SyncCouriersResponse, error)
	GetAllCouriers() ([]dto.CourierResponse, error)
	FindByCodes(codes []string) ([]entity.Courier, error)
}

type courierUsecase struct {
	repo        repository.CourierRepository
	biteshipSvc biteshipSvc.BiteshipService
}

func NewCourierUsecase(repo repository.CourierRepository, biteshipSvc biteshipSvc.BiteshipService) CourierUsecase {
	return &courierUsecase{
		repo:        repo,
		biteshipSvc: biteshipSvc,
	}
}

func (u *courierUsecase) SyncCouriers() (*dto.SyncCouriersResponse, error) {
	biteshipResp, err := u.biteshipSvc.GetCouriers()
	if err != nil {
		return nil, errors.New("gagal mengambil data kurir dari Biteship: " + err.Error())
	}

	couriers := make([]entity.Courier, 0, len(biteshipResp.Couriers))
	now := time.Now()
	for _, c := range biteshipResp.Couriers {
		couriers = append(couriers, entity.Courier{
			ID:                 uuid.New().String(),
			CourierCode:        c.CourierCode,
			CourierName:        c.CourierName,
			CourierServiceCode: c.CourierServiceCode,
			CourierServiceName: c.CourierServiceName,
			Tier:               c.Tier,
			ServiceType:        c.ServiceType,
			ShippingType:       c.ShippingType,
			DurationRange:      c.ShipmentDurationRange,
			DurationUnit:       c.ShipmentDurationUnit,
			IsAvailable:        true,
			CreatedAt:          now,
			UpdatedAt:          now,
		})
	}

	if err := u.repo.UpsertAll(couriers); err != nil {
		return nil, errors.New("gagal menyimpan data kurir ke database: " + err.Error())
	}

	return &dto.SyncCouriersResponse{
		Synced:  len(couriers),
		Message: "Sinkronisasi berhasil",
	}, nil
}

func (u *courierUsecase) GetAllCouriers() ([]dto.CourierResponse, error) {
	couriers, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.ToCourierResponses(couriers), nil
}

func (u *courierUsecase) FindByCodes(codes []string) ([]entity.Courier, error) {
	return u.repo.FindByCodes(codes)
}
