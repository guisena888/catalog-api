package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/repository"
	"github.com/mytheresa/go-hiring-challenge/app/repository/entity"
	"github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/pkg/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CategoriesRepositorySuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repository.CategoriesRepository
	ctx  context.Context
}

func TestCategoriesRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoriesRepositorySuite))
}

func (s *CategoriesRepositorySuite) SetupSuite() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=password dbname=challenge port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	s.Require().NoError(err)

	s.db = db
	s.repo = repository.NewCategoriesRepository(db)
	s.ctx = context.Background()
}

func (s *CategoriesRepositorySuite) SetupTest() {
	s.truncateTables()
	s.seedTestData()
}

func (s *CategoriesRepositorySuite) truncateTables() {
	s.Require().NoError(s.db.Exec("TRUNCATE TABLE product_variants, products, categories RESTART IDENTITY CASCADE").Error)
}

func (s *CategoriesRepositorySuite) seedTestData() {
	categories := []entity.Category{
		{Code: "clothing", Name: "Clothing"},
		{Code: "shoes", Name: "Shoes"},
	}
	s.Require().NoError(s.db.Create(&categories).Error)
}

func (s *CategoriesRepositorySuite) TestGetCategories_Success() {
	categories, err := s.repo.GetCategories(s.ctx)

	s.NoError(err)
	s.Len(categories, 2)
}

func (s *CategoriesRepositorySuite) TestGetCategoryByCode_Success() {
	category, err := s.repo.GetCategoryByCode(s.ctx, "clothing")

	s.NoError(err)
	s.NotNil(category)
	s.Equal("clothing", category.Code)
	s.Equal("Clothing", category.Name)
}

func (s *CategoriesRepositorySuite) TestGetCategoryByCode_NotFound() {
	category, err := s.repo.GetCategoryByCode(s.ctx, "nonexistent")

	s.Nil(category)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *CategoriesRepositorySuite) TestCreateCategory_Success() {
	category := &model.Category{
		Code: "hats",
		Name: "Hats",
	}

	err := s.repo.CreateCategory(s.ctx, category)

	s.NoError(err)
	s.NotZero(category.ID)

	saved, err := s.repo.GetCategoryByCode(s.ctx, "hats")
	s.NoError(err)
	s.Equal("Hats", saved.Name)
}

func (s *CategoriesRepositorySuite) TestCreateCategory_DuplicateCategory() {
	category := &model.Category{
		Code: "clothing",
		Name: "Different Name",
	}

	err := s.repo.CreateCategory(s.ctx, category)

	s.ErrorIs(err, errors.ErrCategoryAlreadyExists)
}
