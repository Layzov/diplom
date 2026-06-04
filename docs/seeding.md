# Stage 7a: Seed Scripts & Test Data

## Обзор

Seed скрипты генерируют реалистичные тестовые данные в БД для демонстрации и тестирования системы.

**Что генерируется:**
- 5 пользователей (по умолчанию)
- По 2 предмета на пользователя
- По 3 темы на предмет
- По 5 задач на тему
- По 5 учебных сессий на пользователя
- TaskAttempts (попытки с результатами: correct, partial, wrong, skipped)
- Repetitions (повторения, рассчитанные по алгоритму)

**Итого:** ~5000+ записей в БД за несколько секунд.

---

## 📦 Структура кода

```
internal/seeder/
  └─ seeder.go          # Config, Seed() функция, генераторы данных

cmd/seed/
  └─ main.go            # CLI точка входа с флагами
```

## 🚀 Использование

### Быстрый старт

```bash
# 1. Запустить Docker (если не запущен)
docker compose up -d

# 2. Подождать, пока PostgreSQL готов (5-10 сек)
sleep 10

# 3. Запустить seed со стандартными параметрами
go run ./cmd/seed

# 4. Запустить API и проверить данные
go run ./cmd/api
# Открыть http://localhost:8080/health
```

### Кастомизация

Генерируй больше данных для нагрузочного тестирования:

```bash
# 10 пользователей, по 3 предмета, по 4 темы, по 10 задач, по 10 сессий
go run ./cmd/seed \
  -users=10 \
  -subjects=3 \
  -topics=4 \
  -tasks=10 \
  -sessions=10
```

### Параметры

| Флаг | Значение по умолчанию | Описание |
|------|-------|---|
| `-users` | 5 | Количество пользователей |
| `-subjects` | 2 | Предметов на пользователя |
| `-topics` | 3 | Тем на предмет |
| `-tasks` | 5 | Задач на тему |
| `-sessions` | 5 | Сессий на пользователя |

---

## 📊 Что генерируется

### Users
```
user1@example.com - password1
user2@example.com - password2
...
user5@example.com - password5
```

Все пароли хешируются через bcrypt.

### Subjects (Предметы)
- Математика, Физика, Английский язык, История, Биология, Химия, Литература, Информатика

### Tasks (Задачи)
- **Type:** Theory, Practice, Flashcard
- **Result:** correct (60%), wrong (20%), partial (15%), skipped (5%)
- **Time:** 5-65 секунд на задачу

### Repetitions (Повторения)
Рассчитываются по алгоритму планировщика:
- **Correct:** Повторить через 3 дня
- **Partial:** Повторить через 2 дня
- **Wrong/Skipped:** Повторить через 1 день

---

## 💾 Пример выходных данных

```
INFO seed: starting users=5 subjects_per_user=2 topics_per_subject=3 tasks_per_topic=5 sessions_per_user=5
INFO database connected host=localhost db=diplom
INFO database schema migrated
INFO seed: completed successfully users=5 subjects=10 topics=30 tasks=150

✅ Database seeded successfully!
📊 Generated data:
  - 5 users
  - 2 subjects per user
  - 3 topics per subject
  - 5 tasks per topic
  - 5 sessions per user

💡 Next: run 'go run ./cmd/api' to start the server
```

---

## 🔍 Проверка данных после seed

### Запрос 1: Получить статистику по пользователю
```bash
#登录
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user1@example.com", "password": "password1"}' | jq -r '.data.token')

# Получить статистику
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/me/stats
```

### Запрос 2: Получить ближайшие повторения
```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/me/repetitions/upcoming
```

### Запрос 3: Получить статистику по темам
```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/me/stats/topics
```

---

## 📈 Масштабирование

Для нагрузочного тестирования:

```bash
# Большой набор данных для анализа производительности
go run ./cmd/seed \
  -users=50 \
  -subjects=5 \
  -topics=10 \
  -tasks=20 \
  -sessions=20

# Это создаст:
# - 50 пользователей
# - 250 предметов
# - 2500 тем
# - 50 000 задач
# - 25 000 сессий
# - ~200 000 попыток и повторений
```

---

## 🐛 Отладка

### Если seed падает на миграции:
```bash
# Убедись, что БД запущена
docker compose ps

# Проверь логи
docker compose logs postgres

# Перезапусти контейнеры
docker compose down -v
docker compose up -d
```

### Если seed медленный:
- Уменьши кол-во записей
- Проверь производительность БД
- Используй `docker compose logs postgres` для мониторинга

---

## 📝 Для защиты ВКР

**Рекомендуемый набор данных:**
```bash
go run ./cmd/seed -users=3 -subjects=2 -topics=3 -tasks=5 -sessions=3
```

Это создаст:
- 3 пользователя
- 6 предметов
- 18 тем
- 90 задач
- 9 сессий

**Достаточно для демонстрации, но быстро запускается (~2 сек).**

