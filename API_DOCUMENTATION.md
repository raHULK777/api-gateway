# API Documentation

Base URL: `https://api-gateway.onrender.com/api/v1`

## Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

Tokens expire after 24 hours.

---

## Auth Endpoints

### Register User

Register a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Response (201 Created):**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors:**
- `400 Bad Request` - Invalid input
- `409 Conflict` - Email already registered
- `500 Internal Server Error` - Server error

---

### Login

Authenticate and receive a JWT token.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors:**
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Invalid credentials

---

## User Endpoints (Protected)

### Get Current User

Retrieve the authenticated user's profile.

**Endpoint:** `GET /users/me`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe"
}
```

---

### Get User by ID

Retrieve a user by their ID.

**Endpoint:** `GET /users/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id` (integer) - User ID

**Response (200 OK):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe"
}
```

**Errors:**
- `404 Not Found` - User not found

---

## Item Endpoints (Protected)

### List Items

Get all items for the authenticated user.

**Endpoint:** `GET /items`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "items": [
    {
      "id": 1,
      "name": "My Item",
      "description": "Item description",
      "user_id": 1,
      "created_at": "2026-10-11T10:00:00Z",
      "updated_at": "2026-10-11T10:00:00Z"
    }
  ]
}
```

---

### Create Item

Create a new item.

**Endpoint:** `POST /items`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "My Item",
  "description": "Item description"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "My Item",
  "description": "Item description",
  "user_id": 1,
  "created_at": "2026-10-11T10:00:00Z",
  "updated_at": "2026-10-11T10:00:00Z"
}
```

**Errors:**
- `400 Bad Request` - Invalid input

---

### Get Item

Get a specific item by ID.

**Endpoint:** `GET /items/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id` (integer) - Item ID

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "My Item",
  "description": "Item description",
  "user_id": 1,
  "created_at": "2026-10-11T10:00:00Z",
  "updated_at": "2026-10-11T10:00:00Z"
}
```

**Errors:**
- `404 Not Found` - Item not found

---

### Update Item

Update an existing item.

**Endpoint:** `PUT /items/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Path Parameters:**
- `id` (integer) - Item ID

**Request Body:**
```json
{
  "name": "Updated Item",
  "description": "Updated description"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Updated Item",
  "description": "Updated description",
  "user_id": 1,
  "created_at": "2026-10-11T10:00:00Z",
  "updated_at": "2026-10-11T11:00:00Z"
}
```

**Errors:**
- `404 Not Found` - Item not found or unauthorized
- `400 Bad Request` - Invalid input

---

### Delete Item

Delete an item.

**Endpoint:** `DELETE /items/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id` (integer) - Item ID

**Response (200 OK):**
```json
{
  "message": "item deleted"
}
```

**Errors:**
- `404 Not Found` - Item not found or unauthorized

---

## System Endpoints

### Health Check

Check if the API is running.

**Endpoint:** `GET /health`

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": 1728638400
}
```

---

### Prometheus Metrics

Get Prometheus-compatible metrics.

**Endpoint:** `GET /metrics`

**Response:** Prometheus text format with metrics like:
- `api_gateway_requests_total` - Total requests by method, path, status
- `api_gateway_request_duration_seconds` - Request duration histogram

---

## Error Responses

All errors follow this format:
```json
{
  "error": "Error message"
}
```

### Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing or invalid token |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Resource already exists |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |

---

## Rate Limiting

- Default: 100 requests per second per IP
- Exceeded requests return `429 Too Many Requests`
- Configurable via `RATE_LIMIT_RPS` environment variable

---

## Examples

### cURL Examples

```bash
# Register
curl -X POST https://api.example.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","name":"John Doe"}'

# Login
curl -X POST https://api.example.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Get profile
curl https://api.example.com/api/v1/users/me \
  -H "Authorization: Bearer TOKEN"

# Create item
curl -X POST https://api.example.com/api/v1/items \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Item","description":"Description"}'

# List items
curl https://api.example.com/api/v1/items \
  -H "Authorization: Bearer TOKEN"
```

### JavaScript (Fetch API)

```javascript
// Register
const register = async () => {
  const response = await fetch('https://api.example.com/api/v1/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: 'user@example.com',
      password: 'password123',
      name: 'John Doe'
    })
  });
  return await response.json();
};

// Login
const login = async () => {
  const response = await fetch('https://api.example.com/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: 'user@example.com',
      password: 'password123'
    })
  });
  return await response.json();
};

// Protected endpoint
const getProfile = async (token) => {
  const response = await fetch('https://api.example.com/api/v1/users/me', {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  return await response.json();
};
```

### Python (requests)

```python
import requests

BASE_URL = "https://api.example.com/api/v1"

# Register
response = requests.post(f"{BASE_URL}/auth/register", json={
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
})
token = response.json()["token"]

# Login
response = requests.post(f"{BASE_URL}/auth/login", json={
    "email": "user@example.com",
    "password": "password123"
})
token = response.json()["token"]

# Protected endpoint
headers = {"Authorization": f"Bearer {token}"}
response = requests.get(f"{BASE_URL}/users/me", headers=headers)
```
