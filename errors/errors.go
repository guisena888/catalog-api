package errors

type Kind int

const (
	KindInvalidInput Kind = iota
	KindNotFound
	KindAlreadyExists
)

type AppError struct {
	Kind    Kind
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrInvalidOffset         = &AppError{Kind: KindInvalidInput, Message: "invalid offset"}
	ErrInvalidLimit          = &AppError{Kind: KindInvalidInput, Message: "invalid limit"}
	ErrInvalidPriceLessThan  = &AppError{Kind: KindInvalidInput, Message: "invalid priceLessThan value"}
	ErrInvalidRequestBody    = &AppError{Kind: KindInvalidInput, Message: "invalid request body"}
	ErrMissingRequiredFields = &AppError{Kind: KindInvalidInput, Message: "code and name are required"}
	ErrMissingProductCode    = &AppError{Kind: KindInvalidInput, Message: "product code is required"}
	ErrProductNotFound       = &AppError{Kind: KindNotFound, Message: "product not found"}
	ErrCategoryAlreadyExists = &AppError{Kind: KindAlreadyExists, Message: "category already exists"}
)
