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
	SessionID      uuid.UUID    `gorm:"type:uuid;not null;index:idx_task_attempts_session_created_at,priority:1"`
	TaskID         uuid.UUID    `gorm:"type:uuid;not null;index:idx_task_attempts_user_task,priority:2"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null;index:idx_task_attempts_user_task,priority:1"`
	AnswerText     string       `gorm:"type:text"`
	SelectedOption string       `gorm:"type:text"`
	IsCorrect      bool         `gorm:"not null;default:false"`
	Result         AnswerResult `gorm:"type:text;not null;default:wrong"`
	ResponseTimeMs int          `gorm:"not null;default:0"`
	ReviewedAt     *time.Time
	CreatedAt      time.Time `gorm:"index:idx_task_attempts_session_created_at,priority:2"`
	UpdatedAt      time.Time

	Session Session `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Task    Task    `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (TaskAttempt) TableName() string {
	return "task_attempts"
}
