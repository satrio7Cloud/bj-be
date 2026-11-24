package usecase

import (
	"context"
	"time"

	"bj/internal/domain/dashboard/dto"
	"bj/internal/domain/dashboard/repository"
	"bj/internal/shared/cache"

	"github.com/redis/go-redis/v9"
)

type DashboardUsecase interface {
	GetDashboardSummary() (dto.DashboardSummaryResponse, error)
	GetDashboardOrders(page, limit int, status string) ([]dto.DashboardOrderResponse, int64, error)
	GetChartData() ([]dto.DashboardChartResponse, error)
	GetRecentActivities() ([]dto.DashboardActivityResponse, error)
}

type dashboardUsecase struct {
	dashboardRepo repository.DashboardRepository
	redisClient   *redis.Client
}

func NewDashboardUsecase(repo repository.DashboardRepository, redisClient *redis.Client) DashboardUsecase {
	return &dashboardUsecase{
		dashboardRepo: repo,
		redisClient:   redisClient,
	}
}

func (u *dashboardUsecase) GetDashboardSummary() (dto.DashboardSummaryResponse, error) {
	return cache.GetOrSet(context.Background(), u.redisClient, "dashboard:summary", 1*time.Minute, func() (dto.DashboardSummaryResponse, error) {
		summary, err := u.dashboardRepo.GetSummary()
		if err != nil {
			return summary, err
		}

		if pendingActions, err := u.dashboardRepo.GetPendingActions(); err == nil {
			summary.PendingActions = pendingActions
		}

		if topServices, err := u.dashboardRepo.GetTopServices(); err == nil {
			summary.TopServices = topServices
		}

		return summary, nil
	})
}

func (u *dashboardUsecase) GetDashboardOrders(page, limit int, status string) ([]dto.DashboardOrderResponse, int64, error) {
	orders, total, err := u.dashboardRepo.GetOrders(page, limit, status)
	if err != nil {
		return nil, 0, err
	}

	return dto.ToDashboardOrderResponses(orders), total, nil
}

func (u *dashboardUsecase) GetChartData() ([]dto.DashboardChartResponse, error) {
	return cache.GetOrSet(context.Background(), u.redisClient, "dashboard:chart_data", 5*time.Minute, func() ([]dto.DashboardChartResponse, error) {
		return u.dashboardRepo.GetChartData()
	})
}

func (u *dashboardUsecase) GetRecentActivities() ([]dto.DashboardActivityResponse, error) {
	return u.dashboardRepo.GetRecentActivities()
}
