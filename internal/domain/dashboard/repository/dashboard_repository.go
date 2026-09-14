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
	r.db.Table("dbo.Orders").Where("Status = ?", "Pending").Count(&summary.TotalPendingOrders)
	r.db.Table("dbo.OrderKTLs").Where("Status = ?", "Pending").Count(&summary.TotalPendingKTLs)

	// Total Express Orders (Only in normal Orders)
	r.db.Table("dbo.Orders").Where("ServiceLevel = ?", "EXPRESS").Count(&summary.TotalExpressOrders)

	// Total Completed Orders (This Month)
	r.db.Table("dbo.Orders").
		Where("Status = ? AND MONTH(OrderDate) = MONTH(GETDATE()) AND YEAR(OrderDate) = YEAR(GETDATE())", "Completed").
		Count(&summary.TotalCompletedOrders)
	r.db.Table("dbo.OrderKTLs").
		Where("Status = ? AND MONTH(OrderDate) = MONTH(GETDATE()) AND YEAR(OrderDate) = YEAR(GETDATE())", "Completed").
		Count(&summary.TotalCompletedKTLs)

	// Total Revenue (This Month)
	var revenueOrders, revenueKTLs *float64
	r.db.Table("dbo.Orders").Select("SUM(TotalCost)").
		Where("Status = ? AND MONTH(OrderDate) = MONTH(GETDATE()) AND YEAR(OrderDate) = YEAR(GETDATE())", "Completed").
		Scan(&revenueOrders)
	
	r.db.Table("dbo.OrderKTLs").Select("SUM(Fee)").
		Where("Status = ? AND MONTH(OrderDate) = MONTH(GETDATE()) AND YEAR(OrderDate) = YEAR(GETDATE())", "Completed").
		Scan(&revenueKTLs)

	if revenueOrders != nil {
		summary.TotalRevenue = *revenueOrders
	}
	if revenueKTLs != nil {
		summary.TotalRevenueKTLs = *revenueKTLs
	}

	return summary, nil
}

func (r *dashboardRepository) GetOrders() ([]entity.DashboardOrder, error) {
	var orders []entity.DashboardOrder
	query := `
		SELECT * FROM (
			SELECT CAST(OrderID AS VARCHAR(36)) AS OrderID, OrderDate, Status, ServiceLevel, TotalCost, CustomerName, PhoneNumber, PlateNumber, ServiceName
			FROM vw_DashboardOrderList
			UNION ALL
			SELECT CAST(ID AS VARCHAR(36)) AS OrderID, OrderDate, Status, 'STANDARD' AS ServiceLevel, Fee AS TotalCost, CustomerName, PhoneNumber, '-' AS PlateNumber, 'Cetak KTL' AS ServiceName
			FROM dbo.OrderKTLs
		) AS CombinedOrders
		ORDER BY OrderDate DESC
	`
	err := r.db.Raw(query).Scan(&orders).Error
	return orders, err
}

func (r *dashboardRepository) GetPendingActions() (dto.PendingActionsResponse, error) {
	var res dto.PendingActionsResponse
	
	// Unverified Documents (IsVerified = 0)
	r.db.Table("dbo.OrderDocuments").Where("IsVerified = ? AND DeletedAt IS NULL", 0).Count(&res.UnverifiedDocuments)

	// Stuck Orders (Status = Processing AND OrderDate > 3 days ago)
	r.db.Table("dbo.Orders").
		Where("Status = ? AND DATEDIFF(day, OrderDate, GETDATE()) >= ?", "Processing", 3).
		Count(&res.StuckOrders)
	r.db.Table("dbo.OrderKTLs").
		Where("Status = ? AND DATEDIFF(day, OrderDate, GETDATE()) >= ?", "Processing", 3).
		Count(&res.StuckKTLs)

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
		WITH CombinedData AS (
			SELECT OrderDate, ID, TotalCost as Revenue FROM dbo.Orders
			UNION ALL
			SELECT OrderDate, ID, Fee as Revenue FROM dbo.OrderKTLs
		)
		SELECT TOP 30
			CAST(OrderDate AS DATE) as Date,
			COUNT(ID) as OrderCount,
			SUM(Revenue) as Revenue
		FROM CombinedData
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
