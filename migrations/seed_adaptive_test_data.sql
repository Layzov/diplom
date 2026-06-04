-- Adaptive Learning Test Data Script
-- This script populates the database with test data demonstrating the adaptive repetition system

BEGIN TRANSACTION;

-- ============================================================================
-- Step 1: Create test user
-- ============================================================================
INSERT INTO users (id, email, password_hash, full_name, created_at, updated_at)
VALUES (
    'a1234567-1234-5678-9012-345678901234'::uuid,
    'testadaptive@exam.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36XQd/Xa', -- bcrypt of 'password'
    'Adaptive Learner',
    NOW(),
    NOW()
);

-- ============================================================================
-- Step 2: Create subject (Math)
-- ============================================================================
INSERT INTO subjects (id, user_id, title, description, created_at, updated_at)
VALUES (
    'b1234567-1234-5678-9012-345678901234'::uuid,
    'a1234567-1234-5678-9012-345678901234'::uuid,
    'Math Basics',
    'Mathematics fundamentals for exam prep',
    NOW(),
    NOW()
);

-- ============================================================================
-- Step 3: Create topic (Algebra)
-- ============================================================================
INSERT INTO topics (id, subject_id, title, description, created_at, updated_at)
VALUES (
    'c1234567-1234-5678-9012-345678901234'::uuid,
    'b1234567-1234-5678-9012-345678901234'::uuid,
    'Algebra Basics',
    'Linear equations and basics',
    NOW(),
    NOW()
);

-- ============================================================================
-- Step 4: Create 5 tasks (simulating EGE-style problems)
-- ============================================================================
INSERT INTO tasks (id, topic_id, title, description, type, difficulty, estimated_time_minutes, created_at, updated_at)
VALUES
    ('d1111111-1111-1111-1111-111111111111'::uuid, 'c1234567-1234-5678-9012-345678901234'::uuid, 'Task 1: Solve simple equation', 'Solve: x + 5 = 10', 'multiple_choice', 20, 2, NOW(), NOW()),
    ('d2222222-2222-2222-2222-222222222222'::uuid, 'c1234567-1234-5678-9012-345678901234'::uuid, 'Task 2: Linear functions', 'Find slope of y=2x+3', 'multiple_choice', 40, 4, NOW(), NOW()),
    ('d3333333-3333-3333-3333-333333333333'::uuid, 'c1234567-1234-5678-9012-345678901234'::uuid, 'Task 3: System of equations', 'Solve x+y=5, x-y=1', 'multiple_choice', 60, 6, NOW(), NOW()),
    ('d4444444-4444-4444-4444-444444444444'::uuid, 'c1234567-1234-5678-9012-345678901234'::uuid, 'Task 4: Quadratic equations (STRUGGLING)', 'Solve x²-5x+6=0', 'multiple_choice', 80, 8, NOW(), NOW()),
    ('d5555555-5555-5555-5555-555555555555'::uuid, 'c1234567-1234-5678-9012-345678901234'::uuid, 'Task 5: Complex expressions (STRUGGLING)', 'Simplify (x+2)²-(x-2)²', 'multiple_choice', 100, 10, NOW(), NOW());

-- ============================================================================
-- Step 5: Create LEARNING SESSION (first session)
-- ============================================================================
INSERT INTO study_sessions (id, user_id, topic_id, task_count, mode, status, started_at, created_at, updated_at)
VALUES (
    'e1111111-1111-1111-1111-111111111111'::uuid,
    'a1234567-1234-5678-9012-345678901234'::uuid,
    'c1234567-1234-5678-9012-345678901234'::uuid,
    5,
    'learning',
    'done',
    NOW() - INTERVAL '2 hours',
    NOW() - INTERVAL '2 hours',
    NOW() - INTERVAL '1.5 hours'
);

-- ============================================================================
-- Step 6: Add all 5 tasks to learning session
-- ============================================================================
INSERT INTO session_tasks (id, session_id, task_id, created_at)
VALUES
    ('f1111111-1111-1111-1111-111111111111'::uuid, 'e1111111-1111-1111-1111-111111111111'::uuid, 'd1111111-1111-1111-1111-111111111111'::uuid, NOW()),
    ('f2222222-2222-2222-2222-222222222222'::uuid, 'e1111111-1111-1111-1111-111111111111'::uuid, 'd2222222-2222-2222-2222-222222222222'::uuid, NOW()),
    ('f3333333-3333-3333-3333-333333333333'::uuid, 'e1111111-1111-1111-1111-111111111111'::uuid, 'd3333333-3333-3333-3333-333333333333'::uuid, NOW()),
    ('f4444444-4444-4444-4444-444444444444'::uuid, 'e1111111-1111-1111-1111-111111111111'::uuid, 'd4444444-4444-4444-4444-444444444444'::uuid, NOW()),
    ('f5555555-5555-5555-5555-555555555555'::uuid, 'e1111111-1111-1111-1111-111111111111'::uuid, 'd5555555-5555-5555-5555-555555555555'::uuid, NOW());

-- ============================================================================
-- Step 7: Record attempts - CORRECT answers for tasks 1,2,3 and WRONG for 4,5
-- ============================================================================
-- Task 1: CORRECT
INSERT INTO task_attempts (id, session_id, task_id, user_answer, answer_result, submitted_at, created_at)
VALUES (
    'g1111111-1111-1111-1111-111111111111'::uuid,
    'e1111111-1111-1111-1111-111111111111'::uuid,
    'd1111111-1111-1111-1111-111111111111'::uuid,
    '5',
    'correct',
    NOW() - INTERVAL '1.9 hours',
    NOW() - INTERVAL '1.9 hours'
);

-- Task 2: CORRECT
INSERT INTO task_attempts (id, session_id, task_id, user_answer, answer_result, submitted_at, created_at)
VALUES (
    'g2222222-2222-2222-2222-222222222222'::uuid,
    'e1111111-1111-1111-1111-111111111111'::uuid,
    'd2222222-2222-2222-2222-222222222222'::uuid,
    '2',
    'correct',
    NOW() - INTERVAL '1.8 hours',
    NOW() - INTERVAL '1.8 hours'
);

-- Task 3: CORRECT
INSERT INTO task_attempts (id, session_id, task_id, user_answer, answer_result, submitted_at, created_at)
VALUES (
    'g3333333-3333-3333-3333-333333333333'::uuid,
    'e1111111-1111-1111-1111-111111111111'::uuid,
    'd3333333-3333-3333-3333-333333333333'::uuid,
    'x=3, y=2',
    'correct',
    NOW() - INTERVAL '1.7 hours',
    NOW() - INTERVAL '1.7 hours'
);

-- Task 4: WRONG (STRUGGLING)
INSERT INTO task_attempts (id, session_id, task_id, user_answer, answer_result, submitted_at, created_at)
VALUES (
    'g4444444-4444-4444-4444-444444444444'::uuid,
    'e1111111-1111-1111-1111-111111111111'::uuid,
    'd4444444-4444-4444-4444-444444444444'::uuid,
    'wrong_answer',
    'wrong',
    NOW() - INTERVAL '1.6 hours',
    NOW() - INTERVAL '1.6 hours'
);

-- Task 5: WRONG (STRUGGLING)
INSERT INTO task_attempts (id, session_id, task_id, user_answer, answer_result, submitted_at, created_at)
VALUES (
    'g5555555-5555-5555-5555-555555555555'::uuid,
    'e1111111-1111-1111-1111-111111111111'::uuid,
    'd5555555-5555-5555-5555-555555555555'::uuid,
    'incomplete',
    'wrong',
    NOW() - INTERVAL '1.5 hours',
    NOW() - INTERVAL '1.5 hours'
);

-- ============================================================================
-- Step 8: Create repetitions ONLY for failed attempts (key to adaptive system!)
-- ============================================================================
-- NO REPETITION for tasks 1,2,3 (correct answers) - that's the trick!

-- Repetition for Task 4 (WRONG - creates repetition with quality='struggling')
INSERT INTO repetitions (id, user_id, task_id, quality, status, planned_for, created_at, updated_at)
VALUES (
    'h4444444-4444-4444-4444-444444444444'::uuid,
    'a1234567-1234-5678-9012-345678901234'::uuid,
    'd4444444-4444-4444-4444-444444444444'::uuid,
    'struggling', -- Quality = STRUGGLING
    'planned',
    NOW() + INTERVAL '1 day', -- Review in 1 day
    NOW(),
    NOW()
);

-- Repetition for Task 5 (WRONG - creates repetition with quality='struggling')
INSERT INTO repetitions (id, user_id, task_id, quality, status, planned_for, created_at, updated_at)
VALUES (
    'h5555555-5555-5555-5555-555555555555'::uuid,
    'a1234567-1234-5678-9012-345678901234'::uuid,
    'd5555555-5555-5555-5555-555555555555'::uuid,
    'struggling', -- Quality = STRUGGLING
    'planned',
    NOW() + INTERVAL '1 day', -- Review in 1 day
    NOW(),
    NOW()
);

-- ============================================================================
-- Step 9: Create REVIEW SESSION (adaptive - should get tasks 4,5 first!)
-- ============================================================================
INSERT INTO study_sessions (id, user_id, topic_id, task_count, mode, status, started_at, created_at, updated_at)
VALUES (
    'e2222222-2222-2222-2222-222222222222'::uuid,
    'a1234567-1234-5678-9012-345678901234'::uuid,
    'c1234567-1234-5678-9012-345678901234'::uuid,
    3, -- Only want 3 tasks for review
    'review', -- REVIEW mode = adaptive selection!
    'done',
    NOW() - INTERVAL '30 minutes',
    NOW() - INTERVAL '30 minutes',
    NOW() - INTERVAL '10 minutes'
);

-- ============================================================================
-- Step 10: Add STRUGGLING tasks to review session (should be 4,5!)
-- ============================================================================
-- These will be selected by ListForReview() query which prioritizes quality='struggling'
INSERT INTO session_tasks (id, session_id, task_id, created_at)
VALUES
    ('f4444444-4444-4444-4444-444444444444'::uuid, 'e2222222-2222-2222-2222-222222222222'::uuid, 'd4444444-4444-4444-4444-444444444444'::uuid, NOW()),
    ('f5555555-5555-5555-5555-555555555555'::uuid, 'e2222222-2222-2222-2222-222222222222'::uuid, 'd5555555-5555-5555-5555-555555555555'::uuid, NOW()),
    ('f3333333-3333-3333-3333-333333333333'::uuid, 'e2222222-2222-2222-2222-222222222222'::uuid, 'd3333333-3333-3333-3333-333333333333'::uuid, NOW()); -- 3rd best task as fallback

-- ============================================================================
-- Summary: This test data demonstrates:
-- ============================================================================
-- 1. LEARNING session created all 5 tasks (mode='learning')
-- 2. User answered CORRECTLY on tasks 1,2,3 (NO repetitions created)
-- 3. User answered WRONG on tasks 4,5 (repetitions created with quality='struggling')
-- 4. REVIEW session (mode='review') got tasks 4,5 FIRST (adaptive prioritization!)
-- 5. Task-level stats will show tasks 4,5 as "struggling" (difficulty identification)
--
-- Key adaptive features:
-- - Only wrong answers → repetitions (reduces noise)
-- - quality='struggling' identifies problem areas
-- - Review mode queries prioritize by quality (struggling > learning > mastered)
-- - SQL view task_learning_stats shows per-task difficulty metrics

COMMIT;

-- ============================================================================
-- Verification queries
-- ============================================================================

-- View all repetitions (should see only tasks 4,5 with quality='struggling')
SELECT id, task_id, quality, status, planned_for 
FROM repetitions 
WHERE user_id = 'a1234567-1234-5678-9012-345678901234'::uuid
ORDER BY quality;

-- View tasks in learning session (should see all 5)
SELECT st.task_id, t.title, t.difficulty
FROM session_tasks st
JOIN tasks t ON t.id = st.task_id
WHERE st.session_id = 'e1111111-1111-1111-1111-111111111111'::uuid
ORDER BY t.difficulty;

-- View tasks in review session (should see 4,5,3 - struggling first!)
SELECT st.task_id, t.title, t.difficulty
FROM session_tasks st
JOIN tasks t ON t.id = st.task_id
WHERE st.session_id = 'e2222222-2222-2222-2222-222222222222'::uuid
ORDER BY t.difficulty DESC;

-- View attempt statistics (shows mastery levels)
SELECT 
    t.id,
    t.title,
    COUNT(DISTINCT ta.id) as attempt_count,
    SUM(CASE WHEN ta.answer_result = 'correct' THEN 1 ELSE 0 END)::float / COUNT(DISTINCT ta.id) as success_rate,
    COUNT(CASE WHEN r.quality = 'struggling' THEN 1 END) as struggling_count,
    COUNT(CASE WHEN r.quality = 'learning' THEN 1 END) as learning_count,
    COUNT(CASE WHEN r.quality = 'mastered' THEN 1 END) as mastered_count
FROM tasks t
LEFT JOIN task_attempts ta ON t.id = ta.task_id
LEFT JOIN repetitions r ON t.id = r.task_id
WHERE t.topic_id = 'c1234567-1234-5678-9012-345678901234'::uuid
GROUP BY t.id, t.title
ORDER BY t.difficulty;
