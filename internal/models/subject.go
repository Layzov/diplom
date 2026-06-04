package models

import (
	"time"

	"github.com/google/uuid"
)

// Subject is a user's course or exam area (e.g. "Математика").
// NumOfTopics is denormalized; keep in sync in service layer on topic CRUD.
type Subject struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index:idx_subjects_user_created_at,priority:1"`
	User        User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title       string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	NumOfTopics int       `gorm:"not null;default:0"`
	CreatedAt   time.Time `gorm:"index:idx_subjects_user_created_at,priority:2"`
	UpdatedAt   time.Time

	Topics []Topic `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
