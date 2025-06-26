#!/bin/bash

# GoChop Fly.io Deployment Script
# This script deploys PostgreSQL, Redis, and the Go backend to Fly.io

set -e

echo "🚀 Starting GoChop deployment to Fly.io..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if flyctl is installed
if ! command -v flyctl &> /dev/null; then
    echo -e "${RED}❌ flyctl is not installed. Please install it first:${NC}"
    echo "   brew install flyctl"
    echo "   # or use the install script from https://fly.io/docs/hands-on/install-flyctl/"
    exit 1
fi

# Check if user is logged in
if ! flyctl auth whoami &> /dev/null; then
    echo -e "${RED}❌ Not logged into Fly.io. Please login first:${NC}"
    echo "   flyctl auth login"
    exit 1
fi

echo -e "${BLUE}📦 Step 1: Creating PostgreSQL database...${NC}"
if flyctl apps list | grep -q "gochop-postgres"; then
    echo -e "${YELLOW}⚠️  PostgreSQL app already exists, skipping creation${NC}"
else
    flyctl postgres create --name gochop-postgres --region ord --vm-size shared-cpu-1x --volume-size 10
fi

echo -e "${BLUE}📦 Step 2: Creating Redis volume...${NC}"
if flyctl volumes list -a gochop-redis | grep -q "redis_data"; then
    echo -e "${YELLOW}⚠️  Redis volume already exists, skipping creation${NC}"
else
    flyctl volumes create redis_data --size 1 --region ord --app gochop-redis
fi

echo -e "${BLUE}📦 Step 3: Deploying Redis...${NC}"
cd redis
if flyctl apps list | grep -q "gochop-redis"; then
    flyctl deploy
else
    flyctl launch --name gochop-redis --no-deploy
    flyctl deploy
fi
cd ..

echo -e "${BLUE}📦 Step 4: Setting up backend secrets...${NC}"
echo -e "${YELLOW}Please set up your secrets:${NC}"
echo "flyctl secrets set DATABASE_URL=postgres://username:password@gochop-postgres.internal:5432/dbname -a gochop-backend"
echo "flyctl secrets set REDIS_URL=redis://gochop-redis.internal:6379 -a gochop-backend"
echo "flyctl secrets set NEXTAUTH_SECRET=your-super-secret-key -a gochop-backend"

echo -e "${BLUE}📦 Step 5: Deploying backend...${NC}"
cd backend
if flyctl apps list | grep -q "gochop-backend"; then
    flyctl deploy
else
    flyctl launch --name gochop-backend --no-deploy
    flyctl deploy
fi
cd ..

echo -e "${GREEN}✅ Deployment complete!${NC}"
echo -e "${BLUE}🔗 Your app should be available at:${NC}"
echo "   https://gochop-backend.fly.dev"
echo ""
echo -e "${YELLOW}💡 Next steps:${NC}"
echo "1. Set up your database secrets (see above)"
echo "2. Update your frontend to point to the new backend URL"
echo "3. Configure your domain (optional)"
echo ""
echo -e "${BLUE}📊 Monitor your apps:${NC}"
echo "   flyctl logs -a gochop-backend"
echo "   flyctl logs -a gochop-redis"
echo "   flyctl logs -a gochop-postgres" 