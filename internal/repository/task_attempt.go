package repository

import (
	"diplom/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskAttemptRepository struct {
	db *gorm.DB
}

func NewTaskAttemptRepository(db *gorm.DB) *TaskAttemptRepository {
	return &TaskAttemptRepository{db: db}
}

func (r *TaskAttemptRepository) CreateTx(tx *gorm.DB, attempt *models.TaskAttempt) error {
	return mapError(tx.Create(attempt).Error)
}

func (r *TaskAttemptRepository) ListBySessionID(sessionID uuid.UUID) ([]models.TaskAttempt, error) {
	var items []models.TaskAttempt
	err := r.db.Where("session_id = ?", sessionID).Order("created_at ASC").Find(&items).Error
	return items, mapError(err)
}

func (r *TaskAttemptRepository) ListBySessionIDTx(tx *gorm.DB, sessionID uuid.UUID) ([]models.TaskAttempt, error) {
	var items []models.TaskAttempt
	err := tx.Where("session_id = ?", sessionID).Order("created_at ASC").Find(&items).Error
	return items, mapError(err)
}
