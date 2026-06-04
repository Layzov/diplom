package db

import (
	"fmt"

	"gorm.io/gorm"
)

// indexStatements are applied after AutoMigrate (partial indexes are not expressible in GORM tags).
// Keep in sync with migrations/0003_indexes.sql.
var indexStatements = []string{
	`CREATE INDEX IF NOT EXISTS idx_subjects_user_created_at
	 ON subjects (user_id, created_at DESC)`,
	`CREATE INDEX IF NOT EXISTS idx_topics_subject_created_at
	 ON topics (subject_id, created_at DESC)`,
	`CREATE INDEX IF NOT EXISTS idx_tasks_topic_created_at
	 ON tasks (topic_id, created_at DESC)`,
	`CREATE INDEX IF NOT EXISTS idx_sessions_user_started_at
	 ON sessions (user_id, started_at DESC)`,
	`CREATE INDEX IF NOT EXISTS idx_task_attempts_session_created_at
	 ON task_attempts (session_id, created_at)`,
	`CREATE INDEX IF NOT EXISTS idx_task_attempts_user_task
	 ON task_attempts (user_id, task_id)`,
	`CREATE INDEX IF NOT EXISTS idx_repetitions_planned_user_due
	 ON repetitions (user_id, repeat_at)
	 WHERE status = 'planned'`,
	`CREATE INDEX IF NOT EXISTS idx_repetitions_session_planned_due
	 ON repetitions (session_id, repeat_at)
	 WHERE status = 'planned'`,
	`CREATE INDEX IF NOT EXISTS idx_repetitions_user_status_due
	 ON repetitions (user_id, status, repeat_at)`,
	`CREATE INDEX IF NOT EXISTS idx_repetitions_user_topic
	 ON repetitions (user_id, topic_id)`,
}

func applyIndexes(gdb *gorm.DB) error {
	for i, stmt := range indexStatements {
		if err := gdb.Exec(stmt).Error; err != nil {
			return fmt.Errorf("index statement %d: %w", i+1, err)
		}
	}
	return nil
}
