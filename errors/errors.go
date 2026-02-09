package errors

import (
	"errors"
)

var (
	ErrInvalidOffset        = errors.New("invalid offset")
	ErrInvalidLimit         = errors.New("invalid limit")
	ErrInvalidPriceLessThan = errors.New("invalid priceLessThan value")
	ErrProductNotFound      = errors.New("product not found")
)
