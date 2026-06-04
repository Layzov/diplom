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

// ListForReview returns tasks from a topic prioritizing failed ones.
// Priority: STRUGGLING (failed) > LEARNING (partial) > others.
func (r *TaskRepository) ListForReview(userID, topicID uuid.UUID, limit int) ([]models.Task, error) {
	var tasks []models.Task

	// Get task IDs sorted by struggle level (failed/wrong first)
	type TaskScore struct {
		TaskID    uuid.UUID
		Priority  int // 1=struggling, 2=learning, 3=other
	}

	var scores []TaskScore
	err := r.db.
		Model(&models.TaskAttempt{}).
		Distinct("task_id").
		Select(`task_id, 
		CASE 
			WHEN r.quality = 'struggling' THEN 1
			WHEN r.quality = 'learning' THEN 2
			ELSE 3
		END as priority`).
		Joins("JOIN repetitions r ON task_attempts.id = r.task_attempt_id").
		Where("task_attempts.user_id = ? AND task_attempts.result IN ('wrong', 'skipped')", userID).
		Order("priority ASC").
		Scan(&scores).Error

	if err != nil {
		return nil, mapError(err)
	}

	if len(scores) > 0 {
		taskIDs := make([]uuid.UUID, 0, len(scores))
		for _, s := range scores {
			taskIDs = append(taskIDs, s.TaskID)
		}

		// Get full task objects, limited
		err = r.db.
			Where("topic_id = ? AND id IN ?", topicID, taskIDs).
			Limit(limit).
			Find(&tasks).Error
		return tasks, mapError(err)
	}

	// If no failed tasks, return regular tasks from topic
	err = r.db.
		Where("topic_id = ?", topicID).
		Order("created_at DESC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, mapError(err)
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
