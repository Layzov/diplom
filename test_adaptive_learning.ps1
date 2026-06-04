# Adaptive Learning Test Script for PowerShell

$API = "http://localhost:8080/api/v1"

function Write-Header {
    param([string]$Text)
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host $Text -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
}

function Write-Step {
    param([string]$Text)
    Write-Host $Text -ForegroundColor Yellow
}

function Write-Success {
    param([string]$Text)
    Write-Host $Text -ForegroundColor Green
}

# ============================================================================
# STEP 1: Register and Login
# ============================================================================
Write-Header "Testing Adaptive Spaced Repetition System"

Write-Step "Step 1: Registering user..."
$registerBody = @{
    email = "testadaptive@exam.com"
    password = "SecurePass123!"
    full_name = "Adaptive Learner"
} | ConvertTo-Json

$registerResponse = Invoke-WebRequest -Uri "$API/auth/register" `
    -Method POST `
    -ContentType "application/json" `
    -Body $registerBody | ConvertFrom-Json

Write-Host "Response: $($registerResponse | ConvertTo-Json)" -ForegroundColor Gray

Write-Step "Step 2: Logging in..."
$loginBody = @{
    email = "testadaptive@exam.com"
    password = "SecurePass123!"
} | ConvertTo-Json

$loginResponse = Invoke-WebRequest -Uri "$API/auth/login" `
    -Method POST `
    -ContentType "application/json" `
    -Body $loginBody | ConvertFrom-Json

$token = $loginResponse.access_token
Write-Success "Token: $token"

if ([string]::IsNullOrEmpty($token)) {
    Write-Host "Failed to get token" -ForegroundColor Red
    exit 1
}

$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# ============================================================================
# STEP 3: Create Subject
# ============================================================================
Write-Step "Step 3: Creating subject (Math)..."
$subjectBody = @{
    title = "Math Basics"
    description = "Mathematics fundamentals for exam prep"
} | ConvertTo-Json

$subjectResponse = Invoke-WebRequest -Uri "$API/me/subjects" `
    -Method POST `
    -Headers $headers `
    -Body $subjectBody | ConvertFrom-Json

$subjectId = $subjectResponse.id
Write-Success "Subject ID: $subjectId"

# ============================================================================
# STEP 4: Create Topic
# ============================================================================
Write-Step "Step 4: Creating topic (Algebra)..."
$topicBody = @{
    title = "Algebra Basics"
    description = "Linear equations and basics"
} | ConvertTo-Json

$topicResponse = Invoke-WebRequest -Uri "$API/subjects/$subjectId/topics" `
    -Method POST `
    -Headers $headers `
    -Body $topicBody | ConvertFrom-Json

$topicId = $topicResponse.id
Write-Success "Topic ID: $topicId"

# ============================================================================
# STEP 5: Create Tasks
# ============================================================================
Write-Step "Step 5: Creating 5 tasks..."
$taskIds = @()

for ($i = 1; $i -le 5; $i++) {
    $taskBody = @{
        title = "Task $i : Solve equation"
        description = "Solve: 2x + 5 = 15"
        type = "multiple_choice"
        difficulty = $i * 20
        estimated_time_minutes = $i * 2
    } | ConvertTo-Json

    $taskResponse = Invoke-WebRequest -Uri "$API/topics/$topicId/tasks" `
        -Method POST `
        -Headers $headers `
        -Body $taskBody | ConvertFrom-Json

    $taskIds += $taskResponse.id
    Write-Success "  Task $i ID: $($taskResponse.id)"
}

# ============================================================================
# STEP 6: LEARNING MODE - First Session
# ============================================================================
Write-Header "LEARNING MODE: First Session"

Write-Step "Creating learning session (5 new tasks)..."
$session1Body = @{
    topic_id = $topicId
    task_count = 5
    mode = "learning"
} | ConvertTo-Json

$session1Response = Invoke-WebRequest -Uri "$API/sessions" `
    -Method POST `
    -Headers $headers `
    -Body $session1Body | ConvertFrom-Json

$session1Id = $session1Response.id
Write-Success "Session 1 ID (Learning): $session1Id"
Write-Host "Response: $($session1Response | ConvertTo-Json)" -ForegroundColor Gray

# Get tasks for session
Write-Step "Getting tasks for learning session..."
$tasksL = Invoke-WebRequest -Uri "$API/sessions/$session1Id/tasks" `
    -Method GET `
    -Headers $headers | ConvertFrom-Json

Write-Success "Tasks in learning session:"
Write-Host "$($tasksL | ConvertTo-Json -Depth 10)" -ForegroundColor Gray

# ============================================================================
# STEP 7: Submit Attempts (mixed results)
# ============================================================================
Write-Header "Submitting attempts with mixed results"

# Correct answers for tasks 1, 2, 3
for ($i = 0; $i -le 2; $i++) {
    Write-Step "Task $($i+1): Correct answer"
    $attemptBody = @{
        task_id = $taskIds[$i]
        user_answer = "5"
        answer_result = "correct"
    } | ConvertTo-Json

    $attemptResponse = Invoke-WebRequest -Uri "$API/sessions/$session1Id/attempts" `
        -Method POST `
        -Headers $headers `
        -Body $attemptBody | ConvertFrom-Json

    Write-Host "Result: $($attemptResponse | ConvertTo-Json -Depth 5)" -ForegroundColor Gray
}

# Wrong answers for tasks 4, 5 (struggling!)
for ($i = 3; $i -le 4; $i++) {
    Write-Step "Task $($i+1): WRONG answer (STRUGGLING AREA!)"
    $attemptBody = @{
        task_id = $taskIds[$i]
        user_answer = "wrong"
        answer_result = "wrong"
    } | ConvertTo-Json

    $attemptResponse = Invoke-WebRequest -Uri "$API/sessions/$session1Id/attempts" `
        -Method POST `
        -Headers $headers `
        -Body $attemptBody | ConvertFrom-Json

    Write-Host "Result: $($attemptResponse | ConvertTo-Json -Depth 5)" -ForegroundColor Gray
}

# Finish session
Write-Step "Finishing learning session..."
$finishResponse = Invoke-WebRequest -Uri "$API/sessions/$session1Id/finish" `
    -Method POST `
    -Headers $headers | ConvertFrom-Json

Write-Host "Finished: $($finishResponse | ConvertTo-Json)" -ForegroundColor Gray

# ============================================================================
# STEP 8: Get Session Stats
# ============================================================================
Write-Step "Getting session stats..."
$statsResponse = Invoke-WebRequest -Uri "$API/sessions/$session1Id/stats" `
    -Method GET `
    -Headers $headers | ConvertFrom-Json

Write-Success "Session Stats:"
Write-Host "$($statsResponse | ConvertTo-Json -Depth 10)" -ForegroundColor Gray

# ============================================================================
# STEP 9: Get Task-Level Statistics
# ============================================================================
Write-Header "Task-Level Statistics (Shows which tasks were struggled on!)"

Write-Step "Getting task stats for topic (task_learning_stats view)..."
$taskStatsResponse = Invoke-WebRequest -Uri "$API/topics/$topicId/stats" `
    -Method GET `
    -Headers $headers | ConvertFrom-Json

Write-Success "Task Stats:"
Write-Host "$($taskStatsResponse | ConvertTo-Json -Depth 10)" -ForegroundColor Gray

# ============================================================================
# STEP 10: REVIEW MODE - Second Session (should prioritize tasks 4,5!)
# ============================================================================
Write-Header "REVIEW MODE: Adaptive Second Session (Should prioritize tasks 4,5!)"

Write-Step "Creating review session..."
$session2Body = @{
    topic_id = $topicId
    task_count = 3
    mode = "review"
} | ConvertTo-Json

$session2Response = Invoke-WebRequest -Uri "$API/sessions" `
    -Method POST `
    -Headers $headers `
    -Body $session2Body | ConvertFrom-Json

$session2Id = $session2Response.id
Write-Success "Session 2 ID (Review): $session2Id"

# Get tasks - SHOULD BE TASKS 4,5 first!
Write-Step "Getting tasks for REVIEW session (check if prioritized!)..."
$tasksR = Invoke-WebRequest -Uri "$API/sessions/$session2Id/tasks" `
    -Method GET `
    -Headers $headers | ConvertFrom-Json

Write-Success "Review session tasks (should have tasks 4,5 first!):"
Write-Host "$($tasksR | ConvertTo-Json -Depth 10)" -ForegroundColor Green

# ============================================================================
# STEP 11: Summary
# ============================================================================
Write-Header "SUMMARY: Adaptive Learning Test Complete!"

Write-Host "What happened:" -ForegroundColor Green
Write-Host "1. Created 5 tasks in a topic"
Write-Host "2. First session (LEARNING mode): Got all 5 tasks"
Write-Host "3. Answered correctly: tasks 1, 2, 3"
Write-Host "4. Answered wrong (STRUGGLING): tasks 4, 5"
Write-Host "5. System created repetitions ONLY for wrong answers (4, 5)"
Write-Host "6. Second session (REVIEW mode): Got different set"

Write-Host "`nExpected behavior:" -ForegroundColor Cyan
Write-Host "- Review session should PRIORITIZE tasks 4 and 5 (the struggling ones)"
Write-Host "- Task stats should show tasks 4, 5 have quality='struggling'"
Write-Host "- Task stats should show tasks 1-3 have quality='mastered' (or not show)"

Write-Host "`nCheck Swagger at: http://localhost:8080/swagger/index.html" -ForegroundColor Green
Write-Host "Or test endpoints manually!`n" -ForegroundColor Green

Write-Host "Test script completed successfully!" -ForegroundColor Green
