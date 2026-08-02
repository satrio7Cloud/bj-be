package routes

import (
	"github.com/gofiber/fiber/v2"

	"bj/internal/di"
	adminRoutes "bj/internal/domain/user/routes"
	brandRoutes "bj/internal/domain/brand/routes"
	customerRoutes "bj/internal/domain/customer/routes"
	documentRoutes "bj/internal/domain/document/routes"
	modelRoutes "bj/internal/domain/model/routes"
	orderRoutes "bj/internal/domain/order/routes"
	serviceRoutes "bj/internal/domain/service/routes"
	vehicleRoutes "bj/internal/domain/vehicle/routes"
	invoiceRoutes "bj/internal/domain/invoice/routes"
	dashboardRoutes "bj/internal/domain/dashboard/routes"
)

// RegisterRoutes adalah pusat pendaftaran semua routing di aplikasi.
// Konsep ini mengadaptasi pola dari proyek CashRecon Anda.
func RegisterRoutes(app fiber.Router, c *di.Container) {
	// Buat grup utama /api/v1
	api := app.Group("/api/v1")

	// Daftarkan route untuk masing-masing domain
	adminRoutes.RegisterAdminRoutes(api, c.AdminHdl)
	brandRoutes.RegisterBrandRoutes(api, c.BrandHdl)
	customerRoutes.RegisterCustomerRoutes(api, c.CustomerHdl)
	documentRoutes.RegisterDocumentRoutes(api, c.DocumentHdl)
	modelRoutes.RegisterModelRoutes(api, c.ModelHdl)
	orderRoutes.RegisterOrderRoutes(api, c.OrderHdl)
	serviceRoutes.RegisterServiceRoutes(api, c.ServiceHdl)
	vehicleRoutes.RegisterVehicleRoutes(api, c.VehicleHdl)
	dashboardRoutes.RegisterDashboardRoutes(api, c.DashboardHdl)
	invoiceRoutes.RegisterInvoiceRoutes(api, c.InvoiceHdl)
}
