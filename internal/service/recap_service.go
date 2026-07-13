package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/admalfrizi/weekly-wrapped-be/internal/model"
	"github.com/admalfrizi/weekly-wrapped-be/internal/repository"
	"github.com/google/uuid"
)

type RecapService interface {
	GenerateRecap(ctx context.Context, userID string, startDate time.Time) (*model.WeeklyRecap, error)
	GetRecapBySlug(ctx context.Context, slug string) (*model.WeeklyRecap, error)
}

type recapService struct {
	repo         repository.RecapRepository
	dashService  DashboardService
}

func NewRecapService(repo repository.RecapRepository, dashService DashboardService) RecapService {
	return &recapService{
		repo:        repo,
		dashService: dashService,
	}
}

func generateRandomSuffix() string {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "1a2b3c" // fallback
	}
	return hex.EncodeToString(bytes)
}

func (s *recapService) GenerateRecap(ctx context.Context, userID string, startDate time.Time) (*model.WeeklyRecap, error) {
	dashboardData, err := s.dashService.GetWeeklyDashboard(ctx, userID, startDate)
	if err != nil {
		return nil, errors.New("failed to aggregate weekly stats")
	}

	statsJSON, err := json.Marshal(dashboardData)
	if err != nil {
		return nil, errors.New("failed to freeze stats snapshot")
	}

	slug := fmt.Sprintf("wk%d-%s", dashboardData.WeekNumber, generateRandomSuffix())

	narrative := fmt.Sprintf(
		"Minggu ini kamu luar biasa! Kamu mencatat total %d entri aktivitas. %s Terus pertahankan konsistensi ini di minggu depan!",
		dashboardData.TotalEntries,
		dashboardData.InsightText,
	)

	userUUID, _ := uuid.Parse(userID)
	endDate := startDate.AddDate(0, 0, 6)

	recap := &model.WeeklyRecap{
		ID:            uuid.New(),
		UserID:        userUUID,
		WeekStart:     startDate,
		WeekEnd:       endDate,
		Slug:          slug,
		StatsSnapshot: statsJSON,
		Narrative:     narrative,
		CreatedAt:     time.Now(),
	}

	if err := s.repo.Upsert(ctx, recap); err != nil {
		return nil, errors.New("failed to save weekly recap")
	}

	return recap, nil
}

func (s *recapService) GetRecapBySlug(ctx context.Context, slug string) (*model.WeeklyRecap, error) {
	recap, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("recap not found")
	}
	return recap, nil
}