package service

import (
	"math"
	"time"

	"diplom/internal/models"

	"github.com/google/uuid"
)

func buildAttemptStats(sessionID uuid.UUID, attempts []models.TaskAttempt, nextRepeat *time.Time) models.AttemptStats {
	now := time.Now().UTC()
	stats := models.AttemptStats{
		SessionID:    sessionID,
		CalculatedAt: now,
		CreatedAt:    now,
		UpdatedAt:    now,
		NextRepeatAt: nextRepeat,
	}

	stats.TotalTasks = len(attempts)
	if stats.TotalTasks == 0 {
		return stats
	}

	var totalTime int
	for _, a := range attempts {
		switch a.Result {
		case models.AnswerResultCorrect:
			stats.CorrectCount++
		case models.AnswerResultWrong:
			stats.WrongCount++
		case models.AnswerResultSkipped:
			stats.SkippedCount++
		case models.AnswerResultPartial:
			stats.PartialCount++
		}
		totalTime += a.ResponseTimeMs
	}

	stats.TotalTimeMs = totalTime
	stats.AverageTimeMs = totalTime / stats.TotalTasks
	stats.SuccessRate = math.Round(float64(stats.CorrectCount)/float64(stats.TotalTasks)*10000) / 100

	return stats
}
