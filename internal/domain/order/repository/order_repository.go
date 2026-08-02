package repository

import (
	"bj/internal/domain/order/entity"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *entity.Order) error
	FindAll(page int, limit int) ([]entity.Order, error)
	FindByID(id string) (*entity.Order, error)
	UpdateStatus(id string, status string) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Create(order *entity.Order) error {
	query := `INSERT INTO dbo.Orders (
		ID, CustomerID, VehicleID, ServiceID, ProcessedByID, Status, 
		PickupMethod, ReturnMethod, CustomerTrackingNumber, BiroJasaTrackingNumber, PickupAddress,
		TaxAmount, ServiceFee, AdminFee, PhysicalCheckFee, DeliveryFee, ExpressFee, TotalCost, 
		ServiceLevel, IsNameTransferRequired, SamsatOrigin, SamsatDestination, Notes, OrderDate
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	isNameTransferBit := 0
	if order.IsNameTransferRequired {
		isNameTransferBit = 1
	}

	return r.db.Exec(query, 
		order.ID, order.CustomerID, order.VehicleID, order.ServiceID, order.ProcessedByID, order.Status,
		order.PickupMethod, order.ReturnMethod, order.CustomerTrackingNumber, order.BiroJasaTrackingNumber, order.PickupAddress,
		order.TaxAmount, order.ServiceFee, order.AdminFee, order.PhysicalCheckFee, order.DeliveryFee, order.ExpressFee, order.TotalCost,
		order.ServiceLevel, isNameTransferBit, order.SamsatOrigin, order.SamsatDestination, order.Notes, order.OrderDate,
	).Error
}

func (r *orderRepository) FindAll(page int, limit int) ([]entity.Order, error) {
	var orders []entity.Order
	
	offset := (page - 1) * limit
	
	query := `SELECT 
		CAST(ID AS VARCHAR(36)) AS ID, CAST(CustomerID AS VARCHAR(36)) AS CustomerID, 
		CAST(VehicleID AS VARCHAR(36)) AS VehicleID, CAST(ServiceID AS VARCHAR(36)) AS ServiceID, 
		CAST(ProcessedByID AS VARCHAR(36)) AS ProcessedByID, Status, 
		PickupMethod, ReturnMethod, CustomerTrackingNumber, BiroJasaTrackingNumber, PickupAddress,
		TaxAmount, ServiceFee, AdminFee, PhysicalCheckFee, DeliveryFee, ExpressFee, TotalCost, 
		ServiceLevel, IsNameTransferRequired, SamsatOrigin, SamsatDestination, Notes, OrderDate
		FROM dbo.Orders 
		ORDER BY OrderDate DESC
		OFFSET ? ROWS
		FETCH NEXT ? ROWS ONLY`
		
	err := r.db.Raw(query, offset, limit).Scan(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByID(id string) (*entity.Order, error) {
	var order entity.Order
	query := `SELECT 
		CAST(ID AS VARCHAR(36)) AS ID, CAST(CustomerID AS VARCHAR(36)) AS CustomerID, 
		CAST(VehicleID AS VARCHAR(36)) AS VehicleID, CAST(ServiceID AS VARCHAR(36)) AS ServiceID, 
		CAST(ProcessedByID AS VARCHAR(36)) AS ProcessedByID, Status, 
		PickupMethod, ReturnMethod, CustomerTrackingNumber, BiroJasaTrackingNumber, PickupAddress,
		TaxAmount, ServiceFee, AdminFee, PhysicalCheckFee, DeliveryFee, ExpressFee, TotalCost, 
		ServiceLevel, IsNameTransferRequired, SamsatOrigin, SamsatDestination, Notes, OrderDate 
	FROM dbo.Orders WHERE ID = ?`
	res := r.db.Raw(query, id).Scan(&order)
	
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	
	return &order, nil
}

func (r *orderRepository) UpdateStatus(id string, status string) error {
	query := `UPDATE dbo.Orders SET Status = ? WHERE ID = ?`
	return r.db.Exec(query, status, id).Error
}
