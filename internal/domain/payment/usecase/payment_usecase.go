package usecase

import (
	"errors"

	"bj/internal/constants"
	orderDto "bj/internal/domain/order/dto"
	orderUsecase "bj/internal/domain/order/usecase"
)

type PaymentUsecase interface {
	HandleWebhookMidtrans(payload map[string]interface{}) error
}

type paymentUsecase struct {
	orderUC orderUsecase.OrderUsecase
}

func NewPaymentUsecase(orderUC orderUsecase.OrderUsecase) PaymentUsecase {
	return &paymentUsecase{
		orderUC: orderUC,
	}
}

func (u *paymentUsecase) HandleWebhookMidtrans(payload map[string]interface{}) error {
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return errors.New("order_id tidak ditemukan di payload")
	}

	transactionStatus, ok := payload["transaction_status"].(string)
	if !ok {
		return errors.New("transaction_status tidak ditemukan di payload")
	}

	// Handle tombol "Test" dari Dashboard Midtrans (mengirimkan dummy order_id)
	if len(orderID) > 19 && orderID[:19] == "payment_notif_test_" {
		return nil // Langsung return success agar dashboard Midtrans centang hijau
	}

	// Ambil data order KTL
	order, err := u.orderUC.GetOrderKTLByID(orderID)
	if err != nil {
		return errors.New("order tidak ditemukan")
	}

	// Logika perubahan status berdasarkan notifikasi Midtrans
	switch transactionStatus {
	case "settlement", "capture":
		order.Status = constants.OrderKTLStatusPending // Menjadi Pending (siap diproses)
	case "cancel", "deny", "expire":
		order.Status = "Batal"
	case "pending":
		order.Status = "Menunggu Pembayaran"
	}

	req := orderDto.UpdateOrderStatusRequest{
		Status: order.Status,
	}
	// Menggunakan "" agar adminID menjadi nil di database
	return u.orderUC.UpdateOrderKTLStatus(orderID, req, "")
}
