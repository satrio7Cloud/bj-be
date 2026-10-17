package middleware

import (
	"os"
	"strings"

	"bj/internal/constants"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  constants.StatusError,
				"message": constants.MsgUnauthorized,
			})
		}

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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  constants.StatusError,
				"message": "Server configuration error: JWT_SECRET tidak diset",
			})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  constants.StatusError,
				"message": constants.MsgUnauthorized,
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("admin_id", claims["id"])
			c.Locals("admin_email", claims["email"])
			c.Locals("admin_role", claims["role"])
		}

		return c.Next()
	}
}
