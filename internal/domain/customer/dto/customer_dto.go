package dto

import (
	"time"

	"bj/internal/domain/customer/entity"
)

type CreateCustomerRequest struct {
	Fullname    string `json:"fullname" validate:"required"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Address     string `json:"address"`
	NIK         string `json:"nik"`
}

func (req *CreateCustomerRequest) ToEntity(id string) *entity.Customer {
	return &entity.Customer{
		ID:          id,
		Fullname:    req.Fullname,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		NIK:         req.NIK,
		CreatedAt:   time.Now(),
	}
}

type UpdateCustomerRequest struct {
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	NIK         string `json:"nik"`
}

func (req *UpdateCustomerRequest) ApplyToEntity(c *entity.Customer) {
	if req.Fullname != "" {
		c.Fullname = req.Fullname
	}
	if req.PhoneNumber != "" {
		c.PhoneNumber = req.PhoneNumber
	}
	if req.Address != "" {
		c.Address = req.Address
	}
	if req.NIK != "" {
		c.NIK = req.NIK
	}
	if req.Email != "" {
		c.Email = req.Email
	}
}

type CustomerResponse struct {
	ID          string    `json:"id"`
	Fullname    string    `json:"fullname"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	NIK         string    `json:"nik"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToCustomerResponse(c entity.Customer) CustomerResponse {
	return CustomerResponse{
		ID:          c.ID,
		Fullname:    c.Fullname,
		Email:       c.Email,
		PhoneNumber: c.PhoneNumber,
		Address:     c.Address,
		NIK:         c.NIK,
		CreatedAt:   c.CreatedAt,
	}
}

func ToCustomerResponses(customers []entity.Customer) []CustomerResponse {
	responses := make([]CustomerResponse, 0, len(customers))
	for _, c := range customers {
		responses = append(responses, ToCustomerResponse(c))
	}
	return responses
}
