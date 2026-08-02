package routes

import (
	"bj/internal/domain/invoice/handler"

	"bj/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterInvoiceRoutes(router fiber.Router, invoiceHandler *handler.InvoiceHandler) {
	// Public route for customers to fetch invoice by Order ID
	router.Get("/invoices/order/:orderId", invoiceHandler.GetInvoiceByOrderID)
	
	// Public route for Midtrans Webhook
	router.Post("/invoices/midtrans-callback", invoiceHandler.MidtransCallback)

	invoiceGroup := router.Group("/invoices", middleware.Protected())

	invoiceGroup.Post("/", invoiceHandler.GenerateInvoice)
	invoiceGroup.Get("/", invoiceHandler.GetAll)
	invoiceGroup.Get("/list", invoiceHandler.GetInvoices)
	invoiceGroup.Patch("/:id/pay", invoiceHandler.PayInvoice)
}
