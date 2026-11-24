package dto

import "bj/internal/domain/courier/entity"

// CourierResponse adalah data satu kurir dan layanannya dari tabel master.
type CourierResponse struct {
	ID                 string `json:"id"`
	CourierCode        string `json:"courier_code"`
	CourierName        string `json:"courier_name"`
	CourierServiceCode string `json:"courier_service_code"`
	CourierServiceName string `json:"courier_service_name"`
	Tier               string `json:"tier"`
	ServiceType        string `json:"service_type"`
	ShippingType       string `json:"shipping_type"`
	DurationRange      string `json:"duration_range"`
	DurationUnit       string `json:"duration_unit"`
	IsAvailable        bool   `json:"is_available"`
}

// SyncCouriersResponse adalah respons setelah sinkronisasi data kurir dari Biteship.
type SyncCouriersResponse struct {
	Synced  int    `json:"synced"`
	Message string `json:"message"`
}

func ToCourierResponse(c entity.Courier) CourierResponse {
	return CourierResponse{
		ID:                 c.ID,
		CourierCode:        c.CourierCode,
		CourierName:        c.CourierName,
		CourierServiceCode: c.CourierServiceCode,
		CourierServiceName: c.CourierServiceName,
		Tier:               c.Tier,
		ServiceType:        c.ServiceType,
		ShippingType:       c.ShippingType,
		DurationRange:      c.DurationRange,
		DurationUnit:       c.DurationUnit,
		IsAvailable:        c.IsAvailable,
	}
}

func ToCourierResponses(couriers []entity.Courier) []CourierResponse {
	responses := make([]CourierResponse, 0, len(couriers))
	for _, c := range couriers {
		responses = append(responses, ToCourierResponse(c))
	}
	return responses
}
