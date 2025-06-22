# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GymLog is a full-stack workout routine manager with React Native frontend (Expo) and Go gRPC backend. Features include exercise search, routine creation, workout tracking, and comprehensive exercise database.

## Architecture

- **Frontend**: React Native with Expo Router (file-based routing)
- **Backend**: Go with gRPC API, MongoDB database
- **Design System**: Uber Base Design System with semantic colors and typography
- **Authentication**: Google Sign-In integration
- **Database**: MongoDB for exercises, users, and workout data

## Development Commands

### Frontend (React Native/Expo)
```bash
cd frontend-native
npm install                    # Install dependencies
npm start                     # Start Expo development server
npm run ios                   # Run on iOS simulator
npm run android               # Run on Android emulator
npm run build:ios             # Build for iOS
npm run build:android         # Build for Android
```

### Backend (Go)
```bash
cd backend
make setup                    # Setup development environment
make run                      # Start the server (port 50051)
make test                     # Run tests
make build                    # Build binary
make docker-compose-up        # Start with MongoDB via Docker
make populate-exercises       # Populate exercise database
```

## Design System Requirements

### Uber Base Semantic Colors
- **MANDATORY**: Always use semantic colors from `getColor()` function
- **Page backgrounds**: `backgroundSecondary` (gray 50)
- **Card backgrounds**: `backgroundPrimary` (white)
- **Card borders**: `borderOpaque` (gray 100)
- **Primary text**: `contentPrimary` (black)
- **Secondary text**: `contentSecondary` (gray 800)
- **Accent actions**: `backgroundAccent` with `contentOnColor`
- **Never use**: Legacy colors, hex codes, or primitive colors directly

### Typography System
- **Use Typography component only**: Never use HTML elements like `<h1>`, `<p>`, `<span>`
- **Four categories**: Display, Headings, Labels, Paragraphs
- **Page titles**: `heading-xlarge` + `contentPrimary`
- **Section headers**: `heading-small` + `contentPrimary`
- **Body text**: `paragraph-medium` + `contentSecondary`
- **Button text**: Handled by Button component
- **Form labels**: `label-small` + `contentPrimary`

### Layout Patterns
- Use section-based layout on `backgroundSecondary` backgrounds
- Individual item cards: `backgroundPrimary` with `borderOpaque` borders
- Section spacing: `paddingHorizontal: 16, marginBottom: 24`
- Button borders: `borderRadius: 3` (Base design system style)

## Data Model Conventions

### Naming Convention (CRITICAL)
- **Always use camelCase**: `muscleGroup`, `exerciseId`, `workoutSession`
- **Never use snake_case**: No `muscle_group`, `exercise_id`, `workout_session`
- **Backend structs**: PascalCase fields with camelCase JSON tags
- **API responses**: All JSON uses camelCase

### Exercise Constants (Backend)
- **MANDATORY**: Use typed constants from `pkg/models/exercise_constants.go`
- **Categories**: `models.CategoryStrength`, `models.CategoryCardio`
- **Equipment**: `models.EquipmentBarbell`, `models.EquipmentDumbbell`
- **Muscles**: `models.MuscleBiceps`, `models.MuscleChest`
- **Always validate**: Use `IsValid*` functions for user input

## Project Structure

### Frontend (`frontend-native/`)
```
app/                     # Expo Router pages
├── (auth)/             # Authentication pages
├── (tabs)/             # Tab navigation pages
├── active-workout.tsx  # Active workout screen
├── create-routine.tsx  # Routine creation
└── exercise-search.tsx # Exercise search
src/
├── components/         # Reusable UI components
├── auth/              # Authentication service
├── hooks/             # Custom React hooks
└── types/             # TypeScript type definitions
```

### Backend (`backend/`)
```
cmd/server/            # Main application entry
internal/services/     # Business logic services
pkg/
├── api/v1/           # API definitions
├── models/           # Data models and constants
├── database/         # MongoDB connection
└── middleware/       # Authentication middleware
proto/gymlog/v1/      # Protocol Buffer definitions
scripts/              # Database population scripts
```

## Key Files and Constants

- `frontend-native/src/components/Colors.ts`: Semantic color system implementation
- `frontend-native/COLOR_SYSTEM.md`: Complete color system documentation
- `frontend-native/TYPOGRAPHY_SYSTEM.md`: Typography hierarchy guide
- `backend/pkg/models/exercise_constants.go`: Exercise type definitions
- `backend/pkg/models/EXERCISE_CONSTANTS_README.md`: Exercise constants guide
- `.cursorrules`: Project-specific development guidelines

## Testing

- **Frontend**: No specific test framework configured yet
- **Backend**: Go testing with `make test`
- **Integration**: Docker Compose for MongoDB testing

## Important Development Notes

1. **Always check documentation**: COLOR_SYSTEM.md and TYPOGRAPHY_SYSTEM.md before UI work
2. **Use ColorDemo component**: Reference available semantic colors
3. **Validate exercise data**: Use backend validation functions for categories/equipment
4. **Follow authentication flow**: Google Sign-In integration already configured
5. **MongoDB collections**: Use lowercase plural names (exercises, users, workouts)
6. **Error handling**: Return meaningful errors with context
7. **Context propagation**: Always pass `context.Context` in Go services