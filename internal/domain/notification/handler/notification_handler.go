package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/notification/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	uc usecase.NotificationUsecase
}

func NewNotificationHandler(uc usecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{uc: uc}
}

// @Summary Get Unread Notifications
// @Description Retrieve a list of unread notifications
// @Tags notification
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/notifications/unread [get]
func (h *NotificationHandler) GetUnread(c *fiber.Ctx) error {
	// Default limit 10
	limit := c.QueryInt("limit", 10)
	notifs, err := h.uc.GetUnreadNotifications(limit)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   notifs,
	})
}

// @Summary Mark Notification as Read
// @Description Mark a notification as read
// @Tags notification
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/notifications/{id}/read [patch]
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.uc.MarkNotificationAsRead(id)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.StatusSuccess,
		"message": "Notification marked as read",
	})
}
