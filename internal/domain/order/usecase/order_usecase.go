package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"regexp"
	"strings"
	"time"

	"bj/internal/constants"
	biteshipUC "bj/internal/domain/biteship/usecase"
	custEntity "bj/internal/domain/customer/entity"
	custRepo "bj/internal/domain/customer/repository"
	"bj/internal/domain/order/dto"
	"bj/internal/domain/order/entity"
	"bj/internal/domain/order/repository"
	logEntity "bj/internal/domain/orderlog/entity"
	logRepo "bj/internal/domain/orderlog/repository"
	srvRepo "bj/internal/domain/service/repository"
	vehEntity "bj/internal/domain/vehicle/entity"
	vehRepo "bj/internal/domain/vehicle/repository"
	"bj/internal/shared/storage"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type OrderUsecase interface {
	CreateOrderAdmin(req dto.CreateOrderAdminRequest, adminID string) (string, error)
	CreateOrderCustomer(req dto.CreateOrderCustomerRequest) (string, string, error)
	GetAllOrders(page int, limit int) ([]dto.OrderResponse, error)
	GetOrderByID(id string) (*dto.OrderResponse, error)
	UpdateOrderStatus(id string, req dto.UpdateOrderStatusRequest, adminID string) error
	CreateOrderKTL(ctx context.Context, req dto.CreateOrderKTLRequest, fileKtl *multipart.FileHeader) (dto.OrderKTLResponse, error)
	GetAllOrderKTLs(page int, limit int) ([]dto.OrderKTLResponse, error)
	GetOrderKTLByID(id string) (*dto.OrderKTLResponse, error)
	UpdateOrderKTLStatus(id string, req dto.UpdateOrderStatusRequest, adminID string) error
}

type orderUsecase struct {
	orderRepo      repository.OrderRepository
	orderLogRepo   logRepo.OrderLogRepository
	serviceRepo    srvRepo.ServiceRepository
	customerRepo   custRepo.CustomerRepository
	vehicleRepo    vehRepo.VehicleRepository
	biteshipUC     biteshipUC.BiteshipUsecase
	storageService storage.StorageService
	snapClient     snap.Client
}

func NewOrderUsecase(
	repo repository.OrderRepository,
	logRepo logRepo.OrderLogRepository,
	serviceRepo srvRepo.ServiceRepository,
	customerRepo custRepo.CustomerRepository,
	vehicleRepo vehRepo.VehicleRepository,
	biteshipUC biteshipUC.BiteshipUsecase,
	storageService storage.StorageService,
	snapClient snap.Client,
) OrderUsecase {
	return &orderUsecase{
		orderRepo:      repo,
		orderLogRepo:   logRepo,
		serviceRepo:    serviceRepo,
		customerRepo:   customerRepo,
		vehicleRepo:    vehicleRepo,
		biteshipUC:     biteshipUC,
		storageService: storageService,
		snapClient:     snapClient,
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

		SamsatOrigin:      o.SamsatOrigin,
		SamsatDestination: o.SamsatDestination,

		Notes:     o.Notes,
		OrderDate: o.OrderDate,

		TrackingNumber: o.TrackingNumber,
		CourierName:    o.CourierName,
		CourierService: o.CourierService,
		ShippingCost:   o.ShippingCost,
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

		SamsatOrigin:      req.SamsatOrigin,
		SamsatDestination: req.SamsatDestination,

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

func (u *orderUsecase) CreateOrderCustomer(req dto.CreateOrderCustomerRequest) (string, string, error) {
	newID := uuid.New().String()

	// 0. Resolve Customer (Pakai ID jika ada, atau Cari / Buat otomatis dari data diri)
	customerID := req.CustomerID
	if customerID == "" {
		if req.PhoneNumber == "" && req.NIK == "" {
			return "", "", errors.New("nomor telepon atau NIK pelanggan wajib diisi jika customer_id tidak disertakan")
		}

		existingCust, err := u.customerRepo.FindByPhoneOrNIK(req.PhoneNumber, req.NIK)
		if err != nil {
			return "", "", err
		}

		if existingCust != nil {
			customerID = existingCust.ID
			// Update informasi jika ada pembaruan data
			if req.Fullname != "" && req.Fullname != existingCust.Fullname {
				existingCust.Fullname = req.Fullname
			}
			if req.Address != "" && req.Address != existingCust.Address {
				existingCust.Address = req.Address
			}
			if req.Email != "" && req.Email != existingCust.Email {
				existingCust.Email = req.Email
			}
			_ = u.customerRepo.Update(existingCust)
		} else {
			if req.Fullname == "" {
				return "", "", errors.New("nama lengkap pelanggan wajib diisi untuk pendaftaran pelanggan baru")
			}
			newCust := &custEntity.Customer{
				ID:          uuid.New().String(),
				Fullname:    req.Fullname,
				PhoneNumber: req.PhoneNumber,
				Address:     req.Address,
				NIK:         req.NIK,
				Email:       req.Email,
				CreatedAt:   time.Now(),
			}
			if err := u.customerRepo.Create(newCust); err != nil {
				return "", "", err
			}
			customerID = newCust.ID
		}
	}

	// 0.1 Resolve Vehicle (Pakai ID jika ada, atau Cari / Buat otomatis dari data kendaraan)
	vehicleID := req.VehicleID
	if vehicleID == "" {
		if req.PlateNumber == "" {
			return "", "", errors.New("plat nomor kendaraan wajib diisi jika vehicle_id tidak disertakan")
		}

		existingVehicle, err := u.vehicleRepo.FindByPlateNumber(req.PlateNumber)
		if err != nil {
			return "", "", err
		}

		if existingVehicle != nil {
			vehicleID = existingVehicle.ID
			// Update relasi pemilik jika kendaraan ini sekarang dipesan atas nama customer ini
			if existingVehicle.CustomerID != customerID {
				existingVehicle.CustomerID = customerID
				_ = u.vehicleRepo.Update(existingVehicle)
			}
		} else {
			if req.ModelID == "" {
				return "", "", errors.New("model kendaraan wajib dipilih untuk mendaftarkan kendaraan baru")
			}
			year := req.Year
			if year == 0 {
				year = time.Now().Year()
			}
			newVehicle := &vehEntity.Vehicle{
				ID:          uuid.New().String(),
				CustomerID:  customerID,
				ModelID:     req.ModelID,
				PlateNumber: req.PlateNumber,
				Year:        year,
			}
			if err := u.vehicleRepo.Create(newVehicle); err != nil {
				return "", "", err
			}
			vehicleID = newVehicle.ID
		}
	}

	// Default fees
	var serviceFee float64 = 0
	var deliveryFee float64 = 0
	var expressFee float64 = 0
	var adminFee float64 = 0
	var physicalCheckFee float64 = 0
	var serviceName string

	// 1. Get Service Base Price
	srv, err := u.serviceRepo.FindByID(req.ServiceID)
	if err != nil {
		return "", "", err
	}
	if srv == nil {
		return "", "", errors.New("layanan tidak ditemukan")
	}

	serviceName = srv.ServiceName
	serviceFee = srv.ServiceFee

	// If service has its own express fee defined, use it as fallback
	if req.ServiceLevel == "EXPRESS" && srv.ExpressFee > 0 {
		expressFee = srv.ExpressFee
	} else if req.ServiceLevel == "EXPRESS" {
		expressFee = 50000 // default fallback
	}

	// Use AdminFee from Service DB if it is set (> 0)
	if srv.AdminFee > 0 {
		adminFee = srv.AdminFee
	}

	// Use PhysicalCheckFee from Service DB if it is set (> 0)
	if srv.PhysicalCheckFee > 0 {
		physicalCheckFee = srv.PhysicalCheckFee
	} else {
		// Fallback: Calculate Physical Check Fee (Rp 30.000 for Mutasi, Balik Nama, 5-Year renewal, or if explicitly requested)
		serviceNameLower := strings.ToLower(srv.ServiceName)
		categoryLower := strings.ToLower(srv.Category)
		if categoryLower == "mutasi kendaraan" || categoryLower == "balik nama" ||
			strings.Contains(serviceNameLower, "5 tahunan") || strings.Contains(serviceNameLower, "ganti plat") ||
			req.IsNameTransferRequired {
			physicalCheckFee = 30000
		}
	}

	// 2. Calculate Delivery Fee (Kalkulasi Otomatis via Biteship & Skenario Pengiriman)
	pickupUpper := strings.ToUpper(req.PickupMethod)
	if pickupUpper == "WALK_IN" || pickupUpper == "DATANG_SENDIRI" {
		deliveryFee = 0
	} else if req.CourierName != "" && req.CourierService != "" {
		destPostalCode := req.PostalCode
		if destPostalCode == "" {
			destPostalCode = req.DestinationPostalCode
		}
		if destPostalCode == "" {
			// Otomatis cari 5 digit kodepos di string alamat jika tidak dikirim spesifik
			fullAddress := req.Address
			if req.PickupAddress != nil && *req.PickupAddress != "" {
				fullAddress = *req.PickupAddress
			}
			re := regexp.MustCompile(`\b\d{5}\b`)
			destPostalCode = re.FindString(fullAddress)
		}

		if destPostalCode != "" && u.biteshipUC != nil {
			calculatedFee, err := u.biteshipUC.CalculateShippingFee(req.CourierName, req.CourierService, destPostalCode, req.PickupMethod)
			if err == nil && calculatedFee > 0 {
				deliveryFee = calculatedFee
			}
		}

		// Fallback jika API Biteship gagal atau kodepos tidak terdeteksi
		if deliveryFee == 0 && req.DeliveryFee > 0 {
			deliveryFee = req.DeliveryFee
		}
	} else if req.DeliveryFee > 0 {
		deliveryFee = req.DeliveryFee
	}

	// 3. Calculate Express Fee (handled above)
	if req.ServiceLevel == "EXPRESS" && expressFee == 0 {
		expressFee = 50000
	}

	totalCost := serviceFee + adminFee + deliveryFee + expressFee + physicalCheckFee

	var courierName *string
	if req.CourierName != "" {
		courierName = &req.CourierName
	}
	var courierService *string
	if req.CourierService != "" {
		courierService = &req.CourierService
	}

	pickupAddress := req.PickupAddress
	if (pickupAddress == nil || *pickupAddress == "") && req.Address != "" {
		pickupAddress = &req.Address
	}

	order := &entity.Order{
		ID:         newID,
		CustomerID: customerID,
		VehicleID:  vehicleID,
		ServiceID:  req.ServiceID,

		PickupMethod:           req.PickupMethod,
		ReturnMethod:           req.ReturnMethod,
		CustomerTrackingNumber: req.CustomerTrackingNumber,
		PickupAddress:          pickupAddress,

		TaxAmount:        0, // Customer doesn't know the exact tax yet
		ServiceFee:       serviceFee,
		AdminFee:         adminFee,
		PhysicalCheckFee: physicalCheckFee,
		DeliveryFee:      deliveryFee,
		ExpressFee:       expressFee,
		TotalCost:        totalCost,

		CourierName:    courierName,
		CourierService: courierService,
		ShippingCost:   deliveryFee,

		ServiceLevel:           req.ServiceLevel,
		IsNameTransferRequired: req.IsNameTransferRequired,

		SamsatOrigin:      req.SamsatOrigin,
		SamsatDestination: req.SamsatDestination,

		Notes:         req.Notes,
		Status:        constants.OrderStatusPending,
		ProcessedByID: nil, // Processed by no one initially (Wait for Admin to pick up)
		OrderDate:     time.Now(),
	}

	err = u.orderRepo.Create(order)
	if err != nil {
		return "", "", err
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

	return newID, serviceName, nil
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

func (u *orderUsecase) CreateOrderKTL(ctx context.Context, req dto.CreateOrderKTLRequest, fileKtl *multipart.FileHeader) (dto.OrderKTLResponse, error) {
	// Validasi input
	if req.Fee <= 0 {
		return dto.OrderKTLResponse{}, errors.New("fee harus lebih dari 0")
	}

	validRegions := map[string]bool{
		"Jakarta":   true,
		"Tangerang": true,
		"Depok":     true,
	}
	if !validRegions[req.Region] {
		return dto.OrderKTLResponse{}, errors.New("region tidak didukung")
	}

	if req.VehicleType != "Motor" && req.VehicleType != "Mobil" {
		return dto.OrderKTLResponse{}, errors.New("vehicle type harus Motor atau Mobil")
	}

	var buktiKTLPath *string

	if fileKtl != nil && u.storageService != nil {
		path, err := u.storageService.UploadFile(ctx, fileKtl, "cetak-ktl")
		if err == nil {
			buktiKTLPath = &path
		} else {
			println("Warning: gagal upload bukti ktl -", err.Error())
		}
	}

	newID := uuid.New().String()

	// Integrasi Midtrans Snap
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  newID,
			GrossAmt: int64(req.Fee),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.CustomerName,
			Phone: req.PhoneNumber,
		},
	}

	snapResp, snapErr := u.snapClient.CreateTransaction(snapReq)
	var snapToken string
	if snapErr != nil {
		// Log error but we might still want to proceed or fail
		println("ERROR GAGAL CREATE SNAP MIDTRANS:", snapErr.GetMessage())
		return dto.OrderKTLResponse{}, errors.New("gagal membuat transaksi pembayaran")
	} else {
		snapToken = snapResp.Token
	}

	orderKTL := &entity.OrderKTL{
		ID:                  newID,
		CustomerName:        req.CustomerName,
		PhoneNumber:         req.PhoneNumber,
		VehicleType:         req.VehicleType,
		Region:              req.Region,
		Fee:                 req.Fee,
		MerchantRef:         nil, // Null
		Status:              "Menunggu Pembayaran",
		BuktiKtlPath:        buktiKTLPath,
		BuktiPembayaranPath: nil, // Null
		OrderDate:           time.Now(),
	}

	err := u.orderRepo.CreateOrderKTL(orderKTL)
	if err != nil {
		return dto.OrderKTLResponse{}, err
	}

	resp := dto.OrderKTLResponse{
		ID:                  orderKTL.ID,
		CustomerName:        orderKTL.CustomerName,
		PhoneNumber:         orderKTL.PhoneNumber,
		VehicleType:         orderKTL.VehicleType,
		Region:              orderKTL.Region,
		Fee:                 orderKTL.Fee,
		MerchantRef:         orderKTL.MerchantRef,
		Status:              orderKTL.Status,
		BuktiKtlPath:        orderKTL.BuktiKtlPath,
		BuktiPembayaranPath: orderKTL.BuktiPembayaranPath,
		SnapToken:           snapToken,
		OrderDate:           orderKTL.OrderDate,
	}

	return resp, nil
}

func (u *orderUsecase) GetAllOrderKTLs(page int, limit int) ([]dto.OrderKTLResponse, error) {
	orders, err := u.orderRepo.FindAllOrderKTL(page, limit)
	if err != nil {
		return nil, err
	}

	var responses []dto.OrderKTLResponse
	for _, o := range orders {
		responses = append(responses, dto.OrderKTLResponse{
			ID:                  o.ID,
			CustomerName:        o.CustomerName,
			PhoneNumber:         o.PhoneNumber,
			VehicleType:         o.VehicleType,
			Region:              o.Region,
			Fee:                 o.Fee,
			MerchantRef:         o.MerchantRef,
			Status:              o.Status,
			BuktiKtlPath:        o.BuktiKtlPath,
			BuktiPembayaranPath: o.BuktiPembayaranPath,
			OrderDate:           o.OrderDate,
		})
	}

	return responses, nil
}

func (u *orderUsecase) UpdateOrderKTLStatus(id string, req dto.UpdateOrderStatusRequest, adminID string) error {
	err := u.orderRepo.UpdateKTLStatus(id, req.Status)
	if err != nil {
		return err
	}

	var adminIDPtr *string
	if adminID != "" {
		adminIDPtr = &adminID
	}

	log := &logEntity.OrderLog{
		ID:        uuid.New().String(),
		OrderID:   id,
		AdminID:   adminIDPtr,
		OldStatus: "",
		NewStatus: req.Status,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
	}

	err = u.orderLogRepo.Create(log)
	if err != nil {
		println("ERROR GAGAL SAVE LOG KTL:", err.Error())
	}

	return nil
}

func (u *orderUsecase) GetOrderKTLByID(id string) (*dto.OrderKTLResponse, error) {
	order, err := u.orderRepo.FindOrderKTLByID(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("pesanan KTL tidak ditemukan")
	}

	resp := &dto.OrderKTLResponse{
		ID:                  order.ID,
		CustomerName:        order.CustomerName,
		PhoneNumber:         order.PhoneNumber,
		VehicleType:         order.VehicleType,
		Region:              order.Region,
		Fee:                 order.Fee,
		MerchantRef:         order.MerchantRef,
		Status:              order.Status,
		BuktiKtlPath:        order.BuktiKtlPath,
		BuktiPembayaranPath: order.BuktiPembayaranPath,
		OrderDate:           order.OrderDate,
	}

	return resp, nil
}
