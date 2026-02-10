package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	apperrors "github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoriesHandler struct {
	categories CategoryService
}

func NewCategoriesHandler(s CategoryService) *CategoriesHandler {
	return &CategoriesHandler{categories: s}
}

func (h *CategoriesHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categories.GetCategories(r.Context())
	if err != nil {
		api.HandleError(w, err)
		return
	}

	response := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		response[i] = toCategoryResponse(c)
	}

	api.OKResponse(w, response)
}

type createCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (h *CategoriesHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.HandleError(w, apperrors.ErrInvalidRequestBody)
		return
	}

	if req.Code == "" || req.Name == "" {
		api.HandleError(w, apperrors.ErrMissingRequiredFields)
		return
	}

	category := &models.Category{
		Code: req.Code,
		Name: req.Name,
	}

	result, err := h.categories.CreateCategory(r.Context(), category)
	if err != nil {
		api.HandleError(w, err)
		return
	}

	api.CreatedResponse(w, toCategoryResponse(*result))
}
