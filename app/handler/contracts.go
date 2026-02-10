package handler

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/pkg/model"
)

type ProductService interface {
	GetProducts(ctx context.Context, filter *model.ProductFilter) ([]model.Product, int64, error)
	GetProductDetails(ctx context.Context, code string) (*model.Product, error)
}

type CategoryService interface {
	GetCategories(ctx context.Context) ([]model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
}
