package service

import (
	"context"
	"errors"

	apperrors "github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/pkg/model"
)

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(r CategoryRepository) *CategoryService {
	return &CategoryService{repo: r}
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]model.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	err := s.repo.CreateCategory(ctx, category)
	if errors.Is(err, apperrors.ErrCategoryAlreadyExists) {
		return s.repo.GetCategoryByCode(ctx, category.Code)
	}
	if err != nil {
		return nil, err
	}
	return category, nil
}
