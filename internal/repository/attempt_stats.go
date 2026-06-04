package repository

import (
	"errors"

	"diplom/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttemptStatsRepository struct {
	db *gorm.DB
}

func NewAttemptStatsRepository(db *gorm.DB) *AttemptStatsRepository {
	return &AttemptStatsRepository{db: db}
}

func (r *AttemptStatsRepository) GetBySessionID(sessionID uuid.UUID) (*models.AttemptStats, error) {
	var s models.AttemptStats
	err := r.db.First(&s, "session_id = ?", sessionID).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &s, nil
}

func (r *AttemptStatsRepository) UpsertTx(tx *gorm.DB, stats *models.AttemptStats) error {
	var existing models.AttemptStats
	err := tx.Where("session_id = ?", stats.SessionID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return mapError(tx.Create(stats).Error)
	}
	if err != nil {
		return mapError(err)
	}
	stats.ID = existing.ID
	return mapError(tx.Save(stats).Error)
}
