package entity

import "time"

type Admin struct {
	ID          string     `gorm:"column:id;primaryKey" json:"id"`
	Fullname    string     `gorm:"column:fullname" json:"fullname"`
	Email       string     `gorm:"column:email" json:"email"`
	Password    string     `gorm:"column:password" json:"-"`
	PhoneNumber string     `gorm:"column:phone_number" json:"phone_number"`
	IsActive    bool       `gorm:"column:is_active" json:"is_active"`
	LastLoginAt *time.Time `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}

// TableName overrides the table name used by Admin to `Admins`
func (Admin) TableName() string {
	return "Admins"
}
