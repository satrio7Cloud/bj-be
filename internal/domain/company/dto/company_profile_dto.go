package dto

import (
	"time"

	"bj/internal/domain/company/entity"
)

type CompanyProfileResponse struct {
	ID               string    `json:"id"`
	OfficeName       string    `json:"office_name"`
	FullAddress      string    `json:"full_address"`
	PostCode         string    `json:"post_code"`
	PhoneNumber      string    `json:"phone_number"`
	OperationalHours string    `json:"operational_hours"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func ToCompanyProfileResponse(profile entity.CompanyProfile) CompanyProfileResponse {
	return CompanyProfileResponse{
		ID:               profile.ID,
		OfficeName:       profile.OfficeName,
		FullAddress:      profile.FullAddress,
		PostCode:         profile.PostCode,
		PhoneNumber:      profile.PhoneNumber,
		OperationalHours: profile.OperationalHours,
		CreatedAt:        profile.CreatedAt,
		UpdatedAt:        profile.UpdatedAt,
	}
}

type CreateCompanyProfileRequest struct {
	OfficeName       string `json:"office_name" validate:"required"`
	FullAddress      string `json:"full_address" validate:"required"`
	PostCode         string `json:"post_code" validate:"required,len=5"`
	PhoneNumber      string `json:"phone_number" validate:"required"`
	OperationalHours string `json:"operational_hours" validate:"required"`
}

func (req *CreateCompanyProfileRequest) ToEntity(id string) *entity.CompanyProfile {
	now := time.Now()
	return &entity.CompanyProfile{
		ID:               id,
		OfficeName:       req.OfficeName,
		FullAddress:      req.FullAddress,
		PostCode:         req.PostCode,
		PhoneNumber:      req.PhoneNumber,
		OperationalHours: req.OperationalHours,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

type UpdateCompanyProfileRequest struct {
	OfficeName       string `json:"office_name" validate:"required"`
	FullAddress      string `json:"full_address" validate:"required"`
	PostCode         string `json:"post_code" validate:"required,len=5"`
	PhoneNumber      string `json:"phone_number" validate:"required"`
	OperationalHours string `json:"operational_hours" validate:"required"`
}

func (req *UpdateCompanyProfileRequest) ApplyToEntity(profile *entity.CompanyProfile) {
	profile.OfficeName = req.OfficeName
	profile.FullAddress = req.FullAddress
	profile.PostCode = req.PostCode
	profile.PhoneNumber = req.PhoneNumber
	profile.OperationalHours = req.OperationalHours
	profile.UpdatedAt = time.Now()
}
