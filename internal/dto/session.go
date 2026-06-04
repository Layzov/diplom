package dto

import (
	"time"

	"diplom/internal/models"

	"github.com/google/uuid"
)

type SessionResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	SubjectID *uuid.UUID `json:"subject_id,omitempty"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Notes     *string    `json:"notes,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateSessionRequest struct {
	SubjectID *uuid.UUID `json:"subject_id,omitempty"`
	Notes     *string    `json:"notes,omitempty"`
}

type FinishSessionRequest struct {
	Notes *string `json:"notes,omitempty"`
}

type SessionListResponse struct {
	Items []SessionResponse `json:"items"`
}

func SessionFromModel(s *models.Session) SessionResponse {
	return SessionResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		SubjectID: s.SubjectID,
		StartedAt: s.StartedAt,
		EndedAt:   s.EndedAt,
		Notes:     s.Notes,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
