package repository

import (
	"context"
	"errors"

	"github.com/admalfrizi/weekly-wrapped-be/internal/model"
	"github.com/admalfrizi/weekly-wrapped-be/internal/query"
	"github.com/jackc/pgx/v5"
)

type ActivityRepository interface {
	Create(ctx context.Context, activity *model.Activity) error
	List(ctx context.Context, userID string, limit, offset int) ([]model.Activity, int, error)
	GetByID(ctx context.Context, activityID, userID string) (*model.Activity, error)
	Update(ctx context.Context, activity *model.Activity) error
	Delete(ctx context.Context, activityID, userID string) error
}

type activityRepository struct {
	*BaseRepository
}

func NewActivityRepository(base *BaseRepository) ActivityRepository {
	return &activityRepository{BaseRepository: base}
}

func (r *activityRepository) Create(ctx context.Context, activity *model.Activity) error {
	err := r.db.QueryRow(ctx, query.InsertActivity,
		activity.UserID,
		activity.CategoryID,
		activity.Value,
		activity.Note,
		activity.OccurredAt,
	).Scan(&activity.ID, &activity.CreatedAt)

	return err
}

func (r *activityRepository) List(ctx context.Context, userID string, limit, offset int) ([]model.Activity, int, error) {
	var totalItems int
	err := r.db.QueryRow(ctx, query.CountActivities, userID).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, query.ListActivities, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var activities []model.Activity
	for rows.Next() {
		var a model.Activity
		var c model.Category 
		
		err := rows.Scan(
			&a.ID, &a.UserID, &a.CategoryID, &a.Value, &a.Note, &a.OccurredAt, &a.CreatedAt,
			&c.ID, &c.Name,
		)
		if err != nil {
			return nil, 0, err
		}
		
		a.Category = &c 
		activities = append(activities, a)
	}

	return activities, totalItems, nil
}

func (r *activityRepository) GetByID(ctx context.Context, activityID, userID string) (*model.Activity, error) {
	var a model.Activity
	var c model.Category

	err := r.db.QueryRow(ctx, query.GetActivityByID, activityID, userID).Scan(
		&a.ID, &a.UserID, &a.CategoryID, &a.Value, &a.Note, &a.OccurredAt, &a.CreatedAt,
		&c.ID, &c.Name,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("activity not found")
		}
		return nil, err
	}

	a.Category = &c
	return &a, nil
}

func (r *activityRepository) Update(ctx context.Context, activity *model.Activity) error {
	var returnedID string
	err := r.db.QueryRow(ctx, query.UpdateActivity,
		activity.CategoryID,
		activity.Value,
		activity.Note,
		activity.OccurredAt,
		activity.ID, 
		activity.UserID,
	).Scan(&returnedID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("activity not found")
		}
		return err
	}
	return nil
}

func (r *activityRepository) Delete(ctx context.Context, activityID, userID string) error {
	tag, err := r.db.Exec(ctx, query.DeleteActivity, activityID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("activity not found")
	}
	return nil
}