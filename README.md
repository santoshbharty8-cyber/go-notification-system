# 🚀 Go Notification & Event Processing System

A production-grade, event-driven backend system built using **Golang + Redis**, designed to handle high-throughput event processing with **concurrency, reliability, and scalability**.

## 🎯 Key Features

- ⚡ Event-driven architecture
- 🔁 Retry mechanism with exponential backoff
- 🧵 Concurrent worker pool (goroutines)
- 📦 Redis-based distributed queue
- 🛑 Rate limiting (Sliding Window algorithm)
- 🧠 Idempotency (atomic SET NX)
- 📬 Dead Letter Queue (DLQ)
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
- **Containerization:** Docker

## 📁 Project Structure

```text
go-notification-system/
│
├── cmd/api/                  # Entry point
├── internal/
│   ├── handlers/             # API handlers
│   ├── services/             # Business logic
│   ├── queue/                # Redis queue
│   ├── workers/              # Worker pool
│   ├── idempotency/          # SETNX logic
│   ├── ratelimiter/          # Sliding window
│   ├── models/               # Data models
│   ├── middleware/           # Rate limiter middleware
│   └── bootstrap/            # App initialization
│
├── pkg/utils/                # Helper utilities
```

## ⚙️ Setup Instructions

### 1. Clone Repo
```bash
git clone <your-repo-url>
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

### Idempotency (SET NX)
- Prevents duplicate processing.
- Ensures only one worker processes an event.

### Dead Letter Queue (DLQ)
- Stores events that fail after retries.
- Prevents data loss.

### Rate Limiting
- Sliding window algorithm.
- Protects API from abuse.

## 🚀 System Flow

```text
API → Redis Queue → Worker → Idempotency → Process → Retry → DLQ
```

## 🔥 Scalability Design

- Stateless API layer.
- Redis enables distributed workers.
- Horizontal scaling supported.

## 📈 Future Improvements

- Kafka integration
- Prometheus + Grafana monitoring
- Distributed rate limiter (Redis)
- WebSocket notifications
- CI/CD pipeline

## 🎯 Interview Highlights

This project demonstrates:

- System design (event-driven architecture)
- Concurrency in Go
- Distributed systems using Redis
- Real-world DSA implementation
- Fault tolerance (retry + DLQ)
- Idempotency handling

## 👨‍💻 Author

**Santosh Kumar Bharty**  
Backend Developer (**Golang | Python | System Design**)