package models

import (
	"time"

	"github.com/google/uuid"
)

// Session is one study run for a user (optionally scoped to a subject).
type Session struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	User      User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubjectID *uuid.UUID `gorm:"type:uuid;index"`
	Subject   *Subject   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	StartedAt time.Time  `gorm:"not null"`
	EndedAt   *time.Time
	Notes     *string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Stats        *AttemptStats  `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TaskAttempts []TaskAttempt  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Session) TableName() string {
	return "sessions"
}
