package repository

import (
	"bj/internal/domain/customer/entity"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *entity.Customer) error
	FindAll() ([]entity.Customer, error)
	FindByID(id string) (*entity.Customer, error)
	FindByPhoneOrNIK(phone string, nik string) (*entity.Customer, error)
	Update(customer *entity.Customer) error
	Delete(id string) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}

func (r *customerRepository) Create(customer *entity.Customer) error {
	query := `
		INSERT INTO dbo.Customers (ID, Fullname, PhoneNumber, Address, NIK, Email, CreatedAt) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	return r.db.Exec(query, customer.ID, customer.Fullname, customer.PhoneNumber, customer.Address, customer.NIK, customer.Email, customer.CreatedAt).Error
}

func (r *customerRepository) FindAll() ([]entity.Customer, error) {
	var customers []entity.Customer
	query := `
		SELECT 
			CAST(ID AS VARCHAR(36)) AS ID, Fullname, PhoneNumber, Address, NIK, Email, CreatedAt 
		FROM dbo.Customers 
		ORDER BY CreatedAt DESC`
	err := r.db.Raw(query).Scan(&customers).Error
	return customers, err
}

func (r *customerRepository) FindByID(id string) (*entity.Customer, error) {
	var customer entity.Customer
	query := `
		SELECT 
			CAST(ID AS VARCHAR(36)) AS ID, Fullname, PhoneNumber, Address, NIK, Email, CreatedAt 
		FROM dbo.Customers 
		WHERE ID = ?
	`
	res := r.db.Raw(query, id).Scan(&customer)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &customer, nil
}

func (r *customerRepository) FindByPhoneOrNIK(
	phone string,
	nik string,
) (*entity.Customer, error) {
	if phone == "" && nik == "" {
		return nil, fmt.Errorf("phone atau NIK wajib diisi")
	}

	var customer entity.Customer

	query := `
        SELECT
            CAST(ID AS VARCHAR(36)) AS ID,
            Fullname,
            PhoneNumber,
            Address,
            NIK,
            Email,
            CreatedAt
        FROM dbo.Customers
        WHERE
    `

	var conditions []string
	var args []interface{}

	if phone != "" {
		conditions = append(conditions, "PhoneNumber = ?")
		args = append(args, phone)
	}

	if nik != "" {
		conditions = append(conditions, "NIK = ?")
		args = append(args, nik)
	}

	query += strings.Join(conditions, " OR ")

	res := r.db.Raw(query, args...).Scan(&customer)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &customer, nil
}

func (r *customerRepository) Update(customer *entity.Customer) error {
	query := `
		UPDATE dbo.Customers 
		SET Fullname = ?, PhoneNumber = ?, Address = ?, NIK = ?, Email = ? 
		WHERE ID = ?
	`
	return r.db.Exec(query, customer.Fullname, customer.PhoneNumber, customer.Address, customer.NIK, customer.Email, customer.ID).Error
}

func (r *customerRepository) Delete(id string) error {
	query := `DELETE FROM Customers WHERE ID = ?`
	return r.db.Exec(query, id).Error
}
