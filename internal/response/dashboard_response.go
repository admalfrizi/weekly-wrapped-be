package response

type StatCard struct {
	CategoryName   string  `json:"category_name"`
	CurrentTotal   float64 `json:"current_total"`
	Unit           string  `json:"unit"`
	TrendPercent   float64 `json:"trend_percent"`
}

// DailyActivity represents the line chart data (Aktivitas per hari)
type DailyActivity struct {
	Date         string  `json:"date"`
	DayName      string  `json:"day_name"`
	CategoryName string  `json:"category_name"`
	TotalValue   float64 `json:"total_value"`
}

// CategoryComposition represents the progress bars (Komposisi kategori)
type CategoryComposition struct {
	CategoryName string  `json:"category_name"`
	Percentage   int     `json:"percentage"` // e.g., 42, 28
}

// DashboardWeeklyResponse is the master payload for the whole page
type DashboardWeeklyResponse struct {
	WeekNumber   int                   `json:"week_number"`
	TotalEntries int                   `json:"total_entries"`
	Cards        []StatCard            `json:"cards"`
	ChartData    []DailyActivity       `json:"chart_data"`
	Compositions []CategoryComposition `json:"compositions"`
	InsightText  string                `json:"insight_text"`
}