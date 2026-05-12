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
		&models.Card{},
		&models.StudySession{},
		&models.AttemptStat{},
		&models.Repetition{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	if err := gdb.Exec(createCalendarViewSQL).Error; err != nil {
		return fmt.Errorf("create calendar view: %w", err)
	}

	return nil
}

// calendar is a read-only projection over repetitions for date-based queries (no duplicated rows).
const createCalendarViewSQL = `
CREATE OR REPLACE VIEW calendar AS
SELECT
	r.id,
	r.user_id,
	r.task_id,
	r.card_id,
	r.next_due_at AS due_at,
	(r.next_due_at AT TIME ZONE 'UTC')::date AS calendar_date,
	r.status,
	r.interval_days,
	r.ease_factor,
	r.review_count,
	r.last_review_at
FROM repetitions r;
`
