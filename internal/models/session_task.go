package models

import (
	"time"

	"github.com/google/uuid"
)

// SessionTask links tasks to sessions (which tasks are in this session).
type SessionTask struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	SessionID uuid.UUID `gorm:"type:uuid;not null;index:idx_session_tasks_session_task,priority:1"`
	TaskID    uuid.UUID `gorm:"type:uuid;not null;index:idx_session_tasks_session_task,priority:2"`
	CreatedAt time.Time

	Session Session `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Task    Task    `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (SessionTask) TableName() string {
	return "session_tasks"
}
