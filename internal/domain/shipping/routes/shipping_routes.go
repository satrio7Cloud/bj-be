package routes

import (
	"bj/internal/domain/shipping/handler"
	"bj/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterShippingRoutes(router fiber.Router, h *handler.ShippingHandler) {
	// ── Public Routes (Tanpa Auth) ───────────────────────────────────────────
	// GET  /api/v1/settings/shipping → Baca konfigurasi aktif (kurir aktif, berat, nilai)
	router.Get("/settings/shipping", h.GetShippingSettings)
	// POST /api/v1/shipping/rates    → Hitung ongkir (difilter kurir aktif)
	router.Post("/shipping/rates", h.CalculateRates)

	// ── Admin Routes (Protected JWT) ─────────────────────────────────────────
	admin := router.Group("/admin/settings", middleware.Protected())
	// PUT /api/v1/admin/settings/shipping → Simpan konfigurasi kurir aktif, berat & nilai berkas
	admin.Put("/shipping", h.SaveShippingSettings)
}
