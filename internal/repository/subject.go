package repository

import (
	"diplom/internal/apperror"
	"diplom/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

func (r *SubjectRepository) Create(subject *models.Subject) error {
	return mapError(r.db.Create(subject).Error)
}

func (r *SubjectRepository) ListByUserID(userID uuid.UUID) ([]models.Subject, error) {
	var items []models.Subject
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error
	return items, mapError(err)
}

func (r *SubjectRepository) GetByID(id uuid.UUID) (*models.Subject, error) {
	var s models.Subject
	err := r.db.First(&s, "id = ?", id).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &s, nil
}

func (r *SubjectRepository) Update(subject *models.Subject) error {
	return mapError(r.db.Save(subject).Error)
}

func (r *SubjectRepository) Delete(id uuid.UUID) error {
	res := r.db.Delete(&models.Subject{}, "id = ?", id)
	if res.Error != nil {
		return mapError(res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *SubjectRepository) AdjustTopicCount(tx *gorm.DB, subjectID uuid.UUID, delta int) error {
	res := tx.Model(&models.Subject{}).
		Where("id = ?", subjectID).
		UpdateColumn("num_of_topics", gorm.Expr("num_of_topics + ?", delta))
	if res.Error != nil {
		return mapError(res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
