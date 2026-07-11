package model

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	Email         string
	Username      string
	Name          string
	ProfileImgURL *string
	PasswordHash  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}