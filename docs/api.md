# 🌐 API Documentation
## Go Notification & Event Processing System

## 🎯 Overview

This document describes the HTTP APIs exposed by the system.

**Base URL:** `http://localhost:8080`

## 📌 API Summary

| Method | Endpoint | Description |
|---|---|---|
| POST | `/event` | Submit a new event |
| GET | `/health` | Health check |

## 🧾 1. POST /event

### 📌 Description
Accepts an event and queues it for asynchronous processing.

### 🧾 Request Format

**Headers**
```http
Content-Type: application/json
```

**Body**
```json
{
  "id": "101",
  "type": "order_created",
  "timestamp": 1710000000,
  "payload": {
    "order_id": "ORD101"
  }
}
```

### 🧠 Field Explanation

| Field | Type | Required | Description |
|---|---|---:|---|
| id | string | ✅ | Unique event identifier |
| type | string | ✅ | Event type |
| timestamp | int64 | ✅ | Unix timestamp |
| payload | object | ✅ | Event-specific data |

### 🎯 Supported Event Types

| Event Type | Description |
|---|---|
| user_registered | Trigger welcome flow |
| order_created | Trigger order processing |

### 🔄 Request Flow

```text
Client Request
 → Decode JSON
 → Validate Event
 → Rate Limit Check
 → Push to Redis Queue
 → Return Response
```

### ✅ Success Response

```json
{
  "status": "success",
  "message": "event queued",
  "data": {
    "event_id": "101"
  }
}
```

### ❌ Error Responses

#### 1. Invalid JSON
```json
{
  "status": "error",
  "message": "invalid request body"
}
```

#### 2. Validation Error
```json
{
  "status": "error",
  "message": "event id is required"
}
```

#### 3. Invalid Event Type
```json
{
  "status": "error",
  "message": "invalid event type"
}
```

#### 4. Rate Limit Exceeded
```json
{
  "status": "error",
  "message": "rate limit exceeded"
}
```

#### 5. Queue Failure (Redis issue)
```json
{
  "status": "error",
  "message": "failed to enqueue event"
}
```

## ⚠️ Important Notes

### 🔁 Asynchronous Processing
- API only queues the event.
- Processing happens in the background.
- The response does **not** guarantee successful processing.

### 🧠 Idempotency
- Duplicate events with the same `id` are ignored.
- Only the first event is processed.

### 🛑 Rate Limiting
- Limits requests per client.
- Uses the Sliding Window algorithm.

## 🧪 Example Usage

### ▶️ cURL
```bash
curl -X POST http://localhost:8080/event \
-H "Content-Type: application/json" \
-d '{
  "id": "200",
  "type": "user_registered",
  "timestamp": 1710000000,
  "payload": {
    "email": "test@example.com"
  }
}'
```

### ▶️ Postman
- **Method:** POST
- **URL:** `http://localhost:8080/event`
- **Body:** JSON

## 🏥 2. GET /health

### 📌 Description
Checks if the service is running.

### ▶️ Request
```bash
curl http://localhost:8080/health
```

### ✅ Response
```text
OK
```

## 🔐 Security Considerations (Future)

- Add authentication (JWT/API key).
- Validate payload schema.
- Input sanitization.

## 📊 Observability (Future)

- Request logging.
- Metrics for success/failure rate.
- Tracing.

## ⚙️ Rate Limiting Details

| Parameter | Value |
|---|---|
| Limit | 5 requests |
| Window | 10 seconds |

## 🚀 API Guarantees

| Feature | Guarantee |
|---|---|
| Availability | High |
| Processing | Asynchronous |
| Duplicate Handling | Idempotent |
| Failure Handling | Retry + DLQ |

