package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VolumeMetrics struct {
	TotalVolumeLoad   float64            `bson:"totalVolumeLoad" json:"totalVolumeLoad"`
	Tonnage           float64            `bson:"tonnage" json:"tonnage"`
	RelativeVolume    float64            `bson:"relativeVolume" json:"relativeVolume"`
	EffectiveReps     int32              `bson:"effectiveReps" json:"effectiveReps"`
	HardSets          int32              `bson:"hardSets" json:"hardSets"`
	MuscleGroupVolume map[string]float64 `bson:"muscleGroupVolume" json:"muscleGroupVolume"`
}

type PerformanceMetrics struct {
	AverageRPE      float64         `bson:"averageRpe" json:"averageRpe"`
	MaxRPE          int32           `bson:"maxRpe" json:"maxRpe"`
	RPEDistribution map[int32]int32 `bson:"rpeDistribution" json:"rpeDistribution"`
}

type IntensityMetrics struct {
	AverageIntensity      float64            `bson:"averageIntensity" json:"averageIntensity"`           // Avg Intensity = Σ(Weight × Reps) / Σ(1RM × Reps) × 100
	RelativeIntensity     float64            `bson:"relativeIntensity" json:"relativeIntensity"`         // Weight Lifted / Body Weight
	RPEAdjustedLoad       float64            `bson:"rpeAdjustedLoad" json:"rpeAdjustedLoad"`             // Weight × (RPE/10)^2
	IntensityDistribution map[string]float64 `bson:"intensityDistribution" json:"intensityDistribution"` // Zone percentages: Light <60%, Moderate 60-80%, Heavy >80%
	LoadDensity           float64            `bson:"loadDensity" json:"loadDensity"`                     // Total Volume / Total Workout Time (minutes)
}

type SetMetrics struct {
	TotalSets      int32   `bson:"totalSets" json:"totalSets"`
	CompletedSets  int32   `bson:"completedSets" json:"completedSets"`
	CompletionRate float64 `bson:"completionRate" json:"completionRate"`
}

type ExerciseMetrics struct {
	ExerciseID       primitive.ObjectID `bson:"exerciseId" json:"exerciseId"`
	ExerciseName     string             `bson:"exerciseName" json:"exerciseName"`
	VolumeLoad       float64            `bson:"volumeLoad" json:"volumeLoad"`
	Tonnage          float64            `bson:"tonnage" json:"tonnage"`
	SetCount         int32              `bson:"setCount" json:"setCount"`
	TotalReps        int32              `bson:"totalReps" json:"totalReps"`
	AverageWeight    float64            `bson:"averageWeight" json:"averageWeight"`
	MaxWeight        float64            `bson:"maxWeight" json:"maxWeight"`
	PrimaryMuscles   []string           `bson:"primaryMuscles" json:"primaryMuscles"`
	SecondaryMuscles []string           `bson:"secondaryMuscles" json:"secondaryMuscles"`
}

type WorkoutMetrics struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID              primitive.ObjectID `bson:"userId" json:"userId"`
	SessionID           primitive.ObjectID `bson:"sessionId" json:"sessionId"`
	RoutineID           primitive.ObjectID `bson:"routineId" json:"routineId"`
	Date                time.Time          `bson:"date" json:"date"`
	VolumeMetrics       VolumeMetrics      `bson:"volumeMetrics" json:"volumeMetrics"`
	PerformanceMetrics  PerformanceMetrics `bson:"performanceMetrics" json:"performanceMetrics"`
	IntensityMetrics    IntensityMetrics   `bson:"intensityMetrics" json:"intensityMetrics"`
	SetMetrics          SetMetrics         `bson:"setMetrics" json:"setMetrics"`
	ExerciseMetrics     []ExerciseMetrics  `bson:"exerciseMetrics" json:"exerciseMetrics"`
	WorkoutDurationSecs int32              `bson:"workoutDurationSecs" json:"workoutDurationSecs"`
	CreatedAt           time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt           time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type PeriodMetrics struct {
	StartDate          time.Time          `bson:"startDate" json:"startDate"`
	EndDate            time.Time          `bson:"endDate" json:"endDate"`
	TotalWorkouts      int32              `bson:"totalWorkouts" json:"totalWorkouts"`
	TotalVolumeLoad    float64            `bson:"totalVolumeLoad" json:"totalVolumeLoad"`
	TotalTonnage       float64            `bson:"totalTonnage" json:"totalTonnage"`
	AverageWorkoutTime float64            `bson:"averageWorkoutTime" json:"averageWorkoutTime"`
	MuscleGroupVolume  map[string]float64 `bson:"muscleGroupVolume" json:"muscleGroupVolume"`
	VolumeProgression  []float64          `bson:"volumeProgression" json:"volumeProgression"`
}

type UserMetrics struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"userId" json:"userId"`
	WeeklyMetrics   PeriodMetrics      `bson:"weeklyMetrics" json:"weeklyMetrics"`
	MonthlyMetrics  PeriodMetrics      `bson:"monthlyMetrics" json:"monthlyMetrics"`
	AllTimeMetrics  PeriodMetrics      `bson:"allTimeMetrics" json:"allTimeMetrics"`
	VolumeLandmarks VolumeLandmarks    `bson:"volumeLandmarks" json:"volumeLandmarks"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type VolumeLandmarks struct {
	MEV map[string]float64 `bson:"mev" json:"mev"` // Minimum Effective Volume per muscle group
	MAV map[string]float64 `bson:"mav" json:"mav"` // Maximum Adaptive Volume per muscle group
	MRV map[string]float64 `bson:"mrv" json:"mrv"` // Maximum Recoverable Volume per muscle group
}

type VolumeProgressionPoint struct {
	Date       time.Time `bson:"date" json:"date"`
	VolumeLoad float64   `bson:"volumeLoad" json:"volumeLoad"`
	Tonnage    float64   `bson:"tonnage" json:"tonnage"`
}

type TrendMetrics struct {
	Period            string                   `bson:"period" json:"period"` // "weekly", "monthly"
	VolumeProgression []VolumeProgressionPoint `bson:"volumeProgression" json:"volumeProgression"`
	VolumeGrowthRate  float64                  `bson:"volumeGrowthRate" json:"volumeGrowthRate"`
	MuscleGroupTrends map[string][]float64     `bson:"muscleGroupTrends" json:"muscleGroupTrends"`
	PerformanceTrends []float64                `bson:"performanceTrends" json:"performanceTrends"`
}

type MuscleGroupConstants struct {
	PrimaryMultiplier   float64
	SecondaryMultiplier float64
}

var DefaultMuscleGroupWeights = MuscleGroupConstants{
	PrimaryMultiplier:   1.0,
	SecondaryMultiplier: 0.6,
}
