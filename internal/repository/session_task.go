package repository

import (
	"time"

	"diplom/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionTaskRepository struct {
	db *gorm.DB
}

func NewSessionTaskRepository(db *gorm.DB) *SessionTaskRepository {
	return &SessionTaskRepository{db: db}
}

func (r *SessionTaskRepository) Create(st *models.SessionTask) error {
	return mapError(r.db.Create(st).Error)
}

func (r *SessionTaskRepository) CreateBatch(sessionID uuid.UUID, taskIDs []uuid.UUID) error {
	if len(taskIDs) == 0 {
		return nil
	}
	now := time.Now()
	sessionTasks := make([]models.SessionTask, 0, len(taskIDs))
	for _, taskID := range taskIDs {
		sessionTasks = append(sessionTasks, models.SessionTask{
			ID:        uuid.New(),
			SessionID: sessionID,
			TaskID:    taskID,
			CreatedAt: now,
		})
	}
	return mapError(r.db.CreateInBatches(sessionTasks, 100).Error)
}

func (r *SessionTaskRepository) ListBySessionID(sessionID uuid.UUID) ([]models.SessionTask, error) {
	var items []models.SessionTask
	err := r.db.Where("session_id = ?", sessionID).Find(&items).Error
	return items, mapError(err)
}

// ListTaskIDsBySessionID returns just the task IDs for a session
func (r *SessionTaskRepository) ListTaskIDsBySessionID(sessionID uuid.UUID) ([]uuid.UUID, error) {
	var taskIDs []uuid.UUID
	err := r.db.Model(&models.SessionTask{}).
		Where("session_id = ?", sessionID).
		Pluck("task_id", &taskIDs).Error
	return taskIDs, mapError(err)
}

func (r *SessionTaskRepository) Delete(sessionID, taskID uuid.UUID) error {
	res := r.db.Where("session_id = ? AND task_id = ?", sessionID, taskID).Delete(&models.SessionTask{})
	return mapError(res.Error)
}
