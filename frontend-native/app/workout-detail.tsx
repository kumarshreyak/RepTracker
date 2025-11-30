import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  ScrollView,
  ActivityIndicator,
  Pressable,
  TouchableOpacity,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { router, useLocalSearchParams } from 'expo-router';
import { Typography, Button, getColor, ExerciseHistory, type SemanticColor } from '../src/components';
import { useAuth } from '../src/hooks/useAuth';
import { MaterialIcons } from '@expo/vector-icons';

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
  exercise?: {
    id: string;
    name: string;
    muscleGroup: string;
    equipment: string;
    category: string;
  };
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
  rpeRating?: number;
}

export default function WorkoutDetailScreen() {
  const { user } = useAuth();
  const { id } = useLocalSearchParams<{ id: string }>();
  const [workout, setWorkout] = useState<WorkoutSession | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [expandedExercises, setExpandedExercises] = useState<Set<number>>(new Set());

  useEffect(() => {
    fetchWorkoutDetails();
  }, [id]);

  const fetchWorkoutDetails = async () => {
    if (!id) return;

    try {
      setLoading(true);
      setError(null);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(`${API_BASE_URL}/api/workout-sessions/${id}`);
      
      if (!response.ok) {
        throw new Error('Failed to fetch workout details');
      }
      
      const data = await response.json();
      
      // Use the correct property or the data directly if it contains the workout
      const workoutData = data.session || data.workout || data.workoutSession || (data.id ? data : null);
      setWorkout(workoutData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load workout');
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { 
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const formatTime = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleTimeString('en-US', { 
      hour: 'numeric',
      minute: '2-digit',
      hour12: true,
    });
  };

  const formatDuration = (seconds: number): string => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    
    if (hours > 0) {
      return `${hours}h ${minutes}m`;
    }
    return `${minutes}m`;
  };

  const calculateTotalVolume = (): number => {
    if (!workout) return 0;
    
    return workout.exercises.reduce((total, exercise) => {
      const exerciseVolume = exercise.sets.reduce((setTotal, set) => {
        if (set.completed) {
          return setTotal + (set.actualWeight * set.actualReps);
        }
        return setTotal;
      }, 0);
      return total + exerciseVolume;
    }, 0);
  };

  const calculateTotalSets = (): number => {
    if (!workout) return 0;
    
    return workout.exercises.reduce((total, exercise) => {
      return total + exercise.sets.filter(set => set.completed).length;
    }, 0);
  };

  const getRPEDescription = (rating: number): string => {
    if (rating <= 5) return 'Easy';
    if (rating <= 7) return 'Moderate';
    if (rating <= 8) return 'Hard';
    return 'Max Effort';
  };

  const getRPEColor = (rating: number): SemanticColor => {
    if (rating <= 5) return 'contentPositive';
    if (rating <= 7) return 'contentWarning';
    return 'contentNegative';
  };

  const getRPEBackgroundColor = (rating: number): string => {
    if (rating <= 5) return getColor('backgroundPositive');
    if (rating <= 7) return getColor('backgroundWarning');
    return getColor('backgroundNegative');
  };

  const toggleExerciseHistory = (exerciseIndex: number) => {
    const newExpandedExercises = new Set(expandedExercises);
    if (newExpandedExercises.has(exerciseIndex)) {
      newExpandedExercises.delete(exerciseIndex);
    } else {
      newExpandedExercises.add(exerciseIndex);
    }
    setExpandedExercises(newExpandedExercises);
  };

  if (loading) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={getColor('accent')} />
        </View>
      </SafeAreaView>
    );
  }

  if (error || !workout) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.errorContainer}>
          <Typography variant="paragraph-medium" color="contentNegative">
            {error || 'Workout not found'}
          </Typography>
          <Button
            variant="primary"
            size="default"
            onPress={() => router.back()}
            style={styles.errorButton}
          >
            Go Back
          </Button>
        </View>
      </SafeAreaView>
    );
  }

  const totalVolume = calculateTotalVolume();
  const totalSets = calculateTotalSets();
  const completedExercises = workout.exercises.filter(ex => ex.completed).length;

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Pressable onPress={() => router.back()} style={styles.backButton}>
          <Typography variant="label-medium" color="contentPrimary">
            ←
          </Typography>
        </Pressable>
        <Typography variant="label-medium" color="contentPrimary">
          Workout Details
        </Typography>
        <View style={styles.headerSpacer} />
      </View>

      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        {/* Workout Info Section */}
        <View style={styles.workoutInfoSection}>
          <Typography variant="label-large" color="contentPrimary" style={styles.workoutName}>
            {workout.name.replace(/ - \d+\/\d+\/\d+/, '')}
          </Typography>
          <Typography variant="paragraph-small" color="contentSecondary">
            {formatDate(workout.finishedAt || workout.startedAt)}
          </Typography>
          <Typography variant="paragraph-xsmall" color="contentTertiary" style={styles.timeInfo}>
            {formatTime(workout.startedAt)} - {workout.finishedAt ? formatTime(workout.finishedAt) : 'In Progress'}
          </Typography>
        </View>

        {/* Summary Stats */}
        <View style={styles.statsGrid}>
          <View style={styles.statCard}>
            <Typography variant="display-xsmall" color="contentPrimary">
              {formatDuration(workout.durationSeconds)}
            </Typography>
            <Typography variant="paragraph-xsmall" color="contentSecondary">
              Duration
            </Typography>
          </View>
          
          <View style={styles.statCard}>
            <Typography variant="display-xsmall" color="contentPrimary">
              {totalVolume.toLocaleString()}
            </Typography>
            <Typography variant="paragraph-xsmall" color="contentSecondary">
                                Volume (kg)
            </Typography>
          </View>
          
          <View style={styles.statCard}>
            <Typography variant="display-xsmall" color="contentPrimary">
              {totalSets}
            </Typography>
            <Typography variant="paragraph-xsmall" color="contentSecondary">
              Sets
            </Typography>
          </View>
          
          <View style={styles.statCard}>
            <Typography variant="display-xsmall" color="contentPrimary">
              {completedExercises}/{workout.exercises.length}
            </Typography>
            <Typography variant="paragraph-xsmall" color="contentSecondary">
              Exercises
            </Typography>
          </View>
        </View>

        {/* RPE Rating if available */}
        {workout.rpeRating && (
          <View style={styles.rpeSection}>
            <View style={styles.rpeHeader}>
              <Typography variant="label-small" color="contentPrimary">
                Workout Intensity
              </Typography>
              <Typography variant="label-medium" color={getRPEColor(workout.rpeRating)}>
                {getRPEDescription(workout.rpeRating)}
              </Typography>
            </View>
                         <View style={styles.rpeBar}>
               <View 
                 style={[
                   styles.rpeFill, 
                   { 
                     width: `${(workout.rpeRating / 10) * 100}%`,
                     backgroundColor: getRPEBackgroundColor(workout.rpeRating)
                   }
                 ]} 
               />
             </View>
            <Typography variant="paragraph-xsmall" color="contentTertiary" style={styles.rpeValue}>
              RPE {workout.rpeRating}/10
            </Typography>
          </View>
        )}

        {/* Exercises List */}
        <View style={styles.exercisesSection}>
          <Typography variant="label-medium" color="contentPrimary" style={styles.sectionTitle}>
            Exercises
          </Typography>
          
          {workout.exercises.map((exercise, index) => {
            const isHistoryExpanded = expandedExercises.has(index);
            return (
              <View key={index} style={styles.exerciseCard}>
                <View style={styles.exerciseHeader}>
                  <View style={styles.exerciseInfo}>
                    <Typography variant="heading-xsmall" color="contentPrimary">
                      {exercise.exercise?.name || `Exercise ${index + 1}`}
                    </Typography>
                    {exercise.exercise?.muscleGroup && (
                      <Typography variant="paragraph-xsmall" color="contentSecondary">
                        {exercise.exercise.muscleGroup}
                      </Typography>
                    )}
                  </View>
                  <View style={styles.exerciseHeaderRight}>
                    {exercise.completed && (
                      <View style={styles.completedBadge}>
                        <Typography variant="label-xsmall" color="contentPositive">
                          ✓
                        </Typography>
                      </View>
                    )}
                    {exercise.exerciseId && user?.id && (
                      <TouchableOpacity
                        style={styles.historyToggle}
                        onPress={() => toggleExerciseHistory(index)}
                      >
                        <MaterialIcons
                          name="history"
                          size={20}
                          color={getColor('contentTertiary')}
                        />
                        <MaterialIcons
                          name={isHistoryExpanded ? "expand-less" : "expand-more"}
                          size={16}
                          color={getColor('contentTertiary')}
                        />
                      </TouchableOpacity>
                    )}
                  </View>
                </View>
              
              {/* Sets */}
              <View style={styles.setsContainer}>
                {exercise.sets.map((set, setIndex) => (
                  <View key={setIndex} style={styles.setRow}>
                    <View style={styles.setNumber}>
                      <Typography variant="label-xsmall" color="contentSecondary">
                        {setIndex + 1}
                      </Typography>
                    </View>
                    
                    <View style={styles.setDetails}>
                      {set.completed ? (
                        <>
                          <Typography variant="paragraph-small" color="contentPrimary">
                            {set.actualWeight} kg × {set.actualReps} reps
                          </Typography>
                          {(set.targetWeight !== set.actualWeight || set.targetReps !== set.actualReps) && (
                            <Typography variant="paragraph-xsmall" color="contentTertiary">
                              Target: {set.targetWeight} kg × {set.targetReps} reps
                            </Typography>
                          )}
                        </>
                      ) : (
                        <Typography variant="paragraph-small" color="contentTertiary">
                          {set.targetWeight} kg × {set.targetReps} reps (skipped)
                        </Typography>
                      )}
                    </View>
                    
                    {set.completed && (
                      <View style={styles.setCheckmark}>
                        <Typography variant="label-xsmall" color="contentPositive">
                          ✓
                        </Typography>
                      </View>
                    )}
                  </View>
                ))}
              </View>
              
              {exercise.notes && (
                <View style={styles.exerciseNotes}>
                  <Typography variant="paragraph-xsmall" color="contentSecondary">
                    {exercise.notes}
                  </Typography>
                </View>
              )}

              {/* Exercise History */}
              {isHistoryExpanded && exercise.exerciseId && user?.id && (
                <View style={styles.exerciseHistoryContainer}>
                  <ExerciseHistory
                    exerciseId={exercise.exerciseId}
                    userId={user.id}
                    exerciseName={exercise.exercise?.name}
                  />
                </View>
              )}
            </View>
            );
          })}
        </View>

        {/* Workout Notes */}
        {workout.notes && (
          <View style={styles.notesSection}>
            <Typography variant="label-small" color="contentPrimary" style={styles.notesTitle}>
              Notes
            </Typography>
            <Typography variant="paragraph-small" color="contentSecondary">
              {workout.notes}
            </Typography>
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundSecondary'),
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 32,
  },
  errorButton: {
    marginTop: 16,
  },
  
  // Header
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: getColor('backgroundPrimary'),
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderOpaque'),
  },
  backButton: {
    padding: 8,
    marginLeft: -8,
  },
  headerSpacer: {
    width: 24,
  },
  
  scrollView: {
    flex: 1,
  },
  
  // Workout Info Section
  workoutInfoSection: {
    paddingHorizontal: 16,
    paddingTop: 24,
    paddingBottom: 16,
  },
  workoutName: {
    marginBottom: 8,
  },
  timeInfo: {
    marginTop: 4,
  },
  
  // Stats Grid
  statsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    paddingHorizontal: 12,
    marginBottom: 16,
  },
  statCard: {
    width: '50%',
    padding: 4,
  },
  statCardInner: {
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    padding: 16,
    alignItems: 'center',
  },
  
  // RPE Section
  rpeSection: {
    marginHorizontal: 16,
    marginBottom: 24,
    padding: 16,
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
  },
  rpeHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  rpeBar: {
    height: 8,
    backgroundColor: getColor('backgroundTertiary'),
    borderRadius: 4,
    overflow: 'hidden',
    marginBottom: 8,
  },
  rpeFill: {
    height: '100%',
    borderRadius: 4,
  },
  rpeValue: {
    textAlign: 'center',
  },
  
  // Exercises Section
  exercisesSection: {
    paddingHorizontal: 16,
    marginBottom: 24,
  },
  sectionTitle: {
    marginBottom: 16,
  },
  exerciseCard: {
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    padding: 16,
    marginBottom: 12,
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
  completedBadge: {
    width: 24,
    height: 24,
    borderRadius: 12,
    backgroundColor: getColor('backgroundLightPositive'),
    justifyContent: 'center',
    alignItems: 'center',
  },
  exerciseHeaderRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  historyToggle: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 8,
    borderRadius: 6,
    backgroundColor: getColor('backgroundSecondary'),
    gap: 4,
  },
  
  // Sets
  setsContainer: {
    marginBottom: 8,
  },
  setRow: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: getColor('borderTransparent'),
  },
  setNumber: {
    width: 24,
    height: 24,
    borderRadius: 12,
    backgroundColor: getColor('backgroundTertiary'),
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  setDetails: {
    flex: 1,
  },
  setCheckmark: {
    marginLeft: 12,
  },
  
  exerciseNotes: {
    marginTop: 8,
    paddingTop: 8,
    borderTopWidth: 1,
    borderTopColor: getColor('borderTransparent'),
  },
  exerciseHistoryContainer: {
    marginTop: 16,
    paddingTop: 16,
    borderTopWidth: 1,
    borderTopColor: getColor('borderOpaque'),
  },
  
  // Notes Section
  notesSection: {
    marginHorizontal: 16,
    marginBottom: 24,
    padding: 16,
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
  },
  notesTitle: {
    marginBottom: 8,
  },
}); 