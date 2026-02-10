package models

import (
	"context"
	stderrors "errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	apperrors "github.com/mytheresa/go-hiring-challenge/errors"
	"gorm.io/gorm"
)

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{db: db}
}

func (r *CategoriesRepository) GetCategories(ctx context.Context) ([]Category, error) {
	var categories []Category
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepository) GetCategoryByCode(ctx context.Context, code string) (*Category, error) {
	var category Category
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoriesRepository) CreateCategory(ctx context.Context, category *Category) error {
	err := r.db.WithContext(ctx).Create(category).Error

	var pgErr *pgconn.PgError
	if stderrors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return apperrors.ErrCategoryAlreadyExists
	}
	if err != nil {
		return err
	}
	return nil
}
