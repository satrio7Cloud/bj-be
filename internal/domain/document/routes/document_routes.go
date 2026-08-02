package routes

import (
	"bj/internal/domain/document/handler"
	"bj/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterDocumentRoutes(router fiber.Router, h *handler.DocumentHandler) {
	// Rute Publik (Customer Web)
	router.Post("/documents/customer", h.CreateDocumentCustomer)

	// Semua rute Document Admin harus diproteksi dengan JWT
	group := router.Group("/documents", middleware.Protected())
	group.Post("/admin", h.CreateDocumentAdmin)
	group.Get("/order/:orderId", h.GetDocumentsByOrder)
	group.Get("/:id/url", h.GetDocumentURL)
	group.Patch("/:id/verify", h.VerifyDocument)
	group.Delete("/:id", h.DeleteDocument)
}
