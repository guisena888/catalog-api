package service

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/pkg/model"
)

type ProductRepository interface {
	GetProducts(ctx context.Context, filter *model.ProductFilter) ([]model.Product, int64, error)
	GetProductDetails(ctx context.Context, code string) (*model.Product, error)
}

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]model.Category, error)
	GetCategoryByCode(ctx context.Context, code string) (*model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) error
}
