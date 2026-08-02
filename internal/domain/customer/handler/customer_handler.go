package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/customer/dto"
	"bj/internal/domain/customer/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	customerUsecase usecase.CustomerUsecase
}

func NewCustomerHandler(uc usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		customerUsecase: uc,
	}
}

// @Summary Create New Customer
// @Description Create a new customer profile
// @Tags customers
// @Accept json
// @Produce json
// @Param request body dto.CreateCustomerRequest true "Customer Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/customers [post]
func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var req dto.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest})
	}
	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest, "errors": validationErrors})
	}

	newID, err := h.customerUsecase.CreateCustomer(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": constants.StatusSuccess, 
		"message": "Pelanggan berhasil ditambahkan",
		"data": fiber.Map{
			"id": newID,
		},
	})
}

// @Summary Get All Customers
// @Description Retrieve a list of all customers
// @Tags customers
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/customers [get]
func (h *CustomerHandler) GetCustomers(c *fiber.Ctx) error {
	customers, err := h.customerUsecase.GetAllCustomers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": constants.StatusSuccess, "data": customers})
}
