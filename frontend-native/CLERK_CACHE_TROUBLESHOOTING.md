# Clerk Cache Troubleshooting Guide

## Problem: Wrong Clerk Instance Tokens Being Sent

### Symptoms:
- Changed `EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY` to production instance
- App still sends bearer tokens from dev environment
- Backend receives tokens from wrong Clerk instance

### Root Cause:
**Clerk tokens are cached in SecureStore and persist across app restarts.**

When you change the publishable key in `.env`, the app will use the new key for the Clerk UI and authentication flow, but **cached tokens from the old instance remain in SecureStore** and continue to be sent to your backend.

## Solutions (Choose One):

### ✅ Solution 1: Clear App Data (Recommended)

#### iOS:
1. Delete the app from your device/simulator
2. Reinstall: `npx expo run:ios` or `eas build --profile preview --platform ios`

#### Android:
1. Go to Settings → Apps → RepTracker
2. Tap "Storage"
3. Tap "Clear Data" (not just "Clear Cache")
4. Reopen the app

### ✅ Solution 2: Programmatic Cache Clearing

Add a debug button or call this on app startup when switching environments:

```typescript
import { clearClerkCache } from '@/app/utils/tokenCache';

// In your component or debug screen
const handleClearCache = async () => {
  await clearClerkCache();
  // Force user to sign in again
  await signOut();
};
```

### ✅ Solution 3: Use Separate Build Profiles

Configure `eas.json` with different Clerk keys per environment:

```json
{
  "build": {
    "development": {
      "env": {
        "EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY": "pk_test_YOUR_DEV_KEY"
      }
    },
    "production": {
      "env": {
        "EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY": "pk_live_YOUR_PROD_KEY"
      }
    }
  }
}
```

Then build with the correct profile:
```bash
# Development build
eas build --profile development --platform ios

# Production build
eas build --profile production --platform ios
```

## Prevention:

### 1. Always Use Build Profiles
Never rely on `.env` alone for environment switching. Use `eas.json` build profiles.

### 2. Separate Clerk Instances
- **Development**: Use Clerk test instance (`pk_test_...`)
- **Production**: Use Clerk production instance (`pk_live_...`)

### 3. Clear Cache When Switching
When switching between environments during development, always clear app data.

## Verification:

After clearing cache, verify the correct tokens are being used:

1. **Check the publishable key**:
   ```typescript
   console.log('Clerk Key:', process.env.EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY);
   ```

2. **Check the token**:
   ```typescript
   const token = await getToken();
   console.log('Token:', token?.substring(0, 50) + '...');
   ```

3. **Verify on backend**:
   - Check backend logs for JWT validation
   - Ensure backend is configured with correct Clerk instance keys

## Common Mistakes:

❌ **Changing `.env` and expecting immediate effect**
- Cached tokens persist in SecureStore

❌ **Not configuring `eas.json` with Clerk keys**
- Build profiles should include all environment variables

❌ **Using same Clerk instance for dev and production**
- Use separate instances to avoid confusion

❌ **Only clearing app cache, not app data**
- Must clear app data (iOS: delete app, Android: clear data)

## Quick Fix Checklist:

- [ ] Update `eas.json` with correct `EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY` for each profile
- [ ] Clear app data (iOS: delete app, Android: clear data in settings)
- [ ] Rebuild app with correct profile: `eas build --profile production`
- [ ] Sign in again with fresh authentication
- [ ] Verify correct tokens in backend logs

## Additional Resources:

- [Clerk Token Cache Documentation](https://clerk.com/docs/references/react-native/clerk-provider#token-cache)
- [Expo SecureStore Documentation](https://docs.expo.dev/versions/latest/sdk/securestore/)
- [EAS Build Configuration](https://docs.expo.dev/build/eas-json/)

