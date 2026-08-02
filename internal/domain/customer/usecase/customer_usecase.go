package usecase

import (
	"errors"

	"time"

	"bj/internal/domain/customer/dto"
	"bj/internal/domain/customer/entity"
	"bj/internal/domain/customer/repository"

	"github.com/google/uuid"
)

type CustomerUsecase interface {
	CreateCustomer(req dto.CreateCustomerRequest) (string, error)
	GetAllCustomers() ([]dto.CustomerResponse, error)
	GetCustomerByID(id string) (*dto.CustomerResponse, error)
	UpdateCustomer(id string, req dto.UpdateCustomerRequest) error
	DeleteCustomer(id string) error
}

type customerUsecase struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerUsecase(repo repository.CustomerRepository) CustomerUsecase {
	return &customerUsecase{repo}
}

// Mapper functions
func mapCustomerToResponse(c entity.Customer) dto.CustomerResponse {
	return dto.CustomerResponse{
		ID:          c.ID,
		Fullname:    c.Fullname,
		Email:       c.Email,
		PhoneNumber: c.PhoneNumber,
		Address:     c.Address,
		NIK:         c.NIK,
		CreatedAt:   c.CreatedAt,
	}
}

func (u *customerUsecase) CreateCustomer(req dto.CreateCustomerRequest) (string, error) {
	newID := uuid.New().String()
	customer := &entity.Customer{
		ID:          newID,
		Fullname:    req.Fullname,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		NIK:         req.NIK,
		CreatedAt:   time.Now(),
	}
	err := u.customerRepo.Create(customer)
	return newID, err
}

func (u *customerUsecase) GetAllCustomers() ([]dto.CustomerResponse, error) {
	customers, err := u.customerRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.CustomerResponse
	for _, c := range customers {
		responses = append(responses, mapCustomerToResponse(c))
	}
	return responses, nil
}

func (u *customerUsecase) GetCustomerByID(id string) (*dto.CustomerResponse, error) {
	customer, err := u.customerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.New("pelanggan tidak ditemukan")
	}

	resp := mapCustomerToResponse(*customer)
	return &resp, nil
}

func (u *customerUsecase) UpdateCustomer(id string, req dto.UpdateCustomerRequest) error {
	customer, err := u.customerRepo.FindByID(id)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("pelanggan tidak ditemukan")
	}

	if req.Fullname != "" {
		customer.Fullname = req.Fullname
	}
	if req.PhoneNumber != "" {
		customer.PhoneNumber = req.PhoneNumber
	}
	if req.Address != "" {
		customer.Address = req.Address
	}
	if req.NIK != "" {
		customer.NIK = req.NIK
	}
	if req.Email != "" {
		customer.Email = req.Email
	}

	return u.customerRepo.Update(customer)
}

func (u *customerUsecase) DeleteCustomer(id string) error {
	return u.customerRepo.Delete(id)
}
