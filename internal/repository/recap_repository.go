package repository

import (
	"context"

	"github.com/admalfrizi/weekly-wrapped-be/internal/model"
)

type RecapRepository interface {
	Upsert(ctx context.Context, recap *model.WeeklyRecap) error
	GetBySlug(ctx context.Context, slug string) (*model.WeeklyRecap, error)
}

type recapRepository struct {
	*BaseRepository
}

func NewRecapRepository(base *BaseRepository) RecapRepository {
	return &recapRepository{BaseRepository: base}
}

func (r *recapRepository) Upsert(ctx context.Context, recap *model.WeeklyRecap) error {
	query := `
		INSERT INTO weekly_recaps 
			(id, user_id, week_start, week_end, slug, stats_snapshot, narrative, created_at)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, week_start) 
		DO UPDATE SET 
			slug = EXCLUDED.slug,
			stats_snapshot = EXCLUDED.stats_snapshot,
			narrative = EXCLUDED.narrative;
	`

	_, err := r.db.Exec(ctx, query, 
		recap.ID, 
		recap.UserID, 
		recap.WeekStart, 
		recap.WeekEnd, 
		recap.Slug, 
		recap.StatsSnapshot, 
		recap.Narrative, 
		recap.CreatedAt,
	)
	
	return err
}

func (r *recapRepository) GetBySlug(ctx context.Context, slug string) (*model.WeeklyRecap, error) {
	query := `
		SELECT id, user_id, week_start, week_end, slug, stats_snapshot, narrative, created_at
		FROM weekly_recaps
		WHERE slug = $1
		LIMIT 1;
	`

	var recap model.WeeklyRecap
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&recap.ID,
		&recap.UserID,
		&recap.WeekStart,
		&recap.WeekEnd,
		&recap.Slug,
		&recap.StatsSnapshot,
		&recap.Narrative,
		&recap.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &recap, nil
}