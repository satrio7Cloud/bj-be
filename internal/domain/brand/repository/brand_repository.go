package repository

import (
	"bj/internal/domain/brand/entity"

	"gorm.io/gorm"
)

type BrandRepository interface {
	Create(brand *entity.Brand) error
	FindAll() ([]entity.Brand, error)
	FindByID(id string) (*entity.Brand, error)
	Update(brand *entity.Brand) error
	Delete(id string) error
}

type brandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) BrandRepository {
	return &brandRepository{db}
}

func (r *brandRepository) Create(brand *entity.Brand) error {
	query := `INSERT INTO dbo.MasterVehicleBrands (ID, BrandName) VALUES (?, ?)`
	return r.db.Exec(query, brand.ID, brand.BrandName).Error
}

func (r *brandRepository) FindAll() ([]entity.Brand, error) {
	var brands []entity.Brand
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, BrandName FROM dbo.MasterVehicleBrands ORDER BY BrandName ASC`
	err := r.db.Raw(query).Scan(&brands).Error
	return brands, err
}

func (r *brandRepository) FindByID(id string) (*entity.Brand, error) {
	var brand entity.Brand
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, BrandName FROM dbo.MasterVehicleBrands WHERE ID = ?`
	res := r.db.Raw(query, id).Scan(&brand)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &brand, nil
}

func (r *brandRepository) Update(brand *entity.Brand) error {
	query := `UPDATE dbo.MasterVehicleBrands SET BrandName = ? WHERE ID = ?`
	return r.db.Exec(query, brand.BrandName, brand.ID).Error
}

func (r *brandRepository) Delete(id string) error {
	query := `DELETE FROM dbo.MasterVehicleBrands WHERE ID = ?`
	return r.db.Exec(query, id).Error
}
