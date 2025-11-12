# GymLog Backend

A Go backend service for the GymLog application with MongoDB integration and Google OAuth authentication.

## Setup

### Prerequisites

- Go 1.23+
- MongoDB (local or cloud)
- Environment variables

### Environment Variables

Create a `.env` file in the backend directory with the following variables:

```env
# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
DB_NAME=gymlog

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here

# Google OAuth Configuration  
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# Google Gemini API Configuration
GEMINI_API_KEY=your-gemini-api-key
```

### Installation

1. Install dependencies:
```bash
go mod tidy
```

2. Start MongoDB (if running locally):
```bash
mongod
```

3. Start the server:
```bash
go run cmd/server/main.go
```

The server will start both:
- gRPC server on port 50051
- HTTP server on port 8080

## API Endpoints

### Authentication
- `POST /auth/google` - Google OAuth authentication
- `POST /auth/validate` - Validate session token  
- `POST /auth/logout` - Logout user

### Users
- `POST /api/users` - Create user
- `GET /api/users/{userId}` - Get user by ID
- `PUT /api/users/{userId}` - Update user
- `GET /api/users` - List users

### Exercises
- `GET /api/exercises` - List exercises with filtering
- `POST /api/exercises` - Create exercise
- `GET /api/exercises/{id}` - Get exercise by ID
- `PUT /api/exercises/{id}` - Update exercise
- `DELETE /api/exercises/{id}` - Delete exercise
- `GET /api/exercises/quick-add` - Get exercises for quick add
- `GET /api/exercises/{id}/history?user_id={userId}` - Get exercise history from last 3 workout sessions

### Workouts (Routines)
- `GET /api/workouts` - List workouts
- `POST /api/workouts` - Create workout
- `GET /api/workouts/{id}` - Get workout by ID
- `GET /api/workouts/{id}/start` - Get workout with exercise details for starting
- `PUT /api/workouts/{id}` - Update workout
- `DELETE /api/workouts/{id}` - Delete workout

### Workout Sessions
- `GET /api/users/{userId}/workout-sessions` - List workout sessions
- `POST /api/users/{userId}/workout-sessions` - Create workout session
- `GET /api/workout-sessions/{id}` - Get workout session
- `PUT /api/workout-sessions/{id}` - Update workout session
- `DELETE /api/workout-sessions/{id}` - Delete workout session
- `POST /api/workout-sessions/{id}/exercises/{exerciseIndex}/start` - Start exercise
- `POST /api/workout-sessions/{id}/exercises/{exerciseIndex}/finish` - Finish exercise
- `PUT /api/workout-sessions/{id}/exercises/{exerciseIndex}/sets/{setIndex}` - Update set
- `POST /api/workout-sessions/{id}/progressive-overload` - Apply progressive overload to workout
- `POST /api/workout-sessions/{id}/progressive-overload/ai` - Apply AI-powered progressive overload (async)
- `GET /api/users/{userId}/ai-progressive-overload-responses` - List stored AI progressive overload responses

#### Progressive Overload API Details

**Apply Progressive Overload**
```
POST /api/workout-sessions/{id}/progressive-overload
```
Applies progressive overload to a workout based on the performance recorded in a workout session. This automatically increases the difficulty of the workout for future use.

**Progressive Overload Logic:**
For each exercise in the workout session:
1. **Increase rep count by 1** for all completed sets
2. **If rep count > 15 after step 1:**
   - If weight is non-zero: increase weight by 2.5 and set rep count to 8
   - If weight is zero: do nothing (bodyweight exercises)
3. **Skip exercises** that don't exist in the target workout
4. **Only process completed sets** from the session

Request body:
```json
{
  "workoutId": "workout_id_to_update"
}
```

Response:
```json
{
  "success": true,
  "message": "Progressive overload applied successfully",
  "updatedWorkout": {
    "id": "workout_id",
    "userId": "user_id",
    "name": "Push Day",
    "description": "Updated workout with progressive overload",
    "exercises": [
      {
        "exerciseId": "exercise_id",
        "exercise": {...},
        "sets": [
          {
            "reps": 9,     // Was 8, increased by 1
            "weight": 50.0,
            "durationSeconds": 0,
            "distance": 0,
            "notes": ""
          },
          {
            "reps": 8,     // Was 16, weight increased and reps reset
            "weight": 22.5, // Was 20, increased by 2.5
            "durationSeconds": 0,
            "distance": 0,
            "notes": ""
          }
        ],
        "notes": "",
        "restSeconds": 90
      }
    ],
    "startedAt": null,
    "finishedAt": null,
    "durationSeconds": 0,
    "notes": "",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:35:00Z"
  }
}
```

**Use Case:**
After completing a workout session, call this API to automatically progress the workout difficulty for the next time the user performs this routine. This ensures continuous progression without manual intervention.

#### AI Progressive Overload API Details

**Apply AI-Powered Progressive Overload (Async)**
```
POST /api/workout-sessions/{id}/progressive-overload/ai
```
Applies AI-powered progressive overload analysis to a workout using Gemini AI. This endpoint now **runs asynchronously** - it immediately returns a response indicating that processing has started, then performs the AI analysis and workout updates in the background.

**Request Body:**
```json
{
  "workoutId": "workout_id_to_update"
}
```

**Immediate Response:**
```json
{
  "success": true,
  "message": "AI progressive overload analysis started. Processing in background...",
  "updatedWorkout": {
    "id": "original_workout_id",
    "name": "Push Day",
    "exercises": [/* original exercises */]
  }
}
```

**Async Processing Flow:**
1. **Immediate Return**: API returns immediately with `success: true` and processing message
2. **Background Processing**: AI analysis, workout updates, and database storage happen asynchronously  
3. **Status Tracking**: Results are stored and can be queried via the List AI Progressive Overload Responses endpoint
4. **Error Handling**: Any errors during async processing are captured and stored with `success: false`

**To Check Processing Status and Results:**
Use the List AI Progressive Overload Responses endpoint below to monitor the status and retrieve results of your async AI analysis.

**Benefits of Async Processing:**
- **Non-blocking**: API responds immediately without waiting for AI processing
- **Better UX**: Frontend can show immediate feedback and poll for results
- **Reliability**: Long-running AI analysis doesn't risk request timeouts
- **Scalability**: Server can handle multiple concurrent AI requests efficiently

**List AI Progressive Overload Responses**
```
GET /api/users/{userId}/ai-progressive-overload-responses
```
Retrieves stored AI progressive overload responses for a user in descending order by creation time. This endpoint allows you to view the history of AI recommendations and their effectiveness.

**Query Parameters:**
- `pageSize` (optional): Number of responses per page (default: 20, max: 100)
- `pageToken` (optional): Token for pagination to get next page of results

**Response:**
```json
{
  "responses": [
    {
      "id": "response_id",
      "userId": "user_id",
      "workoutSessionId": "session_id",
      "workoutId": "workout_id",
      "analysisSummary": "Strong progression shown across all exercises. User consistently completing target reps with room for advancement.",
      "updatedExercises": [
        {
          "exerciseId": "exercise_id",
          "exerciseName": "Bench Press",
          "progressionApplied": true,
          "reasoning": "User completed all sets with 2+ RIR, ready for weight increase",
          "changesMade": "Increased weight from 80kg to 82.5kg based on consistent completion",
          "sets": [
            {
              "reps": 8,
              "weight": 82.5,
              "durationSeconds": 0,
              "distance": 0,
              "notes": ""
            }
          ],
          "notes": "Strong performance, continue progression",
          "restSeconds": 180
        }
      ],
      "success": true,
      "message": "AI progressive overload applied successfully - Strong progression trends identified",
      "recentSessionsCount": 7,
      "aiModel": "gemini-2.5-flash",
      "processingTimeMs": 1245,
      "createdAt": "2024-01-15T10:35:00Z",
      "updatedAt": "2024-01-15T10:35:00Z"
    }
  ],
  "nextPageToken": "next_page_token_if_more_results"
}
```

**Use Case:**
Track the history of AI recommendations to understand progression patterns, analyze AI decision-making, and monitor the effectiveness of different progression strategies over time.

### Metrics
- `GET /api/users/{userId}/metrics` - Get user metrics
- `GET /api/users/{userId}/metrics/trends` - Get volume trends
- `GET /api/workout-sessions/{sessionId}/metrics` - Get workout metrics

### Insights
- `POST /api/users/{userId}/insights` - Generate workout insights
- `GET /api/users/{userId}/insights/recent` - Get recent insights

### Workout Suggestions
- `POST /api/users/{userId}/suggestions/workouts` - Generate AI workout suggestions
- `GET /api/users/{userId}/suggestions/stored` - Get stored workout suggestions (paginated, ordered by creation date desc)

#### Workout Suggestions API Details

**Generate Workout Suggestions**
```
POST /api/users/{userId}/suggestions/workouts
```
Generates AI-powered workout suggestions based on recent workout sessions and stores them in the database.

Request body (optional):
```json
{
  "daysToAnalyze": 14,    // Number of days of workout history to analyze (default: 14)
  "maxSuggestions": 3     // Maximum number of suggestions to generate (default: 3)
}
```

Response:
```json
{
  "suggestions": [
    {
      "originalWorkoutId": "workout_id",
      "name": "Improved Push Day",
      "description": "Enhanced push workout with better progression",
      "exercises": [...],
      "changes": [
        {
          "type": "sets",
          "exerciseId": "exercise_id", 
          "exerciseName": "Bench Press",
          "oldValue": "3 sets",
          "newValue": "4 sets",
          "reason": "Increase volume for better strength gains"
        }
      ],
      "overallReasoning": "Based on your recent progress...",
      "priority": 5
    }
  ],
  "analysisSummary": "Analysis of your recent workout patterns shows..."
}
```

**Get Stored Suggestions**
```
GET /api/users/{userId}/suggestions/stored?pageSize=10&pageToken=...
```
Retrieves previously generated workout suggestions with "pending" status in decreasing order of creation date. Only suggestions that haven't been accepted or rejected are returned.

Query parameters:
- `pageSize` (optional): Number of suggestion sets per page (default: 10, max: 50)
- `pageToken` (optional): Token for pagination (timestamp-based)

Response:
```json
{
  "storedSuggestions": [
    {
      "id": "suggestion_set_id",
      "userId": "user_id", 
      "suggestions": [
        {
          "originalWorkoutId": "workout_id",
          "name": "Improved Push Day",
          "description": "Enhanced push workout",
          "exercises": [...],
          "changes": [...],
          "overallReasoning": "Based on your recent progress...",
          "priority": 5,
          "status": "pending",  // "pending", "accepted", "rejected"
          "statusUpdatedAt": "2024-01-15T10:35:00Z"  // When status was last changed
        }
      ],
      "analysisSummary": "Analysis summary from when suggestions were generated",
      "daysAnalyzed": 14,
      "createdAt": "2024-01-15T10:30:00Z",
      "updatedAt": "2024-01-15T10:30:00Z"
    }
  ],
  "nextPageToken": "2024-01-14T10:30:00Z"  // For pagination
}
```

**Confirm Suggestion (Accept/Reject)**
```
POST /api/users/{userId}/suggestions/{suggestionId}/confirm
```
Accepts or rejects a specific workout suggestion. When accepted, applies the suggested changes to the original workout.

Request body:
```json
{
  "suggestionIndex": 0,  // Index of the suggestion within the stored suggestion set
  "accept": true  // true to accept, false to reject
}
```

Response:
```json
{
  "success": true,
  "message": "Suggestion accepted and applied to workout"
}
```

When a suggestion is accepted:
- The original workout is updated with the suggested changes
- The suggestion status is marked as "accepted"
- The `statusUpdatedAt` timestamp is set

When a suggestion is rejected:
- The suggestion status is marked as "rejected" 
- The `statusUpdatedAt` timestamp is set
- No changes are applied to the original workout

### Health Check
- `GET /health` - Health check endpoint

## Database Schema

### Users Collection
```json
{
  "_id": "ObjectId",
  "email": "string",
  "name": "string", 
  "google_id": "string",
  "picture": "string",
  "created_at": "Date",
  "updated_at": "Date"
}
```

### Sessions Collection
```json
{
  "_id": "ObjectId",
  "user_id": "ObjectId",
  "access_token": "string",
  "refresh_token": "string",
  "expires_at": "Date",
  "created_at": "Date"
}
```

#### Exercise History API Details

**Get Exercise History**
```
GET /api/exercises/{id}/history?user_id={userId}
```
Retrieves the workout session history for a specific exercise from the last 3 completed workout sessions where the exercise was performed. Includes both the actual performance data and AI progressive overload recommendations if available.

**Query Parameters:**
- `user_id` (required): Filter history for a specific user.

**Response:**
```json
{
  "history": [
    {
      "sessionExercise": {
        "exerciseId": "exercise_id",
        "exercise": {
          "id": "exercise_id",
          "name": "Bench Press",
          "description": "...",
          "category": "strength",
          "primaryMuscles": ["chest"],
          "secondaryMuscles": ["triceps", "shoulders"]
        },
        "sets": [
          {
            "targetReps": 10,
            "targetWeight": 135.0,
            "actualReps": 10,
            "actualWeight": 135.0,
            "completed": true,
            "startedAt": "2024-01-15T10:00:00Z",
            "finishedAt": "2024-01-15T10:02:00Z"
          }
        ],
        "notes": "Felt strong today",
        "restSeconds": 180,
        "completed": true,
        "startedAt": "2024-01-15T10:00:00Z",
        "finishedAt": "2024-01-15T10:05:00Z"
      },
      "aiExercise": {
        "exerciseId": "exercise_id",
        "exerciseName": "Bench Press",
        "progressionApplied": true,
        "reasoning": "Strong consistent performance, ready for weight increase",
        "changesMade": "Increased weight from 135 to 137.5 lbs",
        "sets": [
          {
            "reps": 10,
            "weight": 137.5
          }
        ],
        "restSeconds": 180
      },
      "sessionInfo": {
        "id": "session_id",
        "name": "Push Day Workout",
        "finishedAt": "2024-01-15T11:00:00Z"
      }
    }
  ],
  "exerciseName": "Bench Press"
}
```

**Features:**
- Returns last 3 workout sessions where the exercise was performed by the specified user
- Includes complete exercise performance data (target vs actual reps/weights)
- Includes AI progressive overload recommendations when available
- Provides session context (workout name, completion time)
- Requires user-specific filtering for data privacy and performance
- Automatically excludes incomplete/active sessions

## Documentation

- **[METRICS.md](METRICS.md)** - Comprehensive guide to the workout metrics system, including all calculated metrics, formulas, API endpoints, and usage examples
- **[INSIGHTS_SERVICE.md](INSIGHTS_SERVICE.md)** - Guide to the AI-powered insights service using Google Gemini, including setup, usage, and API endpoints

## Authentication Flow

1. Frontend initiates Google OAuth
2. NextAuth handles OAuth callback
3. Frontend calls `/api/auth/google` with user info
4. Backend creates/updates user in MongoDB
5. Backend creates session with access token
6. Frontend stores session token
7. All subsequent API calls include session token
8. Backend validates session on each request
9. Expired sessions redirect to login 