package routes

import (
	"bj/internal/domain/biteship/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterBiteshipRoutes(router fiber.Router, hdl *handler.BiteshipHandler) {
	shippingGroup := router.Group("/shipping")

	shippingGroup.Post("/rates", hdl.CheckRates)
	shippingGroup.Get("/couriers", hdl.GetCouriers) // Endpoint baru untuk melihat list kurir
	
	shippingGroup.Post("/order", hdl.CreateOrder)
	shippingGroup.Get("/tracking/:tracking_id", hdl.GetTracking)
}
