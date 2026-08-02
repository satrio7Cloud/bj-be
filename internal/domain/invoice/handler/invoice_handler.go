package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/invoice/dto"
	"bj/internal/domain/invoice/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type InvoiceHandler struct {
	invoiceUsecase usecase.InvoiceUsecase
}

func NewInvoiceHandler(uc usecase.InvoiceUsecase) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceUsecase: uc,
	}
}

// @Summary Midtrans Webhook
// @Description Callback URL for Midtrans to notify payment status
// @Tags invoices
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/invoices/midtrans-callback [post]
func (h *InvoiceHandler) MidtransCallback(c *fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid payload"})
	}

	if err := h.invoiceUsecase.HandleMidtransNotification(payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "ok"})
}

// @Summary Generate Invoice for an Order
// @Description Generates a new invoice or returns the existing one for a given order
// @Tags invoices
// @Accept json
// @Produce json
// @Param request body dto.GenerateInvoiceRequest true "Generate Invoice Request"
// @Success 201 {object} dto.InvoiceDetailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/invoices [post]
func (h *InvoiceHandler) GenerateInvoice(c *fiber.Ctx) error {
	var req dto.GenerateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest})
	}
	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest, "errors": validationErrors})
	}

	resp, err := h.invoiceUsecase.GenerateInvoice(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   resp,
	})
}

// @Summary Get Invoice By Order ID
// @Description Retrieve a detailed invoice JSON by Order ID
// @Tags invoices
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {object} dto.InvoiceDetailResponse
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/invoices/order/{orderId} [get]
func (h *InvoiceHandler) GetInvoiceByOrderID(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	resp, err := h.invoiceUsecase.GetInvoiceByOrderID(orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   resp,
	})
}

// @Summary Pay Invoice
// @Description Simulate payment of an invoice (Update status to PAID)
// @Tags invoices
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param request body dto.PayInvoiceRequest true "Payment Info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/invoices/{id}/pay [patch]
func (h *InvoiceHandler) PayInvoice(c *fiber.Ctx) error {
	invoiceID := c.Params("id")

	var req dto.PayInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest})
	}
	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest, "errors": validationErrors})
	}

	if err := h.invoiceUsecase.PayInvoice(invoiceID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Invoice berhasil dibayar",
	})
}

// @Summary Get All Invoices
// @Description Retrieve a paginated list of all invoices
// @Tags invoices
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/invoices [get]
func (h *InvoiceHandler) GetAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	resp, err := h.invoiceUsecase.GetAll(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   resp["data"],
		"meta": fiber.Map{
			"total": resp["total"],
			"page":  resp["page"],
			"limit": resp["limit"],
		},
	})
}

// @Summary Get All Invoices
// @Description Retrieve a paginated list of all invoices
// @Tags invoices
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/invoices/list [get]
func (h *InvoiceHandler) GetInvoices(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	invoices, total, err := h.invoiceUsecase.GetAllInvoices(page, limit)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   invoices,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}


