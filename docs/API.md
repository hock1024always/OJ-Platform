# API Reference

Base URL: `http://<host>:8080/api/v1`

All request and response bodies are JSON. Protected endpoints require a `Bearer` token in the `Authorization` header.

---

## Authentication

### Register

```
POST /register
```

**Request body:**

```json
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "YourPassword123"
}
```

**Response `200`:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com"
  }
}
```

**Response `400`:** username or email already exists.

---

### Login

```
POST /login
```

**Request body:**

```json
{
  "username": "alice",
  "password": "YourPassword123"
}
```

**Response `200`:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "<jwt>",
    "username": "alice",
    "user_id": 1
  }
}
```

The token is valid for 24 hours. Pass it in subsequent requests:

```
Authorization: Bearer <token>
```

---

### Get Profile

```
GET /profile
Authorization: Bearer <token>
```

**Response `200`:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "created_at": "2026-03-18T10:00:00+08:00"
  }
}
```

---

## Problems

### List Problems

```
GET /problems?page=1&pageSize=20
```

**Query parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `pageSize` | int | 20 | Items per page |

**Response `200`:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "problems": [
      {
        "id": 1,
        "title": "两数之和",
        "difficulty": "Easy",
        "tags": "数组,哈希表",
        "time_limit": 5000,
        "memory_limit": 256
      }
    ],
    "page": 1,
    "pageSize": 20
  }
}
```

---

### Get Problem

```
GET /problems/:id
```

**Response `200`:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "problem": {
      "id": 1,
      "title": "两数之和",
      "description": "给定一个整数数组 nums ...",
      "difficulty": "Easy",
      "tags": "数组,哈希表",
      "time_limit": 5000,
      "memory_limit": 256,
      "function_template": "func twoSum(nums []int, target int) []int {\n    // 请在此实现你的代码\n}"
    },
    "testCases": [
      {
        "id": 1,
        "input": "2 7 11 15\n9",
        "output": "0 1",
        "is_public": true
      }
    ]
  }
}
```

Only test cases with `is_public: true` are returned. Hidden test cases are used for judging but not exposed.

---

### Create Problem *(admin)*

```
POST /problems
Authorization: Bearer <token>
```

**Request body:**

```json
{
  "title": "Problem Title",
  "description": "Problem description...",
  "difficulty": "Easy",
  "tags": "Array,Hash Table",
  "time_limit": 5000,
  "memory_limit": 256,
  "function_template": "func solve(n int) int {\n    \n}",
  "driver_code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n    var n int\n    fmt.Scan(&n)\n    fmt.Println(solve(n))\n}"
}
```

**Response `200`:** created problem object.

---

## Submissions

### Submit Solution

```
POST /submit
Authorization: Bearer <token>
```

**Request body:**

```json
{
  "problem_id": 1,
  "language": "Go",
  "code": "func twoSum(nums []int, target int) []int {\n    m := make(map[int]int)\n    for i, v := range nums {\n        if j, ok := m[target-v]; ok {\n            return []int{j, i}\n        }\n        m[v] = i\n    }\n    return nil\n}"
}
```

**Response `200`:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 42,
    "status": "Pending",
    "created_at": "2026-03-18T11:39:40+08:00"
  }
}
```

The submission is processed asynchronously. Poll the result endpoint.

---

### Get Submission Result

```
GET /submissions/:id
Authorization: Bearer <token>
```

**Response `200`:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 42,
    "user_id": 1,
    "problem_id": 1,
    "language": "Go",
    "status": "Accepted",
    "result": "",
    "time_used": 312,
    "memory_used": 0,
    "created_at": "2026-03-18T11:39:40+08:00",
    "updated_at": "2026-03-18T11:39:41+08:00"
  }
}
```

**Possible `status` values:**

| Status | Description |
|--------|-------------|
| `Pending` | Queued, not yet judged |
| `Accepted` | All test cases passed |
| `Wrong Answer` | Output does not match expected |
| `Compile Error` | Code failed to compile; error in `result` field |
| `Runtime Error` | Program crashed; stderr in `result` field |
| `Time Limit Exceeded` | Execution exceeded the time limit |
| `System Error` | Internal platform error |

---

## Health Check

```
GET /health
```

**Response `200`:**

```json
{
  "code": 200,
  "message": "success",
  "data": { "status": "ok" }
}
```

---

## Error Responses

All errors follow this format:

```json
{
  "code": 400,
  "message": "error description"
}
```

| HTTP Code | Meaning |
|-----------|---------|
| 400 | Bad request / validation error |
| 401 | Missing or invalid token |
| 404 | Resource not found |
| 500 | Internal server error |
