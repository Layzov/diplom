package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserStatsResponse is aggregated progress for a user (view user_learning_stats).
type UserStatsResponse struct {
	UserID             uuid.UUID `json:"user_id" gorm:"column:user_id"`
	TotalSessions      int64     `json:"total_sessions" gorm:"column:total_sessions"`
	FinishedSessions   int64     `json:"finished_sessions" gorm:"column:finished_sessions"`
	ActiveSessions     int64     `json:"active_sessions" gorm:"column:active_sessions"`
	TotalAttempts      int64     `json:"total_attempts" gorm:"column:total_attempts"`
	CorrectCount       int64     `json:"correct_count" gorm:"column:correct_count"`
	WrongCount         int64     `json:"wrong_count" gorm:"column:wrong_count"`
	PartialCount       int64     `json:"partial_count" gorm:"column:partial_count"`
	SkippedCount       int64     `json:"skipped_count" gorm:"column:skipped_count"`
	OverallSuccessRate float64   `json:"overall_success_rate" gorm:"column:overall_success_rate"`
	PlannedRepetitions int64     `json:"planned_repetitions" gorm:"column:planned_repetitions"`
	OverdueRepetitions int64     `json:"overdue_repetitions" gorm:"column:overdue_repetitions"`
}

type TopicStatsResponse struct {
	TopicID            uuid.UUID `json:"topic_id" gorm:"column:topic_id"`
	SubjectID          uuid.UUID `json:"subject_id" gorm:"column:subject_id"`
	TopicTitle         string    `json:"topic_title" gorm:"column:topic_title"`
	SubjectTitle       string    `json:"subject_title" gorm:"column:subject_title"`
	UserID             uuid.UUID `json:"user_id" gorm:"column:user_id"`
	AttemptCount       int64     `json:"attempt_count" gorm:"column:attempt_count"`
	CorrectCount       int64     `json:"correct_count" gorm:"column:correct_count"`
	WrongCount         int64     `json:"wrong_count" gorm:"column:wrong_count"`
	PartialCount       int64     `json:"partial_count" gorm:"column:partial_count"`
	SkippedCount       int64     `json:"skipped_count" gorm:"column:skipped_count"`
	SuccessRate        float64   `json:"success_rate" gorm:"column:success_rate"`
	PlannedRepetitions int64     `json:"planned_repetitions" gorm:"column:planned_repetitions"`
}

type TopicStatsListResponse struct {
	Items []TopicStatsResponse `json:"items"`
}

type SessionStatsResponse struct {
	SessionID     uuid.UUID  `json:"session_id"`
	UserID        uuid.UUID  `json:"user_id"`
	SubjectID     *uuid.UUID `json:"subject_id,omitempty"`
	StartedAt     time.Time  `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at,omitempty"`
	TotalTasks    int        `json:"total_tasks"`
	CorrectCount  int        `json:"correct_count"`
	WrongCount    int        `json:"wrong_count"`
	PartialCount  int        `json:"partial_count"`
	SkippedCount  int        `json:"skipped_count"`
	TotalTimeMs   int        `json:"total_time_ms"`
	AverageTimeMs int        `json:"average_time_ms"`
	SuccessRate   float64    `json:"success_rate"`
	NextRepeatAt  *time.Time `json:"next_repeat_at,omitempty"`
	CalculatedAt  *time.Time `json:"calculated_at,omitempty"`
	HasStats      bool       `json:"has_stats"`
}

type SessionStatsListResponse struct {
	Items []SessionStatsResponse `json:"items"`
}

type UpcomingRepetitionResponse struct {
	RepetitionResponse
	IsOverdue bool `json:"is_overdue"`
}

type UpcomingListResponse struct {
	Items []UpcomingRepetitionResponse `json:"items"`
}
