package service

import (
	"time"

	"diplom/internal/models"
)

// ScheduleRepeatAt returns the next review time from attempt result (simple spaced repetition).
func ScheduleRepeatAt(result models.AnswerResult, from time.Time) time.Time {
	var days int
	switch result {
	case models.AnswerResultCorrect:
		days = 3
	case models.AnswerResultPartial:
		days = 2
	case models.AnswerResultWrong:
		days = 1
	case models.AnswerResultSkipped:
		days = 1
	default:
		days = 1
	}
	return from.Add(time.Duration(days) * 24 * time.Hour)
}

// DetermineQuality maps answer result to repetition quality level.
func DetermineQuality(result models.AnswerResult) models.RepetitionQuality {
	switch result {
	case models.AnswerResultCorrect:
		return models.QualityMastered
	case models.AnswerResultPartial:
		return models.QualityLearning
	case models.AnswerResultWrong:
		return models.QualityStruggling
	case models.AnswerResultSkipped:
		return models.QualityStruggling
	default:
		return models.QualityLearning
	}
}

