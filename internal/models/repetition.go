package models

import (
	"time"

	"github.com/google/uuid"
)

// Repetition lifecycle for scheduling UI and queries.
const (
	RepetitionScheduled  = "scheduled"
	RepetitionCompleted  = "completed"
	RepetitionSnoozed    = "snoozed"
	RepetitionCancelled  = "cancelled"
)

// Repetition is the next review obligation for a user on a task/card.
type Repetition struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	User         User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TaskID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Task         Task      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CardID       *uuid.UUID `gorm:"type:uuid;index"`
	Card         *Card     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	NextDueAt    time.Time `gorm:"not null;index"`
	IntervalDays float64   `gorm:"not null;default:0"`
	EaseFactor   float64   `gorm:"not null;default:2.5"`
	ReviewCount  int       `gorm:"not null;default:0"`
	Status       string    `gorm:"size:32;not null;default:scheduled;index"`
	LastReviewAt *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
