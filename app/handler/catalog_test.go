package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/handler"
	"github.com/mytheresa/go-hiring-challenge/app/mocks"
	apperrors "github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/pkg/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func sampleProducts() []model.Product {
	return []model.Product{
		{Code: "P001", Price: decimal.NewFromFloat(10.99), Category: model.Category{Code: "clothing"}},
		{Code: "P002", Price: decimal.NewFromFloat(5.50), Category: model.Category{Code: "shoes"}},
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
		GetProducts(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f *model.ProductFilter) ([]model.Product, int64, error) {
			s.Equal(0, f.Pagination.Offset)
			s.Equal(10, f.Pagination.Limit)
			s.Empty(f.CategoryCode)
			s.Nil(f.PriceLessThan)
			return sampleProducts(), int64(2), nil
		})

	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *CatalogHandlerSuite) TestGetCatalog_CustomPagination() {
	s.mock.EXPECT().
		GetProducts(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f *model.ProductFilter) ([]model.Product, int64, error) {
			s.Equal(5, f.Pagination.Offset)
			s.Equal(20, f.Pagination.Limit)
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
		GetProducts(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f *model.ProductFilter) ([]model.Product, int64, error) {
			s.Equal("shoes", f.CategoryCode)
			s.Equal(0, f.Pagination.Offset)
			s.Equal(10, f.Pagination.Limit)
			return nil, int64(0), nil
		})

	req := httptest.NewRequest(http.MethodGet, "/catalog?category=shoes", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *CatalogHandlerSuite) TestGetCatalog_PriceLessThanFilter() {
	s.mock.EXPECT().
		GetProducts(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f *model.ProductFilter) ([]model.Product, int64, error) {
			s.Require().NotNil(f.PriceLessThan)
			s.True(decimal.NewFromInt(10).Equal(*f.PriceLessThan))
			return nil, int64(0), nil
		})

	req := httptest.NewRequest(http.MethodGet, "/catalog?priceLessThan=10", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *CatalogHandlerSuite) TestGetCatalog_ResponseIncludesTotalAndCategory() {
	s.mock.EXPECT().
		GetProducts(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f *model.ProductFilter) ([]model.Product, int64, error) {
			return sampleProducts(), int64(42), nil
		})

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
		GetProducts(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f *model.ProductFilter) ([]model.Product, int64, error) {
			return nil, int64(0), errors.New("db down")
		})

	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	rec := httptest.NewRecorder()
	s.handler.HandleGetCatalog(rec, req)

	s.Equal(http.StatusInternalServerError, rec.Code)
	s.Contains(rec.Body.String(), `"error":"internal server error"`)
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
	product := &model.Product{
		Code:     "PROD001",
		Price:    decimal.NewFromFloat(10.99),
		Category: model.Category{Code: "clothing"},
		Variants: []model.Variant{
			{Name: "Variant A", SKU: "SKU001A", Price: decimal.NewFromFloat(11.99)},
			{Name: "Variant B", SKU: "SKU001B", Price: decimal.Decimal{}},
		},
	}

	s.mock.EXPECT().
		GetProductDetails(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, code string) (*model.Product, error) {
			s.Equal("PROD001", code)
			return product, nil
		})

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
	product := &model.Product{
		Code:     "PROD001",
		Price:    decimal.NewFromFloat(10.99),
		Category: model.Category{Code: "clothing"},
		Variants: []model.Variant{
			{Name: "With Price", SKU: "SKU-A", Price: decimal.NewFromFloat(15.00)},
			{Name: "Inherited Price", SKU: "SKU-B", Price: decimal.NewFromFloat(10.99)},
		},
	}

	s.mock.EXPECT().
		GetProductDetails(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, code string) (*model.Product, error) {
			s.Equal("PROD001", code)
			return product, nil
		})

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
		GetProductDetails(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, code string) (*model.Product, error) {
			s.Equal("INVALID", code)
			return nil, apperrors.ErrProductNotFound
		})

	req := httptest.NewRequest(http.MethodGet, "/catalog/INVALID", nil)
	req.SetPathValue("code", "INVALID")
	rec := httptest.NewRecorder()
	s.handler.HandleGetProductDetails(rec, req)

	s.Equal(http.StatusNotFound, rec.Code)
	s.Contains(rec.Body.String(), `"error":"product not found"`)
}

func (s *CatalogHandlerSuite) TestGetProductDetails_InternalError() {
	s.mock.EXPECT().
		GetProductDetails(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, code string) (*model.Product, error) {
			s.Equal("PROD001", code)
			return nil, errors.New("db down")
		})

	req := httptest.NewRequest(http.MethodGet, "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	rec := httptest.NewRecorder()
	s.handler.HandleGetProductDetails(rec, req)

	s.Equal(http.StatusInternalServerError, rec.Code)
	s.Contains(rec.Body.String(), `"error":"internal server error"`)
}

func (s *CatalogHandlerSuite) TestGetProductDetails_EmptyCode() {
	req := httptest.NewRequest(http.MethodGet, "/catalog/", nil)
	req.SetPathValue("code", "")
	rec := httptest.NewRecorder()
	s.handler.HandleGetProductDetails(rec, req)

	s.Equal(http.StatusBadRequest, rec.Code)
	s.Contains(rec.Body.String(), `"error":"product code is required"`)
}
