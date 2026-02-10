package handler

import (
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
	if err := filter.Parse(r); err != nil {
		api.HandleError(w, err)
		return
	}

	res, total, err := h.products.GetProducts(r.Context(), filter)
	if err != nil {
		api.HandleError(w, err)
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
	if code == "" {
		api.HandleError(w, apperrors.ErrMissingProductCode)
		return
	}

	product, err := h.products.GetProductDetails(r.Context(), code)
	if err != nil {
		api.HandleError(w, err)
		return
	}

	api.OKResponse(w, toProductDetailResponse(*product))
}
