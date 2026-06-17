# Prince University – Student Records RESTful API

## Stack
- **Go** + **Gin** (web framework)
- **GORM** ORM + **PostgreSQL** (`gorm.io/driver/postgres`)
- **Docker Compose** for the PostgreSQL database (and the API)

## Architecture (layered)
```
controller/  HTTP handlers (Gin)        -> request/response only
service/     business logic             -> total credit score, honor-roll rule, enroll
repository/  data access (GORM)         -> queries against PostgreSQL
model/       domain entities + view DTO
database/    connection, migration, seeding
main.go      wires database -> repository -> service -> controller
```

## Domain model
A `Student` can be registered to many `Course`s, and each `Course` can have
many `Student`s — a many-to-many relationship (join table `student_courses`).

## Run

### Option A — Postgres in Docker, API on host
```bash
docker compose up -d db          # starts PostgreSQL (host port 5433 -> container 5432)
DB_PORT=5433 go run .            # API on http://localhost:8080
```

### Option B — everything in Docker
```bash
docker compose up --build        # db + api together; API on http://localhost:8080
```

DB connection is configured via env vars (defaults in parentheses):
`DB_HOST` (localhost), `DB_PORT` (5432), `DB_USER` (prince),
`DB_PASSWORD` (prince), `DB_NAME` (student_records), `PORT` (8080).

The schema is auto-migrated and the University's existing data is seeded on
first run.

## Endpoints

### 1. List all students — `GET /api/students`
All students with their course(s) and the **computed total credit score**,
sorted **ascending by last name**.
```bash
curl http://localhost:8080/api/students
```

### 2. Honor Roll — `GET /api/students/honor-roll`
Students registered for at least one course with `creditScore >= 5`.
```bash
curl http://localhost:8080/api/students/honor-roll
```

### 3. Enroll a new student — `POST /api/students`
Registers a new student into a course by `courseId`.
```bash
curl -X POST http://localhost:8080/api/students \
  -H "Content-Type: application/json" \
  -d '{
        "matricNumber": "E776",
        "firstName": "Jane",
        "lastName": "Dougherty",
        "dateOfAdmission": "2026-06-15",
        "courseId": 3
      }'
```

## Inspect the database (for the DB-tables screenshot)
```bash
docker exec -it prince_pg psql -U prince -d student_records \
  -c "SELECT * FROM students ORDER BY student_id;" \
  -c "SELECT * FROM courses ORDER BY course_id;" \
  -c "SELECT * FROM student_courses;"
```
