package repository

import (
	"bj/internal/domain/user/entity"

	"gorm.io/gorm"
)

type AdminRepository interface {
	FindByEmail(email string) (*entity.Admin, error)
	Create(admin *entity.Admin) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db}
}

func (r *adminRepository) FindByEmail(email string) (*entity.Admin, error) {
	var admin entity.Admin

	query := `SELECT CAST(id AS VARCHAR(36)) AS id, fullname, email, password, phone_number, is_active, created_at FROM dbo.Admins WHERE email = ?`
	res := r.db.Raw(query, email).Scan(&admin)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil // Return nil jika tidak ditemukan
	}

	return &admin, nil
}

func (r *adminRepository) Create(admin *entity.Admin) error {
	query := `INSERT INTO dbo.Admins (id, fullname, email, password, phone_number, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, admin.ID, admin.Fullname, admin.Email, admin.Password, admin.PhoneNumber, admin.IsActive, admin.CreatedAt).Error
}
