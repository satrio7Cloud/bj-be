package routes

import (
	"bj/internal/domain/user/handler"
	"bj/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterAdminRoutes(router fiber.Router, adminHandler *handler.AdminHandler) {
	// Grup endpoint untuk /admins
	adminGroup := router.Group("/admins")

	// Mendaftarkan endpoint login
	// POST /api/v1/admins/login
	adminGroup.Post("/login", adminHandler.Login)
	// Mendaftarkan endpoint register
	// POST /api/v1/admins/register
	adminGroup.Post("/register", adminHandler.Register)

	// -- CONTOH PENGGUNAAN MIDDLEWARE --
	// Semua route di bawah group "protected" wajib melampirkan Token JWT di Header
	protected := adminGroup.Group("/", middleware.Protected())

	// GET /api/v1/admins/me
	protected.Get("/me", func(c *fiber.Ctx) error {
		// Mengambil data yang disimpan oleh middleware di context
		adminID := c.Locals("admin_id")
		adminEmail := c.Locals("admin_email")
		adminRole := c.Locals("admin_role")

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"id":    adminID,
				"email": adminEmail,
				"role":  adminRole,
			},
		})
	})
	
	// POST /api/v1/admins/logout
	protected.Post("/logout", adminHandler.Logout)
}
