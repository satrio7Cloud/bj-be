package repository

import (
	"bj/internal/domain/shipping/entity"
	"gorm.io/gorm"
)

type ShippingSettingsRepository interface {
	Get() (*entity.ShippingSettings, error)
	Upsert(settings *entity.ShippingSettings) error
}

type shippingSettingsRepository struct {
	db *gorm.DB
}

func NewShippingSettingsRepository(db *gorm.DB) ShippingSettingsRepository {
	return &shippingSettingsRepository{db}
}

func (r *shippingSettingsRepository) Get() (*entity.ShippingSettings, error) {
	var s entity.ShippingSettings
	res := r.db.Raw(`SELECT 
		CAST(ID AS VARCHAR(36)) AS ID,
		EnabledCouriers,
		DefaultItemWeight, DefaultItemValue,
		CAST(UpdatedBy AS VARCHAR(36)) AS UpdatedBy,
		CreatedAt, UpdatedAt
	FROM dbo.ShippingSettings WHERE ID = 'default'`).Scan(&s)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &s, nil
}

func (r *shippingSettingsRepository) Upsert(s *entity.ShippingSettings) error {
	var count int64
	r.db.Raw(`SELECT COUNT(1) FROM dbo.ShippingSettings WHERE ID = ?`, s.ID).Scan(&count)

	if count > 0 {
		query := `UPDATE dbo.ShippingSettings SET
			EnabledCouriers = ?,
			DefaultItemWeight = ?,
			DefaultItemValue = ?,
			UpdatedBy = ?,
			UpdatedAt = GETDATE()
		WHERE ID = ?`
		return r.db.Exec(query,
			s.EnabledCouriers,
			s.DefaultItemWeight,
			s.DefaultItemValue,
			s.UpdatedBy,
			s.ID,
		).Error
	}

	query := `INSERT INTO dbo.ShippingSettings (
		ID, EnabledCouriers, DefaultItemWeight, DefaultItemValue, UpdatedBy, CreatedAt, UpdatedAt
	) VALUES (?, ?, ?, ?, ?, GETDATE(), GETDATE())`
	return r.db.Exec(query,
		s.ID,
		s.EnabledCouriers,
		s.DefaultItemWeight,
		s.DefaultItemValue,
		s.UpdatedBy,
	).Error
}
