package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"


	"bj/internal/di"
	biteshipRoutes "bj/internal/domain/biteship/routes"
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
	notificationRoutes "bj/internal/domain/notification/routes"
	paymentRoutes "bj/internal/domain/payment/routes"
)

// RegisterRoutes adalah pusat pendaftaran semua routing di aplikasi.
// Konsep ini mengadaptasi pola dari proyek CashRecon Anda.
func RegisterRoutes(app fiber.Router, c *di.Container) {
	// Proxy redirect ke URL MinIO untuk gambar KTL
	app.Get("/cetak-ktl/:filename", func(ctx *fiber.Ctx) error {
		filename := ctx.Params("filename")
		objectName := "cetak-ktl/" + filename
		
		// Generate presigned URL yang valid selama 1 jam
		url, err := c.StorageService.GetPresignedURL(ctx.Context(), objectName, time.Hour)
		if err != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "File tidak ditemukan",
			})
		}
		// Stream/Proxy gambar asli dari MinIO secara transparan tanpa mengubah URL di browser
		if err := proxy.Do(ctx, url); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Gagal mengambil file dari server",
			})
		}
		
		return nil
	})

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
	notificationRoutes.RegisterNotificationRoutes(api, c.NotificationHdl)
	paymentRoutes.RegisterPaymentRoutes(api, c.PaymentHdl)
	biteshipRoutes.RegisterBiteshipRoutes(api, c.BiteshipHdl)
}
