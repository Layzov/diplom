package dto

import (
	"time"

	"diplom/internal/models"

	"github.com/google/uuid"
)

type TaskResponse struct {
	ID        uuid.UUID       `json:"id"`
	TopicID   uuid.UUID       `json:"topic_id"`
	Title     string          `json:"title"`
	Content   string          `json:"content"`
	Type      models.TaskType `json:"type"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title   string          `json:"title"`
	Content string          `json:"content"`
	Type    models.TaskType `json:"type,omitempty"`
}

type UpdateTaskRequest struct {
	Title   *string          `json:"title,omitempty"`
	Content *string          `json:"content,omitempty"`
	Type    *models.TaskType `json:"type,omitempty"`
}

type TaskListResponse struct {
	Items []TaskResponse `json:"items"`
}

func TaskFromModel(t *models.Task) TaskResponse {
	return TaskResponse{
		ID:        t.ID,
		TopicID:   t.TopicID,
		Title:     t.Title,
		Content:   t.Content,
		Type:      t.Type,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
