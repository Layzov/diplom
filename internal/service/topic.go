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

func (s *TopicService) Create(subjectID uuid.UUID, req dto.CreateTopicRequest) (*dto.TopicResponse, error) {
	if err := validateNonEmpty(req.Title, "title"); err != nil {
		return nil, err
	}

	var topic models.Topic
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if _, err := s.subjects.GetByID(subjectID); err != nil {
			return err
		}
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

func (s *TopicService) ListBySubject(subjectID uuid.UUID) (dto.TopicListResponse, error) {
	if _, err := s.subjects.GetByID(subjectID); err != nil {
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

func (s *TopicService) Get(id uuid.UUID) (*dto.TopicResponse, error) {
	topic, err := s.topics.GetByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.TopicFromModel(topic)
	return &resp, nil
}

func (s *TopicService) Update(id uuid.UUID, req dto.UpdateTopicRequest) (*dto.TopicResponse, error) {
	topic, err := s.topics.GetByID(id)
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

func (s *TopicService) Delete(id uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		topic, err := s.topics.GetByID(id)
		if err != nil {
			return err
		}
		if err := s.topics.DeleteTx(tx, id); err != nil {
			return err
		}
		return s.subjects.AdjustTopicCount(tx, topic.SubjectID, -1)
	})
}
