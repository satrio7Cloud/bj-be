package routes

import (
	"bj/internal/domain/order/handler"
	"bj/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterOrderRoutes(router fiber.Router, orderHandler *handler.OrderHandler) {
	// Rute Publik (Customer Web)
	router.Post("/orders/customer", orderHandler.CreateOrderCustomer)

	// Semua rute Order Admin harus diproteksi dengan JWT
	orderGroup := router.Group("/orders", middleware.Protected())

	orderGroup.Post("/admin", orderHandler.CreateOrderAdmin)
	orderGroup.Get("/", orderHandler.GetOrders)
	orderGroup.Get("/:id", orderHandler.GetOrderByID)
	orderGroup.Patch("/:id/status", orderHandler.UpdateOrderStatus)
}
