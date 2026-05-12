package models

import (
	"time"

	"github.com/google/uuid"
)

// User is an account that owns subjects and study data.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"size:320;uniqueIndex;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Subjects []Subject `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
