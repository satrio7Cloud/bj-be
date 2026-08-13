package main

import (
	"log"

	"strings"

	"bj/internal/config"
	"bj/internal/di"
	appRoutes "bj/internal/routes"
	"bj/middleware"

	openapimiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

// @title Biro Jasa API
// @version 1.0
// @description REST API untuk sistem Biro Jasa
// @host localhost:3000
// @BasePath /
func main() {

	cfg := config.LoadConfig()

	// -- WIRING DEPENDENCY INJECTION MENGGUNAKAN CONTAINER --
	container, err := di.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Gagal membangun container aplikasi: %v", err)
	}

	// Inisialisasi Fiber dengan Custom ErrorHandler
	app := fiber.New(fiber.Config{
		AppName: "BJ-BE",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			errMsg := err.Error()
			// Log error mentah di console server
			log.Printf("[INTERNAL ERROR] %v\n", err)

			// Sanitasi pesan error teknis/TCP agar tidak bocor ke client
			if strings.Contains(errMsg, "tcp") || strings.Contains(errMsg, "dial") || strings.Contains(errMsg, "sqlserver") || strings.Contains(errMsg, "i/o timeout") || strings.Contains(errMsg, "connection refused") {
				errMsg = "Layanan database backend sedang tidak dapat dijangkau. Silakan periksa koneksi atau coba beberapa saat lagi."
			}

			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": errMsg,
			})
		},
	})

	// Mengaktifkan Middleware CORS & Logger
	app.Use(middleware.CorsConfig())
	app.Use(middleware.RequestLogger())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Server is Running Success")
	})

	// Setup API Routes terpusat (seperti di CashRecon)
	appRoutes.RegisterRoutes(app, container)

	// Route untuk melayani file swagger.yaml dari folder docs
	app.Get("/docs/swagger.yaml", func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Set("Pragma", "no-cache")
		c.Set("Expires", "0")
		return c.SendFile("./docs/swagger.yaml")
	})

	// Setup Swagger UI menggunakan go-openapi middleware
	app.Use("/docs", adaptor.HTTPHandler(openapimiddleware.SwaggerUI(openapimiddleware.SwaggerUIOpts{
		BasePath: "/",
		Path:     "docs",
		SpecURL:  "/docs/swagger.yaml",
		Title:    "Biro Jasa API Docs",
	}, nil)))

	port := ":" + cfg.AppPort
	log.Printf("Server di jalankan di port %s...\n", cfg.AppPort)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Error saat menjalankan Server: %v", err)
	}
}
