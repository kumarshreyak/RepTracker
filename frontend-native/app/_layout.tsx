import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import React, { useEffect } from 'react';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { ClerkProvider, useAuth } from '@clerk/clerk-expo';
import { tokenCache } from './utils/tokenCache';
import { setGlobalTokenGetter } from '../src/utils/api';

const publishableKey = process.env.EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY!;

if (!publishableKey) {
  throw new Error('Missing EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY in environment variables');
}

function RootLayoutNav() {
  const { getToken } = useAuth();

  useEffect(() => {
    // Set up global token getter for all API requests
    setGlobalTokenGetter(async () => {
      try {
        return await getToken();
      } catch (error) {
        console.error('Error getting token:', error);
        return null;
      }
    });
  }, [getToken]);

  return (
    <Stack screenOptions={{ headerShown: false }}>
      <Stack.Screen name="index" options={{ headerShown: false }} />
      <Stack.Screen name="(auth)" options={{ headerShown: false }} />
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
      <Stack.Screen 
        name="create-routine" 
        options={{ 
          headerShown: false
        }} 
      />
      <Stack.Screen 
        name="exercise-search" 
        options={{ 
          headerShown: false,
          presentation: 'modal'
        }} 
      />
      <Stack.Screen name="+not-found" />
    </Stack>
  );
}

export default function RootLayout() {
  return (
    <ClerkProvider publishableKey={publishableKey} tokenCache={tokenCache}>
      <SafeAreaProvider>
        <RootLayoutNav />
        <StatusBar style="auto" />
      </SafeAreaProvider>
    </ClerkProvider>
  );
} 