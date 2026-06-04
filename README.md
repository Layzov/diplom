# 🎓 Exam Prep Backend — Веб-приложение для подготовки к экзаменам

Backend-часть системы подготовки к экзаменам на основе метода интервальных повторений (spaced repetition).

---

## 🎯 Особенности

✅ **Метод интервальных повторений** — автоматический расчет дат повторений  
✅ **Учебные сессии** — трекирование попыток и результатов  
✅ **Статистика и аналитика** — прогресс по предметам и темам  
✅ **JWT аутентификация** — защищенные endpoints  
✅ **REST API** — легко интегрируется с фронтенд-приложением  
✅ **SQL views** — оптимизированные запросы для аналитики  
✅ **Индексы** — производительность на больших объемах данных  

---

## 🛠 Стек технологий

| Компонент | Технология |
|-----------|-----------|
| **Language** | Go 1.25 |
| **Router** | chi v5 |
| **ORM** | GORM v1.25 |
| **Database** | PostgreSQL 16 |
| **Auth** | JWT (golang-jwt) |
| **Containerization** | Docker & Docker Compose |
| **Config** | .env файлы |

---

## 📦 Архитектура

```
cmd/api/
  └─ main.go              # Точка входа приложения

internal/
  ├─ handler/             # HTTP handlers (endpoints)
  ├─ service/             # Бизнес-логика
  ├─ repository/          # Работа с БД
  ├─ models/              # GORM модели
  ├─ dto/                 # Request/Response DTO
  ├─ auth/                # JWT менеджер
  ├─ config/              # Конфигурация
  ├─ db/                  # Подключение и миграции
  ├─ apperror/            # Обработка ошибок
  └─ seeder/              # Генерация тестовых данных

pkg/
  ├─ middleware/          # HTTP middleware
  ├─ response/            # JSON ответы
  ├─ httputil/            # HTTP утилиты
  └─ handlers/            # Логирование

migrations/
  └─ *.sql                # SQL views и индексы
```

**Слои:**
```
HTTP Request
    ↓
Handler (тонкий слой, только валидация)
    ↓
Service (бизнес-логика, scheduler)
    ↓
Repository (SQL запросы)
    ↓
PostgreSQL
```

---

## 🚀 Быстрый старт

### 1️⃣ Требования

- Go 1.25+
- Docker & Docker Compose
- PostgreSQL 16 (запускается в Docker)

### 2️⃣ Подготовка

```bash
# Клонировать репозиторий
git clone <repo>
cd diplom.worktrees/agents-web-app-exam-prep-backend-setup

# Скопировать .env (если нужно)
cp .env.example .env

# Проверить переменные окружения
cat .env
```

**Обязательные переменные в .env:**
```env
APP_ENV=local
HTTP_ADDRESS=:8080

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=diplom
POSTGRES_PASSWORD=diplom
POSTGRES_DB=diplom

JWT_SECRET=your-secret-key-min-16-chars
JWT_ACCESS_TTL=24h
```

### 3️⃣ Запуск PostgreSQL

```bash
# Запустить контейнер
docker compose up -d

# Проверить статус
docker compose ps

# Просмотреть логи (если есть ошибки)
docker compose logs postgres
```

### 4️⃣ Заполнение БД тестовыми данными

```bash
# Со стандартными параметрами (5 юзеров, 2 предмета, 3 темы, 5 задач, 5 сессий)
go run ./cmd/seed

# Или с кастомными параметрами
go run ./cmd/seed -users=10 -subjects=3 -topics=5 -tasks=10 -sessions=10

# Подробнее: docs/seeding.md
```

### 5️⃣ Запуск сервера

```bash
go run ./cmd/api

# Успешный запуск выглядит так:
# INFO starting API env=local
# INFO database connected host=localhost db=diplom
# INFO database schema migrated
# INFO api routes registered prefix=/api/v1
# INFO http listening addr=:8080
```

### 6️⃣ Проверка здоровья

```bash
# Health check
curl http://localhost:8080/health
# {"status":"ok"}
```

---

## 📚 API Endpoints

### Публичные
- `GET /health` — Health check

### Аутентификация
- `POST /api/v1/auth/register` — Регистрация
- `POST /api/v1/auth/login` — Вход (выдает JWT token)

### Защищенные (требуют `Authorization: Bearer <token>`)

#### Профиль и статистика
- `GET /api/v1/me` — Получить профиль пользователя
- `GET /api/v1/me/stats` — Общая статистика
- `GET /api/v1/me/stats/topics` — Статистика по темам
- `GET /api/v1/me/stats/sessions` — Статистика по сессиям
- `GET /api/v1/me/repetitions/upcoming` — Ближайшие повторения

#### Предметы (CRUD)
- `GET /api/v1/subjects` — Список предметов
- `POST /api/v1/subjects` — Создать предмет
- `GET /api/v1/subjects/:id` — Получить предмет
- `PUT /api/v1/subjects/:id` — Обновить предмет
- `DELETE /api/v1/subjects/:id` — Удалить предмет

#### Темы (CRUD)
- `GET /api/v1/topics?subject_id=...` — Список тем
- `POST /api/v1/topics` — Создать тему
- `GET /api/v1/topics/:id` — Получить тему
- `PUT /api/v1/topics/:id` — Обновить тему
- `DELETE /api/v1/topics/:id` — Удалить тему

#### Задачи (CRUD)
- `GET /api/v1/tasks?topic_id=...` — Список задач
- `POST /api/v1/tasks` — Создать задачу
- `GET /api/v1/tasks/:id` — Получить задачу
- `PUT /api/v1/tasks/:id` — Обновить задачу
- `DELETE /api/v1/tasks/:id` — Удалить задачу

#### Сессии обучения
- `GET /api/v1/sessions` — Список сессий
- `POST /api/v1/sessions` — Создать сессию
- `GET /api/v1/sessions/:id` — Получить сессию
- `POST /api/v1/sessions/:id/complete` — Завершить сессию

**Полная документация:** `docs/api.md`

---

## 🔐 Аутентификация

### Регистрация

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123",
    "name": "John Doe"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

### Вход

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400
  }
}
```

### Использование токена

```bash
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 📊 Алгоритм планировщика повторений

**Правило:** Дата повторения зависит от результата попытки.

| Результат | Дней до повторения |
|-----------|---------|
| ✅ Correct | 3 дня |
| ⚠️ Partial | 2 дня |
| ❌ Wrong | 1 день |
| ⏭️ Skipped | 1 день |

Реализация: `internal/service/scheduler.go`

---

## 🧪 Тестирование

### Health Check
```bash
curl http://localhost:8080/health
```

### Пример Flow: от регистрации до статистики

```bash
#!/bin/bash

# 1. Регистрация
REGISTER=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "testpass123",
    "name": "Test User"
  }')
echo "Registered: $REGISTER"

# 2. Вход
LOGIN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "testpass123"
  }')
TOKEN=$(echo $LOGIN | jq -r '.data.token')
echo "Token: $TOKEN"

# 3. Создать предмет
SUBJECT=$(curl -s -X POST http://localhost:8080/api/v1/subjects \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Математика",
    "description": "Подготовка к ОГЭ"
  }')
SUBJECT_ID=$(echo $SUBJECT | jq -r '.data.id')
echo "Created Subject: $SUBJECT_ID"

# 4. Получить статистику
STATS=$(curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/me/stats)
echo "Stats: $STATS"
```

**Сохранить в `test-flow.sh` и запустить:**
```bash
bash test-flow.sh
```

---

## 📈 Производительность

### SQL Views для аналитики
- `calendar` — Календарь повторений
- `user_learning_stats` — Общая статистика пользователя
- `topic_learning_stats` — Статистика по темам
- `session_learning_stats` — Статистика по сессиям

### Индексы
```
idx_subjects_user_created_at          — Быстрый поиск предметов по пользователю
idx_topics_subject_created_at         — Быстрый поиск тем по предмету
idx_sessions_user_started_at          — Быстрый поиск сессий
idx_repetitions_planned_user_due      — Ближайшие повторения
idx_task_attempts_user_task           — История попыток
```

Все индексы создаются автоматически при запуске `db.Migrate()`.

---

## 🐛 Отладка

### Логирование

Приложение выводит подробные логи в зависимости от `APP_ENV`:
- `APP_ENV=local` — Pretty печать, DEBUG уровень
- `APP_ENV=dev` — JSON логи, DEBUG уровень
- `APP_ENV=prod` — JSON логи, INFO уровень

### Проверка БД

```bash
# Подключиться к БД через docker
docker compose exec postgres psql -U diplom -d diplom

# SQL запросы в REPL
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM subjects;
SELECT COUNT(*) FROM sessions;
```

### Очистка данных

```bash
# Остановить контейнер и удалить volume
docker compose down -v

# Запустить заново
docker compose up -d
go run ./cmd/seed
```

---

## 📁 Файлы конфигурации

| Файл | Назначение |
|------|-----------|
| `.env.example` | Шаблон переменных окружения |
| `.env` | Локальные переменные (git-ignored) |
| `docker-compose.yml` | Конфигурация Docker сервисов |
| `go.mod` | Go зависимости |
| `go.sum` | Контрольные суммы зависимостей |

---

## 📚 Дополнительная документация

- **API Документация:** `docs/api.md`
- **Аутентификация:** `docs/auth.md`
- **Seed Скрипты:** `docs/seeding.md`
- **Производительность:** `docs/performance.md`

---

## 🚀 Деплой

### На сервер (пример)

```bash
# 1. Собрать бинарник
CGO_ENABLED=1 GOOS=linux go build -o api ./cmd/api

# 2. Создать .env на сервере
cat > .env << EOF
APP_ENV=prod
POSTGRES_HOST=db.example.com
POSTGRES_USER=produser
...
EOF

# 3. Запустить
./api
```

### С Docker

```bash
# Создать Dockerfile
cat > Dockerfile << 'EOF'
FROM golang:1.25-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o api ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/api .
CMD ["./api"]
EOF

# Собрать и запустить
docker build -t exam-prep-api .
docker run -p 8080:8080 --env-file .env exam-prep-api
```

---

## 📞 Контакты и помощь

- **Issues:** GitHub Issues (если есть репо на GitHub)
- **Documentation:** `docs/` папка
- **Logs:** Проверь вывод `go run ./cmd/api`

---

## 📜 Лицензия

MIT License — используй свободно!

---

**Happy coding! 🚀**
