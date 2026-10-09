package handler

import (
	"bj/internal/domain/biteship/dto"
	"bj/internal/domain/biteship/usecase"

	"github.com/gofiber/fiber/v2"
)

type BiteshipHandler struct {
	usecase usecase.BiteshipUsecase
}

func NewBiteshipHandler(usecase usecase.BiteshipUsecase) *BiteshipHandler {
	return &BiteshipHandler{usecase: usecase}
}

func (h *BiteshipHandler) CheckRates(c *fiber.Ctx) error {
	var req dto.RatesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if req.DestinationPostalCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "destination_postal_code is required",
		})
	}

	resp, err := h.usecase.CheckRates(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Success fetch rates",
		"data":    resp.Pricing,
	})
}

func (h *BiteshipHandler) GetCouriers(c *fiber.Ctx) error {
	resp, err := h.usecase.GetAvailableCouriers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Success fetch couriers",
		"data":    resp.Couriers,
	})
}

func (h *BiteshipHandler) CreateOrder(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	resp, err := h.usecase.CreateOrder(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Success create order",
		"data":    resp,
	})
}

func (h *BiteshipHandler) GetTracking(c *fiber.Ctx) error {
	trackingID := c.Params("tracking_id")
	if trackingID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "tracking_id is required",
		})
	}

	resp, err := h.usecase.GetTracking(trackingID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Success fetch tracking",
		"data":    resp,
	})
}
