package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"time"

	"bj/internal/domain/document/dto"
	"bj/internal/domain/document/entity"
	"bj/internal/domain/document/repository"
	orderRepo "bj/internal/domain/order/repository"
	serviceRepo "bj/internal/domain/service/repository"
	"bj/internal/shared/storage"
	"strings"

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

// Mapper functions
func mapDocumentToResponse(d entity.Document) dto.DocumentResponse {
	return dto.DocumentResponse{
		ID:           d.ID,
		OrderID:      d.OrderID,
		DocumentType: d.DocumentType,
		IsVerified:   d.IsVerified,
	}
}

func (u *documentUsecase) CreateDocument(ctx context.Context, req dto.CreateDocumentRequest, file *multipart.FileHeader) error {
	// 0. Ambil informasi Order & Service untuk menentukan nama folder
	folderName := "documents" // default
	order, err := u.orderRepo.FindByID(req.OrderID)
	if err == nil && order != nil {
		service, errSvc := u.serviceRepo.FindByID(order.ServiceID)
		if errSvc == nil && service != nil {
			// Format folder: perpanjangan-motor-1-tahun
			folderName = strings.ToLower(service.ServiceName)
			folderName = strings.ReplaceAll(folderName, " ", "-")
			// Bersihkan karakter aneh jika perlu, tapi replace spasi minimal cukup
		}
	} else {
		// Jika tidak ditemukan di tabel Order reguler, cari di tabel OrderKTL
		orderKTL, errKTL := u.orderRepo.FindOrderKTLByID(req.OrderID)
		if errKTL == nil && orderKTL != nil {
			folderName = "cetak-ktl"
		} else {
			return errors.New("pesanan tidak valid atau tidak ditemukan")
		}
	}

	// 1. Upload ke MinIO
	objectName, err := u.storageService.UploadFile(ctx, file, folderName)
	if err != nil {
		return err
	}

	// 2. Simpan ke database
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

	// Presigned URL valid selama 15 menit
	return u.storageService.GetPresignedURL(ctx, document.FilePath, 15*time.Minute)
}

func (u *documentUsecase) GetDocumentsByOrder(orderID string) ([]dto.DocumentResponse, error) {
	documents, err := u.documentRepo.FindByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.DocumentResponse, 0)
	for _, d := range documents {
		responses = append(responses, mapDocumentToResponse(d))
	}
	return responses, nil
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
		// KARENA INI SOFT DELETE, kita tidak menghapus file fisiknya dari MinIO.
		// Jika suatu saat ingin direstore, file masih ada.
		// _ = u.storageService.DeleteFile(ctx, document.FilePath)
	}

	return u.documentRepo.Delete(id)
}
