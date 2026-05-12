package models

import (
	"time"

	"github.com/google/uuid"
)

// Attachment is a file linked to a task (path or object key in storage).
type Attachment struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TaskID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Task         Task      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StorageKey   string    `gorm:"size:1024;not null"`
	ContentType  *string `gorm:"size:128"`
	OriginalName *string `gorm:"size:512"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
