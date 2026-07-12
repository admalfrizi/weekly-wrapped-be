package repository

import (
	"context"

	"github.com/admalfrizi/weekly-wrapped-be/internal/model"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]model.Category, error)
}

type categoryRepository struct {
	*BaseRepository
}

func NewCategoryRepository(base *BaseRepository) CategoryRepository {
	return &categoryRepository{BaseRepository: base}
}

func (c *categoryRepository) List(ctx context.Context) ([]model.Category, error) {
	query := `
		SELECT id, name, created_at 
		FROM categories 
		ORDER BY name ASC;
	`

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	
	for rows.Next() {
		var cat model.Category
		err := rows.Scan(
			&cat.ID, 
			&cat.Name,
			&cat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	if categories == nil {
		categories = make([]model.Category, 0)
	}

	return categories, nil
}

