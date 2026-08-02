package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/order/dto"
	"bj/internal/domain/order/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(uc usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: uc,
	}
}

// @Summary Create New Order (Admin)
// @Description Create a new order with full control over fees from the admin dashboard
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderAdminRequest true "Order Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/admin [post]
func (h *OrderHandler) CreateOrderAdmin(c *fiber.Ctx) error {
	var req dto.CreateOrderAdminRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
		})
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
			"errors":  validationErrors,
		})
	}

	// Mengambil ID Admin dari JWT Context
	adminID := c.Locals("admin_id").(string)

	newID, err := h.orderUsecase.CreateOrderAdmin(req, adminID)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Pesanan berhasil dibuat oleh Admin",
		"data": fiber.Map{
			"id": newID,
		},
	})
}

// @Summary Create New Order (Customer)
// @Description Create a new order from public web, fees are auto-calculated
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderCustomerRequest true "Order Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/customer [post]
func (h *OrderHandler) CreateOrderCustomer(c *fiber.Ctx) error {
	var req dto.CreateOrderCustomerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
		})
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
			"errors":  validationErrors,
		})
	}

	newID, err := h.orderUsecase.CreateOrderCustomer(req)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Pesanan perpanjangan STNK berhasil dibuat",
		"data": fiber.Map{
			"id": newID,
		},
	})
}

// @Summary Get All Orders
// @Description Retrieve a paginated list of all orders
// @Tags orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(50)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders [get]
func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)

	orders, err := h.orderUsecase.GetAllOrders(page, limit)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   orders,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
		},
	})
}

// @Summary Get Order by ID
// @Description Retrieve detailed information of a specific order
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	order, err := h.orderUsecase.GetOrderByID(id)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   order,
	})
}

// @Summary Update Order Status
// @Description Update the processing status of an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body dto.UpdateOrderStatusRequest true "Status Update Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/{id}/status [patch]
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateOrderStatusRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
		})
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
			"errors":  validationErrors,
		})
	}

	// Mengambil ID Admin dari JWT Context
	adminID := c.Locals("admin_id").(string)

	err := h.orderUsecase.UpdateOrderStatus(id, req, adminID)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Status pesanan berhasil diperbarui",
	})
}
