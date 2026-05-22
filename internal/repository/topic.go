package repository

import (
	"diplom/internal/apperror"
	"diplom/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

func (r *TopicRepository) Create(topic *models.Topic) error {
	return r.CreateTx(r.db, topic)
}

func (r *TopicRepository) CreateTx(tx *gorm.DB, topic *models.Topic) error {
	return mapError(tx.Create(topic).Error)
}

func (r *TopicRepository) ListBySubjectID(subjectID uuid.UUID) ([]models.Topic, error) {
	var items []models.Topic
	err := r.db.Where("subject_id = ?", subjectID).Order("created_at DESC").Find(&items).Error
	return items, mapError(err)
}

func (r *TopicRepository) GetByID(id uuid.UUID) (*models.Topic, error) {
	var t models.Topic
	err := r.db.First(&t, "id = ?", id).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &t, nil
}

func (r *TopicRepository) Update(topic *models.Topic) error {
	return mapError(r.db.Save(topic).Error)
}

func (r *TopicRepository) Delete(id uuid.UUID) error {
	return r.DeleteTx(r.db, id)
}

func (r *TopicRepository) DeleteTx(tx *gorm.DB, id uuid.UUID) error {
	res := tx.Delete(&models.Topic{}, "id = ?", id)
	if res.Error != nil {
		return mapError(res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *TopicRepository) AdjustTaskCount(tx *gorm.DB, topicID uuid.UUID, delta int) error {
	res := tx.Model(&models.Topic{}).
		Where("id = ?", topicID).
		UpdateColumn("num_of_tasks", gorm.Expr("num_of_tasks + ?", delta))
	if res.Error != nil {
		return mapError(res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
