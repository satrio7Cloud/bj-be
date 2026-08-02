package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/vehicle/dto"
	"bj/internal/domain/vehicle/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type VehicleHandler struct {
	vehicleUsecase usecase.VehicleUsecase
}

func NewVehicleHandler(uc usecase.VehicleUsecase) *VehicleHandler {
	return &VehicleHandler{
		vehicleUsecase: uc,
	}
}

// @Summary Create New Vehicle
// @Description Create a new customer vehicle
// @Tags vehicles
// @Accept json
// @Produce json
// @Param request body dto.CreateVehicleRequest true "Vehicle Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/vehicles [post]
func (h *VehicleHandler) CreateVehicle(c *fiber.Ctx) error {
	var req dto.CreateVehicleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest})
	}
	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest, "errors": validationErrors})
	}

	newID, err := h.vehicleUsecase.CreateVehicle(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": constants.StatusSuccess, 
		"message": "Kendaraan berhasil ditambahkan",
		"data": fiber.Map{
			"id": newID,
		},
	})
}

// @Summary Get All Vehicles
// @Description Retrieve a list of all customer vehicles
// @Tags vehicles
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/vehicles [get]
func (h *VehicleHandler) GetVehicles(c *fiber.Ctx) error {
	vehicles, err := h.vehicleUsecase.GetAllVehicles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": constants.StatusSuccess, "data": vehicles})
}
