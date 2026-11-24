package routes

import (
	"bj/internal/domain/courier/handler"
	"bj/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterCourierRoutes(router fiber.Router, h *handler.CourierHandler) {
	admin := router.Group("/admin/couriers", middleware.Protected())

	// GET  /api/v1/admin/couriers      → Ambil seluruh master kurir
	// POST /api/v1/admin/couriers/sync → Sinkronisasi kurir dari Biteship
	admin.Get("/", h.GetAllCouriers)
	admin.Post("/sync", h.SyncCouriers)
}

