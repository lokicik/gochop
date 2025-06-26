# GitHub Actions Setup for Fly.io Deployment

## Quick Setup

1. **Get your Fly.io API token**:

   ```bash
   flyctl auth token
   ```

2. **Add the token to GitHub Secrets**:

   - Go to your repository on GitHub
   - Click Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `FLY_API_TOKEN`
   - Value: [paste your token from step 1]

3. **Initial Setup** (run once):

   ```bash
   # Deploy manually first to create the apps
   chmod +x scripts/deploy.sh
   ./scripts/deploy.sh
   ```

4. **Automatic Deployment**:
   After the initial setup, every push to `main` or `master` will automatically deploy your app!

## Workflow Details

The workflow in `.github/workflows/deploy.yml` will:

1. Deploy Redis first
2. Deploy the backend after Redis is ready
3. Use remote Docker builds (faster and doesn't require local Docker)

## Monitoring

Check your deployments:

- GitHub Actions tab in your repository
- Fly.io dashboard: https://fly.io/dashboard
- Logs: `flyctl logs -a gochop-backend`
