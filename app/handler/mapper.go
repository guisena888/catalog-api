package handler

import (
	"github.com/mytheresa/go-hiring-challenge/pkg/model"
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

func toCatalogResponse(products []model.Product, total int64) CatalogResponse {
	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = toProductResponse(p)
	}
	return CatalogResponse{
		Products: responses,
		Total:    total,
	}
}

func toProductResponse(p model.Product) ProductResponse {
	return ProductResponse{
		Code:     p.Code,
		Price:    p.Price.InexactFloat64(),
		Category: p.Category.Code,
	}
}

func toProductDetailResponse(p model.Product) ProductDetailResponse {
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

func toVariantResponse(v model.Variant) VariantResponse {
	return VariantResponse{
		Name:  v.Name,
		SKU:   v.SKU,
		Price: v.Price.InexactFloat64(),
	}
}

type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func toCategoryResponse(c model.Category) CategoryResponse {
	return CategoryResponse{
		Code: c.Code,
		Name: c.Name,
	}
}
