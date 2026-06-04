-- Stage 6: indexes for FK lookups, filters, and sort columns used by API/views.
-- Applied on startup from internal/db/migrate.go (idempotent).

-- CRUD lists
CREATE INDEX IF NOT EXISTS idx_subjects_user_created_at
    ON subjects (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_topics_subject_created_at
    ON topics (subject_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_tasks_topic_created_at
    ON tasks (topic_id, created_at DESC);

-- Sessions & attempts
CREATE INDEX IF NOT EXISTS idx_sessions_user_started_at
    ON sessions (user_id, started_at DESC);

CREATE INDEX IF NOT EXISTS idx_task_attempts_session_created_at
    ON task_attempts (session_id, created_at);

CREATE INDEX IF NOT EXISTS idx_task_attempts_user_task
    ON task_attempts (user_id, task_id);

-- Repetitions: calendar / upcoming / list
CREATE INDEX IF NOT EXISTS idx_repetitions_planned_user_due
    ON repetitions (user_id, repeat_at)
    WHERE status = 'planned';

CREATE INDEX IF NOT EXISTS idx_repetitions_session_planned_due
    ON repetitions (session_id, repeat_at)
    WHERE status = 'planned';

CREATE INDEX IF NOT EXISTS idx_repetitions_user_status_due
    ON repetitions (user_id, status, repeat_at);

CREATE INDEX IF NOT EXISTS idx_repetitions_user_topic
    ON repetitions (user_id, topic_id);
