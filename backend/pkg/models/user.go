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
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	AccessToken  string             `bson:"access_token" json:"access_token"`
	RefreshToken string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	ExpiresAt    time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type Exercise struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description" json:"description"`
	MuscleGroup  string             `bson:"muscle_group" json:"muscleGroup"`
	Equipment    string             `bson:"equipment" json:"equipment"`
	Difficulty   string             `bson:"difficulty" json:"difficulty"`
	Instructions []string           `bson:"instructions" json:"instructions"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
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
	IsActive        bool                     `bson:"is_active" json:"is_active"`
	CreatedAt       time.Time                `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time                `bson:"updated_at" json:"updated_at"`
}
