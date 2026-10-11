#!/bin/bash

# Railway.app Deployment Script for API Gateway
# This script helps you deploy to Railway with minimal steps

set -e

echo "🚀 API Gateway - Railway Deployment Helper"
echo "==========================================="
echo ""

# Check if railway CLI is installed
if ! command -v railway &> /dev/null; then
    echo "❌ Railway CLI not found. Installing..."
    npm install -g @railway/cli
    echo "✅ Railway CLI installed"
fi

echo ""
echo "📋 RAILWAY DEPLOYMENT STEPS"
echo "==========================="
echo ""
echo "1. LOGIN TO RAILWAY"
echo "   railway login"
echo ""
echo "2. INITIALIZE PROJECT (if not already done)"
echo "   railway init"
echo ""
echo "3. DEPLOY FROM GITHUB"
echo "   a. Push to GitHub first:"
echo "      git add ."
echo "      git commit -m 'Deploy API Gateway to Railway'"
echo "      git push origin main"
echo ""
echo "   b. On Railway Dashboard:"
echo "      → Click 'New Project'"
echo "      → Select 'Deploy from GitHub repo'"
echo "      → Choose your repository"
echo ""
echo "4. ADD DATABASES"
echo "   a. Click '+' → 'New Service' → 'Postgres'"
echo "   b. Click '+' → 'New Service' → 'Redis'"
echo ""
echo "5. SET ENVIRONMENT VARIABLES"
echo "   In your API Gateway service → Variables:"
echo "   • DATABASE_URL=\${postgres.DATABASE_URL}"
echo "   • REDIS_URL=\${redis.URI}"
echo "   • JWT_SECRET=<your-secret-key>"
echo "   • PORT=8080"
echo ""
echo "6. REDEPLOY"
echo "   Railway will auto-redeploy with new variables"
echo ""
echo "🔗 Helpful Links:"
echo "   • Railway Docs: https://docs.railway.app"
echo "   • GitHub Integration: https://docs.railway.app/guides/public-github"
echo ""
echo "💡 Your live URL will be: https://api-gateway.up.railway.app"
echo ""
