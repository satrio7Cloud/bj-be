package routes

import (
	"bj/internal/domain/brand/handler"
	"bj/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterBrandRoutes(router fiber.Router, h *handler.BrandHandler) {
	group := router.Group("/brands")
	group.Post("/", middleware.Protected(), h.CreateBrand)
	group.Get("/", h.GetBrands)
}
