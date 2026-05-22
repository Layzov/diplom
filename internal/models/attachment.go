package models

import (
	"time"

	"github.com/google/uuid"
)

// Attachment is a file linked to a task.
type Attachment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TaskID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Task      Task      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FileName  string    `gorm:"size:1024;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
