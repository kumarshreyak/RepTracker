import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  TextInput,
  Alert,
} from 'react-native';
import { router } from 'expo-router';
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

  const handleSaveRoutine = async () => {
    if (!routineName.trim() || exercises.length === 0 || !user?.id) {
      Alert.alert(
        'Validation Error',
        'Please enter a routine name and add at least one exercise.'
      );
      return;
    }

    setLoading(true);
    
    try {
      // Transform exercises to match backend format
      const workoutExercises = exercises.map(exercise => ({
        exercise_id: exercise.id,
        sets: exercise.sets.map(set => ({
          reps: set.reps,
          weight: set.weight,
          duration_seconds: 0,
          distance: 0,
          notes: "",
        })),
        notes: "",
        rest_seconds: 60,
      }));

      const workoutData = {
        user_id: user.id,
        name: routineName,
        description: `Workout routine with ${exercises.length} exercises`,
        exercises: workoutExercises,
        started_at: null,
        notes: "",
      };

      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(`${API_BASE_URL}/api/workouts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(workoutData),
      });

      if (!response.ok) {
        throw new Error(`Failed to save routine: ${response.statusText}`);
      }

      const result = await response.json();
      console.log('Routine saved successfully:', result);

      // Clear AsyncStorage after successful save
      await AsyncStorage.removeItem('routineExercises');
      
      Alert.alert(
        'Success',
        'Routine saved successfully!',
        [
          {
            text: 'OK',
            onPress: () => router.back(),
          },
        ]
      );
    } catch (error) {
      console.error('Error saving routine:', error);
      Alert.alert(
        'Error',
        'Failed to save routine. Please try again.'
      );
    } finally {
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
          <Typography variant="text-default" color="dark" style={styles.exerciseName}>
            {exercise.name}
          </Typography>
          <Typography variant="text-small" color="light">
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
            <Typography variant="text-small" color="dark" style={styles.setTitle}>
              Set {setIndex + 1}
            </Typography>
            <Typography variant="text-small" color="light">
              {set.reps} reps @ {set.weight}lbs
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
          <Typography variant="heading-xxlarge" color="dark" style={styles.title}>
            Create Routine
          </Typography>
          <Typography variant="text-default" color="light">
            Build your workout routine by adding exercises
          </Typography>
        </View>

        {/* Routine Details */}
        <View style={styles.card}>
          <Typography variant="heading-small" color="dark" style={styles.sectionTitle}>
            Routine Details
          </Typography>
          <View style={styles.inputContainer}>
            <Typography variant="text-small" color="dark" style={styles.inputLabel}>
              Routine Name *
            </Typography>
            <TextInput
              value={routineName}
              onChangeText={setRoutineName}
              placeholder="Enter routine name"
              style={styles.textInput}
              placeholderTextColor={getColor('light')}
            />
          </View>
        </View>

        {/* Exercises Section */}
        <View style={styles.card}>
          <View style={styles.exercisesHeader}>
            <Typography variant="heading-small" color="dark">
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
              <Typography variant="text-default" color="light" style={styles.emptyStateTitle}>
                No exercises added yet
              </Typography>
              <Typography variant="text-small" color="light">
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
  title: {
    marginVertical: 16,
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
  inputContainer: {
    marginBottom: 16,
  },
  inputLabel: {
    marginBottom: 8,
  },
  textInput: {
    borderWidth: 1,
    borderColor: getColor('light-gray-3'),
    borderRadius: 3,
    paddingHorizontal: 12,
    paddingVertical: 8,
    fontSize: 16,
    color: getColor('dark'),
    backgroundColor: 'white',
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
    borderColor: getColor('light-gray-3'),
    borderRadius: 8,
    backgroundColor: 'white',
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
    backgroundColor: getColor('light-gray-1'),
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
  },
  actionButton: {
    flex: 1,
  },
}); 