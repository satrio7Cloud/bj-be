package routes

import (
	"bj/internal/domain/customer/handler"
	"bj/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterCustomerRoutes(router fiber.Router, h *handler.CustomerHandler) {
	group := router.Group("/customers", middleware.Protected())
	group.Post("/", h.CreateCustomer)
	group.Get("/", h.GetCustomers)
}
