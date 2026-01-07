# API Documentation

## Overview

The indexer will expose a REST API for querying indexed blockchain data.

## Endpoints (Planned)

### Health Check

```
GET /health
```

Response:
```json
{
  "status": "ok",
  "current_slot": 12345678,
  "is_running": true
}
```

### Get Block

```
GET /api/v1/blocks/:slot
```

Response:
```json
{
  "slot": 12345678,
  "blockhash": "...",
  "previous_blockhash": "...",
  "parent_slot": 12345677,
  "transactions_count": 150
}
```

### Get Transaction

```
GET /api/v1/transactions/:signature
```

Response:
```json
{
  "signature": "...",
  "slot": 12345678,
  "status": "confirmed",
  "fee": 5000,
  "accounts": [...]
}
```

### Get Current Slot

```
GET /api/v1/slots/current
```

Response:
```json
{
  "current_slot": 12345678,
  "timestamp": "2026-01-07T15:00:00Z"
}
```

### Get Indexer Status

```
GET /api/v1/status
```

Response:
```json
{
  "is_running": true,
  "current_slot": 12345678,
  "start_slot": 0,
  "blocks_processed": 12345678,
  "uptime_seconds": 3600
}
```

## Error Responses

### 404 Not Found
```json
{
  "error": "resource not found",
  "code": "NOT_FOUND"
}
```

### 500 Internal Server Error
```json
{
  "error": "internal server error",
  "code": "INTERNAL_ERROR"
}
```

## Rate Limiting

- 100 requests per minute per IP
- Returns 429 Too Many Requests when exceeded

## Authentication

Currently not implemented. Future versions may include API key authentication.
