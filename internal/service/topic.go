package service

import (
	"strings"

	"diplom/internal/dto"
	"diplom/internal/models"
	"diplom/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopicService struct {
	db       *gorm.DB
	subjects *repository.SubjectRepository
	topics   *repository.TopicRepository
}

func NewTopicService(db *gorm.DB, subjects *repository.SubjectRepository, topics *repository.TopicRepository) *TopicService {
	return &TopicService{db: db, subjects: subjects, topics: topics}
}

func (s *TopicService) Create(subjectID, authUserID uuid.UUID, req dto.CreateTopicRequest) (*dto.TopicResponse, error) {
	if err := validateNonEmpty(req.Title, "title"); err != nil {
		return nil, err
	}
	if err := s.subjectOwned(subjectID, authUserID); err != nil {
		return nil, err
	}

	var topic models.Topic
	err := s.db.Transaction(func(tx *gorm.DB) error {
		topic = models.Topic{
			SubjectID: subjectID,
			Title:     strings.TrimSpace(req.Title),
		}
		if err := s.topics.CreateTx(tx, &topic); err != nil {
			return err
		}
		return s.subjects.AdjustTopicCount(tx, subjectID, 1)
	})
	if err != nil {
		return nil, err
	}
	resp := dto.TopicFromModel(&topic)
	return &resp, nil
}

func (s *TopicService) ListBySubject(subjectID, authUserID uuid.UUID) (dto.TopicListResponse, error) {
	if err := s.subjectOwned(subjectID, authUserID); err != nil {
		return dto.TopicListResponse{}, err
	}
	items, err := s.topics.ListBySubjectID(subjectID)
	if err != nil {
		return dto.TopicListResponse{}, err
	}
	out := make([]dto.TopicResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.TopicFromModel(&items[i]))
	}
	return dto.TopicListResponse{Items: out}, nil
}

func (s *TopicService) Get(id, authUserID uuid.UUID) (*dto.TopicResponse, error) {
	topic, err := s.topicOwned(id, authUserID)
	if err != nil {
		return nil, err
	}
	resp := dto.TopicFromModel(topic)
	return &resp, nil
}

func (s *TopicService) Update(id, authUserID uuid.UUID, req dto.UpdateTopicRequest) (*dto.TopicResponse, error) {
	topic, err := s.topicOwned(id, authUserID)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		if err := validateNonEmpty(*req.Title, "title"); err != nil {
			return nil, err
		}
		topic.Title = strings.TrimSpace(*req.Title)
	}
	if err := s.topics.Update(topic); err != nil {
		return nil, err
	}
	resp := dto.TopicFromModel(topic)
	return &resp, nil
}

func (s *TopicService) Delete(id, authUserID uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		topic, err := s.topicOwned(id, authUserID)
		if err != nil {
			return err
		}
		if err := s.topics.DeleteTx(tx, id); err != nil {
			return err
		}
		return s.subjects.AdjustTopicCount(tx, topic.SubjectID, -1)
	})
}

// GetTaskStatsByTopic returns detailed learning stats for each task in a topic.
func (s *TopicService) GetTaskStatsByTopic(topicID, userID uuid.UUID) (interface{}, error) {
	if _, err := s.topicOwned(topicID, userID); err != nil {
		return nil, err
	}

	var stats []map[string]interface{}
	err := s.db.Raw(`
		SELECT
			task_id,
			task_title,
			topic_id,
			attempt_count,
			unique_users,
			correct_count,
			wrong_count,
			partial_count,
			skipped_count,
			success_rate,
			avg_response_time_ms,
			planned_repetitions,
			struggling_count,
			learning_count,
			mastered_count,
			last_attempted_at
		FROM task_learning_stats
		WHERE topic_id = ? AND user_id = ?
		ORDER BY struggle_count DESC, attempt_count DESC
	`, topicID, userID).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"topic_id": topicID,
		"tasks":    stats,
	}, nil
}

func (s *TopicService) subjectOwned(subjectID, authUserID uuid.UUID) error {
	subject, err := s.subjects.GetByID(subjectID)
	if err != nil {
		return err
	}
	return ensureOwner(subject.UserID, authUserID)
}

func (s *TopicService) topicOwned(topicID, authUserID uuid.UUID) (*models.Topic, error) {
	topic, err := s.topics.GetByID(topicID)
	if err != nil {
		return nil, err
	}
	if err := s.subjectOwned(topic.SubjectID, authUserID); err != nil {
		return nil, err
	}
	return topic, nil
}
