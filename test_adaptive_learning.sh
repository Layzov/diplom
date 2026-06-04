#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API="http://localhost:8080/api/v1"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Testing Adaptive Spaced Repetition System${NC}"
echo -e "${BLUE}========================================${NC}\n"

# ============================================================================
# STEP 1: Register and Login
# ============================================================================
echo -e "${YELLOW}Step 1: Registering user...${NC}"
REGISTER=$(curl -s -X POST "$API/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testadaptive@exam.com",
    "password": "SecurePass123!",
    "full_name": "Adaptive Learner"
  }')

echo "Response: $REGISTER"

echo -e "\n${YELLOW}Step 2: Logging in...${NC}"
LOGIN=$(curl -s -X POST "$API/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testadaptive@exam.com",
    "password": "SecurePass123!"
  }')

TOKEN=$(echo $LOGIN | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
echo -e "${GREEN}Token: $TOKEN${NC}"

if [ -z "$TOKEN" ]; then
  echo "Failed to get token"
  exit 1
fi

AUTH="Authorization: Bearer $TOKEN"

# ============================================================================
# STEP 3: Create Subject
# ============================================================================
echo -e "\n${YELLOW}Step 3: Creating subject (Math)...${NC}"
SUBJECT=$(curl -s -X POST "$API/me/subjects" \
  -H "Content-Type: application/json" \
  -H "$AUTH" \
  -d '{
    "title": "Math Basics",
    "description": "Mathematics fundamentals for exam prep"
  }')

SUBJECT_ID=$(echo $SUBJECT | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
echo -e "${GREEN}Subject ID: $SUBJECT_ID${NC}"

# ============================================================================
# STEP 4: Create Topic
# ============================================================================
echo -e "\n${YELLOW}Step 4: Creating topic (Algebra)...${NC}"
TOPIC=$(curl -s -X POST "$API/subjects/$SUBJECT_ID/topics" \
  -H "Content-Type: application/json" \
  -H "$AUTH" \
  -d '{
    "title": "Algebra Basics",
    "description": "Linear equations and basics"
  }')

TOPIC_ID=$(echo $TOPIC | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
echo -e "${GREEN}Topic ID: $TOPIC_ID${NC}"

# ============================================================================
# STEP 5: Create Tasks (with different types)
# ============================================================================
echo -e "\n${YELLOW}Step 5: Creating 5 tasks...${NC}"

declare -a TASK_IDS

for i in {1..5}; do
  TASK=$(curl -s -X POST "$API/topics/$TOPIC_ID/tasks" \
    -H "Content-Type: application/json" \
    -H "$AUTH" \
    -d "{
      \"title\": \"Task $i: Solve equation\",
      \"description\": \"Solve: 2x + 5 = 15\",
      \"type\": \"multiple_choice\",
      \"difficulty\": $((i * 20)),
      \"estimated_time_minutes\": $((i * 2))
    }")

  TASK_ID=$(echo $TASK | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
  TASK_IDS[$i]=$TASK_ID
  echo -e "${GREEN}  Task $i ID: $TASK_ID${NC}"
done

# ============================================================================
# STEP 6: LEARNING MODE - First Session
# ============================================================================
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}LEARNING MODE: First Session${NC}"
echo -e "${BLUE}========================================${NC}"

echo -e "\n${YELLOW}Creating learning session (5 new tasks)...${NC}"
SESSION1=$(curl -s -X POST "$API/sessions" \
  -H "Content-Type: application/json" \
  -H "$AUTH" \
  -d "{
    \"topic_id\": \"$TOPIC_ID\",
    \"task_count\": 5,
    \"mode\": \"learning\"
  }")

SESSION1_ID=$(echo $SESSION1 | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
echo -e "${GREEN}Session 1 ID (Learning): $SESSION1_ID${NC}"
echo "Response: $SESSION1"

# Get tasks for session
echo -e "\n${YELLOW}Getting tasks for learning session...${NC}"
TASKS_L=$(curl -s -X GET "$API/sessions/$SESSION1_ID/tasks" \
  -H "$AUTH")

echo "Tasks in session: $TASKS_L"

# ============================================================================
# STEP 7: Submit Attempts (mixed results - simulate struggling on tasks 4,5)
# ============================================================================
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}Submitting attempts with mixed results${NC}"
echo -e "${BLUE}========================================${NC}"

# Correct answers for tasks 1, 2, 3
for i in {1..3}; do
  echo -e "\n${YELLOW}Task $i: Correct answer${NC}"
  ATTEMPT=$(curl -s -X POST "$API/sessions/$SESSION1_ID/attempts" \
    -H "Content-Type: application/json" \
    -H "$AUTH" \
    -d "{
      \"task_id\": \"${TASK_IDS[$i]}\",
      \"user_answer\": \"5\",
      \"answer_result\": \"correct\"
    }")

  echo "Result: $ATTEMPT"
done

# Wrong answers for tasks 4, 5 (struggling areas!)
for i in {4..5}; do
  echo -e "\n${YELLOW}Task $i: WRONG answer (struggling)${NC}"
  ATTEMPT=$(curl -s -X POST "$API/sessions/$SESSION1_ID/attempts" \
    -H "Content-Type: application/json" \
    -H "$AUTH" \
    -d "{
      \"task_id\": \"${TASK_IDS[$i]}\",
      \"user_answer\": \"wrong\",
      \"answer_result\": \"wrong\"
    }")

  echo "Result: $ATTEMPT"
done

# Finish session
echo -e "\n${YELLOW}Finishing learning session...${NC}"
FINISH=$(curl -s -X POST "$API/sessions/$SESSION1_ID/finish" \
  -H "$AUTH")

echo "Finished: $FINISH"

# ============================================================================
# STEP 8: Get Session Stats
# ============================================================================
echo -e "\n${YELLOW}Getting session stats...${NC}"
STATS=$(curl -s -X GET "$API/sessions/$SESSION1_ID/stats" \
  -H "$AUTH")

echo "Session Stats: $STATS"

# ============================================================================
# STEP 9: Get Task-Level Statistics (view adaptive priorities!)
# ============================================================================
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}Task-Level Statistics${NC}"
echo -e "${BLUE}(Shows which tasks were struggled on!)${NC}"
echo -e "${BLUE}========================================${NC}"

echo -e "\n${YELLOW}Getting task stats for topic (task_learning_stats view)...${NC}"
TASK_STATS=$(curl -s -X GET "$API/topics/$TOPIC_ID/stats" \
  -H "$AUTH")

echo "Task Stats: $TASK_STATS"

# ============================================================================
# STEP 10: REVIEW MODE - Second Session (should prioritize tasks 4,5!)
# ============================================================================
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}REVIEW MODE: Adaptive Second Session${NC}"
echo -e "${BLUE}(Should prioritize struggling tasks 4,5!)${NC}"
echo -e "${BLUE}========================================${NC}"

echo -e "\n${YELLOW}Creating review session...${NC}"
SESSION2=$(curl -s -X POST "$API/sessions" \
  -H "Content-Type: application/json" \
  -H "$AUTH" \
  -d "{
    \"topic_id\": \"$TOPIC_ID\",
    \"task_count\": 3,
    \"mode\": \"review\"
  }")

SESSION2_ID=$(echo $SESSION2 | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
echo -e "${GREEN}Session 2 ID (Review): $SESSION2_ID${NC}"

# Get tasks - SHOULD BE TASKS 4,5 first (the ones user struggled on)!
echo -e "\n${YELLOW}Getting tasks for REVIEW session (check if prioritized!)...${NC}"
TASKS_R=$(curl -s -X GET "$API/sessions/$SESSION2_ID/tasks" \
  -H "$AUTH")

echo -e "${GREEN}Review session tasks (should have tasks 4,5 first):${NC}"
echo "$TASKS_R"

# ============================================================================
# STEP 11: Summary
# ============================================================================
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}SUMMARY: Adaptive Learning Test Complete!${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "\n${GREEN}What happened:${NC}"
echo "1. Created 5 tasks in a topic"
echo "2. First session (LEARNING mode): Got all 5 tasks"
echo "3. Answered correctly: tasks 1, 2, 3"
echo "4. Answered wrong (STRUGGLING): tasks 4, 5"
echo "5. System created repetitions ONLY for wrong answers (4, 5)"
echo "6. Second session (REVIEW mode): Got different set"
echo -e "\n${BLUE}Expected behavior:${NC}"
echo "- Review session should PRIORITIZE tasks 4 and 5 (the struggling ones)"
echo "- Task stats should show tasks 4, 5 have quality='struggling'"
echo "- Task stats should show tasks 1-3 have quality='mastered' (or not show)"
echo -e "\n${GREEN}Check Swagger at: http://localhost:8080/swagger/index.html${NC}"
echo -e "${GREEN}Or test endpoints manually!${NC}\n"
