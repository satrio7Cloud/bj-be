package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/company/dto"
	"bj/internal/domain/company/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type CompanyProfileHandler struct {
	usecase usecase.CompanyProfileUsecase
}

func NewCompanyProfileHandler(usecase usecase.CompanyProfileUsecase) *CompanyProfileHandler {
	return &CompanyProfileHandler{usecase: usecase}
}

// @Summary Get Company / Office Profile
// @Description Retrieve current office address, postal code, phone number, and operational hours
// @Tags company-profile
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/company-profile [get]
func (h *CompanyProfileHandler) GetProfile(c *fiber.Ctx) error {
	profile, err := h.usecase.GetProfile()
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Profil kantor berhasil diambil",
		"data":    profile,
	})
}

// @Summary Create Company / Office Profile
// @Description Create office address, postal code, phone number, and operational hours (Admin Only)
// @Tags company-profile
// @Accept json
// @Produce json
// @Param request body dto.CreateCompanyProfileRequest true "Office Profile Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/company-profile [post]
func (h *CompanyProfileHandler) CreateProfile(c *fiber.Ctx) error {
	var req dto.CreateCompanyProfileRequest

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

	profile, err := h.usecase.CreateProfile(req)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Profil kantor berhasil dibuat",
		"data":    profile,
	})
}

// @Summary Update Company / Office Profile
// @Description Update office address, postal code, phone number, and operational hours (Admin Only)
// @Tags company-profile
// @Accept json
// @Produce json
// @Param request body dto.UpdateCompanyProfileRequest true "Office Profile Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/company-profile [put]
func (h *CompanyProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	var req dto.UpdateCompanyProfileRequest

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

	profile, err := h.usecase.UpdateProfile(req)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Profil kantor berhasil diperbarui",
		"data":    profile,
	})
}
