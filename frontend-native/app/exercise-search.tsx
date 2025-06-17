import React from 'react';
import { router } from 'expo-router';
import { ExerciseSearchScreen } from '../src/screens';
import { useAuth } from '../src/hooks/useAuth';

export default function ExerciseSearchRoute() {
  const { user } = useAuth();

  const handleNavigateBack = () => {
    router.back();
  };

  if (!user) {
    router.replace('/');
    return null;
  }

  return (
    <ExerciseSearchScreen
      user={user}
      onNavigateBack={handleNavigateBack}
    />
  );
} 