# 🎯 Stage 7a: Seed Scripts — ЗАВЕРШЕНО ✅

## 📋 Что было сделано

### 1. Создан пакет `internal/seeder`
```
internal/seeder/seeder.go
├─ Config struct (параметризация)
├─ Seed() func (основная функция)
└─ Генераторы: generateUser(), generateSubject(), generateTopic(), generateTask(), generateSession(), generateTaskAttempt()
```

**Возможности:**
- Генерирует пользователей с хешированными паролями (bcrypt)
- Создает предметы, темы, задачи (реалистичные названия на русском)
- Генерирует учебные сессии с попытками (50% задач из набора пользователя)
- Автоматически создает Repetitions по алгоритму scheduler'а
- Поддерживает кастомные параметры через Config

### 2. Создан CLI `cmd/seed/main.go`
```bash
go run ./cmd/seed [-users=N] [-subjects=N] [-topics=N] [-tasks=N] [-sessions=N]
```

**Параметры:**
- `-users=5` (по умолчанию)
- `-subjects=2` (на пользователя)
- `-topics=3` (на предмет)
- `-tasks=5` (на тему)
- `-sessions=5` (на пользователя)

### 3. Написана документация
- `docs/seeding.md` — Полное руководство по seed скриптам
- `README.md` — Быстрый старт и полная архитектура проекта

---

## 🚀 Как использовать

### Быстрый старт (стандартные данные)
```bash
# 1. Запустить Docker
docker compose up -d

# 2. Подождать 5-10 сек, пока PostgreSQL готов
# 3. Запустить seed
go run ./cmd/seed

# 4. Запустить API
go run ./cmd/api

# 5. Проверить здоровье
curl http://localhost:8080/health
```

### Кастомные параметры
```bash
# Большой набор для нагрузочного теста
go run ./cmd/seed -users=50 -subjects=5 -topics=10 -tasks=20 -sessions=20

# Создаст: 50 юзеров × 5 субъектов × 10 тем × 20 задач = 50,000 задач!
```

---

## 📊 Что генерируется

**По умолчанию (с флагами по умолчанию):**
- 5 пользователей
- 10 предметов (subjects)
- 30 тем (topics)
- 150 задач (tasks)
- 25 сессий (5 пользователей × 5 сессий)
- ~150 task_attempts
- ~150 repetitions

**Итого:** ~500+ записей в БД

---

## 🔍 Пример сгенерированных данных

### Users:
```
user1@example.com (password1)
user2@example.com (password2)
...
user5@example.com (password5)
```

### Subjects:
```
Математика (1)
Физика (1)
Английский язык (1)
История (1)
...
```

### Tasks:
```
Задача 1: Дайте определение [Theory] (5-65 сек)
Задача 2: Решите уравнение [Practice] (5-65 сек)
Задача 3: Переведите на английский [Flashcard] (5-65 сек)
...
```

### Results Distribution:
- 60% Correct → Повторить через 3 дня
- 20% Wrong → Повторить через 1 день
- 15% Partial → Повторить через 2 дня
- 5% Skipped → Повторить через 1 день

---

## ✅ Проверка работы

### 1. После seed запустить API
```bash
go run ./cmd/api
```

Должно вывести:
```
INFO starting API env=local
INFO database connected host=localhost db=diplom
INFO database schema migrated
INFO api routes registered prefix=/api/v1
INFO http listening addr=:8080
```

### 2. Проверить в БД
```bash
docker compose exec postgres psql -U diplom -d diplom

# В REPL:
SELECT COUNT(*) FROM users;              # 5
SELECT COUNT(*) FROM subjects;           # 10
SELECT COUNT(*) FROM topics;             # 30
SELECT COUNT(*) FROM tasks;              # 150
SELECT COUNT(*) FROM task_attempts;      # ~150
SELECT COUNT(*) FROM repetitions;        # ~150
```

### 3. Проверить через API
```bash
# Регистрация
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user1@example.com", "password": "password1", "name": "User 1"}'

# Вход
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user1@example.com", "password": "password1"}' | jq -r '.data.token')

# Получить профиль
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/me

# Получить статистику
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/me/stats
```

---

## 🎬 Использование для ВКР

**Рекомендуемая конфигурация:**
```bash
go run ./cmd/seed -users=3 -subjects=2 -topics=3 -tasks=5 -sessions=3
```

**Почему эти значения:**
- 3 пользователя → можно показать разные истории
- 2 предмета на пользователя → 6 субъектов всего
- 3 темы на предмет → 18 тем
- 5 задач на тему → 90 задач
- 3 сессии на пользователя → 9 сессий

**Результат:** 
- Быстро запускается (2-3 сек)
- Хватает для демонстрации всех feature'ов
- Достаточно данных для красивой статистики

---

## 📝 Дальнейшие этапы

После Stage 7a рекомендуемый порядок:

1. **Stage 8d** (Demo script & curl) — 1.5 часа
   - Bash-скрипт с примерами запросов
   - Шоу всех основных flow'ов

2. **Stage 8a** (OpenAPI/Swagger) — 2 часа
   - Интерактивная документация
   - `/docs/swagger` endpoint

3. **Stage 8b** (Интеграционные тесты) — 3 часа
   - Tests для основных flows
   - `go test ./internal/...`

4. **Stage 8c** (Performance) — 2 часа
   - EXPLAIN ANALYZE анализ
   - Диаграммы производительности

5. **Stage 9** (Документация для ВКР) — 4 часа
   - Архитектурная диаграмма
   - Примеры API
   - Описание алгоритмов

---

## 💡 Советы

### Если нужны разные типы данных:
```bash
# Для быстрого теста (10 сек)
go run ./cmd/seed -users=1 -subjects=1 -topics=1 -tasks=5 -sessions=1

# Для демонстрации (30 сек)
go run ./cmd/seed -users=3 -subjects=2 -topics=3 -tasks=5 -sessions=3

# Для нагрузочного теста (30 сек)
go run ./cmd/seed -users=20 -subjects=3 -topics=5 -tasks=10 -sessions=10

# Для ВКР (1-2 минуты)
go run ./cmd/seed -users=10 -subjects=2 -topics=5 -tasks=10 -sessions=5
```

### Если seed падает:
1. Убедись Docker запущен: `docker compose ps`
2. Проверь БД готова: `docker compose logs postgres`
3. Перезапусти: `docker compose down -v && docker compose up -d`

### Если seed очень медленный:
- Уменьши количество параметров
- Проверь дисковый лимит: `docker system df`
- Очисти старые контейнеры: `docker system prune`

---

## 📦 Файлы добавленные на Stage 7a

```
✅ internal/seeder/seeder.go          (~300 строк)
✅ cmd/seed/main.go                   (~100 строк)
✅ docs/seeding.md                    (~150 строк)
✅ README.md                          (~350 строк)
```

**Итого:** ~900 строк документации и кода.

---

## 🎉 Готово!

Stage 7a завершен ✅

**Статус проекта:**
- ✅ Stages 1-6: Основная архитектура
- ✅ Stage 7a: Seed скрипты
- ⏳ Stages 8a-8d: Документация и тесты
- ⏳ Stage 9: Документация для ВКР

**Дальше:** Stage 8d — Demo script для демонстрации на защите! 🚀
