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

export default function HomeTab() {
  const { user } = useAuth();
  const [routines, setRoutines] = useState<Routine[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

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

  useEffect(() => {
    if (user?.id) {
      fetchRoutines();
    }
  }, [user?.id]);

  // Refresh routines when screen comes back into focus
  useFocusEffect(
    React.useCallback(() => {
      if (user?.id) {
        fetchRoutines();
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

  if (!user) {
    router.replace('/');
    return null;
  }

  const renderRoutineCard = (routine: Routine) => (
    <View key={routine.id} style={styles.routineCard}>
      <View style={styles.routineInfo}>
        <Typography variant="text-default" color="dark" style={styles.routineName}>
          {routine.name}
        </Typography>
        <Typography variant="text-small" color="light" style={styles.routineDescription}>
          {routine.description}
        </Typography>
        <Typography variant="text-small" color="light" style={styles.routineMeta}>
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

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        {/* Header */}
        <View style={styles.header}>
          <View style={styles.headerContent}>
            <View>
              <Typography variant="heading-large" color="dark" style={styles.welcomeTitle}>
                Welcome to GymLog
              </Typography>
              <Typography variant="text-default" color="light">
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
          <Typography variant="heading-default" color="dark" style={styles.sectionTitle}>
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
            <Typography variant="heading-default" color="dark">
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
              <Typography variant="text-default" color="light">
                Loading routines...
              </Typography>
            </View>
          ) : error ? (
            <View style={styles.centerContent}>
              <Typography variant="text-default" color="red">
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
              <Typography variant="text-default" color="light" style={styles.emptyStateTitle}>
                No routines created yet
              </Typography>
              <Typography variant="text-small" color="light" style={styles.emptyStateSubtitle}>
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


      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('light-gray-1'),
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
    borderColor: getColor('light-gray-3'),
    borderRadius: 8,
    backgroundColor: 'white',
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

}); 