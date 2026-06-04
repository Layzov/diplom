package service

import (
	"diplom/internal/auth"
	"diplom/internal/repository"

	"gorm.io/gorm"
)

// Services groups domain services for HTTP handlers.
type Services struct {
	Users       *UserService
	Subjects    *SubjectService
	Topics      *TopicService
	Tasks       *TaskService
	Sessions    *SessionService
	Repetitions *RepetitionService
	Statistics  *StatisticsService
	Auth        *AuthService
}

func New(db *gorm.DB, jwt *auth.Manager) *Services {
	userRepo := repository.NewUserRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	sessionTaskRepo := repository.NewSessionTaskRepository(db)
	attemptRepo := repository.NewTaskAttemptRepository(db)
	attemptStatsRepo := repository.NewAttemptStatsRepository(db)
	repRepo := repository.NewRepetitionRepository(db)
	calRepo := repository.NewCalendarRepository(db)
	analyticsRepo := repository.NewStatsRepository(db)

	svc := &Services{
		Users:    NewUserService(userRepo),
		Subjects: NewSubjectService(userRepo, subjectRepo),
		Topics:   NewTopicService(db, subjectRepo, topicRepo),
		Tasks:    NewTaskService(db, subjectRepo, topicRepo, taskRepo),
		Sessions: NewSessionService(
			db, userRepo, subjectRepo, taskRepo,
			sessionRepo, attemptRepo, attemptStatsRepo, repRepo, sessionTaskRepo,
		),
		Repetitions: NewRepetitionService(userRepo, repRepo, calRepo),
		Statistics:  NewStatisticsService(analyticsRepo),
	}
	svc.Auth = NewAuthService(svc.Users, jwt)
	return svc
}
