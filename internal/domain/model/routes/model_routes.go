package routes

import (
	"bj/internal/domain/model/handler"
	"bj/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterModelRoutes(router fiber.Router, h *handler.ModelHandler) {
	group := router.Group("/models", middleware.Protected())
	group.Post("/", h.CreateModel)
	group.Get("/", h.GetModels)
}
