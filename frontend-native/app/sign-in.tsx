import React from 'react';
import { router } from 'expo-router';
import { SignInScreen } from '../src/screens';

export default function SignInRoute() {
  const handleSignInSuccess = () => {
    // Navigation will be handled automatically by the auth state change
    // The index.tsx will redirect to /(tabs) when authenticated
    router.replace('/');
  };

  return <SignInScreen onSignInSuccess={handleSignInSuccess} />;
} 