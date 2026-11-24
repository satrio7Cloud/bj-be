package repository

import (
	"bj/internal/domain/courier/entity"
	"gorm.io/gorm"
)

type CourierRepository interface {
	UpsertAll(couriers []entity.Courier) error
	FindAll() ([]entity.Courier, error)
	FindByCodes(codes []string) ([]entity.Courier, error)
}

type courierRepository struct {
	db *gorm.DB
}

func NewCourierRepository(db *gorm.DB) CourierRepository {
	return &courierRepository{db}
}

func (r *courierRepository) UpsertAll(couriers []entity.Courier) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DELETE FROM dbo.ShippingCouriers`).Error; err != nil {
			return err
		}
		if len(couriers) == 0 {
			return nil
		}
		return tx.Exec(buildBulkInsert(couriers), buildArgs(couriers)...).Error
	})
}

func buildBulkInsert(couriers []entity.Courier) string {
	query := `INSERT INTO dbo.ShippingCouriers 
		(ID, CourierCode, CourierName, CourierServiceCode, CourierServiceName, Tier, ServiceType, ShippingType, DurationRange, DurationUnit, IsAvailable, CreatedAt, UpdatedAt)
		VALUES `
	for i := range couriers {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, GETDATE(), GETDATE())"
	}
	return query
}

func buildArgs(couriers []entity.Courier) []interface{} {
	args := make([]interface{}, 0, len(couriers)*11)
	for _, c := range couriers {
		args = append(args,
			c.ID, c.CourierCode, c.CourierName, c.CourierServiceCode, c.CourierServiceName,
			c.Tier, c.ServiceType, c.ShippingType, c.DurationRange, c.DurationUnit, c.IsAvailable,
		)
	}
	return args
}

func (r *courierRepository) FindAll() ([]entity.Courier, error) {
	var couriers []entity.Courier
	query := `SELECT 
		CAST(ID AS VARCHAR(36)) AS ID,
		CourierCode, CourierName, CourierServiceCode, CourierServiceName,
		Tier, ServiceType, ShippingType, DurationRange, DurationUnit,
		CAST(IsAvailable AS BIT) AS IsAvailable,
		CreatedAt, UpdatedAt
	FROM dbo.ShippingCouriers
	ORDER BY CourierName ASC, CourierServiceName ASC`
	err := r.db.Raw(query).Scan(&couriers).Error
	return couriers, err
}

func (r *courierRepository) FindByCodes(codes []string) ([]entity.Courier, error) {
	if len(codes) == 0 {
		return []entity.Courier{}, nil
	}
	var couriers []entity.Courier
	err := r.db.Raw(`SELECT 
		CAST(ID AS VARCHAR(36)) AS ID,
		CourierCode, CourierName, CourierServiceCode, CourierServiceName,
		Tier, ServiceType, ShippingType, DurationRange, DurationUnit,
		CAST(IsAvailable AS BIT) AS IsAvailable,
		CreatedAt, UpdatedAt
	FROM dbo.ShippingCouriers
	WHERE CourierCode IN ?
	ORDER BY CourierName ASC, CourierServiceName ASC`, codes).Scan(&couriers).Error
	return couriers, err
}

