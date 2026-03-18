# Architecture

This document describes the technical design of OJ Platform v1.0.

---

## System Overview

```
┌──────────────────────────────────────────────────────────┐
│                        Browser                           │
│         (HTML + Vanilla JS + CodeMirror editor)          │
└───────────────────────┬──────────────────────────────────┘
                        │ HTTP / REST
┌───────────────────────▼──────────────────────────────────┐
│                   Gin HTTP Server                        │
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐ │
│  │  Middleware  │  │   Handlers   │  │  Static Files   │ │
│  │  JWT / CORS  │  │  (thin layer)│  │  web/ directory │ │
│  └─────────────┘  └──────┬───────┘  └─────────────────┘ │
└─────────────────────────┬┴──────────────────────────────┘
                          │
┌─────────────────────────▼──────────────────────────────┐
│                      Services                          │
│  ┌──────────────────┐    ┌─────────────────────────┐   │
│  │   UserService    │    │     JudgeService         │   │
│  │  register/login  │    │  submit → queue → result │   │
│  └──────────────────┘    └────────────┬────────────┘   │
└───────────────────────────────────────┼────────────────┘
                                        │
┌───────────────────────────────────────▼────────────────┐
│                    Task Queue                          │
│         Worker Pool (20 goroutines)                    │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ...           │
│  │ Worker 1 │ │ Worker 2 │ │ Worker 3 │               │
└──────┬───────────────────────────────────────────────┘
       │
┌──────▼──────────────────────────────────────────────┐
│                  Judge Engine                       │
│  1. Write user_code + driver_code → /tmp/oj_*/      │
│  2. go build -o binary main.go                      │
│  3. Run binary with stdin = test_case.input         │
│  4. Compare stdout with expected output             │
│  5. Return: Accepted / Wrong Answer / Compile Error │
│             Runtime Error / Time Limit Exceeded     │
└──────┬──────────────────────────────────────────────┘
       │
┌──────▼──────────────────────────────────────────────┐
│                Repository / Database                │
│          GORM ORM → SQLite (v1) / PostgreSQL        │
│  Tables: users, problems, test_cases, submissions   │
└─────────────────────────────────────────────────────┘
```

---

## Component Descriptions

### HTTP Layer — Gin

Routes are registered in `internal/routes/routes.go`. The handler layer is intentionally thin: parse and validate the request, call the appropriate service, return a JSON response.

All API routes are prefixed with `/api/v1`. Static frontend files are served from the `web/` directory.

### Authentication — JWT

JWT tokens are issued at login and validated by the `AuthRequired()` middleware on protected routes. Tokens carry `user_id` and `username` claims, expire in 24 hours.

### Judge Engine (`internal/judge/judge.go`)

The core of the system. For each test case:

1. A temporary directory is created under `/tmp/oj_*`
2. The final source file is assembled:
   ```
   package main          ← from driver_code
   import (...)          ← from driver_code
   
   func solution(...) {  ← user code, inserted here
       ...
   }
   
   func main() {         ← from driver_code
       ...               ← reads stdin, calls solution, prints stdout
   }
   ```
3. `go build` compiles the file
4. The binary runs with test case input piped to stdin
5. stdout is compared to the expected output (trimmed)
6. The temporary directory is removed

Timeout is enforced via `time.After` + `Process.Kill()`.

### LeetCode-style Interface

Each problem stores two extra fields:

| Field | Purpose |
|-------|---------|
| `function_template` | Shown in the editor as the starting template |
| `driver_code` | `package main` + I/O boilerplate, invisible to the user |

The user writes only the solution function. The platform assembles the complete compilable file before judging.

### Task Queue (`internal/queue/queue.go`)

A simple channel-based worker pool. Submissions are queued as `Task` structs. Twenty goroutines process tasks concurrently, which is sufficient for 10–100 simultaneous users given typical Go compilation time (~1–2s per task).

### Data Models

```
User
  id, username, email, password_hash

Problem
  id, title, description, difficulty, tags
  time_limit, memory_limit
  function_template, driver_code

TestCase
  id, problem_id, input, output, is_public

Submission
  id, user_id, problem_id, code, language
  status, result, time_used, memory_used
```

---

## Request Flow: Code Submission

```
POST /api/v1/submit
    │
    ▼
JudgeHandler.SubmitCode()
    │  validate token, parse body
    ▼
JudgeService.Submit()
    │  create Submission record (status=Pending)
    │  push Task to queue
    ▼
[async goroutine]
    │
    ▼
TaskQueue → Worker picks up task
    │
    ▼
JudgeService.handleTask()
    │  load Problem + TestCases from DB
    │  for each test case:
    │      judge.RunGo(userCode, input, expectedOutput, driverCode)
    │      if not Accepted → return early
    ▼
Update Submission record (Accepted / Wrong Answer / ...)
    │
    ▼
Client polls GET /api/v1/submissions/:id every 1s
    │  until status != Pending
    ▼
Display result in browser
```

---

## Technology Choices

| Concern | Choice | Rationale |
|---------|--------|-----------|
| Language | Go 1.21 | Single binary, low overhead, native concurrency |
| HTTP | Gin | Minimal, fast, widely used in Go ecosystem |
| ORM | GORM | Reduces boilerplate, easy migration |
| Database | SQLite (v1) | Zero-ops for initial deployment; GORM makes PostgreSQL swap easy |
| Auth | JWT (HS256) | Stateless, simple to implement |
| Frontend | Vanilla JS | No build step, easy to modify |
| Editor | CodeMirror 5 | Syntax highlighting, lightweight, CDN-served |
| Execution | `os/exec` + `go build` | Direct, no Docker overhead in v1 |

---

## Known Limitations (v1)

- **No sandbox isolation** — user code runs as the server process user; malicious code could access the filesystem. Planned fix: Docker-per-submission in v2.
- **SQLite** — not suitable for high write concurrency in production; switch to PostgreSQL for multi-instance deployment.
- **No rate limiting** — a single user can flood the submission queue.
- **Output comparison is exact string match** — whitespace-sensitive; some problems may need custom comparators.
