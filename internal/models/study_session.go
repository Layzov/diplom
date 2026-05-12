package models

import (
	"time"

	"github.com/google/uuid"
)

// StudySession is one study run for a user (may focus on a subject).
type StudySession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubjectID *uuid.UUID `gorm:"type:uuid;index"`
	Subject   *Subject  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	StartedAt time.Time `gorm:"not null"`
	EndedAt   *time.Time
	Notes     *string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Attempts []AttemptStat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName keeps "sessions" out of reserved / ambiguous naming in SQL.
func (StudySession) TableName() string {
	return "study_sessions"
}
