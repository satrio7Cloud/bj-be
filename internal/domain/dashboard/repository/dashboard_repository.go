package repository

import (
	"fmt"
	"time"

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
	GetOrders(page, limit int, status string) ([]entity.DashboardOrder, int64, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetSummary() (dto.DashboardSummaryResponse, error) {
	var summary dto.DashboardSummaryResponse

	now := time.Now()
	startOfMonth := time.Date(
		now.Year(), now.Month(), 1,
		0, 0, 0, 0, now.Location(),
	)
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

	// Aggregate Orders dalam satu query.
	var orders struct {
		TotalPendingOrders   int64   `gorm:"column:total_pending_orders"`
		TotalExpressOrders   int64   `gorm:"column:total_express_orders"`
		TotalCompletedOrders int64   `gorm:"column:total_completed_orders"`
		TotalRevenue         float64 `gorm:"column:total_revenue"`
	}

	err := r.db.Table("dbo.Orders WITH (NOLOCK)").
		Select(`
			COALESCE(SUM(CASE WHEN Status = ? THEN 1 ELSE 0 END), 0) AS total_pending_orders,
			COALESCE(SUM(CASE WHEN ServiceLevel = ? THEN 1 ELSE 0 END), 0) AS total_express_orders,
			COALESCE(SUM(CASE WHEN Status = ?
				AND OrderDate >= ? AND OrderDate < ?
				THEN 1 ELSE 0 END), 0) AS total_completed_orders,
			COALESCE(SUM(CASE WHEN Status = ?
				AND OrderDate >= ? AND OrderDate < ?
				THEN TotalCost ELSE 0 END), 0) AS total_revenue
		`,
			"Pending",
			"EXPRESS",
			"Completed", startOfMonth, startOfNextMonth,
			"Completed", startOfMonth, startOfNextMonth,
		).
		Scan(&orders).Error
	if err != nil {
		return summary, fmt.Errorf("get orders summary: %w", err)
	}

	// Aggregate OrderKTLs dalam satu query.
	var ktls struct {
		TotalPendingKTLs   int64   `gorm:"column:total_pending_ktls"`
		TotalCompletedKTLs int64   `gorm:"column:total_completed_ktls"`
		TotalRevenueKTLs   float64 `gorm:"column:total_revenue_ktls"`
	}

	err = r.db.Table("dbo.OrderKTLs WITH (NOLOCK)").
		Select(`
			COALESCE(SUM(CASE WHEN Status = ? THEN 1 ELSE 0 END), 0) AS total_pending_ktls,
			COALESCE(SUM(CASE WHEN Status = ?
				AND OrderDate >= ? AND OrderDate < ?
				THEN 1 ELSE 0 END), 0) AS total_completed_ktls,
			COALESCE(SUM(CASE WHEN Status = ?
				AND OrderDate >= ? AND OrderDate < ?
				THEN Fee ELSE 0 END), 0) AS total_revenue_ktls
		`,
			"Pending",
			"Completed", startOfMonth, startOfNextMonth,
			"Completed", startOfMonth, startOfNextMonth,
		).
		Scan(&ktls).Error
	if err != nil {
		return summary, fmt.Errorf("get KTL summary: %w", err)
	}

	summary.TotalPendingOrders = orders.TotalPendingOrders
	summary.TotalExpressOrders = orders.TotalExpressOrders
	summary.TotalCompletedOrders = orders.TotalCompletedOrders
	summary.TotalRevenue = orders.TotalRevenue

	summary.TotalPendingKTLs = ktls.TotalPendingKTLs
	summary.TotalCompletedKTLs = ktls.TotalCompletedKTLs
	summary.TotalRevenueKTLs = ktls.TotalRevenueKTLs

	return summary, nil
}

func (r *dashboardRepository) GetOrders(page, limit int, status string) ([]entity.DashboardOrder, int64, error) {
	var orders []entity.DashboardOrder
	var total int64

	baseQuery := `FROM vw_DashboardOrderList`

	countQuery := `SELECT COUNT(*) ` + baseQuery
	if err := r.db.Raw(countQuery).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	dataQuery := `
		SELECT CAST(OrderID AS VARCHAR(36)) AS OrderID, OrderDate, Status, ServiceLevel, TotalCost, CustomerName, PhoneNumber, PlateNumber, ServiceName, UnverifiedDocs 
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
	r.db.Table("dbo.OrderDocuments WITH (NOLOCK)").Where("IsVerified = ? AND DeletedAt IS NULL", 0).Count(&res.UnverifiedDocuments)

	// Stuck Orders: Menggunakan batas waktu (>= 3 hari lalu) tanpa fungsi SQL pada kolom (SARGable)
	threeDaysAgo := time.Now().AddDate(0, 0, -3)

	r.db.Table("dbo.Orders WITH (NOLOCK)").
		Where("Status = ? AND OrderDate <= ?", "Processing", threeDaysAgo).
		Count(&res.StuckOrders)
	r.db.Table("dbo.OrderKTLs WITH (NOLOCK)").
		Where("Status = ? AND OrderDate <= ?", "Processing", threeDaysAgo).
		Count(&res.StuckKTLs)

	return res, nil
}

func (r *dashboardRepository) GetTopServices() ([]dto.TopServiceResponse, error) {
	var res []dto.TopServiceResponse
	query := `
		WITH CombinedServices AS (
			SELECT S.ServiceName
			FROM dbo.Orders O WITH (NOLOCK)
			JOIN dbo.MasterServices S WITH (NOLOCK) ON O.ServiceID = S.ID
			
			UNION ALL
			
			SELECT 'Cetak KTL' AS ServiceName
			FROM dbo.OrderKTLs WITH (NOLOCK)
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

	// SARGable parameter: 30 hari yang lalu (menghindari fungsi GETDATE() berulang di SQL)
	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	// Optimasi: Agregasi per tabel terlebih dahulu sebelum di-UNION ALL
	// Ini mencegah pembentukan temporary table berukuran besar dan mengurangi dataset UNION menjadi maksimal 60 baris saja
	query := `
		WITH DailyOrders AS (
			SELECT 
				CAST(OrderDate AS DATE) AS [Date],
				COUNT(*) AS OrderCount,
				SUM(TotalCost) AS Revenue
			FROM dbo.Orders WITH (NOLOCK)
			WHERE OrderDate >= ?
			GROUP BY CAST(OrderDate AS DATE)
		),
		DailyKTLs AS (
			SELECT 
				CAST(OrderDate AS DATE) AS [Date],
				COUNT(*) AS OrderCount,
				SUM(Fee) AS Revenue
			FROM dbo.OrderKTLs WITH (NOLOCK)
			WHERE OrderDate >= ?
			GROUP BY CAST(OrderDate AS DATE)
		),
		CombinedDaily AS (
			SELECT [Date], OrderCount, Revenue FROM DailyOrders
			UNION ALL
			SELECT [Date], OrderCount, Revenue FROM DailyKTLs
		)
		SELECT TOP 30
			CONVERT(VARCHAR(10), [Date], 120) AS Date,
			SUM(OrderCount) AS OrderCount,
			COALESCE(SUM(Revenue), 0) AS Revenue
		FROM CombinedDaily
		GROUP BY [Date]
		ORDER BY [Date] DESC
	`
	err := r.db.Raw(query, startDate, startDate).Scan(&res).Error
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
		FROM dbo.OrderLogs L WITH (NOLOCK)
		LEFT JOIN dbo.Admins A WITH (NOLOCK) ON L.AdminID = A.ID
		ORDER BY L.CreatedAt DESC
	`
	err := r.db.Raw(query).Scan(&res).Error
	return res, err
}
