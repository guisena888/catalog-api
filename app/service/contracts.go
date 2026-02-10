package service

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductRepository interface {
	GetProducts(ctx context.Context, filter *models.ProductFilter) ([]models.Product, int64, error)
	GetProductDetails(ctx context.Context, code string) (*models.Product, error)
}

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetCategoryByCode(ctx context.Context, code string) (*models.Category, error)
	CreateCategory(ctx context.Context, category *models.Category) error
}
