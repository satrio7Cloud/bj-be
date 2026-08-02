package main

import (
	"log"

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

	// Inisialisasi Fiber
	app := fiber.New(fiber.Config{
		AppName: "BJ-BE",
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
