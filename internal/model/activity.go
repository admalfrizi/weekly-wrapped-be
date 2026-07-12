package model

import (
	"time"
	"github.com/google/uuid"
)

type Activity struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	CategoryID uuid.UUID `json:"category_id"`
	Value      float64   `json:"value"`
	Note       *string   `json:"note,omitempty"`
	Category   *Category `json:"category,omitempty"`
	OccurredAt time.Time `json:"occurred_at"`
	CreatedAt  time.Time `json:"created_at"`
}