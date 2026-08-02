package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/model/dto"
	"bj/internal/domain/model/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type ModelHandler struct {
	modelUsecase usecase.ModelUsecase
}

func NewModelHandler(uc usecase.ModelUsecase) *ModelHandler {
	return &ModelHandler{
		modelUsecase: uc,
	}
}

// @Summary Create New Model
// @Description Create a new vehicle model associated with a brand
// @Tags models
// @Accept json
// @Produce json
// @Param request body dto.CreateModelRequest true "Model Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/models [post]
func (h *ModelHandler) CreateModel(c *fiber.Ctx) error {
	var req dto.CreateModelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest})
	}
	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest, "errors": validationErrors})
	}

	if err := h.modelUsecase.CreateModel(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": constants.StatusSuccess, "message": "Model berhasil ditambahkan"})
}

// @Summary Get All Models
// @Description Retrieve a list of all vehicle models
// @Tags models
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/models [get]
func (h *ModelHandler) GetModels(c *fiber.Ctx) error {
	models, err := h.modelUsecase.GetAllModels()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": constants.StatusSuccess, "data": models})
}
