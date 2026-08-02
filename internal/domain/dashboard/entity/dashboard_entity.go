package entity

import "time"

// DashboardOrder represents a flat structure fetched from vw_DashboardOrderList view
// It is read-only and used for displaying order summary in the dashboard
type DashboardOrder struct {
	OrderID      string    `gorm:"column:OrderID;primaryKey"`
	OrderDate    time.Time `gorm:"column:OrderDate"`
	Status       string    `gorm:"column:Status"`
	ServiceLevel string    `gorm:"column:ServiceLevel"`
	TotalCost    float64   `gorm:"column:TotalCost"`
	
	CustomerName string    `gorm:"column:CustomerName"`
	PhoneNumber  string    `gorm:"column:PhoneNumber"`
	
	PlateNumber  string    `gorm:"column:PlateNumber"`
	ServiceName  string    `gorm:"column:ServiceName"`
}

// TableName overrides the table name used by GORM to map to the SQL View
func (DashboardOrder) TableName() string {
	return "vw_DashboardOrderList"
}
