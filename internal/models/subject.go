package models

import (
	"time"

	"github.com/google/uuid"
)

// Subject groups topics for one user (e.g. "Математика").
type Subject struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	User        User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name        string    `gorm:"size:255;not null"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Topics []Topic `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
