import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  Image,
  ScrollView,
  Alert,
} from 'react-native';
import { Typography, Button } from '../components';
import { getColor } from '../components/Colors';
import { authService, User } from '../auth/AuthService';

interface Routine {
  id: string;
  user_id: string;
  name: string;
  description: string;
  exercises: any[];
  created_at: string;
  updated_at: string;
}

interface HomeScreenProps {
  user: User;
  onSignOut: () => void;
  onNavigateToCreateRoutine: () => void;
  onNavigateToExerciseLibrary: () => void;
  onNavigateToStartWorkout: () => void;
}

export const HomeScreen: React.FC<HomeScreenProps> = ({ 
  user, 
  onSignOut, 
  onNavigateToCreateRoutine, 
  onNavigateToExerciseLibrary, 
  onNavigateToStartWorkout 
}) => {
  const [routines, setRoutines] = useState<Routine[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchRoutines = async () => {
    if (!user?.id) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      // TODO: Configure the API base URL for your React Native app
      // Example: const API_BASE_URL = 'http://localhost:3000' or your deployed backend URL
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:3000';
      const response = await fetch(`${API_BASE_URL}/api/workouts?user_id=${user.id}`);
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || "Failed to fetch routines");
      }
      
      if (data.success) {
        setRoutines(data.workouts || []);
      }
    } catch (err) {
      console.error("Error fetching routines:", err);
      setError(err instanceof Error ? err.message : "Failed to fetch routines");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (user?.id) {
      fetchRoutines();
    }
  }, [user?.id]);

  const handleSignOut = async () => {
    await authService.signOut();
    onSignOut();
  };

  const handleCreateRoutine = () => {
    onNavigateToCreateRoutine();
  };

  const handleStartWorkout = () => {
    onNavigateToStartWorkout();
  };

  const handleExerciseLibrary = () => {
    onNavigateToExerciseLibrary();
  };

  const handleEditRoutine = (routineId: string) => {
    // TODO: Navigate to edit routine screen
    console.log('Edit routine:', routineId);
  };

  const handleStartRoutineWorkout = (routineId: string) => {
    // TODO: Navigate to start workout with routine
    console.log('Start workout with routine:', routineId);
  };

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
          {routine.exercises?.length || 0} exercises • Created {new Date(routine.created_at).toLocaleDateString()}
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
        <View style={styles.headerCard}>
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

        {/* User Profile */}
        <View style={styles.card}>
          <Typography variant="heading-default" color="dark" style={styles.sectionTitle}>
            Your Profile
          </Typography>
          
          <View style={styles.profileGrid}>
            <View style={styles.profileItem}>
              <Typography variant="text-small" color="light" style={styles.profileLabel}>
                Name
              </Typography>
              <Typography variant="text-default" color="dark">
                {user?.name || 'Not provided'}
              </Typography>
            </View>
            
            <View style={styles.profileItem}>
              <Typography variant="text-small" color="light" style={styles.profileLabel}>
                Email
              </Typography>
              <Typography variant="text-default" color="dark">
                {user?.email || 'Not provided'}
              </Typography>
            </View>
            
            <View style={styles.profileItem}>
              <Typography variant="text-small" color="light" style={styles.profileLabel}>
                User ID
              </Typography>
              <Typography variant="text-default" color="dark">
                {user?.id || 'Not provided'}
              </Typography>
            </View>
            
            <View style={styles.profileItem}>
              <Typography variant="text-small" color="light" style={styles.profileLabel}>
                Session Token
              </Typography>
              <Typography variant="text-default" color="dark">
                {user?.sessionToken ? `${user.sessionToken.substring(0, 20)}...` : 'Not provided'}
              </Typography>
            </View>
          </View>

          {user?.picture && (
            <View style={styles.profilePictureSection}>
              <Typography variant="text-small" color="light" style={styles.profileLabel}>
                Profile Picture
              </Typography>
              <Image 
                source={{ uri: user.picture }} 
                style={styles.profilePicture}
              />
            </View>
          )}
        </View>

        {/* Quick Actions */}
        <View style={styles.card}>
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
        <View style={styles.card}>
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

        {/* System Status */}
        <View style={styles.card}>
          <Typography variant="heading-small" color="dark" style={styles.sectionTitle}>
            System Status
          </Typography>
          
          <View style={styles.statusGrid}>
            <View style={styles.statusItem}>
              <Typography variant="text-small" color="light" style={styles.statusLabel}>
                Authentication Status
              </Typography>
              <Typography variant="text-default" color="green">
                ✅ Authenticated with Backend
              </Typography>
            </View>
            
            <View style={styles.statusItem}>
              <Typography variant="text-small" color="light" style={styles.statusLabel}>
                Session Status
              </Typography>
              <Typography variant="text-default" color="green">
                ✅ Active Session
              </Typography>
            </View>
        </View>
      </View>
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('light-gray-1'),
  },
  scrollView: {
    flex: 1,
  },
  headerCard: {
    backgroundColor: 'white',
    marginHorizontal: 16,
    marginTop: 16,
    marginBottom: 16,
    borderRadius: 8,
    padding: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 2,
    elevation: 2,
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  welcomeTitle: {
    marginBottom: 8,
  },
  card: {
    backgroundColor: 'white',
    marginHorizontal: 16,
    marginBottom: 16,
    borderRadius: 8,
    padding: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 2,
    elevation: 2,
  },
  sectionTitle: {
    marginBottom: 16,
  },
  profileGrid: {
    gap: 16,
  },
  profileItem: {
    marginBottom: 8,
  },
  profileLabel: {
    marginBottom: 4,
  },
  profilePictureSection: {
    marginTop: 16,
  },
  profilePicture: {
    width: 64,
    height: 64,
    borderRadius: 32,
    marginTop: 8,
    backgroundColor: getColor('light-gray-3'),
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
    gap: 16,
  },
  routineCard: {
    padding: 16,
    borderWidth: 1,
    borderColor: getColor('light-gray-3'),
    borderRadius: 8,
    backgroundColor: 'white',
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
  statusGrid: {
    gap: 16,
  },
  statusItem: {
    marginBottom: 8,
  },
  statusLabel: {
    marginBottom: 4,
  },
}); 