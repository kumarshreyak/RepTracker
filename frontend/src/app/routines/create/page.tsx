"use client";

import { Typography, Button } from "@/components";
import Link from "next/link";
import { useState, useEffect } from "react";
import { useAuth } from "@/hooks/useAuth";

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

export default function CreateRoutinePage() {
  const { backendUser, isAuthenticated, isLoading } = useAuth();
  const [routineName, setRoutineName] = useState("");
  const [exercises, setExercises] = useState<RoutineExercise[]>([]);
  const [loading, setLoading] = useState(false);

  // Load exercises from localStorage on mount
  useEffect(() => {
    const loadExercises = () => {
      const saved = localStorage.getItem('routineExercises');
      if (saved) {
        try {
          const parsedExercises = JSON.parse(saved);
          setExercises(parsedExercises);
        } catch (error) {
          console.error('Error loading exercises from localStorage:', error);
          setExercises([]);
        }
      }
    };

    loadExercises();

    // Listen for storage changes (when coming back from exercise search)
    const handleStorageChange = () => {
      loadExercises();
    };

    window.addEventListener('focus', handleStorageChange);
    return () => window.removeEventListener('focus', handleStorageChange);
  }, []);

  const handleSaveRoutine = async () => {
    if (!routineName.trim() || exercises.length === 0 || !backendUser?.id) {
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
        user_id: backendUser.id,
        name: routineName,
        description: `Workout routine with ${exercises.length} exercises`,
        exercises: workoutExercises,
        started_at: null,
        notes: "",
      };

      const response = await fetch('/api/routines', {
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

      // Clear localStorage after successful save
      localStorage.removeItem('routineExercises');
      
      // Navigate back to home
      window.location.href = '/home';
    } catch (error) {
      console.error('Error saving routine:', error);
      alert('Failed to save routine. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const removeExercise = (index: number) => {
    const newExercises = exercises.filter((_, i) => i !== index);
    setExercises(newExercises);
    localStorage.setItem('routineExercises', JSON.stringify(newExercises));
  };

  const clearAllExercises = () => {
    setExercises([]);
    localStorage.removeItem('routineExercises');
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-light-gray-1 flex items-center justify-center">
        <Typography variant="text-large" color="light">
          Loading...
        </Typography>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-light-gray-1 flex items-center justify-center">
        <Typography variant="text-large" color="light">
          Please sign in to continue
        </Typography>
      </div>
    );
  }

  return (
    <main className="min-h-screen bg-light-gray-1">
      <div className="max-w-4xl mx-auto px-6 py-8">
        {/* Page Header */}
        <div className="mb-8">
          <Link href="/home" className="inline-block mb-4">
            <Button variant="text" size="default">
              ← Back to Home
            </Button>
          </Link>
          <Typography variant="heading-xxlarge" color="dark" className="mb-2">
            Create Routine
          </Typography>
          <Typography variant="text-default" color="light">
            Build your workout routine by adding exercises
          </Typography>
        </div>

        {/* Routine Details */}
        <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6 mb-6">
          <div className="mb-6">
            <Typography variant="heading-small" color="dark" className="mb-3">
              Routine Details
            </Typography>
            <div>
              <label className="block mb-2">
                <Typography variant="text-small" color="dark">
                  Routine Name *
                </Typography>
              </label>
              <input
                type="text"
                value={routineName}
                onChange={(e) => setRoutineName(e.target.value)}
                placeholder="Enter routine name"
                className="w-full px-3 py-2 border border-light-gray-3 rounded-[3px] focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>
        </div>

        {/* Exercises Section */}
        <div className="bg-white rounded-lg shadow-sm border border-light-gray-3 p-6 mb-6">
          <div className="flex justify-between items-center mb-6">
            <Typography variant="heading-small" color="dark">
              Exercises ({exercises.length})
            </Typography>
            <div className="flex gap-3">
              {exercises.length > 0 && (
                <Button variant="text" size="small" onClick={clearAllExercises}>
                  Clear All
                </Button>
              )}
              <Link href="/exercises/search">
                <Button variant="primary" size="default">
                  Add Exercise
                </Button>
              </Link>
            </div>
          </div>

          {exercises.length === 0 ? (
            <div className="text-center py-12">
              <Typography variant="text-default" color="light" className="mb-4">
                No exercises added yet
              </Typography>
              <Typography variant="text-small" color="light">
                Click "Add Exercise" to start building your routine
              </Typography>
            </div>
          ) : (
            <div className="space-y-4">
              {exercises.map((exercise, index) => (
                <div
                  key={`${exercise.id}-${index}`}
                  className="p-4 border border-light-gray-3 rounded-lg"
                >
                  <div className="flex justify-between items-start mb-3">
                    <div className="flex-1">
                      <Typography variant="text-default" color="dark" className="font-medium mb-1">
                        {exercise.name}
                      </Typography>
                      <Typography variant="text-small" color="light">
                        {exercise.muscle_group} • {exercise.sets.length} sets
                      </Typography>
                    </div>
                    <Button
                      variant="text"
                      size="small"
                      onClick={() => removeExercise(index)}
                    >
                      Remove
                    </Button>
                  </div>
                  
                  {/* Set Details */}
                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
                    {exercise.sets.map((set, setIndex) => (
                      <div
                        key={setIndex}
                        className="px-3 py-2 bg-light-gray-1 rounded-lg"
                      >
                        <Typography variant="text-small" color="dark" className="font-medium">
                          Set {setIndex + 1}
                        </Typography>
                        <Typography variant="text-small" color="light">
                          {set.reps} reps @ {set.weight}lbs
                        </Typography>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Action Buttons */}
        <div className="flex gap-4">
          <Link href="/home" className="flex-1">
            <Button variant="secondary" size="large" className="w-full justify-center">
              Cancel
            </Button>
          </Link>
          <Button
            variant="primary"
            size="large"
            className="flex-1 justify-center"
            onClick={handleSaveRoutine}
            disabled={!routineName.trim() || exercises.length === 0 || loading || !backendUser?.id}
          >
            {loading ? "Saving..." : "Save Routine"}
          </Button>
        </div>
      </div>
    </main>
  );
} 