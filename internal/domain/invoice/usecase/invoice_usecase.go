package usecase

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
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
	"bj/internal/utils"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type InvoiceUsecase interface {
	GenerateInvoice(req invDto.GenerateInvoiceRequest) (*invDto.InvoiceDetailResponse, error)
	GetInvoiceByOrderID(orderID string) (*invDto.InvoiceDetailResponse, error)
	SyncInvoiceStatus(orderIDOrInvoiceNumber string) (*invDto.InvoiceDetailResponse, error)
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

// syncMidtransStatus memeriksa status transaksi secara realtime ke Midtrans Core API
func (u *invoiceUsecase) syncMidtransStatus(invoice *invEntity.Invoice) error {
	if invoice == nil || invoice.PaymentStatus == "PAID" || invoice.InvoiceNumber == "" {
		return nil
	}

	var snapEnv midtrans.EnvironmentType = midtrans.Sandbox
	if u.cfg.Midtrans.Env == "production" {
		snapEnv = midtrans.Production
	}

	client := coreapi.Client{}
	client.New(u.cfg.Midtrans.ServerKey, snapEnv)

	resp, midErr := client.CheckTransaction(invoice.InvoiceNumber)
	if midErr != nil || resp == nil {
		return fmt.Errorf("gagal cek transaksi: %w", midErr)
	}

	switch resp.TransactionStatus {
	case "settlement", "capture":
		invoice.PaymentStatus = "PAID"
		now := utils.NowWIB()
		invoice.PaidAt = &now
		if resp.PaymentType != "" {
			invoice.PaymentMethod = &resp.PaymentType
		}
		err := u.invoiceRepo.Update(invoice)
		if err != nil {
			return err
		}

		// Sinkronkan status pesanan jika masih pending
		if invoice.OrderID != "" {
			order, err := u.orderRepo.FindByID(invoice.OrderID)
			if err == nil && order != nil {
				if order.Status == "Pending" || order.Status == "PENDING" {
					// error handling pada update status order
					errStatus := u.orderRepo.UpdateStatus(invoice.OrderID, "Processing")
					if errStatus != nil {
						return fmt.Errorf("invoice lunas, tapi gagal menaikkan status order ke Processing: %w", errStatus)
					}
				}
			}
		}
	case "expire", "cancel", "deny":
		invoice.PaymentStatus = "FAILED"
		err := u.invoiceRepo.Update(invoice)
		if err != nil {
			return fmt.Errorf("gagal memperbarui data invoice gagal/expired: %w", err)
		}
	}

	return nil
}

func (u *invoiceUsecase) GenerateInvoice(req invDto.GenerateInvoiceRequest) (*invDto.InvoiceDetailResponse, error) {
	// Cek apakah invoice untuk order ini sudah ada dan sudah memiliki token pembayaran Midtrans
	existing, _ := u.invoiceRepo.GetByOrderID(req.OrderID)
	if existing != nil && existing.PaymentToken != nil && *existing.PaymentToken != "" {
		return u.GetInvoiceByOrderID(req.OrderID)
	}

	// Ambil Order
	order, err := u.orderRepo.FindByID(req.OrderID)
	if err != nil {
		return nil, errors.New("order tidak ditemukan")
	}

	// Ambil data Pelanggan, Kendaraan, dan Layanan sekaligus
	cust, err := u.custRepo.FindByID(order.CustomerID)
	if err != nil || cust == nil {
		return nil, fmt.Errorf("gagal mendapatkan data pelanggan: %w", err)
	}
	veh, _ := u.vehRepo.FindByID(order.VehicleID)
	srv, _ := u.srvRepo.FindByID(order.ServiceID)

	// Pembulatan Desimal Nominal (GrossAmt) akurat dengan math.Round
	grossAmt := int64(math.Round(order.TotalCost))
	if grossAmt <= 0 {
		return nil, errors.New("total tagihan harus lebih besar dari 0")
	}

	nowWIB := utils.NowWIB()
	var invNum string
	var invoice *invEntity.Invoice

	if existing != nil {
		invNum = existing.InvoiceNumber
		invoice = existing
	} else {
		// Generate Invoice Number (Contoh: INV-YYYYMMDD-ORDERID_SEBAGIAN)
		dateStr := nowWIB.Format("20060102")
		shortID := order.ID
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}
		invNum = fmt.Sprintf("INV-%s-%s", dateStr, shortID)

		// data invoice di database lokal terlebih dahulu (status UNPAID)
		invoice = &invEntity.Invoice{
			ID:            uuid.New().String(),
			OrderID:       order.ID,
			InvoiceNumber: invNum,
			IssueDate:     nowWIB,
			DueDate:       nowWIB.Add(24 * 7 * time.Hour), // 7 Hari
			PaymentStatus: "UNPAID",
			PaymentToken:  nil,
			PaymentUrl:    nil,
		}

		if err := u.invoiceRepo.Create(invoice); err != nil {
			return nil, fmt.Errorf("gagal menyimpan invoice: %w", err)
		}
	}

	// Hit Midtrans Snap API
	var snapEnv midtrans.EnvironmentType = midtrans.Sandbox
	if u.cfg.Midtrans.Env == "production" {
		snapEnv = midtrans.Production
	}
	snapClient := snap.Client{}
	snapClient.New(u.cfg.Midtrans.ServerKey, snapEnv)

	reqMidtrans := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  invNum,
			GrossAmt: grossAmt,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: cust.Fullname,
			Email: cust.Email,
			Phone: cust.PhoneNumber,
		},
	}

	snapResp, errSnap := snapClient.CreateTransaction(reqMidtrans)
	if errSnap != nil && strings.Contains(errSnap.Error(), "already been taken") {
		invNum = fmt.Sprintf("%s-%d", invNum, time.Now().Unix()%100000)
		invoice.InvoiceNumber = invNum
		_ = u.invoiceRepo.Update(invoice)

		reqMidtrans.TransactionDetails.OrderID = invNum
		snapResp, errSnap = snapClient.CreateTransaction(reqMidtrans)
	}

	if errSnap != nil || snapResp == nil {
		invoice.PaymentStatus = "FAILED"
		_ = u.invoiceRepo.Update(invoice)
		return nil, fmt.Errorf("Midtrans Error: %v", errSnap)
	}

	// Update Invoice dengan Token & Redirect URL dari Midtrans
	paymentToken := snapResp.Token
	paymentUrl := snapResp.RedirectURL
	invoice.PaymentToken = &paymentToken
	invoice.PaymentUrl = &paymentUrl
	invoice.PaymentStatus = "UNPAID"

	if err := u.invoiceRepo.Update(invoice); err != nil {
		return nil, fmt.Errorf("gagal memperbarui token pembayaran di database: %w", err)
	}

	// Mapping langsung di memori tanpa query ulang ke database
	return invDto.ToInvoiceDetailResponse(invoice, order, cust, veh, srv), nil
}

func (u *invoiceUsecase) GetInvoiceByOrderID(orderID string) (*invDto.InvoiceDetailResponse, error) {
	invoice, err := u.invoiceRepo.GetByOrderID(orderID)
	if err != nil || invoice == nil {
		return nil, errors.New("invoice tidak ditemukan")
	}

	// Auto sync ke Midtrans jika belum PAID
	if invoice.PaymentStatus != "PAID" {
		_ = u.syncMidtransStatus(invoice)
	}

	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order tidak ditemukan")
	}

	cust, _ := u.custRepo.FindByID(order.CustomerID)
	veh, _ := u.vehRepo.FindByID(order.VehicleID)
	srv, _ := u.srvRepo.FindByID(order.ServiceID)

	return invDto.ToInvoiceDetailResponse(invoice, order, cust, veh, srv), nil
}

// SyncInvoiceStatus melakukan sinkronisasi eksplisit dengan Midtrans API
func (u *invoiceUsecase) SyncInvoiceStatus(orderIDOrInvoiceNumber string) (*invDto.InvoiceDetailResponse, error) {
	invoice, err := u.invoiceRepo.GetByOrderID(orderIDOrInvoiceNumber)
	if err != nil || invoice == nil {
		invoice, err = u.invoiceRepo.GetByID(orderIDOrInvoiceNumber)
	}
	if invoice == nil {
		// Cari dari list berdasarkan InvoiceNumber
		invoices, _, _ := u.invoiceRepo.GetAll(1, 500)
		for _, inv := range invoices {
			if inv.InvoiceNumber == orderIDOrInvoiceNumber {
				invoice, _ = u.invoiceRepo.GetByID(inv.ID)
				break
			}
		}
	}
	if invoice == nil {
		return nil, errors.New("invoice tidak ditemukan")
	}

	_ = u.syncMidtransStatus(invoice)
	return u.GetInvoiceByOrderID(invoice.OrderID)
}

func (u *invoiceUsecase) PayInvoice(invoiceID string, req invDto.PayInvoiceRequest) error {
	invoice, err := u.invoiceRepo.GetByID(invoiceID)
	if err != nil {
		return errors.New("invoice tidak ditemukan")
	}

	if invoice.PaymentStatus == "PAID" {
		return errors.New("invoice sudah dibayar sebelumnya")
	}

	now := utils.NowWIB()
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
		paymentStatus := inv.PaymentStatus
		paidAt := inv.PaidAt
		paymentMethod := inv.PaymentMethod

		// Realtime sync ke Midtrans jika status belum PAID
		if paymentStatus != "PAID" && inv.InvoiceNumber != "" {
			if invEntityObj, err := u.invoiceRepo.GetByID(inv.ID); err == nil && invEntityObj != nil {
				_ = u.syncMidtransStatus(invEntityObj)
				paymentStatus = invEntityObj.PaymentStatus
				paidAt = invEntityObj.PaidAt
				paymentMethod = invEntityObj.PaymentMethod
			}
		}

		responseList = append(responseList, invDto.FromInvoiceWithDetails(&inv, paymentStatus, paidAt, paymentMethod))
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
		return int64(rand.Intn(99) + 1)
	}

	return available[rand.Intn(len(available))]
}

func (u *invoiceUsecase) GetAllInvoices(page int, limit int) ([]invDto.InvoiceDetailResponse, int64, error) {
	invoices, total, err := u.invoiceRepo.FindAll(page, limit)
	if err == nil {
		for i := range invoices {
			invoices[i].IssueDate = utils.ToWIB(invoices[i].IssueDate)
			invoices[i].DueDate = utils.ToWIB(invoices[i].DueDate)
			invoices[i].PaidAt = utils.ToWIBPtr(invoices[i].PaidAt)

			if invoices[i].PaymentStatus != "PAID" && invoices[i].InvoiceNumber != "" {
				if invEntityObj, err := u.invoiceRepo.GetByID(invoices[i].InvoiceID); err == nil && invEntityObj != nil {
					_ = u.syncMidtransStatus(invEntityObj)
					invoices[i].PaymentStatus = invEntityObj.PaymentStatus
					invoices[i].PaidAt = utils.ToWIBPtr(invEntityObj.PaidAt)
					invoices[i].PaymentMethod = invEntityObj.PaymentMethod
				}
			}
		}
	}
	return invoices, total, err
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
	invoices, _, err := u.invoiceRepo.GetAll(1, 1000)
	if err != nil {
		return err
	}
	var targetInvoice *invEntity.InvoiceWithDetails
	for _, inv := range invoices {
		if inv.InvoiceNumber == orderID || (inv.InvoiceNumber != "" && (strings.HasPrefix(orderID, inv.InvoiceNumber) || strings.HasPrefix(inv.InvoiceNumber, orderID))) {
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

	now := utils.NowWIB()
	switch transactionStatus {
	case "capture", "settlement":
		invToUpdate.PaymentStatus = "PAID"
		paymentType, _ := payload["payment_type"].(string)
		invToUpdate.PaymentMethod = &paymentType
		invToUpdate.PaidAt = &now
	case "cancel", "deny", "expire":
		invToUpdate.PaymentStatus = "FAILED"
	}

	return u.invoiceRepo.Update(invToUpdate)
}
