import React, { useState, useEffect, useRef } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  Alert,
  Vibration,
} from 'react-native';
import { router, useLocalSearchParams } from 'expo-router';
import { Typography, Button } from '../src/components';
import { getColor } from '../src/components/Colors';
import { useAuth } from '../src/hooks/useAuth';

interface WorkoutSet {
  reps: number;
  weight: number;
  durationSeconds?: number;
  distance?: number;
  notes?: string;
  completed?: boolean;
}

interface WorkoutExercise {
  exerciseId: string;
  exercise?: {
    id: string;
    name: string;
    description: string;
    muscleGroup: string;
    equipment: string;
    difficulty: string;
    instructions: string[];
  };
  sets: WorkoutSet[];
  notes?: string;
  restSeconds?: number;
  completed?: boolean;
}

interface ActiveWorkout {
  id?: string;
  sessionId?: string; // Track the workout session ID
  routineId: string;
  routineName: string;
  routineDescription?: string;
  exercises: WorkoutExercise[];
  startedAt: Date;
  notes?: string;
}

export default function ActiveWorkoutScreen() {
  const { user } = useAuth();
  const params = useLocalSearchParams();
  const routineId = params.routineId as string;
  
  // Timer state
  const [startTime] = useState(new Date());
  const [elapsedTime, setElapsedTime] = useState(0);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  
  // Workout state
  const [activeWorkout, setActiveWorkout] = useState<ActiveWorkout | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Timer effect
  useEffect(() => {
    intervalRef.current = setInterval(() => {
      setElapsedTime(Date.now() - startTime.getTime());
    }, 1000);

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [startTime]);

  // Load routine data and start workout
  useEffect(() => {
    if (routineId && user?.id) {
      fetchRoutineAndStartWorkout();
    }
  }, [routineId, user?.id]);

  const fetchRoutineAndStartWorkout = async () => {
    try {
      setLoading(true);
      setError(null);

      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      
      // Fetch the specific routine/workout with populated exercise details
      const routineResponse = await fetch(`${API_BASE_URL}/api/workouts/${routineId}/start`);
      
      if (!routineResponse.ok) {
        throw new Error('Failed to load routine');
      }

      const routineData = await routineResponse.json();
      
      // Create a workout session from the routine
      const sessionData = {
        userId: user?.id,
        routineId: routineData.id,
        name: `${routineData.name} - ${new Date().toLocaleDateString()}`,
        description: `Active workout session for ${routineData.name}`,
        notes: '',
      };

      const sessionResponse = await fetch(`${API_BASE_URL}/api/workout-sessions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(sessionData),
      });

      if (!sessionResponse.ok) {
        throw new Error('Failed to create workout session');
      }

      const sessionResult = await sessionResponse.json();
      console.log('Created workout session:', sessionResult);
      
      // Convert the session data to active workout format
      const workout: ActiveWorkout = {
        sessionId: sessionResult.id,
        routineId: routineData.id,
        routineName: routineData.name,
        routineDescription: routineData.description,
        exercises: sessionResult.exercises.map((ex: any) => ({
          ...ex,
          completed: false,
          sets: ex.sets.map((set: any) => ({
            reps: set.targetReps || set.reps,
            weight: set.targetWeight || set.weight,
            durationSeconds: set.durationSeconds,
            distance: set.distance,
            notes: set.notes,
            completed: false,
          })),
        })),
        startedAt: startTime,
      };

      setActiveWorkout(workout);
    } catch (err) {
      console.error('Error loading routine and starting session:', err);
      setError(err instanceof Error ? err.message : 'Failed to start workout session');
    } finally {
      setLoading(false);
    }
  };

  const formatTime = (milliseconds: number): string => {
    const totalSeconds = Math.floor(milliseconds / 1000);
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const seconds = totalSeconds % 60;

    if (hours > 0) {
      return `${hours}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    }
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  const toggleSetComplete = async (exerciseIndex: number, setIndex: number) => {
    if (!activeWorkout || !activeWorkout.sessionId) return;

    const currentSet = activeWorkout.exercises[exerciseIndex].sets[setIndex];
    const newCompletedState = !currentSet.completed;

    // Update local state immediately for responsive UI
    setActiveWorkout(prev => {
      if (!prev) return prev;

      const updated = { ...prev };
      updated.exercises = [...prev.exercises];
      updated.exercises[exerciseIndex] = { ...prev.exercises[exerciseIndex] };
      updated.exercises[exerciseIndex].sets = [...prev.exercises[exerciseIndex].sets];
      updated.exercises[exerciseIndex].sets[setIndex] = {
        ...prev.exercises[exerciseIndex].sets[setIndex],
        completed: newCompletedState,
      };

      // Check if all sets in exercise are completed
      const allSetsCompleted = updated.exercises[exerciseIndex].sets.every(set => set.completed);
      updated.exercises[exerciseIndex].completed = allSetsCompleted;

      return updated;
    });

    // Haptic feedback
    Vibration.vibrate(50);

    // Update the backend workout session
    try {
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      
      const updateData = {
        actualReps: currentSet.reps,
        actualWeight: currentSet.weight,
        durationSeconds: currentSet.durationSeconds || 0,
        distance: currentSet.distance || 0,
        notes: currentSet.notes || '',
        completed: newCompletedState,
      };

      const response = await fetch(
        `${API_BASE_URL}/api/workout-sessions/${activeWorkout.sessionId}/exercises/${exerciseIndex}/sets/${setIndex}`,
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(updateData),
        }
      );

      if (!response.ok) {
        console.error('Failed to update set in backend');
        // Optionally revert the local state change here
      } else {
        console.log(`Set ${setIndex + 1} of exercise ${exerciseIndex + 1} updated successfully`);
      }
    } catch (error) {
      console.error('Error updating set:', error);
      // Optionally revert the local state change here
    }
  };

  const toggleExerciseComplete = async (exerciseIndex: number) => {
    if (!activeWorkout || !activeWorkout.sessionId) return;

    const currentExercise = activeWorkout.exercises[exerciseIndex];
    const newCompletedState = !currentExercise.completed;

    // Update local state immediately for responsive UI
    setActiveWorkout(prev => {
      if (!prev) return prev;

      const updated = { ...prev };
      updated.exercises = [...prev.exercises];
      updated.exercises[exerciseIndex] = { ...prev.exercises[exerciseIndex] };
      
      updated.exercises[exerciseIndex].completed = newCompletedState;
      
      // Update all sets to match exercise completion state
      updated.exercises[exerciseIndex].sets = updated.exercises[exerciseIndex].sets.map(set => ({
        ...set,
        completed: newCompletedState,
      }));

      return updated;
    });

    // Haptic feedback
    Vibration.vibrate(newCompletedState ? [0, 100, 50, 100] : 50);

    // Update all sets in the backend to match the exercise completion state
    try {
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      
      // Update each set in the exercise
      const setUpdatePromises = currentExercise.sets.map(async (set, setIndex) => {
        const updateData = {
          actualReps: set.reps,
          actualWeight: set.weight,
          durationSeconds: set.durationSeconds || 0,
          distance: set.distance || 0,
          notes: set.notes || '',
          completed: newCompletedState,
        };

        return fetch(
          `${API_BASE_URL}/api/workout-sessions/${activeWorkout.sessionId}/exercises/${exerciseIndex}/sets/${setIndex}`,
          {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(updateData),
          }
        );
      });

      const results = await Promise.all(setUpdatePromises);
      const failedUpdates = results.filter(response => !response.ok);
      
      if (failedUpdates.length > 0) {
        console.error(`Failed to update ${failedUpdates.length} sets in backend`);
      } else {
        console.log(`Exercise ${exerciseIndex + 1} completion state updated successfully`);
      }

      // Mark the exercise as finished/started in the backend
      if (newCompletedState) {
        await fetch(
          `${API_BASE_URL}/api/workout-sessions/${activeWorkout.sessionId}/exercises/${exerciseIndex}/finish`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ notes: '' }),
          }
        );
      }
    } catch (error) {
      console.error('Error updating exercise completion:', error);
      // Optionally revert the local state change here
    }
  };

  const handleEndWorkout = async () => {
    if (!activeWorkout || !user?.id) return;

    Alert.alert(
      'End Workout?',
      'Are you sure you want to end this workout? Your progress will be saved.',
      [
        { text: 'Cancel', style: 'cancel' },
        { text: 'End Workout', style: 'destructive', onPress: confirmEndWorkout },
      ]
    );
  };

  const confirmEndWorkout = async () => {
    if (!activeWorkout || !user?.id) return;

    try {
      setSaving(true);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      const finishedAt = new Date();
      const durationSeconds = Math.floor((finishedAt.getTime() - startTime.getTime()) / 1000);

      // Update the existing session with actual performance data and mark as completed
      if (!activeWorkout.sessionId) {
        throw new Error('No active workout session found');
      }
      const updateData = {
        exercises: activeWorkout.exercises.map(ex => ({
          exerciseId: ex.exerciseId,
          sets: ex.sets.map(set => ({
            targetReps: set.reps,
            targetWeight: set.weight,
            actualReps: set.reps, // In this implementation, we're using target as actual
            actualWeight: set.weight, // In a real app, user would input actual values
            durationSeconds: set.durationSeconds || 0,
            distance: set.distance || 0,
            notes: set.notes || '',
            completed: set.completed || false,
            startedAt: startTime.toISOString(),
            finishedAt: set.completed ? finishedAt.toISOString() : null,
          })),
          notes: ex.notes || '',
          restSeconds: ex.restSeconds || 0,
          completed: ex.completed || false,
          startedAt: startTime.toISOString(),
          finishedAt: ex.completed ? finishedAt.toISOString() : null,
        })),
        finishedAt: finishedAt.toISOString(),
        durationSeconds,
        isActive: false,
        notes: activeWorkout.notes || `Completed ${activeWorkout.exercises.filter(ex => ex.completed).length}/${activeWorkout.exercises.length} exercises`,
      };

      // Update the session to mark it as completed
      const updateResponse = await fetch(`${API_BASE_URL}/api/workout-sessions/${activeWorkout.sessionId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updateData),
      });

      if (!updateResponse.ok) {
        throw new Error('Failed to update workout session');
      }

      console.log('Workout session completed successfully');

      // Navigate back to home with success message
      router.replace('/(tabs)');
      
      // Show success alert after navigation
      setTimeout(() => {
        Alert.alert(
          'Workout Completed!',
          `Great job! You worked out for ${formatTime(elapsedTime)} and completed ${activeWorkout.exercises.filter(ex => ex.completed).length}/${activeWorkout.exercises.length} exercises.`,
          [{ text: 'OK' }]
        );
      }, 500);

    } catch (err) {
      console.error('Error saving workout session:', err);
      Alert.alert(
        'Error',
        'Failed to save workout session. Please try again.',
        [{ text: 'OK' }]
      );
    } finally {
      setSaving(false);
    }
  };

  const handleCancelWorkout = () => {
    Alert.alert(
      'Cancel Workout?',
      'Are you sure you want to cancel this workout? Your progress will be lost.',
      [
        { text: 'Keep Going', style: 'cancel' },
        { text: 'Cancel Workout', style: 'destructive', onPress: () => router.back() },
      ]
    );
  };

  if (loading) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.centerContent}>
          <Typography variant="text-default" color="light">
            Loading workout...
          </Typography>
        </View>
      </SafeAreaView>
    );
  }

  if (error || !activeWorkout) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.centerContent}>
          <Typography variant="text-default" color="red" style={styles.errorText}>
            {error || 'Failed to load workout'}
          </Typography>
          <Button variant="primary" size="default" onPress={() => router.back()}>
            Go Back
          </Button>
        </View>
      </SafeAreaView>
    );
  }

  const completedExercises = activeWorkout.exercises.filter(ex => ex.completed).length;
  const totalExercises = activeWorkout.exercises.length;
  const completedSets = activeWorkout.exercises.reduce((total, ex) => total + ex.sets.filter(set => set.completed).length, 0);
  const totalSets = activeWorkout.exercises.reduce((total, ex) => total + ex.sets.length, 0);

  return (
    <SafeAreaView style={styles.container}>
      {/* Header with Timer */}
      <View style={styles.headerCard}>
        <View style={styles.timerSection}>
          <Typography variant="heading-large" color="dark" style={styles.timerText}>
            {formatTime(elapsedTime)}
          </Typography>
          <Typography variant="text-default" color="light">
            Workout Duration
          </Typography>
        </View>
        
        <View style={styles.progressSection}>
          <Typography variant="text-small" color="light">
            {completedExercises}/{totalExercises} exercises • {completedSets}/{totalSets} sets
          </Typography>
        </View>
      </View>

      {/* Workout Info */}
      <View style={styles.card}>
        <Typography variant="heading-default" color="dark" style={styles.workoutTitle}>
          {activeWorkout.routineName}
        </Typography>
        {activeWorkout.routineDescription && (
          <Typography variant="text-default" color="light" style={styles.workoutDescription}>
            {activeWorkout.routineDescription}
          </Typography>
        )}
      </View>

      {/* Exercises List */}
      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        <View style={styles.exercisesList}>
          {activeWorkout.exercises.map((exercise, exerciseIndex) => (
            <View key={exerciseIndex} style={[
              styles.exerciseCard,
              exercise.completed && styles.exerciseCardCompleted
            ]}>
              {/* Exercise Header */}
              <View style={styles.exerciseHeader}>
                <View style={styles.exerciseInfo}>
                  <Typography variant="text-default" color="dark" style={styles.exerciseName}>
                    {exercise.exercise?.name || `Exercise ${exerciseIndex + 1}`}
                  </Typography>
                  {exercise.exercise?.muscleGroup && (
                    <Typography variant="text-small" color="light">
                      {exercise.exercise.muscleGroup}
                    </Typography>
                  )}
                </View>
                <Button
                  variant={exercise.completed ? "primary" : "secondary"}
                  size="small"
                  onPress={() => toggleExerciseComplete(exerciseIndex)}
                >
                  {exercise.completed ? "✓ Complete" : "Mark Complete"}
                </Button>
              </View>

              {/* Sets List */}
              <View style={styles.setsSection}>
                <Typography variant="text-small" color="light" style={styles.setsHeader}>
                  Sets ({exercise.sets.filter(set => set.completed).length}/{exercise.sets.length} completed)
                </Typography>
                
                {exercise.sets.map((set, setIndex) => (
                  <View key={setIndex} style={[
                    styles.setRow,
                    set.completed && styles.setRowCompleted
                  ]}>
                    <View style={styles.setInfo}>
                      <Typography variant="text-small" color="dark" style={styles.setNumber}>
                        Set {setIndex + 1}
                      </Typography>
                      <Typography variant="text-small" color="light">
                        {set.reps} reps
                        {set.weight > 0 && ` • ${set.weight}lbs`}
                        {set.durationSeconds && set.durationSeconds > 0 && ` • ${Math.floor(set.durationSeconds / 60)}:${(set.durationSeconds % 60).toString().padStart(2, '0')}`}
                        {set.notes && ` • ${set.notes}`}
                      </Typography>
                    </View>
                    <Button
                      variant={set.completed ? "primary" : "text"}
                      size="small"
                      onPress={() => toggleSetComplete(exerciseIndex, setIndex)}
                    >
                      {set.completed ? "✓" : "Done"}
                    </Button>
                  </View>
                ))}
              </View>

              {exercise.notes && (
                <View style={styles.exerciseNotes}>
                  <Typography variant="text-small" color="light">
                    Notes: {exercise.notes}
                  </Typography>
                </View>
              )}
            </View>
          ))}
        </View>
      </ScrollView>

      {/* Bottom Actions */}
      <View style={styles.bottomActions}>
        <Button
          variant="secondary"
          size="large"
          style={styles.actionButton}
          onPress={handleCancelWorkout}
        >
          Cancel Workout
        </Button>
        <Button
          variant="danger"
          size="large"
          style={styles.actionButton}
          onPress={handleEndWorkout}
          disabled={saving}
        >
          {saving ? "Saving..." : "End Workout"}
        </Button>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('light-gray-1'),
  },
  centerContent: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 32,
  },
  errorText: {
    textAlign: 'center',
    marginBottom: 24,
  },
  headerCard: {
    backgroundColor: 'white',
    marginHorizontal: 16,
    marginTop: 16,
    marginBottom: 16,
    borderRadius: 8,
    padding: 20,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
    minHeight: 140,
  },
  timerSection: {
    alignItems: 'center',
    marginBottom: 20,
    flex: 1,
    justifyContent: 'center',
  },
  timerText: {
    fontSize: 56,
    fontWeight: 'bold',
    fontFamily: 'monospace',
    marginBottom: 12,
    lineHeight: 64,
    textAlign: 'center',
  },
  progressSection: {
    alignItems: 'center',
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
  workoutTitle: {
    marginBottom: 8,
  },
  workoutDescription: {
    marginBottom: 0,
  },
  scrollView: {
    flex: 1,
  },
  exercisesList: {
    paddingHorizontal: 16,
    paddingBottom: 100, // Space for bottom actions
  },
  exerciseCard: {
    backgroundColor: 'white',
    borderRadius: 8,
    padding: 16,
    marginBottom: 16,
    borderWidth: 2,
    borderColor: getColor('light-gray-3'),
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 2,
    elevation: 2,
  },
  exerciseCardCompleted: {
    borderColor: getColor('green'),
    backgroundColor: getColor('green-bright'),
  },
  exerciseHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 16,
  },
  exerciseInfo: {
    flex: 1,
    marginRight: 12,
  },
  exerciseName: {
    fontWeight: '600',
    marginBottom: 4,
  },
  setsSection: {
    marginBottom: 12,
  },
  setsHeader: {
    marginBottom: 12,
    fontWeight: '500',
  },
  setRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 12,
    paddingHorizontal: 16,
    marginBottom: 8,
    backgroundColor: getColor('light-gray-1'),
    borderRadius: 6,
    borderWidth: 1,
    borderColor: getColor('light-gray-3'),
  },
  setRowCompleted: {
    backgroundColor: getColor('green-bright'),
    borderColor: getColor('green'),
  },
  setInfo: {
    flex: 1,
  },
  setNumber: {
    fontWeight: '500',
    marginBottom: 2,
  },
  exerciseNotes: {
    marginTop: 8,
    paddingTop: 12,
    borderTopWidth: 1,
    borderTopColor: getColor('light-gray-3'),
  },
  bottomActions: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'white',
    flexDirection: 'row',
    padding: 16,
    gap: 12,
    borderTopWidth: 1,
    borderTopColor: getColor('light-gray-3'),
    shadowColor: '#000',
    shadowOffset: { width: 0, height: -2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 5,
  },
  actionButton: {
    flex: 1,
  },
}); 