import React, { useState, useEffect } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  TextInput,
  TouchableOpacity,
  Alert,
  ActivityIndicator,
} from 'react-native';
import { router } from 'expo-router';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { Typography, Button } from '../src/components';
import { getColor } from '../src/components/Colors';
import { useAuth } from '../src/hooks/useAuth';
import { User } from '../src/auth/AuthService';

interface Exercise {
  id: string;
  name: string;
  description: string;
  muscle_group: string;
  equipment: string;
  difficulty: string;
  instructions: string[];
}

interface WorkoutSet {
  reps: number;
  weight: number;
}

interface WorkoutSetInput {
  reps: string;
  weight: string;
}

interface RoutineExercise {
  id: string;
  name: string;
  muscle_group: string;
  sets: WorkoutSet[];
}

export default function ExerciseSearchRoute() {
  const { user } = useAuth();
  const [searchQuery, setSearchQuery] = useState("");
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedExercise, setSelectedExercise] = useState<Exercise | null>(null);
  const [numberOfSets, setNumberOfSets] = useState("3");
  const [sets, setSets] = useState<WorkoutSetInput[]>([
    { reps: "10", weight: "" },
    { reps: "10", weight: "" },
    { reps: "10", weight: "" },
  ]);

  // Fetch exercises
  const fetchExercises = async (query?: string) => {
    setLoading(true);
    try {
      const params = new URLSearchParams();
      if (query) params.append("search", query);
      
      const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(`${API_BASE_URL}/api/exercises?${params.toString()}`);
      
      if (!response.ok) throw new Error("Failed to fetch exercises");
      
      const data = await response.json();
      setExercises(data.exercises || []);
    } catch (error) {
      console.error("Error fetching exercises:", error);
      setExercises([]);
      Alert.alert('Error', 'Failed to fetch exercises. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  // Load exercises on mount and when search changes
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      fetchExercises(searchQuery);
    }, 300); // Debounce search

    return () => clearTimeout(timeoutId);
  }, [searchQuery]);

  // Update sets array when numberOfSets changes
  useEffect(() => {
    const setsCount = parseInt(numberOfSets) || 1;
    const newSets = Array.from({ length: setsCount }, (_, index) => 
      sets[index] || { reps: "10", weight: "" }
    );
    setSets(newSets);
  }, [numberOfSets]);

  const updateSet = (index: number, field: 'reps' | 'weight', value: string) => {
    const newSets = [...sets];
    newSets[index] = { ...newSets[index], [field]: value };
    setSets(newSets);
  };

  const isFormValid = () => {
    return selectedExercise && 
           sets.length > 0 && 
           sets.every(set => {
             const reps = parseFloat(set.reps) || 0;
             const weight = parseFloat(set.weight) || 0;
             return reps > 0 && weight > 0;
           });
  };

  const handleAddExercise = async () => {
    if (!selectedExercise || !isFormValid()) return;
    
    // Convert string inputs to numbers for storage
    const convertedSets: WorkoutSet[] = sets.map(set => ({
      reps: parseFloat(set.reps) || 0,
      weight: parseFloat(set.weight) || 0,
    }));
    
    const exerciseData: RoutineExercise = {
      id: selectedExercise.id,
      name: selectedExercise.name,
      muscle_group: selectedExercise.muscle_group,
      sets: convertedSets,
    };
    
    try {
      // Save to AsyncStorage
      const existingExercises = await AsyncStorage.getItem('routineExercises');
      const exercisesList = existingExercises ? JSON.parse(existingExercises) : [];
      exercisesList.push(exerciseData);
      await AsyncStorage.setItem('routineExercises', JSON.stringify(exercisesList));
      
      console.log("Adding exercise to routine:", exerciseData);
      
      Alert.alert(
        'Success',
        'Exercise added to routine!',
        [
          {
            text: 'Add Another',
            onPress: () => {
              setSelectedExercise(null);
              setNumberOfSets("3");
              setSets([
                { reps: "10", weight: "" },
                { reps: "10", weight: "" },
                { reps: "10", weight: "" },
              ]);
            },
          },
          {
            text: 'Back to Routine',
            onPress: () => router.back(),
            style: 'default',
          },
        ]
      );
    } catch (error) {
      console.error('Error saving exercise:', error);
      Alert.alert('Error', 'Failed to add exercise. Please try again.');
    }
  };

  const handleBackToSearch = () => {
    setSelectedExercise(null);
    setNumberOfSets("3");
    setSets([
      { reps: "10", weight: "" },
      { reps: "10", weight: "" },
      { reps: "10", weight: "" },
    ]);
  };

  if (!user) {
    router.replace('/');
    return null;
  }

  const renderExerciseItem = (exercise: Exercise) => (
    <TouchableOpacity
      key={exercise.id}
      style={[
        styles.exerciseItem,
        selectedExercise?.id === exercise.id && styles.exerciseItemSelected
      ]}
      onPress={() => setSelectedExercise(exercise)}
    >
      <Typography variant="label-medium" color="text-primary" style={styles.exerciseName}>
        {exercise.name}
      </Typography>
      <Typography variant="paragraph-small" color="text-secondary" style={styles.exerciseDetails}>
        {exercise.muscle_group} • {exercise.equipment} • {exercise.difficulty}
      </Typography>
    </TouchableOpacity>
  );

  const renderSetConfiguration = (set: WorkoutSetInput, index: number) => (
    <View key={index} style={styles.setCard}>
      <Typography variant="paragraph-small" color="text-primary" style={styles.setTitle}>
        Set {index + 1}
      </Typography>
      <View style={styles.setInputs}>
        <View style={styles.setInputContainer}>
          <Typography variant="paragraph-small" color="text-primary" style={styles.inputLabel}>
            Reps *
          </Typography>
          <TextInput
            value={set.reps}
            onChangeText={(text) => updateSet(index, 'reps', text)}
            keyboardType="numeric"
            style={styles.setInput}
            placeholderTextColor={getColor('text-secondary')}
          />
        </View>
        <View style={styles.setInputContainer}>
          <Typography variant="paragraph-small" color="text-primary" style={styles.inputLabel}>
            Weight (kg) *
          </Typography>
          <TextInput
            value={set.weight}
            onChangeText={(text) => updateSet(index, 'weight', text)}
            keyboardType="numeric"
            style={styles.setInput}
            placeholderTextColor={getColor('text-secondary')}
          />
        </View>
      </View>
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        {/* Page Header */}
        <View style={styles.header}>
          <Button variant="text" size="default" onPress={() => router.back()}>
            ← Back to Create Routine
          </Button>
          <Typography variant="heading-xxlarge" color="text-primary" style={styles.title}>
            Add Exercise
          </Typography>
          <Typography variant="paragraph-medium" color="text-secondary">
            {selectedExercise ? 'Configure your exercise sets' : 'Search and select exercises for your routine'}
          </Typography>
        </View>

        {/* Mode 1: Search Mode - Only show when no exercise is selected */}
        {!selectedExercise && (
          <>
            {/* Search Section */}
            <View style={styles.section}>
              <Typography variant="heading-small" color="text-primary" style={styles.sectionTitle}>
                Search Exercises
              </Typography>
              <TextInput
                value={searchQuery}
                onChangeText={setSearchQuery}
                placeholder="Search by name, muscle group, or equipment..."
                style={styles.searchInput}
                placeholderTextColor={getColor('text-secondary')}
              />
            </View>

            {/* Exercise List */}
            <View style={styles.section}>
              <Typography variant="heading-small" color="text-primary" style={styles.sectionTitle}>
                Exercise Library
              </Typography>

              {loading ? (
                <View style={styles.loadingContainer}>
                  <ActivityIndicator size="large" color={getColor('primary')} />
                  <Typography variant="paragraph-medium" color="text-secondary" style={styles.loadingText}>
                    Loading exercises...
                  </Typography>
                </View>
              ) : exercises.length === 0 ? (
                <View style={styles.emptyState}>
                  <Typography variant="paragraph-medium" color="text-secondary" style={styles.emptyStateTitle}>
                    No exercises found
                  </Typography>
                  {searchQuery && (
                    <Typography variant="paragraph-small" color="text-secondary" style={styles.emptyStateSubtitle}>
                      Try a different search term
                    </Typography>
                  )}
                </View>
              ) : (
                <View style={styles.exerciseList}>
                  {exercises.map(renderExerciseItem)}
                </View>
              )}
            </View>
          </>
        )}

        {/* Mode 2: Exercise Configuration Mode - Only show when exercise is selected */}
        {selectedExercise && (
          <View style={styles.section}>
            {/* Selected Exercise Header with Cross Button */}
            <View style={styles.selectedExerciseHeader}>
              <View style={styles.selectedExerciseInfo}>
                <Typography variant="heading-small" color="text-primary" style={styles.selectedExerciseName}>
                  {selectedExercise.name}
                </Typography>
                <Typography variant="paragraph-small" color="text-secondary" style={styles.selectedExerciseDetails}>
                  {selectedExercise.muscle_group} • {selectedExercise.equipment} • {selectedExercise.difficulty}
                </Typography>
              </View>
              <TouchableOpacity
                style={styles.crossButton}
                onPress={handleBackToSearch}
              >
                <Typography variant="paragraph-medium" color="text-secondary" style={styles.crossButtonText}>
                  ✕
                </Typography>
              </TouchableOpacity>
            </View>

            {/* Set Configuration */}
            <View style={styles.setConfiguration}>
              <Typography variant="label-medium" color="text-primary" style={styles.configTitle}>
                Configure Sets
              </Typography>
              
              {/* Number of Sets */}
              <View style={styles.numberOfSetsContainer}>
                <Typography variant="paragraph-small" color="text-primary" style={styles.inputLabel}>
                  Number of Sets *
                </Typography>
                <TextInput
                  value={numberOfSets}
                  onChangeText={setNumberOfSets}
                  keyboardType="numeric"
                  style={styles.numberOfSetsInput}
                  placeholderTextColor={getColor('text-secondary')}
                />
              </View>

              {/* Individual Sets */}
              <View style={styles.setsContainer}>
                {sets.map(renderSetConfiguration)}
              </View>

              {!isFormValid() && selectedExercise && (
                <View style={styles.validationError}>
                  <Typography variant="paragraph-small" color="danger">
                    All fields are required. Please ensure all sets have reps greater than 0 and weight greater than 0.
                  </Typography>
                </View>
              )}
            </View>

            {/* Add Button */}
            <Button
              variant="primary"
              size="large"
              style={styles.addButton}
              onPress={handleAddExercise}
              disabled={!isFormValid()}
            >
              Add to Routine
            </Button>
          </View>
        )}
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
  searchInput: {
    borderWidth: 1,
    borderColor: getColor('border-light'),
    borderRadius: 3,
    paddingHorizontal: 12,
    paddingVertical: 8,
    fontSize: 16,
    color: getColor('text-primary'),
    backgroundColor: getColor('surface'),
  },
  loadingContainer: {
    alignItems: 'center',
    paddingVertical: 32,
  },
  loadingText: {
    marginTop: 16,
  },
  emptyState: {
    alignItems: 'center',
    paddingVertical: 32,
  },
  emptyStateTitle: {
    textAlign: 'center',
  },
  emptyStateSubtitle: {
    textAlign: 'center',
    marginTop: 8,
  },
  exerciseList: {
  },
  exerciseItem: {
    padding: 16,
    borderWidth: 1,
    borderColor: getColor('border-light'),
    borderRadius: 8,
    marginBottom: 12,
    backgroundColor: getColor('surface'),
  },
  exerciseItemSelected: {
    borderColor: getColor('primary'),
    backgroundColor: getColor('surface'),
  },
  exerciseName: {
    fontWeight: '500',
    marginBottom: 4,
  },
  exerciseDetails: {
    marginBottom: 8,
  },
  exerciseDetailsSection: {
    marginBottom: 24,
  },
  exerciseMetadata: {
    marginBottom: 12,
    gap: 4,
  },
  metadataLabel: {
    fontWeight: '500',
  },
  setConfiguration: {
  },
  configTitle: {
    fontWeight: '500',
    marginBottom: 16,
  },
  numberOfSetsContainer: {
    marginBottom: 16,
  },
  inputLabel: {
    marginBottom: 8,
  },
  numberOfSetsInput: {
    borderWidth: 1,
    borderColor: getColor('border-light'),
    borderRadius: 3,
    paddingHorizontal: 12,
    paddingVertical: 8,
    fontSize: 16,
    color: getColor('text-primary'),
    backgroundColor: getColor('surface'),
  },
  setsContainer: {
  },
  setCard: {
    padding: 12,
    borderWidth: 1,
    borderColor: getColor('border-light'),
    borderRadius: 8,
    marginBottom: 12,
    backgroundColor: getColor('background'),
  },
  setTitle: {
    fontWeight: '500',
    marginBottom: 12,
  },
  setInputs: {
    flexDirection: 'row',
    gap: 12,
  },
  setInputContainer: {
    flex: 1,
  },
  setInput: {
    borderWidth: 1,
    borderColor: getColor('border-light'),
    borderRadius: 3,
    paddingHorizontal: 8,
    paddingVertical: 6,
    fontSize: 14,
    color: getColor('text-primary'),
    backgroundColor: getColor('surface'),
  },
  validationError: {
    padding: 12,
    backgroundColor: getColor('surface'),
    borderWidth: 1,
    borderColor: getColor('danger'),
    borderRadius: 8,
    marginTop: 16,
  },
  addButton: {
    marginTop: 24,
  },
  selectedExerciseHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
  },
  selectedExerciseInfo: {
    flex: 1,
  },
  selectedExerciseName: {
    fontWeight: '500',
    marginBottom: 4,
  },
  selectedExerciseDetails: {
    marginBottom: 8,
  },
  crossButton: {
    padding: 8,
  },
  crossButtonText: {
    fontSize: 16,
    fontWeight: 'bold',
  },
}); 