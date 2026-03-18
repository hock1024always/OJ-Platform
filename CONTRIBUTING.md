# Contributing to OJ Platform

Thank you for your interest in contributing. This document outlines how to get involved.

---

## Ways to Contribute

- **Bug reports** — open an issue describing the problem and steps to reproduce
- **Feature requests** — open an issue with the `enhancement` label
- **Code contributions** — fork the repo, make changes, open a pull request
- **Problem set** — add more problems with driver code and test cases via `scripts/import_leetcode.go`
- **Documentation** — improve docs, fix typos, add examples

---

## Development Setup

### Requirements

- Go 1.21+
- Git
- Make (optional)

### Clone and run

```bash
git clone https://github.com/your-org/oj-platform.git
cd oj-platform
go mod download
go run scripts/import_leetcode.go
go run ./cmd/server/
```

Server starts at `http://localhost:8080`.

---

## Code Structure

| Package | Responsibility |
|---------|---------------|
| `cmd/server` | Main entry point, wires all components |
| `internal/handlers` | HTTP layer — parse request, call service, return response |
| `internal/services` | Business logic |
| `internal/judge` | Compile and execute Go code, compare output |
| `internal/queue` | Worker pool for concurrent judging |
| `internal/models` | GORM structs (User, Problem, TestCase, Submission) |
| `internal/repository` | Database queries |
| `internal/middleware` | JWT authentication, CORS |
| `pkg/config` | Config file loading |
| `web/` | Frontend HTML, CSS, JavaScript |
| `scripts/` | Data import scripts |

---

## Submitting a Pull Request

1. Fork the repository and create a branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes. Keep commits focused and descriptive.

3. Ensure the project builds without errors:
   ```bash
   go build ./...
   ```

4. Test your changes manually or add tests if applicable.

5. Open a pull request against `main` with:
   - A clear title
   - Description of what changed and why
   - Any related issue numbers

---

## Adding a New Problem

Each problem requires:

1. **Title, description, difficulty, tags**
2. **`FunctionTemplate`** — the function signature shown in the editor
3. **`DriverCode`** — `package main` + imports + `func main()` that reads stdin, calls the user function, and prints to stdout
4. **Test cases** — at least one public (shown to user) and one or more hidden

Add the problem to `scripts/import_leetcode.go` following the existing pattern, then re-run the import script.

---

## Code Style

- Follow standard Go formatting (`gofmt`)
- Keep handler functions thin — logic belongs in services
- Avoid adding dependencies without discussion

---

## Reporting Issues

When filing a bug report, include:

- Go version (`go version`)
- OS and architecture
- Steps to reproduce
- Expected vs actual behavior
- Relevant log output from `server.log`

---

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
