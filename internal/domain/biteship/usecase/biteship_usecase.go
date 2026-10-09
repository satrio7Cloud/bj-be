package usecase

import (
	"bj/internal/domain/biteship/dto"
	"bj/internal/domain/biteship/service"
)

type BiteshipUsecase interface {
	CheckRates(req dto.RatesRequest) (*dto.RateResponse, error)
	GetAvailableCouriers() (*dto.CourierResponse, error)
	CreateOrder(req dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
	GetTracking(trackingID string) (*dto.TrackingResponse, error)
}

type biteshipUsecase struct {
	biteshipService  service.BiteshipService
	defaultCouriers  string
	originPostalCode string
}

func NewBiteshipUsecase(biteshipService service.BiteshipService, defaultCouriers, originPostalCode string) BiteshipUsecase {
	return &biteshipUsecase{
		biteshipService:  biteshipService,
		defaultCouriers:  defaultCouriers,
		originPostalCode: originPostalCode,
	}
}

func (u *biteshipUsecase) CheckRates(req dto.RatesRequest) (*dto.RateResponse, error) {
	// Di sini bisa ditambahkan validasi atau modifikasi data sebelum dipanggil ke API eksternal
	// Menggunakan nilai default dari config jika tidak dikirim dari request
	if req.OriginPostalCode == "" {
		req.OriginPostalCode = u.originPostalCode
	}

	if req.Couriers == "" {
		req.Couriers = u.defaultCouriers
	}

	return u.biteshipService.GetRates(req)
}

func (u *biteshipUsecase) GetAvailableCouriers() (*dto.CourierResponse, error) {
	return u.biteshipService.GetCouriers()
}

func (u *biteshipUsecase) CreateOrder(req dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	if req.OriginPostalCode == "" {
		req.OriginPostalCode = u.originPostalCode
	}
	// For production, maybe auto-populate office address details here if origin is not provided
	
	return u.biteshipService.CreateOrder(req)
}

func (u *biteshipUsecase) GetTracking(trackingID string) (*dto.TrackingResponse, error) {
	return u.biteshipService.GetTracking(trackingID)
}
