# TaskFlow

TaskFlow is a modern, collaborative Kanban-style task management application built with a high-performance Go backend and a responsive React frontend.

## 1. Overview

TaskFlow enables teams to manage projects through a visual Kanban interface with support for multi-player collaboration, role-based access control, and real-time-friendly pessimistic/optimistic UI updates.

### Tech Stack
- **Backend**: Go 1.22+, Chi Router, pgx (Postgres Driver), golang-migrate.
- **Frontend**: React 19, Vite, Tailwind CSS v4, TanStack Query, Shadcn/UI.
- **Database**: PostgreSQL 15.
- **DevOps**: Docker & Docker Compose.

---

## 2. Architecture Decisions

### Multi-Player Core
Unlike basic task managers, TaskFlow is built from the ground up for teams. Introduced a `team_members` relationship which allows project owners to collaborate with other registered users.

### High-Performance Data Retrieval
To avoid the "N+1 query" problem, the backend uses **Common Table Expressions (CTEs)** and **LEFT JOINs** to fetch tasks along with their assignee metadata (Name, Email, Avatars) in a single database round-trip.

### Hybrid UI Updates
- **Optimistic UI**: Drag-and-drop operations immediately reflect on the board for zero-latency feel, with automatic state rollback on server failures.
- **Pessimistic UI**: Critical operations like "Delete" or "Add Project" use standard loading states to ensure data integrity.

### Role-Based Access Control (RBAC)
The backend enforces strict ownership checks:
- **Project Owners**: Full control over projects and all tasks.
- **Task Assignees**: Can update task status but cannot modify descriptions or delete tasks.
- **Others**: View-only access if they are part of the project team.

---

## 3. Running Locally

You only need **Docker** and **Git** installed.

```bash
# 1. Clone the repository
git clone https://github.com/your-username/taskflow-rahul
cd taskflow-rahul

# 2. Setup environment variables (contains sensible defaults)
cp .env.example .env

# 3. Spin up the full stack (DB, API, Frontend)
docker compose up --build
```

- **Frontend**: [http://localhost:3000](http://localhost:3000)
- **API**: [http://localhost:4000](http://localhost:4000)
- **Database**: `localhost:5432`

---

## 4. Migrations

Migrations run **automatically** on container start via the Go `InitDB` process. 
- Source: `backend/migrations/`
- Tool: `golang-migrate`

To manually run migrations if needed:
```bash
# Up
docker exec -it taskflow-api ./taskflow-api -migrate=up
# Down
docker exec -it taskflow-api ./taskflow-api -migrate=down
```

---

## 5. Test Credentials

The database is pre-seeded with a testing account:

- **Email**: `test@example.com`
- **Password**: `password123`

---

## 6. API Reference

All protected endpoints require an `Authorization: Bearer <token>` header.

### Auth
- `POST /auth/register` - Create account
- `POST /auth/login` - Get JWT token

### Projects
- `GET /projects` - List projects (supports `?page=X&limit=Y` pagination)
- `POST /projects` - Create project
- `GET /projects/{id}` - Get project details + tasks
- `GET /projects/{id}/stats` - Get task distribution stats

### Tasks
- `POST /projects/{id}/tasks` - Create task
- `PATCH /tasks/{id}` - Update task details or status
- `DELETE /tasks/{id}` - Delete task

### Teams
- `GET /team` - List your team members
- `POST /team` - Add member by email

---

## 7. What I'd Do With More Time

1. **E2E Testing**: Implement Playwright or Cypress tests for the drag-and-drop flows.
2. **WebSockets**: Replace optimistic polling with real-time updates so team members see changes instantly.
3. **File Attachments**: Add support for S3-backed task attachments.
4. **Activity Logs**: Implement an audit trail for task history (who moved what and when).
5. **Caching**: Introduce Redis to cache frequent `/stats` and `/projects` requests.
