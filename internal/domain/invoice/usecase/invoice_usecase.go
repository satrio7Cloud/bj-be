package usecase

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"crypto/sha512"
	"encoding/hex"

	"bj/internal/config"
	"bj/internal/domain/customer/repository"
	invDto "bj/internal/domain/invoice/dto"
	invEntity "bj/internal/domain/invoice/entity"
	invRepo "bj/internal/domain/invoice/repository"
	orderRepo "bj/internal/domain/order/repository"
	srvRepo "bj/internal/domain/service/repository"
	vehRepo "bj/internal/domain/vehicle/repository"
	"bj/internal/pkg/qris"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type InvoiceUsecase interface {
	GenerateInvoice(req invDto.GenerateInvoiceRequest) (*invDto.InvoiceDetailResponse, error)
	GetInvoiceByOrderID(orderID string) (*invDto.InvoiceDetailResponse, error)
	PayInvoice(invoiceID string, req invDto.PayInvoiceRequest) error
	GetAll(page, limit int) (map[string]interface{}, error)
	GetAllInvoices(page int, limit int) ([]invDto.InvoiceDetailResponse, int64, error)
	HandleMidtransNotification(payload map[string]interface{}) error
}

type invoiceUsecase struct {
	invoiceRepo invRepo.InvoiceRepository
	orderRepo   orderRepo.OrderRepository
	custRepo    repository.CustomerRepository
	vehRepo     vehRepo.VehicleRepository
	srvRepo     srvRepo.ServiceRepository
	qrisGen     qris.QrisGenerator
	cfg         *config.Config
}

func NewInvoiceUsecase(
	ir invRepo.InvoiceRepository,
	or orderRepo.OrderRepository,
	cr repository.CustomerRepository,
	vr vehRepo.VehicleRepository,
	sr srvRepo.ServiceRepository,
	qg qris.QrisGenerator,
	cfg *config.Config,
) InvoiceUsecase {
	return &invoiceUsecase{
		invoiceRepo: ir,
		orderRepo:   or,
		custRepo:    cr,
		vehRepo:     vr,
		srvRepo:     sr,
		qrisGen:     qg,
		cfg:         cfg,
	}
}

func (u *invoiceUsecase) GenerateInvoice(req invDto.GenerateInvoiceRequest) (*invDto.InvoiceDetailResponse, error) {
	// 1. Cek apakah invoice untuk order ini sudah ada
	existing, _ := u.invoiceRepo.GetByOrderID(req.OrderID)
	if existing != nil {
		return u.GetInvoiceByOrderID(req.OrderID)
	}

	// 2. Ambil Order
	order, err := u.orderRepo.FindByID(req.OrderID)
	if err != nil {
		return nil, errors.New("order tidak ditemukan")
	}

	// 3. Generate Invoice Number (Contoh: INV-YYYYMMDD-ORDERID_SEBAGIAN)
	dateStr := time.Now().Format("20060102")
	shortID := order.ID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}
	invNum := fmt.Sprintf("INV-%s-%s", dateStr, shortID)

	// 4. Hit Midtrans Snap API
	var snapEnv midtrans.EnvironmentType = midtrans.Sandbox
	if u.cfg.Midtrans.Env == "production" {
		snapEnv = midtrans.Production
	}
	snapClient := snap.Client{}
	snapClient.New(u.cfg.Midtrans.ServerKey, snapEnv)

	reqMidtrans := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  invNum,
			GrossAmt: int64(order.TotalCost),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: "Pelanggan",
			Email: "no-reply@birojasa.com",
			Phone: "0800000000",
		},
	}

	if cust, err := u.custRepo.FindByID(order.CustomerID); err == nil {
		reqMidtrans.CustomerDetail.FName = cust.Fullname
		reqMidtrans.CustomerDetail.Phone = cust.PhoneNumber
		reqMidtrans.CustomerDetail.Email = cust.Email
	}

	snapResp, errSnap := snapClient.CreateTransaction(reqMidtrans)

	var paymentToken *string
	var paymentUrl *string
	if errSnap == nil && snapResp != nil {
		paymentToken = &snapResp.Token
		paymentUrl = &snapResp.RedirectURL
	} else {
		return nil, fmt.Errorf("Midtrans Error: %v", errSnap)
	}

	// 5. Buat Entity Invoice
	invoice := &invEntity.Invoice{
		ID:            uuid.New().String(),
		OrderID:       order.ID,
		InvoiceNumber: invNum,
		IssueDate:     time.Now(),
		DueDate:       time.Now().Add(24 * 7 * time.Hour), // 7 Hari
		PaymentStatus: "UNPAID",
		PaymentToken:  paymentToken,
		PaymentUrl:    paymentUrl,
	}

	err = u.invoiceRepo.Create(invoice)
	if err != nil {
		return nil, err
	}

	return u.GetInvoiceByOrderID(req.OrderID)
}

func (u *invoiceUsecase) GetInvoiceByOrderID(orderID string) (*invDto.InvoiceDetailResponse, error) {
	invoice, err := u.invoiceRepo.GetByOrderID(orderID)
	if err != nil || invoice == nil {
		return nil, errors.New("invoice tidak ditemukan")
	}

	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order tidak ditemukan")
	}

	cust, _ := u.custRepo.FindByID(order.CustomerID)
	veh, _ := u.vehRepo.FindByID(order.VehicleID)
	srv, _ := u.srvRepo.FindByID(order.ServiceID)

	resp := &invDto.InvoiceDetailResponse{
		InvoiceID:        invoice.ID,
		OrderID:          invoice.OrderID,
		InvoiceNumber:    invoice.InvoiceNumber,
		IssueDate:        invoice.IssueDate,
		DueDate:          invoice.DueDate,
		PaymentStatus:    invoice.PaymentStatus,
		PaymentMethod:    invoice.PaymentMethod,
		PaidAt:           invoice.PaidAt,
		PaymentToken:     invoice.PaymentToken,
		PaymentUrl:       invoice.PaymentUrl,
		
		QrisPayload:      invoice.QrisPayload,
		QrisAmount:       invoice.QrisAmount,

		TaxAmount:        order.TaxAmount,
		ServiceFee:       order.ServiceFee,
		AdminFee:         order.AdminFee,
		PhysicalCheckFee: order.PhysicalCheckFee,
		DeliveryFee:      order.DeliveryFee,
		ExpressFee:       order.ExpressFee,
		TotalCost:        order.TotalCost,
	}

	if cust != nil {
		resp.CustomerName = cust.Fullname
		resp.CustomerPhone = cust.PhoneNumber
	}
	if veh != nil {
		resp.VehiclePlate = veh.PlateNumber
	}
	if srv != nil {
		resp.ServiceName = srv.ServiceName
	}

	return resp, nil
}

func (u *invoiceUsecase) PayInvoice(invoiceID string, req invDto.PayInvoiceRequest) error {
	invoice, err := u.invoiceRepo.GetByID(invoiceID)
	if err != nil {
		return errors.New("invoice tidak ditemukan")
	}

	if invoice.PaymentStatus == "PAID" {
		return errors.New("invoice sudah dibayar sebelumnya")
	}

	now := time.Now()
	invoice.PaymentStatus = "PAID"
	invoice.PaymentMethod = &req.PaymentMethod
	invoice.PaidAt = &now

	return u.invoiceRepo.Update(invoice)
}

func (u *invoiceUsecase) GetAll(page, limit int) (map[string]interface{}, error) {
	invoices, total, err := u.invoiceRepo.GetAll(page, limit)
	if err != nil {
		return nil, errors.New("terjadi kendala saat mengambil data invoice dari database")
	}

	var responseList []invDto.InvoiceDetailResponse
	for _, inv := range invoices {
		responseList = append(responseList, invDto.InvoiceDetailResponse{
			InvoiceID:        inv.ID,
			OrderID:          inv.OrderID,
			InvoiceNumber:    inv.InvoiceNumber,
			IssueDate:        inv.IssueDate,
			DueDate:          inv.DueDate,
			PaymentStatus:    inv.PaymentStatus,
			PaymentMethod:    inv.PaymentMethod,
			PaidAt:           inv.PaidAt,
			QrisPayload:      inv.QrisPayload,
			QrisAmount:       inv.QrisAmount,
			CustomerName:     inv.CustomerName,
			CustomerEmail:    inv.CustomerEmail,
			CustomerPhone:    inv.CustomerPhone,
			VehiclePlate:     inv.VehiclePlate,
			ServiceName:      inv.ServiceName,
			TotalCost:        inv.TotalCost,
		})
	}

	return map[string]interface{}{
		"data":  responseList,
		"total": total,
		"page":  page,
		"limit": limit,
	}, nil
}

func randomUniqueCode(baseAmount float64, activeAmounts []float64) int64 {
	used := make(map[int64]bool)
	for _, qrisAmt := range activeAmounts {
		code := int64(qrisAmt - baseAmount)
		if code >= 1 && code <= 99 {
			used[code] = true
		}
	}

	var available []int64
	for code := int64(1); code <= 99; code++ {
		if !used[code] {
			available = append(available, code)
		}
	}

	if len(available) == 0 {
		// Fallback to standard random if pool is full (should not happen in ordinary cases)
		return int64(rand.Intn(99) + 1)
	}

	// Select a random available code
	return available[rand.Intn(len(available))]
}

func (u *invoiceUsecase) GetAllInvoices(page int, limit int) ([]invDto.InvoiceDetailResponse, int64, error) {
	return u.invoiceRepo.FindAll(page, limit)
}

func (u *invoiceUsecase) HandleMidtransNotification(payload map[string]interface{}) error {
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return errors.New("order_id missing from payload")
	}

	statusCode, _ := payload["status_code"].(string)
	grossAmount, _ := payload["gross_amount"].(string)
	signatureKey, _ := payload["signature_key"].(string)
	transactionStatus, _ := payload["transaction_status"].(string)

	// Verifikasi signature_key = SHA512(order_id + status_code + gross_amount + server_key)
	rawString := orderID + statusCode + grossAmount + u.cfg.Midtrans.ServerKey
	hash := sha512.New()
	hash.Write([]byte(rawString))
	expectedSignature := hex.EncodeToString(hash.Sum(nil))

	if signatureKey != expectedSignature {
		return errors.New("invalid signature key")
	}

	// Cari invoice berdasarkan InvoiceNumber = order_id di Midtrans
	invoices, _, err := u.invoiceRepo.GetAll(1, 1000) // atau buat method get by InvoiceNumber
	if err != nil {
		return err
	}
	var targetInvoice *invEntity.InvoiceWithDetails
	for _, inv := range invoices {
		if inv.InvoiceNumber == orderID {
			targetInvoice = &inv
			break
		}
	}

	if targetInvoice == nil {
		return errors.New("invoice not found")
	}

	invToUpdate, err := u.invoiceRepo.GetByID(targetInvoice.ID)
	if err != nil {
		return err
	}

	now := time.Now()
	if transactionStatus == "capture" || transactionStatus == "settlement" {
		invToUpdate.PaymentStatus = "PAID"
		paymentType, _ := payload["payment_type"].(string)
		invToUpdate.PaymentMethod = &paymentType
		invToUpdate.PaidAt = &now
	} else if transactionStatus == "cancel" || transactionStatus == "deny" || transactionStatus == "expire" {
		invToUpdate.PaymentStatus = "FAILED"
	}

	return u.invoiceRepo.Update(invToUpdate)
}
