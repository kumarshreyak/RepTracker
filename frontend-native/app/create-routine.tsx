import React from 'react';
import { router } from 'expo-router';
import { CreateRoutineScreen } from '../src/screens';
import { useAuth } from '../src/hooks/useAuth';

export default function CreateRoutineRoute() {
  const { user } = useAuth();

  const handleNavigateBack = () => {
    router.back();
  };

  const handleNavigateToExerciseSearch = () => {
    router.push('/exercise-search');
  };

  if (!user) {
    router.replace('/');
    return null;
  }

  return (
    <CreateRoutineScreen
      user={user}
      onNavigateBack={handleNavigateBack}
      onNavigateToExerciseSearch={handleNavigateToExerciseSearch}
    />
  );
} 