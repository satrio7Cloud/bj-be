package usecase

import (
	"encoding/json"
	"errors"
	"strings"

	biteshipDTO "bj/internal/domain/biteship/dto"
	biteshipSvc "bj/internal/domain/biteship/service"
	courierEntity "bj/internal/domain/courier/entity"
	courierRepo "bj/internal/domain/courier/repository"
	"bj/internal/domain/shipping/dto"
	"bj/internal/domain/shipping/entity"
	"bj/internal/domain/shipping/repository"
)

type ShippingUsecase interface {
	GetSettings() (*dto.ShippingSettingsResponse, error)
	SaveSettings(adminID string, req dto.UpdateShippingSettingsRequest) error
	CalculateRates(req dto.CalculateRatesRequest) (*dto.CalculateRatesResponse, error)
}

type shippingUsecase struct {
	settingsRepo repository.ShippingSettingsRepository
	courierRepo  courierRepo.CourierRepository
	biteshipSvc  biteshipSvc.BiteshipService
	originPostal string
}

func NewShippingUsecase(
	settingsRepo repository.ShippingSettingsRepository,
	courierRepo courierRepo.CourierRepository,
	biteshipSvc biteshipSvc.BiteshipService,
	originPostalCode string,
) ShippingUsecase {
	return &shippingUsecase{
		settingsRepo: settingsRepo,
		courierRepo:  courierRepo,
		biteshipSvc:  biteshipSvc,
		originPostal: originPostalCode,
	}
}

func (u *shippingUsecase) GetSettings() (*dto.ShippingSettingsResponse, error) {
	settings, err := u.settingsRepo.Get()
	if err != nil {
		return nil, err
	}
	if settings == nil {
		return &dto.ShippingSettingsResponse{
			ID:                "default",
			EnabledCouriers:   []dto.EnabledCourierItem{},
			DefaultItemWeight: 200,
			DefaultItemValue:  100000,
		}, nil
	}
	return dto.ToShippingSettingsResponse(*settings)
}

func (u *shippingUsecase) SaveSettings(adminID string, req dto.UpdateShippingSettingsRequest) error {
	enabledJSON, err := json.Marshal(req.EnabledCouriers)
	if err != nil {
		return errors.New("gagal memproses daftar kurir aktif")
	}

	s := &entity.ShippingSettings{
		ID:                "default",
		EnabledCouriers:   string(enabledJSON),
		DefaultItemWeight: req.DefaultItemWeight,
		DefaultItemValue:  req.DefaultItemValue,
		UpdatedBy:         adminID,
	}
	return u.settingsRepo.Upsert(s)
}

func (u *shippingUsecase) CalculateRates(req dto.CalculateRatesRequest) (*dto.CalculateRatesResponse, error) {
	settings, err := u.GetSettings()
	if err != nil {
		return nil, err
	}

	pickupUpper := strings.ToUpper(strings.TrimSpace(req.PickupMethod))
	multiplier := 1.0
	if pickupUpper == "DOOR_TO_DOOR" {
		multiplier = 2.0
	}
	if pickupUpper == "WALK_IN" || pickupUpper == "PICKUP_INTERNAL" {
		return &dto.CalculateRatesResponse{
			PickupMethod: req.PickupMethod,
			Multiplier:   0,
			Rates:        []dto.RateOption{},
		}, nil
	}

	if len(settings.EnabledCouriers) == 0 {
		return nil, errors.New("belum ada kurir yang diaktifkan admin")
	}

	// Build set enabled "courier_code|courier_service_code"
	enabledSet := make(map[string]bool)
	uniqueCodesMap := make(map[string]bool)
	for _, item := range settings.EnabledCouriers {
		code := strings.ToLower(strings.TrimSpace(item.CourierCode))
		svcCode := strings.ToLower(strings.TrimSpace(item.CourierServiceCode))
		if code != "" && svcCode != "" {
			enabledSet[code+"|"+svcCode] = true
			uniqueCodesMap[code] = true
		}
	}

	uniqueCodes := make([]string, 0, len(uniqueCodesMap))
	for code := range uniqueCodesMap {
		uniqueCodes = append(uniqueCodes, code)
	}
	courierStr := strings.Join(uniqueCodes, ",")

	// Hit Biteship API
	ratesReq := biteshipDTO.RatesRequest{
		OriginPostalCode:      u.originPostal,
		DestinationPostalCode: req.DestinationPostalCode,
		Couriers:              courierStr,
		Items: []biteshipDTO.Item{
			{
				Name:     "Dokumen Layanan Biro Jasa",
				Value:    int(settings.DefaultItemValue),
				Weight:   settings.DefaultItemWeight,
				Quantity: 1,
			},
		},
	}

	rateResp, err := u.biteshipSvc.GetRates(ratesReq)
	if err != nil {
		return nil, errors.New("gagal mengambil tarif dari Biteship: " + err.Error())
	}

	// Lookup master kurir dari domain courier (tabel dbo.ShippingCouriers)
	dbCouriers, _ := u.courierRepo.FindByCodes(uniqueCodes)

	// Filter & bangun response rates
	rates := make([]dto.RateOption, 0, len(rateResp.Pricing))
	for _, p := range rateResp.Pricing {
		var matchedCourier *courierEntity.Courier
		for i, c := range dbCouriers {
			courierMatch := strings.EqualFold(c.CourierCode, p.CourierCode) ||
				strings.EqualFold(c.CourierCode, p.Company) ||
				strings.EqualFold(c.CourierCode, p.CourierName) ||
				strings.EqualFold(c.CourierName, p.CourierName) ||
				strings.Contains(strings.ToLower(p.CourierName), strings.ToLower(c.CourierCode)) ||
				strings.Contains(strings.ToLower(c.CourierName), strings.ToLower(p.CourierName)) ||
				strings.Contains(strings.ToLower(p.CourierName), strings.ToLower(c.CourierName))

			if !courierMatch {
				continue
			}

			pSvcName := strings.ToLower(strings.TrimSpace(p.CourierServiceName))
			pSvcCode := strings.ToLower(strings.TrimSpace(p.CourierServiceCode))
			pType := strings.ToLower(strings.TrimSpace(p.ServiceType))
			cSvcName := strings.ToLower(strings.TrimSpace(c.CourierServiceName))
			cSvcCode := strings.ToLower(strings.TrimSpace(c.CourierServiceCode))

			serviceMatch := strings.EqualFold(cSvcCode, pSvcCode) ||
				strings.EqualFold(cSvcName, pSvcName) ||
				strings.EqualFold(cSvcCode, pSvcName) ||
				strings.EqualFold(cSvcCode, pType) ||
				strings.Contains(pSvcName, cSvcName) ||
				strings.Contains(cSvcName, pSvcName) ||
				strings.Contains(pSvcName, "("+cSvcCode+")") ||
				strings.Contains(pSvcName, " "+cSvcCode)

			if serviceMatch {
				matchedCourier = &dbCouriers[i]
				break
			}
		}

		if matchedCourier == nil {
			continue
		}

		filterKey := strings.ToLower(matchedCourier.CourierCode) + "|" + strings.ToLower(matchedCourier.CourierServiceCode)
		if !enabledSet[filterKey] {
			continue
		}

		basePrice := float64(p.Price)
		finalPrice := basePrice * multiplier

		durRange := p.DurationRange
		if durRange == "" {
			durRange = matchedCourier.DurationRange
		}
		durUnit := p.DurationUnit
		if durUnit == "" {
			durUnit = matchedCourier.DurationUnit
		}

		svcName := p.CourierServiceName
		if svcName == "" {
			svcName = matchedCourier.CourierServiceName
		}

		rates = append(rates, dto.RateOption{
			CourierCode:        matchedCourier.CourierCode,
			CourierName:        matchedCourier.CourierName,
			CourierServiceCode: matchedCourier.CourierServiceCode,
			CourierServiceName: svcName,
			DurationRange:      durRange,
			DurationUnit:       durUnit,
			BasePrice:          basePrice,
			FinalPrice:         finalPrice,
		})
	}

	return &dto.CalculateRatesResponse{
		PickupMethod: req.PickupMethod,
		Multiplier:   multiplier,
		Rates:        rates,
	}, nil
}
