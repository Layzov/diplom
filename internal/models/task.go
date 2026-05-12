package models

import (
	"time"

	"github.com/google/uuid"
)

// TaskKind describes how content is practiced (free-form string for extensibility).
const (
	TaskKindTheory   = "theory"
	TaskKindExercise = "exercise"
	TaskKindCard     = "card"
)

// Task is learnable content inside a topic.
type Task struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TopicID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Topic     Topic     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title     string    `gorm:"size:500;not null"`
	Prompt    *string   `gorm:"type:text"`
	TaskKind  string    `gorm:"size:32;not null;default:theory"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Attachments []Attachment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Cards       []Card       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
