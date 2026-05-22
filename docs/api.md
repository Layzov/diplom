# API v1 (этап 3 — CRUD)

Базовый префикс: `/api/v1`

Ошибки: `{"error":{"code":"...","message":"..."}}`

## Users

| Method | Path | Описание |
|--------|------|----------|
| POST | `/users` | Создать пользователя |
| GET | `/users` | Список |
| GET | `/users/{id}` | По id |
| PUT | `/users/{id}` | Обновить |
| DELETE | `/users/{id}` | Удалить |

```json
// POST /users
{
  "email": "student@example.com",
  "name": "Иван",
  "password": "secret123"
}
```

## Subjects

| Method | Path |
|--------|------|
| POST | `/users/{id}/subjects` |
| GET | `/users/{id}/subjects` |
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
