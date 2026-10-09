package routes

import (
	"bj/internal/domain/payment/handler"
	"github.com/gofiber/fiber/v2"
)

func RegisterPaymentRoutes(router fiber.Router, paymentHandler *handler.PaymentHandler) {
	paymentGroup := router.Group("/payment")

	// Rute Webhook Midtrans
	paymentGroup.Post("/webhook", paymentHandler.MidtransWebhook)
}
