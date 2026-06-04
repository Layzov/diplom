# Stage 8d: Demo Script — Complete Demo Flow

## 📋 Обзор

**Demo scripts** — это готовые bash/batch скрипты, которые демонстрируют полный цикл работы системы:

```
Register → Create Subject → Create Topic → Create Task → 
Start Session → Add Attempts → Complete Session → 
View Statistics → Check Repetitions
```

**Время выполнения:** ~15-20 секунд (со всеми output'ами)

---

## 🚀 Быстрый старт

### Linux / macOS / WSL

```bash
# 1. Убедиться что все запущено
docker compose up -d
go run ./cmd/seed
go run ./cmd/api

# 2. В отдельной сессии запустить demo
bash docs/demo.sh

# Со временной задержкой между запросами
bash docs/demo.sh -delay 2
```

### Windows (CMD/PowerShell)

```batch
REM 1. Убедиться что все запущено
docker compose up -d
go run ./cmd/seed
go run ./cmd/api

REM 2. В отдельной сессии запустить demo
docs\demo.bat
```

### Windows (WSL)

```bash
# Используй bash версию (как на Linux)
bash docs/demo.sh -delay 1
```

---

## 📂 Файлы

### `docs/demo.sh` (Unix/Linux/macOS/WSL)
- **Строк:** ~330
- **Зависимости:** bash, curl, jq
- **Параметры:**
  - `-delay N` — пауза N секунд между запросами (для показа)
  - `BASE_URL=...` — если API на другом адресе

### `docs/demo.bat` (Windows)
- **Строк:** ~250
- **Зависимости:** curl (встроен в Windows 10+), jq (опционально)
- **Параметры:**
  - `-delay N` — пауза N секунд между запросами

---

## 🎯 Что происходит в demo.sh

### 1️⃣ Health Check
```bash
GET /health
```
Проверяет что API запущен и готов к работе.

### 2️⃣ Register
```bash
POST /auth/register
{
  "email": "demo_<timestamp>@example.com",
  "password": "demo123456",
  "name": "Demo User"
}
```
Создает нового пользователя (уникальный email за счет timestamp).

### 3️⃣ Login
```bash
POST /auth/login
{
  "email": "...",
  "password": "demo123456"
}
```
Получает JWT token для всех последующих запросов.

### 4️⃣ Create Subject (Предмет)
```bash
POST /subjects (с токеном)
{
  "title": "Демонстрационная Математика",
  "description": "Полный цикл демонстрации..."
}
```

### 5️⃣ Create Topic (Тема)
```bash
POST /topics
{
  "subject_id": "<uuid>",
  "title": "Основные понятия и теория"
}
```

### 6️⃣ Create Task (Задача)
```bash
POST /tasks
{
  "topic_id": "<uuid>",
  "type": "practice",
  "title": "Решите задачу по геометрии",
  "content": "..."
}
```

### 7️⃣ Create Session (Начало учебной сессии)
```bash
POST /sessions
{
  "subject_id": "<uuid>",
  "started_at": "<ISO-8601>",
  "notes": "Demo study session"
}
```

### 8️⃣ Add Task Attempts (Попытки решения)
```bash
POST /task_attempts (2 раза: correct и partial)
{
  "session_id": "<uuid>",
  "task_id": "<uuid>",
  "result": "correct|partial",
  "time_ms": 15000|22000
}
```

### 9️⃣ Complete Session (Завершение сессии)
```bash
POST /sessions/<id>/complete
{
  "ended_at": "<ISO-8601>",
  "notes": "Session completed"
}
```

### 🔟 Get User Statistics
```bash
GET /me/stats
```
Показывает общую статистику: всего сессий, попыток, успеха.

### 1️⃣ 1️⃣ Get Topic Statistics
```bash
GET /me/stats/topics
```
Показывает статистику по темам.

### 1️⃣ 2️⃣ Get Upcoming Repetitions
```bash
GET /me/repetitions/upcoming
```
Показывает запланированные повторения (рассчитаны автоматически).

### 1️⃣ 3️⃣-1️⃣ 7️⃣ Get All Resources
- GET /me
- GET /subjects
- GET /topics?subject_id=...
- GET /tasks?topic_id=...
- GET /sessions

---

## 🎬 Использование для защиты ВКР

### Сценарий 1: Живая демонстрация

```bash
# Терминал 1: Запустить API
go run ./cmd/api

# Терминал 2: Запустить demo с задержками для показа
bash docs/demo.sh -delay 1

# На экране комиссия видит:
# - Каждый запрос и его результат
# - Красивый JSON output
# - Прогресс система в реальном времени
```

### Сценарий 2: Быстрое прохождение

```bash
# Без задержек (все за 3-5 сек)
bash docs/demo.sh

# Потом показываешь статистику:
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/me/stats | jq
```

### Сценарий 3: Кастомный flow

Модифицируй demo.sh для показа специфических feature'ов:

```bash
# Скопируй и отредактируй
cp docs/demo.sh docs/demo-custom.sh
# Добавь свои endpoints, удали ненужные
```

---

## 📊 Пример output'а

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎓 Exam Prep Backend — Complete Demo Flow
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

API URL: http://localhost:8080/api/v1
Delay: 1s between requests

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1️⃣  HEALTH CHECK
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

→ GET /health

{
  "status": "ok"
}

✅ API is healthy!

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
2️⃣  REGISTER NEW USER
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

→ POST /auth/register

Email: demo_1717530035@example.com

{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "email": "demo_1717530035@example.com",
    "name": "Demo User",
    "role": "user",
    "created_at": "2026-06-05T00:15:35.123Z"
  }
}

✅ User registered: 550e8400-e29b-41d4-a716-446655440001

...и так далее
```

---

## 🔧 Требования

### Linux/macOS/WSL
```bash
# Required
bash
curl
jq  # для форматирования JSON (sudo apt install jq)

# Install jq if not present
sudo apt install jq      # Debian/Ubuntu
brew install jq         # macOS
```

### Windows
```
curl           # встроен в Windows 10+
jq             # опционально, но рекомендуется
               # Скачай с https://github.com/jqlang/jq
```

---

## 🐛 Отладка

### Ошибка: `curl: (7) Failed to connect`
```
❌ API is not responding
```

**Решение:**
```bash
# Проверь что API запущен
go run ./cmd/api

# Проверь что на порту 8080 ничего не запущено
lsof -i :8080  # Linux/macOS
netstat -ano | findstr :8080  # Windows
```

### Ошибка: `jq: command not found`
```bash
# Windows: Установи jq
# https://github.com/jqlang/jq/releases

# Linux/macOS
sudo apt install jq  # Debian/Ubuntu
brew install jq     # macOS
```

### Ошибка: `Email already exists`
```
❌ Registration failed!
```

**Причина:** Email из предыдущего запуска еще существует.

**Решение:** Скрипт использует timestamp, поэтому просто запусти его снова. Или очисти БД:
```bash
docker compose down -v
docker compose up -d
go run ./cmd/seed
```

### Медленный demo (много пауз)
```bash
# Запусти без задержек
bash docs/demo.sh

# Или с меньшей задержкой
bash docs/demo.sh -delay 0.5
```

---

## 📝 Примеры модификаций

### Пример 1: Только основной цикл (регистрация + статистика)

```bash
# Скопируй demo.sh, удали шаги 4-8, оставь только:
# 1. Health
# 2. Register
# 3. Login
# 9. Get Stats
```

### Пример 2: Несколько пользователей

```bash
for i in {1..5}; do
  bash docs/demo.sh -delay 0.1
done
```

### Пример 3: Тест нагрузки (100 запросов)

```bash
for i in {1..100}; do
  curl -s http://localhost:8080/api/v1/health | jq .
done
```

---

## 🎯 Для защиты ВКР

**Рекомендуемый flow:**

1. **Подготовка (30 сек до защиты):**
   ```bash
   docker compose up -d
   go run ./cmd/seed -users=3
   go run ./cmd/api
   ```

2. **На защите (во время доклада):**
   ```bash
   # Запусти в отдельной сессии
   bash docs/demo.sh -delay 1
   
   # Комментируй что происходит:
   # "Вот регистрация пользователя..."
   # "Создаем предмет..."
   # "Начинаем учебную сессию..."
   # И т.д.
   ```

3. **Итоги:**
   - Комиссия видит работающую систему в реальном времени
   - JSON responses показывают структуру данных
   - Statistics доказывают что логика работает
   - Repetitions показывают алгоритм spaced repetition

**Эффект:** 💯 Комиссия влюбится в проект!

---

## 📌 Чек-лист перед защитой

- [ ] Docker запущен: `docker compose ps`
- [ ] PostgreSQL готов: `docker compose logs postgres`
- [ ] API запускается: `go run ./cmd/api` (проверь /health)
- [ ] Seed работает: `go run ./cmd/seed` (проверь статистику)
- [ ] Demo script работает: `bash docs/demo.sh`
- [ ] Интернет есть (если нужно curl с внешних API)
- [ ] Все зависимости установлены (jq, bash, curl)

---

## 🚀 Готово!

Stage 8d завершена! 🎉

**Дальше:**
- Stage 8a (OpenAPI/Swagger) — интерактивная документация
- Stage 9 (ВКР документация) — обязательно!
