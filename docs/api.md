# API v1

Базовый префикс: `/api/v1`

Ошибки: `{"error":{"code":"...","message":"..."}}`

**Аутентификация (этап 5.5):** почти все эндпоинты требуют `Authorization: Bearer <token>`.  
Подробнее: [auth.md](./auth.md).

| Method | Path |
|--------|------|
| POST | `/auth/register` |
| POST | `/auth/login` |

## Me (текущий пользователь)

| Method | Path |
|--------|------|
| GET/PUT | `/me` |
| GET/POST | `/me/subjects` |
| GET | `/me/sessions`, `/me/repetitions`, `/me/calendar` |
| GET | `/me/stats`, `/me/stats/topics`, `/me/stats/sessions`, `/me/stats/upcoming` |

## Subjects

| Method | Path |
|--------|------|
| POST | `/me/subjects` |
| GET | `/me/subjects` |
| GET | `/subjects/{id}` |
| PUT | `/subjects/{id}` |
| DELETE | `/subjects/{id}` |

## Topics

| Method | Path |
|--------|------|
| POST | `/subjects/{id}/topics` |
| GET | `/subjects/{id}/topics` |
| GET | `/topics/{id}` |
| PUT | `/topics/{id}` |
| DELETE | `/topics/{id}` |

## Tasks

| Method | Path |
|--------|------|
| POST | `/topics/{id}/tasks` |
| GET | `/topics/{id}/tasks` |
| GET | `/tasks/{id}` |
| PUT | `/tasks/{id}` |
| DELETE | `/tasks/{id}` |

```json
// POST /topics/{id}/tasks
{
  "title": "Определение производной",
  "content": "f'(x) = ...",
  "type": "flashcard"
}
```

Типы заданий: `flashcard`, `test`, `short_answer`, `fill_in_the_blank`, `matching`, `manual_review`.

---

## Этап 4 — учебный цикл

### Сессии

| Method | Path | Описание |
|--------|------|----------|
| POST | `/sessions` | Начать сессию |
| GET | `/sessions/{id}` | Сессия |
| POST | `/sessions/{id}/finish` | Завершить, посчитать `attempt_stats` |
| GET | `/users/{id}/sessions` | Список сессий пользователя |

```json
// POST /sessions
{
  "user_id": "uuid",
  "subject_id": "uuid (optional)",
  "notes": "optional"
}
```

### Попытки (task_attempts + repetition)

| Method | Path |
|--------|------|
| POST | `/sessions/{id}/attempts` |
| GET | `/sessions/{id}/attempts` |
| GET | `/sessions/{id}/stats` | после `finish` |

```json
// POST /sessions/{id}/attempts
{
  "task_id": "uuid",
  "answer_text": "ответ",
  "is_correct": true,
  "result": "correct",
  "response_time_ms": 1500
}
```

`result`: `correct`, `partial`, `wrong`, `skipped`.

Интервалы повторения (упрощённо): correct → 3 дня, partial → 2, wrong/skipped → 1.

### Повторения и календарь

| Method | Path |
|--------|------|
| GET | `/me/repetitions?status=planned` |
| GET | `/me/calendar?from=RFC3339&to=RFC3339` |
| PATCH | `/repetitions/{id}` |

```json
// PATCH /repetitions/{id}
{ "status": "done" }
```

Календарь — SQL view `calendar` (только `status = planned`).

---

## Этап 5 — статистика

SQL views: `user_learning_stats`, `topic_learning_stats`, `session_learning_stats`.

| Method | Path | Описание |
|--------|------|----------|
| GET | `/me/stats` | Сводка по пользователю |
| GET | `/me/stats/topics` | Агрегаты по темам |
| GET | `/me/stats/sessions` | Сессии + attempt_stats |
| GET | `/me/stats/upcoming` | Ближайшие повторения |

Query для upcoming:
- `limit` — по умолчанию 10, макс. 100
- `include_overdue` — `true` / `false` (по умолчанию `false`, только `due_at >= now`)

```json
// GET /me/stats — пример ответа
{
  "user_id": "...",
  "total_sessions": 5,
  "finished_sessions": 4,
  "active_sessions": 1,
  "total_attempts": 42,
  "correct_count": 30,
  "overall_success_rate": 71.43,
  "planned_repetitions": 8,
  "overdue_repetitions": 2
}
```

---

Индексы и проверка планов запросов (этап 6): [performance.md](./performance.md).
