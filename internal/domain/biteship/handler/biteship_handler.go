package handler

import (
	"bj/internal/constants"
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

// CheckRates godoc
// @Summary Check shipping rates
// @Description Check shipping rates for available couriers based on origin and destination
// @Tags shipping
// @Accept json
// @Produce json
// @Param request body dto.RatesRequest true "Rates Request Payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/shipping/rates [post]
func (h *BiteshipHandler) CheckRates(c *fiber.Ctx) error {
	var req dto.RatesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
		})
	}

	if req.DestinationPostalCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgDestinationPostalCodeRequired,
		})
	}

	resp, err := h.usecase.CheckRates(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": constants.MsgSuccessFetchRates,
		"data":    resp.Pricing,
	})
}

// GetCouriers godoc
// @Summary Get available couriers
// @Description Get a list of available couriers and their services
// @Tags shipping
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/shipping/couriers [get]
func (h *BiteshipHandler) GetCouriers(c *fiber.Ctx) error {
	resp, err := h.usecase.GetAvailableCouriers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": constants.MsgSuccessFetchCouriers,
		"data":    resp.Couriers,
	})
}

// CreateOrder godoc
// @Summary Create a new shipping order
// @Description Create a new order via Biteship and get a tracking ID
// @Tags shipping
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderRequest true "Create Order Request Payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/shipping/order [post]
func (h *BiteshipHandler) CreateOrder(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
		})
	}

	resp, err := h.usecase.CreateOrder(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": constants.MsgSuccessCreateOrder,
		"data":    resp,
	})
}

// GetTracking godoc
// @Summary Get shipping tracking
// @Description Get tracking history and status for a specific tracking ID
// @Tags shipping
// @Produce json
// @Param tracking_id path string true "Tracking ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/shipping/tracking/{tracking_id} [get]
func (h *BiteshipHandler) GetTracking(c *fiber.Ctx) error {
	trackingID := c.Params("tracking_id")
	if trackingID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgTrackingIDRequired,
		})
	}

	resp, err := h.usecase.GetTracking(trackingID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": constants.MsgSuccessFetchTracking,
		"data":    resp,
	})
}
