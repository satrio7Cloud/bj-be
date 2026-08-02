package dto

type CreateBrandRequest struct {
	BrandName string `json:"brand_name" validate:"required"`
}

type UpdateBrandRequest struct {
	BrandName string `json:"brand_name" validate:"required"`
}

type BrandResponse struct {
	ID        string `json:"id"`
	BrandName string `json:"brand_name"`
}
