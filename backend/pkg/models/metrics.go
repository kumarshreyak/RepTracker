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

type StrengthMetrics struct {
	EstimatedOneRMEpley   map[string]float64 `bson:"estimatedOneRmEpley" json:"estimatedOneRmEpley"`     // Exercise ID -> 1RM using Epley formula
	EstimatedOneRMBrzycki map[string]float64 `bson:"estimatedOneRmBrzycki" json:"estimatedOneRmBrzycki"` // Exercise ID -> 1RM using Brzycki formula
	WilksScore            float64            `bson:"wilksScore" json:"wilksScore"`                       // Wilks coefficient for powerlifting comparison
	PushPullRatio         float64            `bson:"pushPullRatio" json:"pushPullRatio"`                 // Push volume / Pull volume ratio
	PowerOutput           float64            `bson:"powerOutput" json:"powerOutput"`                     // (Weight × Distance × Reps) / Time
}

type ProgressAdaptationMetrics struct {
	ProgressiveOverloadIndex float64                  `bson:"progressiveOverloadIndex" json:"progressiveOverloadIndex"` // POI = (Current Week Volume × Avg Intensity) / (Previous Week Volume × Avg Intensity)
	WeekOverWeekProgressRate map[string]float64       `bson:"weekOverWeekProgressRate" json:"weekOverWeekProgressRate"` // Exercise ID -> Progress Rate = (Current 1RM - Previous 1RM) / Previous 1RM × 100
	PlateauDetection         map[string]PlateauStatus `bson:"plateauDetection" json:"plateauDetection"`                 // Exercise ID -> Plateau status
	StrengthGainVelocity     map[string]float64       `bson:"strengthGainVelocity" json:"strengthGainVelocity"`         // Exercise ID -> SGV = (Current 1RM - Initial 1RM) / Training Weeks
	AdaptationRate           map[string]float64       `bson:"adaptationRate" json:"adaptationRate"`                     // Exercise ID -> Adaptation Rate = Δ Performance / Δ Volume
}

type PlateauStatus struct {
	IsPlateaued        bool    `bson:"isPlateaued" json:"isPlateaued"`               // True if Progress Rate < 1% for 3+ consecutive weeks
	ConsecutiveWeeks   int32   `bson:"consecutiveWeeks" json:"consecutiveWeeks"`     // Number of consecutive weeks with <1% progress
	LastProgressRate   float64 `bson:"lastProgressRate" json:"lastProgressRate"`     // Most recent progress rate
	WeeksSinceProgress int32   `bson:"weeksSinceProgress" json:"weeksSinceProgress"` // Weeks since last significant progress
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
	ID                        primitive.ObjectID        `bson:"_id,omitempty" json:"id"`
	UserID                    primitive.ObjectID        `bson:"userId" json:"userId"`
	SessionID                 primitive.ObjectID        `bson:"sessionId" json:"sessionId"`
	RoutineID                 primitive.ObjectID        `bson:"routineId" json:"routineId"`
	Date                      time.Time                 `bson:"date" json:"date"`
	VolumeMetrics             VolumeMetrics             `bson:"volumeMetrics" json:"volumeMetrics"`
	PerformanceMetrics        PerformanceMetrics        `bson:"performanceMetrics" json:"performanceMetrics"`
	IntensityMetrics          IntensityMetrics          `bson:"intensityMetrics" json:"intensityMetrics"`
	StrengthMetrics           StrengthMetrics           `bson:"strengthMetrics" json:"strengthMetrics"`
	ProgressAdaptationMetrics ProgressAdaptationMetrics `bson:"progressAdaptationMetrics" json:"progressAdaptationMetrics"`
	RecoveryFatigueMetrics    RecoveryFatigueMetrics    `bson:"recoveryFatigueMetrics" json:"recoveryFatigueMetrics"`
	BodyCompositionMetrics    BodyCompositionMetrics    `bson:"bodyCompositionMetrics" json:"bodyCompositionMetrics"`
	MuscleSpecificMetrics     MuscleSpecificMetrics     `bson:"muscleSpecificMetrics" json:"muscleSpecificMetrics"`
	WorkCapacityMetrics       WorkCapacityMetrics       `bson:"workCapacityMetrics" json:"workCapacityMetrics"`
	TrainingPatternMetrics    TrainingPatternMetrics    `bson:"trainingPatternMetrics" json:"trainingPatternMetrics"`
	PeriodizationMetrics      PeriodizationMetrics      `bson:"periodizationMetrics" json:"periodizationMetrics"`
	SetMetrics                SetMetrics                `bson:"setMetrics" json:"setMetrics"`
	ExerciseMetrics           []ExerciseMetrics         `bson:"exerciseMetrics" json:"exerciseMetrics"`
	WorkoutDurationSecs       int32                     `bson:"workoutDurationSecs" json:"workoutDurationSecs"`
	CreatedAt                 time.Time                 `bson:"createdAt" json:"createdAt"`
	UpdatedAt                 time.Time                 `bson:"updatedAt" json:"updatedAt"`
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

// Push/Pull muscle group classifications for strength ratios
var PushMuscles = map[string]bool{
	"chest":     true,
	"triceps":   true,
	"shoulders": true,
}

var PullMuscles = map[string]bool{
	"lats":        true,
	"middle back": true,
	"rhomboids":   true,
	"traps":       true,
	"biceps":      true,
}

// Wilks coefficient constants for male/female
type WilksCoefficients struct {
	A, B, C, D, E, F float64
}

var WilksMaleCoefficients = WilksCoefficients{
	A: -216.0475144,
	B: 16.2606339,
	C: -0.002388645,
	D: -0.00113732,
	E: 7.01863e-06,
	F: -1.291e-08,
}

var WilksFemaleCoefficients = WilksCoefficients{
	A: 594.31747775582,
	B: -27.23842536447,
	C: 0.82112226871,
	D: -0.00930733913,
	E: 4.731582e-05,
	F: -9.054e-08,
}

// Recovery & Fatigue Metrics
type RecoveryFatigueMetrics struct {
	AcuteChronicWorkloadRatio float64 `bson:"acuteChronicWorkloadRatio" json:"acuteChronicWorkloadRatio"` // ACWR = Acute Load (7 days) / Chronic Load (28 days), optimal: 0.8-1.3
	TrainingStrainScore       float64 `bson:"trainingStrainScore" json:"trainingStrainScore"`             // TSS = (Volume × Intensity × RPE) / 1000
	RecoveryNeedIndex         float64 `bson:"recoveryNeedIndex" json:"recoveryNeedIndex"`                 // RNI = (Volume × Intensity% × RPE) / (Hours Since Last Session × 10)
	FatigueAccumulationIndex  float64 `bson:"fatigueAccumulationIndex" json:"fatigueAccumulationIndex"`   // FAI = Σ(Daily TSS) - Σ(Daily Recovery Score)
	OvertrainingRiskScore     float64 `bson:"overtrainingRiskScore" json:"overtrainingRiskScore"`         // OTS = (ACWR × RPE Trend × Performance Decline) / Recovery Quality, high risk if > 2.0
	RecoveryScore             float64 `bson:"recoveryScore" json:"recoveryScore"`                         // Recovery Score = Hours Rest × Sleep Quality × Nutrition Factor
}

// Body Composition & Anthropometric Metrics
type BodyCompositionMetrics struct {
	BMI                    float64 `bson:"bmi" json:"bmi"`                                       // BMI = Weight (kg) / Height (m)²
	StrengthToWeightRatio  float64 `bson:"strengthToWeightRatio" json:"strengthToWeightRatio"`   // SWR = Total of Big 3 Lifts / Body Weight
	AllometricScalingScore float64 `bson:"allometricScalingScore" json:"allometricScalingScore"` // Scaled Strength = Lift Weight / Body Weight^(2/3)
	PonderalIndex          float64 `bson:"ponderalIndex" json:"ponderalIndex"`                   // PI = Weight (kg) / Height (m)³
}

// Default values for body composition calculations when user profile data is not available
type DefaultUserProfile struct {
	BodyWeight      float64 // kg
	Height          float64 // meters
	IsMale          bool
	SleepQuality    float64 // 1.0-10.0 scale
	NutritionFactor float64 // 1.0-10.0 scale
}

var DefaultProfile = DefaultUserProfile{
	BodyWeight:      75.0, // 75kg default
	Height:          1.75, // 1.75m default (175cm)
	IsMale:          true,
	SleepQuality:    7.0, // Good sleep quality
	NutritionFactor: 7.0, // Good nutrition
}

// Big 3 lifts for strength-to-weight ratio calculation
var Big3Exercises = map[string]bool{
	"squat":    true,
	"bench":    true,
	"deadlift": true,
}

// Muscle-Specific Metrics
type MuscleSpecificMetrics struct {
	MuscleGroupDistribution map[string]float64 `bson:"muscleGroupDistribution" json:"muscleGroupDistribution"` // MG Distribution% = (MG Volume / Total Volume) × 100
	MuscleImbalanceIndex    map[string]float64 `bson:"muscleImbalanceIndex" json:"muscleImbalanceIndex"`       // Imbalance = |Left Side Strength - Right Side Strength| / Average Strength × 100
	AntagonistRatio         map[string]float64 `bson:"antagonistRatio" json:"antagonistRatio"`                 // Antagonist Ratio = Antagonist Strength / Agonist Strength (e.g., Hamstring/Quad ratio ideal: 0.6-0.8)
	StimulusToFatigueRatio  map[string]float64 `bson:"stimulusToFatigueRatio" json:"stimulusToFatigueRatio"`   // SFR = (Performance Gain / Baseline) / (RPE × Volume)
}

// Work Capacity Metrics
type WorkCapacityMetrics struct {
	TotalWorkCapacity      float64 `bson:"totalWorkCapacity" json:"totalWorkCapacity"`           // TWC = Σ(Sets × Reps × Weight × (1 - Rest Time/300))
	DensityTrainingIndex   float64 `bson:"densityTrainingIndex" json:"densityTrainingIndex"`     // Density = Total Volume / Total Time
	DensityProgressPercent float64 `bson:"densityProgressPercent" json:"densityProgressPercent"` // Progress% = (Current Density - Previous Density) / Previous Density × 100
	TimeUnderTension       float64 `bson:"timeUnderTension" json:"timeUnderTension"`             // TUT = Σ(Reps × Tempo in seconds)
	MechanicalTensionScore float64 `bson:"mechanicalTensionScore" json:"mechanicalTensionScore"` // MTS = Weight × TUT × (RPE/10)
}

// Antagonist muscle group pairs for ratio calculations
var AntagonistPairs = map[string]string{
	"quadriceps":         "hamstrings",
	"hamstrings":         "quadriceps",
	"chest":              "middle back",
	"middle back":        "chest",
	"biceps":             "triceps",
	"triceps":            "biceps",
	"anterior deltoids":  "posterior deltoids",
	"posterior deltoids": "anterior deltoids",
}

// Left-Right muscle pairs for imbalance calculations
var LeftRightPairs = map[string]string{
	"left quadriceps":  "right quadriceps",
	"right quadriceps": "left quadriceps",
	"left hamstrings":  "right hamstrings",
	"right hamstrings": "left hamstrings",
	"left chest":       "right chest",
	"right chest":      "left chest",
	"left biceps":      "right biceps",
	"right biceps":     "left biceps",
	"left triceps":     "right triceps",
	"right triceps":    "left triceps",
}

// Default tempo values for exercises (in seconds: eccentric, pause, concentric, pause)
type ExerciseTempo struct {
	Eccentric  float64 // Lowering phase
	Pause1     float64 // Bottom pause
	Concentric float64 // Lifting phase
	Pause2     float64 // Top pause
}

var DefaultExerciseTempos = map[string]ExerciseTempo{
	"squat": {
		Eccentric:  2.0,
		Pause1:     1.0,
		Concentric: 1.0,
		Pause2:     0.0,
	},
	"bench press": {
		Eccentric:  2.0,
		Pause1:     1.0,
		Concentric: 1.0,
		Pause2:     0.0,
	},
	"deadlift": {
		Eccentric:  2.0,
		Pause1:     0.0,
		Concentric: 1.0,
		Pause2:     0.0,
	},
	"default": {
		Eccentric:  2.0,
		Pause1:     0.5,
		Concentric: 1.0,
		Pause2:     0.0,
	},
}

// Default rest time between sets (in seconds)
const DefaultRestTime = 180.0 // 3 minutes

// Training Pattern Analytics Metrics
type TrainingPatternMetrics struct {
	OptimalFrequency       float64 `bson:"optimalFrequency" json:"optimalFrequency"`             // Optimal Frequency = Recovery Time + 24-48 hours
	RecoveryTime           float64 `bson:"recoveryTime" json:"recoveryTime"`                     // Recovery Time = 24 × (Volume/10) × (Intensity/70) × (RPE/7)
	ExerciseSelectionIndex float64 `bson:"exerciseSelectionIndex" json:"exerciseSelectionIndex"` // Diversity Index = Unique Exercises / Total Exercises × 100
	ConsistencyScore       float64 `bson:"consistencyScore" json:"consistencyScore"`             // Consistency = (Actual Workouts / Planned Workouts) × 100
	WorkoutCompletionRate  float64 `bson:"workoutCompletionRate" json:"workoutCompletionRate"`   // Completion Rate = (Completed Sets / Planned Sets) × 100
}

// Advanced Periodization Metrics
type PeriodizationMetrics struct {
	ChronicTrainingLoad   float64 `bson:"chronicTrainingLoad" json:"chronicTrainingLoad"`     // CTL = Exponentially Weighted Average of daily TSS over 42 days
	AcuteTrainingLoad     float64 `bson:"acuteTrainingLoad" json:"acuteTrainingLoad"`         // ATL = Exponentially Weighted Average of daily TSS over 7 days
	TrainingStressBalance float64 `bson:"trainingStressBalance" json:"trainingStressBalance"` // TSB = CTL - ATL
	FormFreshnessIndex    float64 `bson:"formFreshnessIndex" json:"formFreshnessIndex"`       // FFI = TSB / CTL × 100
}

// Constants for Training Pattern Analytics
const (
	MinimumRecoveryHours = 24.0
	MaximumRecoveryHours = 48.0
	BaseVolumeThreshold  = 10.0
	BaseIntensityPercent = 70.0
	BaseRPEThreshold     = 7.0
)

// Constants for Periodization Metrics
const (
	ChronicTrainingLoadDays = 42   // 6 weeks
	AcuteTrainingLoadDays   = 7    // 1 week
	ExponentialDecayFactor  = 0.07 // Standard decay factor for training load
)
