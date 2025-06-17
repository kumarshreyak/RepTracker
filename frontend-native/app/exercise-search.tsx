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
      <Typography variant="text-default" color="dark" style={styles.exerciseName}>
        {exercise.name}
      </Typography>
      <Typography variant="text-small" color="light" style={styles.exerciseDetails}>
        {exercise.muscle_group} • {exercise.equipment} • {exercise.difficulty}
      </Typography>
      <Typography variant="text-small" color="light" style={styles.exerciseDescription}>
        {exercise.description}
      </Typography>
    </TouchableOpacity>
  );

  const renderSetConfiguration = (set: WorkoutSetInput, index: number) => (
    <View key={index} style={styles.setCard}>
      <Typography variant="text-small" color="dark" style={styles.setTitle}>
        Set {index + 1}
      </Typography>
      <View style={styles.setInputs}>
        <View style={styles.setInputContainer}>
          <Typography variant="text-small" color="dark" style={styles.inputLabel}>
            Reps *
          </Typography>
          <TextInput
            value={set.reps}
            onChangeText={(text) => updateSet(index, 'reps', text)}
            keyboardType="numeric"
            style={styles.setInput}
            placeholderTextColor={getColor('light')}
          />
        </View>
        <View style={styles.setInputContainer}>
          <Typography variant="text-small" color="dark" style={styles.inputLabel}>
            Weight (lbs) *
          </Typography>
          <TextInput
            value={set.weight}
            onChangeText={(text) => updateSet(index, 'weight', text)}
            keyboardType="numeric"
            style={styles.setInput}
            placeholderTextColor={getColor('light')}
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
          <Typography variant="heading-xxlarge" color="dark" style={styles.title}>
            Add Exercise
          </Typography>
          <Typography variant="text-default" color="light">
            Search and select exercises for your routine
          </Typography>
        </View>

        {/* Search Section */}
        <View style={styles.card}>
          <Typography variant="heading-small" color="dark" style={styles.sectionTitle}>
            Search Exercises
          </Typography>
          <TextInput
            value={searchQuery}
            onChangeText={setSearchQuery}
            placeholder="Search by name, muscle group, or equipment..."
            style={styles.searchInput}
            placeholderTextColor={getColor('light')}
          />
        </View>

        {/* Exercise List */}
        <View style={styles.card}>
          <Typography variant="heading-small" color="dark" style={styles.sectionTitle}>
            Exercise Library
          </Typography>

          {loading ? (
            <View style={styles.loadingContainer}>
              <ActivityIndicator size="large" color={getColor('blue-bright')} />
              <Typography variant="text-default" color="light" style={styles.loadingText}>
                Loading exercises...
              </Typography>
            </View>
          ) : exercises.length === 0 ? (
            <View style={styles.emptyState}>
              <Typography variant="text-default" color="light" style={styles.emptyStateTitle}>
                No exercises found
              </Typography>
              {searchQuery && (
                <Typography variant="text-small" color="light" style={styles.emptyStateSubtitle}>
                  Try a different search term
                </Typography>
              )}
            </View>
          ) : (
            <ScrollView style={styles.exerciseList} nestedScrollEnabled>
              {exercises.map(renderExerciseItem)}
            </ScrollView>
          )}
        </View>

        {/* Exercise Details & Configuration */}
        {selectedExercise && (
          <View style={styles.card}>
            <Typography variant="heading-small" color="dark" style={styles.sectionTitle}>
              Exercise Details
            </Typography>

            {/* Exercise Info */}
            <View style={styles.exerciseDetailsSection}>
              <Typography variant="text-default" color="dark" style={styles.selectedExerciseName}>
                {selectedExercise.name}
              </Typography>
              <View style={styles.exerciseMetadata}>
                <Typography variant="text-small" color="light">
                  <Typography variant="text-small" color="dark" style={styles.metadataLabel}>Muscle:</Typography> {selectedExercise.muscle_group}
                </Typography>
                <Typography variant="text-small" color="light">
                  <Typography variant="text-small" color="dark" style={styles.metadataLabel}>Equipment:</Typography> {selectedExercise.equipment}
                </Typography>
                <Typography variant="text-small" color="light">
                  <Typography variant="text-small" color="dark" style={styles.metadataLabel}>Difficulty:</Typography> {selectedExercise.difficulty}
                </Typography>
              </View>
              <Typography variant="text-small" color="light" style={styles.selectedExerciseDescription}>
                {selectedExercise.description}
              </Typography>
            </View>

            {/* Set Configuration */}
            <View style={styles.setConfiguration}>
              <Typography variant="text-default" color="dark" style={styles.configTitle}>
                Configure Sets
              </Typography>
              
              {/* Number of Sets */}
              <View style={styles.numberOfSetsContainer}>
                <Typography variant="text-small" color="dark" style={styles.inputLabel}>
                  Number of Sets *
                </Typography>
                <TextInput
                  value={numberOfSets}
                  onChangeText={setNumberOfSets}
                  keyboardType="numeric"
                  style={styles.numberOfSetsInput}
                  placeholderTextColor={getColor('light')}
                />
              </View>

              {/* Individual Sets */}
              <ScrollView style={styles.setsContainer} nestedScrollEnabled>
                {sets.map(renderSetConfiguration)}
              </ScrollView>

              {!isFormValid() && selectedExercise && (
                <View style={styles.validationError}>
                  <Typography variant="text-small" color="red">
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
  searchInput: {
    borderWidth: 1,
    borderColor: getColor('light-gray-3'),
    borderRadius: 3,
    paddingHorizontal: 12,
    paddingVertical: 8,
    fontSize: 16,
    color: getColor('dark'),
    backgroundColor: 'white',
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
    maxHeight: 300,
  },
  exerciseItem: {
    padding: 16,
    borderWidth: 1,
    borderColor: getColor('light-gray-3'),
    borderRadius: 8,
    marginBottom: 12,
    backgroundColor: 'white',
  },
  exerciseItemSelected: {
    borderColor: getColor('blue-bright'),
    backgroundColor: getColor('blue-light-2'),
  },
  exerciseName: {
    fontWeight: '500',
    marginBottom: 4,
  },
  exerciseDetails: {
    marginBottom: 8,
  },
  exerciseDescription: {
    lineHeight: 16,
  },
  exerciseDetailsSection: {
    marginBottom: 24,
  },
  selectedExerciseName: {
    fontWeight: '500',
    marginBottom: 12,
  },
  exerciseMetadata: {
    marginBottom: 12,
    gap: 4,
  },
  metadataLabel: {
    fontWeight: '500',
  },
  selectedExerciseDescription: {
    lineHeight: 18,
  },
  setConfiguration: {
    marginTop: 24,
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
    borderColor: getColor('light-gray-3'),
    borderRadius: 3,
    paddingHorizontal: 12,
    paddingVertical: 8,
    fontSize: 16,
    color: getColor('dark'),
    backgroundColor: 'white',
  },
  setsContainer: {
    maxHeight: 300,
  },
  setCard: {
    padding: 12,
    borderWidth: 1,
    borderColor: getColor('light-gray-3'),
    borderRadius: 8,
    marginBottom: 12,
    backgroundColor: getColor('light-gray-1'),
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
    borderColor: getColor('light-gray-3'),
    borderRadius: 3,
    paddingHorizontal: 8,
    paddingVertical: 6,
    fontSize: 14,
    color: getColor('dark'),
    backgroundColor: 'white',
  },
  validationError: {
    padding: 12,
    backgroundColor: getColor('red-light-2'),
    borderWidth: 1,
    borderColor: getColor('red-light-1'),
    borderRadius: 8,
    marginTop: 16,
  },
  addButton: {
    marginTop: 24,
  },
}); 