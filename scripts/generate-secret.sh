#!/bin/bash

# Generate secure secrets for deployment

echo "🔐 Generate Secure Secrets"
echo "=========================="
echo ""

# Generate JWT Secret
JWT_SECRET=$(openssl rand -hex 32)
echo "JWT_SECRET=$JWT_SECRET"
echo ""

# Generate Database Password
DB_PASSWORD=$(openssl rand -base64 24)
echo "DATABASE_PASSWORD=$DB_PASSWORD"
echo ""

# Generate Redis Password
REDIS_PASSWORD=$(openssl rand -base64 24)
echo "REDIS_PASSWORD=$REDIS_PASSWORD"
echo ""

echo "💾 Save these secrets securely! You'll need them for deployment."
echo ""
echo "Example .env file:"
echo "-------------------"
echo "PORT=8080"
echo "DATABASE_URL=postgres://postgres:${DB_PASSWORD}@localhost:5432/apigateway"
echo "REDIS_URL=redis://:${REDIS_PASSWORD}@localhost:6379"
echo "JWT_SECRET=${JWT_SECRET}"
echo "RATE_LIMIT_RPS=100"
