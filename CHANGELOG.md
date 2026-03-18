# Changelog

All notable changes to this project will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] — 2026-03-18

Initial public release.

### Added

**Backend**
- Go 1.21 HTTP server using Gin framework
- JWT authentication (register, login, profile)
- GORM ORM with SQLite database
- Auto-migration on startup for all models (User, Problem, TestCase, Submission)
- `Problem` model with `function_template` and `driver_code` fields enabling LeetCode-style judging
- Judge engine: assembles complete Go source from user function + driver code, compiles with `go build`, executes with stdin/stdout, enforces time limit via `time.After` + `Process.Kill()`
- Submission status: Accepted / Wrong Answer / Compile Error / Runtime Error / Time Limit Exceeded / System Error
- Channel-based worker pool (20 goroutines) for concurrent submission processing
- REST API: `/api/v1/register`, `/api/v1/login`, `/api/v1/problems`, `/api/v1/submit`, `/api/v1/submissions/:id`
- CORS middleware permitting all origins (development-friendly)
- Health check endpoint `/health`
- Problem importer script (`scripts/import_leetcode.go`)

**Problem Set — LeetCode Hot 100 subset (15 problems)**
- 两数之和 (Easy) — Array, Hash Table
- 爬楼梯 (Easy) — Dynamic Programming
- 最大子数组和 (Medium) — Array, DP
- 买卖股票的最佳时机 (Easy) — Array, DP
- 只出现一次的数字 (Easy) — Bit Manipulation
- 多数元素 (Easy) — Array, Hash Table
- 移动零 (Easy) — Array, Two Pointers
- 合并两个有序数组 (Easy) — Array, Two Pointers
- 验证回文串 (Easy) — Two Pointers, String
- 找到字符串中所有字母异位词 (Medium) — Sliding Window
- 二叉树的最大深度 (Easy) — Tree, DFS
- 二叉树的中序遍历 (Easy) — Tree, DFS
- 对称二叉树 (Easy) — Tree, BFS
- 有效的括号 (Easy) — Stack, String
- 最长公共前缀 (Easy) — String

Each problem includes: function template, driver code (I/O harness), public example test cases, and hidden test cases.

**Frontend**
- Problem list page with difficulty badges and tag display
- Problem detail page with description and public test cases
- CodeMirror 5 editor with Go syntax highlighting
- Editor pre-filled with problem-specific function template on page load
- Submission result modal with polling (1s interval) until judging completes
- JWT stored in `localStorage`; auth guard on protected pages
- Dynamic API base URL (no hardcoded hostname)

**DevOps**
- Dockerfile (multi-stage build, CGO enabled for SQLite)
- `docker-compose.yml`
- `Makefile`
- `deploy.sh` and `start.sh` helper scripts

**Documentation**
- `README.md` — project overview, quick start, problem table, structure
- `docs/ARCHITECTURE.md` — system design, component descriptions, request flow, technology choices, known limitations
- `docs/API.md` — full REST API reference with request/response examples
- `CONTRIBUTING.md` — development setup, PR process, how to add problems
- `CHANGELOG.md` — this file
- `LICENSE` — MIT

---

## Roadmap

Items planned for future versions:

- **v1.1** — Leaderboard, submission history per user, admin problem management UI
- **v1.2** — PostgreSQL support, rate limiting, pagination improvements
- **v2.0** — Docker sandbox per submission (security isolation), more languages, expanded problem set
