import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  Alert,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native';
import { router, useFocusEffect } from 'expo-router';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { Typography, Button, Input } from '../src/components';
import { getColor } from '../src/components/Colors';
import { useAuth } from '../src/hooks/useAuth';
import { RoutineExercise, WorkoutSet } from '@/types/exercise';

// Quick add exercises will be fetched from the API

export default function CreateRoutineRoute() {
  const { user } = useAuth();
  const [routineName, setRoutineName] = useState("");
  const [exercises, setExercises] = useState<RoutineExercise[]>([]);
  const [quickAddExercises, setQuickAddExercises] = useState<RoutineExercise[]>([]);
  const [loading, setLoading] = useState(false);

  // Load exercises from AsyncStorage and fetch quick add exercises on mount
  useEffect(() => {
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

    const fetchQuickAddExercises = async () => {
      console.log('🏋️ fetchQuickAddExercises: Starting fetch process');
      console.log('👤 fetchQuickAddExercises: User object:', user);
      console.log('🆔 fetchQuickAddExercises: User ID:', user?.id);
      
      try {
        const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
        console.log('🌐 fetchQuickAddExercises: API_BASE_URL:', API_BASE_URL);
        
        const url = new URL(`${API_BASE_URL}/api/exercises/quick-add`);
        if (user?.id) {
          url.searchParams.append('userId', user.id);
          console.log('✅ fetchQuickAddExercises: Added userId to URL');
        } else {
          console.log('⚠️ fetchQuickAddExercises: No user ID available, fetching default exercises');
        }
        url.searchParams.append('limit', '5');
        
        const finalUrl = url.toString();
        console.log('📡 fetchQuickAddExercises: Making request to:', finalUrl);

        const response = await fetch(finalUrl);
        console.log('📥 fetchQuickAddExercises: Response received:', {
          status: response.status,
          statusText: response.statusText,
          ok: response.ok,
          headers: Object.fromEntries(response.headers.entries()),
        });

        if (response.ok) {
          const data = await response.json();
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
        } else {
          const errorText = await response.text();
          console.error('❌ fetchQuickAddExercises: Response not ok:', {
            status: response.status,
            statusText: response.statusText,
            errorBody: errorText,
          });
          setQuickAddExercises([]);
        }
      } catch (error) {
        console.error('💥 fetchQuickAddExercises: Error occurred:', error);
        console.error('💥 fetchQuickAddExercises: Error stack:', error instanceof Error ? error.stack : 'No stack trace');
        // Fallback to empty array if API fails
        setQuickAddExercises([]);
      }
    };

    loadExercises();
    fetchQuickAddExercises();
  }, [user?.id]);

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
            setExercises(normalizedExercises);
          } else {
            setExercises([]);
          }
        } catch (error) {
          console.error('Error loading exercises from AsyncStorage on focus:', error);
          setExercises([]);
        }
      };

      const fetchQuickAddExercises = async () => {
        console.log('🔄 fetchQuickAddExercises (focus): Starting fetch process');
        console.log('👤 fetchQuickAddExercises (focus): User object:', user);
        console.log('🆔 fetchQuickAddExercises (focus): User ID:', user?.id);
        
        try {
          const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
          console.log('🌐 fetchQuickAddExercises (focus): API_BASE_URL:', API_BASE_URL);
          
          const url = new URL(`${API_BASE_URL}/api/exercises/quick-add`);
          if (user?.id) {
            url.searchParams.append('userId', user.id);
            console.log('✅ fetchQuickAddExercises (focus): Added userId to URL');
          } else {
            console.log('⚠️ fetchQuickAddExercises (focus): No user ID available, fetching default exercises');
          }
          url.searchParams.append('limit', '5');
          
          const finalUrl = url.toString();
          console.log('📡 fetchQuickAddExercises (focus): Making request to:', finalUrl);

          const response = await fetch(finalUrl);
          console.log('📥 fetchQuickAddExercises (focus): Response received:', {
            status: response.status,
            statusText: response.statusText,
            ok: response.ok,
            headers: Object.fromEntries(response.headers.entries()),
          });

          if (response.ok) {
            const data = await response.json();
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
          } else {
            const errorText = await response.text();
            console.error('❌ fetchQuickAddExercises (focus): Response not ok:', {
              status: response.status,
              statusText: response.statusText,
              errorBody: errorText,
            });
            setQuickAddExercises([]);
          }
        } catch (error) {
          console.error('💥 fetchQuickAddExercises (focus): Error occurred:', error);
          console.error('💥 fetchQuickAddExercises (focus): Error stack:', error instanceof Error ? error.stack : 'No stack trace');
          setQuickAddExercises([]);
        }
      };

      loadExercises();
      fetchQuickAddExercises();
    }, [user?.id])
  );

  const handleQuickAdd = async (exercise: RoutineExercise) => {
    const newExercises = [...exercises, exercise];
    setExercises(newExercises);
    await AsyncStorage.setItem('routineExercises', JSON.stringify(newExercises));
  };

  const removeExercise = async (index: number) => {
    const newExercises = exercises.filter((_, i) => i !== index);
    setExercises(newExercises);
    await AsyncStorage.setItem('routineExercises', JSON.stringify(newExercises));
  };

  const clearAllExercises = async () => {
    setExercises([]);
    await AsyncStorage.removeItem('routineExercises');
  };

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
            New Routine
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
              autoFocus
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
            {loading ? "Creating..." : "Create"}
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