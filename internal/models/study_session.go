package models

import (
	"time"

	"github.com/google/uuid"
)

// SessionMode defines the learning mode
type SessionMode string

const (
	SessionModeLearning SessionMode = "learning"
	SessionModeReview   SessionMode = "review"
)

// Session is one study run for a user (optionally scoped to a subject/topic).
type Session struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID   `gorm:"type:uuid;not null;index:idx_sessions_user_started_at,priority:1"`
	User      User        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubjectID *uuid.UUID  `gorm:"type:uuid;index"`
	Subject   *Subject    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TopicID   *uuid.UUID  `gorm:"type:uuid;index"`
	Topic     *Topic      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TaskCount int         `gorm:"not null;default:0"`
	Mode      SessionMode `gorm:"type:text;not null;default:learning"`
	StartedAt time.Time   `gorm:"not null;index:idx_sessions_user_started_at,priority:2"`
	EndedAt   *time.Time
	Notes     *string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Stats        *AttemptStats  `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TaskAttempts []TaskAttempt  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SessionTasks []SessionTask  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Session) TableName() string {
	return "sessions"
}
