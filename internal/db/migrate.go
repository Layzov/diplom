package db

import (
	"fmt"

	"diplom/internal/models"

	"gorm.io/gorm"
)

// Migrate applies schema changes for local/dev bootstrap: GORM AutoMigrate plus DB views.
func Migrate(gdb *gorm.DB) error {
	if err := gdb.AutoMigrate(
		&models.User{},
		&models.Subject{},
		&models.Topic{},
		&models.Task{},
		&models.Attachment{},
		&models.Session{},
		&models.SessionTask{},
		&models.TaskAttempt{},
		&models.AttemptStats{},
		&models.Repetition{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	views := []struct {
		name string
		sql  string
	}{
		{"calendar", createCalendarViewSQL},
		{"user_learning_stats", createUserLearningStatsViewSQL},
		{"topic_learning_stats", createTopicLearningStatsViewSQL},
		{"session_learning_stats", createSessionLearningStatsViewSQL},
	}
	for _, v := range views {
		if err := gdb.Exec(v.sql).Error; err != nil {
			return fmt.Errorf("create view %s: %w", v.name, err)
		}
	}

	if err := applyIndexes(gdb); err != nil {
		return fmt.Errorf("apply indexes: %w", err)
	}

	return nil
}

const createCalendarViewSQL = `
CREATE OR REPLACE VIEW calendar AS
SELECT
	r.id,
	r.user_id,
	r.session_id,
	r.task_attempt_id,
	r.task_id,
	r.topic_id,
	r.repeat_at AS due_at,
	(r.repeat_at AT TIME ZONE 'UTC')::date AS calendar_date,
	r.status,
	r.created_at,
	r.updated_at
FROM repetitions r
WHERE r.status = 'planned';
`

const createUserLearningStatsViewSQL = `
CREATE OR REPLACE VIEW user_learning_stats AS
SELECT
	u.id AS user_id,
	COUNT(DISTINCT s.id) AS total_sessions,
	COUNT(DISTINCT s.id) FILTER (WHERE s.ended_at IS NOT NULL) AS finished_sessions,
	COUNT(DISTINCT s.id) FILTER (WHERE s.ended_at IS NULL) AS active_sessions,
	COALESCE(SUM(ast.total_tasks), 0)::bigint AS total_attempts,
	COALESCE(SUM(ast.correct_count), 0)::bigint AS correct_count,
	COALESCE(SUM(ast.wrong_count), 0)::bigint AS wrong_count,
	COALESCE(SUM(ast.partial_count), 0)::bigint AS partial_count,
	COALESCE(SUM(ast.skipped_count), 0)::bigint AS skipped_count,
	CASE
		WHEN COALESCE(SUM(ast.total_tasks), 0) > 0 THEN
			ROUND(SUM(ast.correct_count)::numeric / SUM(ast.total_tasks) * 100, 2)
		ELSE 0
	END AS overall_success_rate,
	COUNT(r.id) FILTER (WHERE r.status = 'planned') AS planned_repetitions,
	COUNT(r.id) FILTER (WHERE r.status = 'planned' AND r.repeat_at < NOW()) AS overdue_repetitions
FROM users u
LEFT JOIN sessions s ON s.user_id = u.id
LEFT JOIN attempt_stats ast ON ast.session_id = s.id
LEFT JOIN repetitions r ON r.user_id = u.id
GROUP BY u.id;
`

const createTopicLearningStatsViewSQL = `
CREATE OR REPLACE VIEW topic_learning_stats AS
SELECT
	sub.user_id,
	t.id AS topic_id,
	t.subject_id,
	t.title AS topic_title,
	sub.title AS subject_title,
	COUNT(ta.id) AS attempt_count,
	COUNT(ta.id) FILTER (WHERE ta.result = 'correct') AS correct_count,
	COUNT(ta.id) FILTER (WHERE ta.result = 'wrong') AS wrong_count,
	COUNT(ta.id) FILTER (WHERE ta.result = 'partial') AS partial_count,
	COUNT(ta.id) FILTER (WHERE ta.result = 'skipped') AS skipped_count,
	CASE
		WHEN COUNT(ta.id) > 0 THEN
			ROUND(COUNT(ta.id) FILTER (WHERE ta.result = 'correct')::numeric / COUNT(ta.id) * 100, 2)
		ELSE 0
	END AS success_rate,
	COUNT(rep.id) FILTER (WHERE rep.status = 'planned') AS planned_repetitions
FROM topics t
JOIN subjects sub ON sub.id = t.subject_id
LEFT JOIN tasks tk ON tk.topic_id = t.id
LEFT JOIN task_attempts ta ON ta.task_id = tk.id AND ta.user_id = sub.user_id
LEFT JOIN repetitions rep ON rep.topic_id = t.id AND rep.user_id = sub.user_id
GROUP BY sub.user_id, t.id, t.subject_id, t.title, sub.title;
`

const createSessionLearningStatsViewSQL = `
CREATE OR REPLACE VIEW session_learning_stats AS
SELECT
	s.id AS session_id,
	s.user_id,
	s.subject_id,
	s.started_at,
	s.ended_at,
	ast.id AS stats_id,
	ast.total_tasks,
	ast.correct_count,
	ast.wrong_count,
	ast.partial_count,
	ast.skipped_count,
	ast.total_time_ms,
	ast.average_time_ms,
	ast.success_rate,
	ast.next_repeat_at,
	ast.calculated_at
FROM sessions s
LEFT JOIN attempt_stats ast ON ast.session_id = s.id;
`
