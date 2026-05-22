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
		&models.TaskAttempt{},
		&models.AttemptStats{},
		&models.Repetition{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	if err := gdb.Exec(createCalendarViewSQL).Error; err != nil {
		return fmt.Errorf("create calendar view: %w", err)
	}

	return nil
}

// calendar lists upcoming repetitions only (status = planned).
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
