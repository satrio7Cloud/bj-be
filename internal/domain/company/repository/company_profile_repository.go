package repository

import (
	"bj/internal/config"
	"bj/internal/domain/company/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyProfileRepository interface {
	Get() (*entity.CompanyProfile, error)
	Update(profile *entity.CompanyProfile) error
	Create(profile *entity.CompanyProfile) error
}

type companyProfileRepository struct {
	db        *gorm.DB
	officeCfg config.OfficeConfig
}

func NewCompanyProfileRepository(db *gorm.DB, officeCfg config.OfficeConfig) CompanyProfileRepository {
	return &companyProfileRepository{
		db:        db,
		officeCfg: officeCfg,
	}
}

func (r *companyProfileRepository) Get() (*entity.CompanyProfile, error) {
	var profile entity.CompanyProfile
	query := `
		SELECT TOP 1 CAST(ID AS VARCHAR(36)) AS ID, OfficeName, FullAddress,
        PostCode, PhoneNumber, OperationalHours, CreatedAt, UpdatedAt
        FROM dbo.MasterOffice ORDER BY CreatedAt ASC
	`
	res := r.db.Raw(query).Scan(&profile)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 || profile.ID == "" {
		now := time.Now()
		defaultProfile := &entity.CompanyProfile{
			ID:               uuid.New().String(),
			OfficeName:       r.officeCfg.Name,
			FullAddress:      r.officeCfg.Address,
			PostCode:         r.officeCfg.PostalCode,
			PhoneNumber:      r.officeCfg.Phone,
			OperationalHours: r.officeCfg.OperationalHours,
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		if err := r.Create(defaultProfile); err != nil {
			return nil, err
		}
		return defaultProfile, nil
	}

	return &profile, nil
}

func (r *companyProfileRepository) Create(profile *entity.CompanyProfile) error {
	query := `
		INSERT INTO dbo.MasterOffice
        (ID, OfficeName, FullAddress, PostCode,PhoneNumber, 
		OperationalHours, CreatedAt, UpdatedAt)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	return r.db.Exec(query,
		profile.ID,
		profile.OfficeName,
		profile.FullAddress,
		profile.PostCode,
		profile.PhoneNumber,
		profile.OperationalHours,
		profile.CreatedAt,
		profile.UpdatedAt,
	).Error
}

func (r *companyProfileRepository) Update(profile *entity.CompanyProfile) error {
	query := `
		UPDATE dbo.MasterOffice
        SET OfficeName = ?, FullAddress = ?, PostCode = ?,PhoneNumber = ?, OperationalHours = ?, UpdatedAt = ?
        WHERE ID = ?
	`
	return r.db.Exec(query,
		profile.OfficeName,
		profile.FullAddress,
		profile.PostCode,
		profile.PhoneNumber,
		profile.OperationalHours,
		time.Now(),
		profile.ID,
	).Error
}
