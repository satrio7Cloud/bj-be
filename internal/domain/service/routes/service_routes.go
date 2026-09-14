package routes

import (
	"bj/internal/domain/service/handler"
	"bj/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterServiceRoutes(router fiber.Router, h *handler.ServiceHandler) {
	group := router.Group("/services")
	group.Post("/", middleware.Protected(), h.CreateService)
	group.Get("/", h.GetServices)
	group.Put("/:id", middleware.Protected(), h.UpdateService)
}
