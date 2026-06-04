package service

import (
	"strings"

	"diplom/internal/dto"
	"diplom/internal/models"
	"diplom/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskService struct {
	db       *gorm.DB
	subjects *repository.SubjectRepository
	topics   *repository.TopicRepository
	tasks    *repository.TaskRepository
}

func NewTaskService(
	db *gorm.DB,
	subjects *repository.SubjectRepository,
	topics *repository.TopicRepository,
	tasks *repository.TaskRepository,
) *TaskService {
	return &TaskService{db: db, subjects: subjects, topics: topics, tasks: tasks}
}

func (s *TaskService) Create(topicID, authUserID uuid.UUID, req dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	if err := validateNonEmpty(req.Title, "title"); err != nil {
		return nil, err
	}
	if err := validateNonEmpty(req.Content, "content"); err != nil {
		return nil, err
	}
	if err := s.topicOwned(topicID, authUserID); err != nil {
		return nil, err
	}

	taskType := req.Type
	if taskType == "" {
		taskType = models.TaskTypeFlashcard
	}
	if err := validateTaskType(taskType); err != nil {
		return nil, err
	}

	var task models.Task
	err := s.db.Transaction(func(tx *gorm.DB) error {
		task = models.Task{
			TopicID: topicID,
			Title:   strings.TrimSpace(req.Title),
			Content: strings.TrimSpace(req.Content),
			Type:    taskType,
		}
		if err := s.tasks.CreateTx(tx, &task); err != nil {
			return err
		}
		return s.topics.AdjustTaskCount(tx, topicID, 1)
	})
	if err != nil {
		return nil, err
	}
	resp := dto.TaskFromModel(&task)
	return &resp, nil
}

func (s *TaskService) ListByTopic(topicID, authUserID uuid.UUID) (dto.TaskListResponse, error) {
	if err := s.topicOwned(topicID, authUserID); err != nil {
		return dto.TaskListResponse{}, err
	}
	items, err := s.tasks.ListByTopicID(topicID)
	if err != nil {
		return dto.TaskListResponse{}, err
	}
	out := make([]dto.TaskResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.TaskFromModel(&items[i]))
	}
	return dto.TaskListResponse{Items: out}, nil
}

func (s *TaskService) Get(id, authUserID uuid.UUID) (*dto.TaskResponse, error) {
	task, err := s.tasks.GetByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.topicOwned(task.TopicID, authUserID); err != nil {
		return nil, err
	}
	resp := dto.TaskFromModel(task)
	return &resp, nil
}

func (s *TaskService) Update(id, authUserID uuid.UUID, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.tasks.GetByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.topicOwned(task.TopicID, authUserID); err != nil {
		return nil, err
	}
	if req.Title != nil {
		if err := validateNonEmpty(*req.Title, "title"); err != nil {
			return nil, err
		}
		task.Title = strings.TrimSpace(*req.Title)
	}
	if req.Content != nil {
		if err := validateNonEmpty(*req.Content, "content"); err != nil {
			return nil, err
		}
		task.Content = strings.TrimSpace(*req.Content)
	}
	if req.Type != nil {
		if err := validateTaskType(*req.Type); err != nil {
			return nil, err
		}
		task.Type = *req.Type
	}
	if err := s.tasks.Update(task); err != nil {
		return nil, err
	}
	resp := dto.TaskFromModel(task)
	return &resp, nil
}

func (s *TaskService) Delete(id, authUserID uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		task, err := s.tasks.GetByID(id)
		if err != nil {
			return err
		}
		if err := s.topicOwned(task.TopicID, authUserID); err != nil {
			return err
		}
		if err := s.tasks.DeleteTx(tx, id); err != nil {
			return err
		}
		return s.topics.AdjustTaskCount(tx, task.TopicID, -1)
	})
}

func (s *TaskService) topicOwned(topicID, authUserID uuid.UUID) error {
	topic, err := s.topics.GetByID(topicID)
	if err != nil {
		return err
	}
	subject, err := s.subjects.GetByID(topic.SubjectID)
	if err != nil {
		return err
	}
	return ensureOwner(subject.UserID, authUserID)
}
