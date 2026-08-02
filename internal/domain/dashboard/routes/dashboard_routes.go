package routes

import (
	"bj/internal/domain/dashboard/handler"
	"bj/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterDashboardRoutes(router fiber.Router, h *handler.DashboardHandler) {
	group := router.Group("/dashboard", middleware.Protected())
	group.Get("/summary", h.GetSummary)
	group.Get("/orders", h.GetOrders)
	group.Get("/chart", h.GetChartData)
	group.Get("/activities", h.GetRecentActivities)
}
