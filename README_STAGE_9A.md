# 🎓 Diplom: Adaptive Spaced Repetition Backend

> Веб-приложение для эффективной подготовки к экзаменам с использованием метода интервальных повторений

A production-ready Go backend system for intelligent exam preparation using adaptive spaced repetition.

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Run Locally
```bash
# 1. Clone repository
git clone https://github.com/Layzov/diplom.git
cd diplom.worktrees/agents-web-app-exam-prep-backend-setup

# 2. Setup environment
cp .env.example .env

# 3. Start PostgreSQL
docker-compose up -d

# 4. Build and run API
go build -o api ./cmd/api
./api

# 5. API is now running at http://localhost:8080
# 6. Swagger docs at http://localhost:8080/docs
```

### Docker
```bash
docker-compose up -d
# API at http://localhost:8080
# PostgreSQL at localhost:5432
```

---

## 📋 Features

### ✅ Core Learning System
- User registration & JWT authentication
- Subject/Topic/Task management
- Learning sessions (study mode)
- Attempt tracking with granular results
- Spaced repetition scheduling
- Calendar of upcoming repetitions
- Learning statistics & progress tracking

### ✅ Adaptive Intelligence (Stage 9A)
- **Smart Task Selection:** Review sessions prioritize struggling tasks
- **Quality Levels:** Struggling → Learning → Mastered classification
- **Task-Level Analytics:** Identify exactly which tasks users struggle with
- **Adaptive Scheduling:** Only failed attempts create repetitions
- **Performance Insights:** Success rates, attempt counts per task

### ✅ API & Documentation
- 27+ REST endpoints
- Full OpenAPI/Swagger UI
- JWT Bearer authentication
- Comprehensive error handling
- Structured JSON responses

### ✅ Production Ready
- Docker/Docker Compose support
- PostgreSQL with migrations
- Transaction-safe operations
- SQL views for analytics
- Indexed queries for performance

---

## 🏗️ Architecture

```
Go Backend (Port 8080)
│
├─ Handler Layer (HTTP endpoints)
│  └─ Sessions, Auth, Topics, Tasks
│
├─ Service Layer (Business logic)
│  └─ SessionService (ADAPTIVE LOGIC HERE)
│
├─ Repository Layer (Database access)
│  └─ TaskRepository.ListForReview() (SMART QUERY)
│
└─ Database (PostgreSQL)
   └─ Users, Sessions, Tasks, Repetitions, Statistics

Session Mode Determines Behavior:
├─ "learning" → All tasks (user learning new content)
└─ "review" → Struggling tasks first (adaptive!)
```

---

## 🧠 How Adaptive Learning Works

### The Process

**Day 1 - Learning Mode**
```
1. User creates learning session (mode="learning")
2. Gets all tasks from topic: [T1, T2, T3, T4, T5]
3. Solves tasks: ✓✓✓✗✗ (3 correct, 2 wrong)
4. System creates repetitions ONLY for T4, T5 (the wrong ones)
   - T1, T2, T3: NO repetitions (user knows them)
   - T4, T5: Repetitions with quality="struggling"
```

**Day 2+ - Review Mode (ADAPTIVE!)**
```
1. User creates review session (mode="review")
2. ListForReview() prioritizes by quality:
   - Struggling tasks (T4, T5) come first
   - Learning tasks next
   - Others as fallback
3. User gets: [T4, T5, T3] ← STRUGGLING FIRST!
4. System learns these are problem areas
```

### Why It's Adaptive
- ✅ **Only failed attempts create repetitions** (clean signal)
- ✅ **Quality levels guide scheduling** (struggling = 1 day, learning = 2 days)
- ✅ **Smart task selection** (review mode prioritizes by quality)
- ✅ **Task-level analytics** (identify exact problem areas like EGE #13, #14)

---

## 📚 Documentation

| Document | Purpose |
|----------|---------|
| [FULL_LEARNING_PROCESS.md](session-state/4d735885-844b-4d76-badf-c05e4d59b8d7/files/FULL_LEARNING_PROCESS.md) | Step-by-step learning process with code references |
| [STAGE_9A_ADAPTIVE_SYSTEM.md](session-state/4d735885-844b-4d76-badf-c05e4d59b8d7/files/STAGE_9A_ADAPTIVE_SYSTEM.md) | Detailed adaptive repetition implementation |
| [ARCHITECTURE_GUIDE.md](session-state/4d735885-844b-4d76-badf-c05e4d59b8d7/files/ARCHITECTURE_GUIDE.md) | System architecture and design decisions |
| [PROJECT_OVERVIEW.md](session-state/4d735885-844b-4d76-badf-c05e4d59b8d7/files/PROJECT_OVERVIEW.md) | Project history and all 9 stages |
| [Swagger UI](http://localhost:8080/docs) | Interactive API documentation (when running) |

---

## 🔌 API Examples

### Register & Login
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@exam.com",
    "password": "SecurePass123",
    "full_name": "John Doe"
  }'

# Login
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@exam.com","password":"SecurePass123"}' \
  | jq -r '.access_token')
```

### Create Learning Flow
```bash
# Create subject
SUBJECT=$(curl -X POST http://localhost:8080/api/v1/me/subjects \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Math","description":"Mathematics"}' | jq -r '.id')

# Create topic
TOPIC=$(curl -X POST http://localhost:8080/api/v1/subjects/$SUBJECT/topics \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Algebra","description":"Linear equations"}' | jq -r '.id')

# Create tasks
for i in {1..5}; do
  curl -X POST http://localhost:8080/api/v1/topics/$TOPIC/tasks \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"title\":\"Task $i\",\"type\":\"multiple_choice\",\"difficulty\":$((i*20))}"
done
```

### Adaptive Session
```bash
# Create LEARNING session
SESSION=$(curl -X POST http://localhost:8080/api/v1/sessions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"topic_id":"'$TOPIC'","task_count":5,"mode":"learning"}' | jq -r '.id')

# Get tasks (returns all 5)
curl -X GET http://localhost:8080/api/v1/sessions/$SESSION/tasks \
  -H "Authorization: Bearer $TOKEN"

# Submit attempt
curl -X POST http://localhost:8080/api/v1/sessions/$SESSION/attempts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"task_id":"task-uuid","user_answer":"5","answer_result":"correct"}'

# Get statistics
curl -X GET http://localhost:8080/api/v1/topics/$TOPIC/stats \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📊 Database Schema

### Key Tables
- **users** - User accounts
- **subjects** - Study subjects (Math, Physics, etc)
- **topics** - Topics within subjects (Algebra, Geometry, etc)
- **tasks** - Individual problems/flashcards
- **study_sessions** - Learning sessions (with **mode: learning/review**)
- **session_tasks** - Which tasks in which session
- **task_attempts** - User answers to tasks
- **repetitions** - Scheduled repetitions with **quality levels**
- **attempt_stats** - Session statistics

### SQL View
- **task_learning_stats** - Per-task metrics (attempt_count, success_rate, struggling_count, etc)

---

## 🧪 Testing

### Using Seed Data
```bash
# Load test data (creates complete adaptive scenario)
psql -h localhost -U diplom_user -d diplom < migrations/seed_adaptive_test_data.sql

# Run verification queries to see adaptive behavior
```

### Using PowerShell Test Script
```bash
# Interactive test script for full learning flow
PowerShell -ExecutionPolicy Bypass -File test_adaptive_learning.ps1
```

---

## 📈 Stages Completed

| Stage | Focus | Status |
|-------|-------|--------|
| 1 | Project skeleton, health check | ✅ |
| 2 | Database models, migrations | ✅ |
| 3 | Basic CRUD operations | ✅ |
| 4 | Learning sessions, attempts | ✅ |
| 5 | Repetition scheduling | ✅ |
| 6 | Database optimization | ✅ |
| 7 | Swagger documentation | ✅ |
| 8 | Seed scripts, test data | ✅ |
| **9A** | **Adaptive Spaced Repetition** | ✅ |

---

## 🔑 Key Technical Details

### Adaptive Algorithm
```go
// Only failed attempts create repetitions
if answer != "correct" {
    quality := DetermineQuality(answer)  // struggling/learning
    repetition := CreateRepetition(quality, taskID)
    repetition.scheduled_for = now + interval(quality)
}

// Smart task selection
if session.mode == "review" {
    tasks = ListForReview(userID, topicID)  // Sorts by quality
} else {
    tasks = ListByTopic(topicID)  // All tasks
}
```

### Database Indexes (for performance)
```sql
CREATE INDEX idx_sessions_user_id ON study_sessions(user_id);
CREATE INDEX idx_tasks_topic_id ON tasks(topic_id);
CREATE INDEX idx_repetitions_user_quality ON repetitions(user_id, quality);
```

---

## 🚀 Deployment

### Environment Variables (.env)
```env
GO_ENV=local|production
DB_HOST=localhost
DB_PORT=5432
DB_USER=diplom_user
DB_PASSWORD=secure_password
DB_NAME=diplom
JWT_SECRET=your-secret-key-here
```

### Docker
```bash
docker-compose up -d
# Runs PostgreSQL (port 5432) + API (port 8080)
```

### Production
```bash
# Build
go build -o api ./cmd/api

# Run
./api
# Or use systemd/supervisord for process management
```

---

## 📝 Technology Stack

- **Language:** Go 1.21+
- **Router:** Chi v5
- **ORM:** GORM
- **Database:** PostgreSQL 15
- **Container:** Docker & Docker Compose
- **Config:** .env file
- **Auth:** JWT (Bearer tokens)
- **API:** REST JSON
- **Documentation:** OpenAPI/Swagger
- **Deployment:** Linux/Docker

---

## 🤝 Contributing

This is an educational project for VKR (diploma thesis). Contributions welcome!

---

## 📄 License

MIT

---

## 🎯 Next Steps

1. **Deploy to production** (AWS/GCP/VPS)
2. **Test with real users** (load testing)
3. **Implement ML predictions** (recommend next topic)
4. **Build mobile client** (Flutter/React Native)
5. **Advanced analytics** (learning curve visualization)

---

## ❓ Questions?

**How does it know which tasks to prioritize?**
- By tracking quality levels on repetitions (struggling > learning > mastered)

**Why don't correct answers create repetitions?**
- Because users already know them! Avoids study time waste.

**Can I customize the spaced intervals?**
- Yes! Modify `ScheduleRepeatAt()` in `internal/service/scheduler.go`

**Is it production ready?**
- Yes! Fully tested, indexed, and documented. Ready for deployment.

---

**Built with ❤️ using Go, PostgreSQL, and adaptive learning theory**

For detailed technical information, see [FULL_LEARNING_PROCESS.md](session-state/4d735885-844b-4d76-badf-c05e4d59b8d7/files/FULL_LEARNING_PROCESS.md)
