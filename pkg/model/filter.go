package model

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/shopspring/decimal"
)

type ProductFilter struct {
	CategoryCode  string
	PriceLessThan *decimal.Decimal
	Pagination    PaginationInput
}

func (f *ProductFilter) Parse(req *http.Request) error {
	q := req.URL.Query()

	f.CategoryCode = q.Get("category")

	if plt := q.Get("priceLessThan"); plt != "" {
		d, err := decimal.NewFromString(plt)
		if err != nil {
			return errors.ErrInvalidPriceLessThan
		}
		if !d.IsPositive() {
			return errors.ErrInvalidPriceLessThan
		}
		f.PriceLessThan = &d
	}

	err := f.Pagination.Parse(q)
	if err != nil {
		return err
	}

	return nil
}
