package model

import (
	"github.com/shopspring/decimal"
)

// Product represents a product in the domain.
type Product struct {
	ID       uint
	Code     string
	Price    decimal.Decimal
	Category Category
	Variants []Variant
}
