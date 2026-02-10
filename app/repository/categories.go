package repository

import (
	"context"
	stderrors "errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mytheresa/go-hiring-challenge/app/repository/entity"
	apperrors "github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/mytheresa/go-hiring-challenge/pkg/model"
	"gorm.io/gorm"
)

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{db: db}
}

func (r *CategoriesRepository) GetCategories(ctx context.Context) ([]model.Category, error) {
	var entities []entity.Category
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return toCategoryDomainList(entities), nil
}

func (r *CategoriesRepository) GetCategoryByCode(ctx context.Context, code string) (*model.Category, error) {
	var e entity.Category
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&e).Error; err != nil {
		return nil, err
	}
	return toCategoryDomain(&e), nil
}

func (r *CategoriesRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	e := toCategoryEntity(category)
	err := r.db.WithContext(ctx).Create(e).Error

	var pgErr *pgconn.PgError
	if stderrors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return apperrors.ErrCategoryAlreadyExists
	}
	if err != nil {
		return err
	}

	category.ID = e.ID
	return nil
}

func toCategoryDomain(e *entity.Category) *model.Category {
	return &model.Category{
		ID:   e.ID,
		Code: e.Code,
		Name: e.Name,
	}
}

func toCategoryDomainList(entities []entity.Category) []model.Category {
	result := make([]model.Category, len(entities))
	for i, e := range entities {
		result[i] = *toCategoryDomain(&e)
	}
	return result
}

func toCategoryEntity(c *model.Category) *entity.Category {
	return &entity.Category{
		ID:   c.ID,
		Code: c.Code,
		Name: c.Name,
	}
}
