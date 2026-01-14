import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  ScrollView,
  Alert,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { router, useFocusEffect, useLocalSearchParams } from 'expo-router';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { Typography, Button, Input } from '../src/components';
import { getColor } from '../src/components/Colors';
import { useAuth } from '../src/hooks/useAuth';
import { RoutineExercise, RoutineSet } from '@/types/exercise';
import { apiGet, apiPost, apiPut } from '../src/utils/api';

// Workout interface for API responses
interface Workout {
  id: string;
  userId: string;
  name: string;
  description: string;
  exercises: {
    exerciseId: string;
    exercise?: {
      id: string;
      name: string;
      primaryMuscles: string[];
      secondaryMuscles: string[];
    };
    sets: {
      reps: number;
      weight: number;
      durationSeconds?: number;
      distance?: number;
      notes?: string;
    }[];
    notes?: string;
    restSeconds?: number;
  }[];
  startedAt?: string;
  finishedAt?: string;
  durationSeconds?: number;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

// Quick add exercises will be fetched from the API

export default function CreateRoutineRoute() {
  const { user } = useAuth();
  const { routineId } = useLocalSearchParams<{ routineId?: string }>();
  const isEditing = !!routineId;
  
  const [routineName, setRoutineName] = useState("");
  const [exercises, setExercises] = useState<RoutineExercise[]>([]);
  const [quickAddExercises, setQuickAddExercises] = useState<RoutineExercise[]>([]);
  const [loading, setLoading] = useState(false);
  const [initialDataLoading, setInitialDataLoading] = useState(isEditing);

  // Load existing routine data if editing
  useEffect(() => {
    if (isEditing && routineId) {
      fetchExistingRoutine(routineId);
    }
  }, [isEditing, routineId]);

  const fetchExistingRoutine = async (id: string) => {
    console.log('🔄 fetchExistingRoutine: Starting fetch for routineId:', id);
    try {
      setInitialDataLoading(true);
      const workout: Workout = await apiGet<Workout>(`/api/workouts/${id}`);
      console.log('📦 fetchExistingRoutine: Received workout data:', JSON.stringify(workout, null, 2));
      
      // Set routine name
      setRoutineName(workout.name);
      
      // Fetch detailed exercise information for each exercise
      console.log('🔄 fetchExistingRoutine: Fetching detailed exercise data for each exercise');
      const routineExercises: RoutineExercise[] = await Promise.all(
        workout.exercises.map(async (ex) => {
          try {
            // Fetch detailed exercise data from the exercise API
            const exerciseData = await apiGet<any>(`/api/exercises/${ex.exerciseId}`);
            console.log(`📦 fetchExistingRoutine: Fetched exercise ${ex.exerciseId}:`, exerciseData.name);
            
            return {
              id: ex.exerciseId,
              name: exerciseData.name,
              primaryMuscles: exerciseData.primaryMuscles || [],
              secondaryMuscles: exerciseData.secondaryMuscles || [],
              sets: ex.sets.map(set => ({
                reps: set.reps,
                weight: set.weight,
                durationSeconds: set.durationSeconds,
                distance: set.distance,
                notes: set.notes,
              }))
            };
          } catch (error) {
            // Fallback to workout data if exercise API fails
            console.warn(`⚠️ fetchExistingRoutine: Failed to fetch exercise ${ex.exerciseId}, using fallback data`);
            return {
              id: ex.exerciseId,
              name: ex.exercise?.name || `Exercise ${ex.exerciseId}`,
              primaryMuscles: ex.exercise?.primaryMuscles || [],
              secondaryMuscles: ex.exercise?.secondaryMuscles || [],
              sets: ex.sets.map(set => ({
                reps: set.reps,
                weight: set.weight,
                durationSeconds: set.durationSeconds,
                distance: set.distance,
                notes: set.notes,
              }))
            };
          }
        })
      );
      
      setExercises(routineExercises);
      console.log('✅ fetchExistingRoutine: Successfully loaded routine data with detailed exercise information');
    } catch (error) {
      console.error('❌ fetchExistingRoutine: Error:', error);
      Alert.alert(
        'Error',
        'Failed to load routine data. Please try again.',
        [{ text: 'OK', onPress: () => router.back() }]
      );
    } finally {
      setInitialDataLoading(false);
    }
  };

  // Load exercises from AsyncStorage and fetch quick add exercises on mount
  useEffect(() => {
    if (!isEditing) {
      const loadExercises = async () => {
        try {
          const saved = await AsyncStorage.getItem('routineExercises');
          if (saved) {
            const parsedExercises = JSON.parse(saved);
            // Ensure exercises have the new structure
            const normalizedExercises = parsedExercises.map((ex: any) => ({
              ...ex,
              primaryMuscles: ex.primaryMuscles || (ex.muscleGroup ? [ex.muscleGroup] : []),
              secondaryMuscles: ex.secondaryMuscles || []
            }));
            setExercises(normalizedExercises);
          }
        } catch (error) {
          console.error('Error loading exercises from AsyncStorage:', error);
          setExercises([]);
        }
      };
      loadExercises();
    }

    const fetchQuickAddExercises = async () => {
      console.log('🏋️ fetchQuickAddExercises: Starting fetch process');
      
      try {
        // User ID is now handled by authentication, no need to pass it as query param
        const endpoint = '/api/exercises/quick-add?limit=5';
        console.log('📡 fetchQuickAddExercises: Making request to:', endpoint);

        const data = await apiGet<{ exercises: any[] }>(endpoint);
        console.log('📦 fetchQuickAddExercises: Raw response data:', JSON.stringify(data, null, 2));
        
        if (data.exercises && Array.isArray(data.exercises)) {
          const quickAddData = data.exercises.map((exercise: any) => ({
            id: exercise.id,
            name: exercise.name,
            primaryMuscles: exercise.primaryMuscles || [],
            secondaryMuscles: exercise.secondaryMuscles || [],
            sets: [
              { reps: 10, weight: 0 },
              { reps: 10, weight: 0 },
              { reps: 10, weight: 0 }
            ]
          }));
          console.log('🔧 fetchQuickAddExercises: Transformed quick add data:', JSON.stringify(quickAddData, null, 2));
          setQuickAddExercises(quickAddData);
          console.log('✅ fetchQuickAddExercises: Successfully set quick add exercises');
        } else {
          console.log('⚠️ fetchQuickAddExercises: No exercises array in response or invalid format');
          setQuickAddExercises([]);
        }
      } catch (error) {
        console.error('💥 fetchQuickAddExercises: Error occurred:', error);
        console.error('💥 fetchQuickAddExercises: Error stack:', error instanceof Error ? error.stack : 'No stack trace');
        // Fallback to empty array if API fails
        setQuickAddExercises([]);
      }
    };

    fetchQuickAddExercises();
  }, [user?.id, isEditing]);

  // Refresh exercises when screen comes back into focus
  useFocusEffect(
    React.useCallback(() => {
      const loadExercises = async () => {
        try {
          const saved = await AsyncStorage.getItem('routineExercises');
          if (saved) {
            const parsedExercises = JSON.parse(saved);
            // Ensure exercises have the new structure
            const normalizedExercises = parsedExercises.map((ex: any) => ({
              ...ex,
              primaryMuscles: ex.primaryMuscles || (ex.muscleGroup ? [ex.muscleGroup] : []),
              secondaryMuscles: ex.secondaryMuscles || []
            }));
            
            if (isEditing) {
              // In edit mode, combine existing routine exercises with newly added ones from AsyncStorage
              // Only add exercises that aren't already in the routine (avoid duplicates)
                             setExercises(prevExercises => {
                 const existingIds = new Set(prevExercises.map((ex: RoutineExercise) => ex.id));
                 const newExercises = normalizedExercises.filter((ex: RoutineExercise) => !existingIds.has(ex.id));
                 return [...prevExercises, ...newExercises];
               });
            } else {
              // In create mode, replace exercises with AsyncStorage content
              setExercises(normalizedExercises);
            }
          } else if (!isEditing) {
            // Only clear exercises if we're in create mode and there's nothing in AsyncStorage
            setExercises([]);
          }
        } catch (error) {
          console.error('Error loading exercises from AsyncStorage on focus:', error);
          if (!isEditing) {
            setExercises([]);
          }
        }
      };
      
      loadExercises();

      const fetchQuickAddExercises = async () => {
        console.log('🔄 fetchQuickAddExercises (focus): Starting fetch process');
        
        try {
          // User ID is now handled by authentication, no need to pass it as query param
          const endpoint = '/api/exercises/quick-add?limit=5';
          console.log('📡 fetchQuickAddExercises (focus): Making request to:', endpoint);

          const data = await apiGet<{ exercises: any[] }>(endpoint);
          console.log('📦 fetchQuickAddExercises (focus): Raw response data:', JSON.stringify(data, null, 2));
          
          if (data.exercises && Array.isArray(data.exercises)) {
            const quickAddData = data.exercises.map((exercise: any) => ({
              id: exercise.id,
              name: exercise.name,
              primaryMuscles: exercise.primaryMuscles || [],
              secondaryMuscles: exercise.secondaryMuscles || [],
              sets: [
                { reps: 10, weight: 0 },
                { reps: 10, weight: 0 },
                { reps: 10, weight: 0 }
              ]
            }));
            console.log('🔧 fetchQuickAddExercises (focus): Transformed quick add data:', JSON.stringify(quickAddData, null, 2));
            setQuickAddExercises(quickAddData);
            console.log('✅ fetchQuickAddExercises (focus): Successfully set quick add exercises');
          } else {
            console.log('⚠️ fetchQuickAddExercises (focus): No exercises array in response or invalid format');
            setQuickAddExercises([]);
          }
        } catch (error) {
          console.error('💥 fetchQuickAddExercises (focus): Error occurred:', error);
          console.error('💥 fetchQuickAddExercises (focus): Error stack:', error instanceof Error ? error.stack : 'No stack trace');
          setQuickAddExercises([]);
        }
      };

      fetchQuickAddExercises();
    }, [user?.id, isEditing])
  );

  const handleQuickAdd = async (exercise: RoutineExercise) => {
    const newExercises = [...exercises, exercise];
    setExercises(newExercises);
    
    if (isEditing) {
      // In edit mode, save only the newly added exercise to AsyncStorage
      // (it will be combined with existing routine exercises)
      try {
        const existingAsyncExercises = await AsyncStorage.getItem('routineExercises');
        const exercisesList = existingAsyncExercises ? JSON.parse(existingAsyncExercises) : [];
        
        // Check if this exercise is already in AsyncStorage
        const alreadyExists = exercisesList.some((ex: RoutineExercise) => ex.id === exercise.id);
        if (!alreadyExists) {
          exercisesList.push(exercise);
          await AsyncStorage.setItem('routineExercises', JSON.stringify(exercisesList));
        }
      } catch (error) {
        console.error('Error saving exercise to AsyncStorage in edit mode:', error);
      }
    } else {
      // In create mode, save all exercises to AsyncStorage
      await AsyncStorage.setItem('routineExercises', JSON.stringify(newExercises));
    }
  };

  const removeExercise = async (index: number) => {
    const exerciseToRemove = exercises[index];
    const newExercises = exercises.filter((_, i) => i !== index);
    setExercises(newExercises);
    
    if (isEditing) {
      // In edit mode, remove the exercise from AsyncStorage if it was newly added
      try {
        const existingAsyncExercises = await AsyncStorage.getItem('routineExercises');
        if (existingAsyncExercises) {
          const exercisesList = JSON.parse(existingAsyncExercises);
          const filteredList = exercisesList.filter((ex: RoutineExercise) => ex.id !== exerciseToRemove.id);
          await AsyncStorage.setItem('routineExercises', JSON.stringify(filteredList));
        }
      } catch (error) {
        console.error('Error removing exercise from AsyncStorage in edit mode:', error);
      }
    } else {
      // In create mode, save all remaining exercises to AsyncStorage
      await AsyncStorage.setItem('routineExercises', JSON.stringify(newExercises));
    }
  };

  const clearAllExercises = async () => {
    if (isEditing) {
      // In edit mode, only clear newly added exercises from AsyncStorage
      // but reset the exercises state to original routine exercises
      try {
        await AsyncStorage.removeItem('routineExercises');
        // Reload the original routine data
        if (routineId) {
          await fetchExistingRoutine(routineId);
        }
      } catch (error) {
        console.error('Error clearing AsyncStorage in edit mode:', error);
      }
    } else {
      // In create mode, clear everything
      setExercises([]);
      await AsyncStorage.removeItem('routineExercises');
    }
  };

  const handleSaveRoutine = async () => {
    const actionName = isEditing ? 'update' : 'create';
    console.log(`🚀 handleSaveRoutine: Starting routine ${actionName} process`);
    console.log('📝 handleSaveRoutine: Current state:', {
      routineName: routineName.trim(),
      exercisesCount: exercises.length,
      userId: user?.id,
      hasUser: !!user,
      userObject: user,
      userIdType: typeof user?.id,
      userIdLength: user?.id?.length,
      isEditing,
      routineId,
    });

    // Validate that we have a user ID (Clerk ID format: user_xxxxx)
    const isValidUserId = user?.id && typeof user.id === 'string' && user.id.length > 0;
    
    if (!routineName.trim() || exercises.length === 0 || !user?.id || !isValidUserId) {
      console.log(`❌ handleSaveRoutine: Validation failed:`, {
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

    console.log(`✅ handleSaveRoutine: Validation passed, proceeding with ${actionName}`);
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
          durationSeconds: set.durationSeconds || 0,
          distance: set.distance || 0,
          notes: set.notes || "",
        })),
        notes: "",
        restSeconds: 60,
      }));

      console.log('🔧 handleSaveRoutine: Transformed workout exercises:', JSON.stringify(workoutExercises, null, 2));

      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      
      let url: string;
      let method: string;
      let requestData: any;

      if (isEditing) {
        // Update existing routine
        requestData = {
          name: routineName,
          description: `Workout routine with ${exercises.length} exercises`,
          exercises: workoutExercises,
          notes: "",
        };
        console.log(`📦 handleSaveRoutine: Final ${actionName} data:`, JSON.stringify(requestData, null, 2));
        const result = await apiPut(`/api/workouts/${routineId}`, requestData);
        console.log(`✅ handleSaveRoutine: Routine ${actionName}d successfully:`, JSON.stringify(result, null, 2));
      } else {
        // Create new routine
        // Note: userId is no longer needed in request body - backend gets it from auth context
        requestData = {
          name: routineName,
          description: `Workout routine with ${exercises.length} exercises`,
          exercises: workoutExercises,
          startedAt: null,
          notes: "",
        };
        console.log(`📦 handleSaveRoutine: Final ${actionName} data:`, JSON.stringify(requestData, null, 2));
        const result = await apiPost('/api/workouts', requestData);
        console.log(`✅ handleSaveRoutine: Routine ${actionName}d successfully:`, JSON.stringify(result, null, 2));
      }

      // Clear AsyncStorage after successful save (both create and edit modes)
      console.log('🗑️ handleSaveRoutine: Clearing AsyncStorage');
      await AsyncStorage.removeItem('routineExercises');
      console.log('✅ handleSaveRoutine: AsyncStorage cleared successfully');
      
      Alert.alert(
        'Success',
        `Routine ${isEditing ? 'updated' : 'saved'} successfully!`,
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
      console.error(`💥 handleSaveRoutine: Error occurred during ${actionName}:`, error);
      console.error(`💥 handleSaveRoutine: Error stack:`, error instanceof Error ? error.stack : 'No stack trace');
      Alert.alert(
        'Error',
        `Failed to ${actionName} routine. Please try again.`
      );
    } finally {
      console.log('🏁 handleSaveRoutine: Setting loading to false');
      setLoading(false);
    }
  };

  const renderQuickAddPill = (exercise: RoutineExercise) => (
    <TouchableOpacity
      key={exercise.id}
      style={styles.quickAddPill}
      onPress={() => handleQuickAdd(exercise)}
      activeOpacity={0.7}
    >
      <Typography variant="label-small" color="contentPrimary">
        + {exercise.name}
      </Typography>
      <Typography variant="paragraph-xsmall" color="contentTertiary">
        {exercise.sets.length}×{exercise.sets[0].reps}
      </Typography>
    </TouchableOpacity>
  );

  const renderExerciseItem = (exercise: RoutineExercise, index: number) => (
    <View key={`${exercise.id}-${index}`} style={styles.exerciseItem}>
      <View style={styles.exerciseItemLeft}>
        <Typography variant="label-medium" color="contentPrimary">
          {exercise.name}
        </Typography>
        <View style={styles.exerciseMetaRow}>
          <View style={styles.metaPill}>
            <Typography variant="paragraph-xsmall" color="contentSecondary">
              {exercise.primaryMuscles?.join(", ") || "No muscle groups"}
            </Typography>
          </View>
          <Typography variant="paragraph-xsmall" color="contentTertiary">
            {exercise.sets.length} sets
          </Typography>
        </View>
      </View>
      <TouchableOpacity
        onPress={() => removeExercise(index)}
        style={styles.removeButton}
      >
        <Typography variant="label-small" color="contentNegative">
          ×
        </Typography>
      </TouchableOpacity>
    </View>
  );

  if (!user) {
    router.replace('/');
    return null;
  }

  // Show loading screen while fetching existing routine data
  if (initialDataLoading) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={getColor('contentAccent')} />
          <Typography variant="paragraph-medium" color="contentSecondary" style={styles.loadingText}>
            Loading routine...
          </Typography>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.contentWrapper}>
        {/* Minimal Header */}
        <View style={styles.header}>
          <TouchableOpacity onPress={() => router.back()} style={styles.backButton}>
            <Typography variant="label-medium" color="contentSecondary">
              ←
            </Typography>
          </TouchableOpacity>
          <Typography variant="label-medium" color="contentPrimary">
            {isEditing ? 'Edit Routine' : 'New Routine'}
          </Typography>
          <View style={styles.backButton} />
        </View>

        <ScrollView 
          style={styles.scrollView} 
          showsVerticalScrollIndicator={false}
          contentContainerStyle={styles.scrollContent}
        >
          {/* Routine Name Input */}
          <View style={styles.nameSection}>
            <Input
              value={routineName}
              onChangeText={setRoutineName}
              placeholder="Routine name"
              autoFocus={!isEditing}
            />
          </View>

          {/* Quick Add Section */}
          <View style={styles.quickAddSection}>
            <Typography variant="label-small" color="contentSecondary" style={styles.sectionLabel}>
              QUICK ADD
            </Typography>
            <ScrollView 
              horizontal 
              showsHorizontalScrollIndicator={false}
              style={styles.quickAddScroll}
            >
              {quickAddExercises.map(renderQuickAddPill)}
            </ScrollView>
          </View>

          {/* Current Exercises */}
          {exercises.length > 0 && (
            <View style={styles.exercisesSection}>
              <View style={styles.exercisesHeader}>
                <Typography variant="label-small" color="contentSecondary">
                  EXERCISES · {exercises.length}
                </Typography>
                <TouchableOpacity onPress={clearAllExercises}>
                  <Typography variant="label-small" color="contentTertiary">
                    Clear
                  </Typography>
                </TouchableOpacity>
              </View>
              <View style={styles.exercisesList}>
                {exercises.map(renderExerciseItem)}
              </View>
            </View>
          )}

          {/* Add Exercise Button */}
          <TouchableOpacity 
            style={styles.addExerciseButton}
            onPress={() => router.push('/exercise-search')}
            activeOpacity={0.7}
          >
            <Typography variant="label-medium" color="contentAccent">
              + Add Exercise
            </Typography>
          </TouchableOpacity>
        </ScrollView>

        {/* Fixed Bottom Action */}
        <View style={styles.bottomAction}>
          <Button
            variant="primary"
            size="large"
            style={styles.createButton}
            onPress={handleSaveRoutine}
            disabled={!routineName.trim() || exercises.length === 0 || loading}
          >
            {loading ? (isEditing ? "Updating..." : "Creating...") : (isEditing ? "Update" : "Create")}
          </Button>
        </View>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundSecondary'),
  },
  contentWrapper: {
    flex: 1,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    gap: 16,
  },
  loadingText: {
    textAlign: 'center',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: getColor('backgroundPrimary'),
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderTransparent'),
  },
  backButton: {
    width: 32,
    height: 32,
    alignItems: 'center',
    justifyContent: 'center',
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    paddingBottom: 100, // Space for fixed bottom button
  },
  nameSection: {
    backgroundColor: getColor('backgroundPrimary'),
    paddingHorizontal: 16,
    paddingVertical: 20,
    marginBottom: 2,
  },

  quickAddSection: {
    backgroundColor: getColor('backgroundPrimary'),
    paddingTop: 16,
    paddingBottom: 20,
  },
  sectionLabel: {
    paddingHorizontal: 16,
    marginBottom: 12,
    letterSpacing: 0.5,
  },
  quickAddScroll: {
    paddingHorizontal: 16,
  },
  quickAddPill: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    paddingHorizontal: 16,
    paddingVertical: 10,
    backgroundColor: getColor('backgroundTertiary'),
    borderRadius: 20,
    marginRight: 8,
  },
  exercisesSection: {
    backgroundColor: getColor('backgroundPrimary'),
    paddingVertical: 16,
    marginTop: 2,
  },
  exercisesHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    marginBottom: 12,
  },
  exercisesList: {
    paddingHorizontal: 16,
  },
  exerciseItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderTransparent'),
  },
  exerciseItemLeft: {
    flex: 1,
    gap: 4,
  },
  exerciseMetaRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  metaPill: {
    paddingHorizontal: 8,
    paddingVertical: 2,
    backgroundColor: getColor('backgroundTertiary'),
    borderRadius: 4,
  },
  removeButton: {
    width: 32,
    height: 32,
    alignItems: 'center',
    justifyContent: 'center',
  },
  addExerciseButton: {
    backgroundColor: getColor('backgroundPrimary'),
    paddingVertical: 20,
    alignItems: 'center',
    marginTop: 2,
  },
  bottomAction: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    padding: 16,
    backgroundColor: getColor('backgroundPrimary'),
    borderTopWidth: 1,
    borderTopColor: getColor('borderTransparent'),
  },
  createButton: {
    width: '100%',
  },
}); 