package routes

import (
	"bj/internal/domain/order/handler"
	"bj/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterOrderRoutes(router fiber.Router, orderHandler *handler.OrderHandler) {
	// Route Customer
	router.Post("/orders/customer", orderHandler.CreateOrderCustomer)
	router.Post("/orders/ktl", orderHandler.CreateOrderKTL)

	// Route Admin
	orderGroup := router.Group("/orders", middleware.Protected())
	orderGroup.Post("/admin", orderHandler.CreateOrderAdmin)
	orderGroup.Get("/", orderHandler.GetOrders)
	orderGroup.Get("/ktl", orderHandler.GetAllOrderKTLs)
	orderGroup.Patch("/ktl/:id/status", orderHandler.UpdateOrderKTLStatus)
	orderGroup.Get("/ktl/:id", orderHandler.GetOrderKTLByID)
	orderGroup.Get("/:id", orderHandler.GetOrderByID)
	orderGroup.Patch("/:id/status", orderHandler.UpdateOrderStatus)
}
