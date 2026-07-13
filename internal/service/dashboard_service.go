package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/admalfrizi/weekly-wrapped-be/internal/repository"
	"github.com/admalfrizi/weekly-wrapped-be/internal/response"
)

type DashboardService interface {
	GetWeeklyDashboard(ctx context.Context, userID string, startDate time.Time) (*response.DashboardWeeklyResponse, error)
}

type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetWeeklyDashboard(ctx context.Context, userID string, startDate time.Time) (*response.DashboardWeeklyResponse, error) {
	endDate := startDate.AddDate(0, 0, 6).Add(23 * time.Hour).Add(59 * time.Minute)
	prevStartDate := startDate.AddDate(0, 0, -7)
	prevEndDate := endDate.AddDate(0, 0, -7)
	_, weekNumber := startDate.ISOWeek()

	currentTotals, totalEntries, err := s.repo.GetCategoryTotals(ctx, userID, startDate, endDate)
	if err != nil { return nil, err }
	
	prevTotals, _, err := s.repo.GetCategoryTotals(ctx, userID, prevStartDate, prevEndDate)
	if err != nil { return nil, err }

	chartRows, err := s.repo.GetDailyTotals(ctx, userID, startDate, endDate)
	if err != nil { return nil, err }

	var cards []response.StatCard
	var compositions []response.CategoryComposition
	var totalAllValues float64
	var topCategory string
	var topCategoryTotal float64

	prevMap := make(map[string]float64)
	for _, p := range prevTotals { prevMap[p.CategoryName] = p.TotalValue }

	for _, curr := range currentTotals {
		totalAllValues += curr.TotalValue
		if curr.TotalValue > topCategoryTotal {
			topCategoryTotal = curr.TotalValue
			topCategory = curr.CategoryName
		}

		prevVal := prevMap[curr.CategoryName]
		var trend float64
		if prevVal > 0 {
			trend = ((curr.TotalValue - prevVal) / prevVal) * 100
		} else if curr.TotalValue > 0 {
			trend = 100
		}

		unit := "items"
		if curr.CategoryName == "Workout" { unit = "min" }
		if curr.CategoryName == "Reading" { unit = "pages" }
		if curr.CategoryName == "Coding" { unit = "hours" }

		cards = append(cards, response.StatCard{
			CategoryName: curr.CategoryName,
			CurrentTotal: curr.TotalValue,
			Unit:         unit,
			TrendPercent: math.Round(trend*10) / 10,
		})
	}

	for _, curr := range currentTotals {
		pct := 0
		if totalAllValues > 0 {
			pct = int(math.Round((curr.TotalValue / totalAllValues) * 100))
		}
		compositions = append(compositions, response.CategoryComposition{
			CategoryName: curr.CategoryName,
			Percentage:   pct,
		})
	}

	dayNames := map[time.Weekday]string{
		time.Monday: "Sen", time.Tuesday: "Sel", time.Wednesday: "Rab",
		time.Thursday: "Kam", time.Friday: "Jum", time.Saturday: "Sab", time.Sunday: "Min",
	}
	var chartData []response.DailyActivity
	for _, row := range chartRows {
		chartData = append(chartData, response.DailyActivity{
			Date:         row.Date.Format("2006-01-02"),
			DayName:      dayNames[row.Date.Weekday()],
			CategoryName: row.CategoryName,
			TotalValue:   row.TotalValue,
		})
	}

	insight := "Belum ada data cukup minggu ini."
	if topCategory != "" {
		insight = fmt.Sprintf("%s mendominasi minggu ini — mewakili porsi besar dari total aktivitas.", topCategory)
	}

	return &response.DashboardWeeklyResponse{
		WeekNumber:   weekNumber,
		TotalEntries: totalEntries,
		Cards:        cards,
		ChartData:    chartData,
		Compositions: compositions,
		InsightText:  insight,
	}, nil
}