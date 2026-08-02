package usecase

import (
	"errors"
	"strings"
	"time"

	"bj/internal/constants"
	"bj/internal/domain/order/dto"
	"bj/internal/domain/order/entity"
	"bj/internal/domain/order/repository"
	logEntity "bj/internal/domain/orderlog/entity"
	logRepo "bj/internal/domain/orderlog/repository"
	srvRepo "bj/internal/domain/service/repository"

	"github.com/google/uuid"
)

type OrderUsecase interface {
	CreateOrderAdmin(req dto.CreateOrderAdminRequest, adminID string) (string, error)
	CreateOrderCustomer(req dto.CreateOrderCustomerRequest) (string, error)
	GetAllOrders(page int, limit int) ([]dto.OrderResponse, error)
	GetOrderByID(id string) (*dto.OrderResponse, error)
	UpdateOrderStatus(id string, req dto.UpdateOrderStatusRequest, adminID string) error
}

type orderUsecase struct {
	orderRepo    repository.OrderRepository
	orderLogRepo logRepo.OrderLogRepository
	serviceRepo  srvRepo.ServiceRepository
}

func NewOrderUsecase(repo repository.OrderRepository, logRepo logRepo.OrderLogRepository, serviceRepo srvRepo.ServiceRepository) OrderUsecase {
	return &orderUsecase{
		orderRepo:    repo,
		orderLogRepo: logRepo,
		serviceRepo:  serviceRepo,
	}
}

// Mapper functions
func mapOrderToResponse(o entity.Order) dto.OrderResponse {
	return dto.OrderResponse{
		ID:            o.ID,
		CustomerID:    o.CustomerID,
		VehicleID:     o.VehicleID,
		ServiceID:     o.ServiceID,
		ProcessedByID: o.ProcessedByID,
		Status:        o.Status,

		PickupMethod:           o.PickupMethod,
		ReturnMethod:           o.ReturnMethod,
		CustomerTrackingNumber: o.CustomerTrackingNumber,
		BiroJasaTrackingNumber: o.BiroJasaTrackingNumber,
		PickupAddress:          o.PickupAddress,

		TaxAmount:        o.TaxAmount,
		ServiceFee:       o.ServiceFee,
		AdminFee:         o.AdminFee,
		PhysicalCheckFee: o.PhysicalCheckFee,
		DeliveryFee:      o.DeliveryFee,
		ExpressFee:       o.ExpressFee,
		TotalCost:        o.TotalCost,

		ServiceLevel:           o.ServiceLevel,
		IsNameTransferRequired: o.IsNameTransferRequired,

		SamsatOrigin:           o.SamsatOrigin,
		SamsatDestination:      o.SamsatDestination,

		Notes:     o.Notes,
		OrderDate: o.OrderDate,
	}
}

func (u *orderUsecase) CreateOrderAdmin(req dto.CreateOrderAdminRequest, adminID string) (string, error) {
	newID := uuid.New().String()
	// Kalkulasi Total Cost di Backend
	totalCost := req.TaxAmount + req.ServiceFee + req.AdminFee + req.PhysicalCheckFee + req.DeliveryFee + req.ExpressFee

	order := &entity.Order{
		ID:         newID,
		CustomerID: req.CustomerID,
		VehicleID:  req.VehicleID,
		ServiceID:  req.ServiceID,

		PickupMethod:           req.PickupMethod,
		ReturnMethod:           req.ReturnMethod,
		CustomerTrackingNumber: req.CustomerTrackingNumber,
		PickupAddress:          req.PickupAddress,

		TaxAmount:        req.TaxAmount,
		ServiceFee:       req.ServiceFee,
		AdminFee:         req.AdminFee,
		PhysicalCheckFee: req.PhysicalCheckFee,
		DeliveryFee:      req.DeliveryFee,
		ExpressFee:       req.ExpressFee,
		TotalCost:        totalCost,

		ServiceLevel:           req.ServiceLevel,
		IsNameTransferRequired: req.IsNameTransferRequired,

		SamsatOrigin:           req.SamsatOrigin,
		SamsatDestination:      req.SamsatDestination,

		Notes:         req.Notes,
		Status:        constants.OrderStatusPending,
		ProcessedByID: &adminID,
		OrderDate:     time.Now(),
	}
	err := u.orderRepo.Create(order)
	if err != nil {
		return "", err
	}

	// Buat log aktivitas pertama (Pesanan Dibuat)
	log := &logEntity.OrderLog{
		ID:        uuid.New().String(),
		OrderID:   newID,
		AdminID:   &adminID,
		OldStatus: "",
		NewStatus: constants.OrderStatusPending,
		Notes:     "Pesanan dibuat (oleh Admin)",
		CreatedAt: time.Now(),
	}
	_ = u.orderLogRepo.Create(log)

	return newID, nil
}

func (u *orderUsecase) CreateOrderCustomer(req dto.CreateOrderCustomerRequest) (string, error) {
	newID := uuid.New().String()

	// Default fees
	var serviceFee float64 = 0
	var deliveryFee float64 = 0
	var expressFee float64 = 0
	var adminFee float64 = 25000 // Default admin fee
	var physicalCheckFee float64 = 0

	// 1. Get Service Base Price
	srv, err := u.serviceRepo.FindByID(req.ServiceID)
	if err == nil && srv != nil {
		serviceFee = srv.ServiceFee
		
		// If service has its own express fee defined, use it as fallback
		if req.ServiceLevel == "EXPRESS" && srv.ExpressFee > 0 {
			expressFee = srv.ExpressFee
		} else if req.ServiceLevel == "EXPRESS" {
			expressFee = 50000 // default fallback
		}

		// Calculate Physical Check Fee (Rp 30.000 for Mutasi, Balik Nama, 5-Year renewal, or if explicitly requested)
		serviceNameLower := strings.ToLower(srv.ServiceName)
		categoryLower := strings.ToLower(srv.Category)
		if categoryLower == "mutasi kendaraan" || categoryLower == "balik nama" ||
			strings.Contains(serviceNameLower, "5 tahunan") || strings.Contains(serviceNameLower, "ganti plat") ||
			req.IsNameTransferRequired {
			physicalCheckFee = 30000
		}
	} else {
		// Fallback check if service is not found in DB
		if req.IsNameTransferRequired {
			physicalCheckFee = 30000
		}
	}

	// 2. Calculate Delivery Fee (Di-set 0 karena ongkir kurir eksternal dibayar pelanggan di tempat)
	if req.PickupMethod == "OJOL" || req.PickupMethod == "KURIR" || req.PickupMethod == "GOSEND" {
		deliveryFee = 0
	}

	// 3. Calculate Express Fee (handled above)
	if req.ServiceLevel == "EXPRESS" && expressFee == 0 {
		expressFee = 50000
	}

	totalCost := serviceFee + adminFee + deliveryFee + expressFee + physicalCheckFee

	order := &entity.Order{
		ID:         newID,
		CustomerID: req.CustomerID,
		VehicleID:  req.VehicleID,
		ServiceID:  req.ServiceID,

		PickupMethod:           req.PickupMethod,
		ReturnMethod:           req.ReturnMethod,
		CustomerTrackingNumber: req.CustomerTrackingNumber,
		PickupAddress:          req.PickupAddress,

		TaxAmount:        0, // Customer doesn't know the exact tax yet
		ServiceFee:       serviceFee,
		AdminFee:         adminFee,
		PhysicalCheckFee: physicalCheckFee,
		DeliveryFee:      deliveryFee,
		ExpressFee:       expressFee,
		TotalCost:        totalCost,

		ServiceLevel:           req.ServiceLevel,
		IsNameTransferRequired: req.IsNameTransferRequired,

		SamsatOrigin:           req.SamsatOrigin,
		SamsatDestination:      req.SamsatDestination,

		Notes:         req.Notes,
		Status:        constants.OrderStatusPending,
		ProcessedByID: nil, // Processed by no one initially (Wait for Admin to pick up)
		OrderDate:     time.Now(),
	}

	err = u.orderRepo.Create(order)
	if err != nil {
		return "", err
	}

	// Buat log aktivitas pertama (Pesanan Dibuat)
	// Karena ini pelanggan, kita set AdminID ke nil (sebagai penanda "Sistem/Pelanggan")
	log := &logEntity.OrderLog{
		ID:        uuid.New().String(),
		OrderID:   newID,
		AdminID:   nil,
		OldStatus: "",
		NewStatus: constants.OrderStatusPending,
		Notes:     "Pesanan dibuat (oleh Pelanggan)",
		CreatedAt: time.Now(),
	}
	_ = u.orderLogRepo.Create(log)

	return newID, nil
}

func (u *orderUsecase) GetAllOrders(page int, limit int) ([]dto.OrderResponse, error) {
	orders, err := u.orderRepo.FindAll(page, limit)
	if err != nil {
		return nil, err
	}

	var responses []dto.OrderResponse
	for _, o := range orders {
		responses = append(responses, mapOrderToResponse(o))
	}

	return responses, nil
}

func (u *orderUsecase) GetOrderByID(id string) (*dto.OrderResponse, error) {
	order, err := u.orderRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("pesanan tidak ditemukan")
	}

	resp := mapOrderToResponse(*order)

	// Fetch activity logs (history)
	logs, err := u.orderLogRepo.FindByOrderID(id)
	if err == nil && len(logs) > 0 {
		for _, l := range logs {
			resp.History = append(resp.History, dto.OrderStatusHistoryResponse{
				Status:    l.NewStatus,
				Notes:     l.Notes,
				CreatedAt: l.CreatedAt,
			})
		}
	}

	return &resp, nil
}

func (u *orderUsecase) UpdateOrderStatus(id string, req dto.UpdateOrderStatusRequest, adminID string) error {
	order, err := u.orderRepo.FindByID(id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("pesanan tidak ditemukan")
	}

	oldStatus := order.Status

	err = u.orderRepo.UpdateStatus(id, req.Status)
	if err != nil {
		return err
	}

	var adminIDPtr *string
	if adminID != "" {
		adminIDPtr = &adminID
	}

	// Create log synchronously agar ketika FE melakukan refresh data, log sudah pasti tersimpan di DB
	log := &logEntity.OrderLog{
		ID:        uuid.New().String(),
		OrderID:   id,
		AdminID:   adminIDPtr,
		OldStatus: oldStatus,
		NewStatus: req.Status,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
	}

	err = u.orderLogRepo.Create(log)
	if err != nil {
		// Log error jika insert gagal (misal karena constraint uniqueidentifier)
		println("ERROR GAGAL SAVE LOG:", err.Error())
	}

	return nil
}
