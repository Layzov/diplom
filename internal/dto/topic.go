package dto

import (
	"time"

	"diplom/internal/models"

	"github.com/google/uuid"
)

type TopicResponse struct {
	ID         uuid.UUID `json:"id"`
	SubjectID  uuid.UUID `json:"subject_id"`
	Title      string    `json:"title"`
	NumOfTasks int       `json:"num_of_tasks"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateTopicRequest struct {
	Title string `json:"title"`
}

type UpdateTopicRequest struct {
	Title *string `json:"title,omitempty"`
}

type TopicListResponse struct {
	Items []TopicResponse `json:"items"`
}

func TopicFromModel(t *models.Topic) TopicResponse {
	return TopicResponse{
		ID:         t.ID,
		SubjectID:  t.SubjectID,
		Title:      t.Title,
		NumOfTasks: t.NumOfTasks,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}
