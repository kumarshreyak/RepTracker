# GymLog Frontend

A Next.js frontend application for the GymLog workout tracking platform with Google OAuth authentication and design system integration.

## Setup

### Prerequisites

- Node.js 18+
- Google OAuth credentials
- Backend server running

### Environment Variables

Create a `.env.local` file in the frontend directory:

```env
# NextAuth Configuration
NEXTAUTH_SECRET=your-nextauth-secret-key
NEXTAUTH_URL=http://localhost:3000

# Google OAuth
AUTH_GOOGLE_ID=your-google-oauth-client-id
AUTH_GOOGLE_SECRET=your-google-oauth-client-secret

# Backend Configuration  
BACKEND_HTTP_URL=http://localhost:8080
GRPC_SERVER_URL=localhost:50051
```

### Installation

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

The application will be available at `http://localhost:3000`.

## Authentication Flow

### Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URIs:
   - `http://localhost:3000/api/auth/callback/google` (development)
   - `https://yourdomain.com/api/auth/callback/google` (production)

### User Authentication Process

1. **Sign In**: User clicks "Continue with Google" on `/auth/signin`
2. **OAuth Flow**: NextAuth handles Google OAuth redirect
3. **Backend Integration**: 
   - Frontend calls `/api/auth/google` with user data
   - Backend creates/updates user in MongoDB
   - Backend creates session and returns session token
4. **Session Storage**: Frontend stores session token in localStorage
5. **Session Validation**: 
   - `useAuth` hook validates session with backend
   - Session token included in API requests
6. **Session Expiry**: 
   - Backend validates session on each request
   - Expired sessions trigger automatic logout and redirect

### Protected Routes

All routes except `/auth/signin` are protected by the auth middleware. Users without valid sessions are redirected to the sign-in page.

## Components

### Design System

The application uses the Airtable design system components:

- `Typography` - For all text elements with consistent styling
- `Button` - For interactive elements with multiple variants
- `ColorSwatch` - For color selections

### Authentication Hook

```typescript
import { useAuth } from '@/hooks/useAuth'

function MyComponent() {
  const { 
    session,
    backendUser,
    sessionToken,
    isAuthenticated,
    isLoading,
    logout,
    validateSession 
  } = useAuth()

  if (isLoading) return <div>Loading...</div>
  if (!isAuthenticated) return <div>Please sign in</div>

  return <div>Welcome, {backendUser.name}!</div>
}
```

## API Routes

### Authentication APIs

- `GET /api/auth/[...nextauth]` - NextAuth handlers
- `POST /api/auth/google` - Google OAuth user creation/login
- `GET /api/auth/validate` - Session validation
- `POST /api/auth/logout` - Session cleanup

## Project Structure

```
src/
├── app/
│   ├── api/auth/          # Authentication API routes
│   ├── auth/signin/       # Sign-in page
│   ├── home/              # Protected home page
│   └── layout.tsx         # Root layout
├── components/            # Design system components
├── hooks/
│   └── useAuth.ts         # Authentication hook
├── types/
│   └── next-auth.d.ts     # NextAuth type extensions
├── auth.ts                # NextAuth configuration
└── middleware.ts          # Route protection
```

## Session Management

- **Storage**: Session tokens stored in localStorage
- **Validation**: Automatic validation on app load and route changes
- **Expiry**: 24-hour session duration
- **Cleanup**: Automatic cleanup on logout or expiry
- **Redirect**: Expired sessions redirect to `/auth/signin`

## Development

```bash
# Start development server
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Lint code
npm run lint
```

## Integration with Backend

The frontend communicates with the Go backend via:

1. **HTTP APIs** for authentication and session management
2. **gRPC services** for workout data (exercises, routines, workouts)

Session tokens from authentication are included in all backend requests for authorization.
