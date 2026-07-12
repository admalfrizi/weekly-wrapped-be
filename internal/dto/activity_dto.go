package dto

import "time"

type CreateActivityRequest struct {
	CategoryID string    `json:"category_id" binding:"required,uuid"`
	Value      float64   `json:"value" binding:"required,gt=0"`
	Note       *string   `json:"note"`
	OccurredAt time.Time `json:"occurred_at" binding:"required"`
}

type UpdateActivityRequest struct {
	CategoryID string    `json:"category_id" binding:"required,uuid"`
	Value      float64   `json:"value" binding:"required,gt=0"`
	Note       *string   `json:"note"`
	OccurredAt time.Time `json:"occurred_at" binding:"required"`
}