import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  TextInput,
  Alert,
} from 'react-native';
import { router, useFocusEffect } from 'expo-router';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { Typography, Button } from '../src/components';
import { getColor } from '../src/components/Colors';
import { useAuth } from '../src/hooks/useAuth';
import { authService, User } from '../src/auth/AuthService';

interface WorkoutSet {
  reps: number;
  weight: number;
}

interface RoutineExercise {
  id: string;
  name: string;
  muscle_group: string;
  sets: WorkoutSet[];
}

export default function CreateRoutineRoute() {
  const { user } = useAuth();
  const [routineName, setRoutineName] = useState("");
  const [exercises, setExercises] = useState<RoutineExercise[]>([]);
  const [loading, setLoading] = useState(false);

  // Load exercises from AsyncStorage on mount
  useEffect(() => {
    const loadExercises = async () => {
      try {
        const saved = await AsyncStorage.getItem('routineExercises');
        if (saved) {
          const parsedExercises = JSON.parse(saved);
          setExercises(parsedExercises);
        }
      } catch (error) {
        console.error('Error loading exercises from AsyncStorage:', error);
        setExercises([]);
      }
    };

    loadExercises();
  }, []);

  // Refresh exercises when screen comes back into focus
  useFocusEffect(
    React.useCallback(() => {
      const loadExercises = async () => {
        try {
          const saved = await AsyncStorage.getItem('routineExercises');
          if (saved) {
            const parsedExercises = JSON.parse(saved);
            setExercises(parsedExercises);
          } else {
            setExercises([]);
          }
        } catch (error) {
          console.error('Error loading exercises from AsyncStorage on focus:', error);
          setExercises([]);
        }
      };

      loadExercises();
    }, [])
  );

  const handleSaveRoutine = async () => {
    console.log('🚀 handleSaveRoutine: Starting routine save process');
    console.log('📝 handleSaveRoutine: Current state:', {
      routineName: routineName.trim(),
      exercisesCount: exercises.length,
      userId: user?.id,
      hasUser: !!user,
      userObject: user,
      userIdType: typeof user?.id,
      userIdLength: user?.id?.length,
    });

    // Check if user_id is in proper MongoDB ObjectID format (24 characters)
    const isValidUserId = user?.id && typeof user.id === 'string' && user.id.length === 24;
    
    if (!routineName.trim() || exercises.length === 0 || !user?.id || !isValidUserId) {
      console.log('❌ handleSaveRoutine: Validation failed:', {
        hasRoutineName: !!routineName.trim(),
        hasExercises: exercises.length > 0,
        hasUserId: !!user?.id,
        isValidUserId: isValidUserId,
        userIdActual: user?.id,
        userIdLength: user?.id?.length,
      });
      
      if (!isValidUserId) {
        Alert.alert(
          'Authentication Error',
          'Please sign out and sign in again to refresh your session.'
        );
      } else {
        Alert.alert(
          'Validation Error',
          'Please enter a routine name and add at least one exercise.'
        );
      }
      return;
    }

    console.log('✅ handleSaveRoutine: Validation passed, proceeding with save');
    setLoading(true);
    
    try {
      console.log('🔄 handleSaveRoutine: Transforming exercises data');
      console.log('📊 handleSaveRoutine: Raw exercises:', JSON.stringify(exercises, null, 2));

      // Transform exercises to match backend format
      const workoutExercises = exercises.map(exercise => ({
        exerciseId: exercise.id,
        sets: exercise.sets.map(set => ({
          reps: set.reps,
          weight: set.weight,
          durationSeconds: 0,
          distance: 0,
          notes: "",
        })),
        notes: "",
        restSeconds: 60,
      }));

      console.log('🔧 handleSaveRoutine: Transformed workout exercises:', JSON.stringify(workoutExercises, null, 2));

      const workoutData = {
        userId: user.id,
        name: routineName,
        description: `Workout routine with ${exercises.length} exercises`,
        exercises: workoutExercises,
        startedAt: null,
        notes: "",
      };

      console.log('📦 handleSaveRoutine: Final workout data:', JSON.stringify(workoutData, null, 2));

      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      console.log('🌐 handleSaveRoutine: Making API request to:', `${API_BASE_URL}/api/workouts`);

      const response = await fetch(`${API_BASE_URL}/api/workouts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(workoutData),
      });

      console.log('📡 handleSaveRoutine: Response received:', {
        status: response.status,
        statusText: response.statusText,
        ok: response.ok,
        headers: Object.fromEntries(response.headers.entries()),
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('❌ handleSaveRoutine: Response not ok:', {
          status: response.status,
          statusText: response.statusText,
          errorBody: errorText,
        });
        throw new Error(`Failed to save routine: ${response.statusText} - ${errorText}`);
      }

      const result = await response.json();
      console.log('✅ handleSaveRoutine: Routine saved successfully:', JSON.stringify(result, null, 2));

      // Clear AsyncStorage after successful save
      console.log('🗑️ handleSaveRoutine: Clearing AsyncStorage');
      await AsyncStorage.removeItem('routineExercises');
      console.log('✅ handleSaveRoutine: AsyncStorage cleared successfully');
      
      Alert.alert(
        'Success',
        'Routine saved successfully!',
        [
          {
            text: 'OK',
            onPress: () => {
              console.log('🔙 handleSaveRoutine: Navigating back');
              router.back();
            },
          },
        ]
      );
    } catch (error) {
      console.error('💥 handleSaveRoutine: Error occurred during save:', error);
      console.error('💥 handleSaveRoutine: Error stack:', error instanceof Error ? error.stack : 'No stack trace');
      Alert.alert(
        'Error',
        'Failed to save routine. Please try again.'
      );
    } finally {
      console.log('🏁 handleSaveRoutine: Setting loading to false');
      setLoading(false);
    }
  };

  const removeExercise = async (index: number) => {
    const newExercises = exercises.filter((_, i) => i !== index);
    setExercises(newExercises);
    await AsyncStorage.setItem('routineExercises', JSON.stringify(newExercises));
  };

  const clearAllExercises = async () => {
    Alert.alert(
      'Clear All Exercises',
      'Are you sure you want to remove all exercises?',
      [
        {
          text: 'Cancel',
          style: 'cancel',
        },
        {
          text: 'Clear All',
          style: 'destructive',
          onPress: async () => {
            setExercises([]);
            await AsyncStorage.removeItem('routineExercises');
          },
        },
      ]
    );
  };

  if (!user) {
    router.replace('/');
    return null;
  }

  const renderExerciseCard = (exercise: RoutineExercise, index: number) => (
    <View key={`${exercise.id}-${index}`} style={styles.exerciseCard}>
      <View style={styles.exerciseHeader}>
        <View style={styles.exerciseInfo}>
          <Typography variant="label-medium" color="text-primary" style={styles.exerciseName}>
            {exercise.name}
          </Typography>
          <Typography variant="paragraph-small" color="text-secondary">
            {exercise.muscle_group} • {exercise.sets.length} sets
          </Typography>
        </View>
        <Button
          variant="text"
          size="small"
          onPress={() => removeExercise(index)}
        >
          Remove
        </Button>
      </View>
      
      {/* Set Details */}
      <View style={styles.setsContainer}>
        {exercise.sets.map((set, setIndex) => (
          <View key={setIndex} style={styles.setCard}>
            <Typography variant="label-small" color="text-primary" style={styles.setTitle}>
              Set {setIndex + 1}
            </Typography>
            <Typography variant="paragraph-small" color="text-secondary">
              {set.reps} reps @ {set.weight}kg
            </Typography>
          </View>
        ))}
      </View>
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        {/* Page Header */}
        <View style={styles.header}>
          <Button variant="text" size="default" onPress={() => router.back()}>
            ← Back to Home
          </Button>
          <Typography variant="heading-xxlarge" color="text-primary" style={styles.title}>
            Create Routine
          </Typography>
          <Typography variant="paragraph-medium" color="text-secondary">
            Build your workout routine by adding exercises
          </Typography>
        </View>

        {/* Routine Details */}
        <View style={styles.section}>
          <Typography variant="heading-small" color="text-primary" style={styles.sectionTitle}>
            Routine Details
          </Typography>
          <View style={styles.inputContainer}>
            <Typography variant="label-small" color="text-primary" style={styles.inputLabel}>
              Routine Name *
            </Typography>
            <TextInput
              value={routineName}
              onChangeText={setRoutineName}
              placeholder="Enter routine name"
              style={styles.textInput}
              placeholderTextColor={getColor('text-secondary')}
            />
          </View>
        </View>

        {/* Exercises Section */}
        <View style={styles.section}>
          <View style={styles.exercisesHeader}>
            <Typography variant="heading-small" color="text-primary">
              Exercises ({exercises.length})
            </Typography>
            <View style={styles.exercisesActions}>
              {exercises.length > 0 && (
                <Button variant="text" size="small" onPress={clearAllExercises}>
                  Clear All
                </Button>
              )}
              <Button variant="primary" size="default" onPress={() => router.push('/exercise-search')}>
                Add Exercise
              </Button>
            </View>
          </View>

          {exercises.length === 0 ? (
            <View style={styles.emptyState}>
              <Typography variant="paragraph-medium" color="text-secondary" style={styles.emptyStateTitle}>
                No exercises added yet
              </Typography>
              <Typography variant="paragraph-small" color="text-secondary">
                Click "Add Exercise" to start building your routine
              </Typography>
            </View>
          ) : (
            <View style={styles.exercisesList}>
              {exercises.map(renderExerciseCard)}
            </View>
          )}
        </View>

        {/* Action Buttons */}
        <View style={styles.actionButtons}>
          <Button
            variant="secondary"
            size="large"
            style={styles.actionButton}
            onPress={() => router.back()}
          >
            Cancel
          </Button>
          <Button
            variant="primary"
            size="large"
            style={styles.actionButton}
            onPress={handleSaveRoutine}
            disabled={!routineName.trim() || exercises.length === 0 || loading || !user?.id}
          >
            {loading ? "Saving..." : "Save Routine"}
          </Button>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('background'),
  },
  scrollView: {
    flex: 1,
  },
  header: {
    paddingHorizontal: 16,
    paddingTop: 16,
    paddingBottom: 24,
  },
  title: {
    marginVertical: 16,
  },
  section: {
    paddingHorizontal: 16,
    marginBottom: 24,
  },
  sectionTitle: {
    marginBottom: 16,
  },
  inputContainer: {
    marginBottom: 16,
  },
  inputLabel: {
    marginBottom: 8,
  },
  textInput: {
    borderWidth: 1,
    borderColor: getColor('border-light'),
    borderRadius: 3,
    paddingHorizontal: 12,
    paddingVertical: 8,
    fontSize: 16,
    color: getColor('text-primary'),
    backgroundColor: getColor('surface'),
  },
  exercisesHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
    flexWrap: 'wrap',
    gap: 8,
  },
  exercisesActions: {
    flexDirection: 'row',
    gap: 12,
    alignItems: 'center',
  },
  emptyState: {
    alignItems: 'center',
    paddingVertical: 48,
  },
  emptyStateTitle: {
    marginBottom: 16,
    textAlign: 'center',
  },
  exercisesList: {
    gap: 16,
  },
  exerciseCard: {
    padding: 16,
    borderWidth: 1,
    borderColor: getColor('border-light'),
    borderRadius: 8,
    backgroundColor: getColor('surface'),
  },
  exerciseHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 12,
  },
  exerciseInfo: {
    flex: 1,
  },
  exerciseName: {
    fontWeight: '500',
    marginBottom: 4,
  },
  setsContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  setCard: {
    paddingHorizontal: 12,
    paddingVertical: 8,
    backgroundColor: getColor('background'),
    borderRadius: 8,
    minWidth: 100,
  },
  setTitle: {
    fontWeight: '500',
    marginBottom: 2,
  },
  actionButtons: {
    flexDirection: 'row',
    gap: 16,
    paddingHorizontal: 16,
    paddingBottom: 32,
    marginTop: 8,
  },
  actionButton: {
    flex: 1,
  },
}); 