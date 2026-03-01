# RepTracker Backend

Go backend service providing HTTP REST and gRPC APIs for the RepTracker application, with MongoDB storage, Clerk authentication, and Google Gemini AI integration.

## Table of Contents

- [Setup](#setup)
- [Running the Server](#running-the-server)
- [Authentication](#authentication)
- [API Reference](#api-reference)
- [Database Schema](#database-schema)
- [Metrics System](#metrics-system)
- [AI Services](#ai-services)
- [Exercise Constants](#exercise-constants)
- [Scripts](#scripts)
- [Deployment](#deployment)

---

## Setup

### Prerequisites

- Go 1.23+
- MongoDB (local or Atlas)
- Clerk account (publishable key)
- Google Gemini API key (for AI features)

### Environment Variables

Create a `.env` file in the `backend/` directory:

```env
MONGODB_URI=mongodb://localhost:27017
DB_NAME=gymlog
EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_your-clerk-publishable-key
GEMINI_API_KEY=your-gemini-api-key
```

| Variable | Required | Description |
|----------|----------|-------------|
| `MONGODB_URI` | Yes | MongoDB connection string |
| `DB_NAME` | Yes | MongoDB database name |
| `EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY` | Yes | Clerk publishable key for JWT validation |
| `GEMINI_API_KEY` | No | Google Gemini API key (AI features disabled without it) |

### Install Dependencies

```bash
go mod tidy
```

---

## Running the Server

### Locally

```bash
# Start MongoDB
mongod

# Initialize collections (first time only)
mongosh < scripts/init-mongo.js

# Populate exercise database (first time only)
go run scripts/populate_exercises.go

# Start the server
go run cmd/server/main.go
```

The server starts on:
- HTTP: `http://localhost:8080`
- gRPC: `localhost:50051`

### With Docker Compose

```bash
make docker-compose-up
```

### Makefile Commands

```bash
make setup              # Setup development environment
make run                # Start the server
make test               # Run tests
make build              # Build binary
make docker-compose-up  # Start with MongoDB via Docker
make populate-exercises # Populate exercise database
```

---

## Authentication

All API endpoints (except `GET /health`) require a valid Clerk JWT token:

```
Authorization: Bearer <clerk-jwt-token>
```

**Flow:**
1. User signs in via Clerk on the frontend
2. Frontend obtains a JWT from the Clerk session
3. Frontend includes the JWT in the `Authorization` header
4. Backend validates the token against Clerk's JWKS endpoint
5. Backend extracts the Clerk user ID from token claims
6. Backend creates or retrieves the user record in MongoDB

**Notes:**
- No passwords are stored in the backend (handled by Clerk)
- Tokens are short-lived and automatically refreshed by Clerk
- User records are created automatically on first authenticated request

---

## API Reference

### Health Check

```
GET /health
```

### Users

```
GET  /api/users      # Get authenticated user profile
PUT  /api/users      # Create or update user profile (upsert)
```

The user ID is extracted from the JWT token — no need to pass it in the URL.

**GET /api/users response:**
```json
{
  "id": "user_id",
  "clerkId": "clerk_user_id",
  "email": "user@example.com",
  "name": "John Doe",
  "firstName": "John",
  "lastName": "Doe",
  "height": 180,
  "weight": 75,
  "age": 30,
  "goal": "gain_muscle",
  "picture": "https://...",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:35:00Z"
}
```

Returns `404` if user hasn't completed onboarding.

**PUT /api/users request body:**
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "height": 180,
  "weight": 75,
  "age": 30,
  "goal": "gain_muscle"
}
```

Goal must be one of: `lose_fat`, `gain_muscle`, `maintain`.

---

### Exercises

```
GET    /api/exercises                    # List exercises (with filtering)
POST   /api/exercises                    # Create exercise
GET    /api/exercises/{id}               # Get exercise by ID
PUT    /api/exercises/{id}               # Update exercise
DELETE /api/exercises/{id}               # Delete exercise
GET    /api/exercises/quick-add          # Get exercises for quick-add
GET    /api/exercises/{id}/history       # Get exercise history (last 3 sessions)
```

**GET /api/exercises query parameters:**
- `search` — Search by name
- `category` — Filter by category
- `equipment` — Filter by equipment
- `muscleGroup` — Filter by muscle group

**GET /api/exercises/{id}/history** returns the last 3 completed workout sessions where the exercise was performed, including actual vs. target performance and any AI progressive overload recommendations.

---

### Workouts (Routines)

```
GET    /api/workouts          # List workouts for authenticated user
POST   /api/workouts          # Create workout
GET    /api/workouts/{id}     # Get workout by ID
GET    /api/workouts/{id}/start  # Get workout with exercise details for starting
PUT    /api/workouts/{id}     # Update workout
DELETE /api/workouts/{id}     # Delete workout
```

---

### Workout Sessions

```
GET    /api/workout-sessions                                          # List sessions
POST   /api/workout-sessions                                          # Create session
GET    /api/workout-sessions/{id}                                     # Get session
PUT    /api/workout-sessions/{id}                                     # Update session
DELETE /api/workout-sessions/{id}                                     # Delete session
POST   /api/workout-sessions/{id}/exercises/{exerciseIndex}/start     # Start exercise
POST   /api/workout-sessions/{id}/exercises/{exerciseIndex}/finish    # Finish exercise
PUT    /api/workout-sessions/{id}/exercises/{exerciseIndex}/sets/{setIndex}  # Update set
POST   /api/workout-sessions/{id}/progressive-overload                # Apply progressive overload
POST   /api/workout-sessions/{id}/progressive-overload/ai             # Apply AI progressive overload (async)
```

**Progressive Overload Logic** (`POST /api/workout-sessions/{id}/progressive-overload`):

For each completed exercise in the session:
1. Increase rep count by 1 for all completed sets
2. If rep count exceeds 15 and weight is non-zero: increase weight by 2.5 and reset reps to 8
3. Bodyweight exercises (weight = 0) are not weight-progressed

Request body: `{ "workoutId": "workout_id_to_update" }`

**AI Progressive Overload** (`POST /api/workout-sessions/{id}/progressive-overload/ai`):

Runs asynchronously — returns immediately with `success: true`, then processes AI analysis in the background. Results are stored and retrievable via `GET /api/users/ai-progressive-overload-responses`.

Request body: `{ "workoutId": "workout_id_to_update" }`

---

### Metrics

```
GET /api/users/metrics            # Get user metrics (period: weekly|monthly|all)
GET /api/users/metrics/trends     # Get volume trends (period: weekly|monthly)
GET /api/workout-sessions/{id}/metrics  # Get metrics for a specific session
```

---

### Insights

```
POST /api/users/insights    # Generate AI workout insights
GET  /api/users/insights    # Get recent insights (query: limit=5)
```

---

### Workout Suggestions

```
POST /api/users/suggestions/workouts              # Generate AI workout suggestions
GET  /api/users/suggestions/stored                # Get stored suggestions (paginated)
POST /api/users/suggestions/{suggestionId}/confirm  # Accept or reject a suggestion
```

**POST /api/users/suggestions/workouts request body (optional):**
```json
{
  "daysToAnalyze": 14,
  "maxSuggestions": 3
}
```

**POST /api/users/suggestions/{suggestionId}/confirm request body:**
```json
{
  "suggestionIndex": 0,
  "accept": true
}
```

When accepted, the original workout is updated with the suggested changes.

---

### AI Progressive Overload History

```
GET /api/users/ai-progressive-overload-responses  # List AI analysis history (paginated)
```

Query parameters: `pageSize` (default 20, max 100), `pageToken` (for pagination).

---

## Database Schema

### Collections

The following collections are defined in `scripts/init-mongo.js`:

| Collection | Description |
|------------|-------------|
| `users` | User accounts and profiles |
| `sessions` | Authentication sessions |
| `exercises` | Exercise definitions and metadata |
| `workouts` | Workout routines/templates |
| `workout_sessions` | Active and completed workout sessions |
| `workout_suggestions` | AI-generated workout suggestions |
| `ai_progressive_overload_recommendations` | AI progressive overload analysis |

**MongoDB field naming:** All MongoDB queries use `snake_case` field names (e.g., `user_id`, `finished_at`). Go struct tags define both BSON (snake_case) and JSON (camelCase) representations.

### Users Collection

```json
{
  "_id": "ObjectId",
  "clerk_id": "string",
  "email": "string",
  "name": "string",
  "picture": "string",
  "height": "number",
  "weight": "number",
  "age": "number",
  "goal": "string",
  "created_at": "Date",
  "updated_at": "Date"
}
```

---

## Metrics System

Metrics are automatically calculated when a workout session is completed (`FinishedAt` is set). The service is implemented in `internal/services/metrics_service.go`.

### Metric Categories

| Category | Key Metrics |
|----------|-------------|
| Volume | Total Volume Load, Tonnage, Volume per Muscle Group, Effective Reps, Hard Sets |
| Performance | RPE Distribution, Set Completion Rate |
| Intensity | Average Intensity, RPE-Adjusted Load, Intensity Distribution, Load Density |
| Strength | Estimated 1RM (Epley & Brzycki), Wilks Score, Push:Pull Ratio, Power Output |
| Progress & Adaptation | Progressive Overload Index, Week-over-Week Progress Rate, Plateau Detection, Strength Gain Velocity |
| Recovery & Fatigue | Acute/Chronic Training Load (ATL/CTL), Training Stress Balance (TSB), Form/Freshness Index |
| Muscle-Specific | Muscle Group Distribution, Imbalance Index, Antagonist Ratios, Stimulus-to-Fatigue Ratio |
| Work Capacity | Total Work Capacity, Density Training Index, Time Under Tension, Mechanical Tension Score |
| Training Patterns | Optimal Frequency, Recovery Time, Exercise Diversity Index, Consistency Score |
| Injury Risk | Injury Risk Score, Load Spike Alert, Asymmetry Development |
| Efficiency | Strength Efficiency, Volume Efficiency, RPE-Performance Correlation, Technique Consistency |
| Periodization | CTL (42-day), ATL (7-day), TSB, Form/Freshness Index |
| Compound | Fitness Age, Overall Fitness Score, Training Efficiency Quotient |
| Predictive | Plateau Probability, Performance Trajectory, Goal Achievement Timeline, Detraining Risk |

### Volume Landmarks (MEV/MAV/MRV)

Calculated from 8 weeks of historical data per muscle group:
- **MEV** (Minimum Effective Volume): `max(avg × 0.75, avg - std_dev)`
- **MAV** (Maximum Adaptive Volume): `avg + (std_dev × 0.5)`
- **MRV** (Maximum Recoverable Volume): `avg + (std_dev × 1.5)`

### Key Formulas

- **1RM Epley**: `Weight × (1 + Reps/30)`
- **1RM Brzycki**: `Weight × (36 / (37 - Reps))`
- **Progressive Overload Index**: `(Current Week Volume × Avg Intensity) / (Previous Week Volume × Avg Intensity)`
- **Injury Risk Score**: `(ACWR × Imbalance Index × Fatigue Score) / Recovery Quality` — high risk if > 1.5
- **Load Spike Alert**: triggered when weekly volume > 150% of 4-week rolling average

---

## AI Services

All AI features use **Google Gemini 2.5 Flash**. Set `GEMINI_API_KEY` in `.env` to enable them.

### Insights Service (`internal/services/insights_service.go`)

Analyzes the last 7 days of workout metrics and generates prioritized, single-sentence actionable insights.

**Insight types:**
- `progress` — Progressive overload, plateau detection, 1RM trends
- `volume` — MEV/MAV/MRV comparisons, muscle group distribution
- `recovery` — Training stress balance, optimal frequency
- `balance` — Push/pull ratios, muscle imbalances
- `risk` — Injury risk, load spikes, asymmetries

**Priority levels:** 5 (critical) → 1 (informational)

### Workout Suggestions Service (`internal/services/suggestions_service.go`)

Analyzes recent workout history and generates suggested routine improvements. Suggestions have a `pending` status until accepted or rejected via the confirm endpoint.

### AI Progressive Overload Service

Runs asynchronously. Uses Gemini to analyze recent session performance and generate exercise-specific progression recommendations with reasoning. Results stored in `ai_progressive_overload_recommendations` collection.

---

## Exercise Constants

Defined in `pkg/models/exercise_constants.go`. Always use these typed constants instead of raw strings.

### Categories

```go
models.CategoryStrength             // "strength"
models.CategoryCardio               // "cardio"
models.CategoryStretching           // "stretching"
models.CategoryPlyometrics          // "plyometrics"
models.CategoryStrongman            // "strongman"
models.CategoryOlympicWeightlifting // "olympic weightlifting"
models.CategoryCrossfit             // "crossfit"
models.CategoryCalisthenics         // "calisthenics"
```

### Equipment

```go
models.EquipmentNone         // "none"
models.EquipmentBarbell      // "barbell"
models.EquipmentDumbbell     // "dumbbell"
models.EquipmentKettlebell   // "kettlebell"
models.EquipmentMachine      // "machine"
models.EquipmentCable        // "cable"
models.EquipmentBands        // "bands"
models.EquipmentPullUpBar    // "pull-up bar"
models.EquipmentBench        // "bench"
// ... and more
```

### Muscle Groups

```go
models.MuscleGroupArms       // "arms"
models.MuscleGroupBack       // "back"
models.MuscleGroupChest      // "chest"
models.MuscleGroupCore       // "core"
models.MuscleGroupLegs       // "legs"
models.MuscleGroupShoulders  // "shoulders"
models.MuscleGroupCalves     // "calves"
```

### Helper Functions

```go
models.GetAllCategories()                          // All category values
models.GetAllEquipment()                           // All equipment values
models.GetAllMuscles()                             // All muscle values
models.GetMusclesInGroup(models.MuscleGroupLegs)   // Muscles in a group
models.GetMuscleGroup(models.MuscleBiceps)         // Group for a muscle
models.IsValidCategory(category)                   // Validate category
models.IsValidEquipment(equipment)                 // Validate equipment
models.IsValidMuscle(muscle)                       // Validate muscle
```

### Exercise Filter Model

```go
query := models.NewExerciseQuery()
query.Filter.Categories = []models.ExerciseCategory{models.CategoryStrength}
query.Filter.Equipment = []models.EquipmentType{models.EquipmentBarbell}
query.Filter.SearchTerm = "bench press"
query.Pagination.Page = 1
query.Pagination.Limit = 20
```

---

## Scripts

Located in `backend/scripts/`.

### init-mongo.js

Initializes MongoDB collections with indexes. Run once on first setup:

```bash
mongosh < scripts/init-mongo.js
```

### populate_exercises.go

Clears and repopulates the `exercises` collection from `scripts/gymlog.exercises.json` (325 exercises from StrengthLog).

```bash
# Using the shell wrapper (recommended)
./scripts/populate_exercises.sh

# Or directly
go run scripts/populate_exercises.go
```

**Warning:** This deletes all existing exercises before inserting. Safe to run multiple times (idempotent).

Default connection: `mongodb://localhost:27017`, database `gymlog`. Override with `MONGODB_URI` and `DB_NAME` environment variables.

---

## Deployment

The backend is deployed to Google Cloud Run. The Docker image is built and pushed via Cloud Build.

### Prerequisites

```bash
gcloud auth login
gcloud config set project YOUR_PROJECT_ID
gcloud auth configure-docker
```

### Deploy

```bash
# Manual deployment using deploy.sh
./deploy.sh production YOUR_GCP_PROJECT_ID
```

The script automatically pulls secrets from Google Secret Manager and deploys to Cloud Run.

### Cloud Build (CI/CD)

```bash
gcloud builds triggers create github \
  --repo-name=YOUR_GITHUB_REPO \
  --repo-owner=YOUR_GITHUB_USERNAME \
  --branch-pattern="^main$" \
  --build-config=backend/cloudbuild.yaml
```

### Manual gcloud Deployment

```bash
# Build and push image
docker build -t gcr.io/YOUR_PROJECT_ID/gymlog-backend .
docker push gcr.io/YOUR_PROJECT_ID/gymlog-backend

# Deploy to Cloud Run
gcloud run deploy gymlog-backend \
  --image gcr.io/YOUR_PROJECT_ID/gymlog-backend \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10
```

### Google Secret Manager

Secrets are stored in Secret Manager and injected at deploy time. The GCP project is `gymlog-462803`.

```bash
# List secrets
gcloud secrets list

# Create a secret
echo "value" | gcloud secrets create SECRET_NAME --data-file=-

# Update a secret
echo "new-value" | gcloud secrets versions add SECRET_NAME --data-file=-

# Read a secret (for debugging)
gcloud secrets versions access latest --secret="SECRET_NAME"
```

Required secrets: `mongodb-uri`, `db-name`, `EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY`, `GEMINI_API_KEY`.

### Production Environment Variables

| Variable | Secret Name | Notes |
|----------|-------------|-------|
| `MONGODB_URI` | `mongodb-uri` | MongoDB Atlas connection string |
| `DB_NAME` | `db-name` | `gymlog` |
| `EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY` | `clerk-publishable-key` | Use `pk_live_...` in production |
| `GEMINI_API_KEY` | `gemini-api-key` | Google AI Studio key |

### Scaling

```bash
gcloud run services update gymlog-backend \
  --min-instances=1 \
  --max-instances=10 \
  --concurrency=80 \
  --cpu=1 \
  --memory=512Mi \
  --region us-central1
```

### Viewing Logs

```bash
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=gymlog-backend" --limit=50
```

### Troubleshooting Deployment

- **MongoDB connection issues**: Check Atlas network access allows Cloud Run IP ranges
- **Secret not found**: Verify secret exists with `gcloud secrets list` and Cloud Build SA has `roles/secretmanager.secretAccessor`
- **Build failures**: Check `go.mod` dependencies and Dockerfile syntax
