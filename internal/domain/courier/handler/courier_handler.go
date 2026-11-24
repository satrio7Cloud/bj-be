package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/courier/usecase"

	"github.com/gofiber/fiber/v2"
)

type CourierHandler struct {
	uc usecase.CourierUsecase
}

func NewCourierHandler(uc usecase.CourierUsecase) *CourierHandler {
	return &CourierHandler{uc: uc}
}

// GetAllCouriers godoc
// @Summary Ambil seluruh master kurir & layanan
// @Description Mengambil daftar lengkap kurir hasil sinkronisasi dari database lokal (admin only)
// @Tags couriers
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/couriers [get]
func (h *CourierHandler) GetAllCouriers(c *fiber.Ctx) error {
	couriers, err := h.uc.GetAllCouriers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   couriers,
	})
}

// SyncCouriers godoc
// @Summary Sinkronisasi kurir dari Biteship
// @Description Menarik & memperbarui master kurir terbaru dari API Biteship ke database lokal
// @Tags couriers
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/couriers/sync [post]
func (h *CourierHandler) SyncCouriers(c *fiber.Ctx) error {
	result, err := h.uc.SyncCouriers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": result.Message,
		"data": fiber.Map{
			"synced": result.Synced,
		},
	})
}

