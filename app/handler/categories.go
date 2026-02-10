package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mytheresa/go-hiring-challenge/app/api"
	apperrors "github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/pkg/model"
)

var validate = validator.New()

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
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func (h *CategoriesHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.HandleError(w, apperrors.ErrInvalidRequestBody)
		return
	}

	if err := validate.Struct(req); err != nil {
		api.HandleError(w, apperrors.ErrInvalidRequestBody)
		return
	}

	category := &model.Category{
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
