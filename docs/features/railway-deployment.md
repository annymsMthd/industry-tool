# Railway Deployment Plan for Industry-Tool

## Context

This plan deploys the industry-tool EVE Online application to Railway's cloud platform. The application is a full-stack monorepo with:
- **Backend**: Go 1.25 web server (port 8081) with automatic PostgreSQL migrations
- **Frontend**: Next.js 16 application (port 3000) with OAuth authentication
- **Database**: PostgreSQL 17 with 6 migrations for asset tracking, market prices, and stockpiles

**Deployment Goals:**
- Production environment only (no PR previews)
- Private backend (only accessible via frontend proxy)
- Use Railway's free *.up.railway.app domain
- Leverage existing multi-stage Dockerfile
- Auto-deploy on push to main branch

**Key Challenge:** The `next.config.ts` currently hardcodes `localhost:8081` for backend proxying, which must be made dynamic for Railway's internal networking.

## Service Architecture

Deploy as **three separate Railway services** in a single project:

1. **PostgreSQL Database** (Railway Plugin)
   - Managed PostgreSQL 17 instance
   - Provides automatic connection variables
   - Handles backups and monitoring

2. **Backend Service** (Go)
   - Builds from root `Dockerfile` targeting `final-backend` stage
   - Connects to database via Railway's internal variables
   - Runs migrations automatically on startup
   - Exposes on Railway's internal network only (private)

3. **Frontend Service** (Next.js)
   - Builds from root `Dockerfile` targeting `publish-ui` stage
   - Connects to backend via Railway's internal DNS
   - Proxies `/backend/*` requests to private backend
   - Publicly accessible via Railway's domain

## Required Code Changes

### 1. Fix Frontend Backend URL Configuration

**File**: `/mounts/applications/sources/industry-tool/frontend/next.config.ts`

**Current (line 9):**
```typescript
destination: "http://localhost:8081/:path*",
```

**Change to:**
```typescript
destination: `${process.env.BACKEND_URL || "http://localhost:8081"}/:path*`,
```

This allows the frontend to use Railway's internal backend URL in production while falling back to localhost for local development.

### 2. Create Railway Configuration File

**File**: `/mounts/applications/sources/industry-tool/railway.json` (NEW)

```json
{
  "$schema": "https://railway.app/railway.schema.json",
  "build": {
    "builder": "DOCKERFILE",
    "dockerfilePath": "Dockerfile"
  },
  "deploy": {
    "restartPolicyType": "ON_FAILURE",
    "restartPolicyMaxRetries": 10
  }
}
```

### 3. Create .railwayignore (Optional but Recommended)

**File**: `/mounts/applications/sources/industry-tool/.railwayignore` (NEW)

```
build/
artifacts/
.github/
*.md
.env.local
.env.ci
docker-compose*.yaml
```

## Environment Variables Configuration

### Backend Service

```bash
# Application Configuration
PORT=8081

# Database Connection (using Railway reference variables)
DATABASE_HOST=${{Postgres.PGHOST}}
DATABASE_PORT=${{Postgres.PGPORT}}
DATABASE_USER=${{Postgres.PGUSER}}
DATABASE_PASSWORD=${{Postgres.PGPASSWORD}}
DATABASE_NAME=${{Postgres.PGDATABASE}}

# Security (generate with: openssl rand -hex 32)
BACKEND_KEY=<your-generated-secret-here>

# EVE Online ESI OAuth
OAUTH_CLIENT_ID=<your-eve-client-id>
OAUTH_CLIENT_SECRET=<your-eve-client-secret>

# Optional
LOG_LEVEL=INFO
```

### Frontend Service

```bash
# Backend Communication (Railway's internal DNS)
BACKEND_URL=http://backend.production.railway.internal:8081/

# Shared Secret (reference backend's value)
BACKEND_KEY=${{Backend.BACKEND_KEY}}

# NextAuth Configuration
NEXTAUTH_URL=https://$RAILWAY_PUBLIC_DOMAIN/
NEXTAUTH_SECRET=<generate-with-openssl-rand-base64-32>

# EVE Online OAuth for NextAuth Provider
EVE_CLIENT_ID=<your-eve-client-id>
EVE_CLIENT_SECRET=<your-eve-client-secret>
```

## Pre-Deployment Checklist

### 1. Generate Required Secrets

```bash
# Generate BACKEND_KEY (32-byte hex string)
openssl rand -hex 32

# Generate NEXTAUTH_SECRET (base64 encoded)
openssl rand -base64 32
```

### 2. Register EVE Online OAuth Application

1. Go to https://developers.eveonline.com/
2. Create new application
3. **Scopes needed:**
   - `esi-assets.read_assets.v1`
   - `esi-characters.read_corporation_roles.v1`
   - `esi-corporations.read_divisions.v1`
   - `esi-universe.read_structures.v1`
4. **Callback URL**: `https://<your-railway-url>/api/auth/callback/eveonline` (update after Railway generates URL)
5. Save the `Client ID` and `Client Secret`

## Step-by-Step Deployment Guide

### Phase 1: Create Railway Project and Database

1. **Create Railway Account** at https://railway.app/ and sign up with GitHub
2. **Create New Project**: Dashboard → "New Project" → "Deploy from GitHub repo"
3. **Add PostgreSQL Database**: Click "+ New" → "Database" → "PostgreSQL" (service name: `Postgres`)

### Phase 2: Deploy Backend Service

1. **Add Backend Service**: "+ New" → "GitHub Repo" → Select `industry-tool` (service name: `Backend`)
2. **Configure Build**: Settings → Build → Builder: Docker, Dockerfile Path: `Dockerfile`, Docker Build Args: `--target final-backend`
3. **Set Environment Variables**: Add variables from "Backend Service" section above
4. **Deploy**: Railway auto-deploys, monitor logs for "starting services" (success indicator)

### Phase 3: Deploy Frontend Service

1. **Add Frontend Service**: "+ New" → "GitHub Repo" → Select `industry-tool` (service name: `Frontend`)
2. **Configure Build**: Settings → Build → Builder: Docker, Dockerfile Path: `Dockerfile`, Docker Build Args: `--target publish-ui`
3. **Set Environment Variables (Part 1)**: Add all variables except `NEXTAUTH_URL`
4. **Generate Public Domain**: Settings → Networking → "Generate Domain" (copy the generated URL)
5. **Complete Environment Variables**: Add `NEXTAUTH_URL=https://<railway-generated-url>/`
6. **Deploy**: Railway auto-redeploys with new variables

### Phase 4: Update OAuth Callback URL

1. Get frontend public URL from Railway
2. Update EVE Online application callback to: `https://<railway-url>/api/auth/callback/eveonline`

### Phase 5: Verify Deployment

1. **Check Backend Logs**: Look for "starting services", "services started"
2. **Check Frontend Logs**: Look for "ready - started server on 0.0.0.0:3000"
3. **Test Application**:
   - Visit Railway-generated frontend URL
   - Click "Sign In with EVE Online"
   - Complete OAuth flow
   - Verify dashboard loads

## Networking Architecture

```
Internet (HTTPS)
    ↓
Railway Load Balancer
    ↓
Frontend Service (public)
  - Domain: *.up.railway.app
  - Port: 3000
  - Proxies /backend/* to internal backend
    ↓
Backend Service (private)
  - Internal DNS: backend.production.railway.internal:8081
  - Not publicly accessible
    ↓
PostgreSQL Service (private)
  - Internal DNS only
  - Only accessible by backend
```

## Common Issues and Solutions

### Issue 1: Backend Cannot Connect to Database

**Solutions:**
- Verify database service name is `Postgres` (case-sensitive)
- Check DATABASE_* variables use `${{Postgres.PGHOST}}` syntax
- Ensure backend has network access to database

### Issue 2: Frontend Gets 401 from Backend

**Solutions:**
- Verify `BACKEND_KEY` matches exactly in both services
- Use variable reference: `BACKEND_KEY=${{Backend.BACKEND_KEY}}`
- Check header name is `BACKEND-KEY` (with hyphen)

### Issue 3: OAuth Callback Fails

**Solutions:**
- Verify `NEXTAUTH_URL` matches Railway's public domain exactly (including https:// and trailing /)
- Update EVE Online app callback URL to match
- Ensure `NEXTAUTH_SECRET` is set

### Issue 4: Frontend Cannot Reach Backend

**Solutions:**
- Verify `BACKEND_URL` uses internal DNS: `http://backend.production.railway.internal:8081/`
- Check backend service name is exactly `Backend`
- Ensure backend listens on `0.0.0.0:8081`, not `127.0.0.1:8081`

### Issue 5: Docker Build Fails

**Solutions:**
- Verify Docker build target: Backend: `--target final-backend`, Frontend: `--target publish-ui`
- Ensure build args are set in Railway UI under Settings → Build

## Critical Files

**Files to modify:**
- `/mounts/applications/sources/industry-tool/frontend/next.config.ts` - Make backend URL dynamic
- `/mounts/applications/sources/industry-tool/railway.json` - NEW: Railway configuration
- `/mounts/applications/sources/industry-tool/.railwayignore` - NEW: Optimize build

**Files to reference (no changes needed):**
- `/mounts/applications/sources/industry-tool/Dockerfile` - Multi-stage build configuration
- `/mounts/applications/sources/industry-tool/cmd/industry-tool/cmd/settings.go` - Backend environment loading
- `/mounts/applications/sources/industry-tool/internal/database/postgres.go` - Database migrations
- `/mounts/applications/sources/industry-tool/frontend/pages/api/auth/[...nextauth].ts` - OAuth configuration

## Verification Checklist

- [ ] Backend service running (green status in Railway)
- [ ] Frontend service running (green status in Railway)
- [ ] PostgreSQL service running
- [ ] Backend logs show successful database connection
- [ ] Backend logs show migrations completed
- [ ] Frontend logs show Next.js started successfully
- [ ] Can access frontend URL in browser
- [ ] OAuth flow works (sign in with EVE Online)
- [ ] Can see authenticated dashboard
- [ ] No errors in browser console or Railway logs

## Next Steps After Deployment

### Custom Domain (Optional)
- Frontend → Settings → Networking → "Custom Domain"
- Add CNAME in DNS provider
- Update `NEXTAUTH_URL` and EVE OAuth callback

### Monitoring
- Monitor Railway dashboard for CPU/memory usage
- Set up alerts for service failures
- Review application logs regularly

### Database Backups
- Railway includes automatic daily backups
- 7-day retention on free tier
- Upgrade to Pro for longer retention

### Scaling
- **Vertical Scaling**: Settings → Resources → Increase CPU/Memory
- **Horizontal Scaling**: Pro plan supports replicas with automatic load balancing

---

**Time Estimate:** ~1 hour for first deployment
**Post-Deployment:** Automatic deployments on every push to main
**Last Updated:** 2026-02-14
