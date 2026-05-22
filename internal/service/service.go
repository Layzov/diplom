package service

import (
	"diplom/internal/repository"

	"gorm.io/gorm"
)

// Services groups domain services for HTTP handlers.
type Services struct {
	Users    *UserService
	Subjects *SubjectService
	Topics   *TopicService
	Tasks    *TaskService
}

func New(db *gorm.DB) *Services {
	userRepo := repository.NewUserRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	return &Services{
		Users:    NewUserService(userRepo),
		Subjects: NewSubjectService(userRepo, subjectRepo),
		Topics:   NewTopicService(db, subjectRepo, topicRepo),
		Tasks:    NewTaskService(db, topicRepo, taskRepo),
	}
}
