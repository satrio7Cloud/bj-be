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
	custRepo "bj/internal/domain/customer/repository"
	"bj/internal/domain/order/dto"
	"bj/internal/domain/order/repository"
	logEntity "bj/internal/domain/orderlog/entity"
	logRepo "bj/internal/domain/orderlog/repository"
	srvRepo "bj/internal/domain/service/repository"
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

func newOrderLog(orderID string, adminID *string, oldStatus, newStatus, notes string) *logEntity.OrderLog {
	return &logEntity.OrderLog{
		ID:        uuid.New().String(),
		OrderID:   orderID,
		AdminID:   adminID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Notes:     notes,
		CreatedAt: time.Now(),
	}
}

func (u *orderUsecase) CreateOrderAdmin(req dto.CreateOrderAdminRequest, adminID string) (string, error) {
	newID := uuid.New().String()
	order := req.ToEntity(newID, adminID)
	if err := u.orderRepo.Create(order); err != nil {
		return "", err
	}

	_ = u.orderLogRepo.Create(newOrderLog(newID, &adminID, "", constants.OrderStatusPending, "Pesanan dibuat (oleh Admin)"))
	return newID, nil
}

func (u *orderUsecase) CreateOrderCustomer(req dto.CreateOrderCustomerRequest) (string, string, error) {
	newID := uuid.New().String()

	// 1. Resolve Customer (Pakai ID jika ada, atau Cari / Buat otomatis dari data diri)
	customerID := req.CustomerID
	if customerID == "" {
		if req.Fullname == "" {
			return "", "", errors.New("nama lengkap pelanggan wajib diisi")
		}
		if req.PhoneNumber == "" && req.NIK == "" {
			return "", "", errors.New("nomor telepon atau NIK pelanggan wajib diisi")
		}

		existingCust, err := u.customerRepo.FindByPhoneOrNIKAndName(req.PhoneNumber, req.NIK, req.Fullname)
		if err != nil {
			return "", "", err
		}

		if existingCust != nil {
			customerID = existingCust.ID
		} else {
			newCust := req.ToCustomerEntity(uuid.New().String())
			if err := u.customerRepo.Create(newCust); err != nil {
				return "", "", err
			}
			customerID = newCust.ID
		}
	}

	// 2. Resolve Vehicle (Pakai ID jika ada, atau Cari / Buat otomatis dari data kendaraan)
	vehicleID := req.VehicleID
	if vehicleID == "" {
		if req.PlateNumber == "" {
			return "", "", errors.New("plat nomor kendaraan wajib di isi")
		}

		existingVehicle, err := u.vehicleRepo.FindByPlateNumber(req.PlateNumber)
		if err != nil {
			return "", "", err
		}

		if existingVehicle != nil {
			vehicleID = existingVehicle.ID
			if existingVehicle.CustomerID != customerID {
				existingVehicle.CustomerID = customerID
				_ = u.vehicleRepo.Update(existingVehicle)
			}
		} else {
			if req.ModelID == "" {
				return "", "", errors.New("model kendaraan wajib dipilih")
			}
			newVehicle := req.ToVehicleEntity(uuid.New().String(), customerID)
			if err := u.vehicleRepo.Create(newVehicle); err != nil {
				return "", "", err
			}
			vehicleID = newVehicle.ID
		}
	}

	// 3. Kalkulasi Biaya & Layanan
	srv, err := u.serviceRepo.FindByID(req.ServiceID)
	if err != nil {
		return "", "", err
	}
	if srv == nil {
		return "", "", errors.New("layanan tidak ditemukan")
	}

	serviceName := srv.ServiceName
	var serviceFee = srv.ServiceFee
	var adminFee float64 = 0
	var physicalCheckFee float64 = 0
	var expressFee float64 = 0
	var deliveryFee float64 = 0

	if req.ServiceLevel == "EXPRESS" && srv.ExpressFee > 0 {
		expressFee = srv.ExpressFee
	} else if req.ServiceLevel == "EXPRESS" {
		expressFee = 50000
	}

	if srv.AdminFee > 0 {
		adminFee = srv.AdminFee
	}

	if srv.PhysicalCheckFee > 0 {
		physicalCheckFee = srv.PhysicalCheckFee
	} else {
		serviceNameLower := strings.ToLower(srv.ServiceName)
		categoryLower := strings.ToLower(srv.Category)
		if categoryLower == "mutasi kendaraan" || categoryLower == "balik nama" ||
			strings.Contains(serviceNameLower, "5 tahunan") || strings.Contains(serviceNameLower, "ganti plat") ||
			req.IsNameTransferRequired {
			physicalCheckFee = 30000
		}
	}

	// Kalkulasi Biaya Pengiriman via Biteship
	pickupUpper := strings.ToUpper(req.PickupMethod)
	if pickupUpper == "WALK_IN" || pickupUpper == "DATANG_SENDIRI" {
		deliveryFee = 0
	} else if req.CourierName != "" && req.CourierService != "" {
		destPostalCode := req.PostalCode
		if destPostalCode == "" {
			destPostalCode = req.DestinationPostalCode
		}
		if destPostalCode == "" {
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

		if deliveryFee == 0 && req.DeliveryFee > 0 {
			deliveryFee = req.DeliveryFee
		}
	} else if req.DeliveryFee > 0 {
		deliveryFee = req.DeliveryFee
	}

	fees := dto.OrderFees{
		ServiceFee:       serviceFee,
		AdminFee:         adminFee,
		PhysicalCheckFee: physicalCheckFee,
		DeliveryFee:      deliveryFee,
		ExpressFee:       expressFee,
		TotalCost:        serviceFee + adminFee + deliveryFee + expressFee + physicalCheckFee,
	}

	order := req.ToOrderEntity(newID, customerID, vehicleID, fees)
	if err := u.orderRepo.Create(order); err != nil {
		return "", "", err
	}

	_ = u.orderLogRepo.Create(newOrderLog(newID, nil, "", constants.OrderStatusPending, "Pesanan dibuat (oleh Pelanggan)"))
	return newID, serviceName, nil
}

func (u *orderUsecase) GetAllOrders(page int, limit int) ([]dto.OrderResponse, error) {
	orders, err := u.orderRepo.FindAll(page, limit)
	if err != nil {
		return nil, err
	}
	return dto.ToOrderResponses(orders), nil
}

func (u *orderUsecase) GetOrderByID(id string) (*dto.OrderResponse, error) {
	order, err := u.orderRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("pesanan tidak ditemukan")
	}

	resp := dto.ToOrderResponse(*order)

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

	log := newOrderLog(id, adminIDPtr, oldStatus, req.Status, req.Notes)
	if err := u.orderLogRepo.Create(log); err != nil {
		println("ERROR GAGAL SAVE LOG:", err.Error())
	}

	return nil
}

func (u *orderUsecase) CreateOrderKTL(ctx context.Context, req dto.CreateOrderKTLRequest, fileKtl *multipart.FileHeader) (dto.OrderKTLResponse, error) {
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
		println("ERROR GAGAL CREATE SNAP MIDTRANS:", snapErr.GetMessage())
		return dto.OrderKTLResponse{}, errors.New("gagal membuat transaksi pembayaran")
	} else {
		snapToken = snapResp.Token
	}

	orderKTL := req.ToEntity(newID, buktiKTLPath)
	if err := u.orderRepo.CreateOrderKTL(orderKTL); err != nil {
		return dto.OrderKTLResponse{}, err
	}

	return dto.ToOrderKTLResponse(*orderKTL, snapToken), nil
}

func (u *orderUsecase) GetAllOrderKTLs(page int, limit int) ([]dto.OrderKTLResponse, error) {
	orders, err := u.orderRepo.FindAllOrderKTL(page, limit)
	if err != nil {
		return nil, err
	}
	return dto.ToOrderKTLResponses(orders), nil
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

	log := newOrderLog(id, adminIDPtr, "", req.Status, req.Notes)
	if err := u.orderLogRepo.Create(log); err != nil {
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

	resp := dto.ToOrderKTLResponse(*order)
	return &resp, nil
}
