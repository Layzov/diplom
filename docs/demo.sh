#!/bin/bash

###############################################################################
# Exam Prep Backend — Complete Demo Flow
#
# Usage:
#   bash docs/demo.sh
#   bash docs/demo.sh -delay 2        # Add 2-second pause between requests
#
# Prerequisites:
#   1. PostgreSQL running: docker compose up -d
#   2. Seed data: go run ./cmd/seed
#   3. API running: go run ./cmd/api
#
# This script demonstrates the complete learning cycle:
# Register → Create Subject → Create Topic → Create Task → Start Session →
# Add Attempts → Complete Session → View Statistics → Check Repetitions
###############################################################################

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
API_PREFIX="$BASE_URL/api/v1"
DELAY="${DELAY:-0}"

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -delay)
      DELAY="$2"
      shift 2
      ;;
    *)
      shift
      ;;
  esac
done

# Helper functions
print_header() {
  echo -e "\n${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${CYAN}$1${NC}"
  echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_step() {
  echo -e "\n${BLUE}→ $1${NC}"
}

print_success() {
  echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
  echo -e "${RED}❌ $1${NC}"
}

print_json() {
  echo "$1" | jq '.' 2>/dev/null || echo "$1"
}

pause_if_needed() {
  if [ "$DELAY" -gt 0 ]; then
    sleep "$DELAY"
  fi
}

# API call helper
call_api() {
  local method=$1
  local endpoint=$2
  local data=$3
  local token=$4
  
  local headers="-H 'Content-Type: application/json'"
  if [ -n "$token" ]; then
    headers="$headers -H 'Authorization: Bearer $token'"
  fi
  
  if [ -z "$data" ] || [ "$data" == "null" ]; then
    curl -s -X "$method" "$API_PREFIX$endpoint" $headers
  else
    curl -s -X "$method" "$API_PREFIX$endpoint" $headers -d "$data"
  fi
}

###############################################################################
# DEMO START
###############################################################################

print_header "🎓 Exam Prep Backend — Complete Demo Flow"

echo -e "${YELLOW}API URL:${NC} $API_PREFIX"
echo -e "${YELLOW}Delay:${NC} ${DELAY}s between requests"

# 1. Health Check
print_header "1️⃣  HEALTH CHECK"
print_step "GET /health"
HEALTH=$(call_api GET "/health" "")
print_json "$HEALTH"
pause_if_needed

if ! echo "$HEALTH" | jq -e '.status' > /dev/null 2>&1; then
  print_error "API is not responding. Make sure: docker compose up -d && go run ./cmd/api"
  exit 1
fi
print_success "API is healthy!"

# 2. Register User
print_header "2️⃣  REGISTER NEW USER"
EMAIL="demo_$(date +%s)@example.com"
print_step "POST /auth/register"
echo -e "${YELLOW}Email:${NC} $EMAIL"
REGISTER=$(call_api POST "/auth/register" "{
  \"email\": \"$EMAIL\",
  \"password\": \"demo123456\",
  \"name\": \"Demo User\"
}")
print_json "$REGISTER"
pause_if_needed

USER_ID=$(echo "$REGISTER" | jq -r '.data.id // empty')
if [ -z "$USER_ID" ]; then
  print_error "Registration failed!"
  exit 1
fi
print_success "User registered: $USER_ID"

# 3. Login
print_header "3️⃣  LOGIN & GET TOKEN"
print_step "POST /auth/login"
LOGIN=$(call_api POST "/auth/login" "{
  \"email\": \"$EMAIL\",
  \"password\": \"demo123456\"
}")
print_json "$LOGIN"
pause_if_needed

TOKEN=$(echo "$LOGIN" | jq -r '.data.token // empty')
if [ -z "$TOKEN" ]; then
  print_error "Login failed!"
  exit 1
fi
print_success "Token received!"

# 4. Create Subject
print_header "4️⃣  CREATE SUBJECT (Предмет)"
print_step "POST /subjects"
SUBJECT=$(call_api POST "/subjects" "{
  \"title\": \"Демонстрационная Математика\",
  \"description\": \"Полный цикл демонстрации системы подготовки к экзаменам\"
}" "$TOKEN")
print_json "$SUBJECT"
pause_if_needed

SUBJECT_ID=$(echo "$SUBJECT" | jq -r '.data.id // empty')
if [ -z "$SUBJECT_ID" ]; then
  print_error "Failed to create subject!"
  exit 1
fi
print_success "Subject created: $SUBJECT_ID"

# 5. Create Topic
print_header "5️⃣  CREATE TOPIC (Тема)"
print_step "POST /topics"
TOPIC=$(call_api POST "/topics" "{
  \"subject_id\": \"$SUBJECT_ID\",
  \"title\": \"Основные понятия и теория\"
}" "$TOKEN")
print_json "$TOPIC"
pause_if_needed

TOPIC_ID=$(echo "$TOPIC" | jq -r '.data.id // empty')
if [ -z "$TOPIC_ID" ]; then
  print_error "Failed to create topic!"
  exit 1
fi
print_success "Topic created: $TOPIC_ID"

# 6. Create Task
print_header "6️⃣  CREATE TASK (Задача)"
print_step "POST /tasks"
TASK=$(call_api POST "/tasks" "{
  \"topic_id\": \"$TOPIC_ID\",
  \"type\": \"practice\",
  \"title\": \"Решите задачу по геометрии\",
  \"content\": \"Найдите площадь треугольника со сторонами a=5, b=6, c=7\"
}" "$TOKEN")
print_json "$TASK"
pause_if_needed

TASK_ID=$(echo "$TASK" | jq -r '.data.id // empty')
if [ -z "$TASK_ID" ]; then
  print_error "Failed to create task!"
  exit 1
fi
print_success "Task created: $TASK_ID"

# 7. Create Session
print_header "7️⃣  START STUDY SESSION"
print_step "POST /sessions"
NOW=$(date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u +"%Y-%m-%dT%H:%M:%SZ")
SESSION=$(call_api POST "/sessions" "{
  \"subject_id\": \"$SUBJECT_ID\",
  \"started_at\": \"$NOW\",
  \"notes\": \"Demo study session\"
}" "$TOKEN")
print_json "$SESSION"
pause_if_needed

SESSION_ID=$(echo "$SESSION" | jq -r '.data.id // empty')
if [ -z "$SESSION_ID" ]; then
  print_error "Failed to create session!"
  exit 1
fi
print_success "Session started: $SESSION_ID"

# 8. Add Task Attempts
print_header "8️⃣  ADD TASK ATTEMPTS"

# Attempt 1: Correct
print_step "POST /task_attempts (Correct)"
ATTEMPT_1=$(call_api POST "/task_attempts" "{
  \"session_id\": \"$SESSION_ID\",
  \"task_id\": \"$TASK_ID\",
  \"result\": \"correct\",
  \"time_ms\": 15000
}" "$TOKEN")
print_json "$ATTEMPT_1"
pause_if_needed

ATTEMPT_1_ID=$(echo "$ATTEMPT_1" | jq -r '.data.id // empty')
print_success "Correct attempt recorded: $ATTEMPT_1_ID"

# Attempt 2: Partial
print_step "POST /task_attempts (Partial)"
ATTEMPT_2=$(call_api POST "/task_attempts" "{
  \"session_id\": \"$SESSION_ID\",
  \"task_id\": \"$TASK_ID\",
  \"result\": \"partial\",
  \"time_ms\": 22000
}" "$TOKEN")
print_json "$ATTEMPT_2"
pause_if_needed

ATTEMPT_2_ID=$(echo "$ATTEMPT_2" | jq -r '.data.id // empty')
print_success "Partial attempt recorded: $ATTEMPT_2_ID"

# 9. Complete Session
print_header "9️⃣  COMPLETE STUDY SESSION"
print_step "POST /sessions/{id}/complete"
END_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u +"%Y-%m-%dT%H:%M:%SZ")
COMPLETE=$(call_api POST "/sessions/$SESSION_ID/complete" "{
  \"ended_at\": \"$END_TIME\",
  \"notes\": \"Session completed successfully\"
}" "$TOKEN")
print_json "$COMPLETE"
pause_if_needed

print_success "Session completed!"

# 10. Get User Statistics
print_header "🔟 GET USER STATISTICS"
print_step "GET /me/stats"
STATS=$(call_api GET "/me/stats" "" "$TOKEN")
print_json "$STATS"
pause_if_needed

print_success "Statistics retrieved!"

# 11. Get Topic Statistics
print_header "1️⃣ 1️⃣  GET TOPIC STATISTICS"
print_step "GET /me/stats/topics"
TOPIC_STATS=$(call_api GET "/me/stats/topics" "" "$TOKEN")
print_json "$TOPIC_STATS"
pause_if_needed

print_success "Topic statistics retrieved!"

# 12. Get Upcoming Repetitions
print_header "1️⃣ 2️⃣  GET UPCOMING REPETITIONS"
print_step "GET /me/repetitions/upcoming"
REPETITIONS=$(call_api GET "/me/repetitions/upcoming" "" "$TOKEN")
print_json "$REPETITIONS"
pause_if_needed

REPETITION_COUNT=$(echo "$REPETITIONS" | jq '.data | length')
print_success "Found $REPETITION_COUNT upcoming repetitions!"

# 13. Get User Profile
print_header "1️⃣ 3️⃣  GET USER PROFILE"
print_step "GET /me"
PROFILE=$(call_api GET "/me" "" "$TOKEN")
print_json "$PROFILE"
pause_if_needed

print_success "Profile retrieved!"

# 14. Get All Subjects
print_header "1️⃣ 4️⃣  GET ALL SUBJECTS"
print_step "GET /subjects"
SUBJECTS=$(call_api GET "/subjects" "" "$TOKEN")
print_json "$SUBJECTS"
pause_if_needed

SUBJECT_COUNT=$(echo "$SUBJECTS" | jq '.data | length')
print_success "Found $SUBJECT_COUNT subjects!"

# 15. Get All Topics for Subject
print_header "1️⃣ 5️⃣  GET ALL TOPICS FOR SUBJECT"
print_step "GET /topics?subject_id=$SUBJECT_ID"
TOPICS=$(call_api GET "/topics?subject_id=$SUBJECT_ID" "" "$TOKEN")
print_json "$TOPICS"
pause_if_needed

print_success "Topics retrieved!"

# 16. Get All Tasks for Topic
print_header "1️⃣ 6️⃣  GET ALL TASKS FOR TOPIC"
print_step "GET /tasks?topic_id=$TOPIC_ID"
TASKS=$(call_api GET "/tasks?topic_id=$TOPIC_ID" "" "$TOKEN")
print_json "$TASKS"
pause_if_needed

print_success "Tasks retrieved!"

# 17. Get All Sessions
print_header "1️⃣ 7️⃣  GET ALL SESSIONS"
print_step "GET /sessions"
SESSIONS=$(call_api GET "/sessions" "" "$TOKEN")
print_json "$SESSIONS"
pause_if_needed

SESSION_COUNT=$(echo "$SESSIONS" | jq '.data | length')
print_success "Found $SESSION_COUNT sessions!"

# Final Summary
print_header "🎉 DEMO COMPLETE!"

echo -e "\n${GREEN}Summary of Operations:${NC}"
echo "  ✅ Health Check"
echo "  ✅ User Registration"
echo "  ✅ User Authentication (JWT)"
echo "  ✅ Subject Creation"
echo "  ✅ Topic Creation"
echo "  ✅ Task Creation"
echo "  ✅ Study Session Creation"
echo "  ✅ Task Attempts Recording"
echo "  ✅ Session Completion"
echo "  ✅ Statistics Generation"
echo "  ✅ Repetitions Calculation"
echo "  ✅ Data Retrieval (All Resources)"

echo -e "\n${CYAN}Demo User Details:${NC}"
echo "  Email: $EMAIL"
echo "  User ID: $USER_ID"
echo "  Subject: $SUBJECT_ID"
echo "  Topic: $TOPIC_ID"
echo "  Task: $TASK_ID"
echo "  Session: $SESSION_ID"

echo -e "\n${YELLOW}Key Features Demonstrated:${NC}"
echo "  🔐 JWT Authentication with Bearer tokens"
echo "  📚 Full subject → topic → task hierarchy"
echo "  🎯 Study session with multiple task attempts"
echo "  📊 Automatic statistics calculation"
echo "  🔄 Spaced repetition scheduling"
echo "  📈 User learning progress tracking"

echo -e "\n${GREEN}Next Steps:${NC}"
echo "  1. Review statistics in the API responses"
echo "  2. Check database: SELECT COUNT(*) FROM repetitions;"
echo "  3. Modify script to test error cases"
echo "  4. Integrate with frontend"

echo ""
print_header "✨ End of Demo ✨"
