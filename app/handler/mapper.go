package handler

import (
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CatalogResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
}

type ProductResponse struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type ProductDetailResponse struct {
	Code     string            `json:"code"`
	Price    float64           `json:"price"`
	Category string            `json:"category"`
	Variants []VariantResponse `json:"variants"`
}

type VariantResponse struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

func toProductResponse(p models.Product) ProductResponse {
	return ProductResponse{
		Code:     p.Code,
		Price:    p.Price.InexactFloat64(),
		Category: p.Category.Code,
	}
}

func toProductDetailResponse(p models.Product) ProductDetailResponse {
	variants := make([]VariantResponse, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = toVariantResponse(v)
	}

	return ProductDetailResponse{
		Code:     p.Code,
		Price:    p.Price.InexactFloat64(),
		Category: p.Category.Code,
		Variants: variants,
	}
}

func toVariantResponse(v models.Variant) VariantResponse {
	return VariantResponse{
		Name:  v.Name,
		SKU:   v.SKU,
		Price: v.Price.InexactFloat64(),
	}
}
