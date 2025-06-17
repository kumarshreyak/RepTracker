import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import React from 'react';

export default function RootLayout() {
  return (
    <>
      <Stack screenOptions={{ headerShown: false }}>
        <Stack.Screen name="index" options={{ headerShown: false }} />
        <Stack.Screen name="(auth)" options={{ headerShown: false }} />
        <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
        <Stack.Screen 
          name="create-routine" 
          options={{ 
            headerShown: true,
            presentation: 'modal',
            title: 'Create Routine'
          }} 
        />
        <Stack.Screen 
          name="exercise-search" 
          options={{ 
            headerShown: true,
            presentation: 'modal',
            title: 'Search Exercises'
          }} 
        />
        <Stack.Screen name="+not-found" />
      </Stack>
      <StatusBar style="auto" />
    </>
  );
} 