package dto

import (
	"time"

	"diplom/internal/models"

	"github.com/google/uuid"
)

type SubjectResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	NumOfTopics int       `json:"num_of_topics"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSubjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type UpdateSubjectRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type SubjectListResponse struct {
	Items []SubjectResponse `json:"items"`
}

func SubjectFromModel(s *models.Subject) SubjectResponse {
	return SubjectResponse{
		ID:          s.ID,
		UserID:      s.UserID,
		Title:       s.Title,
		Description: s.Description,
		NumOfTopics: s.NumOfTopics,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
