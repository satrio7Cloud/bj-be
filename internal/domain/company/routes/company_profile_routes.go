package routes

import (
	"bj/internal/domain/company/handler"
	"bj/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterCompanyProfileRoutes(router fiber.Router, h *handler.CompanyProfileHandler) {
	group := router.Group("/company-profile")
	group.Get("/", h.GetProfile)
	group.Post("/", middleware.Protected(), h.CreateProfile)
	group.Put("/", middleware.Protected(), h.UpdateProfile)
}
