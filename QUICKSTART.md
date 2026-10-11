# Quickstart Guide

Get the API Gateway running in 5 minutes.

## Prerequisites

- Docker Desktop installed
- Git installed

## Option 1: Local Development (Docker)

```bash
# Navigate to project
cd D:\Remote-job\api-gateway

# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f api-gateway
```

### Test the API

```bash
# Health check
curl http://localhost:8080/health

# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Get profile (replace TOKEN with the token from login)
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer TOKEN"

# Create item
curl -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"My Item","description":"Test item"}'

# List items
curl http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer TOKEN"

# View metrics
curl http://localhost:8080/metrics
```

### Access Dashboards

- **Prometheus:** http://localhost:9090
- **Grafana:** http://localhost:3000 (admin/admin)

### Stop Services

```bash
docker-compose down
```

## Option 2: Run Locally (Go)

```bash
# Install dependencies
go mod tidy

# Start PostgreSQL and Redis manually
# Or use Docker for just these services
docker-compose up -d db redis

# Run the application
go run cmd/server/main.go
```

## Option 3: Deploy to Cloud

See [DEPLOYMENT.md](./DEPLOYMENT.md) for detailed cloud deployment guides.

## API Endpoints Summary

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/health` | No | Health check |
| POST | `/api/v1/auth/register` | No | Register user |
| POST | `/api/v1/auth/login` | No | Login |
| GET | `/api/v1/users/me` | Yes | Get current user |
| GET | `/api/v1/users/:id` | Yes | Get user by ID |
| GET | `/api/v1/items` | Yes | List items |
| POST | `/api/v1/items` | Yes | Create item |
| GET | `/api/v1/items/:id` | Yes | Get item |
| PUT | `/api/v1/items/:id` | Yes | Update item |
| DELETE | `/api/v1/items/:id` | Yes | Delete item |
| GET | `/metrics` | No | Prometheus metrics |

## Troubleshooting

### Connection refused
Make sure PostgreSQL and Redis are running:
```bash
docker-compose ps
```

### Database migration errors
The application uses GORM auto-migration. Tables are created automatically on first run.

### Rate limiting
Default is 100 requests/second per IP. Adjust with `RATE_LIMIT_RPS` environment variable.

## Next Steps

1. Deploy to Render/Railway (see DEPLOYMENT.md)
2. Add load testing (see tests/load_test.js)
3. Customize for your use case
4. Add to your portfolio
