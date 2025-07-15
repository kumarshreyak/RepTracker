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
  Dimensions,
  Vibration,
} from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { router } from 'expo-router';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { Typography, Button } from '../src/components';
import { getColor } from '../src/components/Colors';
import { useAuth } from '../src/hooks/useAuth';
import { User } from '../src/auth/AuthService';
import { Exercise, RoutineExercise, RoutineSet, RoutineSetInput } from '@/types/exercise';

export default function ExerciseSearchRoute() {
  const { user } = useAuth();
  const [searchQuery, setSearchQuery] = useState("");
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedExercise, setSelectedExercise] = useState<Exercise | null>(null);
  const [numberOfSets, setNumberOfSets] = useState("3");
  const [sets, setSets] = useState<RoutineSetInput[]>([
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

  const copySetToAll = (sourceIndex: number) => {
    const sourceSet = sets[sourceIndex];
    const newSets = sets.map(() => ({ ...sourceSet }));
    setSets(newSets);
  };

  const isFormValid = () => {
    return selectedExercise && 
           sets.length > 0 && 
           sets.every(set => {
             const reps = parseFloat(set.reps) || 0;
             return reps > 0;
           });
  };

  const handleAddExercise = async () => {
    if (!selectedExercise || !isFormValid()) return;
    
    // Convert string inputs to numbers for storage
    const convertedSets: RoutineSet[] = sets.map(set => ({
      reps: parseFloat(set.reps) || 0,
      weight: parseFloat(set.weight) || 0,
    }));
    
    const exerciseData: RoutineExercise = {
      id: selectedExercise.id,
      name: selectedExercise.name,
      primaryMuscles: selectedExercise.primaryMuscles,
      secondaryMuscles: selectedExercise.secondaryMuscles,
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

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity onPress={() => router.back()} style={styles.backButton}>
          <Typography variant="label-medium" color="contentSecondary">←</Typography>
        </TouchableOpacity>
        <Typography variant="label-medium" color="contentPrimary">
          Add Exercise
        </Typography>
        <View style={styles.headerSpacer} />
      </View>

      <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false}>
        {/* Search Bar - Only visible when no exercise is selected */}
        {!selectedExercise && (
          <View style={styles.searchContainer}>
            <TextInput
              value={searchQuery}
              onChangeText={setSearchQuery}
              placeholder="Search exercises..."
              style={styles.searchInput}
              placeholderTextColor={getColor('contentTertiary')}
            />
          </View>
        )}

        {/* Selected Exercise Chip */}
        {selectedExercise && (
          <TouchableOpacity style={styles.selectedChip} onPress={handleBackToSearch}>
            <View style={styles.chipContent}>
              <View style={styles.chipText}>
                <Typography variant="label-medium" color="contentPrimary">
                  {selectedExercise.name}
                </Typography>
                <Typography variant="paragraph-xsmall" color="contentSecondary">
                  {selectedExercise.primaryMuscles?.join(", ") || "No muscle groups"}
                </Typography>
              </View>
              <MaterialIcons 
                name="close" 
                size={20} 
                color={getColor('contentSecondary')} 
              />
            </View>
          </TouchableOpacity>
        )}

        {/* Exercise List - Compact grid view when searching */}
        {!selectedExercise && !loading && exercises.length > 0 && (
          <View style={styles.exerciseGrid}>
            {exercises.map((exercise) => (
              <TouchableOpacity
                key={exercise.id}
                style={styles.exerciseGridItem}
                onPress={() => setSelectedExercise(exercise)}
              >
                <Typography 
                  variant="label-small" 
                  color="contentPrimary"
                  style={styles.exerciseName}
                >
                  {exercise.name}
                </Typography>
                <Typography 
                  variant="paragraph-xsmall" 
                  color="contentTertiary"
                >
                  {exercise.primaryMuscles?.join(", ") || "No muscle groups"}
                </Typography>
              </TouchableOpacity>
            ))}
          </View>
        )}

        {/* Loading State */}
        {loading && (
          <View style={styles.centerContainer}>
            <ActivityIndicator size="large" color={getColor('accent')} />
          </View>
        )}

        {/* Empty State */}
        {!selectedExercise && !loading && exercises.length === 0 && searchQuery && (
          <View style={styles.centerContainer}>
            <Typography variant="paragraph-medium" color="contentTertiary">
              No results
            </Typography>
          </View>
        )}

        {/* Sets Configuration - Streamlined */}
        {selectedExercise && (
          <View style={styles.configSection}>
            {/* Quick Set Selector */}
            <View style={styles.setSelector}>
              <Typography variant="label-small" color="contentPrimary">
                Sets
              </Typography>
              <View style={styles.setButtons}>
                {[1, 2, 3, 4, 5].map((num) => (
                  <TouchableOpacity
                    key={num}
                    style={[
                      styles.setButton,
                      parseInt(numberOfSets) === num && styles.setButtonActive
                    ]}
                    onPress={() => setNumberOfSets(num.toString())}
                  >
                    <Typography 
                      variant="label-small" 
                      color={parseInt(numberOfSets) === num ? 'contentOnColor' : 'contentPrimary'}
                    >
                      {num}
                    </Typography>
                  </TouchableOpacity>
                ))}
                <TextInput
                  value={numberOfSets}
                  onChangeText={setNumberOfSets}
                  keyboardType="numeric"
                  style={styles.customSetInput}
                  placeholder="..."
                  placeholderTextColor={getColor('contentTertiary')}
                />
              </View>
            </View>

            {/* Compact Sets Input */}
            <View style={styles.setsInputContainer}>
              {sets.map((set, index) => (
                <View key={index} style={styles.setRow}>
                  <Typography 
                    variant="label-xsmall" 
                    color="contentSecondary"
                    style={styles.setNumber}
                  >
                    {index + 1}
                  </Typography>
                  
                  <View style={styles.setInputGroup}>
                    <TextInput
                      value={set.reps}
                      onChangeText={(text) => updateSet(index, 'reps', text)}
                      keyboardType="numeric"
                      placeholder="Reps"
                      style={[
                        styles.compactInput,
                        !set.reps && styles.inputError
                      ]}
                      placeholderTextColor={getColor('contentTertiary')}
                    />
                    
                    <Typography variant="label-xsmall" color="contentTertiary">×</Typography>
                    
                    <TextInput
                      value={set.weight}
                      onChangeText={(text) => updateSet(index, 'weight', text)}
                      keyboardType="numeric"
                      placeholder="Weight"
                      style={styles.compactInput}
                      placeholderTextColor={getColor('contentTertiary')}
                    />
                    
                    <Typography variant="label-xsmall" color="contentTertiary">kg</Typography>
                  </View>
                  
                  {/* Action area - always present for uniform spacing */}
                  <View style={styles.actionArea}>
                    {index === 0 && sets.length > 1 ? (
                      <TouchableOpacity
                        style={styles.copyButton}
                        onPress={() => copySetToAll(index)}
                      >
                        <MaterialIcons 
                          name="content-copy" 
                          size={20} 
                          color={getColor('contentSecondary')} 
                        />
                      </TouchableOpacity>
                    ) : null}
                  </View>
                </View>
              ))}
            </View>

            {/* Action Button */}
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

const { width } = Dimensions.get('window');
const gridItemWidth = (width - 48) / 2; // 2 columns with padding

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: getColor('backgroundSecondary'),
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
    padding: 8,
  },
  headerSpacer: {
    width: 40, // Balance the header
  },
  scrollView: {
    flex: 1,
  },
  searchContainer: {
    margin: 16,
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderTransparent'),
  },
  searchInput: {
    flex: 1,
    paddingHorizontal: 16,
    paddingVertical: 12,
    fontSize: 16,
    color: getColor('contentPrimary'),
  },

  selectedChip: {
    marginHorizontal: 16,
    padding: 12,
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderSelected'),
  },
  chipContent: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    width: '100%',
  },
  chipText: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  exerciseGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    paddingHorizontal: 16,
    gap: 8,
  },
  exerciseGridItem: {
    width: gridItemWidth,
    padding: 12,
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    borderWidth: 1,
    borderColor: getColor('borderTransparent'),
  },
  exerciseName: {
    marginBottom: 4,
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: 48,
  },
  configSection: {
    padding: 16,
  },
  setSelector: {
    marginBottom: 24,
  },
  setButtons: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginTop: 8,
  },
  setButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: getColor('backgroundPrimary'),
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    justifyContent: 'center',
    alignItems: 'center',
  },
  setButtonActive: {
    backgroundColor: getColor('backgroundAccent'),
    borderColor: getColor('backgroundAccent'),
  },
  customSetInput: {
    width: 50,
    height: 40,
    borderRadius: 8,
    backgroundColor: getColor('backgroundPrimary'),
    borderWidth: 1,
    borderColor: getColor('borderOpaque'),
    paddingHorizontal: 8,
    fontSize: 14,
    color: getColor('contentPrimary'),
    textAlign: 'center',
  },
  setsInputContainer: {
    gap: 12,
  },
  setRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  setNumber: {
    width: 20,
    textAlign: 'center',
  },
  setInputGroup: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    backgroundColor: getColor('backgroundPrimary'),
    borderRadius: 8,
    paddingHorizontal: 12,
    paddingVertical: 8,
  },
  compactInput: {
    flex: 1,
    fontSize: 16,
    color: getColor('contentPrimary'),
    textAlign: 'center',
    paddingVertical: 4,
  },
  inputError: {
    color: getColor('contentNegative'),
  },
  addButton: {
    marginTop: 32,
  },
  actionArea: {
    width: 44,
    alignItems: 'center',
    justifyContent: 'center',
  },
  copyButton: {
    padding: 8,
  },
}); 