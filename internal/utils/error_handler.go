package utils

import (
	"strings"

	"bj/internal/constants"

	"github.com/gofiber/fiber/v2"
)

// HandleError parses the error string and returns an appropriate HTTP status code and formatted JSON response.
func HandleError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	// 1. Foreign Key Constraints (Data tidak ditemukan)
	if strings.Contains(errMsg, "FOREIGN KEY constraint") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": "Data referensi tidak valid atau tidak ditemukan (misal: Customer, Vehicle, atau Service ID salah).",
		})
	}

	// 2. Unique Constraints (Data duplikat)
	if strings.Contains(errMsg, "UNIQUE KEY constraint") || strings.Contains(errMsg, "Violation of UNIQUE KEY") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": "Data yang dimasukkan sudah ada atau duplikat.",
		})
	}

	// 3. Not Found (Data tidak ada)
	if strings.Contains(strings.ToLower(errMsg), "not found") || strings.Contains(strings.ToLower(errMsg), "tidak ditemukan") {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": errMsg,
		})
	}

	// 4. Data Type / Cast Errors (Format UUID salah)
	if strings.Contains(errMsg, "Conversion failed") || strings.Contains(errMsg, "uniqueidentifier") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": "Format ID tidak valid. Pastikan menggunakan format UUID yang benar.",
		})
	}

	// 5. Default Internal Server Error
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  constants.StatusError,
		"message": constants.MsgInternalError,
		"debug":   errMsg, // Opsional: hapus jika ingin menyembunyikan detail error dari user di production
	})
}
