-- Optional: apply manually if you manage schema without AutoMigrate.
-- The app runs the same statement from internal/db/migrate.go on startup.

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
