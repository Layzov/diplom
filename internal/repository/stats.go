package repository

import (
	"time"

	"diplom/internal/apperror"
	"diplom/internal/dto"
	"diplom/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StatsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetUserStats(userID uuid.UUID) (*dto.UserStatsResponse, error) {
	var row dto.UserStatsResponse
	err := r.db.Table("user_learning_stats").
		Where("user_id = ?", userID).
		First(&row).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &row, nil
}

func (r *StatsRepository) ListTopicStats(userID uuid.UUID) ([]dto.TopicStatsResponse, error) {
	var items []dto.TopicStatsResponse
	err := r.db.Table("topic_learning_stats").
		Where("user_id = ?", userID).
		Order("success_rate DESC, attempt_count DESC").
		Find(&items).Error
	return items, mapError(err)
}

func (r *StatsRepository) ListSessionStats(userID uuid.UUID) ([]dto.SessionStatsResponse, error) {
	type row struct {
		SessionID     uuid.UUID  `gorm:"column:session_id"`
		UserID        uuid.UUID  `gorm:"column:user_id"`
		SubjectID     *uuid.UUID `gorm:"column:subject_id"`
		StartedAt     time.Time  `gorm:"column:started_at"`
		EndedAt       *time.Time `gorm:"column:ended_at"`
		StatsID       *uuid.UUID `gorm:"column:stats_id"`
		TotalTasks    *int       `gorm:"column:total_tasks"`
		CorrectCount  *int       `gorm:"column:correct_count"`
		WrongCount    *int       `gorm:"column:wrong_count"`
		PartialCount  *int       `gorm:"column:partial_count"`
		SkippedCount  *int       `gorm:"column:skipped_count"`
		TotalTimeMs   *int       `gorm:"column:total_time_ms"`
		AverageTimeMs *int       `gorm:"column:average_time_ms"`
		SuccessRate   *float64   `gorm:"column:success_rate"`
		NextRepeatAt  *time.Time `gorm:"column:next_repeat_at"`
		CalculatedAt  *time.Time `gorm:"column:calculated_at"`
	}

	var rows []row
	err := r.db.Table("session_learning_stats").
		Where("user_id = ?", userID).
		Order("started_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, mapError(err)
	}

	out := make([]dto.SessionStatsResponse, 0, len(rows))
	for _, rw := range rows {
		item := dto.SessionStatsResponse{
			SessionID: rw.SessionID,
			UserID:    rw.UserID,
			SubjectID: rw.SubjectID,
			StartedAt: rw.StartedAt,
			EndedAt:   rw.EndedAt,
			HasStats:  rw.StatsID != nil,
		}
		if rw.StatsID != nil {
			item.TotalTasks = derefInt(rw.TotalTasks)
			item.CorrectCount = derefInt(rw.CorrectCount)
			item.WrongCount = derefInt(rw.WrongCount)
			item.PartialCount = derefInt(rw.PartialCount)
			item.SkippedCount = derefInt(rw.SkippedCount)
			item.TotalTimeMs = derefInt(rw.TotalTimeMs)
			item.AverageTimeMs = derefInt(rw.AverageTimeMs)
			item.SuccessRate = derefFloat(rw.SuccessRate)
			item.NextRepeatAt = rw.NextRepeatAt
			item.CalculatedAt = rw.CalculatedAt
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *StatsRepository) ListUpcoming(userID uuid.UUID, limit int, includeOverdue bool) ([]dto.UpcomingRepetitionResponse, error) {
	now := time.Now().UTC()
	q := r.db.Table("calendar").Where("user_id = ?", userID)
	if !includeOverdue {
		q = q.Where("due_at >= ?", now)
	}
	var cal []dto.CalendarEntry
	err := q.Order("due_at ASC").Limit(limit).Find(&cal).Error
	if err != nil {
		return nil, mapError(err)
	}

	out := make([]dto.UpcomingRepetitionResponse, 0, len(cal))
	for _, c := range cal {
		rep := dto.RepetitionResponse{
			ID:            c.ID,
			SessionID:     c.SessionID,
			TaskAttemptID: c.TaskAttemptID,
			TaskID:        c.TaskID,
			UserID:        c.UserID,
			TopicID:       c.TopicID,
			RepeatAt:      c.DueAt,
			Status:        c.Status,
			CreatedAt:     c.CreatedAt,
			UpdatedAt:     c.UpdatedAt,
		}
		out = append(out, dto.UpcomingRepetitionResponse{
			RepetitionResponse: rep,
			IsOverdue:          c.DueAt.Before(now),
		})
	}
	return out, nil
}

func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func derefFloat(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

// EnsureUserStatsRow verifies user exists before reading stats view.
func (r *StatsRepository) EnsureUserExists(userID uuid.UUID) error {
	var n int64
	err := r.db.Model(&models.User{}).Where("id = ?", userID).Count(&n).Error
	if err != nil {
		return mapError(err)
	}
	if n == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
