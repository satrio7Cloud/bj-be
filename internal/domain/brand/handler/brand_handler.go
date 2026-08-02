package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/brand/dto"
	"bj/internal/domain/brand/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type BrandHandler struct {
	brandUsecase usecase.BrandUsecase
}

func NewBrandHandler(uc usecase.BrandUsecase) *BrandHandler {
	return &BrandHandler{
		brandUsecase: uc,
	}
}

// @Summary Create New Brand
// @Description Create a new vehicle brand
// @Tags brands
// @Accept json
// @Produce json
// @Param request body dto.CreateBrandRequest true "Brand Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/brands [post]
func (h *BrandHandler) CreateBrand(c *fiber.Ctx) error {
	var req dto.CreateBrandRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest})
	}
	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest, "errors": validationErrors})
	}

	if err := h.brandUsecase.CreateBrand(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": constants.StatusSuccess, "message": "Merek berhasil ditambahkan"})
}

// @Summary Get All Brands
// @Description Retrieve a list of all vehicle brands
// @Tags brands
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/brands [get]
func (h *BrandHandler) GetBrands(c *fiber.Ctx) error {
	brands, err := h.brandUsecase.GetAllBrands()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": constants.StatusSuccess, "data": brands})
}
