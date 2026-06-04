package dto

import (
	"time"

	"diplom/internal/models"

	"github.com/google/uuid"
)

type RepetitionResponse struct {
	ID            uuid.UUID               `json:"id"`
	SessionID     uuid.UUID               `json:"session_id"`
	TaskAttemptID uuid.UUID               `json:"task_attempt_id"`
	TaskID        uuid.UUID               `json:"task_id"`
	UserID        uuid.UUID               `json:"user_id"`
	TopicID       uuid.UUID               `json:"topic_id"`
	RepeatAt      time.Time               `json:"repeat_at"`
	Status        models.RepetitionStatus `json:"status"`
	Quality       models.RepetitionQuality `json:"quality"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

type UpdateRepetitionRequest struct {
	Status models.RepetitionStatus `json:"status"`
}

type RepetitionListResponse struct {
	Items []RepetitionResponse `json:"items"`
}

func RepetitionFromModel(r *models.Repetition) RepetitionResponse {
	return RepetitionResponse{
		ID:            r.ID,
		SessionID:     r.SessionID,
		TaskAttemptID: r.TaskAttemptID,
		TaskID:        r.TaskID,
		UserID:        r.UserID,
		TopicID:       r.TopicID,
		RepeatAt:      r.RepeatAt,
		Status:        r.Status,
		Quality:       r.Quality,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

type CalendarEntry struct {
	ID            uuid.UUID               `json:"id" gorm:"column:id"`
	UserID        uuid.UUID               `json:"user_id" gorm:"column:user_id"`
	SessionID     uuid.UUID               `json:"session_id" gorm:"column:session_id"`
	TaskAttemptID uuid.UUID               `json:"task_attempt_id" gorm:"column:task_attempt_id"`
	TaskID        uuid.UUID               `json:"task_id" gorm:"column:task_id"`
	TopicID       uuid.UUID               `json:"topic_id" gorm:"column:topic_id"`
	DueAt         time.Time               `json:"due_at" gorm:"column:due_at"`
	CalendarDate  string                  `json:"calendar_date" gorm:"column:calendar_date"`
	Status        models.RepetitionStatus `json:"status" gorm:"column:status"`
	CreatedAt     time.Time               `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time               `json:"updated_at" gorm:"column:updated_at"`
}

type CalendarListResponse struct {
	Items []CalendarEntry `json:"items"`
}
