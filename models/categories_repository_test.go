package models_test

import (
	"context"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CategoriesRepositorySuite struct {
	suite.Suite
	db   *gorm.DB
	repo *models.CategoriesRepository
	ctx  context.Context
}

func TestCategoriesRepositorySuite(t *testing.T) {
	suite.Run(t, new(CategoriesRepositorySuite))
}

func (s *CategoriesRepositorySuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)

	s.Require().NoError(db.AutoMigrate(&models.Category{}))

	s.db = db
	s.repo = models.NewCategoriesRepository(db)
	s.ctx = context.Background()

	s.seedTestData()
}

func (s *CategoriesRepositorySuite) seedTestData() {
	categories := []models.Category{
		{ID: 1, Code: "clothing", Name: "Clothing"},
		{ID: 2, Code: "shoes", Name: "Shoes"},
	}
	s.Require().NoError(s.db.Create(&categories).Error)
}

func (s *CategoriesRepositorySuite) TestGetCategories_Success() {
	categories, err := s.repo.GetCategories(s.ctx)

	s.NoError(err)
	s.Len(categories, 2)
}

func (s *CategoriesRepositorySuite) TestGetCategories_Empty() {
	s.db.Exec("DELETE FROM categories")

	categories, err := s.repo.GetCategories(s.ctx)

	s.NoError(err)
	s.Empty(categories)
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
	s.Error(err)
}

func (s *CategoriesRepositorySuite) TestCreateCategory_Success() {
	category := &models.Category{
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

func (s *CategoriesRepositorySuite) TestCreateCategory_DuplicateCode() {
	category := &models.Category{
		Code: "clothing", // Already exists
		Name: "Different Name",
	}

	err := s.repo.CreateCategory(s.ctx, category)

	// Note: With PostgreSQL this returns ErrCategoryAlreadyExists,
	// but SQLite returns a different error. We just verify an error occurs.
	s.Error(err)
}
