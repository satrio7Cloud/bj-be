package dto

import "time"

type CreateCustomerRequest struct {
	Fullname    string `json:"fullname" validate:"required"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Address     string `json:"address"`
	NIK         string `json:"nik"`
}

type UpdateCustomerRequest struct {
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	NIK         string `json:"nik"`
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
