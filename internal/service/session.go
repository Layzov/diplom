package service

import (
	"time"

	"diplom/internal/apperror"
	"diplom/internal/dto"
	"diplom/internal/models"
	"diplom/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionService struct {
	db        *gorm.DB
	users     *repository.UserRepository
	subjects  *repository.SubjectRepository
	tasks     *repository.TaskRepository
	sessions  *repository.SessionRepository
	attempts  *repository.TaskAttemptRepository
	stats     *repository.AttemptStatsRepository
	reps      *repository.RepetitionRepository
}

func NewSessionService(
	db *gorm.DB,
	users *repository.UserRepository,
	subjects *repository.SubjectRepository,
	tasks *repository.TaskRepository,
	sessions *repository.SessionRepository,
	attempts *repository.TaskAttemptRepository,
	stats *repository.AttemptStatsRepository,
	reps *repository.RepetitionRepository,
) *SessionService {
	return &SessionService{
		db: db, users: users, subjects: subjects, tasks: tasks,
		sessions: sessions, attempts: attempts, stats: stats, reps: reps,
	}
}

func (s *SessionService) Create(userID uuid.UUID, req dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	if req.SubjectID != nil {
		subject, err := s.subjects.GetByID(*req.SubjectID)
		if err != nil {
			return nil, err
		}
		if err := ensureOwner(subject.UserID, userID); err != nil {
			return nil, err
		}
	}

	now := time.Now().UTC()
	session := &models.Session{
		UserID:    userID,
		SubjectID: req.SubjectID,
		StartedAt: now,
		Notes:     req.Notes,
	}
	if err := s.sessions.Create(session); err != nil {
		return nil, err
	}
	resp := dto.SessionFromModel(session)
	return &resp, nil
}

func (s *SessionService) Get(id, authUserID uuid.UUID) (*dto.SessionResponse, error) {
	session, err := s.sessionOwned(id, authUserID)
	if err != nil {
		return nil, err
	}
	resp := dto.SessionFromModel(session)
	return &resp, nil
}

func (s *SessionService) ListByUser(userID uuid.UUID) (dto.SessionListResponse, error) {
	ok, err := s.users.Exists(userID)
	if err != nil {
		return dto.SessionListResponse{}, err
	}
	if !ok {
		return dto.SessionListResponse{}, apperror.ErrNotFound
	}
	items, err := s.sessions.ListByUserID(userID)
	if err != nil {
		return dto.SessionListResponse{}, err
	}
	out := make([]dto.SessionResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.SessionFromModel(&items[i]))
	}
	return dto.SessionListResponse{Items: out}, nil
}

func (s *SessionService) Finish(id, authUserID uuid.UUID, req dto.FinishSessionRequest) (*dto.AttemptStatsResponse, error) {
	if _, err := s.sessionOwned(id, authUserID); err != nil {
		return nil, err
	}
	var statsOut models.AttemptStats

	err := s.db.Transaction(func(tx *gorm.DB) error {
		session, err := s.sessions.GetByIDTx(tx, id)
		if err != nil {
			return err
		}
		if session.EndedAt != nil {
			return apperror.Validation("session already finished")
		}

		now := time.Now().UTC()
		session.EndedAt = &now
		if req.Notes != nil {
			session.Notes = req.Notes
		}
		if err := s.sessions.SaveTx(tx, session); err != nil {
			return err
		}

		attempts, err := s.attempts.ListBySessionIDTx(tx, id)
		if err != nil {
			return err
		}

		nextRepeat, err := s.reps.EarliestPlannedForSessionTx(tx, id)
		if err != nil {
			return err
		}

		stats := buildAttemptStats(id, attempts, nextRepeat)
		if err := s.stats.UpsertTx(tx, &stats); err != nil {
			return err
		}
		statsOut = stats
		return nil
	})
	if err != nil {
		return nil, err
	}
	resp := dto.AttemptStatsFromModel(&statsOut)
	return &resp, nil
}

func (s *SessionService) SubmitAttempt(sessionID, authUserID uuid.UUID, req dto.CreateTaskAttemptRequest) (*dto.SubmitAttemptResponse, error) {
	if err := validateAnswerResult(req.Result); err != nil {
		return nil, err
	}

	var attempt models.TaskAttempt
	var repetition models.Repetition

	err := s.db.Transaction(func(tx *gorm.DB) error {
		session, err := s.sessions.GetByIDTx(tx, sessionID)
		if err != nil {
			return err
		}
		if err := ensureOwner(session.UserID, authUserID); err != nil {
			return err
		}
		if session.EndedAt != nil {
			return apperror.Validation("session already finished")
		}

		task, err := s.tasks.GetByID(req.TaskID)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		attempt = models.TaskAttempt{
			SessionID:      sessionID,
			TaskID:         req.TaskID,
			UserID:         session.UserID,
			AnswerText:     req.AnswerText,
			SelectedOption: req.SelectedOption,
			IsCorrect:      req.IsCorrect,
			Result:         req.Result,
			ResponseTimeMs: req.ResponseTimeMs,
			ReviewedAt:     &now,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := s.attempts.CreateTx(tx, &attempt); err != nil {
			return err
		}

		repeatAt := ScheduleRepeatAt(req.Result, now)
		repetition = models.Repetition{
			SessionID:     sessionID,
			TaskAttemptID: attempt.ID,
			TaskID:        task.ID,
			UserID:        session.UserID,
			TopicID:       task.TopicID,
			RepeatAt:      repeatAt,
			Status:        models.RepetitionStatusPlanned,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		return s.reps.CreateTx(tx, &repetition)
	})
	if err != nil {
		return nil, err
	}

	return &dto.SubmitAttemptResponse{
		Attempt:    dto.TaskAttemptFromModel(&attempt),
		Repetition: dto.RepetitionFromModel(&repetition),
	}, nil
}

func (s *SessionService) ListAttempts(sessionID, authUserID uuid.UUID) (dto.TaskAttemptListResponse, error) {
	if _, err := s.sessionOwned(sessionID, authUserID); err != nil {
		return dto.TaskAttemptListResponse{}, err
	}
	items, err := s.attempts.ListBySessionID(sessionID)
	if err != nil {
		return dto.TaskAttemptListResponse{}, err
	}
	out := make([]dto.TaskAttemptResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.TaskAttemptFromModel(&items[i]))
	}
	return dto.TaskAttemptListResponse{Items: out}, nil
}

func (s *SessionService) GetStats(sessionID, authUserID uuid.UUID) (*dto.AttemptStatsResponse, error) {
	if _, err := s.sessionOwned(sessionID, authUserID); err != nil {
		return nil, err
	}
	stats, err := s.stats.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}
	resp := dto.AttemptStatsFromModel(stats)
	return &resp, nil
}

func validateAnswerResult(r models.AnswerResult) error {
	switch r {
	case models.AnswerResultCorrect,
		models.AnswerResultPartial,
		models.AnswerResultWrong,
		models.AnswerResultSkipped:
		return nil
	default:
		return apperror.Validation("invalid answer result")
	}
}

func (s *SessionService) sessionOwned(sessionID, authUserID uuid.UUID) (*models.Session, error) {
	session, err := s.sessions.GetByID(sessionID)
	if err != nil {
		return nil, err
	}
	if err := ensureOwner(session.UserID, authUserID); err != nil {
		return nil, err
	}
	return session, nil
}
