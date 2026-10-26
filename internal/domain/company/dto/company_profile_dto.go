package dto

import "time"

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

type CreateCompanyProfileRequest struct {
	OfficeName       string `json:"office_name" validate:"required"`
	FullAddress      string `json:"full_address" validate:"required"`
	PostCode         string `json:"post_code" validate:"required,len=5"`
	PhoneNumber      string `json:"phone_number" validate:"required"`
	OperationalHours string `json:"operational_hours" validate:"required"`
}

type UpdateCompanyProfileRequest struct {
	OfficeName       string `json:"office_name" validate:"required"`
	FullAddress      string `json:"full_address" validate:"required"`
	PostCode         string `json:"post_code" validate:"required,len=5"`
	PhoneNumber      string `json:"phone_number" validate:"required"`
	OperationalHours string `json:"operational_hours" validate:"required"`
}
