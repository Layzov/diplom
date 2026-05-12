package models

import (
	"time"

	"github.com/google/uuid"
)

// AttemptStat records one answer during a session (input for spaced repetition).
type AttemptStat struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	SessionID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Session    StudySession `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	User       User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TaskID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Task       Task      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CardID     *uuid.UUID `gorm:"type:uuid;index"`
	Card       *Card     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	IsCorrect  bool      `gorm:"not null"`
	Quality    *uint8    // SM-2 style grade 0–5, optional until scheduler uses it
	AnswerMs   *int
	AnsweredAt time.Time `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName maps to attempt_stats.
func (AttemptStat) TableName() string {
	return "attempt_stats"
}
