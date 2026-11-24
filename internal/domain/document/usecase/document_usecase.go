package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"
	"time"

	"bj/internal/domain/document/dto"
	"bj/internal/domain/document/entity"
	"bj/internal/domain/document/repository"
	orderRepo "bj/internal/domain/order/repository"
	serviceRepo "bj/internal/domain/service/repository"
	"bj/internal/shared/storage"

	"github.com/google/uuid"
)

type DocumentUsecase interface {
	CreateDocument(ctx context.Context, req dto.CreateDocumentRequest, file *multipart.FileHeader) error
	GetDocumentURL(ctx context.Context, id string) (string, error)
	GetDocumentsByOrder(orderID string) ([]dto.DocumentResponse, error)
	VerifyDocument(id string, req dto.VerifyDocumentRequest) error
	DeleteDocument(ctx context.Context, id string) error
}

type documentUsecase struct {
	documentRepo   repository.DocumentRepository
	orderRepo      orderRepo.OrderRepository
	serviceRepo    serviceRepo.ServiceRepository
	storageService storage.StorageService
}

func NewDocumentUsecase(repo repository.DocumentRepository, orderRepo orderRepo.OrderRepository, serviceRepo serviceRepo.ServiceRepository, storageService storage.StorageService) DocumentUsecase {
	return &documentUsecase{repo, orderRepo, serviceRepo, storageService}
}

func (u *documentUsecase) CreateDocument(ctx context.Context, req dto.CreateDocumentRequest, file *multipart.FileHeader) error {
	folderName := "documents"
	order, err := u.orderRepo.FindByID(req.OrderID)
	if err == nil && order != nil {
		service, errSvc := u.serviceRepo.FindByID(order.ServiceID)
		if errSvc == nil && service != nil {
			folderName = strings.ToLower(service.ServiceName)
			folderName = strings.ReplaceAll(folderName, " ", "-")
		}
	} else {
		orderKTL, errKTL := u.orderRepo.FindOrderKTLByID(req.OrderID)
		if errKTL == nil && orderKTL != nil {
			folderName = "cetak-ktl"
		} else {
			return errors.New("pesanan tidak valid atau tidak ditemukan")
		}
	}

	objectName, err := u.storageService.UploadFile(ctx, file, folderName)
	if err != nil {
		return err
	}

	document := &entity.Document{
		ID:           uuid.New().String(),
		OrderID:      req.OrderID,
		DocumentType: req.DocumentType,
		FilePath:     objectName,
		IsVerified:   false,
	}
	return u.documentRepo.Create(document)
}

func (u *documentUsecase) GetDocumentURL(ctx context.Context, id string) (string, error) {
	document, err := u.documentRepo.FindByID(id)
	if err != nil {
		return "", err
	}
	if document == nil {
		return "", errors.New("dokumen tidak ditemukan")
	}

	return u.storageService.GetPresignedURL(ctx, document.FilePath, 15*time.Minute)
}

func (u *documentUsecase) GetDocumentsByOrder(orderID string) ([]dto.DocumentResponse, error) {
	documents, err := u.documentRepo.FindByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	return dto.ToDocumentResponses(documents), nil
}

func (u *documentUsecase) VerifyDocument(id string, req dto.VerifyDocumentRequest) error {
	document, err := u.documentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if document == nil {
		return errors.New("dokumen tidak ditemukan")
	}

	document.IsVerified = req.IsVerified
	return u.documentRepo.Update(document)
}

func (u *documentUsecase) DeleteDocument(ctx context.Context, id string) error {
	document, err := u.documentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if document != nil {
		u.storageService.DeleteFile(ctx, document.FilePath)
	}

	return u.documentRepo.Delete(id)
}
