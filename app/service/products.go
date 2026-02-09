package service

import (
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductService struct {
	repo ProductRepository
}

func NewProductService(r ProductRepository) *ProductService {
	return &ProductService{
		repo: r,
	}
}

func (s *ProductService) GetProducts(filter *models.ProductFilter) ([]models.Product, int64, error) {
	return s.repo.GetProducts(filter)
}

func (s *ProductService) GetProductDetails(code string) (*models.Product, error) {
	product, err := s.repo.GetProductDetails(code)
	if err != nil {
		return nil, err
	}

	normalizePrices(product)

	return product, nil
}

func normalizePrices(product *models.Product) {
	for i, v := range product.Variants {
		if !v.Price.IsPositive() {
			product.Variants[i].Price = product.Price
		}
	}
}
