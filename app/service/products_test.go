package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/mocks"
	"github.com/mytheresa/go-hiring-challenge/app/service"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ProductServiceSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	repo    *mocks.MockProductRepository
	service *service.ProductService
}

func (s *ProductServiceSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mocks.NewMockProductRepository(s.ctrl)
	s.service = service.NewProductService(s.repo)
}

func (s *ProductServiceSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestProductServiceSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceSuite))
}

func (s *ProductServiceSuite) TestGetProductDetails_VariantInheritsPrice() {
	productPrice := decimal.NewFromFloat(10.99)
	product := &models.Product{
		Code:  "PROD001",
		Price: productPrice,
		Variants: []models.Variant{
			{Name: "With Price", SKU: "SKU-A", Price: decimal.NewFromFloat(15.00)},
			{Name: "No Price", SKU: "SKU-B", Price: decimal.Decimal{}},
		},
	}

	s.repo.EXPECT().
		GetProductDetails(gomock.Any(), "PROD001").
		Return(product, nil)

	result, err := s.service.GetProductDetails(context.Background(), "PROD001")

	s.NoError(err)
	s.True(decimal.NewFromFloat(15.00).Equal(result.Variants[0].Price))
	s.True(productPrice.Equal(result.Variants[1].Price))
}

func (s *ProductServiceSuite) TestGetProductDetails_Error() {
	s.repo.EXPECT().
		GetProductDetails(gomock.Any(), "INVALID").
		Return(nil, errors.New("not found"))

	result, err := s.service.GetProductDetails(context.Background(), "INVALID")

	s.Nil(result)
	s.EqualError(err, "not found")
}
