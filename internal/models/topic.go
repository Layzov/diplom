package models

import (
	"time"

	"github.com/google/uuid"
)

// Topic is a chapter or block inside a subject.
type Topic struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	SubjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	Subject   Subject   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name      string    `gorm:"size:255;not null"`
	SortOrder *int
	CreatedAt time.Time
	UpdatedAt time.Time

	Tasks []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
