package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// RequestLogger adalah middleware untuk mencatat (log) semua aktivitas request HTTP
// ke terminal. Berguna untuk memantau request yang masuk beserta waktu responnya.
func RequestLogger() fiber.Handler {
	return logger.New(logger.Config{
		// Format tampilan log di terminal:
		// [Waktu] METHOD /path - Status Code - Waktu Proses
		Format:     "[${time}] ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Local",
	})
}
