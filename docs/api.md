# Scorpes API Documentation

## Base URL

```
http://localhost:{PORT}
```

## Response Format

All API responses follow this structure:

```json
{
  "success": true,
  "data": { ... },
  "error": "error message (only on failure)"
}
```

---

## Health Check

### GET /health

Check API server health status.

**Response**

```json
{
  "success": true,
  "data": "OK"
}
```

---

## Targets API

### Target Object

| Field | Type | Description |
|-------|------|-------------|
| `id` | `uuid` | Unique identifier |
| `name` | `string` | Target name |
| `url` | `string` | Target URL to monitor |
| `method` | `string` | HTTP method (GET, POST, PUT, DELETE, HEAD, PATCH, OPTIONS) |
| `interval_seconds` | `integer` | Check interval in seconds (min: 60) |
| `timeout_seconds` | `integer` | Request timeout in seconds (default: 10) |
| `is_active` | `boolean` | Whether the target is active |
| `created_at` | `timestamp` | Creation timestamp |
| `deleted_at` | `timestamp` | Deletion timestamp (null if not deleted) |

---

### GET /api/targets

List all active targets.

**Response**

```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Example API",
      "url": "https://api.example.com/health",
      "method": "GET",
      "interval_seconds": 60,
      "timeout_seconds": 10,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "deleted_at": null
    }
  ]
}
```

---

### POST /api/targets

Create a new monitoring target.

**Request Body**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | `string` | Yes | Target name |
| `url` | `string` | Yes | Valid URL to monitor |
| `method` | `string` | No | HTTP method (default: GET) |
| `interval_seconds` | `integer` | Yes | Check interval (min: 60) |
| `timeout_seconds` | `integer` | No | Request timeout (default: 10) |

**Example Request**

```json
{
  "name": "Production API",
  "url": "https://api.example.com/health",
  "method": "GET",
  "interval_seconds": 60,
  "timeout_seconds": 10
}
```

**Response (201 Created)**

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Production API",
    "url": "https://api.example.com/health",
    "method": "GET",
    "interval_seconds": 60,
    "timeout_seconds": 10,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "deleted_at": null
  }
}
```

**Error Responses**

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `name is required` | Name field is empty |
| 400 | `url is required` | URL field is empty |
| 400 | `invalid URL format` | URL is not valid |
| 400 | `invalid HTTP method` | Method is not supported |
| 400 | `interval seconds must be greater than or equal to 60` | Interval is less than 60 |
| 500 | `Failed to create target` | Database error |

---

### PUT /api/targets/{id}

Update an existing target.

**Path Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `uuid` | Target ID |

**Request Body**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | `string` | Yes | Target name |
| `url` | `string` | Yes | Valid URL to monitor |
| `method` | `string` | No | HTTP method (default: GET) |
| `interval_seconds` | `integer` | Yes | Check interval (min: 60) |
| `timeout_seconds` | `integer` | No | Request timeout (default: 10) |

**Example Request**

```json
{
  "name": "Updated API",
  "url": "https://api.example.com/v2/health",
  "method": "POST",
  "interval_seconds": 120,
  "timeout_seconds": 15
}
```

**Response (200 OK)**

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Updated API",
    "url": "https://api.example.com/v2/health",
    "method": "POST",
    "interval_seconds": 120,
    "timeout_seconds": 15,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "deleted_at": null
  }
}
```

**Error Responses**

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `Target ID is required` | ID path parameter is missing |
| 400 | `Invalid target ID format` | ID is not a valid UUID |
| 400 | `Invalid request body` | JSON parsing failed |
| 400 | `name is required` | Name field is empty |
| 400 | `url is required` | URL field is empty |
| 400 | `invalid URL format` | URL is not valid |
| 400 | `invalid HTTP method` | Method is not supported |
| 400 | `interval seconds must be greater than or equal to 60` | Interval is less than 60 |
| 500 | `Failed to update target` | Database error |

---

### DELETE /api/targets/{id}

Soft delete a target.

**Path Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `uuid` | Target ID |

**Response (200 OK)**

```json
{
  "success": true,
  "data": "Target deleted successfully"
}
```

**Error Responses**

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `Target ID is required` | ID path parameter is missing |
| 400 | `Invalid target ID format` | ID is not a valid UUID |
| 404 | `Target not found` | Target does not exist or already deleted |
| 500 | `Failed to delete target` | Database error |

---

## Status API

### GET /api/status

Get uptime status summary.

**Response**

```json
{
  "success": true,
  "data": "OK"
}
```
