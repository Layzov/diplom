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
