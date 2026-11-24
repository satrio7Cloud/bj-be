package routes

import (
	"bj/internal/domain/biteship/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterBiteshipRoutes(router fiber.Router, hdl *handler.BiteshipHandler) {
	shippingGroup := router.Group("/shipping")

	// Endpoint operasional kurir Biteship (buat order & tracking resi)
	shippingGroup.Post("/order", hdl.CreateOrder)
	shippingGroup.Get("/tracking/:tracking_id", hdl.GetTracking)
}
