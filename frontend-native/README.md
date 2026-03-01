# GymLog Frontend

React Native (Expo) iOS/Android app for GymLog. Uses Expo Router for file-based navigation and Clerk for authentication.

## Table of Contents

- [Setup](#setup)
- [Running the App](#running-the-app)
- [Navigation Structure](#navigation-structure)
- [Authentication](#authentication)
- [API Utility](#api-utility)
- [Design System](#design-system)
- [Building for Production](#building-for-production)
- [Troubleshooting](#troubleshooting)

---

## Setup

### Prerequisites

- Node.js 18+
- Expo CLI (`npm install -g expo-cli`)
- iOS: macOS with Xcode
- Android: Android Studio with an emulator
- [Clerk account](https://clerk.com) (free tier)
- Backend server running (see [backend README](../backend/README.md))

### 1. Install dependencies

```bash
cd frontend-native
npm install
```

### 2. Create environment file

```bash
cat > .env << 'EOF'
EXPO_PUBLIC_API_BASE_URL=http://localhost:8080
EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_YOUR_KEY_HERE
EOF
```

**For Android emulator**, use `http://10.0.2.2:8080` instead of `localhost`.

### 3. Get your Clerk key

1. Go to [clerk.com](https://clerk.com) and create an application named "GymLog"
2. Select **Email** as the authentication method
3. Copy the **Publishable Key** from API Keys (starts with `pk_test_`)
4. Use the same key in both `frontend-native/.env` and `backend/.env`

### Environment Variables Reference

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `EXPO_PUBLIC_API_BASE_URL` | Yes | Backend API base URL | `http://localhost:8080` |
| `EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY` | Yes | Clerk publishable key | `pk_test_...` |

All variables must be prefixed with `EXPO_PUBLIC_` to be accessible in the app. Variables are embedded at build time, not runtime.

---

## Running the App

```bash
npm start           # Start Expo dev server
npm run ios         # Run on iOS simulator
npm run android     # Run on Android emulator
```

After `npm start`, press:
- `i` — Open iOS simulator
- `a` — Open Android emulator
- Scan QR code with Expo Go on a physical device

**Note:** Google Sign-in and some native features require a development build, not Expo Go.

---

## Navigation Structure

The app uses [Expo Router](https://docs.expo.dev/router/introduction/) for file-based routing.

```
app/
├── _layout.tsx              # Root layout — initializes Clerk and global token getter
├── index.tsx                # Entry point — routes to auth or tabs based on auth state
├── (auth)/
│   ├── _layout.tsx          # Auth stack layout
│   └── sign-in.tsx          # Sign-in/sign-up screen
├── (tabs)/
│   ├── _layout.tsx          # Tab bar configuration
│   ├── index.tsx            # Home tab — routines and past workouts
│   ├── insights.tsx         # Insights tab — AI analysis
│   └── profile.tsx          # Profile tab — user settings
├── active-workout.tsx       # Active workout session screen
├── ai-analysis-detail.tsx   # AI analysis detail view
├── create-routine.tsx       # Create/edit routine modal
├── exercise-search.tsx      # Exercise search modal
├── onboarding.tsx           # User onboarding (height, weight, goal)
├── suggestion-detail.tsx    # Workout suggestion detail
├── workout-detail.tsx       # Past workout detail view
└── +not-found.tsx           # 404 screen
```

### Authentication Flow

- Unauthenticated users → `/(auth)/sign-in`
- Authenticated users without profile → `/onboarding`
- Authenticated users with profile → `/(tabs)`

### Navigation Methods

```typescript
import { router } from 'expo-router';

router.push('/create-routine');   // Navigate forward
router.back();                    // Navigate back
router.replace('/(auth)/sign-in'); // Replace (used for sign-out)
```

### Adding New Routes

Create a new file in `app/`:

```typescript
// app/my-screen.tsx
import { View } from 'react-native';
import { Typography } from '@/components';

export default function MyScreen() {
  return (
    <View>
      <Typography variant="heading-large" color="contentPrimary">My Screen</Typography>
    </View>
  );
}
```

Then navigate to it: `router.push('/my-screen')`

---

## Authentication

Authentication is handled by [Clerk](https://clerk.com). The app uses email/password with email verification.

### How It Works

1. User signs in/up via Clerk's built-in UI components
2. Clerk issues a short-lived JWT
3. The app stores the token in SecureStore
4. All API requests include `Authorization: Bearer <token>`
5. Backend validates the token against Clerk's JWKS endpoint

### Auth Hook

```typescript
import { useAuth } from '@/hooks/useAuth';

const { user, isSignedIn, getToken, signOut } = useAuth();
```

### Protected Routes

All routes except `/(auth)/sign-in` require authentication. The root `app/index.tsx` handles redirection based on auth state.

---

## API Utility

All API calls go through `src/utils/api.ts`, which automatically attaches the Clerk bearer token to every request.

### Setup (already done in `app/_layout.tsx`)

```typescript
import { setGlobalTokenGetter } from '@/utils/api';
import { useAuth } from '@clerk/clerk-expo';

const { getToken } = useAuth();
setGlobalTokenGetter(async () => await getToken());
```

### Usage

```typescript
import { apiGet, apiPost, apiPut, apiDelete } from '@/utils/api';

// GET
const data = await apiGet<{ workouts: Routine[] }>('/api/workouts');

// POST
const session = await apiPost('/api/workout-sessions', { routineId, name });

// PUT
await apiPut(`/api/workout-sessions/${sessionId}`, updateData);

// DELETE
await apiDelete(`/api/workouts/${routineId}`);
```

The base URL is read from `EXPO_PUBLIC_API_BASE_URL`.

---

## Design System

The app uses the **Uber Base Design System** with semantic colors and typography. See the dedicated documentation files:

- [`COLOR_SYSTEM.md`](COLOR_SYSTEM.md) — Semantic color system, three-layer architecture, usage patterns
- [`TYPOGRAPHY_SYSTEM.md`](TYPOGRAPHY_SYSTEM.md) — Typography hierarchy, variants, usage rules

### Quick Reference

```typescript
import { Typography, Button, getColor } from '@/components';

// Colors
getColor('backgroundPrimary')   // white — card backgrounds
getColor('backgroundSecondary') // gray 50 — page backgrounds
getColor('contentPrimary')      // black — primary text
getColor('contentSecondary')    // gray 800 — secondary text
getColor('borderOpaque')        // gray 100 — card borders
getColor('backgroundAccent')    // blue — accent actions

// Typography
<Typography variant="heading-large" color="contentPrimary">Title</Typography>
<Typography variant="label-medium" color="contentPrimary">Card title</Typography>
<Typography variant="label-small" color="contentSecondary">Subtitle</Typography>
<Typography variant="label-xsmall" color="contentTertiary">Caption</Typography>

// Buttons
<Button variant="primary" size="default" onPress={handlePress}>Save</Button>
<Button variant="danger" size="small" onPress={handleDelete}>Delete</Button>
```

---

## Building for Production

### EAS Cloud Builds (Recommended)

[Expo Application Services (EAS)](https://expo.dev/eas) handles cloud builds.

```bash
# iOS
eas build --platform ios

# Android
eas build --platform android

# Both platforms
eas build --profile production
```

Build profiles are configured in `eas.json`. The production profile uses the production backend URL (`https://gymlog-973367800029.asia-south1.run.app`).

**For production builds**, use a production Clerk key (`pk_live_...`) in `eas.json`:

```json
{
  "build": {
    "production": {
      "env": {
        "EXPO_PUBLIC_API_BASE_URL": "https://gymlog-973367800029.asia-south1.run.app",
        "EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY": "pk_live_YOUR_PROD_KEY"
      }
    }
  }
}
```

### Local iOS Build (Development)

```bash
# Release build (no Metro bundler required)
npm run ios:release

# Or manually
xcodebuild -workspace ios/RepTracker.xcworkspace \
           -scheme RepTracker \
           -configuration Release \
           -destination generic/platform=iOS \
           -archivePath ios/build/RepTracker.xcarchive \
           archive

xcodebuild -exportArchive \
           -archivePath ios/build/RepTracker.xcarchive \
           -exportPath ios/build/export \
           -exportOptionsPlist ios/RepTracker/ExportOptions.plist
```

Install the resulting `ios/build/export/RepTracker.ipa` via Xcode (Window > Devices and Simulators) or Finder.

### Local Android Build

```bash
npm run android:release
# APK output: android/app/build/outputs/apk/release/
```

---

## Troubleshooting

### Clerk token cache issues (wrong environment tokens)

**Symptom:** Changed `EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY` but app still sends tokens from the old Clerk instance.

**Cause:** Clerk tokens are cached in SecureStore and persist across app restarts.

**Fix:**
- **iOS:** Delete and reinstall the app
- **Android:** Settings → Apps → RepTracker → Storage → Clear Data (not just cache)
- **Programmatic:** Call `clearClerkCache()` from `app/utils/tokenCache.ts` before signing out

**Prevention:** Use separate `eas.json` build profiles with the correct key per environment. Never rely on `.env` alone when switching between dev and production Clerk instances.

### "Missing EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY"

- Verify `.env` exists in `frontend-native/` (not the project root)
- Confirm the key starts with `pk_test_` (dev) or `pk_live_` (production)
- Restart the Expo dev server after creating or modifying `.env`

### "Network request failed"

- Confirm the backend server is running on port 8080
- For Android emulator, use `EXPO_PUBLIC_API_BASE_URL=http://10.0.2.2:8080`
- For physical device, use your machine's local IP address (e.g., `http://192.168.1.x:8080`)

### "Invalid or expired token"

- Sign out and sign back in to get a fresh token
- Verify both frontend and backend use the same Clerk publishable key
- Check that the device has internet access (Clerk validates tokens remotely)

### "User not found" after sign-in

- The user record is created automatically on the first authenticated API call
- Check backend logs for MongoDB connection errors
- Verify `MONGODB_URI` is correct in `backend/.env`

### Build fails with Pod errors

```bash
cd ios
pod deintegrate
pod install
cd ..
```

### "No filename found" (Expo Router)

- Verify `package.json` has `"main": "expo-router/entry"`
- Ensure `app/index.tsx` exists as the entry point
- Check `metro.config.js` is present with default Expo configuration
