package seeder

import (
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"diplom/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Config controls how much test data to generate.
type Config struct {
	Users            int // number of users
	SubjectsPerUser  int // subjects per user
	TopicsPerSubject int // topics per subject
	TasksPerTopic    int // tasks per topic
	SessionsPerUser  int // study sessions per user
}

// DefaultConfig provides reasonable defaults.
func DefaultConfig() Config {
	return Config{
		Users:            5,
		SubjectsPerUser:  2,
		TopicsPerSubject: 3,
		TasksPerTopic:    5,
		SessionsPerUser:  5,
	}
}

// Seed generates test data in the database.
func Seed(db *gorm.DB, cfg Config, log *slog.Logger) error {
	log.Info("seeding database",
		slog.Int("users", cfg.Users),
		slog.Int("subjects_per_user", cfg.SubjectsPerUser),
		slog.Int("topics_per_subject", cfg.TopicsPerSubject),
		slog.Int("tasks_per_topic", cfg.TasksPerTopic),
		slog.Int("sessions_per_user", cfg.SessionsPerUser),
	)

	users := make([]*models.User, 0, cfg.Users)
	for i := 0; i < cfg.Users; i++ {
		u := generateUser(i)
		users = append(users, u)
		if err := db.Create(u).Error; err != nil {
			return fmt.Errorf("create user %d: %w", i, err)
		}
		log.Debug("created user", slog.String("email", u.Email))
	}

	subjects := make([]*models.Subject, 0, cfg.Users*cfg.SubjectsPerUser)
	for _, user := range users {
		for j := 0; j < cfg.SubjectsPerUser; j++ {
			s := generateSubject(user.ID, j)
			subjects = append(subjects, s)
			if err := db.Create(s).Error; err != nil {
				return fmt.Errorf("create subject: %w", err)
			}
		}
		log.Debug("created subjects", slog.Int("count", cfg.SubjectsPerUser), slog.String("user", user.Email))
	}

	topics := make([]*models.Topic, 0, len(subjects)*cfg.TopicsPerSubject)
	for _, subject := range subjects {
		for j := 0; j < cfg.TopicsPerSubject; j++ {
			t := generateTopic(subject.ID, j)
			topics = append(topics, t)
			if err := db.Create(t).Error; err != nil {
				return fmt.Errorf("create topic: %w", err)
			}
		}
	}
	log.Debug("created topics", slog.Int("count", len(topics)))

	tasks := make([]*models.Task, 0, len(topics)*cfg.TasksPerTopic)
	for _, topic := range topics {
		for j := 0; j < cfg.TasksPerTopic; j++ {
			t := generateTask(topic.ID, j)
			tasks = append(tasks, t)
			if err := db.Create(t).Error; err != nil {
				return fmt.Errorf("create task: %w", err)
			}
		}
	}
	log.Debug("created tasks", slog.Int("count", len(tasks)))

	// Sessions: create attempts and repetitions
	for _, user := range users {
		for s := 0; s < cfg.SessionsPerUser; s++ {
			session := generateSession(user.ID, s)
			if err := db.Create(session).Error; err != nil {
				return fmt.Errorf("create session: %w", err)
			}

			// Generate task attempts for this session
			userTasks := make([]*models.Task, 0)
			if err := db.Where("id IN (?)", db.Table("topics").
				Where("subject_id IN (?)", db.Table("subjects").Where("user_id = ?", user.ID).Select("id")).
				Select("id")).
				Find(&userTasks).Error; err != nil {
				return fmt.Errorf("fetch user tasks: %w", err)
			}

			// 50% of user tasks get attempted
			var attemptedTasks []*models.Task
			for _, t := range userTasks {
				if rand.Float64() < 0.5 {
					attemptedTasks = append(attemptedTasks, t)
				}
			}
			if len(attemptedTasks) == 0 && len(userTasks) > 0 {
				attemptedTasks = append(attemptedTasks, userTasks[0])
			}

			stats := &models.AttemptStats{
				ID:        uuid.New(),
				SessionID: session.ID,
			}
			correctCount := 0
			wrongCount := 0
			partialCount := 0
			skippedCount := 0
			totalMS := 0

			for _, task := range attemptedTasks {
				attempt := generateTaskAttempt(session.ID, user.ID, task.ID)
				if err := db.Create(attempt).Error; err != nil {
					return fmt.Errorf("create task attempt: %w", err)
				}

				// Update stats
				switch attempt.Result {
				case models.AnswerResultCorrect:
					correctCount++
				case models.AnswerResultWrong:
					wrongCount++
				case models.AnswerResultPartial:
					partialCount++
				case models.AnswerResultSkipped:
					skippedCount++
				}
				totalMS += attempt.ResponseTimeMs
			}

			stats.TotalTasks = len(attemptedTasks)
			stats.CorrectCount = correctCount
			stats.WrongCount = wrongCount
			stats.PartialCount = partialCount
			stats.SkippedCount = skippedCount
			stats.TotalTimeMs = totalMS
			if len(attemptedTasks) > 0 {
				stats.AverageTimeMs = totalMS / len(attemptedTasks)
			}
			if len(attemptedTasks) > 0 {
				stats.SuccessRate = (float64(correctCount) / float64(len(attemptedTasks))) * 100
			}
			nextRepeat := time.Now().AddDate(0, 0, 1)
			stats.NextRepeatAt = &nextRepeat
			stats.CalculatedAt = time.Now()

			if err := db.Create(stats).Error; err != nil {
				return fmt.Errorf("create attempt stats: %w", err)
			}

			// Generate repetitions based on attempts
			for _, attempt := range attemptedTasks {
				var ta models.TaskAttempt
				if err := db.Where("session_id = ? AND task_id = ?", session.ID, attempt.ID).
					First(&ta).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						continue
					}
					return fmt.Errorf("find task attempt: %w", err)
				}

				// Calculate repeat_at based on result
				var repeatAt time.Time
				switch ta.Result {
				case models.AnswerResultCorrect:
					repeatAt = time.Now().AddDate(0, 0, 3) // 3 days
				case models.AnswerResultPartial:
					repeatAt = time.Now().AddDate(0, 0, 2) // 2 days
				case models.AnswerResultWrong, models.AnswerResultSkipped:
					repeatAt = time.Now().AddDate(0, 0, 1) // 1 day
				}

				rep := &models.Repetition{
					ID:            uuid.New(),
					UserID:        user.ID,
					SessionID:     session.ID,
					TaskAttemptID: ta.ID,
					TaskID:        attempt.ID,
					TopicID:       attempt.TopicID,
					RepeatAt:      repeatAt,
					Status:        models.RepetitionStatusPlanned,
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				if err := db.Create(rep).Error; err != nil {
					return fmt.Errorf("create repetition: %w", err)
				}
			}

			log.Debug("created session with attempts",
				slog.String("session_id", session.ID.String()),
				slog.Int("attempts", len(attemptedTasks)),
			)
		}
	}

	log.Info("seeding completed",
		slog.Int("users", len(users)),
		slog.Int("subjects", len(subjects)),
		slog.Int("topics", len(topics)),
		slog.Int("tasks", len(tasks)),
	)

	return nil
}

func generateUser(index int) *models.User {
	email := fmt.Sprintf("user%d@example.com", index+1)
	name := fmt.Sprintf("Test User %d", index+1)
	password := fmt.Sprintf("password%d", index+1)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return &models.User{
		ID:           uuid.New(),
		Email:        email,
		Name:         name,
		Role:         models.UserRoleUser,
		PasswordHash: string(hashedPassword),
		Settings: models.UserSettings{
			Timezone:             "UTC",
			Language:             "ru-RU",
			Theme:                "light",
			NotificationsEnabled: true,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func generateSubject(userID uuid.UUID, index int) *models.Subject {
	subjects := []string{
		"Математика",
		"Физика",
		"Английский язык",
		"История",
		"Биология",
		"Химия",
		"Литература",
		"Информатика",
	}
	idx := index % len(subjects)
	title := fmt.Sprintf("%s (%d)", subjects[idx], index+1)

	descriptions := []string{
		"Подготовка к ОГЭ по предмету",
		"Заключительный курс перед экзаменом",
		"Интенсивная подготовка",
		"Базовый уровень",
		"Продвинутый уровень",
	}
	desc := descriptions[index%len(descriptions)]

	return &models.Subject{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       title,
		Description: desc,
		NumOfTopics: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func generateTopic(subjectID uuid.UUID, index int) *models.Topic {
	topics := []string{
		"Основные понятия",
		"Методология и теория",
		"Практические применения",
		"Типичные ошибки",
		"Расширенные темы",
		"Интеграция и синтез",
	}
	idx := index % len(topics)
	title := fmt.Sprintf("%s (Тема %d)", topics[idx], index+1)

	return &models.Topic{
		ID:        uuid.New(),
		SubjectID: subjectID,
		Title:     title,
		NumOfTasks: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func generateTask(topicID uuid.UUID, index int) *models.Task {
	taskTypes := []models.TaskType{
		models.TaskTypeFlashcard,
		models.TaskTypeTest,
		models.TaskTypeFillInTheBlank,
	}
	typeIdx := index % len(taskTypes)
	taskType := taskTypes[typeIdx]

	title := fmt.Sprintf("Задача %d: %s", index+1, []string{
		"Дайте определение",
		"Решите уравнение",
		"Переведите на английский",
		"Укажите правильный ответ",
		"Объясните явление",
	}[index%5])

	content := fmt.Sprintf("Содержание задачи номер %d. Требуется ответить на вопрос и объяснить решение.", index+1)

	return &models.Task{
		ID:        uuid.New(),
		TopicID:   topicID,
		Type:      taskType,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func generateSession(userID uuid.UUID, index int) *models.Session {
	now := time.Now()
	startedAt := now.AddDate(0, 0, -index) // offset by days
	endedAt := startedAt.Add(30 * time.Minute)

	return &models.Session{
		ID:        uuid.New(),
		UserID:    userID,
		SubjectID: nil,
		StartedAt: startedAt,
		EndedAt:   &endedAt,
		Notes:     nil,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func generateTaskAttempt(sessionID, userID, taskID uuid.UUID) *models.TaskAttempt {
	results := []models.AnswerResult{
		models.AnswerResultCorrect,
		models.AnswerResultWrong,
		models.AnswerResultPartial,
		models.AnswerResultSkipped,
	}
	resultIdx := rand.Intn(len(results))
	result := results[resultIdx]

	responseTimeMs := rand.Intn(60000) + 5000 // 5-65 seconds

	return &models.TaskAttempt{
		ID:             uuid.New(),
		SessionID:      sessionID,
		UserID:         userID,
		TaskID:         taskID,
		Result:         result,
		ResponseTimeMs: responseTimeMs,
		IsCorrect:      result == models.AnswerResultCorrect,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

