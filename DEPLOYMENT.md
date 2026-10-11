# Deployment Guide

This guide covers deploying the API Gateway to production environments.

## Quick Deploy - Render.com (Recommended)

Render offers the easiest deployment path with free tier including PostgreSQL and Redis.

### Step 1: Push to GitHub
```bash
git init
git add .
git commit -m "Initial commit: API Gateway"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/api-gateway.git
git push -u origin main
```

### Step 2: Create Render Account
1. Sign up at [render.com](https://render.com)
2. Connect your GitHub account

### Step 3: Create PostgreSQL Database
1. Dashboard → New → PostgreSQL
2. Name: `api-gateway-db`
3. Region: Select closest to you
4. Click "Create Database"
5. Copy the **Internal Database URL** (looks like: `postgres://user:pass@ep-xxx.us-east-1.aws.neon.tech/apigateway`)

### Step 4: Create Redis Instance
1. Dashboard → New → Redis
2. Name: `api-gateway-redis`
3. Region: Same as database
4. Click "Create Redis"
5. Copy the **Internal Connection String** (looks like: `redis://default:pass@red-xxx.us-east-1.aws.neon.tech:6379`)

### Step 5: Deploy Web Service
1. Dashboard → New → Web Service
2. Connect your repository
3. Configure:
   - Name: `api-gateway`
   - Region: Same as database
   - Root Directory: `/`
   - Runtime: Docker
   - Dockerfile: `Dockerfile`
4. Add Environment Variables:
   ```
   PORT=8080
   DATABASE_URL=<From Step 3>
   REDIS_URL=<From Step 4>
   JWT_SECRET=<Generate a random 64-char string>
   RATE_LIMIT_RPS=100
   ```
5. Click "Create Web Service"

### Step 6: Get Your Live URL
Once deployed, Render provides:
```
https://api-gateway.onrender.com
```

---

## Alternative: Railway.app

Railway offers $5 free credits monthly.

### Step 1: Push to GitHub
```bash
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/api-gateway.git
git push -u origin main
```

### Step 2: Deploy on Railway
1. Go to [railway.app](https://railway.app)
2. Click "New Project" → "Deploy from GitHub repo"
3. Select your repository
4. Railway auto-detects Dockerfile and deploys

### Step 3: Add Database & Redis
1. Click "+" → "New Service" → "Postgres"
2. Click "+" → "New Service" → "Redis"
3. Go to API Gateway → "Variables"
4. Add:
   ```
   DATABASE_URL=${postgres.DATABASE_URL}
   REDIS_URL=${redis.URI}
   JWT_SECRET=<your-secret>
   PORT=8080
   ```
5. Redeploy

### Step 4: Get Your URL
Railway provides:
```
https://api-gateway.up.railway.app
```

---

## Alternative: Fly.io

Fly.io offers free allowances for small apps.

### Step 1: Install Fly CLI
```bash
# macOS
brew install fly

# Windows
winget install fly

# Linux
curl -L https://fly.io/install.sh | sh
```

### Step 2: Login
```bash
fly auth login
```

### Step 3: Initialize App
```bash
fly launch
# Follow prompts to create app
```

### Step 4: Add PostgreSQL
```bash
fly postgres attach --app api-gateway-db
```

### Step 5: Add Redis
```bash
fly redis attach --app api-gateway-redis
```

### Step 6: Set Secrets
```bash
fly secrets set DATABASE_URL=postgres://...
fly secrets set REDIS_URL=redis://...
fly secrets set JWT_SECRET=your-secret
```

### Step 7: Deploy
```bash
fly deploy
```

---

## Testing Deployment

Once deployed, test the API:

```bash
# Health check
curl https://YOUR_DOMAIN/health

# Register
curl -X POST https://YOUR_DOMAIN/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# Login
curl -X POST https://YOUR_DOMAIN/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Access protected endpoint
curl https://YOUR_DOMAIN/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN"

# Check metrics
curl https://YOUR_DOMAIN/metrics
```

---

## Environment Variables Reference

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| PORT | Yes | Server port | `8080` |
| DATABASE_URL | Yes | PostgreSQL connection | `postgres://user:pass@host:5432/db` |
| REDIS_URL | Yes | Redis connection | `redis://host:6379/0` |
| JWT_SECRET | Yes | JWT signing key | `your-64-char-secret-here` |
| RATE_LIMIT_RPS | No | Requests per second limit | `100` |

---

## Monitoring

### Prometheus Metrics
Access at: `https://YOUR_DOMAIN/metrics`

### Health Check
Access at: `https://YOUR_DOMAIN/health`

---

## Troubleshooting

### Build fails
Check Dockerfile and ensure all Go dependencies are in go.mod

### Database connection fails
Verify DATABASE_URL format:
```
postgres://user:password@host:5432/database?sslmode=require
```

### Redis connection fails
Verify REDIS_URL format:
```
redis://default:password@host:port
```

### Health check failing
Ensure PORT matches what's exposed in Dockerfile

---

## Cost Estimates

| Platform | Free Tier | Paid (Starter) |
|----------|-----------|----------------|
| Render | 750 hrs/month | $7/month |
| Railway | $5 credits | Usage-based |
| Fly.io | Small VM | $5/month min |

---

## Production Checklist

- [ ] Use strong JWT_SECRET (64+ characters)
- [ ] Enable HTTPS (provided by platform)
- [ ] Set up monitoring alerts
- [ ] Configure backup schedules
- [ ] Set appropriate rate limits
- [ ] Add error tracking (Sentry, etc.)
- [ ] Document operational runbook
- [ ] Test failover scenarios
