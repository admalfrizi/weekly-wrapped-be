package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type WeeklyRecap struct {
	ID            uuid.UUID       `json:"id"`
	UserID        uuid.UUID       `json:"user_id"`
	WeekStart     time.Time       `json:"week_start"`
	WeekEnd       time.Time       `json:"week_end"`
	Slug          string          `json:"slug"`
	StatsSnapshot json.RawMessage `json:"stats_snapshot"` // Native JSONB support
	Narrative     string          `json:"narrative"`
	CreatedAt     time.Time       `json:"created_at"`
}