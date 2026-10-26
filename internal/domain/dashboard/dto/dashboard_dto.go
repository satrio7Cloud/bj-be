package dto

import "time"

type DashboardSummaryResponse struct {
	// Existing Metrics
	TotalPendingOrders   int64   `json:"total_pending_orders"`
	TotalExpressOrders   int64   `json:"total_express_orders"`
	TotalCompletedOrders int64   `json:"total_completed_orders"`
	TotalRevenue         float64 `json:"total_revenue"`

	// KTL Metrics
	TotalPendingKTLs   int64   `json:"total_pending_ktls"`
	TotalCompletedKTLs int64   `json:"total_completed_ktls"`
	TotalRevenueKTLs   float64 `json:"total_revenue_ktls"`

	// New Metrics
	PendingActions PendingActionsResponse `json:"pending_actions"`
	TopServices    []TopServiceResponse   `json:"top_services"`
}

type PendingActionsResponse struct {
	UnverifiedDocuments int64 `json:"unverified_documents"`
	StuckOrders         int64 `json:"stuck_orders"`
	StuckKTLs           int64 `json:"stuck_ktls"`
}

type TopServiceResponse struct {
	ServiceName string `json:"service_name"`
	OrderCount  int64  `json:"order_count"`
}

type DashboardChartResponse struct {
	Date       string  `json:"date"` // Format: YYYY-MM-DD
	OrderCount int64   `json:"order_count"`
	Revenue    float64 `json:"revenue"`
}

type DashboardActivityResponse struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	AdminName string    `json:"admin_name"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type DashboardOrderResponse struct {
	OrderID      string    `json:"order_id"`
	OrderDate    time.Time `json:"order_date"`
	Status       string    `json:"status"`
	ServiceLevel string    `json:"service_level"`
	TotalCost    float64   `json:"total_cost"`
	CustomerName string    `json:"customer_name"`
	PhoneNumber  string    `json:"phone_number"`
	PlateNumber  string    `json:"plate_number"`
	ServiceName    string    `json:"service_name"`
	UnverifiedDocs int       `json:"unverified_docs"`
}

