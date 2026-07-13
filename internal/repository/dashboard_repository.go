package repository

import (
	"context"
	"time"
)

type DashboardRepository interface {
	GetCategoryTotals(ctx context.Context, userID string, startDate, endDate time.Time) ([]CategoryTotalRow, int, error)
	GetDailyTotals(ctx context.Context, userID string, startDate, endDate time.Time) ([]DailyTotalRow, error)
}

type dashboardRepository struct {
	*BaseRepository
}

func NewDashboardRepository(base *BaseRepository) DashboardRepository {
	return &dashboardRepository{BaseRepository: base}
}

type CategoryTotalRow struct {
	CategoryID   string
	CategoryName string
	Icon         *string
	ColorHex     *string
	TotalValue   float64
}

func (r *dashboardRepository) GetCategoryTotals(ctx context.Context, userID string, startDate, endDate time.Time) ([]CategoryTotalRow, int, error) {
	query := `
		SELECT 
			c.id, 
			c.name, 
			SUM(a.value) as total_value, 
			SUM(COUNT(a.id)) OVER() as total_entries
		FROM activities a
		JOIN categories c ON a.category_id = c.id
		WHERE a.user_id = $1 AND a.occurred_at >= $2 AND a.occurred_at <= $3
		GROUP BY c.id, c.name;
	`

	rows, err := r.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []CategoryTotalRow
	var totalEntries int

	for rows.Next() {
		var row CategoryTotalRow
		err := rows.Scan(&row.CategoryID, &row.CategoryName, &row.TotalValue, &totalEntries)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, row)
	}
	return results, totalEntries, nil
}

type DailyTotalRow struct {
	Date         time.Time
	CategoryName string
	TotalValue   float64
}

func (r *dashboardRepository) GetDailyTotals(ctx context.Context, userID string, startDate, endDate time.Time) ([]DailyTotalRow, error) {
	query := `
		SELECT 
			DATE_TRUNC('day', a.occurred_at) as activity_date,
			c.name,
			SUM(a.value)
		FROM activities a
		JOIN categories c ON a.category_id = c.id
		WHERE a.user_id = $1 AND a.occurred_at >= $2 AND a.occurred_at <= $3
		GROUP BY activity_date, c.name
		ORDER BY activity_date ASC;
	`

	rows, err := r.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []DailyTotalRow
	for rows.Next() {
		var row DailyTotalRow
		if err := rows.Scan(&row.Date, &row.CategoryName, &row.TotalValue); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}