package repository

import (
	"bj/internal/domain/vehicle/entity"
	"gorm.io/gorm"
)

type VehicleRepository interface {
	Create(vehicle *entity.Vehicle) error
	FindAll() ([]entity.Vehicle, error)
	FindByCustomerID(customerID string) ([]entity.Vehicle, error)
	FindByID(id string) (*entity.Vehicle, error)
	FindByPlateNumber(plateNumber string) (*entity.Vehicle, error)
	Update(vehicle *entity.Vehicle) error
	Delete(id string) error
}

type vehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) VehicleRepository {
	return &vehicleRepository{db}
}

func (r *vehicleRepository) Create(vehicle *entity.Vehicle) error {
	query := `INSERT INTO dbo.Vehicles (ID, CustomerID, ModelID, PlateNumber, Year) VALUES (?, ?, ?, ?, ?)`
	return r.db.Exec(query, vehicle.ID, vehicle.CustomerID, vehicle.ModelID, vehicle.PlateNumber, vehicle.Year).Error
}

func (r *vehicleRepository) FindAll() ([]entity.Vehicle, error) {
	var vehicles []entity.Vehicle
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, CAST(CustomerID AS VARCHAR(36)) AS CustomerID, CAST(ModelID AS VARCHAR(36)) AS ModelID, PlateNumber, Year FROM dbo.Vehicles ORDER BY PlateNumber ASC`
	err := r.db.Raw(query).Scan(&vehicles).Error
	return vehicles, err
}

func (r *vehicleRepository) FindByCustomerID(customerID string) ([]entity.Vehicle, error) {
	var vehicles []entity.Vehicle
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, CAST(CustomerID AS VARCHAR(36)) AS CustomerID, CAST(ModelID AS VARCHAR(36)) AS ModelID, PlateNumber, Year FROM dbo.Vehicles WHERE CustomerID = ? ORDER BY PlateNumber ASC`
	err := r.db.Raw(query, customerID).Scan(&vehicles).Error
	return vehicles, err
}

func (r *vehicleRepository) FindByID(id string) (*entity.Vehicle, error) {
	var vehicle entity.Vehicle
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, CAST(CustomerID AS VARCHAR(36)) AS CustomerID, CAST(ModelID AS VARCHAR(36)) AS ModelID, PlateNumber, Year FROM dbo.Vehicles WHERE ID = ?`
	res := r.db.Raw(query, id).Scan(&vehicle)
	
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	
	return &vehicle, nil
}

func (r *vehicleRepository) FindByPlateNumber(plateNumber string) (*entity.Vehicle, error) {
	var vehicle entity.Vehicle
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, CAST(CustomerID AS VARCHAR(36)) AS CustomerID, CAST(ModelID AS VARCHAR(36)) AS ModelID, PlateNumber, Year FROM dbo.Vehicles WHERE PlateNumber = ?`
	res := r.db.Raw(query, plateNumber).Scan(&vehicle)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &vehicle, nil
}

func (r *vehicleRepository) Update(vehicle *entity.Vehicle) error {
	query := `UPDATE Vehicles SET ModelID = ?, PlateNumber = ?, Year = ? WHERE ID = ?`
	return r.db.Exec(query, vehicle.ModelID, vehicle.PlateNumber, vehicle.Year, vehicle.ID).Error
}

func (r *vehicleRepository) Delete(id string) error {
	query := `DELETE FROM Vehicles WHERE ID = ?`
	return r.db.Exec(query, id).Error
}
