#!/bin/bash

# Test script for GymLog Insights Service
# Usage: ./test_insights.sh <user_id>

if [ -z "$1" ]; then
    echo "Usage: $0 <user_id>"
    echo "Example: $0 507f1f77bcf86cd799439011"
    exit 1
fi

USER_ID=$1
BASE_URL=${BASE_URL:-"http://localhost:8080"}

echo "🧪 Testing Insights Service for User: $USER_ID"
echo "📍 Base URL: $BASE_URL"
echo ""

# Generate insights
echo "1️⃣ Generating insights..."
curl -X POST "$BASE_URL/api/users/$USER_ID/insights" \
  -H "Content-Type: application/json" \
  -w "\n\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo ""
echo "2️⃣ Getting recent insights (limit=3)..."
curl -X GET "$BASE_URL/api/users/$USER_ID/insights?limit=3" \
  -H "Content-Type: application/json" \
  -w "\n\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo ""
echo "✅ Test completed!" 