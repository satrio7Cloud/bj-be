package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/service/dto"
	"bj/internal/domain/service/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type ServiceHandler struct {
	serviceUsecase usecase.ServiceUsecase
}

func NewServiceHandler(uc usecase.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{
		serviceUsecase: uc,
	}
}

// @Summary Create New Service
// @Description Create a new master service catalog item
// @Tags services
// @Accept json
// @Produce json
// @Param request body dto.CreateServiceRequest true "Service Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/services [post]
func (h *ServiceHandler) CreateService(c *fiber.Ctx) error {
	var req dto.CreateServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest})
	}
	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest, "errors": validationErrors})
	}

	if err := h.serviceUsecase.CreateService(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": constants.StatusSuccess, "message": "Layanan berhasil ditambahkan"})
}

// @Summary Get All Services
// @Description Retrieve a list of all master services
// @Tags services
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/services [get]
func (h *ServiceHandler) GetServices(c *fiber.Ctx) error {
	services, err := h.serviceUsecase.GetAllServices()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": constants.StatusSuccess, "data": services})
}
