# 🚀 Go Notification & Event Processing System

A production-grade, event-driven backend system built using **Golang + Redis**, designed to handle high-throughput event processing with **concurrency, reliability, and scalability**.

# 🟢 CI / Coverage Status

![CI](https://github.com/santoshbharty8-cyber/go-notification-system/actions/workflows/ci.yml/badge.svg)
![Coverage](https://codecov.io/gh/santoshbharty8-cyber/go-notification-system/branch/main/graph/badge.svg)

> 📊 Test Coverage: **94%+**

## 🎯 Key Features

- ⚡ Event-driven architecture
- 🔁 Retry mechanism with exponential backoff
- 🧵 Concurrent worker pool (goroutines)
- 📦 Redis-based distributed queue
- 🛑 Rate limiting (Sliding Window algorithm)
- 🧠 Idempotency (atomic SET NX)
- 📬 Dead Letter Queue (DLQ)
- 🔒 Security scanning (gosec + govulncheck)
- 🧪 High test coverage (~94%)
- 🚀 CI/CD pipeline (GitHub Actions)
- 🐳 Docker-ready setup

## 🧠 System Architecture

```text
Client
   ↓
API Layer (Go HTTP Server)
   ↓
Rate Limiter (Sliding Window)
   ↓
Redis Queue (LPUSH)
   ↓
Worker Pool (Goroutines)
   ↓
Idempotency Check (SET NX)
   ↓
Processing Logic
   ↓
Retry + Backoff
   ↓
Success OR DLQ
```

## 🧱 Tech Stack

- **Language:** Go (Golang)
- **Queue:** Redis
- **Concurrency:** Goroutines + Channels
- **Rate Limiting:** Sliding Window (DSA)
- **Testing:** Go testing + Integration tests
- **CI/CD:** GitHub Actions
- **Security:** gosec + govulncheck
- **Containerization:** Docker

# Project Structure

```text
go-notification-system/
│
├── cmd/api/                  # Entry point
├── internal/
│   ├── handlers/             # API handlers
│   ├── services/             # Business logic
│   ├── queue/                # Redis queue + DLQ
│   ├── workers/              # Worker pool + retry
│   ├── idempotency/          # SETNX logic
│   ├── ratelimiter/          # Sliding window
│   ├── models/               # Data models
│   ├── middleware/           # Logging + rate limit
│   ├── redisclient/          # Redis connection
│   ├── config/               # App config
│   └── bootstrap/            # App initialization
│
├── tests/
│   ├── integration/          # End-to-end tests
│   └── helpers/              # Test utilities
│
├── deployments/docker/       # Docker setup
├── scripts/                  # Test scripts
```

## ⚙️ Setup Instructions

### 1. Clone Repo
```bash
git clone https://github.com/santoshbharty8-cyber/go-notification-system.git
cd go-notification-system
```

### 2. Start Redis
```bash
docker run -d -p 6379:6379 redis
```

### 3. Run Application
```bash
go run cmd/api/main.go
```
## 🧪 Run Tests (Local)

```bash
./scripts/test.sh
```
### What it runs
- Tests.
- Coverage check.

## 📊 Coverage Report

```bash
go tool cover -html=coverage.out
```

👉 Opens a detailed HTML coverage visualization in your browser.

## 🧪 API Usage

### Create Event
```bash
curl -X POST http://localhost:8080/event \
-H "Content-Type: application/json" \
-d '{
  "id": "101",
  "type": "order_created",
  "timestamp": 1710000000,
  "payload": {
    "order_id": "ORD101"
  }
}'
```

## 🧠 Key Concepts Explained

### Retry Mechanism
- Retries failed events up to `maxRetries`.
- Uses exponential backoff: `1s → 2s → 4s`.
- Prevents system overload.

### Idempotency (SET NX)
- Prevents duplicate processing.
- Uses Redis atomic operation.
- Ensures only one worker processes an event.

### Dead Letter Queue (DLQ)
- Stores events that fail after retries.
- Prevents data loss.
- Enables debugging & replay.

### Rate Limiting
- Sliding window algorithm.
- Protects API from abuse.
- Per-IP request control

## 🚀 System Flow

```text
API → Rate Limit → Redis Queue → Worker → Idempotency → Process → Retry → DLQ
```

## 🔥 Scalability Design

- Stateless API layer.
- Redis enables distributed workers.
- Horizontal scaling supported.
- Fault-tolerant processing

## 🔐 Security
- Static analysis via `gosec`.
- Vulnerability scanning via `govulncheck`.
- Safe retry logic with no overflow risk.
- Secure randomness using `crypto/rand`.

## 🚀 CI/CD Pipeline
GitHub Actions pipeline includes:
- ✔ Lint (`golangci-lint`)
- ✔ Security Scan (`gosec`)
- ✔ Vulnerability Scan (`govulncheck`)
- ✔ Tests + Coverage (~94%)
- ✔ Docker Build

## 📈 Future Improvements

- Kafka integration
- Prometheus + Grafana monitoring
- Distributed rate limiter (Redis)
- WebSocket notifications
- Kubernetes deployment


## 👨‍💻 Author

**Santosh Kumar Bharty**  
Backend Developer (**Golang | Python | System Design**)