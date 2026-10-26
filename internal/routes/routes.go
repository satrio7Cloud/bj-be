package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"

	"bj/internal/constants"
	"bj/internal/di"
	biteshipRoutes "bj/internal/domain/biteship/routes"
	brandRoutes "bj/internal/domain/brand/routes"
	customerRoutes "bj/internal/domain/customer/routes"
	dashboardRoutes "bj/internal/domain/dashboard/routes"
	documentRoutes "bj/internal/domain/document/routes"
	invoiceRoutes "bj/internal/domain/invoice/routes"
	modelRoutes "bj/internal/domain/model/routes"
	notificationRoutes "bj/internal/domain/notification/routes"
	orderRoutes "bj/internal/domain/order/routes"
	paymentRoutes "bj/internal/domain/payment/routes"
	serviceRoutes "bj/internal/domain/service/routes"
	adminRoutes "bj/internal/domain/user/routes"
	vehicleRoutes "bj/internal/domain/vehicle/routes"
	companyRoutes "bj/internal/domain/company/routes"
)

func RegisterRoutes(app fiber.Router, c *di.Container) {
	app.Get("/cetak-ktl/:filename", func(ctx *fiber.Ctx) error {
		filename := ctx.Params("filename")
		objectName := "cetak-ktl/" + filename

		url, err := c.StorageService.GetPresignedURL(ctx.Context(), objectName, time.Hour)
		if err != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  constants.StatusError,
				"message": constants.MsgFileNotFound,
			})
		}
		if err := proxy.Do(ctx, url); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  constants.StatusError,
				"message": constants.MsgFetchFileFailed,
			})
		}

		return nil
	})

	api := app.Group("/api/v1")
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
	companyRoutes.RegisterCompanyProfileRoutes(api, c.CompanyHdl)
}
