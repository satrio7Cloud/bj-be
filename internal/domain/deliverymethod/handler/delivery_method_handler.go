package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/deliverymethod/dto"
	"bj/internal/domain/deliverymethod/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type DeliveryMethodHandler struct {
	uc usecase.DeliveryMethodUsecase
}

func NewDeliveryMethodHandler(uc usecase.DeliveryMethodUsecase) *DeliveryMethodHandler {
	return &DeliveryMethodHandler{uc: uc}
}

// GetActiveMethods godoc
// @Summary Ambil metode pengantaran yang aktif
// @Description Mengambil opsi metode penyerahan/pengembalian aktif untuk checkout (bisa difilter direction=PICKUP, RETURN, atau BOTH)
// @Tags delivery-methods
// @Produce json
// @Param direction query string false "Filter Direction: PICKUP, RETURN, atau BOTH"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/delivery-methods [get]
func (h *DeliveryMethodHandler) GetActiveMethods(c *fiber.Ctx) error {
	direction := c.Query("direction")
	data, err := h.uc.GetActive(direction)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   data,
	})
}

// GetAllMethods godoc
// @Summary Ambil seluruh master metode pengantaran (Admin)
// @Description Mengambil semua metode pengantaran termasuk yang nonaktif
// @Tags delivery-methods
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/delivery-methods [get]
func (h *DeliveryMethodHandler) GetAllMethods(c *fiber.Ctx) error {
	data, err := h.uc.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   data,
	})
}

// GetMethodByID godoc
// @Summary Ambil detail metode pengantaran berdasarkan ID
// @Tags delivery-methods
// @Produce json
// @Param id path string true "Delivery Method ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/admin/delivery-methods/{id} [get]
func (h *DeliveryMethodHandler) GetMethodByID(c *fiber.Ctx) error {
	id := c.Params("id")
	data, err := h.uc.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   data,
	})
}

// CreateMethod godoc
// @Summary Tambah metode pengantaran baru
// @Tags delivery-methods
// @Accept json
// @Produce json
// @Param request body dto.CreateDeliveryMethodRequest true "Data Metode Pengantaran"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/delivery-methods [post]
func (h *DeliveryMethodHandler) CreateMethod(c *fiber.Ctx) error {
	var req dto.CreateDeliveryMethodRequest
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

	created, err := h.uc.Create(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Metode pengantaran berhasil ditambahkan",
		"data":    created,
	})
}

// UpdateMethod godoc
// @Summary Ubah metode pengantaran
// @Tags delivery-methods
// @Accept json
// @Produce json
// @Param id path string true "Delivery Method ID"
// @Param request body dto.UpdateDeliveryMethodRequest true "Data Perubahan"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/delivery-methods/{id} [put]
func (h *DeliveryMethodHandler) UpdateMethod(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateDeliveryMethodRequest
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

	if err := h.uc.Update(id, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Metode pengantaran berhasil diperbarui",
	})
}

// ToggleMethod godoc
// @Summary Aktifkan / Nonaktifkan metode pengantaran
// @Tags delivery-methods
// @Produce json
// @Param id path string true "Delivery Method ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/admin/delivery-methods/{id}/toggle [patch]
func (h *DeliveryMethodHandler) ToggleMethod(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.uc.ToggleActive(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Status metode pengantaran berhasil diubah",
	})
}

// DeleteMethod godoc
// @Summary Hapus metode pengantaran
// @Tags delivery-methods
// @Produce json
// @Param id path string true "Delivery Method ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/admin/delivery-methods/{id} [delete]
func (h *DeliveryMethodHandler) DeleteMethod(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.uc.Delete(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Metode pengantaran berhasil dihapus",
	})
}

