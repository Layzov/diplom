package service

import (
	"strconv"

	"diplom/internal/apperror"
	"diplom/internal/dto"
	"diplom/internal/repository"

	"github.com/google/uuid"
)

const (
	defaultUpcomingLimit = 10
	maxUpcomingLimit     = 100
)

type StatisticsService struct {
	stats *repository.StatsRepository
}

func NewStatisticsService(stats *repository.StatsRepository) *StatisticsService {
	return &StatisticsService{stats: stats}
}

func (s *StatisticsService) UserOverview(userID uuid.UUID) (*dto.UserStatsResponse, error) {
	if err := s.stats.EnsureUserExists(userID); err != nil {
		return nil, err
	}
	return s.stats.GetUserStats(userID)
}

func (s *StatisticsService) TopicsByUser(userID uuid.UUID) (dto.TopicStatsListResponse, error) {
	if err := s.stats.EnsureUserExists(userID); err != nil {
		return dto.TopicStatsListResponse{}, err
	}
	items, err := s.stats.ListTopicStats(userID)
	if err != nil {
		return dto.TopicStatsListResponse{}, err
	}
	if items == nil {
		items = []dto.TopicStatsResponse{}
	}
	return dto.TopicStatsListResponse{Items: items}, nil
}

func (s *StatisticsService) SessionsByUser(userID uuid.UUID) (dto.SessionStatsListResponse, error) {
	if err := s.stats.EnsureUserExists(userID); err != nil {
		return dto.SessionStatsListResponse{}, err
	}
	items, err := s.stats.ListSessionStats(userID)
	if err != nil {
		return dto.SessionStatsListResponse{}, err
	}
	if items == nil {
		items = []dto.SessionStatsResponse{}
	}
	return dto.SessionStatsListResponse{Items: items}, nil
}

func (s *StatisticsService) UpcomingRepetitions(userID uuid.UUID, limit int, includeOverdue bool) (dto.UpcomingListResponse, error) {
	if err := s.stats.EnsureUserExists(userID); err != nil {
		return dto.UpcomingListResponse{}, err
	}
	if limit <= 0 {
		limit = defaultUpcomingLimit
	}
	if limit > maxUpcomingLimit {
		limit = maxUpcomingLimit
	}
	items, err := s.stats.ListUpcoming(userID, limit, includeOverdue)
	if err != nil {
		return dto.UpcomingListResponse{}, err
	}
	if items == nil {
		items = []dto.UpcomingRepetitionResponse{}
	}
	return dto.UpcomingListResponse{Items: items}, nil
}

func ParseUpcomingLimit(raw string) (int, error) {
	if raw == "" {
		return defaultUpcomingLimit, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 0, apperror.Validation("limit must be a positive integer")
	}
	return n, nil
}
