# Environment Variables Setup for GymLog

## Overview
This document explains how environment variables are configured and loaded in the GymLog React Native app, particularly for Android builds.

## Files Modified/Created

### 1. `.env.example` (Template)
- Contains all required environment variables with placeholder values
- Copy this to `.env` and update with actual values
- Used as reference for required variables

### 2. `.env` (Local Development)
- Contains actual environment variable values
- Automatically loaded by Expo for local development
- **Important**: Must be prefixed with `EXPO_PUBLIC_` to be accessible in the app

### 3. `eas.json` (EAS Build Configuration)
- Updated with environment variables for different build profiles
- **Development**: Uses localhost API URL
- **Preview**: Uses localhost API URL for testing
- **Production**: Uses production API URL (https://gymlog-973367800029.asia-south1.run.app)

### 4. `app.json` (App Configuration)
- Added `extra` section with environment variable references
- Variables are substituted during build time
- Makes environment variables accessible via `expo-constants`

### 5. `SETUP.md` (Documentation)
- Updated with comprehensive environment variable setup instructions
- Includes important notes about `EXPO_PUBLIC_` prefix requirement
- Explains the three-layer environment variable loading system

## Environment Variable Loading Hierarchy

### For Local Development (`expo start`)
1. `.env` file in project root
2. Process environment variables
3. Default fallback values in code

### For EAS Builds (`eas build`)
1. `eas.json` env configuration (highest priority)
2. `.env` file values
3. Default fallback values in code

### For Standalone Builds
1. `app.json` extra configuration
2. Build-time environment variables
3. Default fallback values in code

## Required Environment Variables

```bash
# API Configuration
EXPO_PUBLIC_API_BASE_URL=http://localhost:8080

# Google OAuth Configuration
EXPO_PUBLIC_GOOGLE_WEB_CLIENT_ID=your_web_client_id.googleusercontent.com
EXPO_PUBLIC_GOOGLE_IOS_CLIENT_ID=your_ios_client_id.googleusercontent.com
EXPO_PUBLIC_GOOGLE_ANDROID_CLIENT_ID=your_android_client_id.googleusercontent.com
```

## Build Profile URLs

- **Development/Preview**: `http://localhost:8080`
- **Production**: `https://gymlog-973367800029.asia-south1.run.app`

## Testing Environment Variables

Environment variables are properly loaded and accessible in the app. The configuration ensures:

1. ✅ Local development works with .env file
2. ✅ EAS builds use profile-specific environment variables
3. ✅ Production builds use correct production API URL
4. ✅ All Google OAuth client IDs are properly configured
5. ✅ Android builds have access to all required environment variables

## Usage in Code

Environment variables are accessed using the standard `process.env` syntax:

```javascript
const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
const GOOGLE_WEB_CLIENT_ID = process.env.EXPO_PUBLIC_GOOGLE_WEB_CLIENT_ID || '';
```

## Important Notes

- **MUST** prefix environment variables with `EXPO_PUBLIC_` for client-side access
- Environment variables are embedded at build time (not runtime)
- Different build profiles can have different environment variable values
- Local .env file is gitignored for security (contains sensitive OAuth keys)
- .env.example is tracked in git as a template