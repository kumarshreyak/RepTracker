// Exercise types matching the backend API schema

export interface Exercise {
  id: string;
  name: string;
  description: string;
  category: string;                     // Changed from difficulty
  equipment: string[];                  // Changed to array
  primaryMuscles: string[];            // Changed from muscleGroup
  secondaryMuscles: string[];          // New field
  instructions: string[];
  video?: string;                      // New field
  variationsOn?: string[];             // New field
  variationOn?: string[];              // New field
  createdAt: string;
  updatedAt: string;
}

export interface RoutineSet {
  reps: number;
  weight: number;
  durationSeconds?: number;
  distance?: number;
  notes?: string;
}

export interface RoutineExercise {
  id: string;
  name: string;
  primaryMuscles: string[];            // Updated from muscleGroup
  secondaryMuscles?: string[];         // New optional field
  sets: RoutineSet[];
}

export interface WorkoutSet {
  targetReps: number;
  targetWeight: number;
  actualReps: number;
  actualWeight: number;
  durationSeconds?: number;
  distance?: number;
  notes?: string;
  completed?: boolean;
}

export interface RoutineSetInput {
  reps: string;
  weight: string;
}

export interface WorkoutExercise {
  exerciseId: string;
  exercise?: Exercise;                 // Now uses updated Exercise interface
  sets: WorkoutSet[];
  notes?: string;
  restSeconds?: number;
  completed?: boolean;
}

export interface ActiveWorkout {
  id?: string;
  sessionId?: string;
  routineId: string;
  routineName: string;
  routineDescription?: string;
  exercises: WorkoutExercise[];
  startedAt: Date;
  notes?: string;
}

// Exercise category constants (matching backend)
export const ExerciseCategories = {
  STRENGTH: "strength",
  STRETCHING: "stretching",
  PLYOMETRICS: "plyometrics",
  STRONGMAN: "strongman",
  CARDIO: "cardio",
  OLYMPIC_WEIGHTLIFTING: "olympic weightlifting",
  CROSSFIT: "crossfit",
  CALISTHENICS: "calisthenics",
} as const;

export type ExerciseCategory = typeof ExerciseCategories[keyof typeof ExerciseCategories];

// Equipment type constants (matching backend)
export const EquipmentTypes = {
  NONE: "none",
  BARBELL: "barbell",
  DUMBBELL: "dumbbell",
  KETTLEBELL: "kettlebell",
  MACHINE: "machine",
  CABLE: "cable",
  BODYWEIGHT: "bodyweight",
  BANDS: "bands",
  MEDICINE_BALL: "medicine ball",
  EZ_BAR: "ez-bar",
  FOAM_ROLL: "foam roll",
  PLATE: "plate",
  RESISTANCE_BAND: "resistance band",
  TRAP_BAR: "trap bar",
  SMITH_MACHINE: "smith machine",
  PULL_UP_BAR: "pull-up bar",
} as const;

export type EquipmentType = typeof EquipmentTypes[keyof typeof EquipmentTypes];

// Exercise History Types (for history component)
export interface HistorySet {
  targetReps: number;
  targetWeight: number;
  actualReps: number;
  actualWeight: number;
  completed: boolean;
  startedAt: string;
  finishedAt?: string;
}

export interface HistorySessionExercise {
  exerciseId: string;
  exercise: Exercise;
  sets: HistorySet[];
  notes?: string;
  restSeconds?: number;
  completed: boolean;
  startedAt: string;
  finishedAt?: string;
}

export interface HistorySessionInfo {
  id: string;
  name: string;
  finishedAt: Date | string;
}

export interface HistoryAIExercise {
  exerciseId: string;
  exerciseName: string;
  progressionApplied: boolean;
  reasoning: string;
  changesMade: string;
  sets: Array<{
    reps: number;
    weight: number;
  }>;
  restSeconds?: number;
}

export interface ExerciseHistoryEntry {
  sessionExercise: HistorySessionExercise;
  sessionInfo: HistorySessionInfo;
  aiExercise?: HistoryAIExercise;
}

export interface ExerciseHistoryResponse {
  history: ExerciseHistoryEntry[];
  exerciseName: string;
} 