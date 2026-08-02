package repository

import (
	"bj/internal/domain/dashboard/dto"
	"bj/internal/domain/dashboard/entity"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetSummary() (dto.DashboardSummaryResponse, error)
	GetPendingActions() (dto.PendingActionsResponse, error)
	GetTopServices() ([]dto.TopServiceResponse, error)
	GetChartData() ([]dto.DashboardChartResponse, error)
	GetRecentActivities() ([]dto.DashboardActivityResponse, error)
	GetOrders() ([]entity.DashboardOrder, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetSummary() (dto.DashboardSummaryResponse, error) {
	var summary dto.DashboardSummaryResponse

	// Total Pending Orders
	err := r.db.Table("dbo.Orders").Where("Status = ?", "Pending").Count(&summary.TotalPendingOrders).Error
	if err != nil {
		return summary, err
	}

	// Total Express Orders
	err = r.db.Table("dbo.Orders").Where("ServiceLevel = ?", "EXPRESS").Count(&summary.TotalExpressOrders).Error
	if err != nil {
		return summary, err
	}

	// Total Completed Orders (This Month)
	// Using SQL Server specific MONTH() and YEAR() functions
	err = r.db.Table("dbo.Orders").
		Where("Status = ? AND MONTH(OrderDate) = MONTH(GETDATE()) AND YEAR(OrderDate) = YEAR(GETDATE())", "Completed").
		Count(&summary.TotalCompletedOrders).Error
	if err != nil {
		return summary, err
	}

	// Total Revenue (This Month) - Sum of TotalCost for Completed Orders
	// We use NullFloat64 in case there are no completed orders this month to avoid errors
	var totalRevenue *float64
	err = r.db.Table("dbo.Orders").
		Select("SUM(TotalCost)").
		Where("Status = ? AND MONTH(OrderDate) = MONTH(GETDATE()) AND YEAR(OrderDate) = YEAR(GETDATE())", "Completed").
		Scan(&totalRevenue).Error
	if err != nil {
		return summary, err
	}
	
	if totalRevenue != nil {
		summary.TotalRevenue = *totalRevenue
	} else {
		summary.TotalRevenue = 0
	}

	return summary, nil
}

func (r *dashboardRepository) GetOrders() ([]entity.DashboardOrder, error) {
	var orders []entity.DashboardOrder
	err := r.db.Select("CAST(OrderID AS VARCHAR(36)) AS OrderID, OrderDate, Status, ServiceLevel, TotalCost, CustomerName, PhoneNumber, PlateNumber, ServiceName").
		Order("OrderDate DESC").Find(&orders).Error
	return orders, err
}

func (r *dashboardRepository) GetPendingActions() (dto.PendingActionsResponse, error) {
	var res dto.PendingActionsResponse
	
	// Unverified Documents (IsVerified = 0)
	err := r.db.Table("dbo.OrderDocuments").Where("IsVerified = ? AND DeletedAt IS NULL", 0).Count(&res.UnverifiedDocuments).Error
	if err != nil {
		return res, err
	}

	// Stuck Orders (Status = Processing AND OrderDate > 3 days ago)
	// SQL Server specific DATEDIFF
	err = r.db.Table("dbo.Orders").
		Where("Status = ? AND DATEDIFF(day, OrderDate, GETDATE()) >= ?", "Processing", 3).
		Count(&res.StuckOrders).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *dashboardRepository) GetTopServices() ([]dto.TopServiceResponse, error) {
	var res []dto.TopServiceResponse
	query := `
		SELECT TOP 5 S.ServiceName, COUNT(O.ID) as OrderCount
		FROM dbo.Orders O
		JOIN dbo.MasterServices S ON O.ServiceID = S.ID
		GROUP BY S.ServiceName
		ORDER BY OrderCount DESC
	`
	err := r.db.Raw(query).Scan(&res).Error
	return res, err
}

func (r *dashboardRepository) GetChartData() ([]dto.DashboardChartResponse, error) {
	var res []dto.DashboardChartResponse
	// Last 30 days of data grouped by date
	query := `
		SELECT TOP 30
			CAST(OrderDate AS DATE) as Date,
			COUNT(ID) as OrderCount,
			SUM(TotalCost) as Revenue
		FROM dbo.Orders
		GROUP BY CAST(OrderDate AS DATE)
		ORDER BY Date DESC
	`
	err := r.db.Raw(query).Scan(&res).Error
	return res, err
}

func (r *dashboardRepository) GetRecentActivities() ([]dto.DashboardActivityResponse, error) {
	var res []dto.DashboardActivityResponse
	query := `
		SELECT TOP 10 
			CAST(L.ID AS VARCHAR(36)) AS ID, 
			CAST(L.OrderID AS VARCHAR(36)) AS OrderID, 
			A.Fullname AS AdminName, 
			L.OldStatus, 
			L.NewStatus, 
			L.Notes, 
			L.CreatedAt
		FROM dbo.OrderLogs L
		LEFT JOIN dbo.Admins A ON L.AdminID = A.ID
		ORDER BY L.CreatedAt DESC
	`
	err := r.db.Raw(query).Scan(&res).Error
	return res, err
}
