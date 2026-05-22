package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	UserRoleAdmin Role = "admin"
	UserRoleUser  Role = "user"
)

type UserSettings struct {
	Timezone             string `gorm:"size:255;not null;default:UTC"`
	Language             string `gorm:"size:255;not null;default:ru-RU"`
	Theme                string `gorm:"size:255;not null;default:light"`
	NotificationsEnabled bool   `gorm:"not null;default:true"`
}

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"size:320;uniqueIndex;not null"`
	Name         string    `gorm:"size:255;not null"`
	Role         Role      `gorm:"size:32;not null;default:user"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Settings     UserSettings `gorm:"embedded"`

	Subjects []Subject `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Sessions []Session `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
