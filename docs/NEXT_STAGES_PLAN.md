# 📅 План работ — Этапы 8-9

## Обзор оставшихся этапов

**Статус:** 7 из 13 этапов завершено (54%)

```
Завершено ✅
├─ Stage 1: Скелет проекта
├─ Stage 2: Модели и миграции
├─ Stage 3: CRUD основных сущностей
├─ Stage 4: Учебный цикл
├─ Stage 5: Статистика и views
├─ Stage 5.5: Аутентификация
├─ Stage 6: Индексы и оптимизация
└─ Stage 7a: Seed скрипты ✨ НОВОЕ

В работе ⏳
├─ Stage 8a: OpenAPI/Swagger документация (2 часа)
├─ Stage 8b: Интеграционные тесты (3 часа)
├─ Stage 8c: Performance analysis (2 часа)
├─ Stage 8d: Demo script & curl примеры (1.5 часа)
└─ Stage 9: Документация для ВКР (4 часа)
```

---

## 🎯 Stage 8d: Demo Script (ПРИОРИТЕТ)

**Зачем:** На защите нужно что-то показывать. Curl-скрипты это делают за секунды.

**Время:** ~1.5 часа

### Что нужно

Создать `docs/demo.sh` — bash-скрипт со следующим flow:

1. **Health Check** (1 сек)
   ```bash
   curl http://localhost:8080/health
   ```

2. **Регистрация пользователя** (1 сек)
   ```bash
   curl -X POST .../auth/register \
     -d '{"email":"demo@example.com", "password":"demo123", "name":"Demo User"}'
   ```

3. **Вход и получение токена** (1 сек)
   ```bash
   TOKEN=$(curl -X POST .../auth/login ...)
   ```

4. **Создание предмета** (1 сек)
   ```bash
   SUBJECT=$(curl -X POST .../subjects \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"title":"Математика",...}')
   ```

5. **Создание темы** (1 сек)
   ```bash
   TOPIC=$(curl -X POST .../topics \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"subject_id":"...", "title":"Тема 1",...}')
   ```

6. **Создание задачи** (1 сек)
   ```bash
   TASK=$(curl -X POST .../tasks \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"topic_id":"...", "type":"practice", ...}')
   ```

7. **Создание сессии** (1 сек)
   ```bash
   SESSION=$(curl -X POST .../sessions \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"started_at":"...", "subject_id":"..."}')
   ```

8. **Добавление попыток** (2 сек)
   ```bash
   curl -X POST .../task_attempts \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"task_id":"...", "result":"correct", "time_ms":5000}'
   ```

9. **Завершение сессии** (1 сек)
   ```bash
   curl -X POST .../sessions/.../complete \
     -H "Authorization: Bearer $TOKEN"
   ```

10. **Получение статистики** (1 сек)
    ```bash
    curl .../me/stats \
      -H "Authorization: Bearer $TOKEN"
    ```

11. **Получение ближайших повторений** (1 сек)
    ```bash
    curl .../me/repetitions/upcoming \
      -H "Authorization: Bearer $TOKEN"
    ```

**Итого:** ~15-20 секунд для полного цикла

### Структура скрипта

```bash
#!/bin/bash
set -e

BASE_URL="http://localhost:8080"
API_PREFIX="$BASE_URL/api/v1"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Helper function
call_api() {
  local method=$1
  local endpoint=$2
  local data=$3
  local token=$4
  
  if [ -z "$token" ]; then
    curl -s -X $method "$API_PREFIX$endpoint" \
      -H "Content-Type: application/json" \
      -d "$data"
  else
    curl -s -X $method "$API_PREFIX$endpoint" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $token" \
      -d "$data"
  fi
}

# 1. Health check
echo -e "${BLUE}1. Health check${NC}"
call_api GET "/health" ""
echo ""

# 2. Register
echo -e "${BLUE}2. Register new user${NC}"
REGISTER=$(call_api POST "/auth/register" '{"email":"demo@example.com","password":"demo123","name":"Demo User"}')
echo $REGISTER | jq '.'
echo ""

# 3. Login
echo -e "${BLUE}3. Login${NC}"
LOGIN=$(call_api POST "/auth/login" '{"email":"demo@example.com","password":"demo123"}')
TOKEN=$(echo $LOGIN | jq -r '.data.token')
echo "Token: $TOKEN"
echo ""

# ... и т.д. для каждого шага
```

---

## 🔐 Stage 8a: OpenAPI/Swagger (2 часа)

**Зачем:** Интерактивная документация, стандарт в индустрии.

### Опции реализации

**Вариант 1 (Рекомендуется):** swaggo/swag
```bash
# Установить
go install github.com/swaggo/swag/cmd/swag@latest

# Аннотировать handlers
```

**Вариант 2 (Простой):** Ручной OpenAPI YAML
```
docs/openapi.yaml  (~500 строк вручную)
```

**Вариант 3 (Hybrid):** Предоставить JSON с примерами
```
docs/api-examples.json
```

### Что генерируется

```
GET /docs/swagger
├─ Swagger UI
└─ Interactive endpoints для тестирования
```

---

## 🧪 Stage 8b: Integration Tests (3 часа)

**Зачем:** Доказать, что система работает.

### Тесты для

1. **Auth Flow**
   - Register + Login
   - Token validation
   - Ownership checks

2. **CRUD Flow**
   - Create Subject, Topic, Task
   - Read операции
   - Update операции
   - Delete с проверкой каскада

3. **Session Flow**
   - Create session
   - Add task attempts
   - Complete session
   - Calculate stats

4. **Stats Flow**
   - Fetch user stats
   - Fetch topic stats
   - Fetch upcoming repetitions

### Структура

```
internal/handler/...handler_test.go
internal/service/...service_test.go
internal/repository/...repository_test.go
```

### Запуск

```bash
go test ./internal/... -v

# С покрытием
go test ./internal/... -cover

# Конкретный тест
go test ./internal/handler -run TestAuthRegister -v
```

---

## 📊 Stage 8c: Performance Analysis (2 часа)

**Зачем:** Показать оптимизацию, проверить масштабируемость.

### Анализ

1. **EXPLAIN ANALYZE** для сложных запросов
   ```sql
   EXPLAIN ANALYZE
   SELECT * FROM user_learning_stats WHERE user_id = '...';
   ```

2. **Timing tests**
   - Время загрузки статистики
   - Скорость фильтрации repetitions
   - Производительность при 10K записей

3. **Документирование результатов**
   - Скриншоты EXPLAIN ANALYZE
   - Таблицы с метриками
   - Выводы

### Файл

```
docs/performance.md
├─ Query Analysis
├─ Metrics
├─ Optimization Recommendations
└─ Scaling Strategy
```

---

## 📝 Stage 9: Документация для ВКР (4 часа)

**Зачем:** Полная архитектурная документация для защиты.

### Содержание

1. **Архитектурная диаграмма** (Mermaid или Draw.io)
   ```
   User
    ↓ HTTP
   Handler Layer
    ↓ Business Logic
   Service Layer
    ↓ Database Queries
   Repository Layer
    ↓ SQL
   PostgreSQL
   ```

2. **Диаграмма баз данных** (ERD)
   ```
   Users
   ├─ Subjects
   │  ├─ Topics
   │  │  └─ Tasks
   │  │     └─ TaskAttempts
   │  │        ├─ AttemptStats
   │  │        └─ Repetitions
   └─ Sessions
      └─ TaskAttempts
   ```

3. **API примеры** (для каждого endpoint)
   ```
   POST /api/v1/auth/register
   Request: {...}
   Response: {...}
   ```

4. **Алгоритмы**
   - Scheduler для спaced repetition
   - Расчет статистики
   - Ownership checks

5. **Безопасность**
   - JWT токены
   - Password hashing (bcrypt)
   - SQL injection prevention (GORM)

### Файлы

```
docs/
├─ ARCHITECTURE.md         (Архитектура)
├─ DATABASE_SCHEMA.md      (Schema ERD)
├─ ALGORITHMS.md           (Алгоритмы)
├─ API_EXAMPLES.md         (Примеры)
└─ SECURITY.md             (Безопасность)
```

---

## 🎯 Рекомендуемый порядок

**Общее время:** ~12.5 часов

```
1. Stage 8d  (1.5 часа)  ← START HERE (нужно для защиты)
2. Stage 8a  (2 часа)    ← Высокий приоритет
3. Stage 8b  (3 часа)    ← Желательно
4. Stage 8c  (2 часа)    ← Опционально
5. Stage 9   (4 часа)    ← ОБЯЗАТЕЛЬНО для ВКР
```

**Быстрый вариант (8 часов):**
```
1. Stage 8d  (1.5 часа) ← Demo для защиты
2. Stage 8a  (2 часа)   ← Swagger UI
3. Stage 9   (4.5 часа) ← Документация ВКР

Пропускаем: 8b (tests), 8c (performance)
```

---

## 📌 Для защиты ВКР нужны

✅ **Обязательно:**
1. Seed-скрипты (Stage 7a) — ГОТОВО
2. Demo script (Stage 8d) — 1.5 часа
3. README (уже создан)
4. Документация (Stage 9) — 4 часа

⚠️ **Желательно:**
- OpenAPI/Swagger (Stage 8a)
- Тесты (Stage 8b)

---

## 💡 Советы

1. **Stage 8d первым** — это даст максимум на защите (живая демонстрация)
2. **Stage 9 обязателен** — без документации не будет оценки за ВКР
3. **Stages 8b и 8c опционально** — если будет время

---

## 🚀 Начинаем?

Готов ли ты начать Stage 8d (Demo Script)?

Это займет ~1.5 часа, но даст максимум на защите! 🎬
