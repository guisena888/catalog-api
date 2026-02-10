package repository

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/repository/entity"
	"github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/pkg/model"
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetProducts(ctx context.Context, filter *model.ProductFilter) ([]model.Product, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Product{})

	if filter.CategoryCode != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", filter.CategoryCode)
	}

	if filter.PriceLessThan != nil {
		query = query.Where("products.price < ?", *filter.PriceLessThan)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var entities []entity.Product
	if err := query.Preload("Category").Preload("Variants").
		Offset(filter.Pagination.Offset).Limit(filter.Pagination.Limit).
		Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return toProductDomainList(entities), total, nil
}

func (r *ProductsRepository) GetProductDetails(ctx context.Context, code string) (*model.Product, error) {
	var e entity.Product
	err := r.db.WithContext(ctx).Preload("Category").Preload("Variants").
		Where("code = ?", code).First(&e).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return toProductDomain(&e), nil
}

func toProductDomain(e *entity.Product) *model.Product {
	variants := make([]model.Variant, len(e.Variants))
	for i, v := range e.Variants {
		variants[i] = *toVariantDomain(&v)
	}

	return &model.Product{
		ID:    e.ID,
		Code:  e.Code,
		Price: e.Price,
		Category: model.Category{
			ID:   e.Category.ID,
			Code: e.Category.Code,
			Name: e.Category.Name,
		},
		Variants: variants,
	}
}

func toProductDomainList(entities []entity.Product) []model.Product {
	result := make([]model.Product, len(entities))
	for i, e := range entities {
		result[i] = *toProductDomain(&e)
	}
	return result
}

func toVariantDomain(e *entity.Variant) *model.Variant {
	return &model.Variant{
		ID:    e.ID,
		Name:  e.Name,
		SKU:   e.SKU,
		Price: e.Price,
	}
}
