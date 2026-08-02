package repository

import (
	"bj/internal/domain/model/entity"

	"gorm.io/gorm"
)

type ModelRepository interface {
	Create(model *entity.Model) error
	FindAll() ([]entity.Model, error)
	FindByBrandID(brandID string) ([]entity.Model, error)
	FindByID(id string) (*entity.Model, error)
	Update(model *entity.Model) error
	Delete(id string) error
}

type modelRepository struct {
	db *gorm.DB
}

func NewModelRepository(db *gorm.DB) ModelRepository {
	return &modelRepository{db}
}

func (r *modelRepository) Create(model *entity.Model) error {
	query := `INSERT INTO dbo.MasterVehicleModels (ID, BrandID, ModelName, VehicleType) VALUES (?, ?, ?, ?)`
	return r.db.Exec(query, model.ID, model.BrandID, model.ModelName, model.VehicleType).Error
}

func (r *modelRepository) FindAll() ([]entity.Model, error) {
	var models []entity.Model
	query := `SELECT CAST(M.ID AS VARCHAR(36)) AS ID, CAST(M.BrandID AS VARCHAR(36)) AS BrandID, B.BrandName, M.ModelName, M.VehicleType FROM dbo.MasterVehicleModels M JOIN dbo.MasterVehicleBrands B ON M.BrandID = B.ID ORDER BY M.ModelName ASC`
	err := r.db.Raw(query).Scan(&models).Error
	return models, err
}

func (r *modelRepository) FindByBrandID(brandID string) ([]entity.Model, error) {
	var models []entity.Model
	query := `SELECT CAST(M.ID AS VARCHAR(36)) AS ID, CAST(M.BrandID AS VARCHAR(36)) AS BrandID, B.BrandName, M.ModelName, M.VehicleType FROM dbo.MasterVehicleModels M JOIN dbo.MasterVehicleBrands B ON M.BrandID = B.ID WHERE M.BrandID = ? ORDER BY M.ModelName ASC`
	err := r.db.Raw(query, brandID).Scan(&models).Error
	return models, err
}

func (r *modelRepository) FindByID(id string) (*entity.Model, error) {
	var model entity.Model
	query := `SELECT CAST(M.ID AS VARCHAR(36)) AS ID, CAST(M.BrandID AS VARCHAR(36)) AS BrandID, B.BrandName, M.ModelName, M.VehicleType FROM dbo.MasterVehicleModels M JOIN dbo.MasterVehicleBrands B ON M.BrandID = B.ID WHERE M.ID = ?`
	res := r.db.Raw(query, id).Scan(&model)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &model, nil
}

func (r *modelRepository) Update(model *entity.Model) error {
	query := `UPDATE dbo.MasterVehicleModels SET BrandID = ?, ModelName = ?, VehicleType = ? WHERE ID = ?`
	return r.db.Exec(query, model.BrandID, model.ModelName, model.VehicleType, model.ID).Error
}

func (r *modelRepository) Delete(id string) error {
	query := `DELETE FROM dbo.MasterVehicleModels WHERE ID = ?`
	return r.db.Exec(query, id).Error
}
