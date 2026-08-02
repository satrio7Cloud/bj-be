package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CorsConfig mengatur kebijakan CORS (Cross-Origin Resource Sharing)
// agar API bisa diakses dari Frontend (React/Vue/Angular) yang beda port.
func CorsConfig() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*", // Wildcard origin
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false, // Wajib false jika AllowOrigins menggunakan "*"
	})
}
