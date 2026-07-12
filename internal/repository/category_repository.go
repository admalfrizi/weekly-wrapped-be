package repository

import (
	"context"

	"github.com/admalfrizi/weekly-wrapped-be/internal/model"
)

type CategoryRepository interface {
	List(ctx context.Context, userID string, limit, offset int) ([]model.Category, int, error)
}

type categoryRepository struct {
	*BaseRepository
}

func NewCategoryRepository(base *BaseRepository) CategoryRepository {
	return &categoryRepository{BaseRepository: base}
}

func (c *categoryRepository) List(ctx context.Context, userID string, limit int, offset int) ([]model.Category, int, error) {
	panic("unimplemented")
}
