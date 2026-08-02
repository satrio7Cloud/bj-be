package dto

type CreateDocumentRequest struct {
	OrderID      string `form:"order_id" validate:"required,uuid"`
	DocumentType string `form:"document_type" validate:"required"`
}

type VerifyDocumentRequest struct {
	IsVerified bool `json:"is_verified"`
}

type DocumentResponse struct {
	ID           string `json:"id"`
	OrderID      string `json:"order_id"`
	DocumentType string `json:"document_type"`
	FilePath     string `json:"file_path"`
	IsVerified   bool   `json:"is_verified"`
}
