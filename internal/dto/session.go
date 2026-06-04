package dto

import (
	"time"

	"diplom/internal/models"

	"github.com/google/uuid"
)

type SessionResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	SubjectID *uuid.UUID `json:"subject_id,omitempty"`
	TopicID   *uuid.UUID `json:"topic_id,omitempty"`
	TaskCount int       `json:"task_count"`
	Mode      string    `json:"mode"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSessionRequest struct {
	SubjectID *uuid.UUID `json:"subject_id,omitempty"`
	TopicID   *uuid.UUID `json:"topic_id,omitempty"`
	TaskCount int        `json:"task_count,omitempty"` // Default 15 if omitted
	Mode      string     `json:"mode,omitempty"`       // Default "learning" if omitted
	Notes     *string    `json:"notes,omitempty"`
}

type FinishSessionRequest struct {
	Notes *string `json:"notes,omitempty"`
}

type SessionListResponse struct {
	Items []SessionResponse `json:"items"`
}

type SessionTasksResponse struct {
	SessionID uuid.UUID   `json:"session_id"`
	Tasks     []TaskResponse `json:"tasks"`
}

func SessionFromModel(s *models.Session) SessionResponse {
	return SessionResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		SubjectID: s.SubjectID,
		TopicID:   s.TopicID,
		TaskCount: s.TaskCount,
		Mode:      string(s.Mode),
		StartedAt: s.StartedAt,
		EndedAt:   s.EndedAt,
		Notes:     s.Notes,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
