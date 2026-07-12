package service

import (
	"context"
	"errors"

	"github.com/admalfrizi/weekly-wrapped-be/internal/dto"
	"github.com/admalfrizi/weekly-wrapped-be/internal/model"
	"github.com/admalfrizi/weekly-wrapped-be/internal/repository"
	"github.com/google/uuid"
)

type ActivityService interface {
	Create(ctx context.Context, userID string, req dto.CreateActivityRequest) (*model.Activity, error)
	List(ctx context.Context, userID string, page, limit int) ([]model.Activity, int, error)
	ListCategory(ctx context.Context) ([]model.Category, error)
	GetByID(ctx context.Context, userID, activityID string) (*model.Activity, error)
	Update(ctx context.Context, userID, activityID string, req dto.UpdateActivityRequest) (*model.Activity, error)
	Delete(ctx context.Context, userID, activityID string) error
}

type activityService struct {
	repo         repository.ActivityRepository
	categoryRepo repository.CategoryRepository
}

func NewActivityService(repo repository.ActivityRepository, category repository.CategoryRepository) ActivityService {
	return &activityService{
		repo:         repo,
		categoryRepo: category,
	}
}

func (s *activityService) Create(ctx context.Context, userID string, req dto.CreateActivityRequest) (*model.Activity, error) {
	userUUID, _ := uuid.Parse(userID)
	categoryUUID, _ := uuid.Parse(req.CategoryID)

	activity := &model.Activity{
		UserID:     userUUID,
		CategoryID: categoryUUID,
		Value:      req.Value,
		Note:       req.Note,
		OccurredAt: req.OccurredAt,
	}

	if err := s.repo.Create(ctx, activity); err != nil {
		return nil, errors.New("wdqwdqwf")
	}

	return activity, nil
}

func (s *activityService) List(ctx context.Context, userID string, page, limit int) ([]model.Activity, int, error) {
	offset := (page - 1) * limit

	activities, totalItems, err := s.repo.List(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, errors.New("failed to fetch activities")
	}

	if activities == nil {
		activities = make([]model.Activity, 0)
	}

	return activities, totalItems, nil
}

func (s *activityService) ListCategory(ctx context.Context) ([]model.Category, error) {
	categories, err := s.categoryRepo.List(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch categories")
	}
	return categories, nil
}

func (s *activityService) GetByID(ctx context.Context, userID, activityID string) (*model.Activity, error) {
	activity, err := s.repo.GetByID(ctx, activityID, userID)
	if err != nil {
		if err.Error() == "activity not found" {
			return nil, errors.New("activity does not exist or you do not have permission")
		}
		return nil, errors.New("failed to fetch activity")
	}
	return activity, nil
}

func (s *activityService) Update(ctx context.Context, userID, activityID string, req dto.UpdateActivityRequest) (*model.Activity, error) {
	activity, err := s.GetByID(ctx, userID, activityID)
	if err != nil {
		return nil, err
	}

	categoryUUID, _ := uuid.Parse(req.CategoryID)
	activity.CategoryID = categoryUUID
	activity.Value = req.Value
	activity.Note = req.Note
	activity.OccurredAt = req.OccurredAt

	if err := s.repo.Update(ctx, activity); err != nil {
		return nil, errors.New("failed to update activity")
	}

	return activity, nil
}

func (s *activityService) Delete(ctx context.Context, userID, activityID string) error {
	err := s.repo.Delete(ctx, activityID, userID)
	if err != nil {
		if err.Error() == "activity not found" {
			return errors.New("activity does not exist or you do not have permission")
		}
		return errors.New("failed to delete activity")
	}
	return nil
}
