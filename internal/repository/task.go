package repository

import (
	"diplom/internal/apperror"
	"diplom/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	return r.CreateTx(r.db, task)
}

func (r *TaskRepository) CreateTx(tx *gorm.DB, task *models.Task) error {
	return mapError(tx.Create(task).Error)
}

func (r *TaskRepository) ListByTopicID(topicID uuid.UUID) ([]models.Task, error) {
	var items []models.Task
	err := r.db.Where("topic_id = ?", topicID).Order("created_at DESC").Find(&items).Error
	return items, mapError(err)
}

func (r *TaskRepository) GetByID(id uuid.UUID) (*models.Task, error) {
	var t models.Task
	err := r.db.First(&t, "id = ?", id).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &t, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	return mapError(r.db.Save(task).Error)
}

func (r *TaskRepository) Delete(id uuid.UUID) error {
	return r.DeleteTx(r.db, id)
}

func (r *TaskRepository) DeleteTx(tx *gorm.DB, id uuid.UUID) error {
	res := tx.Delete(&models.Task{}, "id = ?", id)
	if res.Error != nil {
		return mapError(res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
