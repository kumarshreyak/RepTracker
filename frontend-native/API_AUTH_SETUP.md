# API Authentication Setup

## Overview
This document describes the centralized API authentication setup for the GymLog React Native app using Clerk tokens.

## Problem Solved
Previously, API calls were made using direct `fetch()` calls without proper authentication headers, resulting in "401 Unauthorized" errors. Each API call would need to manually handle token retrieval and header setup.

## Solution
Implemented a centralized authentication system that automatically attaches Clerk bearer tokens to all API requests.

## Architecture

### 1. Global Token Getter (`src/utils/api.ts`)
The API utility now supports a global token getter function that is automatically called for all requests:

```typescript
// Set once at app initialization
setGlobalTokenGetter(async () => {
  return await getToken(); // Clerk's getToken function
});

// All subsequent API calls automatically include the bearer token
await apiGet('/api/workouts');
await apiPost('/api/workout-sessions', data);
```

### 2. App Initialization (`app/_layout.tsx`)
The global token getter is initialized in the root layout component:

```typescript
function RootLayoutNav() {
  const { getToken } = useAuth();

  useEffect(() => {
    setGlobalTokenGetter(async () => {
      try {
        return await getToken();
      } catch (error) {
        console.error('Error getting token:', error);
        return null;
      }
    });
  }, [getToken]);

  return <Stack>...</Stack>;
}
```

### 3. API Request Flow
1. Component makes API call: `apiGet('/api/workouts')`
2. API utility checks for token in options parameter
3. If no token provided, calls global token getter
4. Adds `Authorization: Bearer <token>` header
5. Makes authenticated request

## Updated Files

### Core API Files
- `src/utils/api.ts` - Added global token getter support
- `app/_layout.tsx` - Initialize global token getter with Clerk

### Screen Files Updated
All screens now use the centralized API utilities instead of direct fetch calls:

- `app/(tabs)/index.tsx` - Home tab (routines & past workouts)
- `app/(tabs)/insights.tsx` - Insights tab (AI analyses)
- `app/active-workout.tsx` - Active workout session
- `app/create-routine.tsx` - Create/edit routine
- `app/exercise-search.tsx` - Exercise search modal
- `app/workout-detail.tsx` - Workout detail view

### Changes Made
**Before:**
```typescript
const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL;
const response = await fetch(`${API_BASE_URL}/api/workouts`, {
  headers: {
    'Content-Type': 'application/json',
  },
});
```

**After:**
```typescript
import { apiGet } from '../src/utils/api';

const data = await apiGet<{ workouts: Routine[] }>('/api/workouts');
```

## Benefits

1. **Automatic Authentication**: All API calls automatically include bearer tokens
2. **Centralized Token Management**: Token retrieval logic in one place
3. **Cleaner Code**: No need to manually handle tokens in every component
4. **Type Safety**: TypeScript support for request/response types
5. **Error Handling**: Consistent error handling across all API calls
6. **Maintainability**: Easy to update authentication logic globally

## API Utility Functions

### Available Functions
- `apiGet<T>(endpoint, token?, options?)` - GET requests
- `apiPost<T>(endpoint, data, token?, options?)` - POST requests
- `apiPut<T>(endpoint, data, token?, options?)` - PUT requests
- `apiDelete<T>(endpoint, token?, options?)` - DELETE requests

### Usage Examples

```typescript
// GET request
const workouts = await apiGet<{ workouts: Routine[] }>('/api/workouts');

// POST request
const session = await apiPost('/api/workout-sessions', {
  userId: user.id,
  routineId: routineId,
  name: 'My Workout',
});

// PUT request
await apiPut(`/api/workout-sessions/${sessionId}`, updateData);

// DELETE request
await apiDelete(`/api/workouts/${routineId}`);
```

## Token Override
You can still manually provide a token if needed:

```typescript
const customToken = await getCustomToken();
await apiGet('/api/workouts', customToken);
```

## Environment Variables
Ensure `EXPO_PUBLIC_API_BASE_URL` is set in your `.env` file:

```
EXPO_PUBLIC_API_BASE_URL=http://192.168.1.4:8080
```

## Testing
1. Start the backend server
2. Launch the React Native app
3. Sign in with Clerk
4. All API calls should now include proper authentication headers
5. Check backend logs to verify `Authorization: Bearer <token>` headers are present

## Troubleshooting

### "Missing authorization header" error
- Ensure you've signed in with Clerk
- Check that `setGlobalTokenGetter` is called in `_layout.tsx`
- Verify Clerk is properly configured

### Token not being sent
- Check that you're using the API utility functions (`apiGet`, `apiPost`, etc.)
- Verify the global token getter is initialized before making API calls
- Check console logs for token retrieval errors

## Future Improvements
- Add token refresh logic
- Implement request retry on 401 errors
- Add request/response interceptors for logging
- Cache tokens to reduce Clerk API calls

