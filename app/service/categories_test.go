package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/mocks"
	"github.com/mytheresa/go-hiring-challenge/app/service"
	apperrors "github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CategoryServiceSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	repo    *mocks.MockCategoryRepository
	service *service.CategoryService
}

func (s *CategoryServiceSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mocks.NewMockCategoryRepository(s.ctrl)
	s.service = service.NewCategoryService(s.repo)
}

func (s *CategoryServiceSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestCategoryServiceSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceSuite))
}

func (s *CategoryServiceSuite) TestCreateCategory_Success() {
	category := &models.Category{Code: "hats", Name: "Hats"}

	s.repo.EXPECT().
		CreateCategory(gomock.Any(), category).
		Return(nil)

	result, err := s.service.CreateCategory(context.Background(), category)

	s.NoError(err)
	s.Equal("hats", result.Code)
	s.Equal("Hats", result.Name)
}

func (s *CategoryServiceSuite) TestCreateCategory_Idempotent() {
	category := &models.Category{Code: "clothing", Name: "Clothing"}
	existing := &models.Category{ID: 1, Code: "clothing", Name: "Clothing"}

	s.repo.EXPECT().
		CreateCategory(gomock.Any(), category).
		Return(apperrors.ErrCategoryAlreadyExists)

	s.repo.EXPECT().
		GetCategoryByCode(gomock.Any(), "clothing").
		Return(existing, nil)

	result, err := s.service.CreateCategory(context.Background(), category)

	s.NoError(err)
	s.Equal(uint(1), result.ID)
	s.Equal("clothing", result.Code)
}

func (s *CategoryServiceSuite) TestCreateCategory_Error() {
	category := &models.Category{Code: "hats", Name: "Hats"}

	s.repo.EXPECT().
		CreateCategory(gomock.Any(), category).
		Return(errors.New("db down"))

	result, err := s.service.CreateCategory(context.Background(), category)

	s.Nil(result)
	s.EqualError(err, "db down")
}
