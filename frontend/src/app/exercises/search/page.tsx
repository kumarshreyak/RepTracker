"use client";

import { Typography, Button } from "@/components";
import Link from "next/link";
import { useState, useEffect } from "react";

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

interface RoutineExercise {
  id: string;
  name: string;
  muscle_group: string;
  sets: WorkoutSet[];
}

export default function ExerciseSearchPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedExercise, setSelectedExercise] = useState<Exercise | null>(null);
  const [numberOfSets, setNumberOfSets] = useState(3);
  const [sets, setSets] = useState<WorkoutSet[]>([
    { reps: 10, weight: 0 },
    { reps: 10, weight: 0 },
    { reps: 10, weight: 0 },
  ]);

  // Fetch exercises
  const fetchExercises = async (query?: string) => {
    setLoading(true);
    try {
      const params = new URLSearchParams();
      if (query) params.append("search", query);
      
      const response = await fetch(`/api/exercises?${params.toString()}`);
      if (!response.ok) throw new Error("Failed to fetch exercises");
      
      const data = await response.json();
      setExercises(data.exercises || []);
    } catch (error) {
      console.error("Error fetching exercises:", error);
      setExercises([]);
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
    const newSets = Array.from({ length: numberOfSets }, (_, index) => 
      sets[index] || { reps: 10, weight: 0 }
    );
    setSets(newSets);
  }, [numberOfSets]);

  const updateSet = (index: number, field: 'reps' | 'weight', value: number) => {
    const newSets = [...sets];
    newSets[index] = { ...newSets[index], [field]: value };
    setSets(newSets);
  };

  const isFormValid = () => {
    return selectedExercise && 
           sets.length > 0 && 
           sets.every(set => set.reps > 0 && set.weight > 0);
  };

  const handleAddExercise = () => {
    if (!selectedExercise || !isFormValid()) return;
    
    const exerciseData: RoutineExercise = {
      id: selectedExercise.id,
      name: selectedExercise.name,
      muscle_group: selectedExercise.muscle_group,
      sets: sets,
    };
    
    // Save to localStorage
    const existingExercises = JSON.parse(localStorage.getItem('routineExercises') || '[]');
    existingExercises.push(exerciseData);
    localStorage.setItem('routineExercises', JSON.stringify(existingExercises));
    
    console.log("Adding exercise to routine:", exerciseData);
    
    // Navigate back to create routine page
    window.history.back();
  };

  return (
    <main className="min-h-screen bg-light-gray-1">
      <div className="max-w-4xl mx-auto px-6 py-8">
        {/* Page Header */}
        <div className="mb-8">
          <Link href="/routines/create" className="inline-block mb-4">
            <Button variant="text" size="default">
              ← Back to Create Routine
            </Button>
          </Link>
          <Typography variant="heading-xxlarge" color="dark" className="mb-2">
            Add Exercise
          </Typography>
          <Typography variant="text-default" color="light">
            Search and select exercises for your routine
          </Typography>
        </div>

        {/* Search Section */}
        <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6 mb-6">
          <div className="mb-4">
            <Typography variant="heading-small" color="dark" className="mb-3">
              Search Exercises
            </Typography>
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search by name, muscle group, or equipment..."
              className="w-full px-4 py-3 border border-light-gray-3 rounded-[3px] focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Exercise List */}
          <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6">
            <Typography variant="heading-small" color="dark" className="mb-4">
              Exercise Library
            </Typography>

            {loading ? (
              <div className="text-center py-8">
                <Typography variant="text-default" color="light">
                  Loading exercises...
                </Typography>
              </div>
            ) : exercises.length === 0 ? (
              <div className="text-center py-8">
                <Typography variant="text-default" color="light">
                  No exercises found
                </Typography>
                {searchQuery && (
                  <Typography variant="text-small" color="light" className="mt-2">
                    Try a different search term
                  </Typography>
                )}
              </div>
            ) : (
              <div className="space-y-3 max-h-96 overflow-y-auto">
                {exercises.map((exercise) => (
                  <div
                    key={exercise.id}
                    className={`p-4 border rounded-lg cursor-pointer transition-colors ${
                      selectedExercise?.id === exercise.id
                        ? "border-blue-500 bg-blue-50"
                        : "border-light-gray-3 hover:border-light-gray-4"
                    }`}
                    onClick={() => setSelectedExercise(exercise)}
                  >
                    <Typography variant="text-default" color="dark" className="font-medium mb-1">
                      {exercise.name}
                    </Typography>
                    <Typography variant="text-small" color="light" className="mb-2">
                      {exercise.muscle_group} • {exercise.equipment} • {exercise.difficulty}
                    </Typography>
                    <Typography variant="text-small" color="light" className="line-clamp-2">
                      {exercise.description}
                    </Typography>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Exercise Details & Configuration */}
          <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6">
            <Typography variant="heading-small" color="dark" className="mb-4">
              Exercise Details
            </Typography>

            {selectedExercise ? (
              <div className="space-y-6">
                {/* Exercise Info */}
                <div>
                  <Typography variant="text-default" color="dark" className="font-medium mb-2">
                    {selectedExercise.name}
                  </Typography>
                  <div className="space-y-2 mb-4">
                    <div className="flex gap-4">
                      <Typography variant="text-small" color="light">
                        <span className="font-medium">Muscle:</span> {selectedExercise.muscle_group}
                      </Typography>
                      <Typography variant="text-small" color="light">
                        <span className="font-medium">Equipment:</span> {selectedExercise.equipment}
                      </Typography>
                    </div>
                    <Typography variant="text-small" color="light">
                      <span className="font-medium">Difficulty:</span> {selectedExercise.difficulty}
                    </Typography>
                  </div>
                  <Typography variant="text-small" color="light">
                    {selectedExercise.description}
                  </Typography>
                </div>

                {/* Set Configuration */}
                <div className="space-y-4">
                  <Typography variant="text-default" color="dark" className="font-medium">
                    Configure Sets
                  </Typography>
                  
                  {/* Number of Sets */}
                  <div>
                    <label className="block mb-2">
                      <Typography variant="text-small" color="dark">Number of Sets *</Typography>
                    </label>
                    <input
                      type="number"
                      value={numberOfSets}
                      onChange={(e) => setNumberOfSets(Math.max(1, parseInt(e.target.value) || 1))}
                      min="1"
                      max="10"
                      className="w-full px-3 py-2 border border-light-gray-3 rounded-[3px] focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>

                  {/* Individual Sets */}
                  <div className="space-y-3 max-h-64 overflow-y-auto">
                    {sets.map((set, index) => (
                      <div key={index} className="p-3 border border-light-gray-3 rounded-lg">
                        <Typography variant="text-small" color="dark" className="font-medium mb-2">
                          Set {index + 1}
                        </Typography>
                        <div className="grid grid-cols-2 gap-3">
                          <div>
                            <label className="block mb-1">
                              <Typography variant="text-small" color="dark">Reps *</Typography>
                            </label>
                            <input
                              type="number"
                              value={set.reps}
                              onChange={(e) => updateSet(index, 'reps', parseInt(e.target.value) || 0)}
                              min="1"
                              max="100"
                              className="w-full px-2 py-1 border border-light-gray-3 rounded-[3px] focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                          </div>
                          <div>
                            <label className="block mb-1">
                              <Typography variant="text-small" color="dark">Weight (lbs) *</Typography>
                            </label>
                            <input
                              type="number"
                              value={set.weight}
                              onChange={(e) => updateSet(index, 'weight', parseFloat(e.target.value) || 0)}
                              min="0"
                              step="0.5"
                              className="w-full px-2 py-1 border border-light-gray-3 rounded-[3px] focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>

                  {!isFormValid() && (
                    <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
                      <Typography variant="text-small" color="red">
                        All fields are required. Please ensure all sets have reps greater than 0 and weight greater than 0.
                      </Typography>
                    </div>
                  )}
                </div>

                {/* Add Button */}
                <Button
                  variant="primary"
                  size="large"
                  className="w-full justify-center"
                  onClick={handleAddExercise}
                  disabled={!isFormValid()}
                >
                  Add to Routine
                </Button>
              </div>
            ) : (
              <div className="text-center py-12">
                <Typography variant="text-default" color="light">
                  Select an exercise to view details
                </Typography>
              </div>
            )}
          </div>
        </div>
      </div>
    </main>
  );
} 