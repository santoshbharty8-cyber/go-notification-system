# 🏗️ System Architecture — Go Notification & Event Processing System

## 🎯 Overview

This system is a distributed, event-driven backend service designed to process high-volume events reliably using:

- ⚡ Asynchronous processing
- 🧵 Concurrent workers (goroutines)
- 📦 Redis as a distributed queue
- 🔁 Retry + exponential backoff
- 🧠 Idempotency (atomic SET NX)
- 📬 Dead Letter Queue (DLQ)
- 🛑 Rate limiting (Sliding Window)

## 🧠 Architectural Goals

- High throughput → handle a large number of events
- Low latency → API responds quickly
- Fault tolerance → retry + DLQ
- Scalability → horizontal scaling
- Consistency → idempotent processing
- Resilience → system continues even on failures

## 🏗️ High-Level Architecture

```text
          ┌──────────────┐
          │   Client     │
          └──────┬───────┘
                 │
                 ▼
        ┌──────────────────┐
        │   API Layer      │
        │ (Go HTTP Server) │
        └──────┬───────────┘
               │
               ▼
      ┌────────────────────┐
      │  Rate Limiter      │
      │ (Sliding Window)   │
      └──────┬─────────────┘
             │
             ▼
      ┌────────────────────┐
      │   Redis Queue      │
      │   (LPUSH / BRPOP)  │
      └──────┬─────────────┘
             │
   ┌─────────▼─────────┐
   │   Worker Pool     │
   │ (Goroutines)      │
   └─────────┬─────────┘
             │
             ▼
   ┌────────────────────┐
   │  Idempotency Layer │
   │ (SET NX in Redis)  │
   └─────────┬──────────┘
             │
             ▼
   ┌────────────────────┐
   │ Processing Logic   │
   └─────────┬──────────┘
             │
     ┌───────▼─────────┐
     │ Retry Mechanism │
     │ (Backoff)       │
     └───────┬────────┘
             │
      ┌──────▼───────┐
      │ Success      │
      │ OR           │
      │ DLQ          │
      └──────────────┘
```

## 🔄 End-to-End Flow

1. Client sends an event to the `/event` API.
2. API validates the input.
3. Rate limiter checks request frequency.
4. Event is pushed to the Redis queue.
5. Worker pulls the event using `BRPOP`.
6. Worker performs an idempotency check with `SET NX`.
7. If the event is a duplicate, it is skipped.
8. If it is new, it is processed.
9. On failure, the system retries with backoff.
10. After max retries, the event is sent to the DLQ.

## 🧩 Core Components

### 1. API Layer
- Accept incoming events.
- Validate the request.
- Return a fast response.
- Push the event to Redis queue.

**Design choice:** keep the API stateless and lightweight.

### 2. Rate Limiter (Sliding Window)
- Prevent abuse and traffic spikes.
- Track timestamps per client.
- Allow only `N` requests per time window.

**Why Sliding Window:**
- More accurate than fixed window.
- Prevents burst traffic.

### 3. Redis Queue
- Decouple API from processing.
- Enable async architecture.

**Operations:**
- `LPUSH` → enqueue event
- `BRPOP` → blocking dequeue

**Why Redis:**
- Fast in-memory performance.
- Simple setup.
- Supports distributed systems.

### 4. Worker Pool
- Process events concurrently.
- Multiple goroutines consume from the queue.

**Benefits:**
- Parallel processing
- Efficient CPU utilization
- Scalable with more workers

### 5. Idempotency Layer
Duplicate events can occur due to retries, network failures, or duplicate requests.

**Solution:**
```text
SET key value NX EX ttl
```

**Behavior:**
- First worker → allowed.
- Other workers → skipped.

**Result:**
- No duplicate processing.
- Safe distributed execution.

### 6. Event Processing Layer
Handles business logic such as:

- `user_registered` → send welcome email
- `order_created` → trigger order workflow

This layer should be easy to extend for new event types.

### 7. Retry Mechanism
- Max retries: `3`
- Backoff: `1s → 2s → 4s`

**Why:**
- Handles transient failures.
- Avoids immediate system overload.

### 8. Dead Letter Queue (DLQ)
- Stores permanently failed events.
- Prevents data loss.
- Helps with debugging and reprocessing.

## ⚡ Concurrency Model

- Uses goroutines as lightweight workers.
- Workers run independently.
- Shared Redis queue ensures load distribution.

## 📈 Scalability Design

### Horizontal Scaling
```text
Instance 1 ─┐
Instance 2 ─┼── Redis Queue ── Workers
Instance 3 ─┘
```

- Multiple instances can run.
- All share the same Redis queue.

### Stateless API
- No session storage.
- Easy to scale behind a load balancer.

## 🧠 Data Flow Summary

```text
API → Redis → Worker → Idempotency → Process → Retry → DLQ
```

## ⚖️ Trade-offs

| Component | Choice | Trade-off |
|---|---|---|
| Queue | Redis | Not as powerful as Kafka |
| Idempotency | SET NX | Needs TTL tuning |
| Rate Limiter | In-memory | Not distributed |
| Retry | In-worker | Simpler but less flexible |

## 🔒 Reliability Guarantees

| Feature | Guarantee |
|---|---|
| Retry | Handles transient failures |
| DLQ | Prevents data loss |
| Idempotency | Prevents duplicates |
| Redis | Durable queue (if configured) |

## 🚀 Future Improvements

- Kafka instead of Redis
- Distributed rate limiter (Redis-based)
- Observability with Prometheus + Grafana
- OpenTelemetry tracing
- WebSocket notifications
- Circuit breaker pattern

