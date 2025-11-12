#!/bin/bash

# Test script for GymLog Async AI Progressive Overload
# Usage: ./test_async_ai_progressive_overload.sh <session_id> <workout_id>

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: $0 <session_id> <workout_id>"
    echo "Example: $0 507f1f77bcf86cd799439011 507f1f77bcf86cd799439012"
    exit 1
fi

SESSION_ID=$1
WORKOUT_ID=$2
BASE_URL=${BASE_URL:-"http://localhost:8080"}

echo "🤖 Testing Async AI Progressive Overload"
echo "🔍 Session ID: $SESSION_ID"
echo "🏋️  Workout ID: $WORKOUT_ID"
echo "📍 Base URL: $BASE_URL"
echo ""

# Apply AI Progressive Overload (now async)
echo "1️⃣ Triggering AI Progressive Overload (async)..."
RESPONSE=$(curl -X POST "$BASE_URL/api/workout-sessions/$SESSION_ID/apply-progressive-overload" \
  -H "Content-Type: application/json" \
  -d "{\"workoutId\": \"$WORKOUT_ID\"}" \
  -w "\n\nHTTP Status: %{http_code}\n" \
  -s)

echo "$RESPONSE" | jq '.'
HTTP_CODE=$(echo "$RESPONSE" | tail -n1 | grep -o '[0-9]*')

if [ "$HTTP_CODE" = "200" ]; then
    echo ""
    echo "✅ AI processing started successfully! Now checking status..."
    echo ""
    
    # Extract user ID from session (you might need to adjust this based on your session structure)
    # For this test, we'll assume you provide it or we can extract it
    echo "2️⃣ Fetching AI Progressive Overload responses to check status..."
    echo "   (Note: You may need to provide the correct user_id for this endpoint)"
    echo ""
    
    # Wait a moment for processing to potentially complete
    echo "⏳ Waiting 2 seconds to allow processing..."
    sleep 2
    
    # List AI responses - you'll need to replace USER_ID with actual user ID
    echo "3️⃣ Checking AI Progressive Overload responses..."
    echo "   Usage: curl -X GET \"$BASE_URL/api/users/{user_id}/ai-progressive-overload-responses\""
    echo ""
    echo "💡 To check the status of your async AI processing:"
    echo "   Replace {user_id} with the actual user ID and run:"
    echo "   curl -X GET \"$BASE_URL/api/users/{user_id}/ai-progressive-overload-responses\" | jq '.'"
    
else
    echo ""
    echo "❌ AI Progressive Overload failed with HTTP status: $HTTP_CODE"
fi

echo ""
echo "✅ Test completed!"
echo ""
echo "📝 How to interpret results:"
echo "   - The initial call should return immediately with success:true and a 'processing' message"
echo "   - Use the ai-progressive-overload-responses endpoint to check status and results"
echo "   - Look for the 'success' field and 'message' in the response to see processing status"
