package routes

import (
	"bj/internal/domain/deliverymethod/handler"
	"bj/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterDeliveryMethodRoutes(router fiber.Router, h *handler.DeliveryMethodHandler) {

	router.Get("/delivery-methods", h.GetActiveMethods)
	admin := router.Group("/admin/delivery-methods", middleware.Protected())

	admin.Get("/", h.GetAllMethods)
	admin.Get("/:id", h.GetMethodByID)
	admin.Post("/", h.CreateMethod)
	admin.Put("/:id", h.UpdateMethod)
	admin.Patch("/:id/toggle", h.ToggleMethod)
	admin.Delete("/:id", h.DeleteMethod)
}
