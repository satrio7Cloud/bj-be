package usecase

import (
	"bj/internal/domain/dashboard/dto"
	"bj/internal/domain/dashboard/repository"
)

type DashboardUsecase interface {
	GetDashboardSummary() (dto.DashboardSummaryResponse, error)
	GetDashboardOrders(page, limit int) ([]dto.DashboardOrderResponse, int64, error)
	GetChartData() ([]dto.DashboardChartResponse, error)
	GetRecentActivities() ([]dto.DashboardActivityResponse, error)
}

type dashboardUsecase struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardUsecase(repo repository.DashboardRepository) DashboardUsecase {
	return &dashboardUsecase{dashboardRepo: repo}
}

func (u *dashboardUsecase) GetDashboardSummary() (dto.DashboardSummaryResponse, error) {
	summary, err := u.dashboardRepo.GetSummary()
	if err != nil {
		return summary, err
	}

	pendingActions, err := u.dashboardRepo.GetPendingActions()
	if err == nil {
		summary.PendingActions = pendingActions
	}

	topServices, err := u.dashboardRepo.GetTopServices()
	if err == nil {
		summary.TopServices = topServices
	}

	return summary, nil
}

func (u *dashboardUsecase) GetDashboardOrders(page, limit int) ([]dto.DashboardOrderResponse, int64, error) {
	orders, total, err := u.dashboardRepo.GetOrders(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.DashboardOrderResponse
	for _, o := range orders {
		responses = append(responses, dto.DashboardOrderResponse{
			OrderID:      o.OrderID,
			OrderDate:    o.OrderDate,
			Status:       o.Status,
			ServiceLevel: o.ServiceLevel,
			TotalCost:    o.TotalCost,
			CustomerName: o.CustomerName,
			PhoneNumber:  o.PhoneNumber,
			PlateNumber:  o.PlateNumber,
			ServiceName:  o.ServiceName,
		})
	}

	return responses, total, nil
}

func (u *dashboardUsecase) GetChartData() ([]dto.DashboardChartResponse, error) {
	return u.dashboardRepo.GetChartData()
}

func (u *dashboardUsecase) GetRecentActivities() ([]dto.DashboardActivityResponse, error) {
	return u.dashboardRepo.GetRecentActivities()
}
