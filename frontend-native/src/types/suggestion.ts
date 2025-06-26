// Suggestion types matching the backend protobuf schema with camelCase naming

import { WorkoutExercise } from './exercise';

export interface WorkoutChange {
  type: string; // "exercise", "sets", "reps", "weight", "rest", "remove", "add"
  exerciseId?: string;
  exerciseName: string;
  oldValue: string;
  newValue: string;
  reason: string; // One-line explanation for the change
}

export interface SuggestedWorkout {
  originalWorkoutId: string;
  name: string;
  description: string;
  exercises: WorkoutExercise[];
  changes: WorkoutChange[];
  overallReasoning: string;
  priority: number; // 1-5, 5 being most recommended
  status?: string; // "pending", "accepted", "rejected"
  statusUpdatedAt?: string;
}

export interface StoredSuggestedWorkout {
  id: string;
  userId: string;
  originalWorkoutId: string;
  name: string;
  description: string;
  exercises: WorkoutExercise[];
  changes: WorkoutChange[];
  overallReasoning: string;
  priority: number; // 1-5, 5 being most recommended
  status?: string; // "pending", "accepted", "rejected"
  statusUpdatedAt?: string;
  analysisSummary: string;
  daysAnalyzed: number;
  createdAt: string;
  updatedAt: string;
}

export interface GenerateWorkoutSuggestionsRequest {
  daysToAnalyze?: number; // default 14 days
  maxSuggestions?: number; // default 3
}

export interface GenerateWorkoutSuggestionsResponse {
  suggestions: SuggestedWorkout[];
  analysisSummary: string;
}

export interface GetStoredSuggestionsResponse {
  suggestedWorkouts: StoredSuggestedWorkout[];
  nextPageToken?: string;
}

export interface ConfirmSuggestionRequest {
  accept: boolean; // true for accept, false for reject
}

export interface ConfirmSuggestionResponse {
  success: boolean;
  message: string;
}

// Insight types matching the backend protobuf schema
export interface WorkoutInsight {
  id: string;
  userId: string;
  type: string; // progress, volume, strength, recovery, balance, efficiency, motivation, risk
  insight: string;
  basedOn: string; // e.g., "Last 7 days", "Previous workout"
  priority: number; // 1-5, 5 being most important
  createdAt: string;
}

export interface GetRecentInsightsResponse {
  insights: WorkoutInsight[];
}

export interface GenerateInsightsResponse {
  insights: WorkoutInsight[];
} 