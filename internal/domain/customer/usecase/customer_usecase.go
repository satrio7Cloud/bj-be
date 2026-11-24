package usecase

import (
	"errors"

	"bj/internal/domain/customer/dto"
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

func (u *customerUsecase) CreateCustomer(req dto.CreateCustomerRequest) (string, error) {
	newID := uuid.New().String()
	customer := req.ToEntity(newID)
	err := u.customerRepo.Create(customer)
	return newID, err
}

func (u *customerUsecase) GetAllCustomers() ([]dto.CustomerResponse, error) {
	customers, err := u.customerRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.ToCustomerResponses(customers), nil
}

func (u *customerUsecase) GetCustomerByID(id string) (*dto.CustomerResponse, error) {
	customer, err := u.customerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.New("pelanggan tidak ditemukan")
	}

	resp := dto.ToCustomerResponse(*customer)
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

	req.ApplyToEntity(customer)
	return u.customerRepo.Update(customer)
}

func (u *customerUsecase) DeleteCustomer(id string) error {
	return u.customerRepo.Delete(id)
}
