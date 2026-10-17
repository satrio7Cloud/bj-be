package usecase

import (
	"errors"
	"strings"

	"bj/internal/domain/biteship/dto"
	"bj/internal/domain/biteship/service"
)

type BiteshipUsecase interface {
	CheckRates(req dto.RatesRequest) (*dto.RateResponse, error)
	GetAvailableCouriers() (*dto.CourierResponse, error)
	CreateOrder(req dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
	GetTracking(trackingID string) (*dto.TrackingResponse, error)
	CalculateShippingFee(courierName, courierService, destPostalCode, pickupMethod string) (float64, error)
}

type biteshipUsecase struct {
	biteshipService   service.BiteshipService
	defaultCouriers   string
	originPostalCode  string
	defaultItemName   string
	defaultItemWeight int
	defaultItemValue  int
}

func NewBiteshipUsecase(
	biteshipService service.BiteshipService,
	defaultCouriers, originPostalCode string,
	defaultItemName string,
	defaultItemWeight, defaultItemValue int,
) BiteshipUsecase {
	return &biteshipUsecase{
		biteshipService:   biteshipService,
		defaultCouriers:   defaultCouriers,
		originPostalCode:  originPostalCode,
		defaultItemName:   defaultItemName,
		defaultItemWeight: defaultItemWeight,
		defaultItemValue:  defaultItemValue,
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

	if len(req.Items) == 0 {
		itemName := u.defaultItemName
		if itemName == "" {
			itemName = "Dokumen Layanan Biro Jasa"
		}
		itemWeight := u.defaultItemWeight
		if itemWeight <= 0 {
			itemWeight = 200
		}
		itemValue := u.defaultItemValue
		if itemValue <= 0 {
			itemValue = 100000
		}

		req.Items = []dto.Item{
			{
				Name:     itemName,
				Value:    itemValue,
				Weight:   itemWeight,
				Quantity: 1,
			},
		}
	}

	return u.biteshipService.GetRates(req)
}

func (u *biteshipUsecase) GetAvailableCouriers() (*dto.CourierResponse, error) {
	resp, err := u.biteshipService.GetCouriers()
	if err != nil {
		return nil, err
	}

	cleanedCouriers := strings.Trim(u.defaultCouriers, "\" '")
	if cleanedCouriers == "" {
		return resp, nil
	}

	allowedList := strings.Split(cleanedCouriers, ",")
	allowedMap := make(map[string]bool)
	for _, c := range allowedList {
		trimmed := strings.ToLower(strings.Trim(c, "\" '"))
		if trimmed != "" {
			allowedMap[trimmed] = true
		}
	}

	filtered := make([]dto.CourierData, 0)
	for _, courier := range resp.Couriers {
		if allowedMap[strings.ToLower(courier.CourierCode)] {
			filtered = append(filtered, courier)
		}
	}
	resp.Couriers = filtered

	return resp, nil
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

func (u *biteshipUsecase) CalculateShippingFee(courierName, courierService, destPostalCode, pickupMethod string) (float64, error) {
	pickupUpper := strings.ToUpper(strings.TrimSpace(pickupMethod))
	if pickupUpper == "WALK_IN" || pickupUpper == "DATANG_SENDIRI" {
		return 0, nil
	}

	courierName = strings.TrimSpace(courierName)
	courierService = strings.TrimSpace(courierService)
	destPostalCode = strings.TrimSpace(destPostalCode)

	if courierName == "" || courierService == "" || destPostalCode == "" {
		return 0, errors.New("data kurir atau kodepos tujuan tidak lengkap")
	}

	ratesReq := dto.RatesRequest{
		OriginPostalCode:      u.originPostalCode,
		DestinationPostalCode: destPostalCode,
		Couriers:              strings.ToLower(courierName),
	}

	rateResp, err := u.CheckRates(ratesReq)
	if err != nil {
		return 0, err
	}
	if rateResp == nil || len(rateResp.Pricing) == 0 {
		return 0, errors.New("tarif pengiriman tidak ditemukan dari Biteship")
	}

	var foundPrice float64 = 0
	for _, p := range rateResp.Pricing {
		sameCourier := strings.EqualFold(p.CourierName, courierName)
		sameService := strings.EqualFold(p.CourierService, courierService) ||
			strings.Contains(strings.ToLower(p.CourierService), strings.ToLower(courierService)) ||
			strings.Contains(strings.ToLower(courierService), strings.ToLower(p.CourierService))

		if sameCourier && sameService {
			foundPrice = float64(p.Price)
			break
		}
	}

	if foundPrice == 0 {
		return 0, errors.New("layanan kurir yang dipilih tidak tersedia untuk rute ini")
	}

	// Skenario A (Full Kurir: Jemput + Antar) = 2 Arah (price * 2)
	// Skenario B/C (Kirim Sendiri/Mandiri ke Kantor, Antar Pulang via Kurir) = 1 Arah (price * 1)
	if pickupUpper == "MANDIRI" || pickupUpper == "SEMI_MANDIRI" || pickupUpper == "KIRIM_EKSTERNAL" || pickupUpper == "ANTAR_SENDIRI" {
		return foundPrice, nil
	}

	return foundPrice * 2, nil
}
