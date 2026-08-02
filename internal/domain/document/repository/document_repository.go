package repository

import (
	"bj/internal/domain/document/entity"

	"gorm.io/gorm"
)

type DocumentRepository interface {
	Create(document *entity.Document) error
	FindAll() ([]entity.Document, error)
	FindByOrderID(orderID string) ([]entity.Document, error)
	FindByID(id string) (*entity.Document, error)
	Update(document *entity.Document) error
	Delete(id string) error
}

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db}
}

func (r *documentRepository) Create(document *entity.Document) error {
	query := `INSERT INTO dbo.OrderDocuments (ID, OrderID, DocumentType, FilePath, IsVerified) VALUES (?, ?, ?, ?, ?)`
	// Konversi boolean ke bit jika diperlukan, walau driver biasanya handle ini
	isVerifiedBit := 0
	if document.IsVerified {
		isVerifiedBit = 1
	}
	return r.db.Exec(query, document.ID, document.OrderID, document.DocumentType, document.FilePath, isVerifiedBit).Error
}

func (r *documentRepository) FindAll() ([]entity.Document, error) {
	var documents []entity.Document
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, CAST(OrderID AS VARCHAR(36)) AS OrderID, DocumentType, FilePath, IsVerified FROM dbo.OrderDocuments WHERE DeletedAt IS NULL`
	err := r.db.Raw(query).Scan(&documents).Error
	return documents, err
}

func (r *documentRepository) FindByOrderID(orderID string) ([]entity.Document, error) {
	var documents []entity.Document
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, CAST(OrderID AS VARCHAR(36)) AS OrderID, DocumentType, FilePath, IsVerified FROM dbo.OrderDocuments WHERE OrderID = ? AND DeletedAt IS NULL`
	err := r.db.Raw(query, orderID).Scan(&documents).Error
	return documents, err
}

func (r *documentRepository) FindByID(id string) (*entity.Document, error) {
	var document entity.Document
	query := `SELECT CAST(ID AS VARCHAR(36)) AS ID, CAST(OrderID AS VARCHAR(36)) AS OrderID, DocumentType, FilePath, IsVerified FROM dbo.OrderDocuments WHERE ID = ? AND DeletedAt IS NULL`
	res := r.db.Raw(query, id).Scan(&document)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &document, nil
}

func (r *documentRepository) Update(document *entity.Document) error {
	query := `UPDATE dbo.OrderDocuments SET IsVerified = ? WHERE ID = ?`
	isVerifiedBit := 0
	if document.IsVerified {
		isVerifiedBit = 1
	}
	return r.db.Exec(query, isVerifiedBit, document.ID).Error
}

func (r *documentRepository) Delete(id string) error {
	// Soft delete dengan menset DeletedAt ke waktu saat ini
	query := `UPDATE dbo.OrderDocuments SET DeletedAt = GETDATE() WHERE ID = ?`
	return r.db.Exec(query, id).Error
}
