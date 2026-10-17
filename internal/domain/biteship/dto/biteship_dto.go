package dto

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RatesRequest struct {
	OriginPostalCode      string `json:"origin_postal_code"`
	DestinationPostalCode string `json:"destination_postal_code"`
	Couriers              string `json:"couriers"`
	Items                 []Item `json:"items"`
}

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int    `json:"value"`
	Length      int    `json:"length"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Weight      int    `json:"weight"`
	Quantity    int    `json:"quantity"`
}

type RateResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Object  string   `json:"object"`
	Pricing []Pricing `json:"pricing"`
}

type Pricing struct {
	CourierName    string `json:"courier_name"`
	CourierService string `json:"courier_service_name"`
	Price          int    `json:"price"`
}

type CreateOrderRequest struct {
	ShipperContactName      string `json:"shipper_contact_name"`
	ShipperContactPhone     string `json:"shipper_contact_phone"`
	ShipperContactEmail     string `json:"shipper_contact_email,omitempty"`
	OriginContactName       string `json:"origin_contact_name"`
	OriginContactPhone      string `json:"origin_contact_phone"`
	OriginAddress           string `json:"origin_address"`
	OriginPostalCode        string `json:"origin_postal_code"`
	DestinationContactName  string `json:"destination_contact_name"`
	DestinationContactPhone string `json:"destination_contact_phone"`
	DestinationAddress      string `json:"destination_address"`
	DestinationPostalCode   string `json:"destination_postal_code"`
	CourierCompany          string `json:"courier_company"`
	CourierType             string `json:"courier_type"`
	DeliveryType            string `json:"delivery_type"`
	DeliveryDate            string `json:"delivery_date,omitempty"`
	DeliveryTime            string `json:"delivery_time,omitempty"`
	Items                   []Item `json:"items"`
}

type CreateOrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Object  string `json:"object"`
	ID      string `json:"id"`
	Status  string `json:"status"`
	Courier struct {
		TrackingID string `json:"tracking_id"`
		WaybillID  string `json:"waybill_id"`
		Company    string `json:"company"`
		Name       string `json:"name"`
		Type       string `json:"type"`
	} `json:"courier"`
	Price float64 `json:"price"`
}

type TrackingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
	Courier struct {
		Company   string `json:"company"`
		Name      string `json:"name"`
		WaybillID string `json:"waybill_id"`
	} `json:"courier"`
	History []struct {
		Note      string `json:"note"`
		UpdatedAt string `json:"updated_at"`
		Status    string `json:"status"`
	} `json:"history"`
}

type CourierResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Object  string        `json:"object"`
	Couriers []CourierData `json:"couriers"`
}

type CourierData struct {
	AvailableForCashOnDelivery   bool   `json:"available_for_cash_on_delivery"`
	AvailableForProofOfDelivery  bool   `json:"available_for_proof_of_delivery"`
	AvailableForInstantWaybillID bool   `json:"available_for_instant_waybill_id"`
	AvailableForInsurance        bool   `json:"available_for_insurance"`
	CourierCode                  string `json:"courier_code"`
	CourierName                  string `json:"courier_name"`
	CourierServiceCode           string `json:"courier_service_code"`
	CourierServiceName           string `json:"courier_service_name"`
	Tier                         string `json:"tier"`
	Description                  string `json:"description"`
	ServiceType                  string `json:"service_type"`
	ShippingType                 string `json:"shipping_type"`
	ShipmentDurationRange        string `json:"shipment_duration_range"`
	ShipmentDurationUnit         string `json:"shipment_duration_unit"`
}
