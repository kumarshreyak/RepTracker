# GymLog Backend Deployment Guide

This guide covers deploying the GymLog backend to Google Cloud Platform using Cloud Run.

## Prerequisites

1. **Google Cloud Platform Account**
   - Create a GCP project
   - Enable billing
   - Enable the following APIs:
     - Cloud Run API
     - Cloud Build API
     - Container Registry API
     - Cloud SQL Admin API (if using Cloud SQL)

2. **Local Development Tools**
   - [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
   - [Docker](https://docs.docker.com/get-docker/)
   - Go 1.23+

3. **Authentication**
   ```bash
   gcloud auth login
   gcloud config set project YOUR_PROJECT_ID
   gcloud auth configure-docker
   ```

## Environment Setup

### 1. Local Development

1. **Copy environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Update `.env` with your values:**
   ```env
   MONGODB_URI=mongodb://localhost:27017
   DB_NAME=gymlog
   JWT_SECRET=your-local-jwt-secret
   GOOGLE_CLIENT_ID=your-google-client-id
   GOOGLE_CLIENT_SECRET=your-google-client-secret
   ```

3. **Run with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

4. **Or run locally:**
   ```bash
   # Start MongoDB
   mongod
   
   # Run the server
   go run cmd/server/main.go
   ```

### 2. Production Setup

1. **MongoDB Atlas (Recommended)**
   - Create a MongoDB Atlas cluster
   - Get the connection string
   - Add your Cloud Run IP ranges to network access

2. **Google Cloud SQL for MongoDB (Alternative)**
   - Create a Cloud SQL instance
   - Configure networking

## Deployment Options

### Option 1: Manual Deployment (Recommended for first-time setup)

1. **Build and deploy:**
   ```bash
   ./deploy.sh production YOUR_GCP_PROJECT_ID
   ```

2. **Secrets are automatically configured via Secret Manager:**
   The deployment script automatically uses secrets from Google Secret Manager:
   - `MONGODB_URI` from `mongodb-uri` secret
   - `DB_NAME` from `db-name` secret
   - `JWT_SECRET` from `jwt-secret` secret
   
   No manual environment variable setup needed!

### Option 2: Cloud Build (CI/CD)

1. **Set up Cloud Build trigger:**
   ```bash
   gcloud builds triggers create github \
     --repo-name=YOUR_GITHUB_REPO \
     --repo-owner=YOUR_GITHUB_USERNAME \
     --branch-pattern="^main$" \
     --build-config=backend/cloudbuild.yaml
   ```

2. **Set up Cloud Build substitutions for environment variables**

### Option 3: Manual gcloud Commands

1. **Build the image:**
   ```bash
   cd backend
   docker build -t gcr.io/YOUR_PROJECT_ID/gymlog-backend .
   docker push gcr.io/YOUR_PROJECT_ID/gymlog-backend
   ```

2. **Deploy to Cloud Run:**
   ```bash
   gcloud run deploy gymlog-backend \
     --image gcr.io/YOUR_PROJECT_ID/gymlog-backend \
     --platform managed \
     --region us-central1 \
     --allow-unauthenticated \
     --memory 512Mi \
     --cpu 1 \
     --max-instances 10
   ```

## Secret Manager Setup

### Prerequisites

1. **Enable Secret Manager API:**
   ```bash
   gcloud services enable secretmanager.googleapis.com
   ```

2. **Grant Cloud Build access to secrets:**
   ```bash
   gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
     --member="serviceAccount:YOUR_PROJECT_NUMBER@cloudbuild.gserviceaccount.com" \
     --role="roles/secretmanager.secretAccessor"
   ```

### Creating Secrets

1. **MongoDB Atlas URI:**
   ```bash
   echo "mongodb+srv://username:password@cluster.mongodb.net/" | \
   gcloud secrets create mongodb-uri --data-file=-
   ```

2. **Database Name:**
   ```bash
   echo "gymlog" | gcloud secrets create db-name --data-file=-
   ```

3. **JWT Secret (auto-generated):**
   ```bash
   openssl rand -base64 32 | gcloud secrets create jwt-secret --data-file=-
   ```

### Updating Secrets

To update any secret with a new value:
```bash
echo "new-value" | gcloud secrets versions add SECRET_NAME --data-file=-
```

### Verifying Secrets

List all secrets:
```bash
gcloud secrets list
```

View secret metadata (not the actual value):
```bash
gcloud secrets describe SECRET_NAME
```

## Environment Variables Configuration

### Required Variables (via Secret Manager)

| Variable | Description | Secret Name | Example |
|----------|-------------|-------------|---------|
| `MONGODB_URI` | MongoDB connection string | `mongodb-uri` | `mongodb+srv://user:pass@cluster.net/` |
| `DB_NAME` | Database name | `db-name` | `gymlog` |
| `JWT_SECRET` | JWT signing secret | `jwt-secret` | `auto-generated-32-char-string` |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | `google-client-id` | `xxx.apps.googleusercontent.com` |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret | `google-client-secret` | `GOCSPX-xxx` |

**Note:** These variables are automatically injected from Google Secret Manager during deployment.

### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ENV` | Environment | `production` |
| `LOG_LEVEL` | Logging level | `info` |

## Security Configuration

### 1. Cloud Run Security

```bash
# Update service with security settings
gcloud run services update gymlog-backend \
  --cpu-throttling \
  --ingress=all \
  --max-instances=10 \
  --concurrency=80 \
  --timeout=300 \
  --region us-central1
```

### 2. IAM Permissions

Create a service account for Cloud Run:

```bash
# Create service account
gcloud iam service-accounts create gymlog-backend-sa \
  --display-name="GymLog Backend Service Account"

# Assign minimal permissions
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:gymlog-backend-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"
```

### 3. Network Security

For MongoDB Atlas:
- Add Cloud Run IP ranges to Atlas network access
- Use connection string with authentication

For Cloud SQL:
- Enable private IP
- Configure VPC connector if needed

## Monitoring and Logging

### 1. Cloud Logging

View logs:
```bash
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=gymlog-backend" --limit=50
```

### 2. Cloud Monitoring

Set up alerts for:
- High error rates
- Memory usage
- Response time
- Instance count

### 3. Health Checks

Add a health check endpoint to your application and configure:

```bash
gcloud run services update gymlog-backend \
  --set-env-vars HEALTH_CHECK_PATH="/health" \
  --region us-central1
```

## Scaling Configuration

### Auto-scaling settings:

```bash
gcloud run services update gymlog-backend \
  --min-instances=1 \
  --max-instances=10 \
  --concurrency=80 \
  --cpu=1 \
  --memory=512Mi \
  --region us-central1
```

## Domain Configuration

### 1. Custom Domain

```bash
# Map custom domain
gcloud run domain-mappings create \
  --service=gymlog-backend \
  --domain=api.yourdomain.com \
  --region=us-central1
```

### 2. SSL Certificate

Cloud Run automatically provisions SSL certificates for custom domains.

## Troubleshooting

### Common Issues:

1. **MongoDB Connection Issues**
   - Check MongoDB URI format
   - Verify network access rules
   - Check authentication credentials

2. **Environment Variables**
   - Verify all required variables are set
   - Check variable names (case-sensitive)
   - Restart service after changes

3. **Build Issues**
   - Check Dockerfile syntax
   - Verify go.mod dependencies
   - Check build logs in Cloud Build

4. **Memory/CPU Issues**
   - Increase memory allocation
   - Add CPU limits
   - Check application logs

### Debug Commands:

```bash
# Check service status
gcloud run services describe gymlog-backend --region=us-central1

# View recent logs
gcloud logging read "resource.type=cloud_run_revision" --limit=20

# Check environment variables
gcloud run services describe gymlog-backend --region=us-central1 --format="value(spec.template.spec.template.spec.containers[0].env[].name,spec.template.spec.template.spec.containers[0].env[].value)"
```

## Cost Optimization

1. **Use minimum instances** during development
2. **Set CPU throttling** when not in use
3. **Monitor usage** with Cloud Monitoring
4. **Use preemptible instances** for non-critical workloads

## Backup Strategy

1. **MongoDB Backups**
   - Atlas automated backups
   - Manual exports for critical data

2. **Code Backups**
   - Git repository
   - Container images in GCR

## Next Steps

After deployment:

1. Update frontend API endpoints
2. Configure CORS for your domain
3. Set up monitoring alerts
4. Configure backup strategies
5. Set up staging environment
6. Configure CI/CD pipeline 