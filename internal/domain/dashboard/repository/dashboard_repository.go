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
	GetOrders(page, limit int) ([]entity.DashboardOrder, int64, error)
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

func (r *dashboardRepository) GetOrders(page, limit int) ([]entity.DashboardOrder, int64, error) {
	var orders []entity.DashboardOrder
	var total int64

	baseQuery := `FROM vw_DashboardOrderList`

	countQuery := `SELECT COUNT(*) ` + baseQuery
	if err := r.db.Raw(countQuery).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	dataQuery := `
		SELECT CAST(OrderID AS VARCHAR(36)) AS OrderID, OrderDate, Status, ServiceLevel, TotalCost, CustomerName, PhoneNumber, PlateNumber, ServiceName 
		` + baseQuery + ` 
		ORDER BY OrderDate DESC 
		OFFSET ? ROWS FETCH NEXT ? ROWS ONLY
	`
	
	if err := r.db.Raw(dataQuery, offset, limit).Scan(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
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
		WITH CombinedServices AS (
			SELECT S.ServiceName
			FROM dbo.Orders O
			JOIN dbo.MasterServices S ON O.ServiceID = S.ID
			
			UNION ALL
			
			SELECT 'Cetak KTL' AS ServiceName
			FROM dbo.OrderKTLs
		)
		SELECT TOP 5 ServiceName, COUNT(*) as OrderCount
		FROM CombinedServices
		GROUP BY ServiceName
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
			WHERE OrderDate >= CAST(DATEADD(day, -30, GETDATE()) AS DATE)
			UNION ALL
			SELECT OrderDate, ID, Fee as Revenue FROM dbo.OrderKTLs
			WHERE OrderDate >= CAST(DATEADD(day, -30, GETDATE()) AS DATE)
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
