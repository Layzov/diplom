package dto

import "time"

type RegisterRequest struct {
	Email    string        `json:"email"`
	Name     string        `json:"name"`
	Password string        `json:"password"`
	Settings *UserSettings `json:"settings,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	ExpiresAt   time.Time    `json:"expires_at"`
	User        UserResponse `json:"user"`
}
