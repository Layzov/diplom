# Этап 6 — индексы и проверка планов запросов

## Что добавлено

После `AutoMigrate` при старте API применяются дополнительные индексы (`internal/db/indexes.go`, дубликат в `migrations/0003_indexes.sql`).

| Индекс | Таблица | Назначение |
|--------|---------|------------|
| `idx_subjects_user_created_at` | subjects | `GET /me/subjects` |
| `idx_topics_subject_created_at` | topics | список тем предмета |
| `idx_tasks_topic_created_at` | tasks | список заданий темы |
| `idx_sessions_user_started_at` | sessions | сессии пользователя |
| `idx_task_attempts_session_created_at` | task_attempts | попытки в сессии |
| `idx_task_attempts_user_task` | task_attempts | join в `topic_learning_stats` |
| `idx_repetitions_planned_user_due` | repetitions (partial) | календарь / upcoming (`status = planned`) |
| `idx_repetitions_session_planned_due` | repetitions (partial) | `next_repeat_at` при завершении сессии |
| `idx_repetitions_user_status_due` | repetitions | список повторений с фильтром статуса |
| `idx_repetitions_user_topic` | repetitions | агрегация по темам |

Частичные индексы (`WHERE status = 'planned'`) уменьшают размер индекса и точнее покрывают view `calendar` и эндпоинты upcoming.

Составные индексы на FK + поле сортировки дублируются в GORM-тегах моделей (для новых установок через `AutoMigrate`) и в SQL (для уже существующих БД).

## Проверка планов

```powershell
cd f:\diplom
docker compose up -d
go run ./cmd/api          # применит индексы
go run ./cmd/explain      # EXPLAIN без выполнения запросов
go run ./cmd/explain -analyze   # после seed: реальные тайминги
go run ./cmd/explain -user=<uuid> -analyze
```

Ожидаемые признаки хорошего плана на больших данных:

- `Index Scan` / `Bitmap Index Scan` по индексам выше, а не `Seq Scan` на `repetitions`, `sessions`, `task_attempts`
- для upcoming/calendar — `idx_repetitions_planned_user_due`
- для списка сессий — `idx_sessions_user_started_at`

Представления `*_learning_stats` агрегируют много строк: на больших объёмах узкое место — сама агрегация, а не PK lookup. Для отчётов в проде позже можно добавить материализованные view или кэш; для учебного проекта достаточно индексов на базовых таблицах.

## Ручной EXPLAIN в psql

```sql
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM calendar
WHERE user_id = '...' AND due_at >= NOW()
ORDER BY due_at ASC LIMIT 20;
```

## Список индексов в БД

```sql
SELECT tablename, indexname, indexdef
FROM pg_indexes
WHERE schemaname = 'public'
  AND tablename IN (
    'subjects','topics','tasks','sessions',
    'task_attempts','repetitions','attempt_stats'
  )
ORDER BY tablename, indexname;
```

Утилита `go run ./cmd/schemacheck` после миграции выводит число индексов на ключевых таблицах.

## Этап 7

Для осмысленного `-analyze` нужны seed-данные (следующий этап). Без данных планы корректны, но времена не показательны.
