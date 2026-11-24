package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiter mengkonfigurasi rate limiting untuk mencegah spam, brute force, dan DDoS
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,             // Maksimal 100 request(100 request/menit/IP belum tentu cocok untuk semua endpoint.)
		Expiration: 1 * time.Minute, // Per 1 menit per IP
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Terlalu banyak permintaan. Silakan coba beberapa saat lagi.",
			})
		},
		Next: func(c *fiber.Ctx) bool {
			// Lewati preflight OPTIONS, health check, dan swagger docs agar tidak terblokir
			if c.Method() == fiber.MethodOptions {
				return true
			}
			path := c.Path()
			return path == "/health" ||
				strings.HasPrefix(path, "/docs") ||
				path == "/swagger.yaml" ||
				strings.HasPrefix(path, "/uploads") ||
				path == "/api/v1/invoices/midtrans-callback"
		},
	})
}
