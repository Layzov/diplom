package dto

import (
	"time"

	"diplom/internal/models"

	"github.com/google/uuid"
)

type TaskAttemptResponse struct {
	ID             uuid.UUID           `json:"id"`
	SessionID      uuid.UUID           `json:"session_id"`
	TaskID         uuid.UUID           `json:"task_id"`
	UserID         uuid.UUID           `json:"user_id"`
	AnswerText     string              `json:"answer_text,omitempty"`
	SelectedOption string              `json:"selected_option,omitempty"`
	IsCorrect      bool                `json:"is_correct"`
	Result         models.AnswerResult `json:"result"`
	ResponseTimeMs int                 `json:"response_time_ms"`
	ReviewedAt     *time.Time          `json:"reviewed_at,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

type CreateTaskAttemptRequest struct {
	TaskID         uuid.UUID           `json:"task_id"`
	AnswerText     string              `json:"answer_text,omitempty"`
	SelectedOption string              `json:"selected_option,omitempty"`
	IsCorrect      bool                `json:"is_correct"`
	Result         models.AnswerResult `json:"result"`
	ResponseTimeMs int                 `json:"response_time_ms"`
}

type SubmitAttemptResponse struct {
	Attempt    TaskAttemptResponse `json:"attempt"`
	Repetition RepetitionResponse  `json:"repetition"`
}

type TaskAttemptListResponse struct {
	Items []TaskAttemptResponse `json:"items"`
}

func TaskAttemptFromModel(a *models.TaskAttempt) TaskAttemptResponse {
	return TaskAttemptResponse{
		ID:             a.ID,
		SessionID:      a.SessionID,
		TaskID:         a.TaskID,
		UserID:         a.UserID,
		AnswerText:     a.AnswerText,
		SelectedOption: a.SelectedOption,
		IsCorrect:      a.IsCorrect,
		Result:         a.Result,
		ResponseTimeMs: a.ResponseTimeMs,
		ReviewedAt:     a.ReviewedAt,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

type AttemptStatsResponse struct {
	ID            uuid.UUID  `json:"id"`
	SessionID     uuid.UUID  `json:"session_id"`
	TotalTasks    int        `json:"total_tasks"`
	CorrectCount  int        `json:"correct_count"`
	WrongCount    int        `json:"wrong_count"`
	SkippedCount  int        `json:"skipped_count"`
	PartialCount  int        `json:"partial_count"`
	TotalTimeMs   int        `json:"total_time_ms"`
	AverageTimeMs int        `json:"average_time_ms"`
	SuccessRate   float64    `json:"success_rate"`
	NextRepeatAt  *time.Time `json:"next_repeat_at,omitempty"`
	CalculatedAt  time.Time  `json:"calculated_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func AttemptStatsFromModel(s *models.AttemptStats) AttemptStatsResponse {
	return AttemptStatsResponse{
		ID:            s.ID,
		SessionID:     s.SessionID,
		TotalTasks:    s.TotalTasks,
		CorrectCount:  s.CorrectCount,
		WrongCount:    s.WrongCount,
		SkippedCount:  s.SkippedCount,
		PartialCount:  s.PartialCount,
		TotalTimeMs:   s.TotalTimeMs,
		AverageTimeMs: s.AverageTimeMs,
		SuccessRate:   s.SuccessRate,
		NextRepeatAt:  s.NextRepeatAt,
		CalculatedAt:  s.CalculatedAt,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
