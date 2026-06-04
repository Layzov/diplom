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
	db          *gorm.DB
	users       *repository.UserRepository
	subjects    *repository.SubjectRepository
	tasks       *repository.TaskRepository
	sessions    *repository.SessionRepository
	attempts    *repository.TaskAttemptRepository
	stats       *repository.AttemptStatsRepository
	reps        *repository.RepetitionRepository
	sessionTasks *repository.SessionTaskRepository
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
	sessionTasks *repository.SessionTaskRepository,
) *SessionService {
	return &SessionService{
		db: db, users: users, subjects: subjects, tasks: tasks,
		sessions: sessions, attempts: attempts, stats: stats, reps: reps,
		sessionTasks: sessionTasks,
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

	if req.TopicID != nil {
		topic, err := s.tasks.GetByID(*req.TopicID) // TODO: should be TopicRepository
		if err != nil {
			return nil, err
		}
		// Verify topic belongs to subject if provided
		if req.SubjectID != nil && topic.TopicID != *req.TopicID {
			return nil, apperror.Validation("topic not in subject")
		}
	}

	// Default values
	taskCount := req.TaskCount
	if taskCount == 0 {
		taskCount = 15
	}
	mode := models.SessionModeLearning
	if req.Mode == "review" {
		mode = models.SessionModeReview
	}

	now := time.Now().UTC()
	session := &models.Session{
		UserID:    userID,
		SubjectID: req.SubjectID,
		TopicID:   req.TopicID,
		TaskCount: taskCount,
		Mode:      mode,
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
	var repetition *models.Repetition

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

		// Only create repetition for failed/partial/skipped, not for correct
		if req.Result != models.AnswerResultCorrect {
			quality := DetermineQuality(req.Result)
			repeatAt := ScheduleRepeatAt(req.Result, now)
			rep := models.Repetition{
				SessionID:     sessionID,
				TaskAttemptID: attempt.ID,
				TaskID:        task.ID,
				UserID:        session.UserID,
				TopicID:       task.TopicID,
				RepeatAt:      repeatAt,
				Quality:       quality,
				Status:        models.RepetitionStatusPlanned,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			if err := s.reps.CreateTx(tx, &rep); err != nil {
				return err
			}
			repetition = &rep
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var repResp *dto.RepetitionResponse
	if repetition != nil {
		r := dto.RepetitionFromModel(repetition)
		repResp = &r
	}

	return &dto.SubmitAttemptResponse{
		Attempt:    dto.TaskAttemptFromModel(&attempt),
		Repetition: repResp,
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

// GetTasksForSession returns tasks for a session, applying smart selection based on mode.
// For "review" mode: prioritizes failed/struggling tasks.
// For "learning" mode: returns regular tasks from the topic.
func (s *SessionService) GetTasksForSession(sessionID, authUserID uuid.UUID) (dto.SessionTasksResponse, error) {
	session, err := s.sessionOwned(sessionID, authUserID)
	if err != nil {
		return dto.SessionTasksResponse{}, err
	}

	if session.TopicID == nil {
		return dto.SessionTasksResponse{}, apperror.Validation("session must have a topic")
	}

	var tasks []models.Task

	if session.Mode == models.SessionModeReview {
		// Smart selection: prioritize failed tasks
		tasks, err = s.tasks.ListForReview(session.UserID, *session.TopicID, session.TaskCount)
	} else {
		// Learning mode: regular tasks from topic
		allTasks, err := s.tasks.ListByTopicID(*session.TopicID)
		if err != nil {
			return dto.SessionTasksResponse{}, err
		}
		if len(allTasks) > session.TaskCount {
			tasks = allTasks[:session.TaskCount]
		} else {
			tasks = allTasks
		}
	}

	if err != nil {
		return dto.SessionTasksResponse{}, err
	}

	// Bind tasks to session
	taskIDs := make([]uuid.UUID, 0, len(tasks))
	for _, t := range tasks {
		taskIDs = append(taskIDs, t.ID)
	}

	sessionTasks := make([]models.SessionTask, 0, len(taskIDs))
	for _, tid := range taskIDs {
		sessionTasks = append(sessionTasks, models.SessionTask{
			ID:        uuid.New(),
			SessionID: sessionID,
			TaskID:    tid,
			CreatedAt: time.Now().UTC(),
		})
	}

	if err := s.db.CreateInBatches(sessionTasks, 100).Error; err != nil {
		return dto.SessionTasksResponse{}, err
	}

	taskResponses := make([]dto.TaskResponse, 0, len(tasks))
	for i := range tasks {
		taskResponses = append(taskResponses, dto.TaskFromModel(&tasks[i]))
	}

	return dto.SessionTasksResponse{
		SessionID: sessionID,
		Tasks:     taskResponses,
	}, nil
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
