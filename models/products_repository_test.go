package models_test

import (
	"context"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProductsRepositorySuite struct {
	suite.Suite
	db   *gorm.DB
	repo *models.ProductsRepository
	ctx  context.Context
}

func TestProductsRepositorySuite(t *testing.T) {
	suite.Run(t, new(ProductsRepositorySuite))
}

func (s *ProductsRepositorySuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)

	s.Require().NoError(db.AutoMigrate(&models.Category{}, &models.Product{}, &models.Variant{}))

	s.db = db
	s.repo = models.NewProductsRepository(db)
	s.ctx = context.Background()

	s.seedTestData()
}

func (s *ProductsRepositorySuite) seedTestData() {
	categories := []models.Category{
		{ID: 1, Code: "clothing", Name: "Clothing"},
		{ID: 2, Code: "shoes", Name: "Shoes"},
		{ID: 3, Code: "accessories", Name: "Accessories"},
	}
	s.Require().NoError(s.db.Create(&categories).Error)

	products := []models.Product{
		{ID: 1, Code: "PROD001", Price: decimal.NewFromFloat(100.00), CategoryID: 1},
		{ID: 2, Code: "PROD002", Price: decimal.NewFromFloat(50.00), CategoryID: 2},
		{ID: 3, Code: "PROD003", Price: decimal.NewFromFloat(200.00), CategoryID: 3},
		{ID: 4, Code: "PROD004", Price: decimal.NewFromFloat(75.00), CategoryID: 1},
	}
	s.Require().NoError(s.db.Create(&products).Error)

	variants := []models.Variant{
		{ID: 1, ProductID: 1, Name: "Small", SKU: "PROD001-S", Price: decimal.NewFromFloat(95.00)},
		{ID: 2, ProductID: 1, Name: "Large", SKU: "PROD001-L", Price: decimal.Decimal{}},
	}
	s.Require().NoError(s.db.Create(&variants).Error)
}

func (s *ProductsRepositorySuite) TestGetProducts_DefaultPagination() {
	filter := &models.ProductFilter{
		Pagination: models.PaginationInput{Offset: 0, Limit: 10},
	}

	products, total, err := s.repo.GetProducts(s.ctx, filter)

	s.NoError(err)
	s.Equal(int64(4), total)
	s.Len(products, 4)
}

func (s *ProductsRepositorySuite) TestGetProducts_WithPagination() {
	filter := &models.ProductFilter{
		Pagination: models.PaginationInput{Offset: 1, Limit: 2},
	}

	products, total, err := s.repo.GetProducts(s.ctx, filter)

	s.NoError(err)
	s.Equal(int64(4), total)
	s.Len(products, 2)
}

func (s *ProductsRepositorySuite) TestGetProducts_FilterByCategory() {
	filter := &models.ProductFilter{
		CategoryCode: "clothing",
		Pagination:   models.PaginationInput{Offset: 0, Limit: 10},
	}

	products, total, err := s.repo.GetProducts(s.ctx, filter)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(products, 2)
	for _, p := range products {
		s.Equal("clothing", p.Category.Code)
	}
}

func (s *ProductsRepositorySuite) TestGetProducts_FilterByPriceLessThan() {
	price := decimal.NewFromFloat(80.00)
	filter := &models.ProductFilter{
		PriceLessThan: &price,
		Pagination:    models.PaginationInput{Offset: 0, Limit: 10},
	}

	products, total, err := s.repo.GetProducts(s.ctx, filter)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(products, 2)
	for _, p := range products {
		s.True(p.Price.LessThan(price))
	}
}

func (s *ProductsRepositorySuite) TestGetProducts_CombinedFilters() {
	price := decimal.NewFromFloat(80.00)
	filter := &models.ProductFilter{
		CategoryCode:  "clothing",
		PriceLessThan: &price,
		Pagination:    models.PaginationInput{Offset: 0, Limit: 10},
	}

	products, total, err := s.repo.GetProducts(s.ctx, filter)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(products, 1)
	s.Equal("PROD004", products[0].Code)
}

func (s *ProductsRepositorySuite) TestGetProducts_PreloadsCategory() {
	filter := &models.ProductFilter{
		Pagination: models.PaginationInput{Offset: 0, Limit: 10},
	}

	products, _, err := s.repo.GetProducts(s.ctx, filter)

	s.NoError(err)
	for _, p := range products {
		s.NotEmpty(p.Category.Code)
	}
}

func (s *ProductsRepositorySuite) TestGetProductDetails_Success() {
	product, err := s.repo.GetProductDetails(s.ctx, "PROD001")

	s.NoError(err)
	s.NotNil(product)
	s.Equal("PROD001", product.Code)
	s.Equal("clothing", product.Category.Code)
	s.Len(product.Variants, 2)
}

func (s *ProductsRepositorySuite) TestGetProductDetails_NotFound() {
	product, err := s.repo.GetProductDetails(s.ctx, "NONEXISTENT")

	s.Nil(product)
	s.ErrorIs(err, errors.ErrProductNotFound)
}

func (s *ProductsRepositorySuite) TestGetProductDetails_IncludesVariants() {
	product, err := s.repo.GetProductDetails(s.ctx, "PROD001")

	s.NoError(err)
	s.Len(product.Variants, 2)
	s.Equal("Small", product.Variants[0].Name)
	s.Equal("PROD001-S", product.Variants[0].SKU)
}
