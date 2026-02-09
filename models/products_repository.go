package models

import (
	"github.com/mytheresa/go-hiring-challenge/errors"
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

func (r *ProductsRepository) GetProducts(filter *ProductFilter) ([]Product, int64, error) {
	query := r.db.Model(&Product{})

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

	var products []Product
	if err := query.Preload("Category").Preload("Variants").
		Offset(filter.Pagination.Offset).Limit(filter.Pagination.Limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductsRepository) GetProductDetails(code string) (*Product, error) {
	var product Product
	err := r.db.Preload("Category").Preload("Variants").
		Where("code = ?", code).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}
