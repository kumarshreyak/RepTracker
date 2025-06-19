#!/bin/bash

# GymLog Backend Deployment Script for Google Cloud Run
# Usage: ./deploy.sh [environment] [project-id]

set -e

# Configuration
ENVIRONMENT=${1:-production}
PROJECT_ID=${2:-your-gcp-project-id}
REGION="asia-south1"
SERVICE_NAME="gymlog"
IMAGE_NAME="gcr.io/$PROJECT_ID/$SERVICE_NAME"

echo "🚀 Deploying GymLog Backend to Google Cloud Run"
echo "Environment: $ENVIRONMENT"
echo "Project ID: $PROJECT_ID"
echo "Region: $REGION"
echo "Service: $SERVICE_NAME"
echo ""

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo "❌ gcloud CLI is not installed. Please install it first."
    echo "Visit: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Verify secrets exist
echo "🔐 Verifying Secret Manager secrets exist..."
REQUIRED_SECRETS=("mongodb-uri" "db-name" "jwt-secret")
for secret in "${REQUIRED_SECRETS[@]}"; do
    if ! gcloud secrets describe $secret --project $PROJECT_ID &> /dev/null; then
        echo "❌ Secret '$secret' not found. Please create it first."
        echo "Run: echo 'your-value' | gcloud secrets create $secret --data-file=-"
        exit 1
    fi
done
echo "✅ All required secrets found"

# Build and tag the image
echo "🔨 Building Docker image..."
docker build -t $IMAGE_NAME:latest .

# Push to Google Container Registry
echo "📤 Pushing image to Google Container Registry..."
docker push $IMAGE_NAME:latest

# Deploy to Cloud Run with secrets
echo "🚀 Deploying to Cloud Run with Secret Manager integration..."
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_NAME:latest \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10 \
  --set-env-vars ENV=$ENVIRONMENT \
  --update-secrets MONGODB_URI=mongodb-uri:latest,DB_NAME=db-name:latest,JWT_SECRET=jwt-secret:latest \
  --project $PROJECT_ID

# Get the service URL
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --platform managed --region $REGION --format 'value(status.url)' --project $PROJECT_ID)

echo ""
echo "✅ Deployment completed successfully!"
echo "🌐 Service URL: $SERVICE_URL"
echo "🔐 Using secrets from Secret Manager:"
echo "   - MONGODB_URI from mongodb-uri"
echo "   - DB_NAME from db-name" 
echo "   - JWT_SECRET from jwt-secret"
echo ""
echo "Next steps:"
echo "1. Verify your MongoDB Atlas connection string is correct"
echo "2. Test the deployment with a health check"
echo "3. Update your frontend to use the new API endpoint"
echo "" 