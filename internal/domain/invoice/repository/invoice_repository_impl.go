package repository

import (
	"bj/internal/domain/invoice/dto"
	"bj/internal/domain/invoice/entity"
	"gorm.io/gorm"
)

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(invoice *entity.Invoice) error {
	query := `INSERT INTO dbo.invoices (ID, OrderID, InvoiceNumber, IssueDate, DueDate, PaymentStatus, PaymentMethod, PaidAt, QrisPayload, QrisAmount, PaymentToken, PaymentUrl)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, 
		invoice.ID, 
		invoice.OrderID, 
		invoice.InvoiceNumber, 
		invoice.IssueDate, 
		invoice.DueDate, 
		invoice.PaymentStatus, 
		invoice.PaymentMethod, 
		invoice.PaidAt,
		invoice.QrisPayload,
		invoice.QrisAmount,
		invoice.PaymentToken,
		invoice.PaymentUrl).Error
}

func (r *invoiceRepository) GetByOrderID(orderID string) (*entity.Invoice, error) {
	var invoice entity.Invoice
	query := `SELECT 
		CAST(ID AS VARCHAR(36)) AS ID,
		CAST(OrderID AS VARCHAR(36)) AS OrderID,
		InvoiceNumber, IssueDate, DueDate, PaymentStatus, PaymentMethod, PaidAt, QrisPayload, QrisAmount, PaymentToken, PaymentUrl
		FROM dbo.invoices WHERE OrderID = ?`
		
	res := r.db.Raw(query, orderID).Scan(&invoice)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil // Return nil if not found
	}
	return &invoice, nil
}

func (r *invoiceRepository) GetByID(id string) (*entity.Invoice, error) {
	var invoice entity.Invoice
	query := `SELECT 
		CAST(ID AS VARCHAR(36)) AS ID,
		CAST(OrderID AS VARCHAR(36)) AS OrderID,
		InvoiceNumber, IssueDate, DueDate, PaymentStatus, PaymentMethod, PaidAt, QrisPayload, QrisAmount, PaymentToken, PaymentUrl
		FROM dbo.invoices WHERE ID = ?`
		
	res := r.db.Raw(query, id).Scan(&invoice)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil // Return nil if not found
	}
	return &invoice, nil
}

func (r *invoiceRepository) Update(invoice *entity.Invoice) error {
	query := `UPDATE dbo.invoices 
			  SET PaymentStatus = ?, PaymentMethod = ?, PaidAt = ? 
			  WHERE ID = ?`
	return r.db.Exec(query, 
		invoice.PaymentStatus, 
		invoice.PaymentMethod, 
		invoice.PaidAt, 
		invoice.ID).Error
}

func (r *invoiceRepository) GetAll(page, limit int) ([]entity.InvoiceWithDetails, int64, error) {
	var invoices []entity.InvoiceWithDetails
	var total int64

	// Count total records
	countQuery := `SELECT COUNT(*) FROM dbo.invoices`
	if err := r.db.Raw(countQuery).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated records with joins
	query := `SELECT 
		CAST(i.ID AS VARCHAR(36)) AS ID,
		CAST(i.OrderID AS VARCHAR(36)) AS OrderID,
		i.InvoiceNumber, i.IssueDate, i.DueDate, i.PaymentStatus, i.PaymentMethod, i.PaidAt, i.QrisPayload, i.QrisAmount, i.PaymentToken, i.PaymentUrl,
		c.Fullname AS CustomerName, c.Email AS CustomerEmail, c.PhoneNumber AS CustomerPhone,
		v.PlateNumber AS VehiclePlate,
		s.ServiceName AS ServiceName,
		o.TotalCost AS TotalCost
		FROM dbo.invoices i
		LEFT JOIN dbo.Orders o ON i.OrderID = o.ID
		LEFT JOIN dbo.Customers c ON o.CustomerID = c.ID
		LEFT JOIN dbo.Vehicles v ON o.VehicleID = v.ID
		LEFT JOIN dbo.MasterServices s ON o.ServiceID = s.ID
		ORDER BY i.IssueDate DESC
		OFFSET ? ROWS FETCH NEXT ? ROWS ONLY`
	if err := r.db.Raw(query, offset, limit).Scan(&invoices).Error; err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}
func (r *invoiceRepository) GetActiveQrisAmounts(baseAmount float64) ([]float64, error) {
	var amounts []float64
	query := `SELECT i.QrisAmount 
			  FROM dbo.invoices i
			  JOIN dbo.Orders o ON i.OrderID = o.ID
			  WHERE o.TotalCost = ? AND i.PaymentStatus = 'UNPAID' AND i.QrisAmount IS NOT NULL`
	err := r.db.Raw(query, baseAmount).Scan(&amounts).Error
	return amounts, err
}

func (r *invoiceRepository) FindAll(page int, limit int) ([]dto.InvoiceDetailResponse, int64, error) {
	var invoices []dto.InvoiceDetailResponse
	var total int64

	countQuery := `SELECT COUNT(*) FROM dbo.invoices`
	if err := r.db.Raw(countQuery).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	query := `SELECT 
		CAST(i.ID AS VARCHAR(36)) AS invoice_id,
		CAST(i.OrderID AS VARCHAR(36)) AS order_id,
		i.InvoiceNumber AS invoice_number,
		i.IssueDate AS issue_date,
		i.DueDate AS due_date,
		i.PaymentStatus AS payment_status,
		i.PaymentMethod AS payment_method,
		i.PaidAt AS paid_at,
		i.QrisPayload AS qris_payload,
		i.QrisAmount AS qris_amount,
		i.PaymentToken AS payment_token,
		i.PaymentUrl AS payment_url,
		c.Fullname AS customer_name,
		'' AS customer_email,
		c.PhoneNumber AS customer_phone,
		v.PlateNumber AS vehicle_plate,
		s.ServiceName AS service_name,
		o.TaxAmount AS tax_amount,
		o.ServiceFee AS service_fee,
		o.AdminFee AS admin_fee,
		o.PhysicalCheckFee AS physical_check_fee,
		o.DeliveryFee AS delivery_fee,
		o.ExpressFee AS express_fee,
		o.TotalCost AS total_cost
	FROM dbo.invoices i
	LEFT JOIN dbo.Orders o ON i.OrderID = o.ID
	LEFT JOIN dbo.Customers c ON o.CustomerID = c.ID
	LEFT JOIN dbo.Vehicles v ON o.VehicleID = v.ID
	LEFT JOIN dbo.MasterServices s ON o.ServiceID = s.ID
	ORDER BY i.IssueDate DESC
	OFFSET ? ROWS
	FETCH NEXT ? ROWS ONLY`

	if err := r.db.Raw(query, offset, limit).Scan(&invoices).Error; err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}
