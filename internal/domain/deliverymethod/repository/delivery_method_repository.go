package repository

import (
	"bj/internal/domain/deliverymethod/entity"
	"gorm.io/gorm"
)

type DeliveryMethodRepository interface {
	Create(m *entity.DeliveryMethod) error
	FindAll() ([]entity.DeliveryMethod, error)
	FindActive(direction string) ([]entity.DeliveryMethod, error)
	FindByID(id string) (*entity.DeliveryMethod, error)
	FindByCode(code string) (*entity.DeliveryMethod, error)
	Update(m *entity.DeliveryMethod) error
	ToggleActive(id string) error
	Delete(id string) error
}

type deliveryMethodRepository struct {
	db *gorm.DB
}

func NewDeliveryMethodRepository(db *gorm.DB) DeliveryMethodRepository {
	return &deliveryMethodRepository{db}
}

func (r *deliveryMethodRepository) Create(m *entity.DeliveryMethod) error {
	query := `INSERT INTO dbo.MasterShippingMethods 
		(ID, MethodCode, MethodName, Direction, ProviderType, Description, FlatFee, RequiresAddress, RequiresCourier, IsActive, DisplayOrder, CreatedAt, UpdatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, GETDATE(), GETDATE())`

	return r.db.Exec(query,
		m.ID, m.MethodCode, m.MethodName, m.Direction, m.ProviderType,
		m.Description, m.FlatFee, m.RequiresAddress, m.RequiresCourier,
		m.IsActive, m.DisplayOrder,
	).Error
}

func (r *deliveryMethodRepository) FindAll() ([]entity.DeliveryMethod, error) {
	var list []entity.DeliveryMethod
	query := `SELECT 
		CAST(ID AS VARCHAR(36)) AS ID,
		MethodCode, MethodName, Direction, ProviderType, Description,
		FlatFee,
		CAST(RequiresAddress AS BIT) AS RequiresAddress,
		CAST(RequiresCourier AS BIT) AS RequiresCourier,
		CAST(IsActive AS BIT) AS IsActive,
		DisplayOrder, CreatedAt, UpdatedAt
	FROM dbo.MasterShippingMethods
	ORDER BY Direction ASC, DisplayOrder ASC`

	err := r.db.Raw(query).Scan(&list).Error
	return list, err
}

func (r *deliveryMethodRepository) FindActive(direction string) ([]entity.DeliveryMethod, error) {
	var list []entity.DeliveryMethod
	query := `SELECT 
		CAST(ID AS VARCHAR(36)) AS ID,
		MethodCode, MethodName, Direction, ProviderType, Description,
		FlatFee,
		CAST(RequiresAddress AS BIT) AS RequiresAddress,
		CAST(RequiresCourier AS BIT) AS RequiresCourier,
		CAST(IsActive AS BIT) AS IsActive,
		DisplayOrder, CreatedAt, UpdatedAt
	FROM dbo.MasterShippingMethods
	WHERE IsActive = 1 `

	var args []interface{}
	if direction != "" {
		query += `AND (Direction = ? OR Direction = 'BOTH') `
		args = append(args, direction)
	}
	query += `ORDER BY Direction ASC, DisplayOrder ASC`

	err := r.db.Raw(query, args...).Scan(&list).Error
	return list, err
}

func (r *deliveryMethodRepository) FindByID(id string) (*entity.DeliveryMethod, error) {
	var m entity.DeliveryMethod
	query := `SELECT 
		CAST(ID AS VARCHAR(36)) AS ID,
		MethodCode, MethodName, Direction, ProviderType, Description,
		FlatFee,
		CAST(RequiresAddress AS BIT) AS RequiresAddress,
		CAST(RequiresCourier AS BIT) AS RequiresCourier,
		CAST(IsActive AS BIT) AS IsActive,
		DisplayOrder, CreatedAt, UpdatedAt
	FROM dbo.MasterShippingMethods
	WHERE ID = ?`

	res := r.db.Raw(query, id).Scan(&m)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &m, nil
}

func (r *deliveryMethodRepository) FindByCode(code string) (*entity.DeliveryMethod, error) {
	var m entity.DeliveryMethod
	query := `SELECT 
		CAST(ID AS VARCHAR(36)) AS ID,
		MethodCode, MethodName, Direction, ProviderType, Description,
		FlatFee,
		CAST(RequiresAddress AS BIT) AS RequiresAddress,
		CAST(RequiresCourier AS BIT) AS RequiresCourier,
		CAST(IsActive AS BIT) AS IsActive,
		DisplayOrder, CreatedAt, UpdatedAt
	FROM dbo.MasterShippingMethods
	WHERE MethodCode = ?`

	res := r.db.Raw(query, code).Scan(&m)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &m, nil
}

func (r *deliveryMethodRepository) Update(m *entity.DeliveryMethod) error {
	query := `UPDATE dbo.MasterShippingMethods SET
		MethodName = ?,
		Direction = ?,
		ProviderType = ?,
		Description = ?,
		FlatFee = ?,
		RequiresAddress = ?,
		RequiresCourier = ?,
		IsActive = ?,
		DisplayOrder = ?,
		UpdatedAt = GETDATE()
	WHERE ID = ?`

	return r.db.Exec(query,
		m.MethodName, m.Direction, m.ProviderType, m.Description,
		m.FlatFee, m.RequiresAddress, m.RequiresCourier,
		m.IsActive, m.DisplayOrder, m.ID,
	).Error
}

func (r *deliveryMethodRepository) ToggleActive(id string) error {
	query := `UPDATE dbo.MasterShippingMethods SET
		IsActive = CASE WHEN IsActive = 1 THEN 0 ELSE 1 END,
		UpdatedAt = GETDATE()
	WHERE ID = ?`

	return r.db.Exec(query, id).Error
}

func (r *deliveryMethodRepository) Delete(id string) error {
	query := `DELETE FROM dbo.MasterShippingMethods WHERE ID = ?`
	return r.db.Exec(query, id).Error
}

