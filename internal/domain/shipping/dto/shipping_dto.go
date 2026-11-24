package dto

import (
	"encoding/json"
	"errors"

	"bj/internal/domain/shipping/entity"
)

// EnabledCourierItem merepresentasikan satu pasangan kurir + layanan yang diaktifkan admin.
type EnabledCourierItem struct {
	CourierCode        string `json:"courier_code"`
	CourierName        string `json:"courier_name,omitempty"`
	CourierServiceCode string `json:"courier_service_code"`
	CourierServiceName string `json:"courier_service_name,omitempty"`
}

// ShippingSettingsResponse adalah respons konfigurasi pengaturan pengiriman.
type ShippingSettingsResponse struct {
	ID                string               `json:"id"`
	EnabledCouriers   []EnabledCourierItem `json:"enabled_couriers"`
	DefaultItemWeight int                  `json:"default_item_weight"`
	DefaultItemValue  int64                `json:"default_item_value"`
}

// UpdateShippingSettingsRequest adalah body request untuk menyimpan konfigurasi admin.
type UpdateShippingSettingsRequest struct {
	EnabledCouriers   []EnabledCourierItem `json:"enabled_couriers" validate:"required,min=1,dive"`
	DefaultItemWeight int                  `json:"default_item_weight" validate:"required,gt=0"`
	DefaultItemValue  int64                `json:"default_item_value" validate:"required,gt=0"`
}

// ─── Calculate Rates ─────────────────────────────────────────────────────────

// CalculateRatesRequest adalah body request untuk menghitung ongkir.
type CalculateRatesRequest struct {
	DestinationPostalCode string `json:"destination_postal_code" validate:"required"`
	PickupMethod          string `json:"pickup_method" validate:"required"` // DOOR_TO_DOOR, ANTAR_SENDIRI, KIRIM_EKSTERNAL, WALK_IN
}

// RateOption adalah satu pilihan kurir + layanan beserta harganya.
type RateOption struct {
	CourierCode        string  `json:"courier_code"`
	CourierName        string  `json:"courier_name"`
	CourierServiceCode string  `json:"courier_service_code"`
	CourierServiceName string  `json:"courier_service_name"`
	DurationRange      string  `json:"duration_range"`
	DurationUnit       string  `json:"duration_unit"`
	BasePrice          float64 `json:"base_price"`
	FinalPrice         float64 `json:"final_price"`
}

// CalculateRatesResponse adalah respons daftar tarif ongkir yang sudah difilter.
type CalculateRatesResponse struct {
	PickupMethod string       `json:"pickup_method"`
	Multiplier   float64      `json:"multiplier"`
	Rates        []RateOption `json:"rates"`
}

func ToShippingSettingsResponse(s entity.ShippingSettings) (*ShippingSettingsResponse, error) {
	var enabledCouriers []EnabledCourierItem
	if s.EnabledCouriers != "" {
		if err := json.Unmarshal([]byte(s.EnabledCouriers), &enabledCouriers); err != nil {
			return nil, errors.New("gagal membaca data kurir aktif")
		}
	}
	if enabledCouriers == nil {
		enabledCouriers = []EnabledCourierItem{}
	}
	return &ShippingSettingsResponse{
		ID:                s.ID,
		EnabledCouriers:   enabledCouriers,
		DefaultItemWeight: s.DefaultItemWeight,
		DefaultItemValue:  s.DefaultItemValue,
	}, nil
}
