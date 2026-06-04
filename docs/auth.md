# Этап 5.5 — JWT аутентификация

## Настройка

В `.env`:

```
JWT_SECRET=change-me-to-a-long-random-secret   # минимум 16 символов
JWT_ACCESS_TTL=24h
```

## Публичные эндпоинты

| Method | Path |
|--------|------|
| POST | `/api/v1/auth/register` |
| POST | `/api/v1/auth/login` |
| GET | `/health` |

## Ответ login / register

```json
{
  "access_token": "eyJhbG...",
  "expires_at": "2026-05-23T12:00:00Z",
  "user": { "id": "...", "email": "...", "name": "...", "role": "user", ... }
}
```

## Защищённые запросы

Заголовок:

```
Authorization: Bearer <access_token>
```

## Профиль и «свои» данные (`/me`)

| Method | Path |
|--------|------|
| GET | `/api/v1/me` |
| PUT | `/api/v1/me` |
| GET/POST | `/api/v1/me/subjects` |
| GET | `/api/v1/me/sessions` |
| GET | `/api/v1/me/repetitions` |
| GET | `/api/v1/me/calendar` |
| GET | `/api/v1/me/stats`, `/stats/topics`, `/stats/sessions`, `/stats/upcoming` |

CRUD по `subjects`, `topics`, `tasks`, `sessions` — только для ресурсов **владельца** (проверка в service).

`GET /api/v1/users` — только роль **admin**.

## Пример (PowerShell)

```powershell
$auth = Invoke-RestMethod -Method Post http://localhost:8080/api/v1/auth/login `
  -ContentType application/json `
  -Body '{"email":"you@example.com","password":"password1"}'

$headers = @{ Authorization = "Bearer $($auth.access_token)" }

Invoke-RestMethod http://localhost:8080/api/v1/me/stats -Headers $headers
```

## Breaking changes

- `POST /api/v1/users` убран — регистрация через `/auth/register`.
- `POST /api/v1/sessions` — без `user_id` в теле, пользователь из токена.
- Старые пути `/users/{id}/...` заменены на `/me/...` для личных данных.
