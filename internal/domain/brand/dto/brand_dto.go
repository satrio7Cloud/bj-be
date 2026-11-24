package dto

import "bj/internal/domain/brand/entity"

type CreateBrandRequest struct {
	BrandName string `json:"brand_name" validate:"required"`
}

func (req *CreateBrandRequest) ToEntity(id string) *entity.Brand {
	return &entity.Brand{
		ID:        id,
		BrandName: req.BrandName,
	}
}

type UpdateBrandRequest struct {
	BrandName string `json:"brand_name" validate:"required"`
}

func (req *UpdateBrandRequest) ApplyToEntity(b *entity.Brand) {
	b.BrandName = req.BrandName
}

type BrandResponse struct {
	ID        string `json:"id"`
	BrandName string `json:"brand_name"`
}

func ToBrandResponse(b entity.Brand) BrandResponse {
	return BrandResponse{
		ID:        b.ID,
		BrandName: b.BrandName,
	}
}

func ToBrandResponses(brands []entity.Brand) []BrandResponse {
	responses := make([]BrandResponse, 0, len(brands))
	for _, b := range brands {
		responses = append(responses, ToBrandResponse(b))
	}
	return responses
}
