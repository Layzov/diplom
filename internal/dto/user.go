package dto

import (
	"time"

	"diplom/internal/models"

	"github.com/google/uuid"
)

type UserSettings struct {
	Timezone             string `json:"timezone"`
	Language             string `json:"language"`
	Theme                string `json:"theme"`
	NotificationsEnabled bool   `json:"notifications_enabled"`
}

type UserResponse struct {
	ID        uuid.UUID    `json:"id"`
	Email     string       `json:"email"`
	Name      string       `json:"name"`
	Role      models.Role  `json:"role"`
	Settings  UserSettings `json:"settings"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CreateUserRequest struct {
	Email    string       `json:"email"`
	Name     string       `json:"name"`
	Password string       `json:"password"`
	Role     *models.Role `json:"role,omitempty"`
	Settings *UserSettings `json:"settings,omitempty"`
}

type UpdateUserRequest struct {
	Email    *string       `json:"email,omitempty"`
	Name     *string       `json:"name,omitempty"`
	Password *string       `json:"password,omitempty"`
	Role     *models.Role  `json:"role,omitempty"`
	Settings *UserSettings `json:"settings,omitempty"`
}

type UserListResponse struct {
	Items []UserResponse `json:"items"`
}

func UserFromModel(u *models.User) UserResponse {
	return UserResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
		Role:  u.Role,
		Settings: UserSettings{
			Timezone:             u.Settings.Timezone,
			Language:             u.Settings.Language,
			Theme:                u.Settings.Theme,
			NotificationsEnabled: u.Settings.NotificationsEnabled,
		},
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
