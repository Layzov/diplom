package repository

import (
	"errors"
	"time"

	"diplom/internal/apperror"
	"diplom/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepetitionRepository struct {
	db *gorm.DB
}

func NewRepetitionRepository(db *gorm.DB) *RepetitionRepository {
	return &RepetitionRepository{db: db}
}

func (r *RepetitionRepository) CreateTx(tx *gorm.DB, rep *models.Repetition) error {
	return mapError(tx.Create(rep).Error)
}

func (r *RepetitionRepository) GetByID(id uuid.UUID) (*models.Repetition, error) {
	var rep models.Repetition
	err := r.db.First(&rep, "id = ?", id).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &rep, nil
}

func (r *RepetitionRepository) Update(rep *models.Repetition) error {
	return mapError(r.db.Save(rep).Error)
}

func (r *RepetitionRepository) ListByUserID(userID uuid.UUID, status *models.RepetitionStatus) ([]models.Repetition, error) {
	q := r.db.Where("user_id = ?", userID)
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var items []models.Repetition
	err := q.Order("repeat_at ASC").Find(&items).Error
	return items, mapError(err)
}

func (r *RepetitionRepository) EarliestPlannedForSessionTx(tx *gorm.DB, sessionID uuid.UUID) (*time.Time, error) {
	var rep models.Repetition
	err := tx.Where("session_id = ? AND status = ?", sessionID, models.RepetitionStatusPlanned).
		Order("repeat_at ASC").
		First(&rep).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, mapError(err)
	}
	t := rep.RepeatAt
	return &t, nil
}

func (r *RepetitionRepository) Delete(id uuid.UUID) error {
	res := r.db.Delete(&models.Repetition{}, "id = ?", id)
	if res.Error != nil {
		return mapError(res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
