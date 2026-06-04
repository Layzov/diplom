@echo off
REM ===========================================================================
REM Exam Prep Backend — Complete Demo Flow (Windows Version)
REM
REM Usage:
REM   demo.bat
REM   demo.bat -delay 2    # Add 2-second pause between requests
REM
REM Prerequisites:
REM   1. PostgreSQL running: docker compose up -d
REM   2. Seed data: go run .\cmd\seed
REM   3. API running: go run .\cmd\api
REM   4. curl installed (or use WSL)
REM
REM ===========================================================================

setlocal enabledelayedexpansion

REM Configuration
set BASE_URL=http://localhost:8080
set API_PREFIX=%BASE_URL%/api/v1
set DELAY=0

REM Parse arguments
if "%1"=="-delay" (
  set DELAY=%2
)

REM Colors setup (Windows 10+)
for /F %%A in ('echo prompt $H ^| cmd') do set "BS=%%A"

setlocal enabledelayedexpansion

echo.
echo ===========================================================================
echo  Demo: Exam Prep Backend - Complete Flow
echo ===========================================================================
echo.
echo  API URL: %API_PREFIX%
echo  Delay: %DELAY%s between requests
echo.

REM 1. Health Check
echo.
echo [1] Health Check
echo ----
curl -s -X GET "%API_PREFIX%/health" -H "Content-Type: application/json"
echo.
timeout /t %DELAY% /nobreak > nul

REM Generate unique email
for /f "tokens=2-4 delims=/ " %%a in ('date /t') do (set mydate=%%c%%a%%b)
for /f "tokens=1-2 delims=/:" %%a in ('time /t') do (set mytime=%%a%%b)
set EMAIL=demo_%mydate%%mytime%@example.com

REM 2. Register User
echo.
echo [2] Register User
echo ----
echo Email: %EMAIL%
set REGISTER_JSON={^
  "email": "%EMAIL%",^
  "password": "demo123456",^
  "name": "Demo User"^
}
for /f "delims=" %%i in ('curl -s -X POST "%API_PREFIX%/auth/register" -H "Content-Type: application/json" -d "%REGISTER_JSON%" ^| jq -r ".data.id"') do set USER_ID=%%i
echo User ID: %USER_ID%
timeout /t %DELAY% /nobreak > nul

REM 3. Login
echo.
echo [3] Login
echo ----
set LOGIN_JSON={^
  "email": "%EMAIL%",^
  "password": "demo123456"^
}
for /f "delims=" %%i in ('curl -s -X POST "%API_PREFIX%/auth/login" -H "Content-Type: application/json" -d "%LOGIN_JSON%" ^| jq -r ".data.token"') do set TOKEN=%%i
echo Token: %TOKEN:~0,50%...
timeout /t %DELAY% /nobreak > nul

REM 4. Create Subject
echo.
echo [4] Create Subject
echo ----
set SUBJECT_JSON={^
  "title": "Демонстрационная Математика",^
  "description": "Полный цикл демонстрации"^
}
for /f "delims=" %%i in ('curl -s -X POST "%API_PREFIX%/subjects" -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "%SUBJECT_JSON%" ^| jq -r ".data.id"') do set SUBJECT_ID=%%i
echo Subject ID: %SUBJECT_ID%
timeout /t %DELAY% /nobreak > nul

REM 5. Create Topic
echo.
echo [5] Create Topic
echo ----
set TOPIC_JSON={^
  "subject_id": "%SUBJECT_ID%",^
  "title": "Основные понятия"^
}
for /f "delims=" %%i in ('curl -s -X POST "%API_PREFIX%/topics" -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "%TOPIC_JSON%" ^| jq -r ".data.id"') do set TOPIC_ID=%%i
echo Topic ID: %TOPIC_ID%
timeout /t %DELAY% /nobreak > nul

REM 6. Create Task
echo.
echo [6] Create Task
echo ----
set TASK_JSON={^
  "topic_id": "%TOPIC_ID%",^
  "type": "practice",^
  "title": "Решите задачу",^
  "content": "Найдите x в уравнении"^
}
for /f "delims=" %%i in ('curl -s -X POST "%API_PREFIX%/tasks" -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "%TASK_JSON%" ^| jq -r ".data.id"') do set TASK_ID=%%i
echo Task ID: %TASK_ID%
timeout /t %DELAY% /nobreak > nul

REM 7. Create Session
echo.
echo [7] Create Session
echo ----
for /f %%i in ('powershell -Command "Get-Date -AsUTC -Format o"') do set NOW=%%i
set SESSION_JSON={^
  "subject_id": "%SUBJECT_ID%",^
  "started_at": "%NOW%"^
}
for /f "delims=" %%i in ('curl -s -X POST "%API_PREFIX%/sessions" -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "%SESSION_JSON%" ^| jq -r ".data.id"') do set SESSION_ID=%%i
echo Session ID: %SESSION_ID%
timeout /t %DELAY% /nobreak > nul

REM 8. Add Task Attempts
echo.
echo [8] Add Task Attempts
echo ----

set ATTEMPT_JSON={^
  "session_id": "%SESSION_ID%",^
  "task_id": "%TASK_ID%",^
  "result": "correct",^
  "time_ms": 15000^
}
echo Adding: Correct attempt...
curl -s -X POST "%API_PREFIX%/task_attempts" -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "%ATTEMPT_JSON%" | jq .
timeout /t %DELAY% /nobreak > nul

echo.
set ATTEMPT_JSON={^
  "session_id": "%SESSION_ID%",^
  "task_id": "%TASK_ID%",^
  "result": "partial",^
  "time_ms": 22000^
}
echo Adding: Partial attempt...
curl -s -X POST "%API_PREFIX%/task_attempts" -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "%ATTEMPT_JSON%" | jq .
timeout /t %DELAY% /nobreak > nul

REM 9. Complete Session
echo.
echo [9] Complete Session
echo ----
set COMPLETE_JSON={^
  "ended_at": "%NOW%"^
}
curl -s -X POST "%API_PREFIX%/sessions/%SESSION_ID%/complete" -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "%COMPLETE_JSON%" | jq .
timeout /t %DELAY% /nobreak > nul

REM 10. Get Statistics
echo.
echo [10] Get User Statistics
echo ----
curl -s -X GET "%API_PREFIX%/me/stats" -H "Authorization: Bearer %TOKEN%" | jq .
timeout /t %DELAY% /nobreak > nul

REM 11. Get Topic Statistics
echo.
echo [11] Get Topic Statistics
echo ----
curl -s -X GET "%API_PREFIX%/me/stats/topics" -H "Authorization: Bearer %TOKEN%" | jq .
timeout /t %DELAY% /nobreak > nul

REM 12. Get Upcoming Repetitions
echo.
echo [12] Get Upcoming Repetitions
echo ----
curl -s -X GET "%API_PREFIX%/me/repetitions/upcoming" -H "Authorization: Bearer %TOKEN%" | jq .
timeout /t %DELAY% /nobreak > nul

REM 13. Get Profile
echo.
echo [13] Get User Profile
echo ----
curl -s -X GET "%API_PREFIX%/me" -H "Authorization: Bearer %TOKEN%" | jq .
timeout /t %DELAY% /nobreak > nul

REM Summary
echo.
echo ===========================================================================
echo  DEMO COMPLETE!
echo ===========================================================================
echo.
echo  User Email: %EMAIL%
echo  User ID: %USER_ID%
echo  Subject: %SUBJECT_ID%
echo  Topic: %TOPIC_ID%
echo  Task: %TASK_ID%
echo  Session: %SESSION_ID%
echo.
echo ===========================================================================
pause
