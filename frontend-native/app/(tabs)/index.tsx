import React from 'react';
import { router } from 'expo-router';
import { HomeScreen } from '../../src/screens';
import { useAuth } from '../../src/hooks/useAuth';
import { authService } from '../../src/auth/AuthService';

export default function HomeTab() {
  const { user } = useAuth();

  const handleSignOut = async () => {
    await authService.signOut();
    router.replace('/');
  };

  const handleNavigateToCreateRoutine = () => {
    router.push('/create-routine');
  };

  const handleNavigateToExerciseLibrary = () => {
    // TODO: Implement exercise library route
    console.log('Navigate to exercise library');
  };

  const handleNavigateToStartWorkout = () => {
    // TODO: Implement start workout route
    console.log('Navigate to start workout');
  };

  if (!user) {
    router.replace('/');
    return null;
  }

  return (
    <HomeScreen
      user={user}
      onSignOut={handleSignOut}
      onNavigateToCreateRoutine={handleNavigateToCreateRoutine}
      onNavigateToExerciseLibrary={handleNavigateToExerciseLibrary}
      onNavigateToStartWorkout={handleNavigateToStartWorkout}
    />
  );
} 