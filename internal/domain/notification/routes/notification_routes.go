package routes

import (
	"bj/internal/domain/notification/handler"
	"bj/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterNotificationRoutes(router fiber.Router, hdl *handler.NotificationHandler) {
	notificationGroup := router.Group("/notifications", middleware.Protected())
	notificationGroup.Get("/unread", hdl.GetUnread)
	notificationGroup.Patch("/:id/read", hdl.MarkAsRead)
}
