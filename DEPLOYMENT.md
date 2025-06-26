# 🚀 GoChop Fly.io Deployment Guide

This guide will help you deploy your GoChop backend, Redis, and PostgreSQL to Fly.io with GitHub Actions for automatic deployment.

## 📋 Prerequisites

1. **Install Fly CLI**:

   ```bash
   # macOS
   brew install flyctl

   # Linux/Windows
   curl -L https://fly.io/install.sh | sh
   ```

2. **Login to Fly.io**:

   ```bash
   flyctl auth login
   ```

3. **Set up GitHub Secrets** (for automatic deployment):
   - Go to your GitHub repository
   - Navigate to Settings → Secrets and variables → Actions
   - Add a new secret: `FLY_API_TOKEN`
   - Get your token by running: `flyctl auth token`

## 🛠️ Manual Deployment

### Option 1: Use the Deployment Script (Recommended)

```bash
# Make the script executable
chmod +x scripts/deploy.sh

# Run the deployment script
./scripts/deploy.sh
```

### Option 2: Step-by-Step Manual Deployment

#### 1. Deploy PostgreSQL

```bash
flyctl postgres create --name gochop-postgres --region ord --vm-size shared-cpu-1x --volume-size 10
```

Note down the connection details provided.

#### 2. Deploy Redis

```bash
# Create volume for Redis data persistence
flyctl volumes create redis_data --size 1 --region ord --app gochop-redis

# Deploy Redis
cd redis
flyctl launch --name gochop-redis --no-deploy
flyctl deploy
cd ..
```

#### 3. Deploy Backend

```bash
cd backend

# Launch the app (creates fly.toml if it doesn't exist)
flyctl launch --name gochop-backend --no-deploy

# Set environment secrets
flyctl secrets set DATABASE_URL="postgres://username:password@gochop-postgres.internal:5432/dbname"
flyctl secrets set REDIS_URL="redis://gochop-redis.internal:6379"
flyctl secrets set NEXTAUTH_SECRET="your-super-secret-key-change-this"

# Deploy
flyctl deploy

cd ..
```

## 🔄 Automatic Deployment with GitHub Actions

Once you've set up the `FLY_API_TOKEN` secret in GitHub, every push to the `main` or `master` branch will automatically deploy your application.

The workflow will:

1. Deploy Redis first
2. Deploy the backend after Redis is ready

## 🔧 Environment Variables

Your backend will have access to these environment variables:

- `DATABASE_URL`: Connection string to your PostgreSQL database
- `REDIS_URL`: Connection string to your Redis instance
- `BASE_URL`: Your app's public URL (automatically set)
- `NEXTAUTH_SECRET`: Secret for NextAuth.js

## 📊 Monitoring and Debugging

### View Logs

```bash
# Backend logs
flyctl logs -a gochop-backend

# Redis logs
flyctl logs -a gochop-redis

# PostgreSQL logs
flyctl logs -a gochop-postgres
```

### SSH into Containers

```bash
# SSH into backend
flyctl ssh console -a gochop-backend

# SSH into Redis
flyctl ssh console -a gochop-redis
```

### Health Checks

Your backend has a health check endpoint at `/api/health` that monitors:

- Database connectivity
- Redis connectivity
- Overall service status

## 🌐 Service Communication

All services communicate over Fly.io's private network using `.internal` hostnames:

- **Backend**: `gochop-backend.internal:3001`
- **Redis**: `gochop-redis.internal:6379`
- **PostgreSQL**: `gochop-postgres.internal:5432`

## 🔒 Security Notes

1. **Database Initialization**: Your `init.sql` script runs automatically when PostgreSQL starts
2. **Environment Secrets**: Use `flyctl secrets` instead of environment variables for sensitive data
3. **Network Security**: All internal communication uses Fly.io's private network
4. **HTTPS**: All external traffic is automatically encrypted with Fly.io's TLS termination

## 🎯 Production Checklist

- [ ] Update `NEXTAUTH_SECRET` to a secure random value
- [ ] Configure proper CORS origins for production
- [ ] Set up monitoring and alerts
- [ ] Configure custom domain (optional)
- [ ] Enable auto-scaling if needed
- [ ] Set up backup strategy for PostgreSQL
- [ ] Review and adjust resource allocations

## 🐛 Troubleshooting

### App Won't Start

```bash
# Check logs for errors
flyctl logs -a gochop-backend

# Check app status
flyctl status -a gochop-backend
```

### Database Connection Issues

```bash
# Test database connectivity
flyctl ssh console -a gochop-backend
# Inside the container:
psql $DATABASE_URL
```

### Redis Connection Issues

```bash
# Test Redis connectivity
flyctl ssh console -a gochop-backend
# Inside the container:
redis-cli -h gochop-redis.internal ping
```

## 📚 Additional Resources

- [Fly.io Documentation](https://fly.io/docs/)
- [Fly.io PostgreSQL Guide](https://fly.io/docs/postgres/)
- [Fly.io Networking](https://fly.io/docs/networking/)
- [GitHub Actions with Fly.io](https://fly.io/docs/app-guides/continuous-deployment-with-github-actions/)
