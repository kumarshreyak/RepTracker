// AI Progressive Overload Types
// Matches backend AIProgressiveOverloadResponseHTTP structure

export interface WorkoutSet {
  reps: number;
  weight: number;
  durationSeconds: number;
  distance: number;
  notes: string;
}

export interface AIProgressiveOverloadExercise {
  exerciseId: string;
  exerciseName: string;
  progressionApplied: boolean;
  reasoning: string;
  changesMade: string;
  sets: WorkoutSet[];
  notes: string;
  restSeconds: number;
}

export interface AIProgressiveOverloadResponse {
  id: string;
  userId: string;
  workoutSessionId: string;
  workoutId: string;
  analysisSummary: string;
  updatedExercises: AIProgressiveOverloadExercise[];
  success: boolean;
  message: string;
  recentSessionsCount: number;
  aiModel: string;
  processingTimeMs: number;
  createdAt: string;
  updatedAt: string;
}

export interface ListAIProgressiveOverloadResponsesResponse {
  responses: AIProgressiveOverloadResponse[];
  nextPageToken?: string;
} 