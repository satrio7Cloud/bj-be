package routes

import (
	"bj/internal/domain/invoice/handler"

	"bj/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterInvoiceRoutes(router fiber.Router, invoiceHandler *handler.InvoiceHandler) {
	// Public route for customers to fetch invoice by Order ID
	router.Get("/invoices/order/:orderId", invoiceHandler.GetInvoiceByOrderID)
	router.Get("/invoices/order/:orderId/sync", invoiceHandler.SyncInvoice)
	router.Post("/invoices/order/:orderId/sync", invoiceHandler.SyncInvoice)
	
	// Public route for customers to generate invoice for their order
	router.Post("/invoices/customer", invoiceHandler.GenerateInvoice)
	router.Post("/invoices", invoiceHandler.GenerateInvoice)

	// Public sync endpoint by Order ID or Invoice ID
	router.Get("/invoices/:id/sync", invoiceHandler.SyncInvoice)
	router.Post("/invoices/:id/sync", invoiceHandler.SyncInvoice)

	// Public route for Midtrans Webhook
	router.Post("/invoices/midtrans-callback", invoiceHandler.MidtransCallback)

	invoiceGroup := router.Group("/invoices", middleware.Protected())

	invoiceGroup.Get("/", invoiceHandler.GetAll)
	invoiceGroup.Get("/list", invoiceHandler.GetInvoices)
	invoiceGroup.Patch("/:id/pay", invoiceHandler.PayInvoice)
	invoiceGroup.Post("/:id/sync", invoiceHandler.SyncInvoice)
}
