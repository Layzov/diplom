package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskType string

const (
	TaskTypeFlashcard      TaskType = "flashcard"
	TaskTypeTest           TaskType = "test"
	TaskTypeTheory         TaskType = "short_answer"
	TaskTypeFillInTheBlank TaskType = "fill_in_the_blank"
	TaskTypeMatching       TaskType = "matching"
	TaskTypeManualReview   TaskType = "manual_review"
)

type Task struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TopicID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Topic     Topic     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title     string    `gorm:"size:500;not null"`
	Content   string    `gorm:"type:text;not null"`
	Type      TaskType  `gorm:"size:32;not null;default:flashcard"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Attachments []Attachment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
