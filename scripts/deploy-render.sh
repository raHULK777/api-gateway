#!/bin/bash

# Render.com Deployment Script for API Gateway
# This script helps you deploy to Render with minimal steps

set -e

echo "🚀 API Gateway - Render Deployment Helper"
echo "=========================================="
echo ""

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo "❌ Git is not installed. Please install git first."
    exit 1
fi

# Check if we're in a git repo
if [ ! -d ".git" ]; then
    echo "❌ Not in a git repository. Run 'git init' first."
    exit 1
fi

echo "✅ Git repository found"
echo ""

# Step 1: Initialize git if needed
if [ ! -f ".gitignore" ]; then
    echo "📝 Creating .gitignore..."
    cat > .gitignore << 'EOF'
.env
.env.local
*.log
node_modules/
EOF
fi

# Step 2: Add all files
echo "📦 Adding files to git..."
git add .

# Step 3: Check for uncommitted changes
if git diff --staged --quiet; then
    echo "⚠️  No changes to commit. Use 'git add <file>' to stage changes."
else
    echo "✅ Changes staged. Commit with:"
    echo "   git commit -m \"Deploy API Gateway to Render\""
fi

echo ""
echo "📋 DEPLOYMENT CHECKLIST"
echo "======================="
echo ""
echo "1. CREATE RENDER ACCOUNT"
echo "   → Go to https://render.com and sign up"
echo ""
echo "2. CREATE POSTGRESQL DATABASE"
echo "   → Dashboard → New → PostgreSQL"
echo "   → Name: api-gateway-db"
echo "   → Copy the Internal Database URL"
echo ""
echo "3. CREATE REDIS INSTANCE"
echo "   → Dashboard → New → Redis"
echo "   → Name: api-gateway-redis"
echo "   → Copy the Connection String"
echo ""
echo "4. DEPLOY WEB SERVICE"
echo "   → Dashboard → New → Web Service"
echo "   → Connect your GitHub repository"
echo "   → Environment Variables:"
echo "     • PORT=8080"
echo "     • DATABASE_URL=<From step 2>"
echo "     • REDIS_URL=<From step 3>"
echo "     • JWT_SECRET=<Generate a random 64-char string>"
echo "     • RATE_LIMIT_RPS=100"
echo ""
echo "5. GET YOUR LIVE URL"
echo "   → Render will provide: https://api-gateway.onrender.com"
echo ""
echo "🔗 Helpful Links:"
echo "   • Render Docs: https://render.com/docs"
echo "   • Docker Deploy: https://render.com/docs/deploy-dockerfile"
echo ""
echo "💡 TIP: Use 'openssl rand -hex 32' to generate a secure JWT_SECRET"
echo ""
echo "Ready to deploy? Follow the steps above!"
