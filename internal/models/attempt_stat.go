package models

import (
	"time"

	"github.com/google/uuid"
)

// AttemptStats is aggregated metrics for one session (1:1 with Session).
type AttemptStats struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	SessionID     uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex"`
	TotalTasks    int        `gorm:"not null;default:0"`
	CorrectCount  int        `gorm:"not null;default:0"`
	WrongCount    int        `gorm:"not null;default:0"`
	SkippedCount  int        `gorm:"not null;default:0"`
	PartialCount  int        `gorm:"not null;default:0"`
	TotalTimeMs   int        `gorm:"not null;default:0"`
	AverageTimeMs int        `gorm:"not null;default:0"`
	SuccessRate   float64    `gorm:"type:numeric(5,2);not null;default:0"`
	NextRepeatAt  *time.Time
	CalculatedAt  time.Time `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Session Session `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (AttemptStats) TableName() string {
	return "attempt_stats"
}
