package handler

import (
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductService interface {
	GetProducts(filter *models.ProductFilter) ([]models.Product, int64, error)
	GetProductDetails(code string) (*models.Product, error)
}
