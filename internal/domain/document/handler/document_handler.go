package handler

import (
	"bj/internal/constants"
	"bj/internal/domain/document/dto"
	"bj/internal/domain/document/usecase"
	"bj/internal/utils"
	"bytes"
	"encoding/json"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type DocumentHandler struct {
	documentUsecase usecase.DocumentUsecase
}

func NewDocumentHandler(uc usecase.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{
		documentUsecase: uc,
	}
}

// @Summary Upload Document (Admin)
// @Description Upload a new document for an order from admin dashboard
// @Tags documents
// @Accept multipart/form-data
// @Produce json
// @Param order_id formData string true "Order ID"
// @Param document_type formData string true "Document Type (e.g. KTP, STNK)"
// @Param document_file formData file true "File to upload (Max 5MB, JPG/PNG/PDF)"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/documents/admin [post]
func (h *DocumentHandler) CreateDocumentAdmin(c *fiber.Ctx) error {
	var req dto.CreateDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest})
	}

	file, err := c.FormFile("document_file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": "File dokumen wajib diunggah (key: document_file)"})
	}

	// Validasi ukuran max 5MB
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": "Ukuran file maksimal 5MB"})
	}

	// Validasi ekstensi file
	ext := filepath.Ext(file.Filename)
	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".pdf": true}
	if !validExts[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": "Ekstensi file tidak didukung. Hanya menerima JPG, JPEG, PNG, atau PDF"})
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest, "errors": validationErrors})
	}

	ctx := c.Context() // Mendapatkan fasthttp context
	if err := h.documentUsecase.CreateDocument(ctx, req, file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": constants.StatusSuccess, "message": "Dokumen berhasil diunggah oleh Admin"})
}

// @Summary Upload Document (Customer)
// @Description Upload a new document for an order from public web
// @Tags documents
// @Accept multipart/form-data
// @Produce json
// @Param order_id formData string true "Order ID"
// @Param document_type formData string true "Document Type (e.g. KTP, STNK)"
// @Param document_file formData file true "File to upload (Max 5MB, JPG/PNG/PDF)"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/documents/customer [post]
func (h *DocumentHandler) CreateDocumentCustomer(c *fiber.Ctx) error {
	var req dto.CreateDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest})
	}

	file, err := c.FormFile("document_file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": "File dokumen wajib diunggah (key: document_file)"})
	}

	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": "Ukuran file maksimal 5MB"})
	}

	ext := filepath.Ext(file.Filename)
	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".pdf": true}
	if !validExts[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": "Ekstensi file tidak didukung. Hanya menerima JPG, JPEG, PNG, atau PDF"})
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": constants.MsgInvalidRequest, "errors": validationErrors})
	}

	ctx := c.Context()
	if err := h.documentUsecase.CreateDocument(ctx, req, file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": constants.StatusSuccess, "message": "Dokumen berhasil diunggah. Menunggu verifikasi Admin."})
}

// @Summary Get Document Pre-signed URL
// @Description Get a temporary URL to download/view the document
// @Tags documents
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/documents/{id}/url [get]
func (h *DocumentHandler) GetDocumentURL(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.Context()

	url, err := h.documentUsecase.GetDocumentURL(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "data": fiber.Map{"url": url}})
}

// @Summary Get Order Documents
// @Description Retrieve a list of documents attached to a specific order
// @Tags documents
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/documents/order/{orderId} [get]
func (h *DocumentHandler) GetDocumentsByOrder(c *fiber.Ctx) error {
	orderID := c.Params("orderId")
	docs, err := h.documentUsecase.GetDocumentsByOrder(orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "data": docs})
}

// @Summary Verify Document
// @Description Update verification status of a document
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Param request body dto.VerifyDocumentRequest true "Verification Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/documents/{id}/verify [patch]
func (h *DocumentHandler) VerifyDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.VerifyDocumentRequest

	// Gunakan custom JSON decoder agar menolak unknown fields (strict)
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusError, "message": "Payload tidak valid atau mengandung field yang tidak dikenali"})
	}

	if err := h.documentUsecase.VerifyDocument(id, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": "Status verifikasi dokumen berhasil diperbarui"})
}

// @Summary Delete Document
// @Description Delete a document and its file from storage
// @Tags documents
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/documents/{id} [delete]
func (h *DocumentHandler) DeleteDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.Context()

	if err := h.documentUsecase.DeleteDocument(ctx, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusError, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": "Dokumen berhasil dihapus"})
}
