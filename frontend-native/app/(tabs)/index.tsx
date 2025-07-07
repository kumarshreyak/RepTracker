import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  Image,
  ScrollView,
  Alert,
  TouchableOpacity,
  Pressable,
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
  const [workoutsPageSize, setWorkoutsPageSize] = useState(5);

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
      const url = `${API_BASE_URL}/api/workouts?userId=${user.id}`;
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

  const fetchPastWorkouts = async (pageSize: number = 5, reset: boolean = false) => {
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
      fetchPastWorkouts(5, true); // Initial load with 5 workouts
    }
  }, [user?.id]);

  // Refresh routines and past workouts when screen comes back into focus
  useFocusEffect(
    React.useCallback(() => {
      if (user?.id) {
        fetchRoutines();
        fetchPastWorkouts(5, true); // Reset to initial load
        setWorkoutsPageSize(5); // Reset page size
      }
    }, [user?.id])
  );

  const handleCreateRoutine = () => {
    router.push('/create-routine');
  };

  const handleEditRoutine = (routineId: string) => {
    router.push(`/create-routine?routineId=${routineId}`);
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

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - date.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 1) return 'Today';
    if (diffDays === 2) return 'Yesterday';
    if (diffDays <= 7) return `${diffDays - 1}d ago`;
    
    return date.toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric' 
    });
  };

  const getCompletedExercisesCount = (exercises: WorkoutSessionExercise[]): number => {
    return exercises.filter(ex => ex.completed).length;
  };

  if (!user) {
    router.replace('/');
    return null;
  }

  const renderRoutineCard = (routine: Routine, index: number) => (
    <Pressable 
      key={routine.id} 
      style={styles.routineCard}
      onPress={() => handleStartRoutineWorkout(routine.id)}
    >
      <View style={styles.routineCardHeader}>
        <View style={styles.routineCardContent}>
          <Typography variant="label-medium" color="contentPrimary" style={styles.routineName}>
            {routine.name}
          </Typography>
          <Typography variant="paragraph-xsmall" color="contentSecondary" style={styles.routineExerciseCount}>
            {routine.exercises?.length || 0} exercises
          </Typography>
        </View>
        
        <Pressable 
          style={styles.editButton}
          onPress={(e) => {
            e.stopPropagation();
            handleEditRoutine(routine.id);
          }}
        >
          <Typography variant="paragraph-xsmall" color="contentSecondary" style={styles.editIcon}>
            ✏️
          </Typography>
        </Pressable>
      </View>
      
      <View style={styles.playButton}>
        <Typography variant="label-medium" color="contentOnColor" style={styles.playIcon}>
          ▶
        </Typography>
      </View>
    </Pressable>
  );

  const renderWorkoutChip = (session: WorkoutSession) => {
    const completedExercises = getCompletedExercisesCount(session.exercises);
    const totalExercises = session.exercises.length;
    const completionPercentage = totalExercises > 0 ? (completedExercises / totalExercises) * 100 : 0;
    
    return (
      <Pressable 
        key={session.id} 
        style={styles.workoutChip}
        onPress={() => router.push(`/workout-detail?id=${session.id}`)}
      >
        <View style={styles.workoutChipContent}>
          <Typography variant="label-xsmall" color="contentPrimary" style={styles.workoutChipName}>
            {session.name.replace(/ - \d+\/\d+\/\d+/, '')}
          </Typography>
          <Typography variant="paragraph-xsmall" color="contentSecondary" style={styles.workoutChipMeta}>
            {formatDate(session.finishedAt || session.startedAt)} • {formatDuration(session.durationSeconds)}
          </Typography>
        </View>
        
        <View style={styles.workoutProgressBar}>
          <View 
            style={[
              styles.workoutProgressFill, 
              { width: `${completionPercentage}%` }
            ]} 
          />
        </View>
      </Pressable>
    );
  };

  return (
    <SafeAreaView style={styles.container}>
      {/* App Bar */}
      <View style={styles.appBar}>
        <Typography variant="heading-xsmall" color="contentPrimary">
          RepTracker
        </Typography>
      </View>

      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        {/* Section A - My Routines */}
        <View style={styles.section}>
          <View style={styles.sectionHeader}>
            <Typography variant="label-medium" color="contentPrimary">
              My Routines
            </Typography>
            <Button
              variant="primary"
              size="small"
              onPress={handleCreateRoutine}
            >
              Create
            </Button>
          </View>

          {loading ? (
            <View style={styles.centerContent}>
              <Typography variant="paragraph-small" color="contentSecondary">
                Loading...
              </Typography>
            </View>
          ) : error ? (
            <View style={styles.centerContent}>
              <Typography variant="paragraph-small" color="contentNegative">
                {error}
              </Typography>
            </View>
          ) : routines.length === 0 ? (
            <View style={styles.emptyState}>
              <Typography variant="paragraph-small" color="contentSecondary" style={styles.emptyStateText}>
                No routines yet
              </Typography>
              <Button 
                variant="primary" 
                size="small"
                style={styles.emptyStateButton}
                onPress={handleCreateRoutine}
              >
                Create First Routine
              </Button>
            </View>
          ) : (
            <ScrollView 
              horizontal 
              showsHorizontalScrollIndicator={false}
              style={styles.routinesCarousel}
              contentContainerStyle={styles.routinesCarouselContent}
            >
              {routines.map(renderRoutineCard)}
            </ScrollView>
          )}
        </View>

        {/* Section B - Past Workouts */}
        <View style={styles.section}>
          <Typography variant="label-medium" color="contentPrimary" style={styles.sectionTitle}>
            Past Workouts
          </Typography>

          {workoutsLoading && pastWorkouts.length === 0 ? (
            <View style={styles.centerContent}>
              <Typography variant="paragraph-small" color="contentSecondary">
                Loading...
              </Typography>
            </View>
          ) : workoutsError ? (
            <View style={styles.centerContent}>
              <Typography variant="paragraph-small" color="contentNegative">
                {workoutsError}
              </Typography>
            </View>
          ) : pastWorkouts.length === 0 ? (
            <View style={styles.emptyState}>
              <Typography variant="paragraph-small" color="contentSecondary" style={styles.emptyStateText}>
                No workouts completed yet
              </Typography>
            </View>
          ) : (
            <View style={styles.workoutsList}>
              {pastWorkouts.map(renderWorkoutChip)}
              
              {hasMoreWorkouts && (
                <TouchableOpacity 
                  style={styles.viewMoreButton}
                  onPress={handleViewMoreWorkouts}
                  disabled={workoutsLoading}
                >
                  <Typography variant="paragraph-small" color="contentSecondary">
                    {workoutsLoading ? 'Loading...' : 'View More'}
                  </Typography>
                </TouchableOpacity>
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
    backgroundColor: getColor('backgroundPrimary'),
  },
  
  // App Bar
  appBar: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: getColor('backgroundPrimary'),
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderOpaque'),
  },


  scrollView: {
    flex: 1,
  },

  // Sections
  section: {
    paddingTop: 24,
    paddingBottom: 8,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    marginBottom: 16,
  },
  sectionTitle: {
    paddingHorizontal: 16,
    marginBottom: 16,
  },

  // Routines Carousel
  routinesCarousel: {
    paddingLeft: 16,
  },
  routinesCarouselContent: {
    paddingRight: 16,
  },
  routineCard: {
    width: 200,
    height: 120,
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    marginRight: 12,
    padding: 16,
    justifyContent: 'space-between',
  },
  routineCardHeader: {
    flex: 1,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
  },
  routineCardContent: {
    flex: 1,
  },
  routineName: {
    marginBottom: 4,
  },
  editButton: {
    padding: 4,
    marginTop: -4,
    marginRight: -4,
  },
  editIcon: {
    fontSize: 16,
  },
  routineExerciseCount: {
    marginBottom: 0,
  },
  playButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: getColor('backgroundAccent'),
    justifyContent: 'center',
    alignItems: 'center',
    alignSelf: 'flex-end',
  },
  playIcon: {
    marginLeft: 2, // Optical centering for play icon
  },

  // Past Workouts List
  workoutsList: {
    paddingHorizontal: 16,
  },
  workoutChip: {
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    padding: 16,
    marginBottom: 8,
  },
  workoutChipContent: {
    marginBottom: 12,
  },
  workoutChipName: {
    marginBottom: 4,
  },
  workoutChipMeta: {
    marginBottom: 0,
  },
  workoutProgressBar: {
    height: 4,
    backgroundColor: getColor('borderOpaque'),
    borderRadius: 2,
    overflow: 'hidden',
  },
  workoutProgressFill: {
    height: '100%',
    backgroundColor: getColor('accent'),
    borderRadius: 2,
  },

  // States
  centerContent: {
    alignItems: 'center',
    paddingVertical: 32,
    paddingHorizontal: 16,
  },
  emptyState: {
    alignItems: 'center',
    paddingVertical: 32,
    paddingHorizontal: 16,
  },
  emptyStateText: {
    marginBottom: 16,
    textAlign: 'center',
  },
  emptyStateButton: {
    minWidth: 180,
  },
  viewMoreButton: {
    alignItems: 'center',
    paddingVertical: 16,
  },


}); 