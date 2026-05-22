package service

import (
	"strings"

	"diplom/internal/apperror"
	"diplom/internal/dto"
	"diplom/internal/models"
	"diplom/internal/repository"

	"github.com/google/uuid"
)

type SubjectService struct {
	users    *repository.UserRepository
	subjects *repository.SubjectRepository
}

func NewSubjectService(users *repository.UserRepository, subjects *repository.SubjectRepository) *SubjectService {
	return &SubjectService{users: users, subjects: subjects}
}

func (s *SubjectService) Create(userID uuid.UUID, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
	if err := validateNonEmpty(req.Title, "title"); err != nil {
		return nil, err
	}
	ok, err := s.users.Exists(userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrNotFound
	}

	subject := &models.Subject{
		UserID:      userID,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
	}
	if err := s.subjects.Create(subject); err != nil {
		return nil, err
	}
	resp := dto.SubjectFromModel(subject)
	return &resp, nil
}

func (s *SubjectService) ListByUser(userID uuid.UUID) (dto.SubjectListResponse, error) {
	ok, err := s.users.Exists(userID)
	if err != nil {
		return dto.SubjectListResponse{}, err
	}
	if !ok {
		return dto.SubjectListResponse{}, apperror.ErrNotFound
	}

	items, err := s.subjects.ListByUserID(userID)
	if err != nil {
		return dto.SubjectListResponse{}, err
	}
	out := make([]dto.SubjectResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.SubjectFromModel(&items[i]))
	}
	return dto.SubjectListResponse{Items: out}, nil
}

func (s *SubjectService) Get(id uuid.UUID) (*dto.SubjectResponse, error) {
	subject, err := s.subjects.GetByID(id)
	if err != nil {
		return nil, err
	}
	resp := dto.SubjectFromModel(subject)
	return &resp, nil
}

func (s *SubjectService) Update(id uuid.UUID, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error) {
	subject, err := s.subjects.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		if err := validateNonEmpty(*req.Title, "title"); err != nil {
			return nil, err
		}
		subject.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		subject.Description = strings.TrimSpace(*req.Description)
	}
	if err := s.subjects.Update(subject); err != nil {
		return nil, err
	}
	resp := dto.SubjectFromModel(subject)
	return &resp, nil
}

func (s *SubjectService) Delete(id uuid.UUID) error {
	return s.subjects.Delete(id)
}
