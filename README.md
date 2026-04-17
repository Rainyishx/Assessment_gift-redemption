# Gift Redemption System

A lightweight Go REST API for managing department Christmas gift redemptions.

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Running Locally](#running-locally)
- [Running with Docker](#running-with-docker)
- [API Reference](#api-reference)
- [Running Tests](#running-tests)
- [Assumptions & Design Decisions](#assumptions--design-decisions)

---

## Overview

Each team sends one representative to redeem their gift. The representative shows their staff pass (unique ID). The system:

1. **Looks up** the staff pass ID to find the team name
2. **Verifies** the team hasn't already redeemed their gift
3. **Records** the redemption if eligible, or rejects it if not

---

## Architecture

The project follows a clean layered architecture:

```
HTTP Request
     │
     ▼
┌─────────────┐
│   Handler   │  Parses HTTP requests, writes HTTP responses
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Service   │  Core business logic (look up → verify → record)
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────┐
│          Repository             │
│  StaffRepo    RedemptionRepo    │  Data access (read & write CSV files)
└─────────────────────────────────┘
```

Each layer depends only on interfaces, making the system easy to test and extend (e.g. swap CSV for a database without touching service or handler code)

---

## Project Structure

```
gift-redemption/
├── cmd/
|    ├──main.go                              # Entry point, wires everything together
├── data/
│   ├── staffmapping.csv                     # Input: staff → team mappings
│   └── redemptions.csv                      # Auto-created: redemption records
└── internal/
    ├── model/
    │   └── model.go                         # Shared data structs
    ├── repository/
    │   ├── staff_repo.go                    # Loads staff CSV into memory
    │   ├── staff_repo_test.go
    │   ├── redemption_repo.go               # Reads/writes redemptions CSV
    │   └── redemption_repo_test.go
    ├── service/
    │   ├── redemption_service.go            # Business logic
    │   └── redemption_service_test.go
    └── handler/
        ├── handler.go                       # HTTP handlers
        └── handler_test.go
    ├── Dockerfile
    ├── README.md
    ├── go.mod
```

---

## Prerequisites

- Go 1.22 or later 
- (Optional) Docker

---

## Running Locally

### 1. Clone the repository

```bash
git clone https://github.com/Rainyishx/Assessment_gift-redemption.git
cd Assessment_gift-redemption
```

### 2. Run the server

```bash
go run ./cmd/main.go
```

The server starts on **http://localhost:8080** by default.

## Running with Docker

### Build the image

```bash
docker build -t gift-redemption-app .
```

### Run the container

```bash
docker run -p 8080:8080 gift-redemption-app
```

## API Reference

### `GET /health`

Health check endpoint.

**Response `200`:**
```json
{ "status": "ok" }
```

---

### `POST /redeem`

Attempt to redeem a gift using a staff pass ID.

**Request body:**
```json
{ "staff_pass_id": "STAFF_001" }
```

**Success `201`:**
```json
{
  "team_name": "TEAM_A",
  "redeemed_at": 1700000000000,
  "message": "Gift redeemed successfully for team TEAM_A"
}
```

**Error responses:**

| Status | Reason |
|--------|--------|
| `400` | Missing or invalid `staff_pass_id` |
| `404` | Staff pass ID not found in the system |
| `409` | Team has already redeemed their gift |
| `500` | Internal server error |

**Error body example:**
```json
{ "error": "TEAM_A has already redeemed their gift" }
```

---

### Example `curl` commands

In a separate terminal, run these commands to test the implementation.

```bash
# Health check
curl http://localhost:8080/health

# Valid redemption
curl -X POST http://localhost:8080/redeem \
  -H "Content-Type: application/json" \
  -d '{"staff_pass_id": "STAFF_001"}'

# Try same team again → 409
curl -X POST http://localhost:8080/redeem \
  -H "Content-Type: application/json" \
  -d '{"staff_pass_id": "STAFF_002"}'

# Unknown staff pass → 404
curl -X POST http://localhost:8080/redeem \
  -H "Content-Type: application/json" \
  -d '{"staff_pass_id": "unknown"}'

# Missing staff_pass_id → 400
curl -X POST http://localhost:8080/redeem \
  -H "Content-Type: application/json" \
  -d '{}'
```

---

### Example `PowerShell` commands
```PowerShell
Invoke-RestMethod -Uri "http://localhost:8080/redeem" -Method POST -ContentType "application/json" -Body '{"staff_pass_id": "STAFF_001"}'
```
---

## Running Tests

```bash
# Run all tests
go test ./... -v

# Run a specific package
go test ./internal/service/... -v
```

---

## Assumptions & Design Decisions

### Storage
- **Staff mappings** are loaded from CSV into an in-memory `map[string]StaffMapping` at startup. Lookups are O(1).
- **Redemptions** are stored in a CSV file (`data/redemptions.csv`), loaded into memory at startup, and flushed to disk on every new redemption. This keeps the approach simple and auditable without requiring a database.

### Concurrency
- The `RedemptionRepository` uses a `sync.Mutex` to protect in-memory state and file writes. This makes the API safe for concurrent requests.

### Graceful Shutdown
- It is assumed that the server may be stopped at any point during operation. For example, 
via `Ctrl+C` locally or a `docker stop` in a containerised environment. This function is implemented for easy testing as it prevents redemption writes from being interrupted mid-flight, which could otherwise leave the `redemptions.csv` file in a partially written state.

### HTTP method routing
- Uses Go 1.22's enhanced `http.ServeMux` which supports method-based routing (`POST /redeem`) without any third-party router.

### Idempotency
- Redemptions are **not idempotent** — submitting the same staff pass twice returns `409 Conflict` on the second call (since the team already redeemed).

### Environment configuration
- File paths and port are configurable via environment variables with sensible defaults, making the app easy to run in both local and containerised environments.


---

## Administrative Note: Commit History
Due to a local Git CLI misconfiguration related to a university email address, early commits in this repository appear under the GitHub handle `rainyishz`. The CLI was reconfigured mid-development, and subsequent commits reflect my primary handle `Rainyishx`. All code was authored locally by me, and I can provide local Git logs to substantiate this if required.
