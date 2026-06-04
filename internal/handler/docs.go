package handler

import (
	"net/http"
	"strings"
)

// DocsHandler returns a handler that serves both Swagger UI and OpenAPI spec
func DocsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Serve OpenAPI spec as JSON
		if strings.HasSuffix(r.URL.Path, "/swagger.json") || strings.HasSuffix(r.URL.Path, ".json") {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(openAPISpec))
			return
		}

		// Serve Swagger UI HTML (default)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(swaggerHTML))
	}
}

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Exam Prep API - Swagger UI</title>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui-bundle.js"></script>
    <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui.css">
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin: 0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script>
        SwaggerUIBundle({
            url: '/docs/swagger.json',
            dom_id: '#swagger-ui',
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.SwaggerUIStandalonePreset
            ],
            layout: 'BaseLayout',
            requestInterceptor: (request) => {
                request.headers['X-Request-ID'] = Math.random().toString(36).substr(2, 9);
                return request;
            }
        });
    </script>
</body>
</html>`

// openAPISpec contains the embedded OpenAPI 3.0 specification
// This is the same content as docs/openapi.yaml
const openAPISpec = `openapi: 3.0.0
info:
  title: Exam Prep Backend API
  description: >
    Веб-приложение для эффективной подготовки к экзаменам с использованием 
    метода интервальных повторений.
  version: 1.0.0
  contact:
    name: API Support
    url: https://github.com/Layzov/diplom
  license:
    name: MIT

servers:
  - url: http://localhost:8080
    description: Local development server
  - url: https://api.example.com
    description: Production server

tags:
  - name: Health
    description: System health checks
  - name: Authentication
    description: User registration and authentication (JWT)
  - name: Profile
    description: User profile and personal statistics
  - name: Subjects
    description: Study subjects (Математика, Физика, etc)
  - name: Topics
    description: Topics within subjects
  - name: Tasks
    description: Tasks/flashcards for learning
  - name: Sessions
    description: Study sessions and attempt tracking
  - name: Repetitions
    description: Spaced repetition scheduling

paths:
  /health:
    get:
      tags: [Health]
      summary: Health check
      description: Check if API is running and healthy
      operationId: getHealth
      responses:
        '200':
          description: API is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: ok

  /api/v1/auth/register:
    post:
      tags: [Authentication]
      summary: Register new user
      description: Create a new user account and return JWT token
      operationId: registerUser
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [email, password]
              properties:
                email:
                  type: string
                  format: email
                  example: user@example.com
                password:
                  type: string
                  format: password
                  minLength: 6
                  example: SecurePassword123
                name:
                  type: string
                  example: John Doe
      responses:
        '200':
          description: User registered successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'
        '400':
          description: Invalid request or user already exists

  /api/v1/auth/login:
    post:
      tags: [Authentication]
      summary: Login user
      description: Authenticate user with email and password, return JWT token
      operationId: loginUser
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [email, password]
              properties:
                email:
                  type: string
                  format: email
                password:
                  type: string
                  format: password
      responses:
        '200':
          description: Login successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'
        '401':
          description: Invalid credentials

  /api/v1/me/profile:
    get:
      tags: [Profile]
      summary: Get user profile
      operationId: getProfile
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Profile retrieved successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    $ref: '#/components/schemas/UserProfile'
        '401':
          description: Unauthorized

    put:
      tags: [Profile]
      summary: Update user profile
      operationId: updateProfile
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
      responses:
        '200':
          description: Profile updated successfully

  /api/v1/me/subjects:
    get:
      tags: [Profile]
      summary: List user's subjects
      operationId: listUserSubjects
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Subjects retrieved

    post:
      tags: [Subjects]
      summary: Create subject
      operationId: createSubject
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [title]
              properties:
                title:
                  type: string
                description:
                  type: string
      responses:
        '201':
          description: Subject created

  /api/v1/subjects/{id}:
    get:
      tags: [Subjects]
      summary: Get subject
      operationId: getSubject
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Subject retrieved

    put:
      tags: [Subjects]
      summary: Update subject
      operationId: updateSubject
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                title:
                  type: string
      responses:
        '200':
          description: Subject updated

    delete:
      tags: [Subjects]
      summary: Delete subject
      operationId: deleteSubject
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: Subject deleted

  /api/v1/subjects/{id}/topics:
    get:
      tags: [Topics]
      summary: List topics by subject
      operationId: listTopicsBySubject
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Topics retrieved

    post:
      tags: [Topics]
      summary: Create topic
      operationId: createTopic
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [title]
              properties:
                title:
                  type: string
      responses:
        '201':
          description: Topic created

  /api/v1/topics/{id}:
    get:
      tags: [Topics]
      summary: Get topic
      operationId: getTopic
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Topic retrieved

    put:
      tags: [Topics]
      summary: Update topic
      operationId: updateTopic
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                title:
                  type: string
      responses:
        '200':
          description: Topic updated

    delete:
      tags: [Topics]
      summary: Delete topic
      operationId: deleteTopic
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: Topic deleted

  /api/v1/topics/{id}/tasks:
    get:
      tags: [Tasks]
      summary: List tasks by topic
      operationId: listTasksByTopic
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Tasks retrieved

    post:
      tags: [Tasks]
      summary: Create task
      operationId: createTask
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [title, content, type]
              properties:
                title:
                  type: string
                content:
                  type: string
                type:
                  type: string
      responses:
        '201':
          description: Task created

  /api/v1/tasks/{id}:
    get:
      tags: [Tasks]
      summary: Get task
      operationId: getTask
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Task retrieved

    put:
      tags: [Tasks]
      summary: Update task
      operationId: updateTask
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                title:
                  type: string
      responses:
        '200':
          description: Task updated

    delete:
      tags: [Tasks]
      summary: Delete task
      operationId: deleteTask
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: Task deleted

  /api/v1/sessions:
    post:
      tags: [Sessions]
      summary: Create study session
      operationId: createSession
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                subject_id:
                  type: string
                  format: uuid
      responses:
        '201':
          description: Session created

  /api/v1/sessions/{id}:
    get:
      tags: [Sessions]
      summary: Get session
      operationId: getSession
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Session retrieved

  /api/v1/sessions/{id}/finish:
    post:
      tags: [Sessions]
      summary: Finish session
      operationId: finishSession
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Session finished

  /api/v1/sessions/{id}/attempts:
    get:
      tags: [Sessions]
      summary: Get session attempts
      operationId: listAttempts
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Attempts retrieved

    post:
      tags: [Sessions]
      summary: Submit task attempt
      operationId: submitAttempt
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [task_id, result]
              properties:
                task_id:
                  type: string
                  format: uuid
                result:
                  type: string
                  enum: [correct, wrong, partial, skipped]
      responses:
        '201':
          description: Attempt recorded

  /api/v1/sessions/{id}/stats:
    get:
      tags: [Sessions]
      summary: Get session statistics
      operationId: getSessionStats
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Statistics retrieved

  /api/v1/me/stats:
    get:
      tags: [Profile]
      summary: Get user learning statistics overview
      operationId: getStatsOverview
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Statistics retrieved

  /api/v1/me/stats/topics:
    get:
      tags: [Profile]
      summary: Get statistics by topic
      operationId: getTopicStats
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Topic statistics retrieved

  /api/v1/me/stats/sessions:
    get:
      tags: [Profile]
      summary: Get statistics by session
      operationId: getSessionStats
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Session statistics retrieved

  /api/v1/me/stats/upcoming:
    get:
      tags: [Profile]
      summary: Get upcoming repetitions
      operationId: getUpcomingRepetitions
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Upcoming repetitions retrieved

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  schemas:
    AuthResponse:
      type: object
      properties:
        status:
          type: string
        data:
          type: object
          properties:
            id:
              type: string
              format: uuid
            email:
              type: string
            access_token:
              type: string

    UserProfile:
      type: object
      properties:
        id:
          type: string
          format: uuid
        email:
          type: string
        name:
          type: string
`
