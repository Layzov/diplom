package service

import (
	"strings"

	"diplom/internal/apperror"
	"diplom/internal/models"
)

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" || !strings.Contains(email, "@") {
		return apperror.Validation("email is required and must be valid")
	}
	return nil
}

func validateNonEmpty(field, name string) error {
	if strings.TrimSpace(field) == "" {
		return apperror.Validation(name + " is required")
	}
	return nil
}

func validateTaskType(t models.TaskType) error {
	switch t {
	case models.TaskTypeFlashcard,
		models.TaskTypeTest,
		models.TaskTypeTheory,
		models.TaskTypeFillInTheBlank,
		models.TaskTypeMatching,
		models.TaskTypeManualReview:
		return nil
	case "":
		return apperror.Validation("task type is required")
	default:
		return apperror.Validation("unknown task type")
	}
}
