import React, { useState, useEffect, useRef } from 'react';
import {
  View,
  StyleSheet,
  ScrollView,
  Alert,
  Vibration,
  TouchableOpacity,
  Dimensions,
  TextInput,
  ActivityIndicator,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { router, useLocalSearchParams } from 'expo-router';
import { Typography, Button, NumberInput, ExerciseHistory } from '../src/components';
import { getColor } from '../src/components/Colors';
import { useAuth } from '../src/hooks/useAuth';
import { MaterialIcons } from '@expo/vector-icons';
import { Exercise, WorkoutSet, WorkoutExercise, ActiveWorkout } from '@/types/exercise';
import { apiGet, apiPost, apiPut } from '../src/utils/api';

const SCREEN_WIDTH = Dimensions.get('window').width;
const CARD_WIDTH = 260;
const CARD_MARGIN = 10;
const CARD_TOTAL_WIDTH = CARD_WIDTH + (CARD_MARGIN * 2);

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
  
  // UI state
  const [currentExerciseIndex, setCurrentExerciseIndex] = useState(0);
  const [currentSetIndex, setCurrentSetIndex] = useState(0);
  const [isPaused, setIsPaused] = useState(false);
  const [showRPEModal, setShowRPEModal] = useState(false);
  const [rpeRating, setRpeRating] = useState<number | null>(null);
  const setsScrollRef = useRef<ScrollView>(null);

  // Timer effect
  useEffect(() => {
    if (!isPaused) {
      intervalRef.current = setInterval(() => {
        setElapsedTime(Date.now() - startTime.getTime());
      }, 1000);
    } else {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    }

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [startTime, isPaused]);

  // Load routine data and start workout
  useEffect(() => {
    if (routineId && user?.id) {
      fetchRoutineAndStartWorkout();
    }
  }, [routineId, user?.id]);

  // Reset sets scroll position when exercise changes
  useEffect(() => {
    if (setsScrollRef.current) {
      setsScrollRef.current.scrollTo({ x: 0, animated: true });
      setCurrentSetIndex(0);
    }
  }, [currentExerciseIndex]);

  const scrollToNextSetOrExercise = () => {
    if (!activeWorkout) return;

    const currentExercise = activeWorkout.exercises[currentExerciseIndex];
    const isLastSet = currentSetIndex >= currentExercise.sets.length - 1;
    const isLastExercise = currentExerciseIndex >= activeWorkout.exercises.length - 1;

    if (isLastSet && !isLastExercise) {
      // Move to next exercise
      setCurrentExerciseIndex(prev => prev + 1);
      // setCurrentSetIndex(0) will be handled by the useEffect above
    } else if (!isLastSet) {
      // Move to next set in current exercise
      const nextSetIndex = currentSetIndex + 1;
      setCurrentSetIndex(nextSetIndex);
      
      // Scroll to the next set
      if (setsScrollRef.current) {
        const scrollX = nextSetIndex * CARD_TOTAL_WIDTH;
        setsScrollRef.current.scrollTo({ x: scrollX, animated: true });
      }
    }
    // If it's the last set of the last exercise, don't auto-scroll
  };

  const handleScroll = (event: any) => {
    const scrollX = event.nativeEvent.contentOffset.x;
    const newSetIndex = Math.round(scrollX / CARD_TOTAL_WIDTH);
    if (newSetIndex !== currentSetIndex && newSetIndex >= 0) {
      setCurrentSetIndex(newSetIndex);
    }
  };

  const fetchRoutineAndStartWorkout = async () => {
    try {
      setLoading(true);
      setError(null);

      const routineData = await apiGet<any>(`/api/workouts/${routineId}/start`);
      
      // Note: userId is no longer needed in request body - backend gets it from auth context
      const sessionData = {
        routineId: routineData.id,
        name: `${routineData.name} - ${new Date().toLocaleDateString()}`,
        description: `Active workout session for ${routineData.name}`,
        notes: '',
      };

      const sessionResult = await apiPost<any>('/api/workout-sessions', sessionData);
      
      const workout: ActiveWorkout = {
        sessionId: sessionResult.id,
        routineId: routineData.id,
        routineName: routineData.name,
        routineDescription: routineData.description,
        exercises: sessionResult.exercises.map((ex: any) => ({
          ...ex,
          completed: false,
          sets: ex.sets.map((set: any) => ({
            targetReps: set.targetReps || set.reps,
            targetWeight: set.targetWeight || set.weight,
            actualReps: set.actualReps || set.targetReps || set.reps,
            actualWeight: set.actualWeight || set.targetWeight || set.weight,
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

  const handleSetComplete = async (exerciseIndex: number, setIndex: number) => {
    if (!activeWorkout || !activeWorkout.sessionId) return;

    const currentSet = activeWorkout.exercises[exerciseIndex].sets[setIndex];
    const newCompletedState = !currentSet.completed;

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

      const allSetsCompleted = updated.exercises[exerciseIndex].sets.every(set => set.completed);
      updated.exercises[exerciseIndex].completed = allSetsCompleted;

      return updated;
    });

    Vibration.vibrate(newCompletedState ? [0, 100, 50, 100] : 50);

    // Auto-scroll to next set/exercise when completing a set
    if (newCompletedState && setIndex === currentSetIndex) {
      setTimeout(() => {
        scrollToNextSetOrExercise();
      }, 300); // Small delay to let the completion animation finish
    }

    try {
      const updateData = {
        actualReps: currentSet.actualReps,
        actualWeight: currentSet.actualWeight,
        durationSeconds: currentSet.durationSeconds || 0,
        distance: currentSet.distance || 0,
        notes: currentSet.notes || '',
        completed: newCompletedState,
      };

      await apiPut(
        `/api/workout-sessions/${activeWorkout.sessionId}/exercises/${exerciseIndex}/sets/${setIndex}`,
        updateData
      );
    } catch (error) {
      console.error('Error updating set:', error);
    }
  };

  const handleValueEdit = (exerciseIndex: number, setIndex: number, field: 'actualReps' | 'actualWeight', value: string) => {
    const numValue = parseInt(value) || 0;
    
    setActiveWorkout(prev => {
      if (!prev) return prev;

      const updated = { ...prev };
      updated.exercises = [...prev.exercises];
      updated.exercises[exerciseIndex] = { ...prev.exercises[exerciseIndex] };
      updated.exercises[exerciseIndex].sets = [...prev.exercises[exerciseIndex].sets];
      updated.exercises[exerciseIndex].sets[setIndex] = {
        ...prev.exercises[exerciseIndex].sets[setIndex],
        [field]: numValue,
      };

      return updated;
    });
  };

  const handleEndWorkout = async () => {
    Alert.alert(
      'End Workout?',
      'Your progress will be saved.',
      [
        { text: 'Cancel', style: 'cancel' },
        { text: 'End Workout', style: 'destructive', onPress: confirmEndWorkout },
      ]
    );
  };

  const confirmEndWorkout = async () => {
    if (!activeWorkout || !user?.id) return;

    // Show RPE modal instead of immediately saving
    setShowRPEModal(true);
  };

  const saveWorkoutWithRPE = async (selectedRPE: number | null) => {
    if (!activeWorkout || !user?.id) return;

    try {
      setSaving(true);
      setShowRPEModal(false);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      const finishedAt = new Date();
      const durationSeconds = Math.floor((finishedAt.getTime() - startTime.getTime()) / 1000);

      if (!activeWorkout.sessionId) {
        throw new Error('No active workout session found');
      }

      const updateData = {
        exercises: activeWorkout.exercises.map(ex => ({
          exerciseId: ex.exerciseId,
          sets: ex.sets.map(set => ({
            targetReps: set.targetReps,
            targetWeight: set.targetWeight,
            actualReps: set.actualReps,
            actualWeight: set.actualWeight,
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
        rpeRating: selectedRPE || 0,
        notes: activeWorkout.notes || `Completed ${activeWorkout.exercises.filter(ex => ex.completed).length}/${activeWorkout.exercises.length} exercises`,
      };

      await apiPut(`/api/workout-sessions/${activeWorkout.sessionId}`, updateData);

      // Apply progressive overload after successfully updating the session
      if (activeWorkout.routineId) {
        try {
          await apiPost(
            `/api/workout-sessions/${activeWorkout.sessionId}/progressive-overload/ai`,
            { workoutId: activeWorkout.routineId }
          );
          console.log('Progressive overload applied successfully');
        } catch (progressiveOverloadError) {
          console.warn('Progressive overload error:', progressiveOverloadError);
          // Don't throw here - workout was saved successfully, progressive overload is optional
        }
      }

      router.replace('/(tabs)');
      
      setTimeout(() => {
        Alert.alert(
          'Workout Completed!',
          `Great job! You worked out for ${formatTime(elapsedTime)}.`,
          [{ text: 'OK' }]
        );
      }, 500);

    } catch (err) {
      console.error('Error saving workout session:', err);
      Alert.alert('Error', 'Failed to save workout session. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handleCancelWorkout = () => {
    Alert.alert(
      'Cancel Workout?',
      'Your progress will be lost.',
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
          <Typography variant="paragraph-medium" color="contentSecondary">
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
          <Typography variant="paragraph-medium" color="contentNegative" style={styles.errorText}>
            {error || 'Failed to load workout'}
          </Typography>
          <Button variant="primary" size="default" onPress={() => router.back()}>
            Go Back
          </Button>
        </View>
      </SafeAreaView>
    );
  }

  const currentExercise = activeWorkout.exercises[currentExerciseIndex];
  const completedExercises = activeWorkout.exercises.filter(ex => ex.completed).length;
  const progressPercentage = ((completedExercises + (currentExercise?.sets.filter(set => set.completed).length || 0) / (currentExercise?.sets.length || 1)) / activeWorkout.exercises.length) * 100;

  return (
    <SafeAreaView style={styles.container}>
      {/* Fixed Header - Minimal Timer Bar */}
      <View style={styles.timerBar}>
        <View style={styles.timerLeft}>
          <Typography variant="heading-xsmall" color="contentPrimary" style={styles.timerText}>
            {formatTime(elapsedTime)}
          </Typography>
        </View>
        
        <View style={styles.timerRight}>
          <TouchableOpacity 
            style={styles.pauseButton} 
            onPress={() => setIsPaused(!isPaused)}
          >
            <MaterialIcons 
              name={isPaused ? "play-arrow" : "pause"} 
              size={24} 
              color={getColor('contentTertiary')} 
            />
          </TouchableOpacity>
          <TouchableOpacity style={styles.endButton} onPress={handleCancelWorkout}>
            <MaterialIcons 
              name="close" 
              size={24} 
              color={getColor('contentNegative')} 
            />
          </TouchableOpacity>
        </View>
      </View>

      {/* Progress Indicator - Visual, No Text */}
      <View style={styles.progressBar}>
        <View style={[styles.progressFill, { width: `${progressPercentage}%` }]} />
      </View>

      {/* Main Content - Vertically Scrollable */}
      <ScrollView 
        style={styles.mainScrollView}
        showsVerticalScrollIndicator={false}
        contentContainerStyle={styles.mainScrollContent}
      >
        {/* Current Exercise Focus Area */}
        <View style={styles.currentExerciseContainer}>
          {/* Exercise Header - Minimal */}
          <View style={styles.exerciseHeader}>
            <View>
              <Typography variant="label-medium" color="contentPrimary">
                {currentExercise.exercise?.name || `Exercise ${currentExerciseIndex + 1}`}
              </Typography>
              <Typography variant="paragraph-xsmall" color="contentTertiary">
                {currentExercise.exercise?.primaryMuscles?.join(", ") || ''}
              </Typography>
            </View>
            
            {/* Visual Set Counter */}
            <View style={styles.setCounter}>
              {currentExercise.sets.map((set, index) => (
                <View key={index} style={styles.setIndicator}>
                  {set.completed ? (
                    <MaterialIcons 
                      name="check-circle" 
                      size={16} 
                      color={getColor('positive')} 
                    />
                  ) : (
                    <MaterialIcons 
                      name="radio-button-unchecked" 
                      size={16} 
                      color={getColor('backgroundTertiary')} 
                    />
                  )}
                </View>
              ))}
            </View>
          </View>

          {/* Active Sets - Carousel Based, Editable */}
          <ScrollView 
            ref={setsScrollRef}
            horizontal 
            showsHorizontalScrollIndicator={false}
            style={styles.setsScroller}
            contentContainerStyle={[
              styles.setsScrollContent,
              { paddingHorizontal: (SCREEN_WIDTH - CARD_WIDTH) / 2 }
            ]}
            snapToInterval={CARD_TOTAL_WIDTH}
            decelerationRate="fast"
            snapToAlignment="start"
            pagingEnabled={false}
            onScroll={handleScroll}
            scrollEventThrottle={16}
          >
            {currentExercise.sets.map((set, setIndex) => (
              <View key={setIndex} style={styles.setCard}>
                {/* Set Number - Small Label */}
                <Typography variant="label-xsmall" color="contentTertiary" style={styles.setLabel}>
                  SET {setIndex + 1}
                </Typography>
                
                {/* Editable Values - Large, Tappable */}
                <View style={styles.setValues}>
                  <NumberInput
                    value={set.actualReps}
                    onValueChange={(value) => handleValueEdit(currentExerciseIndex, setIndex, 'actualReps', value.toString())}
                    unit="reps"
                    displayVariant="display-xsmall"
                    unitVariant="label-small"
                    valueColor="contentPrimary"
                    unitColor="contentTertiary"
                    containerStyle={styles.valueGroup}
                    minValue={0}
                    maxValue={999}
                  />

                  <View style={styles.valueSeparator} />

                  <NumberInput
                    value={set.actualWeight}
                    onValueChange={(value) => handleValueEdit(currentExerciseIndex, setIndex, 'actualWeight', value.toString())}
                    unit="kg"
                    displayVariant="display-xsmall"
                    unitVariant="label-small"
                    valueColor="contentPrimary"
                    unitColor="contentTertiary"
                    containerStyle={styles.valueGroup}
                    minValue={0}
                    maxValue={9999}
                    allowDecimal={true}
                  />
                </View>

                {/* Complete Set Button - Large, Primary Action */}
                <TouchableOpacity
                  style={[
                    styles.completeSetButton,
                    set.completed && styles.completeSetButtonDone
                  ]}
                  onPress={() => handleSetComplete(currentExerciseIndex, setIndex)}
                >
                  {set.completed ? (
                    <MaterialIcons 
                      name="check" 
                      size={20} 
                      color={getColor("contentPositive")} 
                    />
                  ) : (
                    <Typography 
                      variant="label-medium" 
                      color="contentPrimary"
                    >
                      DONE
                    </Typography>
                  )}
                </TouchableOpacity>
              </View>
            ))}
          </ScrollView>

        </View>

        {/* Exercise Chips - All Exercises */}
        <View style={styles.exerciseChipsSection}>
          <ScrollView horizontal showsHorizontalScrollIndicator={false} contentContainerStyle={styles.exerciseChipsContainer}>
            {activeWorkout.exercises.map((exercise, index) => (
              <TouchableOpacity
                key={index}
                style={[
                  styles.exerciseChip,
                  index === currentExerciseIndex && styles.exerciseChipSelected
                ]}
                onPress={() => setCurrentExerciseIndex(index)}
              >
                <View style={styles.exerciseChipContent}>
                  <Typography 
                    variant="label-small" 
                    color={index === currentExerciseIndex ? "contentOnColor" : "contentPrimary"}
                    style={styles.exerciseChipName}
                  >
                    {exercise.exercise?.name || `Exercise ${index + 1}`}
                  </Typography>
                  <View style={styles.exerciseChipMeta}>
                    <Typography 
                      variant="paragraph-xsmall" 
                      color={index === currentExerciseIndex ? "contentOnColor" : "contentSecondary"}
                    >
                      {exercise.sets.length} sets
                    </Typography>
                    {exercise.completed && (
                      <MaterialIcons 
                        name="check-circle" 
                        size={12} 
                        color={getColor(index === currentExerciseIndex ? "contentOnColor" : "positive")} 
                        style={styles.exerciseChipCheck}
                      />
                    )}
                  </View>
                </View>
              </TouchableOpacity>
            ))}
          </ScrollView>
        </View>

        {/* Exercise History - For Current Exercise */}
        {currentExercise.exerciseId && user?.id && (
          <ExerciseHistory 
            exerciseId={currentExercise.exerciseId}
            exerciseName={currentExercise.exercise?.name}
          />
        )}
      </ScrollView>

      {/* Floating Action Button - End Workout */}
      <TouchableOpacity 
        style={styles.fab} 
        onPress={handleEndWorkout}
        disabled={saving}
      >
        <View style={styles.fabContent}>
          {!saving && (
            <MaterialIcons 
              name="flag" 
              size={20} 
              color={getColor("contentOnColor")} 
              style={styles.fabIcon}
            />
          )}
          <Typography variant="label-medium" color="contentOnColor">
            {saving ? "SAVING..." : "FINISH"}
          </Typography>
        </View>
      </TouchableOpacity>

      {/* RPE Rating Modal */}
      {showRPEModal && (
        <View style={styles.modalOverlay}>
          <View style={styles.rpeModal}>
            <Typography variant="heading-small" color="contentPrimary" style={styles.rpeTitle}>
              Rate workout difficulty
            </Typography>
            <Typography variant="paragraph-medium" color="contentSecondary" style={styles.rpeSubtitle}>
              How hard was this workout? (1-10 scale)
            </Typography>
            
            <View style={styles.rpeButtonGrid}>
              {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((rating) => (
                <TouchableOpacity
                  key={rating}
                  style={[
                    styles.rpeButton,
                    rpeRating === rating && styles.rpeButtonSelected
                  ]}
                  onPress={() => setRpeRating(rating)}
                >
                  <Typography 
                    variant="heading-small" 
                    color={rpeRating === rating ? "contentOnColor" : "contentPrimary"}
                  >
                    {rating}
                  </Typography>
                </TouchableOpacity>
              ))}
            </View>

            <View style={styles.rpeButtonRow}>
              <TouchableOpacity 
                style={styles.rpeSkipButton} 
                onPress={() => saveWorkoutWithRPE(null)}
              >
                <Typography variant="label-medium" color="contentSecondary">
                  Skip
                </Typography>
              </TouchableOpacity>
              
              <TouchableOpacity 
                style={[
                  styles.rpeSaveButton,
                  !rpeRating && styles.rpeSaveButtonDisabled
                ]} 
                onPress={() => saveWorkoutWithRPE(rpeRating)}
                disabled={!rpeRating}
              >
                <Typography variant="label-medium" color="contentOnColor">
                  {saving ? "SAVING..." : "SAVE RATING"}
                </Typography>
              </TouchableOpacity>
            </View>
          </View>
        </View>
      )}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundSecondary'),
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
  
  // Timer Bar - Minimal, Always Visible
  timerBar: {
    backgroundColor: getColor('backgroundPrimary'),
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderOpaque'),
  },
  timerLeft: {
    flex: 1,
  },
  timerText: {
    fontFamily: 'monospace',
  },
  timerRight: {
    flexDirection: 'row',
    gap: 16,
  },
  pauseButton: {
    padding: 8,
  },
  endButton: {
    padding: 8,
  },
  
  // Progress Bar - Visual Only
  progressBar: {
    height: 4,
    backgroundColor: getColor('backgroundTertiary'),
  },
  progressFill: {
    height: '100%',
    backgroundColor: getColor('accent'),
  },
  
  // Main Scroll View
  mainScrollView: {
    flex: 1,
  },
  mainScrollContent: {
    paddingBottom: 100, // Space for floating action button
  },
  
  // Current Exercise Area
  currentExerciseContainer: {
    backgroundColor: getColor('backgroundPrimary'),
    marginTop: 2,
  },
  exerciseHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 20,
    paddingTop: 20,
    paddingBottom: 16,
  },
  setCounter: {
    flexDirection: 'row',
    gap: 6,
  },
  setIndicator: {
    alignItems: 'center',
    justifyContent: 'center',
  },
  
  // Sets Scroller - Carousel Based
  setsScroller: {
    height: 300, // Fixed height to prevent unnecessary expansion
  },
  setsScrollContent: {
    // Padding is now dynamically set based on screen width
  },
  setCard: {
    width: CARD_WIDTH,
    height: 260, // Fixed height for consistent card sizing
    marginHorizontal: CARD_MARGIN, // Space between cards for carousel effect
    paddingHorizontal: 20,
    paddingVertical: 24,
    alignItems: 'center',
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    borderRadius: 8,
    backgroundColor: getColor('backgroundPrimary'),
  },
  setLabel: {
    marginBottom: 24,
  },
  setValues: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 40,
  },
  valueGroup: {
    alignItems: 'center',
    paddingHorizontal: 32,
  },
  valueSeparator: {
    width: 1,
    height: 60,
    backgroundColor: getColor('borderOpaque'),
  },
  
  // Complete Button - Primary Action
  completeSetButton: {
    backgroundColor: getColor('backgroundPrimary'),
    borderWidth: 3,
    borderColor: getColor('contentPrimary'),
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 32,
  },
  completeSetButtonDone: {
    backgroundColor: getColor('backgroundLightPositive'),
    borderColor: getColor('positive'),
  },
  
  // Exercise Chips Section
  exerciseChipsSection: {
    paddingVertical: 16,
    borderTopWidth: 1,
    borderTopColor: getColor('borderOpaque'),
  },
  exerciseChipsContainer: {
    paddingHorizontal: 16,
    gap: 8,
  },
  exerciseChip: {
    backgroundColor: getColor('backgroundPrimary'),
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    borderRadius: 8,
    paddingHorizontal: 12,
    paddingVertical: 8,
    minWidth: 120,
  },
  exerciseChipSelected: {
    backgroundColor: getColor('backgroundAccent'),
    borderColor: getColor('backgroundAccent'),
  },
  exerciseChipContent: {
    alignItems: 'center',
  },
  exerciseChipName: {
    marginBottom: 2,
    textAlign: 'center',
  },
  exerciseChipMeta: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 4,
  },
  exerciseChipCheck: {
    marginLeft: 2,
  },
  
  // Floating Action Button
  fab: {
    position: 'absolute',
    bottom: 24,
    right: 20,
    backgroundColor: getColor('backgroundAccent'),
    paddingHorizontal: 32,
    paddingVertical: 16,
    borderRadius: 28,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 8,
    elevation: 8,
  },
  fabContent: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  fabIcon: {
    marginRight: 8,
  },
  
  // RPE Modal Styles
  modalOverlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 1000,
  },
  rpeModal: {
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 16,
    paddingHorizontal: 24,
    paddingVertical: 32,
    marginHorizontal: 20,
    minWidth: 300,
    maxWidth: 400,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.3,
    shadowRadius: 16,
    elevation: 16,
  },
  rpeTitle: {
    textAlign: 'center',
    marginBottom: 8,
  },
  rpeSubtitle: {
    textAlign: 'center',
    marginBottom: 32,
  },
  rpeButtonGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'center',
    gap: 12,
    marginBottom: 32,
  },
  rpeButton: {
    width: 48,
    height: 48,
    backgroundColor: getColor('backgroundPrimary'),
    borderWidth: 2,
    borderColor: getColor('borderOpaque'),
    borderRadius: 8,
    justifyContent: 'center',
    alignItems: 'center',
  },
  rpeButtonSelected: {
    backgroundColor: getColor('backgroundAccent'),
    borderColor: getColor('backgroundAccent'),
  },
  rpeButtonRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: 16,
  },
  rpeSkipButton: {
    flex: 1,
    paddingVertical: 16,
    paddingHorizontal: 24,
    backgroundColor: getColor('backgroundSecondary'),
    borderRadius: 8,
    alignItems: 'center',
  },
  rpeSaveButton: {
    flex: 2,
    paddingVertical: 16,
    paddingHorizontal: 24,
    backgroundColor: getColor('backgroundAccent'),
    borderRadius: 8,
    alignItems: 'center',
  },
  rpeSaveButtonDisabled: {
    backgroundColor: getColor('backgroundTertiary'),
  },
}); 