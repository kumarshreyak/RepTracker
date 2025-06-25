package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Name      string             `bson:"name" json:"name"`
	GoogleID  string             `bson:"google_id,omitempty" json:"google_id,omitempty"`
	Picture   string             `bson:"picture,omitempty" json:"picture,omitempty"`
	Height    float64            `bson:"height,omitempty" json:"height,omitempty"` // Height in cm
	Weight    float64            `bson:"weight,omitempty" json:"weight,omitempty"` // Weight in kg
	Age       int32              `bson:"age,omitempty" json:"age,omitempty"`       // Age in years
	Goal      string             `bson:"goal,omitempty" json:"goal,omitempty"`     // Goal: lose_fat, gain_muscle, maintain
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Goal constants
const (
	GoalLoseFat    = "lose_fat"
	GoalGainMuscle = "gain_muscle"
	GoalMaintain   = "maintain"
)

type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	AccessToken  string             `bson:"access_token" json:"access_token"`
	RefreshToken string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	ExpiresAt    time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type Exercise struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Description      string             `bson:"description" json:"description"`
	Category         string             `bson:"category" json:"category"`
	Equipment        []string           `bson:"equipment" json:"equipment"`
	PrimaryMuscles   []string           `bson:"primaryMuscles" json:"primaryMuscles"`
	SecondaryMuscles []string           `bson:"secondaryMuscles" json:"secondaryMuscles"`
	Instructions     []string           `bson:"instructions" json:"instructions"`
	Video            string             `bson:"video,omitempty" json:"video,omitempty"`
	VariationsOn     []string           `bson:"variationsOn,omitempty" json:"variationsOn,omitempty"`
	VariationOn      []string           `bson:"variationOn,omitempty" json:"variationOn,omitempty"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type WorkoutSet struct {
	Reps            int32   `bson:"reps" json:"reps"`
	Weight          float32 `bson:"weight" json:"weight"`
	DurationSeconds int32   `bson:"duration_seconds,omitempty" json:"duration_seconds,omitempty"`
	Distance        float32 `bson:"distance,omitempty" json:"distance,omitempty"`
	Notes           string  `bson:"notes,omitempty" json:"notes,omitempty"`
}

type WorkoutExercise struct {
	ExerciseID  primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	Sets        []WorkoutSet       `bson:"sets" json:"sets"`
	Notes       string             `bson:"notes,omitempty" json:"notes,omitempty"`
	RestSeconds int32              `bson:"rest_seconds,omitempty" json:"rest_seconds,omitempty"`
}

type Workout struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name            string             `bson:"name" json:"name"`
	Description     string             `bson:"description,omitempty" json:"description,omitempty"`
	Exercises       []WorkoutExercise  `bson:"exercises" json:"exercises"`
	StartedAt       *time.Time         `bson:"started_at,omitempty" json:"started_at,omitempty"`
	FinishedAt      *time.Time         `bson:"finished_at,omitempty" json:"finished_at,omitempty"`
	DurationSeconds int32              `bson:"duration_seconds,omitempty" json:"duration_seconds,omitempty"`
	Notes           string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

type WorkoutSessionSet struct {
	TargetReps      int32      `bson:"target_reps" json:"target_reps"`
	TargetWeight    float32    `bson:"target_weight" json:"target_weight"`
	ActualReps      int32      `bson:"actual_reps" json:"actual_reps"`
	ActualWeight    float32    `bson:"actual_weight" json:"actual_weight"`
	DurationSeconds int32      `bson:"duration_seconds,omitempty" json:"duration_seconds,omitempty"`
	Distance        float32    `bson:"distance,omitempty" json:"distance,omitempty"`
	Notes           string     `bson:"notes,omitempty" json:"notes,omitempty"`
	Completed       bool       `bson:"completed" json:"completed"`
	StartedAt       *time.Time `bson:"started_at,omitempty" json:"started_at,omitempty"`
	FinishedAt      *time.Time `bson:"finished_at,omitempty" json:"finished_at,omitempty"`
}

type WorkoutSessionExercise struct {
	ExerciseID  primitive.ObjectID  `bson:"exercise_id" json:"exercise_id"`
	Sets        []WorkoutSessionSet `bson:"sets" json:"sets"`
	Notes       string              `bson:"notes,omitempty" json:"notes,omitempty"`
	RestSeconds int32               `bson:"rest_seconds,omitempty" json:"rest_seconds,omitempty"`
	Completed   bool                `bson:"completed" json:"completed"`
	StartedAt   *time.Time          `bson:"started_at,omitempty" json:"started_at,omitempty"`
	FinishedAt  *time.Time          `bson:"finished_at,omitempty" json:"finished_at,omitempty"`
}

type WorkoutSession struct {
	ID              primitive.ObjectID       `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID       `bson:"user_id" json:"user_id"`
	RoutineID       primitive.ObjectID       `bson:"routine_id" json:"routine_id"`
	Name            string                   `bson:"name" json:"name"`
	Description     string                   `bson:"description,omitempty" json:"description,omitempty"`
	Exercises       []WorkoutSessionExercise `bson:"exercises" json:"exercises"`
	StartedAt       time.Time                `bson:"started_at" json:"started_at"`
	FinishedAt      *time.Time               `bson:"finished_at,omitempty" json:"finished_at,omitempty"`
	DurationSeconds int32                    `bson:"duration_seconds" json:"duration_seconds"`
	Notes           string                   `bson:"notes,omitempty" json:"notes,omitempty"`
	RPERating       int32                    `bson:"rpe_rating,omitempty" json:"rpe_rating,omitempty"`
	IsActive        bool                     `bson:"is_active" json:"is_active"`
	CreatedAt       time.Time                `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time                `bson:"updated_at" json:"updated_at"`
}

// Stored Workout Suggestions
type StoredWorkoutChange struct {
	Type         string `bson:"type" json:"type"` // "exercise", "sets", "reps", "weight", "rest", "remove", "add"
	ExerciseID   string `bson:"exerciseId" json:"exerciseId"`
	ExerciseName string `bson:"exerciseName" json:"exerciseName"`
	OldValue     string `bson:"oldValue" json:"oldValue"`
	NewValue     string `bson:"newValue" json:"newValue"`
	Reason       string `bson:"reason" json:"reason"`
}

type StoredSuggestedWorkout struct {
	OriginalWorkoutID string                `bson:"originalWorkoutId" json:"originalWorkoutId"`
	Name              string                `bson:"name" json:"name"`
	Description       string                `bson:"description" json:"description"`
	Exercises         []WorkoutExercise     `bson:"exercises" json:"exercises"`
	Changes           []StoredWorkoutChange `bson:"changes" json:"changes"`
	OverallReasoning  string                `bson:"overallReasoning" json:"overallReasoning"`
	Priority          int32                 `bson:"priority" json:"priority"`
}

type StoredWorkoutSuggestion struct {
	ID              primitive.ObjectID       `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID       `bson:"userId" json:"userId"`
	Suggestions     []StoredSuggestedWorkout `bson:"suggestions" json:"suggestions"`
	AnalysisSummary string                   `bson:"analysisSummary" json:"analysisSummary"`
	DaysAnalyzed    int32                    `bson:"daysAnalyzed" json:"daysAnalyzed"`
	CreatedAt       time.Time                `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time                `bson:"updatedAt" json:"updatedAt"`
}
