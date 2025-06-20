import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  Image,
  ScrollView,
  Alert,
} from 'react-native';
import { router, useFocusEffect } from 'expo-router';
import { Typography, Button } from '../../src/components';
import { getColor } from '../../src/components/Colors';
import { useAuth } from '../../src/hooks/useAuth';
import { authService, User } from '../../src/auth/AuthService';

interface Routine {
  id: string;
  userId: string;
  name: string;
  description: string;
  exercises: any[];
  createdAt: string;
  updatedAt: string;
}

interface WorkoutSessionSet {
  targetReps: number;
  targetWeight: number;
  actualReps: number;
  actualWeight: number;
  durationSeconds: number;
  distance: number;
  notes: string;
  completed: boolean;
  startedAt?: string;
  finishedAt?: string;
}

interface WorkoutSessionExercise {
  exerciseId: string;
  exercise?: any;
  sets: WorkoutSessionSet[];
  notes: string;
  restSeconds: number;
  completed: boolean;
  startedAt?: string;
  finishedAt?: string;
}

interface WorkoutSession {
  id: string;
  userId: string;
  routineId: string;
  name: string;
  description: string;
  exercises: WorkoutSessionExercise[];
  startedAt: string;
  finishedAt?: string;
  durationSeconds: number;
  notes: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export default function HomeTab() {
  const { user } = useAuth();
  const [routines, setRoutines] = useState<Routine[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Past workouts state
  const [pastWorkouts, setPastWorkouts] = useState<WorkoutSession[]>([]);
  const [workoutsLoading, setWorkoutsLoading] = useState(false);
  const [workoutsError, setWorkoutsError] = useState<string | null>(null);
  const [hasMoreWorkouts, setHasMoreWorkouts] = useState(true);
  const [workoutsPageSize, setWorkoutsPageSize] = useState(3);

  const fetchRoutines = async () => {
    console.log('[fetchRoutines] Starting fetch routine...');
    console.log('[fetchRoutines] User ID:', user?.id);
    
    if (!user?.id) {
      console.log('[fetchRoutines] No user ID found, exiting early');
      setLoading(false);
      return;
    }

    try {
      console.log('[fetchRoutines] Setting loading state to true');
      setLoading(true);
      setError(null);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:3000';
      const url = `${API_BASE_URL}/api/workouts?user_id=${user.id}`;
      console.log('[fetchRoutines] API Base URL:', API_BASE_URL);
      console.log('[fetchRoutines] Full URL:', url);
      
      console.log('[fetchRoutines] Making fetch request...');
      const response = await fetch(url);
      console.log('[fetchRoutines] Response received - Status:', response.status);
      console.log('[fetchRoutines] Response OK:', response.ok);
      console.log('[fetchRoutines] Response headers:', Object.fromEntries(response.headers.entries()));
      
      console.log('[fetchRoutines] Parsing response JSON...');
      const data = await response.json();
      console.log('[fetchRoutines] Response data:', JSON.stringify(data, null, 2));
      
      if (!response.ok) {
        console.log('[fetchRoutines] Response not OK, throwing error');
        throw new Error(data.error || "Failed to fetch routines");
      }
      
      console.log('[fetchRoutines] Success! Setting routines:', data.workouts?.length || 0, 'routines');
      setRoutines(data.workouts || []);
    } catch (err) {
      console.error("[fetchRoutines] Error caught:", err);
      console.error("[fetchRoutines] Error message:", err instanceof Error ? err.message : 'Unknown error');
      console.error("[fetchRoutines] Error stack:", err instanceof Error ? err.stack : 'No stack trace');
      setError(err instanceof Error ? err.message : "Failed to fetch routines");
    } finally {
      console.log('[fetchRoutines] Setting loading state to false');
      setLoading(false);
    }
  };

  const fetchPastWorkouts = async (pageSize: number = 3, reset: boolean = false) => {
    console.log('[fetchPastWorkouts] Starting fetch past workouts...');
    console.log('[fetchPastWorkouts] User ID:', user?.id);
    console.log('[fetchPastWorkouts] Page size:', pageSize);
    console.log('[fetchPastWorkouts] Reset:', reset);
    
    if (!user?.id) {
      console.log('[fetchPastWorkouts] No user ID found, exiting early');
      return;
    }

    try {
      setWorkoutsLoading(true);
      setWorkoutsError(null);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      // Fetch completed workout sessions (not active ones) ordered by most recent
      const url = `${API_BASE_URL}/api/workout-sessions?userId=${user.id}&pageSize=${pageSize}&activeOnly=false`;
      console.log('[fetchPastWorkouts] URL:', url);
      
      const response = await fetch(url);
      console.log('[fetchPastWorkouts] Response received - Status:', response.status);
      
      const data = await response.json();
      console.log('[fetchPastWorkouts] Response data:', JSON.stringify(data, null, 2));
      
      if (!response.ok) {
        throw new Error(data.error || "Failed to fetch past workouts");
      }
      
      const fetchedSessions = data.sessions || [];
      // Filter out active sessions and only show completed ones
      const completedSessions = fetchedSessions.filter((session: WorkoutSession) => 
        !session.isActive && session.finishedAt
      );
      
      console.log('[fetchPastWorkouts] Completed sessions:', completedSessions.length);
      
      if (reset) {
        setPastWorkouts(completedSessions);
      } else {
        setPastWorkouts(prev => [...prev, ...completedSessions]);
      }
      
      // Check if there are more workouts to load
      setHasMoreWorkouts(!!data.nextPageToken && completedSessions.length === pageSize);
      
    } catch (err) {
      console.error("[fetchPastWorkouts] Error caught:", err);
      setWorkoutsError(err instanceof Error ? err.message : "Failed to fetch past workouts");
    } finally {
      setWorkoutsLoading(false);
    }
  };

  useEffect(() => {
    if (user?.id) {
      fetchRoutines();
      fetchPastWorkouts(3, true); // Initial load with 3 workouts
    }
  }, [user?.id]);

  // Refresh routines and past workouts when screen comes back into focus
  useFocusEffect(
    React.useCallback(() => {
      if (user?.id) {
        fetchRoutines();
        fetchPastWorkouts(3, true); // Reset to initial load
        setWorkoutsPageSize(3); // Reset page size
      }
    }, [user?.id])
  );

  const handleSignOut = () => {
    router.replace('/sign-in');
  };

  const handleCreateRoutine = () => {
    router.push('/create-routine');
  };

  const handleStartWorkout = () => {
    // TODO: Implement start workout route
    console.log('Navigate to start workout');
  };

  const handleExerciseLibrary = () => {
    // TODO: Implement exercise library route
    console.log('Navigate to exercise library');
  };

  const handleEditRoutine = (routineId: string) => {
    // TODO: Navigate to edit routine screen
    console.log('Edit routine:', routineId);
  };

  const handleStartRoutineWorkout = (routineId: string) => {
    // Navigate to active workout screen with routine ID
    router.push(`/active-workout?routineId=${routineId}`);
  };

  const handleViewMoreWorkouts = () => {
    const newPageSize = workoutsPageSize + 10;
    setWorkoutsPageSize(newPageSize);
    fetchPastWorkouts(newPageSize, true);
  };

  const formatDuration = (seconds: number): string => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const remainingSeconds = seconds % 60;
    
    if (hours > 0) {
      return `${hours}h ${minutes}m`;
    }
    if (minutes > 0) {
      return remainingSeconds > 0 ? `${minutes}m ${remainingSeconds}s` : `${minutes}m`;
    }
    return `${remainingSeconds}s`;
  };

  const getCompletedExercisesCount = (exercises: WorkoutSessionExercise[]): number => {
    return exercises.filter(ex => ex.completed).length;
  };

  if (!user) {
    router.replace('/');
    return null;
  }

  const renderRoutineCard = (routine: Routine) => (
    <View key={routine.id} style={styles.routineCard}>
      <View style={styles.routineInfo}>
        <Typography variant="label-medium" color="text-primary" style={styles.routineName}>
          {routine.name}
        </Typography>
        <Typography variant="paragraph-small" color="text-secondary" style={styles.routineDescription}>
          {routine.description}
        </Typography>
        <Typography variant="paragraph-small" color="text-secondary" style={styles.routineMeta}>
          {routine.exercises?.length || 0} exercises • Created {new Date(routine.createdAt).toLocaleDateString()}
        </Typography>
      </View>
      
      <View style={styles.routineActions}>
        <Button 
          variant="primary" 
          size="small" 
          style={styles.routineActionButton}
          onPress={() => handleStartRoutineWorkout(routine.id)}
        >
          Start Workout
        </Button>
        <Button 
          variant="text" 
          size="small"
          onPress={() => handleEditRoutine(routine.id)}
        >
          Edit
        </Button>
      </View>
    </View>
  );

  const renderWorkoutSessionCard = (session: WorkoutSession) => {
    const completedExercises = getCompletedExercisesCount(session.exercises);
    const totalExercises = session.exercises.length;
    const completionPercentage = totalExercises > 0 ? (completedExercises / totalExercises) * 100 : 0;
    
    return (
      <View key={session.id} style={styles.workoutCard}>
        <View style={styles.workoutInfo}>
          <Typography variant="label-medium" color="text-primary" style={styles.workoutName}>
            {session.name}
          </Typography>
          <Typography variant="paragraph-small" color="text-secondary" style={styles.workoutDate}>
            {session.finishedAt ? new Date(session.finishedAt).toLocaleDateString() : new Date(session.startedAt).toLocaleDateString()}
          </Typography>
          
          <View style={styles.workoutStats}>
            <Typography variant="paragraph-small" color="text-secondary">
              {completedExercises}/{totalExercises} exercises completed
            </Typography>
            <Typography variant="paragraph-small" color="text-secondary">
              Duration: {formatDuration(session.durationSeconds)}
            </Typography>
          </View>
          
          {/* Progress bar */}
          <View style={styles.progressContainer}>
            <View style={styles.progressBar}>
              <View 
                style={[
                  styles.progressFill, 
                  { width: `${completionPercentage}%` }
                ]} 
              />
            </View>
            <Typography variant="paragraph-small" color="text-secondary">
              {Math.round(completionPercentage)}%
            </Typography>
          </View>
        </View>
      </View>
    );
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        {/* Header */}
        <View style={styles.header}>
          <View style={styles.headerContent}>
            <View>
                          <Typography variant="heading-large" color="text-primary" style={styles.welcomeTitle}>
              Welcome to GymLog
            </Typography>
            <Typography variant="paragraph-medium" color="text-secondary">
              Track your workouts and progress
            </Typography>
            </View>
            <Button
              variant="secondary"
              size="default"
              onPress={handleSignOut}
            >
              Sign Out
            </Button>
          </View>
        </View>



        {/* Quick Actions */}
        <View style={styles.section}>
          <Typography variant="heading-medium" color="text-primary" style={styles.sectionTitle}>
            Quick Actions
          </Typography>
          
          <View style={styles.quickActions}>
            <Button
              variant="primary"
              size="large"
              style={styles.quickActionButton}
              onPress={handleCreateRoutine}
            >
              Create Routine
            </Button>
            <Button 
              variant="secondary" 
              size="large" 
              style={styles.quickActionButton}
              onPress={handleStartWorkout}
            >
              Start Workout
            </Button>
            <Button
              variant="secondary"
              size="large"
              style={styles.quickActionButton}
              onPress={handleExerciseLibrary}
            >
              Exercise Library
            </Button>
          </View>
        </View>

        {/* My Routines */}
        <View style={styles.section}>
          <View style={styles.routinesHeader}>
            <Typography variant="heading-medium" color="text-primary">
              My Routines
            </Typography>
            <Button
              variant="text"
              size="default"
              onPress={handleCreateRoutine}
            >
              Create New
            </Button>
          </View>

          {loading ? (
            <View style={styles.centerContent}>
              <Typography variant="paragraph-medium" color="text-secondary">
                Loading routines...
              </Typography>
            </View>
          ) : error ? (
            <View style={styles.centerContent}>
              <Typography variant="paragraph-medium" color="danger">
                {error}
              </Typography>
              <Button 
                variant="text" 
                size="small" 
                style={styles.retryButton}
                onPress={fetchRoutines}
              >
                Try Again
              </Button>
            </View>
          ) : routines.length === 0 ? (
            <View style={styles.emptyState}>
              <Typography variant="paragraph-medium" color="text-secondary" style={styles.emptyStateTitle}>
                No routines created yet
              </Typography>
              <Typography variant="paragraph-small" color="text-secondary" style={styles.emptyStateSubtitle}>
                Create your first workout routine to get started
              </Typography>
              <Button 
                variant="primary" 
                size="default"
                style={styles.emptyStateButton}
                onPress={handleCreateRoutine}
              >
                Create Your First Routine
              </Button>
            </View>
          ) : (
            <View style={styles.routinesList}>
              {routines.map(renderRoutineCard)}
            </View>
          )}
        </View>

        {/* Past Workouts */}
        <View style={styles.section}>
          <Typography variant="heading-medium" color="text-primary" style={styles.sectionTitle}>
            Past Workouts
          </Typography>

          {workoutsLoading && pastWorkouts.length === 0 ? (
            <View style={styles.centerContent}>
              <Typography variant="paragraph-medium" color="text-secondary">
                Loading past workouts...
              </Typography>
            </View>
          ) : workoutsError ? (
            <View style={styles.centerContent}>
              <Typography variant="paragraph-medium" color="danger">
                {workoutsError}
              </Typography>
              <Button 
                variant="text" 
                size="small" 
                style={styles.retryButton}
                onPress={() => fetchPastWorkouts(3, true)}
              >
                Try Again
              </Button>
            </View>
          ) : pastWorkouts.length === 0 ? (
            <View style={styles.emptyState}>
              <Typography variant="paragraph-medium" color="text-secondary" style={styles.emptyStateTitle}>
                No past workouts yet
              </Typography>
              <Typography variant="paragraph-small" color="text-secondary" style={styles.emptyStateSubtitle}>
                Complete your first workout to see it here
              </Typography>
            </View>
          ) : (
            <View style={styles.workoutsList}>
              {pastWorkouts.map(renderWorkoutSessionCard)}
              
              {hasMoreWorkouts && (
                <Button 
                  variant="text" 
                  size="default"
                  style={styles.viewMoreButton}
                  onPress={handleViewMoreWorkouts}
                  disabled={workoutsLoading}
                >
                  {workoutsLoading ? 'Loading...' : 'View More'}
                </Button>
              )}
            </View>
          )}
        </View>


      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundSecondary'),
  },
  scrollView: {
    flex: 1,
  },
  header: {
    paddingHorizontal: 16,
    paddingTop: 16,
    paddingBottom: 24,
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  welcomeTitle: {
    marginBottom: 8,
  },
  section: {
    paddingHorizontal: 16,
    marginBottom: 24,
  },
  sectionTitle: {
    marginBottom: 16,
  },

  quickActions: {
    gap: 16,
  },
  quickActionButton: {
    width: '100%',
  },
  routinesHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  centerContent: {
    alignItems: 'center',
    paddingVertical: 32,
  },
  retryButton: {
    marginTop: 8,
  },
  emptyState: {
    alignItems: 'center',
    paddingVertical: 48,
  },
  emptyStateTitle: {
    marginBottom: 16,
  },
  emptyStateSubtitle: {
    marginBottom: 24,
    textAlign: 'center',
  },
  emptyStateButton: {
    minWidth: 200,
  },
  routinesList: {
  },
  routineCard: {
    padding: 16,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    borderRadius: 8,
    backgroundColor: getColor('backgroundPrimary'),
    marginBottom: 12,
  },
  routineInfo: {
    marginBottom: 12,
  },
  routineName: {
    fontWeight: '500',
    marginBottom: 4,
  },
  routineDescription: {
    marginBottom: 8,
  },
  routineMeta: {
    marginBottom: 0,
  },
  routineActions: {
    flexDirection: 'row',
    gap: 8,
  },
  routineActionButton: {
    flex: 1,
  },

  // Past workouts styles
  workoutsList: {
  },
  workoutCard: {
    padding: 16,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    borderRadius: 8,
    backgroundColor: getColor('backgroundPrimary'),
    marginBottom: 12,
  },
  workoutInfo: {
  },
  workoutName: {
    fontWeight: '500',
    marginBottom: 4,
  },
  workoutDate: {
    marginBottom: 12,
  },
  workoutStats: {
    marginBottom: 12,
    gap: 4,
  },
  progressContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  progressBar: {
    flex: 1,
    height: 8,
    backgroundColor: getColor('borderOpaque'),
    borderRadius: 4,
    overflow: 'hidden',
  },
  progressFill: {
    height: '100%',
    backgroundColor: getColor('accent'),
    borderRadius: 4,
  },
  viewMoreButton: {
    marginTop: 16,
    alignSelf: 'center',
  },

}); 