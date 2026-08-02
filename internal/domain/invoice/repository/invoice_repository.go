package repository

import (
	"bj/internal/domain/invoice/dto"
	"bj/internal/domain/invoice/entity"
)

type InvoiceRepository interface {
	Create(invoice *entity.Invoice) error
	GetByOrderID(orderID string) (*entity.Invoice, error)
	GetByID(id string) (*entity.Invoice, error)
	GetAll(page, limit int) ([]entity.InvoiceWithDetails, int64, error)
	Update(invoice *entity.Invoice) error
	GetActiveQrisAmounts(baseAmount float64) ([]float64, error)
	FindAll(page int, limit int) ([]dto.InvoiceDetailResponse, int64, error)
}

