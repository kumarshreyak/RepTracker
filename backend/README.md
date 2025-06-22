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

### HTTP Endpoints

- `POST /api/auth/google` - Authenticate with Google OAuth
- `GET /api/auth/validate` - Validate session token
- `POST /api/auth/logout` - Logout and delete session

### gRPC Services

- `UserService` - User management
- `ExerciseService` - Exercise library
- `WorkoutService` - Workout tracking
- `WorkoutSessionService` - Active workout session management
- `MetricsService` - Workout analytics and volume tracking

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

## Documentation

- **[METRICS.md](METRICS.md)** - Comprehensive guide to the workout metrics system, including all calculated metrics, formulas, API endpoints, and usage examples

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