package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/order/dto"
	"bj/internal/domain/order/usecase"
	"bj/internal/utils"

	"context"
	"strconv"
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

// @Summary Create Order KTL
// @Description Create a new KTL print request with file uploads in one step
// @Tags orders
// @Accept multipart/form-data
// @Produce json
// @Param customer_name formData string true "Customer Name"
// @Param phone_number formData string true "Phone Number"
// @Param vehicle_type formData string true "Vehicle Type"
// @Param region formData string true "Region"
// @Param fee formData number true "Fee"
// @Param merchant_ref formData string false "Merchant Ref"
// @Param bukti_ktl formData file false "Foto KTL Lama"
// @Param bukti_pembayaran formData file false "Bukti Pembayaran"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/ktl [post]
func (h *OrderHandler) CreateOrderKTL(c *fiber.Ctx) error {
	var req dto.CreateOrderKTLRequest

	// Manual parsing for multipart/form-data
	req.CustomerName = c.FormValue("customer_name")
	req.PhoneNumber = c.FormValue("phone_number")
	req.VehicleType = c.FormValue("vehicle_type")
	req.Region = c.FormValue("region")
	req.MerchantRef = c.FormValue("merchant_ref")
	
	feeStr := c.FormValue("fee")
	if feeStr != "" {
		feeFloat, err := strconv.ParseFloat(feeStr, 64)
		if err == nil {
			req.Fee = feeFloat
		}
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
			"errors":  validationErrors,
		})
	}

	ctx := context.Background()
	fileKtl, _ := c.FormFile("bukti_ktl")
	filePembayaran, _ := c.FormFile("bukti_pembayaran")

	resp, err := h.orderUsecase.CreateOrderKTL(ctx, req, fileKtl, filePembayaran)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Permintaan cetak KTL berhasil dibuat",
		"data":    resp,
	})
}

// @Summary Get All Order KTLs
// @Description Get all KTL orders with pagination
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/ktl [get]
func (h *OrderHandler) GetAllOrderKTLs(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	orders, err := h.orderUsecase.GetAllOrderKTLs(page, limit)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   orders,
	})
}

// @Summary Get Order KTL by ID
// @Description Get specific KTL order details
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order KTL ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/orders/ktl/{id} [get]
func (h *OrderHandler) GetOrderKTLByID(c *fiber.Ctx) error {
	id := c.Params("id")
	order, err := h.orderUsecase.GetOrderKTLByID(id)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   order,
	})
}

// @Summary Update Order KTL Status
// @Description Update the processing status of a KTL order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order KTL ID"
// @Param request body dto.UpdateOrderStatusRequest true "Status Update Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/orders/ktl/{id}/status [patch]
func (h *OrderHandler) UpdateOrderKTLStatus(c *fiber.Ctx) error {
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

	adminID := c.Locals("admin_id").(string)

	err := h.orderUsecase.UpdateOrderKTLStatus(id, req, adminID)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Status KTL berhasil diperbarui",
	})
}
