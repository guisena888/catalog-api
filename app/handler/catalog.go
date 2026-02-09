package handler

import (
	"errors"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	apperrors "github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CatalogHandler struct {
	products ProductService
}

func NewCatalogHandler(s ProductService) *CatalogHandler {
	return &CatalogHandler{
		products: s,
	}
}

func (h *CatalogHandler) HandleGetCatalog(w http.ResponseWriter, r *http.Request) {
	filter := new(models.ProductFilter)
	err := filter.Parse(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, total, err := h.products.GetProducts(filter)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	products := make([]ProductResponse, len(res))
	for i, p := range res {
		products[i] = toProductResponse(p)
	}

	api.OKResponse(w, CatalogResponse{
		Products: products,
		Total:    total,
	})
}

func (h *CatalogHandler) HandleGetProductDetails(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	product, err := h.products.GetProductDetails(code)
	if err != nil {
		if errors.Is(err, apperrors.ErrProductNotFound) {
			api.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, toProductDetailResponse(*product))
}
