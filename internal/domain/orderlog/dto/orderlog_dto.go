package dto

import "time"

type OrderLogResponse struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	AdminID   string    `json:"admin_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}
