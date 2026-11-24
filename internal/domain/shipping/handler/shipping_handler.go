package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/shipping/dto"
	"bj/internal/domain/shipping/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type ShippingHandler struct {
	shippingUsecase usecase.ShippingUsecase
}

func NewShippingHandler(uc usecase.ShippingUsecase) *ShippingHandler {
	return &ShippingHandler{shippingUsecase: uc}
}

// GetShippingSettings godoc
// @Summary Ambil konfigurasi pengaturan pengiriman
// @Description Membaca konfigurasi aktif: kurir aktif, berat standar berkas & nilai asuransi
// @Tags shipping
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/settings/shipping [get]
func (h *ShippingHandler) GetShippingSettings(c *fiber.Ctx) error {
	settings, err := h.shippingUsecase.GetSettings()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   settings,
	})
}

// SaveShippingSettings godoc
// @Summary Simpan konfigurasi pengaturan pengiriman
// @Description Admin menyimpan kurir aktif, berat standar berkas (gr) & nilai asuransi (Rp)
// @Tags shipping
// @Accept json
// @Produce json
// @Param request body dto.UpdateShippingSettingsRequest true "Konfigurasi Pengiriman"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/settings/shipping [put]
func (h *ShippingHandler) SaveShippingSettings(c *fiber.Ctx) error {
	var req dto.UpdateShippingSettingsRequest
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

	adminID, _ := c.Locals("admin_id").(string)

	if err := h.shippingUsecase.SaveSettings(adminID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Pengaturan pengiriman berhasil disimpan",
	})
}

// CalculateRates godoc
// @Summary Hitung ongkir
// @Description Menghitung tarif pengiriman berdasarkan kurir aktif dan metode pickup
// @Tags shipping
// @Accept json
// @Produce json
// @Param request body dto.CalculateRatesRequest true "Data Perhitungan Ongkir"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/shipping/rates [post]
func (h *ShippingHandler) CalculateRates(c *fiber.Ctx) error {
	var req dto.CalculateRatesRequest
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

	result, err := h.shippingUsecase.CalculateRates(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   result,
	})
}
