package service

import (
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductRepository interface {
	GetProducts(filter *models.ProductFilter) ([]models.Product, int64, error)
	GetProductDetails(code string) (*models.Product, error)
}

