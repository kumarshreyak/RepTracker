# Expo Router Navigation Setup

This app now uses [Expo Router](https://docs.expo.dev/tutorial/add-navigation/) for file-based navigation, providing a seamless navigation experience across iOS, Android, and web.

## Navigation Structure

```
app/
├── _layout.tsx          # Root layout with Stack navigator
├── index.tsx            # Entry point with authentication routing
├── sign-in.tsx          # Sign-in screen (unauthenticated)
├── create-routine.tsx   # Create routine screen (authenticated)
├── +not-found.tsx       # 404 error page
└── (tabs)/              # Tab navigator group
    ├── _layout.tsx      # Tab layout configuration
    ├── index.tsx        # Home tab (/)
    └── profile.tsx      # Profile tab (/profile)
```

## Key Features

### 🔐 Authentication Flow
- **Entry Point**: `app/index.tsx` handles initial routing based on auth state
- **Unauthenticated users**: Redirected to `/sign-in` 
- **Authenticated users**: Redirected to `/(tabs)` (tab navigator)
- **Route Protection**: Individual routes check auth state and redirect to `/` if needed
- **Centralized Routing**: All navigation goes through `/` which redirects appropriately

### 📱 Tab Navigation
- **Home Tab** (`/`): Main dashboard with routines and quick actions
- **Profile Tab** (`/profile`): User profile and account settings
- Custom icons using Ionicons from `@expo/vector-icons`
- Airtable design system styling

### 🧭 Navigation Methods

#### Using router.push()
```typescript
import { router } from 'expo-router';

// Navigate to create routine
router.push('/create-routine');

// Navigate back
router.back();

// Replace current route (for sign out)
router.replace('/sign-in');
```

#### Using Link component
```typescript
import { Link } from 'expo-router';

<Link href="/create-routine" asChild>
  <Button>Create Routine</Button>
</Link>
```

### 🎨 Design System Integration
- All screens use the Airtable design system components
- Consistent styling with `getColor()` utility
- Typography and Button components throughout

## Adding New Routes

### 1. Create a new file in the app directory
```typescript
// app/new-screen.tsx
import React from 'react';
import { View } from 'react-native';
import { Typography } from '../src/components';

export default function NewScreen() {
  return (
    <View>
      <Typography variant="heading-large">New Screen</Typography>
    </View>
  );
}
```

### 2. Navigate to the new route
```typescript
router.push('/new-screen');
```

### 3. Add to tab navigator (if needed)
Edit `app/(tabs)/_layout.tsx` to add a new tab:

```typescript
<Tabs.Screen
  name="new-tab"
  options={{
    title: 'New Tab',
    tabBarIcon: ({ color, focused }) => (
      <Ionicons name={focused ? 'icon-filled' : 'icon-outline'} color={color} size={24} />
    ),
  }}
/>
```

## Protected Routes

All routes except `/sign-in` are automatically protected by the authentication check in the root layout. If a user is not authenticated, they'll be redirected to the sign-in screen.

## Error Handling

The `+not-found.tsx` route handles any invalid URLs and provides a way to navigate back to the home screen.

## Dependencies

- `expo-router`: File-based routing system
- `@expo/vector-icons`: Icon set for tab navigation
- `react-native-gesture-handler`: Required for navigation gestures
- `react-native-safe-area-context`: Safe area handling
- `react-native-screens`: Native screen optimization

## Migration from React Navigation

This app was migrated from React Navigation to Expo Router, removing the following dependencies:
- `@react-navigation/native`
- `@react-navigation/bottom-tabs` 
- `@react-navigation/stack`

The new system provides better TypeScript support, simpler navigation patterns, and automatic deep linking support.

## Troubleshooting

### "No filename found" Error
This error was resolved by:

1. **Setting correct entry point** in `package.json`:
   ```json
   "main": "expo-router/entry"
   ```

2. **Creating proper file structure** with `app/index.tsx` as the main entry point

3. **Avoiding conditional Stack.Screen rendering** in `_layout.tsx` - let file structure define routes

4. **Adding metro.config.js** with default Expo configuration

5. **Centralizing auth routing** in `app/index.tsx` instead of `_layout.tsx`

### Key Configuration Files
- `package.json`: Entry point must be `"expo-router/entry"`
- `metro.config.js`: Required for proper module resolution
- `app.json`: Include scheme for deep linking
- `app/index.tsx`: Main entry point handling auth-based routing 