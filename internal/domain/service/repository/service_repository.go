package repository

import (
	"bj/internal/domain/service/entity"

	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(service *entity.Service) error
	FindAll() ([]entity.Service, error)
	FindByID(id string) (*entity.Service, error)
	Update(service *entity.Service) error
	Delete(id string) error
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db}
}

func (r *serviceRepository) Create(service *entity.Service) error {
	query := `INSERT INTO dbo.MasterServices (ID, ServiceName, Category, ServiceFee, ExpressFee, AdminFee, PhysicalCheckFee, Description, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, service.ID, service.ServiceName, service.Category, service.ServiceFee, service.ExpressFee, service.AdminFee, service.PhysicalCheckFee, service.Description, service.CreatedAt).Error
}

func (r *serviceRepository) FindAll() ([]entity.Service, error) {
	var services []entity.Service
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, ServiceName, Category, ServiceFee, ExpressFee, AdminFee, PhysicalCheckFee, Description, CreatedAt FROM dbo.MasterServices ORDER BY CreatedAt DESC`
	err := r.db.Raw(query).Scan(&services).Error
	return services, err
}

func (r *serviceRepository) FindByID(id string) (*entity.Service, error) {
	var service entity.Service
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, ServiceName, Category, ServiceFee, ExpressFee, AdminFee, PhysicalCheckFee, Description, CreatedAt FROM dbo.MasterServices WHERE ID = ?`
	res := r.db.Raw(query, id).Scan(&service)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &service, nil
}

func (r *serviceRepository) Update(service *entity.Service) error {
	query := `UPDATE dbo.MasterServices SET ServiceName = ?, Category = ?, ServiceFee = ?, ExpressFee = ?, AdminFee = ?, PhysicalCheckFee = ?, Description = ? WHERE ID = ?`
	return r.db.Exec(query, service.ServiceName, service.Category, service.ServiceFee, service.ExpressFee, service.AdminFee, service.PhysicalCheckFee, service.Description, service.ID).Error
}

func (r *serviceRepository) Delete(id string) error {
	query := `DELETE FROM dbo.MasterServices WHERE ID = ?`
	return r.db.Exec(query, id).Error
}
