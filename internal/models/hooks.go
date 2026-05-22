package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func assignUUID(id *uuid.UUID) error {
	if id == nil {
		return nil
	}
	if *id == uuid.Nil {
		*id = uuid.New()
	}
	return nil
}

// BeforeCreate hooks ensure primary keys are set (UUID v4 from Go).

func (u *User) BeforeCreate(_ *gorm.DB) error {
	return assignUUID(&u.ID)
}

func (s *Subject) BeforeCreate(_ *gorm.DB) error {
	return assignUUID(&s.ID)
}

func (t *Topic) BeforeCreate(_ *gorm.DB) error {
	return assignUUID(&t.ID)
}

func (t *Task) BeforeCreate(_ *gorm.DB) error {
	return assignUUID(&t.ID)
}

func (a *Attachment) BeforeCreate(_ *gorm.DB) error {
	return assignUUID(&a.ID)
}

func (s *Session) BeforeCreate(_ *gorm.DB) error {
	return assignUUID(&s.ID)
}

func (t *TaskAttempt) BeforeCreate(_ *gorm.DB) error {
	return assignUUID(&t.ID)
}

func (a *AttemptStats) BeforeCreate(_ *gorm.DB) error {
	return assignUUID(&a.ID)
}

func (r *Repetition) BeforeCreate(_ *gorm.DB) error {
	return assignUUID(&r.ID)
}
