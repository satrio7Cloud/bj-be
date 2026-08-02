package middleware

import (
	"os"
	"strings"

	"bj/internal/constants"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected adalah middleware untuk memvalidasi JWT token
// pada setiap endpoint yang membutuhkan autentikasi (harus login).
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  constants.StatusError,
				"message": constants.MsgUnauthorized,
			})
		}

		// 2. Pastikan formatnya "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  constants.StatusError,
				"message": "Format token tidak valid. Gunakan format: Bearer <token>",
			})
		}

		tokenString := parts[1]
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			// Fail-closed: jangan pernah pakai secret default di production.
			// Jika JWT_SECRET tidak diset, tolak semua request (aman daripada celah).
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  constants.StatusError,
				"message": "Server configuration error: JWT_SECRET tidak diset",
			})
		}

		// 3. Parse dan validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing yang digunakan adalah HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})

		// 4. Cek apakah ada error atau token sudah kadaluarsa
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  constants.StatusError,
				"message": constants.MsgUnauthorized,
			})
		}

		// 5. Jika sukses, ekstrak informasi (claims) dari token
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Simpan ID, Email, dan Role ke dalam Fiber Context (c.Locals)
			// agar bisa dibaca oleh Handler selanjutnya tanpa harus hit database lagi!
			c.Locals("admin_id", claims["id"])
			c.Locals("admin_email", claims["email"])
			c.Locals("admin_role", claims["role"])
		}

		// 6. Lanjut ke handler berikutnya (proses bisnis)
		return c.Next()
	}
}
