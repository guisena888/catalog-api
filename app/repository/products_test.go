package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/repository"
	"github.com/mytheresa/go-hiring-challenge/app/repository/entity"
	"github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/pkg/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProductsRepositorySuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repository.ProductsRepository
	ctx  context.Context
}

func TestProductsRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductsRepositorySuite))
}

func (s *ProductsRepositorySuite) SetupSuite() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=password dbname=challenge port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	s.Require().NoError(err)

	s.db = db
	s.repo = repository.NewProductsRepository(db)
	s.ctx = context.Background()
}

func (s *ProductsRepositorySuite) SetupTest() {
	s.truncateTables()
	s.seedTestData()
}

func (s *ProductsRepositorySuite) truncateTables() {
	s.Require().NoError(s.db.Exec("TRUNCATE TABLE product_variants, products, categories RESTART IDENTITY CASCADE").Error)
}

func (s *ProductsRepositorySuite) seedTestData() {
	categories := []entity.Category{
		{Code: "clothing", Name: "Clothing"},
		{Code: "shoes", Name: "Shoes"},
		{Code: "accessories", Name: "Accessories"},
	}
	s.Require().NoError(s.db.Create(&categories).Error)

	products := []entity.Product{
		{Code: "PROD001", Price: decimal.NewFromFloat(100.00), CategoryID: categories[0].ID},
		{Code: "PROD002", Price: decimal.NewFromFloat(50.00), CategoryID: categories[1].ID},
		{Code: "PROD003", Price: decimal.NewFromFloat(200.00), CategoryID: categories[2].ID},
		{Code: "PROD004", Price: decimal.NewFromFloat(75.00), CategoryID: categories[0].ID},
	}
	s.Require().NoError(s.db.Create(&products).Error)

	variants := []entity.Variant{
		{ProductID: products[0].ID, Name: "Small", SKU: "PROD001-S", Price: decimal.NewFromFloat(95.00)},
		{ProductID: products[0].ID, Name: "Large", SKU: "PROD001-L", Price: decimal.Decimal{}},
	}
	s.Require().NoError(s.db.Create(&variants).Error)
}

func (s *ProductsRepositorySuite) TestGetProducts_DefaultPagination() {
	filter := &model.ProductFilter{
		Pagination: model.PaginationInput{Offset: 0, Limit: 10},
	}

	products, total, err := s.repo.GetProducts(s.ctx, filter)

	s.NoError(err)
	s.Equal(int64(4), total)
	s.Len(products, 4)
}

func (s *ProductsRepositorySuite) TestGetProducts_WithPagination() {
	filter := &model.ProductFilter{
		Pagination: model.PaginationInput{Offset: 1, Limit: 2},
	}

	products, total, err := s.repo.GetProducts(s.ctx, filter)

	s.NoError(err)
	s.Equal(int64(4), total)
	s.Len(products, 2)
}

func (s *ProductsRepositorySuite) TestGetProducts_FilterByCategory() {
	filter := &model.ProductFilter{
		CategoryCode: "clothing",
		Pagination:   model.PaginationInput{Offset: 0, Limit: 10},
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
	filter := &model.ProductFilter{
		PriceLessThan: &price,
		Pagination:    model.PaginationInput{Offset: 0, Limit: 10},
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
	filter := &model.ProductFilter{
		CategoryCode:  "clothing",
		PriceLessThan: &price,
		Pagination:    model.PaginationInput{Offset: 0, Limit: 10},
	}

	products, total, err := s.repo.GetProducts(s.ctx, filter)

	s.NoError(err)
	s.Equal(int64(1), total)
	s.Len(products, 1)
	s.Equal("PROD004", products[0].Code)
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
	s.Equal("100", product.Price.StringFixed(0))

	variantNames := make(map[string]bool)
	for _, v := range product.Variants {
		variantNames[v.Name] = true
	}
	s.True(variantNames["Small"])
	s.True(variantNames["Large"])
}
