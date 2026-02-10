package handler

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductService interface {
	GetProducts(ctx context.Context, filter *models.ProductFilter) ([]models.Product, int64, error)
	GetProductDetails(ctx context.Context, code string) (*models.Product, error)
}

type CategoryService interface {
	GetCategories(ctx context.Context) ([]models.Category, error)
	CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error)
}
