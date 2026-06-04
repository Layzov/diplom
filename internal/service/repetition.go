package service

import (
	"time"

	"diplom/internal/apperror"
	"diplom/internal/dto"
	"diplom/internal/models"
	"diplom/internal/repository"

	"github.com/google/uuid"
)

type RepetitionService struct {
	users *repository.UserRepository
	reps  *repository.RepetitionRepository
	cal   *repository.CalendarRepository
}

func NewRepetitionService(
	users *repository.UserRepository,
	reps *repository.RepetitionRepository,
	cal *repository.CalendarRepository,
) *RepetitionService {
	return &RepetitionService{users: users, reps: reps, cal: cal}
}

func (s *RepetitionService) ListByUser(userID uuid.UUID, status *models.RepetitionStatus) (dto.RepetitionListResponse, error) {
	ok, err := s.users.Exists(userID)
	if err != nil {
		return dto.RepetitionListResponse{}, err
	}
	if !ok {
		return dto.RepetitionListResponse{}, apperror.ErrNotFound
	}
	items, err := s.reps.ListByUserID(userID, status)
	if err != nil {
		return dto.RepetitionListResponse{}, err
	}
	out := make([]dto.RepetitionResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.RepetitionFromModel(&items[i]))
	}
	return dto.RepetitionListResponse{Items: out}, nil
}

func (s *RepetitionService) UpdateStatus(id, authUserID uuid.UUID, req dto.UpdateRepetitionRequest) (*dto.RepetitionResponse, error) {
	switch req.Status {
	case models.RepetitionStatusPlanned,
		models.RepetitionStatusDone,
		models.RepetitionStatusSkipped:
	default:
		return nil, apperror.Validation("invalid repetition status")
	}

	rep, err := s.reps.GetByID(id)
	if err != nil {
		return nil, err
	}
	if err := ensureOwner(rep.UserID, authUserID); err != nil {
		return nil, err
	}
	rep.Status = req.Status
	rep.UpdatedAt = time.Now().UTC()
	if err := s.reps.Update(rep); err != nil {
		return nil, err
	}
	resp := dto.RepetitionFromModel(rep)
	return &resp, nil
}

func (s *RepetitionService) Calendar(userID uuid.UUID, from, to *time.Time) (dto.CalendarListResponse, error) {
	ok, err := s.users.Exists(userID)
	if err != nil {
		return dto.CalendarListResponse{}, err
	}
	if !ok {
		return dto.CalendarListResponse{}, apperror.ErrNotFound
	}
	items, err := s.cal.List(userID, from, to)
	if err != nil {
		return dto.CalendarListResponse{}, err
	}
	return dto.CalendarListResponse{Items: items}, nil
}
