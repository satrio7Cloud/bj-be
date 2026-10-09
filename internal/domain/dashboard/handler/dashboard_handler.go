package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/dashboard/usecase"
	"bj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	dashboardUsecase usecase.DashboardUsecase
}

func NewDashboardHandler(uc usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{
		dashboardUsecase: uc,
	}
}

// @Summary Get Dashboard Summary
// @Description Retrieve summary statistics for the dashboard
// @Tags dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/dashboard/summary [get]
func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	summary, err := h.dashboardUsecase.GetDashboardSummary()
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   summary,
	})
}

// @Summary Get Dashboard Orders
// @Description Retrieve a list of recent orders for the dashboard
// @Tags dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/dashboard/orders [get]
func (h *DashboardHandler) GetOrders(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	orders, total, err := h.dashboardUsecase.GetDashboardOrders(page, limit)
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  constants.StatusSuccess,
		"message": "Successfully retrieved dashboard orders",
		"data":    orders,
		"meta": fiber.Map{
			"page":          page,
			"limit":         limit,
			"total_records": total,
		},
	})
}

// @Summary Get Chart Data
// @Description Retrieve data for rendering dashboard charts
// @Tags dashboard
// @Produce json

// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/dashboard/chart [get]
func (h *DashboardHandler) GetChartData(c *fiber.Ctx) error {
	data, err := h.dashboardUsecase.GetChartData()
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   data,
	})
}

// @Summary Get Recent Activities
// @Description Retrieve a list of recent activities/logs
// @Tags dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/dashboard/activities [get]
func (h *DashboardHandler) GetRecentActivities(c *fiber.Ctx) error {
	activities, err := h.dashboardUsecase.GetRecentActivities()
	if err != nil {
		return utils.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": constants.StatusSuccess,
		"data":   activities,
	})
}
