package dto

import "bj/internal/domain/document/entity"

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

func ToDocumentResponse(d entity.Document) DocumentResponse {
	return DocumentResponse{
		ID:           d.ID,
		OrderID:      d.OrderID,
		DocumentType: d.DocumentType,
		FilePath:     d.FilePath,
		IsVerified:   d.IsVerified,
	}
}

func ToDocumentResponses(documents []entity.Document) []DocumentResponse {
	responses := make([]DocumentResponse, 0, len(documents))
	for _, d := range documents {
		responses = append(responses, ToDocumentResponse(d))
	}
	return responses
}
