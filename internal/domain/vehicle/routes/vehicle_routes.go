package routes

import (
	"bj/internal/domain/vehicle/handler"
	"bj/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterVehicleRoutes(router fiber.Router, h *handler.VehicleHandler) {
	group := router.Group("/vehicles", middleware.Protected())
	group.Post("/", h.CreateVehicle)
	group.Get("/", h.GetVehicles)
}
