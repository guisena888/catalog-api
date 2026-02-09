package models

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/shopspring/decimal"
)

// Product represents a product in the catalog.
// It includes a unique code and a price.
type Product struct {
	ID         uint            `gorm:"primaryKey"`
	Code       string          `gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	CategoryID uint            `gorm:"not null"`
	Category   Category        `gorm:"foreignKey:CategoryID"`
	Variants   []Variant       `gorm:"foreignKey:ProductID"`
}

func (p *Product) TableName() string {
	return "products"
}

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
