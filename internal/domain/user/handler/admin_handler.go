package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/user/dto"
	"bj/internal/domain/user/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	adminUsecase usecase.AdminUsecase
}

func NewAdminHandler(usecase usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{usecase}
}

// Login godoc
// @Summary Admin Login
// @Description Authenticate admin and get JWT token
// @Tags admins
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/admins/login [post]
func (h *AdminHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	// Parsing JSON Body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
		})
	}

	// Validasi Input Otomatis menggunakan Struct Validator
	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
			"errors":  validationErrors,
		})
	}

	// Panggil Usecase Login
	resp, err := h.adminUsecase.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	// Response Sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": constants.MsgLoginSuccess,
		"data":    resp,
	})
}

// Register godoc
// @Summary Register New Admin
// @Description Create a new admin account
// @Tags admins
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admins/register [post]
func (h *AdminHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	// Parsing JSON Body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
		})
	}

	// Validasi Input Otomatis menggunakan Struct Validator
	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": constants.MsgInvalidRequest,
			"errors":  validationErrors,
		})
	}

	// Panggil Usecase
	err := h.adminUsecase.Register(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  constants.StatusError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Admin berhasil didaftarkan",
	})
}

// Logout godoc
// @Summary Admin Logout
// @Description Logout admin (invalidate token on client side)
// @Tags admins
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admins/logout [post]
func (h *AdminHandler) Logout(c *fiber.Ctx) error {
	// Karena autentikasi menggunakan JWT secara stateless (tanpa Redis blacklist),
	// proses "Logout" yang sesungguhnya terjadi di sisi Frontend dengan menghapus Token dari LocalStorage.
	// Endpoint ini hanya memberikan respons sukses (200 OK) untuk mengkonfirmasi request dari klien.

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Berhasil logout",
	})
}
