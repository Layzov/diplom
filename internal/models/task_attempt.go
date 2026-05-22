package models

import (
	"time"

	"github.com/google/uuid"
)

type AnswerResult string

const (
	AnswerResultCorrect AnswerResult = "correct"
	AnswerResultPartial AnswerResult = "partial"
	AnswerResultWrong   AnswerResult = "wrong"
	AnswerResultSkipped AnswerResult = "skipped"
)

// TaskAttempt is a single answer to a task within a session.
type TaskAttempt struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey"`
	SessionID      uuid.UUID    `gorm:"type:uuid;not null;index"`
	TaskID         uuid.UUID    `gorm:"type:uuid;not null;index"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null;index"`
	AnswerText     string       `gorm:"type:text"`
	SelectedOption string       `gorm:"type:text"`
	IsCorrect      bool         `gorm:"not null;default:false"`
	Result         AnswerResult `gorm:"type:text;not null;default:wrong"`
	ResponseTimeMs int          `gorm:"not null;default:0"`
	ReviewedAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Session Session `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Task    Task    `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (TaskAttempt) TableName() string {
	return "task_attempts"
}
