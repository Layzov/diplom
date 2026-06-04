package repository

import (
	"diplom/internal/apperror"
	"diplom/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *models.Session) error {
	return mapError(r.db.Create(session).Error)
}

func (r *SessionRepository) GetByID(id uuid.UUID) (*models.Session, error) {
	var s models.Session
	err := r.db.First(&s, "id = ?", id).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &s, nil
}

func (r *SessionRepository) ListByUserID(userID uuid.UUID) ([]models.Session, error) {
	var items []models.Session
	err := r.db.Where("user_id = ?", userID).Order("started_at DESC").Find(&items).Error
	return items, mapError(err)
}

func (r *SessionRepository) Update(session *models.Session) error {
	return mapError(r.db.Save(session).Error)
}

func (r *SessionRepository) IsActive(id uuid.UUID) (bool, error) {
	s, err := r.GetByID(id)
	if err != nil {
		return false, err
	}
	return s.EndedAt == nil, nil
}

func (r *SessionRepository) SaveTx(tx *gorm.DB, session *models.Session) error {
	return mapError(tx.Save(session).Error)
}

func (r *SessionRepository) GetByIDTx(tx *gorm.DB, id uuid.UUID) (*models.Session, error) {
	var s models.Session
	err := tx.First(&s, "id = ?", id).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &s, nil
}

func (r *SessionRepository) Exists(id uuid.UUID) (bool, error) {
	var n int64
	err := r.db.Model(&models.Session{}).Where("id = ?", id).Count(&n).Error
	return n > 0, mapError(err)
}

func (r *SessionRepository) DB() *gorm.DB {
	return r.db
}

// EnsureExists returns ErrNotFound if session missing.
func (r *SessionRepository) EnsureExists(id uuid.UUID) error {
	ok, err := r.Exists(id)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrNotFound
	}
	return nil
}
