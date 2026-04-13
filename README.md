# TaskFlow - Full Stack

A modern, collaborative Kanban-style task management application where users can register, manage projects, create tasks, and assign them to team members. Built with a high-performance Go backend and a responsive React frontend, fully containerized for a seamless setup.

---

## Tech Stack

| Layer          | Technology                                  |
| -------------- | ------------------------------------------- |
| **Backend**    | Go 1.22+, Chi Router, pgx (Postgres Driver) |
| **Frontend**   | React 19, Vite, Tailwind CSS v4, Shadcn/UI  |
| **State Mgt.** | TanStack Query (React Query)                |
| **Database**   | PostgreSQL 15                               |
| **Migrations** | golang-migrate                              |
| **Auth**       | JWT (golang-jwt) + bcrypt                   |
| **Logging**    | slog (structured JSON)                      |
| **Deployment** | Docker, Docker Compose (Multi-stage builds) |

---

## Project Structure

### Backend Architecture
```text
backend/
├── cmd/api/main.go          # Entry point — server bootstrap & graceful shutdown
├── internal/
│   ├── database/            # Connection pooling & migration runner
│   ├── handler/             # HTTP handlers (REST logic & validation)
│   ├── middleware/          # JWT authentication middleware
│   ├── models/              # Request/Response & Database structs
│   └── repository/          # Data access layer (Postgres queries)
├── migrations/              # SQL migration files (up + down)
├── pkg/utils/               # Constants & Auth utilities (JWT/Bcrypt)
└── Dockerfile               # Multi-stage Go build
```

### Frontend Architecture
```text
frontend/
├── src/
│   ├── components/
│   │   ├── layout/          # Navigation & shared wrappers
│   │   ├── ui/              # Atomized Shadcn/UI components
│   │   └── ...              # Domain components (TaskCard, TeamManager)
│   ├── context/             # Auth & Global state management
│   ├── lib/                 # API client (Axios) & Tailwind utilities
│   ├── pages/               # Routed view components (Dashboard, Login, etc.)
│   ├── App.tsx              # Root routing mapping
│   └── main.tsx             # React mount point
├── public/                  # Favicons & static assets
├── index.html               # SPA entry point
└── Dockerfile               # Multi-stage Node/Nginx build
```

---

## Architecture Decisions

### Collaborative Core
- **Multi-Player Support**: Introduced a `team_members` relationship which allows project owners to collaborate with other registered users securely.
- **High-Performance Joins**: To avoid the "N+1 query" problem, the backend uses **Common Table Expressions (CTEs)** and **LEFT JOINs** to fetch tasks along with their assignee metadata in a single database round-trip.
- **Role-Based Access Control (RBAC)**: Strict backend ownership checks. Project owners have full control; task assignees can only update statuses; team members gain view-only access.
- **No ORM**: Using `database/sql` patterns (via `pgxpool`) directly with parameterized queries avoids ORM overhead, prevents hidden N+1 issues, and keeps SQL explicit and readable.

### UI/UX Implementation
- **Hybrid State Updates**:
  - **Optimistic UI**: Drag-and-drop operations immediately reflect on the board for a zero-latency feel, with automatic state rollback on server failures (via TanStack Query).
  - **Pessimistic UI**: Critical operations like creating projects or deleting tasks use standard loading states to ensure data integrity.

### Tradeoffs

- **Simple custom pagination**: Using `LIMIT` and `OFFSET` instead of keyset/cursor-based pagination. While cursor-based is strictly superior for huge datasets, `LIMIT`/`OFFSET` fits this medium scale nicely and is much easier to implement and reason about.
- **No repository interfaces for mocking**: Repositories are concrete types injected into handlers. Traded interface-based mocking for immediate simplicity and reduced abstraction overhead.
- **No Websockets**: Real-time multi-player sync relies on polling/refetching rather than WebSockets, intentionally done to simplify the current scope's infrastructure.
- **No Graceful DB Retry on Go Startup**: The API relies on Docker Compose's `restart: on-failure` pattern rather than polling in code with exponential backoff on startup.

---

## Running Locally

Assumes Docker and Docker Compose are installed. No Go or Node.js installation required.

```bash
git clone https://github.com/HLrahul/taskflow-rahul.git
cd taskflow-rahul
cp .env.example .env
docker compose up --build -d
```

- **Frontend**: [http://localhost:3000](http://localhost:3000)
- **API**: [http://localhost:4000](http://localhost:4000)
- **Database**: `localhost:5432`

Verify API health:

```bash
curl http://localhost:4000/
# {"message":"TaskFlow API is running!"}
```

To stop:

```bash
docker compose down
```

To reset everything (wipe DB and start entirely fresh):

```bash
docker compose down -v
docker compose up --build -d
```

---

## Running Migrations

Migrations run **automatically on container startup**. No manual steps required. 

The Go backend spins up, ping checks the database, and calls `golang-migrate` to execute pending migrations and seeds before the HTTP server ever begins listening.

Migration files are located at `backend/migrations/`:

| File                        | Description                                                     |
| --------------------------- | --------------------------------------------------------------- |
| `000001_init_schema.up.sql` | Creates users, projects, tasks, enums, & fkeys (pgcrypto/uuid)  |
| `000002_add_teams.up.sql`   | Creates `team_members` for multi-player tracking                |
| `000003_seed_data.up.sql`   | Inserts 1 test user, 1 project, and 3 tasks safely (idempotent) |

Each migration has a corresponding `.down.sql` for rollback processes.

---

## Test Credentials

A set of seed data is inserted automatically on first run, meaning you don't even have to register to start testing:

```text
Email:    test@example.com
Password: password123
```

The seed inherently provides:

- **1 User:** The test user above.
- **1 Project:** "Ship TaskFlow".
- **3 Tasks:** Demo tasks with differing statuses (`todo`, `in_progress`, and `done`).

---

## API Reference

All endpoints return `Content-Type: application/json`. Protected endpoints require an `Authorization: Bearer <token>` header. A Postman collection (`TaskFlow.postman_collection.json`) is included in the project root.

### Authentication

#### Register (`POST /auth/register`)
```json
// Request
{ "name": "John Doe", "email": "john@example.com", "password": "securepassword1" }

// Response (201 Created)
{ "token": "jwt-token", "user": { "id": "uuid", "name": "John Doe", "email": "..." } }
```

#### Login (`POST /auth/login`)
```json
// Request
{ "email": "test@example.com", "password": "password123" }

// Response (200 OK)
{ "token": "jwt-token", "user": { "id": "uuid", "name": "Test User", "email": "..." } }
```

### Projects

#### List Projects (`GET /projects?page=1&limit=10`)
Returns projects owned by the user or projects where they have assigned tasks.
```json
// Response (200 OK)
{
  "projects": [{ "id": "uuid", "name": "Ship TaskFlow", "description": "..." }],
  "total": 1, "page": 1, "limit": 10
}
```

#### Create Project (`POST /projects`)
```json
// Request
{ "name": "New Project", "description": "Optional" }
```

#### Get Project Details (`GET /projects/{id}`)
Returns project metadata along with all associated tasks.

#### Update Project (`PATCH /projects/{id}`)
*Owner only.* Partially update project name or description.

#### Delete Project (`DELETE /projects/{id}`)
*Owner only.* Cascading delete for project and all its tasks.

#### Project Stats (`GET /projects/{id}/stats`)
```json
// Response (200 OK)
{
  "total": 5,
  "by_status": { "todo": 2, "in_progress": 1, "done": 2 },
  "by_assignee": { "user-uuid": 4, "unassigned": 1 }
}
```

### Tasks

#### List Project Tasks (`GET /projects/{id}/tasks`)
Detail view of tasks within a specific project.

#### Create Task (`POST /projects/{id}/tasks`)
```json
// Request
{
  "title": "Design UI",
  "status": "todo",
  "priority": "high",
  "assignee_id": "user-uuid",
  "due_date": "2026-04-15T00:00:00Z"
}
```

#### Update Task (`PATCH /tasks/{id}`)
Assignee or Owner only. Update status, priority, title, etc.

#### Delete Task (`DELETE /tasks/{id}`)
*Owner only.* Permanently remove a task.

### Team Collaboration

#### List Your Team (`GET /team`)
Returns all users currently in your broad collaborative team.

#### Add Team Member (`POST /team`)
```json
// Request
{ "email": "collaborator@example.com" }
```

---

## Error Responses

All errors follow a standardized format to simplify frontend handling:

```json
// 400 - Validation Failed
{
  "error": "validation failed",
  "fields": { "email": "is required", "password": "too short" }
}

// 401 - Unauthorized
{ "error": "unauthorized" }

// 403 - Forbidden (Insufficient Permissions)
{ "error": "forbidden" }

// 404 - Not Found
{ "error": "not found" }
```

---

## Environment Variables

| Variable | Default | Description |
| :--- | :--- | :--- |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `taskflow` | Database name |
| `JWT_SECRET` | `supersecret123` | Secret for signing JWT tokens |
| `PORT` | `4000` | API server port |
| `BCRYPT_COST` | `12` | Hashing cost (range 4-31) |
| `VITE_API_URL` | `http://localhost:4000` | API endpoint for Frontend build |

---

## What I'd Do With More Time

### Things I'd add next

- **Testing:** Implement testing for both frontend and backend.
- **WebSockets / Server-Sent Events:** Replace optimistic UI / periodic fetching with fully real-time backend updates so multiple team members editing the same project see changes instantaneously.
- **Request Validation Library:** Right now, validation is relatively manual mapped in the handlers. Adopting a library like `go-playground/validator` would reduce boilerplate for more complex schemas.
- **Activity Feed / Audit Logs:** Record changes like "User X moved Task Y to Done" or "User Z updated Description" to provide a full paper trail for teams tracking history.
- **File Attachments:** Allow users to upload screenshots and attach them to specific tasks (storing assets in an AWS S3 bucket or a local MinIO container).
- **Graceful DB Retries on Startup:** Right now, Docker Compose's `restart: on-failure` handles Postgres startup delays gracefully via container restarts. Building a reliable exponential backoff routine directly into Go's connection logic would be more resilient in standalone deployments or orchestrators.
- **Structured Request IDs:** Enhance the structured `slog` logging by injecting a correlation UUID into context via middleware, making tracing a single request end-to-end trivial in production debug scenarios.
- **DB indices:** DB call latency increases in large datasets ( Basically on scale ), so adding indices to the DB would help improve performance.
