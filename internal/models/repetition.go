package models

import (
	"time"

	"github.com/google/uuid"
)

type RepetitionStatus string

const (
	RepetitionStatusPlanned RepetitionStatus = "planned"
	RepetitionStatusDone    RepetitionStatus = "done"
	RepetitionStatusSkipped RepetitionStatus = "skipped"
)

type Repetition struct {
	ID            uuid.UUID        `gorm:"type:uuid;primaryKey"`
	SessionID     uuid.UUID        `gorm:"type:uuid;not null;index"`
	TaskAttemptID uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex"`
	TaskID        uuid.UUID        `gorm:"type:uuid;not null;index"`
	UserID        uuid.UUID        `gorm:"type:uuid;not null;index:idx_repetitions_user_repeat_at,priority:1;index:idx_repetitions_user_status_due,priority:1;index:idx_repetitions_user_topic,priority:1"`
	TopicID       uuid.UUID        `gorm:"type:uuid;not null;index:idx_repetitions_user_topic,priority:2"`
	RepeatAt      time.Time        `gorm:"not null;index:idx_repetitions_user_repeat_at,priority:2;index:idx_repetitions_user_status_due,priority:3"`
	Status        RepetitionStatus `gorm:"type:text;not null;default:planned;index:idx_repetitions_user_status_due,priority:2"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Session     Session     `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TaskAttempt TaskAttempt `gorm:"foreignKey:TaskAttemptID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Task        Task        `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User        User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Topic       Topic       `gorm:"foreignKey:TopicID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Repetition) TableName() string {
	return "repetitions"
}
