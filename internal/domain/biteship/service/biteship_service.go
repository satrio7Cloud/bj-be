package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"bj/internal/config"
	"bj/internal/domain/biteship/dto"
)

type BiteshipService interface {
	GetRates(req dto.RatesRequest) (*dto.RateResponse, error)
	GetCouriers() (*dto.CourierResponse, error)
	CreateOrder(req dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
	GetTracking(trackingID string) (*dto.TrackingResponse, error)
}

type biteshipService struct {
	apiKey  string
	baseURL string
}

func NewBiteshipService(cfg *config.Config) BiteshipService {
	return &biteshipService{
		apiKey:  cfg.Biteship.APIKey,
		baseURL: cfg.Biteship.BaseURL,
	}
}

func (s *biteshipService) GetRates(req dto.RatesRequest) (*dto.RateResponse, error) {
	url := fmt.Sprintf("%s/rates/couriers", s.baseURL)
	
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("biteship API error: %s", string(bodyBytes))
	}

	var rateResp dto.RateResponse
	if err := json.Unmarshal(bodyBytes, &rateResp); err != nil {
		return nil, err
	}

	return &rateResp, nil
}

func (s *biteshipService) GetCouriers() (*dto.CourierResponse, error) {
	url := fmt.Sprintf("%s/couriers", s.baseURL)
	
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("biteship API error: %s", string(bodyBytes))
	}

	var courierResp dto.CourierResponse
	if err := json.Unmarshal(bodyBytes, &courierResp); err != nil {
		return nil, err
	}

	return &courierResp, nil
}

func (s *biteshipService) CreateOrder(req dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	url := fmt.Sprintf("%s/orders", s.baseURL)
	
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("biteship API error: %s", string(bodyBytes))
	}

	var orderResp dto.CreateOrderResponse
	if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
		return nil, err
	}

	return &orderResp, nil
}

func (s *biteshipService) GetTracking(trackingID string) (*dto.TrackingResponse, error) {
	url := fmt.Sprintf("%s/trackings/%s", s.baseURL, trackingID)
	
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("biteship tracking API error: %s", string(bodyBytes))
	}

	var trackResp dto.TrackingResponse
	if err := json.Unmarshal(bodyBytes, &trackResp); err != nil {
		return nil, err
	}

	return &trackResp, nil
}
