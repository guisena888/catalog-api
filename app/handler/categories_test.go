package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/handler"
	"github.com/mytheresa/go-hiring-challenge/app/mocks"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CategoriesHandlerSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	mock    *mocks.MockCategoryService
	handler *handler.CategoriesHandler
}

func (s *CategoriesHandlerSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mock = mocks.NewMockCategoryService(s.ctrl)
	s.handler = handler.NewCategoriesHandler(s.mock)
}

func (s *CategoriesHandlerSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestCategoriesHandlerSuite(t *testing.T) {
	suite.Run(t, new(CategoriesHandlerSuite))
}

// --- GetCategories tests ---

func (s *CategoriesHandlerSuite) TestGetCategories_Success() {
	categories := []models.Category{
		{Code: "clothing", Name: "Clothing"},
		{Code: "shoes", Name: "Shoes"},
	}

	s.mock.EXPECT().
		GetCategories(gomock.Any()).
		Return(categories, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCategories(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	body := rec.Body.String()
	s.Contains(body, `"code":"clothing"`)
	s.Contains(body, `"name":"Clothing"`)
	s.Contains(body, `"code":"shoes"`)
}

func (s *CategoriesHandlerSuite) TestGetCategories_Empty() {
	s.mock.EXPECT().
		GetCategories(gomock.Any()).
		Return([]models.Category{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCategories(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	s.Contains(rec.Body.String(), `[]`)
}

func (s *CategoriesHandlerSuite) TestGetCategories_InternalError() {
	s.mock.EXPECT().
		GetCategories(gomock.Any()).
		Return(nil, errors.New("db down"))

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCategories(rec, req)

	s.Equal(http.StatusInternalServerError, rec.Code)
	s.Contains(rec.Body.String(), `"error":"internal server error"`)
}

// --- CreateCategory tests ---

func (s *CategoriesHandlerSuite) TestCreateCategory_Success() {
	s.mock.EXPECT().
		CreateCategory(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx any, c *models.Category) (*models.Category, error) {
			return c, nil
		})

	body := `{"code":"hats","name":"Hats"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	rec := httptest.NewRecorder()
	s.handler.HandleCreateCategory(rec, req)

	s.Equal(http.StatusCreated, rec.Code)
	s.Contains(rec.Body.String(), `"code":"hats"`)
	s.Contains(rec.Body.String(), `"name":"Hats"`)
}

func (s *CategoriesHandlerSuite) TestCreateCategory_Idempotent() {
	existing := &models.Category{Code: "clothing", Name: "Clothing"}

	s.mock.EXPECT().
		CreateCategory(gomock.Any(), gomock.Any()).
		Return(existing, nil)

	body := `{"code":"clothing","name":"Clothing"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	rec := httptest.NewRecorder()
	s.handler.HandleCreateCategory(rec, req)

	s.Equal(http.StatusCreated, rec.Code)
	s.Contains(rec.Body.String(), `"code":"clothing"`)
	s.Contains(rec.Body.String(), `"name":"Clothing"`)
}

func (s *CategoriesHandlerSuite) TestCreateCategory_InvalidJSON() {
	body := `{invalid`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	rec := httptest.NewRecorder()
	s.handler.HandleCreateCategory(rec, req)

	s.Equal(http.StatusBadRequest, rec.Code)
	s.Contains(rec.Body.String(), `"error":"invalid request body"`)
}

func (s *CategoriesHandlerSuite) TestCreateCategory_MissingFields() {
	body := `{"code":"hats"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	rec := httptest.NewRecorder()
	s.handler.HandleCreateCategory(rec, req)

	s.Equal(http.StatusBadRequest, rec.Code)
	s.Contains(rec.Body.String(), `"error":"code and name are required"`)
}

func (s *CategoriesHandlerSuite) TestCreateCategory_InternalError() {
	s.mock.EXPECT().
		CreateCategory(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("db down"))

	body := `{"code":"hats","name":"Hats"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	rec := httptest.NewRecorder()
	s.handler.HandleCreateCategory(rec, req)

	s.Equal(http.StatusInternalServerError, rec.Code)
	s.Contains(rec.Body.String(), `"error":"internal server error"`)
}
