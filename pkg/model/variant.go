package model

import (
	"github.com/shopspring/decimal"
)

// Variant represents a product variant in the domain.
type Variant struct {
	ID    uint
	Name  string
	SKU   string
	Price decimal.Decimal
}
