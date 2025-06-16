# GymLog - Workout Routine Manager

A full-stack application for creating and managing workout routines, built with Next.js frontend and Go gRPC backend.

## Features

- ✅ **Create Routine Flow**: Navigate from home page to create workout routines
- ✅ **Exercise Search**: Browse and search through a library of exercises
- ✅ **Exercise Details**: View exercise descriptions, muscle groups, equipment, and instructions  
- ✅ **Routine Builder**: Add exercises with custom sets, reps, and weights
- ✅ **Modern UI**: Clean design using Airtable design system components

## Architecture

- **Frontend**: Next.js 14 with TypeScript and Tailwind CSS
- **Backend**: Go with gRPC API
- **Design System**: Custom Airtable-inspired components
- **API**: RESTful HTTP endpoints bridging to gRPC services

## Getting Started

### Prerequisites

- Node.js 18+
- Go 1.23+

### Backend Setup

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Start the gRPC server:
```bash
go run cmd/server/main.go
```

The backend will start on port `50051`.

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`.

## Usage Flow

1. **Home Page** (`/home`): Start here and click "Create Routine"
2. **Create Routine** (`/routines/create`): Name your routine and click "Add Exercise"
3. **Exercise Search** (`/exercises/search`): Search and select exercises, configure sets/reps
4. **Add to Routine**: Exercises are added back to the routine builder

## API Endpoints

### Frontend API Routes
- `GET /api/exercises` - List exercises with optional search/filtering

### Backend gRPC Services  
- `ExerciseService` - CRUD operations for exercises
- `UserService` - User management (planned)
- `WorkoutService` - Workout tracking (planned)

## Development Notes

- Currently using mock data for exercises
- Backend provides sample exercises through gRPC
- Frontend uses HTTP API routes that bridge to gRPC
- Ready for database integration

## Next Steps

- [ ] Database integration (PostgreSQL)
- [ ] User authentication
- [ ] Workout tracking and history
- [ ] Exercise video/image assets
- [ ] Mobile responsive optimizations
