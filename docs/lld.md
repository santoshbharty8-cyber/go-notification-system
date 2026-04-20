# 🔧 Low-Level Design (LLD)
## Go Notification & Event Processing System

## 🎯 Objective

Provide a code-level design of the system, covering:

- Data models
- Module structure
- Function-level flow
- Concurrency handling
- Idempotency implementation
- Retry and DLQ logic

## 📁 Project Structure (Detailed)

```text
internal/
├── handlers/        # HTTP layer
├── services/        # Business logic
├── queue/           # Redis queue producer
├── workers/         # Worker consumers
├── idempotency/     # SET NX logic
├── ratelimiter/     # Sliding window
├── middleware/      # HTTP middleware
├── models/          # Data structures
├── bootstrap/       # App initialization
```

## 🧱 Core Data Models

### 📌 Event Model

```go
type Event struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Timestamp int64                  `json:"timestamp"`
    Payload   map[string]interface{} `json:"payload"`
}
```

**Purpose:**
- Represents an incoming event
- Supports a flexible payload design

### 📌 API Response

```go
type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

## 🌐 API Layer (Handler)

### 📌 EventHandler

**Responsibilities:**
- Decode request
- Validate input
- Push to Redis

**Flow:**

```text
HTTP Request
 → Decode JSON
 → Validate Event
 → Push to Redis
 → Return Response
```

**Key Code Flow:**

```go
event := decodeRequest()
validate(event)
queue.PushToRedisQueue(event)
return success
```

## 📦 Queue Layer (Redis Producer)

### 📌 Function: PushToRedisQueue

```go
func PushToRedisQueue(event Event) error {
    data := json.Marshal(event)
    redis.LPush(queueName, data)
}
```

**Design:**
- Uses Redis List
- Non-blocking push

## 🧵 Worker Layer (Consumer)

### 📌 Worker Loop

```text
loop forever:
  BRPOP Redis
  → deserialize event
  → processWithRetry()
```

### 📌 StartRedisWorker

```go
for {
    data := BRPOP()
    event := unmarshal(data)
    processWithRetry(event)
}
```

## 🔁 Retry Mechanism

### 📌 Function: processWithRetry

```text
for attempt in maxRetries:
  try process
  if success → return
  else → sleep(backoff)
send to DLQ
```

### 📌 Backoff Logic

```text
sleep = base * (2 ^ attempt)
```

**Example:**
- Attempt 1 → 1s
- Attempt 2 → 2s
- Attempt 3 → 4s

## 🧠 Idempotency Layer (Critical)

### 📌 Function: TryMarkProcessing

```text
SET key value NX EX ttl
```

### 📌 Flow

```text
SETNX success → process
SETNX fail → skip
```

### 📌 Redis Key Design

```text
processed_event:{event_id}
```

### 📌 Why Atomic?

Without `SETNX`:

```text
Worker 1 → check
Worker 2 → check
→ both process ❌
```

With `SETNX`:

```text
Only one worker wins ✅
```

## 📬 Dead Letter Queue (DLQ)

### 📌 Flow

```text
retry fails → push to DLQ
```

### 📌 Function

```go
func PushToDLQ(event Event)
```

## 🧠 Rate Limiter (Sliding Window)

### 📌 Data Structure

```go
map[clientID][]timestamp
```

### 📌 Logic

- Remove old timestamps
- If count > limit → reject
- Else → allow

### 📌 Thread Safety

Uses `sync.Mutex`.

## 🔄 Full Execution Flow (Detailed)

```text
Client
 → API Handler
 → Validation
 → Rate Limiter
 → Redis Queue (LPUSH)
 → Worker (BRPOP)
 → Idempotency (SET NX)
 → Process Event
 → Retry (if fail)
 → Success OR DLQ
```

## 🧵 Concurrency Handling

**Worker Model:**
- Multiple goroutines
- Each worker runs an infinite loop
- Shared Redis queue

**Why Safe?**
- Redis handles synchronization
- No shared memory conflicts

## ⚠️ Edge Cases Handling

### 1. Duplicate Event
Handled by idempotency.

### 2. Worker Crash
Event stays in Redis.

### 3. Retry Failure
Event goes to DLQ.

### 4. Redis Failure
System stops.

**Mitigation:**
- Redis cluster
- Retry connection

## ⚙️ Configuration (Suggested)

```env
MAX_RETRIES=3
WORKER_COUNT=3
RATE_LIMIT=5
WINDOW_SIZE=10s
REDIS_URL=localhost:6379
```

## 🔥 Performance Considerations

| Area | Optimization |
|---|---|
| Redis | Use pipelining |
| Workers | Increase count |
| JSON | Optimize serialization |
| Rate limiter | Use Redis version |

## 🧠 Trade-offs

| Feature | Decision | Reason |
|---|---|---|
| Queue | Redis | Simplicity |
| Retry | In worker | Low latency |
| Idempotency | SET NX | Atomic |
| Rate limiter | In-memory | Fast |

