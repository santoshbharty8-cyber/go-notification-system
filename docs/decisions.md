# 🧠 Design Decisions & Trade-offs
## Go Notification & Event Processing System

## 🎯 Purpose
This document explains the key architectural decisions, alternatives considered, and trade-offs made while building this system.

## 🧱 1. Event-Driven Architecture

### ✅ Decision
Use an event-driven architecture instead of synchronous processing.

### 🎯 Why?
- Decouples API from processing.
- Improves scalability.
- Enables asynchronous workflows.

### ⚖️ Trade-offs

| Pros | Cons |
|---|---|
| High scalability | Increased complexity |
| Better performance | Harder debugging |

## 🧵 2. Goroutines for Concurrency

### ✅ Decision
Use goroutines with a worker pool.

### 🎯 Why?
- Lightweight threads.
- Efficient parallel execution.
- Native Go feature.

### ⚖️ Trade-offs

| Pros | Cons |
|---|---|
| Low overhead | Debugging concurrency is harder |
| High performance | Requires careful design |

## 📦 3. Redis as Queue

### ✅ Decision
Use Redis (`LPUSH` / `BRPOP`) as a message queue.

### 🎯 Why?
- Simple setup.
- Fast in-memory operations.
- Supports blocking operations.

### ❌ Alternatives Considered

| Alternative | Redis | Kafka |
|---|---|---|
| Complexity | Easy | Complex setup |
| Durability | Moderate | Highly durable |
| Throughput | Fast | Very high throughput |

👉 Redis was chosen for simplicity and fast development.

### ⚖️ Trade-offs

| Pros | Cons |
|---|---|
| Easy to implement | Less durable than Kafka |
| Fast | Limited advanced features |

## 🔁 4. Retry with Exponential Backoff

### ✅ Decision
Retry failed events using exponential backoff.

### 🎯 Why?
- Handles transient failures.
- Prevents system overload.
- Improves success rate.

### ⚖️ Trade-offs

| Pros | Cons |
|---|---|
| Reliable | Increased latency |
| Smooth load | Slight complexity |

## 📬 5. Dead Letter Queue (DLQ)

### ✅ Decision
Use DLQ for failed events.

### 🎯 Why?
- Prevent data loss.
- Enable debugging.
- Allow reprocessing.

### ⚖️ Trade-offs

| Pros | Cons |
|---|---|
| Reliable failure handling | Requires monitoring |
| Better debugging | Extra storage |

## 🧠 6. Idempotency using SET NX

### ✅ Decision
Use Redis `SET key value NX EX ttl`.

### 🎯 Why?
- Atomic operation.
- Prevent duplicate processing.
- Works in distributed systems.

### ❌ Alternative
`check → process → set`

This is **not safe** because it creates a race condition.

### ⚖️ Trade-offs

| Pros | Cons |
|---|---|
| Precludes duplicates | Needs TTL tuning |
| Atomic | Extra Redis usage |

## 🛑 7. Sliding Window Rate Limiter

### ✅ Decision
Use the Sliding Window algorithm.

### 🎯 Why?
- More accurate than fixed window.
- Prevents burst traffic.

### ❌ Alternatives

| Alternative | Notes |
|---|---|
| Fixed Window | Simpler but inaccurate |
| Token Bucket | Efficient but more complex |

### ⚖️ Trade-offs

| Pros | Cons |
|---|---|
| Accurate | Slightly higher memory |
| Fair | More computation |

## 🧱 8. Stateless API Design

### ✅ Decision
Keep API stateless.

### 🎯 Why?
- Easy horizontal scaling.
- Works with load balancers.
- No session management.

### ⚖️ Trade-offs

| Pros | Cons |
|---|---|
| Scalable | Requires external storage |
| Simple | Slight overhead |

## 🧠 9. In-Worker Retry

### ✅ Decision
Perform retry inside the worker.

### 🎯 Why?
- Faster retry.
- Simpler design.
- Less Redis load.

### ❌ Alternative
Requeue the event in Redis.

### ⚖️ Trade-offs

| Pros | Cons |
|---|---|
| Low latency | Less flexible |
| Simple | Harder tracking |

## 🔒 10. Redis for Idempotency Store

### ✅ Decision
Store processed event IDs in Redis.

### 🎯 Why?
- Fast lookup.
- Shared across instances.
- Supports TTL.

### ⚖️ Trade-offs

| Pros | Cons |
|---|---|
| Distributed | Memory usage |
| Fast | Needs cleanup strategy |

## ⚠️ Known Limitations

- Single Redis instance can become a single point of failure.
  - Solution: Redis Cluster.
- In-memory rate limiter is not distributed.
  - Solution: Redis-based limiter.
- Redis durability is limited compared to log-based brokers.
  - Solution: Enable AOF persistence.

## 🔮 Future Improvements

- Kafka integration.
- Distributed rate limiter.
- Prometheus monitoring.
- OpenTelemetry tracing.
- Circuit breaker pattern.

