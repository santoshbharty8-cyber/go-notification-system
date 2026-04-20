# 🏗️ High-Level Design (HLD)  
## Go Notification & Event Processing System

## 🎯 Objective

Design a scalable, fault-tolerant, event-driven system that:

- Accepts events via API
- Processes them asynchronously
- Handles failures safely with retry + DLQ
- Prevents duplicate processing with idempotency
- Supports horizontal scaling

## 🧠 Functional Requirements

- Accept events via `/event`
- Support multiple event types:
  - `user_registered`
  - `order_created`
- Process events asynchronously
- Retry failed events
- Avoid duplicate processing
- Store failed events in a DLQ

## ⚙️ Non-Functional Requirements

| Requirement | Target |
|---|---|
| Throughput | High (1000+ events/sec) |
| Latency | Low API response time (<100ms) |
| Scalability | Horizontal scaling |
| Reliability | No event loss |
| Consistency | Idempotent processing |
| Availability | High |

## 🏗️ System Components

### 1. API Layer
**Responsibilities:**
- Accept HTTP requests
- Validate input
- Enqueue events to Redis
- Return immediate response

**Design:**
- Stateless
- Lightweight
- Fast response

### 2. Rate Limiter
**Type:** Sliding Window (DSA-based)

**Purpose:**
- Protect API from abuse
- Control traffic spikes

### 3. Redis Queue (Message Broker)
**Role:**
- Acts as a buffer between API and workers
- Decouples ingestion and processing

**Operations:**
- `LPUSH` → Producer
- `BRPOP` → Consumer

### 4. Worker Pool
**Design:**
- Multiple goroutines
- Each worker processes events independently

**Purpose:**
- Parallel processing
- Efficient CPU usage

### 5. Idempotency Layer
**Problem:**
- Duplicate events may occur due to retries or network issues

**Solution:**
```text
SET key value NX EX ttl
```

**Result:**
- Only one worker processes the event
- Others skip it

### 6. Event Processing Layer
Handles business logic:

- `user_registered` → send email
- `order_created` → process order

### 7. Retry Mechanism
**Strategy:**
- Max retries: 3
- Exponential backoff

**Purpose:**
- Handle transient failures

### 8. Dead Letter Queue (DLQ)
**Purpose:**
- Store permanently failed events

**Benefit:**
- Prevent data loss
- Enable debugging

## 🔄 Data Flow Diagram

```text
Client
  ↓
API Layer
  ↓
Rate Limiter
  ↓
Redis Queue
  ↓
Worker Pool
  ↓
Idempotency Check
  ↓
Processing Logic
  ↓
Retry / DLQ
```

## 🧵 Concurrency Model

- Workers use goroutines
- Shared Redis queue distributes load
- No locking required at app level

## 📈 Scalability Design

### Horizontal Scaling
```text
Multiple API Instances
        ↓
   Load Balancer
        ↓
     Redis Queue
        ↓
   Worker Instances
```

**Benefits:**
- Add more instances easily
- Redis acts as a central broker

### Stateless Design
- No session stored in API
- Easy to scale

## ⚖️ Trade-offs

| Decision | Benefit | Trade-off |
|---|---|---|
| Redis queue | Simple & fast | Not as durable as Kafka |
| In-memory rate limiter | Fast | Not distributed |
| Retry in worker | Simple | Less flexible |
| SET NX idempotency | Atomic | Requires TTL tuning |

## 🔒 Reliability Strategy

| Feature | Purpose |
|---|---|
| Retry | Handle transient errors |
| DLQ | Prevent data loss |
| Idempotency | Avoid duplicates |
| Redis | Durable queue |

## 🧠 Failure Scenarios

### Case 1: Worker crashes
Event remains in Redis and is picked by another worker.

### Case 2: Duplicate events
Idempotency prevents reprocessing.

### Case 3: Processing failure
Retry is attempted, then the event moves to DLQ.

### Case 4: Redis down
System becomes unavailable because Redis is a single point of failure.

**Mitigation:**
- Redis cluster
- Persistence using AOF/RDB

## 🔥 Bottlenecks & Optimizations

| Bottleneck | Solution |
|---|---|
| Redis throughput | Use clustering |
| Worker CPU | Increase workers |
| API overload | Rate limiter |
| Duplicate processing | Idempotency |

## 🧠 Alternative Designs

### Replace Redis with Kafka

| Feature | Redis | Kafka |
|---|---|---|
| Setup | Easy | Complex |
| Throughput | Medium | Very high |
| Durability | Medium | High |
| Use case | Small/medium | Large-scale |

### Distributed Rate Limiter
Replace the in-memory limiter with a Redis-based limiter for better distributed control.

## 🔮 Future Enhancements

- Kafka integration
- Prometheus metrics
- Distributed tracing
- Circuit breaker
- WebSocket notifications

