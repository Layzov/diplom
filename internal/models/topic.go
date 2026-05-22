package models

import (
	"time"

	"github.com/google/uuid"
)

// Topic is a chapter inside a subject.
// NumOfTasks is denormalized; keep in sync in service layer on task CRUD.
type Topic struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	SubjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	Subject   Subject   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title     string    `gorm:"size:255;not null"`
	NumOfTasks int      `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Tasks []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
