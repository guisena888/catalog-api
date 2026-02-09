package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/handler"
	"github.com/mytheresa/go-hiring-challenge/app/mocks"
	apperrors "github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func sampleProducts() []models.Product {
	return []models.Product{
		{Code: "P001", Price: decimal.NewFromFloat(10.99), Category: models.Category{Code: "clothing"}},
		{Code: "P002", Price: decimal.NewFromFloat(5.50), Category: models.Category{Code: "shoes"}},
	}
}

type CatalogHandlerSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	mock    *mocks.MockProductService
	handler *handler.CatalogHandler
}

func (s *CatalogHandlerSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mock = mocks.NewMockProductService(s.ctrl)
	s.handler = handler.NewCatalogHandler(s.mock)
}

func (s *CatalogHandlerSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestCatalogHandlerSuite(t *testing.T) {
	suite.Run(t, new(CatalogHandlerSuite))
}

func (s *CatalogHandlerSuite) TestGetCatalog_DefaultPagination() {
	s.mock.EXPECT().
		GetProducts(gomock.Any()).
		DoAndReturn(func(filter *models.ProductFilter) ([]models.Product, int64, error) {
			assert.Equal(s.T(), 0, filter.Pagination.Offset)
			assert.Equal(s.T(), 10, filter.Pagination.Limit)
			return sampleProducts(), int64(2), nil
		})

	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *CatalogHandlerSuite) TestGetCatalog_CustomPagination() {
	s.mock.EXPECT().
		GetProducts(gomock.Any()).
		DoAndReturn(func(filter *models.ProductFilter) ([]models.Product, int64, error) {
			assert.Equal(s.T(), 5, filter.Pagination.Offset)
			assert.Equal(s.T(), 20, filter.Pagination.Limit)
			return sampleProducts(), int64(2), nil
		})

	req := httptest.NewRequest(http.MethodGet, "/catalog?offset=5&limit=20", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *CatalogHandlerSuite) TestGetCatalog_InvalidLimit() {
	req := httptest.NewRequest(http.MethodGet, "/catalog?limit=0", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/catalog?limit=200", nil)
	rec = httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)
	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *CatalogHandlerSuite) TestGetCatalog_CategoryFilter() {
	s.mock.EXPECT().
		GetProducts(gomock.Any()).
		DoAndReturn(func(filter *models.ProductFilter) ([]models.Product, int64, error) {
			assert.Equal(s.T(), "shoes", filter.CategoryCode)
			return nil, int64(0), nil
		})

	req := httptest.NewRequest(http.MethodGet, "/catalog?category=shoes", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *CatalogHandlerSuite) TestGetCatalog_PriceLessThanFilter() {
	s.mock.EXPECT().
		GetProducts(gomock.Any()).
		DoAndReturn(func(filter *models.ProductFilter) ([]models.Product, int64, error) {
			expected := decimal.NewFromInt(10)
			assert.True(s.T(), expected.Equal(*filter.PriceLessThan))
			return nil, int64(0), nil
		})

	req := httptest.NewRequest(http.MethodGet, "/catalog?priceLessThan=10", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *CatalogHandlerSuite) TestGetCatalog_ResponseIncludesTotalAndCategory() {
	s.mock.EXPECT().
		GetProducts(gomock.Any()).
		Return(sampleProducts(), int64(42), nil)

	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	body := rec.Body.String()
	s.Contains(body, `"total":42`)
	s.Contains(body, `"category":"clothing"`)
	s.Contains(body, `"category":"shoes"`)
}

func (s *CatalogHandlerSuite) TestGetCatalog_RepositoryError() {
	s.mock.EXPECT().
		GetProducts(gomock.Any()).
		Return(nil, int64(0), errors.New("db down"))

	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusInternalServerError, rec.Code)
	s.Contains(rec.Body.String(), `"error":"db down"`)
}

func (s *CatalogHandlerSuite) TestGetCatalog_InvalidPriceLessThan() {
	req := httptest.NewRequest(http.MethodGet, "/catalog?priceLessThan=abc", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusBadRequest, rec.Code)
	s.Contains(rec.Body.String(), `"error"`)
}

func (s *CatalogHandlerSuite) TestGetCatalog_NegativePriceLessThan() {
	req := httptest.NewRequest(http.MethodGet, "/catalog?priceLessThan=-5", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *CatalogHandlerSuite) TestGetProductDetails_Success() {
	product := &models.Product{
		Code:     "PROD001",
		Price:    decimal.NewFromFloat(10.99),
		Category: models.Category{Code: "clothing"},
		Variants: []models.Variant{
			{Name: "Variant A", SKU: "SKU001A", Price: decimal.NewFromFloat(11.99)},
			{Name: "Variant B", SKU: "SKU001B", Price: decimal.Decimal{}},
		},
	}

	s.mock.EXPECT().
		GetProductDetails("PROD001").
		Return(product, nil)

	req := httptest.NewRequest(http.MethodGet, "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	rec := httptest.NewRecorder()
	s.handler.HandleGetProductDetails(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	body := rec.Body.String()
	s.Contains(body, `"code":"PROD001"`)
	s.Contains(body, `"category":"clothing"`)
	s.Contains(body, `"variants"`)
}

func (s *CatalogHandlerSuite) TestGetProductDetails_VariantPrices() {
	product := &models.Product{
		Code:     "PROD001",
		Price:    decimal.NewFromFloat(10.99),
		Category: models.Category{Code: "clothing"},
		Variants: []models.Variant{
			{Name: "With Price", SKU: "SKU-A", Price: decimal.NewFromFloat(15.00)},
			{Name: "Inherited Price", SKU: "SKU-B", Price: decimal.NewFromFloat(10.99)},
		},
	}

	s.mock.EXPECT().
		GetProductDetails("PROD001").
		Return(product, nil)

	req := httptest.NewRequest(http.MethodGet, "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	rec := httptest.NewRecorder()
	s.handler.HandleGetProductDetails(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	body := rec.Body.String()
	s.Contains(body, `"price":15`)
	s.Contains(body, `"price":10.99`)
}

func (s *CatalogHandlerSuite) TestGetProductDetails_NotFound() {
	s.mock.EXPECT().
		GetProductDetails("INVALID").
		Return(nil, apperrors.ErrProductNotFound)

	req := httptest.NewRequest(http.MethodGet, "/catalog/INVALID", nil)
	req.SetPathValue("code", "INVALID")
	rec := httptest.NewRecorder()
	s.handler.HandleGetProductDetails(rec, req)

	s.Equal(http.StatusNotFound, rec.Code)
	s.Contains(rec.Body.String(), `"error":"product not found"`)
}

func (s *CatalogHandlerSuite) TestGetProductDetails_InternalError() {
	s.mock.EXPECT().
		GetProductDetails("PROD001").
		Return(nil, errors.New("db down"))

	req := httptest.NewRequest(http.MethodGet, "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	rec := httptest.NewRecorder()
	s.handler.HandleGetProductDetails(rec, req)

	s.Equal(http.StatusInternalServerError, rec.Code)
	s.Contains(rec.Body.String(), `"error":"db down"`)
}
