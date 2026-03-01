# GymLog

A full-stack fitness tracking application for creating workout routines, tracking live sessions, and analyzing performance with AI-powered insights.

## Features

- **Exercise Library** — 325+ exercises with descriptions, muscle groups, equipment, and instructions
- **Routine Builder** — Create custom workout routines with sets, reps, weights, and rest periods
- **Live Workout Tracking** — Real-time session tracking with set-by-set completion and timers
- **Comprehensive Analytics** — 60+ workout metrics including volume, strength, recovery, and periodization
- **AI Insights** — Google Gemini-powered analysis of training patterns, volume landmarks (MEV/MAV/MRV), and injury risk
- **AI Progressive Overload** — Async AI-driven workout progression recommendations
- **Workout Suggestions** — AI-generated routine improvements based on recent history
- **Authentication** — Clerk-based email/password auth with JWT validation

## Architecture

```
GymLog/
├── frontend-native/    # React Native (Expo) iOS/Android app
└── backend/            # Go HTTP + gRPC server with MongoDB
```

```
┌─────────────────────┐     HTTPS/JWT      ┌──────────────────────┐
│  React Native App   │ ─────────────────► │   Go Backend         │
│  (Expo Router)      │                    │   HTTP :8080         │
│                     │                    │   gRPC :50051        │
│  Clerk Auth SDK     │ ◄── Clerk JWKS ──► │   Auth Middleware    │
└─────────────────────┘                    └──────────┬───────────┘
                                                       │
                              ┌────────────────────────┼────────────────────────┐
                              │                        │                        │
                    ┌─────────▼──────┐    ┌────────────▼──────┐    ┌───────────▼──────┐
                    │   MongoDB      │    │  Google Gemini     │    │  Clerk (SaaS)    │
                    │   (gymlog db)  │    │  2.5 Flash AI      │    │  Auth Provider   │
                    └────────────────┘    └───────────────────┘    └──────────────────┘
```

**Stack:**
- Frontend: React Native, Expo Router, TypeScript, Uber Base Design System
- Backend: Go 1.23, gRPC, HTTP REST gateway
- Database: MongoDB (7 collections)
- Auth: Clerk (JWT/JWKS validation)
- AI: Google Gemini 2.5 Flash

## Quick Start

### Prerequisites

- Node.js 18+
- Go 1.23+
- MongoDB (local or Atlas)
- [Clerk account](https://clerk.com) (free tier)
- [Google AI Studio API key](https://makersuite.google.com/app/apikey) (for AI features)

### 1. Clone and configure

```bash
git clone <repo-url>
cd GymLog
```

### 2. Start the backend

```bash
cd backend

# Create .env
cat > .env << 'EOF'
MONGODB_URI=mongodb://localhost:27017
DB_NAME=gymlog
EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_YOUR_CLERK_KEY
GEMINI_API_KEY=YOUR_GEMINI_KEY
EOF

# Install dependencies and start MongoDB
go mod tidy
mongod &

# Initialize database and populate exercises
mongosh < scripts/init-mongo.js
go run scripts/populate_exercises.go

# Start server (HTTP :8080, gRPC :50051)
go run cmd/server/main.go
```

Or with Docker:

```bash
make docker-compose-up
```

### 3. Start the frontend

```bash
cd frontend-native

# Create .env
cat > .env << 'EOF'
EXPO_PUBLIC_API_BASE_URL=http://localhost:8080
EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_YOUR_CLERK_KEY
EOF

npm install
npm start
# Press i for iOS simulator, a for Android emulator
```

### 4. Get your Clerk key

1. Go to [clerk.com](https://clerk.com) and create an application named "GymLog"
2. Select **Email** as the authentication method
3. Copy the **Publishable Key** from API Keys (starts with `pk_test_`)
4. Use the same key in both `backend/.env` and `frontend-native/.env`

## Development Commands

### Backend

```bash
cd backend
make run              # Start server
make test             # Run tests
make build            # Build binary
make docker-compose-up  # Start with MongoDB via Docker
make populate-exercises # Populate exercise database
```

### Frontend

```bash
cd frontend-native
npm start             # Start Expo dev server
npm run ios           # Run on iOS simulator
npm run android       # Run on Android emulator
npm run build:ios     # EAS cloud build for iOS
npm run build:android # EAS cloud build for Android
```

## Project Documentation

- [Backend README](backend/README.md) — API reference, database schema, metrics system, deployment
- [Frontend README](frontend-native/README.md) — Setup, navigation, authentication, building for production
