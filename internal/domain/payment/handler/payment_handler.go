package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/payment/usecase"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	paymentUsecase usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase: paymentUsecase,
	}
}

// MidtransWebhook handles notifications from Midtrans
func (h *PaymentHandler) MidtransWebhook(c *fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": "Payload webhook tidak valid",
		})
	}

	err := h.paymentUsecase.HandleWebhookMidtrans(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Webhook berhasil diproses",
	})
}
