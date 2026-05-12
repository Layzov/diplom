package models

import (
	"time"

	"github.com/google/uuid"
)

// Card is flashcard-style content for a task (front / back).
type Card struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TaskID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Task      Task      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Front     string    `gorm:"type:text;not null"`
	Back      string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
