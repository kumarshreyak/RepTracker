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
	ID                             primitive.ObjectID             `bson:"_id,omitempty" json:"id"`
	UserID                         primitive.ObjectID             `bson:"userId" json:"userId"`
	SessionID                      primitive.ObjectID             `bson:"sessionId" json:"sessionId"`
	RoutineID                      primitive.ObjectID             `bson:"routineId" json:"routineId"`
	Date                           time.Time                      `bson:"date" json:"date"`
	VolumeMetrics                  VolumeMetrics                  `bson:"volumeMetrics" json:"volumeMetrics"`
	PerformanceMetrics             PerformanceMetrics             `bson:"performanceMetrics" json:"performanceMetrics"`
	IntensityMetrics               IntensityMetrics               `bson:"intensityMetrics" json:"intensityMetrics"`
	StrengthMetrics                StrengthMetrics                `bson:"strengthMetrics" json:"strengthMetrics"`
	ProgressAdaptationMetrics      ProgressAdaptationMetrics      `bson:"progressAdaptationMetrics" json:"progressAdaptationMetrics"`
	RecoveryFatigueMetrics         RecoveryFatigueMetrics         `bson:"recoveryFatigueMetrics" json:"recoveryFatigueMetrics"`
	BodyCompositionMetrics         BodyCompositionMetrics         `bson:"bodyCompositionMetrics" json:"bodyCompositionMetrics"`
	MuscleSpecificMetrics          MuscleSpecificMetrics          `bson:"muscleSpecificMetrics" json:"muscleSpecificMetrics"`
	WorkCapacityMetrics            WorkCapacityMetrics            `bson:"workCapacityMetrics" json:"workCapacityMetrics"`
	TrainingPatternMetrics         TrainingPatternMetrics         `bson:"trainingPatternMetrics" json:"trainingPatternMetrics"`
	PeriodizationMetrics           PeriodizationMetrics           `bson:"periodizationMetrics" json:"periodizationMetrics"`
	InjuryRiskPreventionMetrics    InjuryRiskPreventionMetrics    `bson:"injuryRiskPreventionMetrics" json:"injuryRiskPreventionMetrics"`
	EfficiencyTechniqueMetrics     EfficiencyTechniqueMetrics     `bson:"efficiencyTechniqueMetrics" json:"efficiencyTechniqueMetrics"`
	ComparativeNormativeMetrics    ComparativeNormativeMetrics    `bson:"comparativeNormativeMetrics" json:"comparativeNormativeMetrics"`
	PsychologicalBehavioralMetrics PsychologicalBehavioralMetrics `bson:"psychologicalBehavioralMetrics" json:"psychologicalBehavioralMetrics"`
	TimeBasedAnalyticsMetrics      TimeBasedAnalyticsMetrics      `bson:"timeBasedAnalyticsMetrics" json:"timeBasedAnalyticsMetrics"`
	CompoundMetrics                CompoundMetrics                `bson:"compoundMetrics" json:"compoundMetrics"`
	PredictiveAnalyticsMetrics     PredictiveAnalyticsMetrics     `bson:"predictiveAnalyticsMetrics" json:"predictiveAnalyticsMetrics"`
	SetMetrics                     SetMetrics                     `bson:"setMetrics" json:"setMetrics"`
	ExerciseMetrics                []ExerciseMetrics              `bson:"exerciseMetrics" json:"exerciseMetrics"`
	WorkoutDurationSecs            int32                          `bson:"workoutDurationSecs" json:"workoutDurationSecs"`
	CreatedAt                      time.Time                      `bson:"createdAt" json:"createdAt"`
	UpdatedAt                      time.Time                      `bson:"updatedAt" json:"updatedAt"`
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

// Injury Risk & Prevention Metrics
type InjuryRiskPreventionMetrics struct {
	InjuryRiskScore      float64 `bson:"injuryRiskScore" json:"injuryRiskScore"`           // IRS = (ACWR × Imbalance Index × Fatigue Score) / Recovery Quality, high risk if > 1.5
	LoadSpikeAlert       bool    `bson:"loadSpikeAlert" json:"loadSpikeAlert"`             // Spike = True if Weekly Volume > 1.5 × Average of Last 4 Weeks
	AsymmetryDevelopment float64 `bson:"asymmetryDevelopment" json:"asymmetryDevelopment"` // Asymmetry = |(Left Performance - Right Performance)| / Average × 100
}

// Efficiency & Technique Metrics
type EfficiencyTechniqueMetrics struct {
	StrengthEfficiency        float64 `bson:"strengthEfficiency" json:"strengthEfficiency"`               // SE = (1RM Gain / Total Volume) × 1000
	VolumeEfficiency          float64 `bson:"volumeEfficiency" json:"volumeEfficiency"`                   // VE = Performance Improvement / Total Training Volume
	RPEPerformanceCorrelation float64 `bson:"rpePerformanceCorrelation" json:"rpePerformanceCorrelation"` // Correlation = Pearson(Actual Reps, 10 - RPE + Expected Reps at RPE)
	TechniqueConsistency      float64 `bson:"techniqueConsistency" json:"techniqueConsistency"`           // Consistency = 1 - (Standard Deviation of Rep Times / Average Rep Time)
}

// Constants for Injury Risk & Prevention Metrics
const (
	InjuryRiskThreshold    = 1.5  // High risk if injury risk score > 1.5
	LoadSpikeMultiplier    = 1.5  // Load spike if weekly volume > 1.5x 4-week average
	AsymmetryRiskThreshold = 15.0 // High asymmetry risk if > 15%
)

// Constants for Efficiency & Technique Metrics
const (
	StrengthEfficiencyMultiplier = 1000.0 // Multiplier for strength efficiency calculation
	MinCorrelationSampleSize     = 5      // Minimum sets needed for correlation calculation
	DefaultRepTime               = 3.0    // Default rep time in seconds for technique consistency
)

// Comparative & Normative Metrics
type ComparativeNormativeMetrics struct {
	PercentileRanking           float64            `bson:"percentileRanking" json:"percentileRanking"`                     // Percentile = (Number of people below score / Total people) × 100
	TrainingAgeAdjustedExpected map[string]float64 `bson:"trainingAgeAdjustedExpected" json:"trainingAgeAdjustedExpected"` // Exercise ID -> Expected Strength = Base Strength × (1 + 0.1 × sqrt(Training Years))
	GeneticPotentialEstimate    map[string]float64 `bson:"geneticPotentialEstimate" json:"geneticPotentialEstimate"`       // Exercise ID -> Potential = Current Max × (1.5 - 0.5 × (Current/Elite Standard))
	PerformanceCategory         string             `bson:"performanceCategory" json:"performanceCategory"`                 // Beginner, Intermediate, Advanced, Elite based on percentile and training age
	RelativeToExpectations      float64            `bson:"relativeToExpectations" json:"relativeToExpectations"`           // Current Performance / Expected Performance for training age
}

// Psychological & Behavioral Metrics
type PsychologicalBehavioralMetrics struct {
	RPEAccuracy       float64 `bson:"rpeAccuracy" json:"rpeAccuracy"`             // RPE Accuracy = 1 - |Predicted Reps at RPE - Actual Reps| / Actual Reps
	MotivationIndex   float64 `bson:"motivationIndex" json:"motivationIndex"`     // MI = (Consistency × Voluntary Extra Sets × (10 - Avg RPE)) / 100
	BurnoutRiskScore  float64 `bson:"burnoutRiskScore" json:"burnoutRiskScore"`   // BRS = (Decreasing Performance + Increasing RPE + Decreased Frequency) / 3
	TrainingAdherence float64 `bson:"trainingAdherence" json:"trainingAdherence"` // Adherence = (Completed Workouts / Planned Workouts) × 100
	ConsistencyTrend  float64 `bson:"consistencyTrend" json:"consistencyTrend"`   // Trend in workout consistency over time (positive = improving)
}

// Time-Based Analytics Metrics
type TimeBasedAnalyticsMetrics struct {
	OptimalTrainingTimeHour   int32              `bson:"optimalTrainingTimeHour" json:"optimalTrainingTimeHour"`     // Hour of day with best performance (0-23)
	PerformanceByHour         map[string]float64 `bson:"performanceByHour" json:"performanceByHour"`                 // Hour -> Performance Score = (1RM% × Completion Rate × (10 - RPE))
	OptimalRestPeriods        map[string]float64 `bson:"optimalRestPeriods" json:"optimalRestPeriods"`               // Exercise ID -> Optimal Rest = 2^(Intensity%/25) minutes
	SessionDurationEfficiency float64            `bson:"sessionDurationEfficiency" json:"sessionDurationEfficiency"` // SDE = Total Effective Volume / Session Duration
	TimeOfDayPreference       string             `bson:"timeOfDayPreference" json:"timeOfDayPreference"`             // Morning, Afternoon, Evening based on performance patterns
	WorkoutTimingConsistency  float64            `bson:"workoutTimingConsistency" json:"workoutTimingConsistency"`   // Consistency in workout timing (1.0 = always same time, 0.0 = random)
}

// Performance categorization thresholds for comparative metrics
type PerformanceThresholds struct {
	BeginnerPercentile     float64 // Below this percentile = Beginner
	IntermediatePercentile float64 // Below this percentile = Intermediate
	AdvancedPercentile     float64 // Below this percentile = Advanced
	// Above Advanced = Elite
}

var DefaultPerformanceThresholds = PerformanceThresholds{
	BeginnerPercentile:     25.0, // Bottom 25%
	IntermediatePercentile: 60.0, // Bottom 60%
	AdvancedPercentile:     85.0, // Bottom 85%
}

// Elite strength standards for genetic potential calculation (in kg for major lifts)
// Based on IPF world records and elite powerlifting standards adjusted for body weight
type EliteStandards struct {
	SquatCoefficient    float64 // Multiplier × body weight = elite squat
	BenchCoefficient    float64 // Multiplier × body weight = elite bench
	DeadliftCoefficient float64 // Multiplier × body weight = elite deadlift
	OverheadCoefficient float64 // Multiplier × body weight = elite overhead press
}

var EliteStandardsMale = EliteStandards{
	SquatCoefficient:    2.5, // Elite male: 2.5x body weight squat
	BenchCoefficient:    2.0, // Elite male: 2.0x body weight bench
	DeadliftCoefficient: 3.0, // Elite male: 3.0x body weight deadlift
	OverheadCoefficient: 1.5, // Elite male: 1.5x body weight overhead press
}

var EliteStandardsFemale = EliteStandards{
	SquatCoefficient:    2.0, // Elite female: 2.0x body weight squat
	BenchCoefficient:    1.5, // Elite female: 1.5x body weight bench
	DeadliftCoefficient: 2.5, // Elite female: 2.5x body weight deadlift
	OverheadCoefficient: 1.0, // Elite female: 1.0x body weight overhead press
}

// Age and gender categories for percentile ranking
type DemographicCategory struct {
	AgeGroup string // "18-29", "30-39", "40-49", "50-59", "60+"
	Gender   string // "male", "female"
	Weight   string // "light", "medium", "heavy" based on body weight
}

// Constants for Psychological & Behavioral Metrics
const (
	MinWorkoutsForTrend      = 5    // Minimum workouts needed for trend analysis
	BurnoutRiskThreshold     = 0.7  // High burnout risk if score > 0.7
	LowMotivationThreshold   = 0.3  // Low motivation if index < 0.3
	OptimalConsistencyTarget = 0.85 // Target consistency rate (85%)
	VoluntaryExtraSetsBonus  = 1.2  // Bonus multiplier for extra sets
)

// Constants for Time-Based Analytics
const (
	MinSessionsForTimeAnalysis = 10   // Minimum sessions needed for time-of-day analysis
	OptimalRestTimeMultiplier  = 2.0  // Base multiplier for rest time calculation
	RestTimeIntensityDivisor   = 25.0 // Divisor for intensity in rest time formula
	MinRestTimeMins            = 1.0  // Minimum rest time in minutes
	MaxRestTimeMins            = 10.0 // Maximum rest time in minutes
	TimingConsistencyWindow    = 2.0  // Hours window for considering "consistent" timing
)

// Age category mappings for demographic analysis
var AgeCategoryMap = map[int]string{
	18: "18-29", 19: "18-29", 20: "18-29", 21: "18-29", 22: "18-29", 23: "18-29", 24: "18-29", 25: "18-29", 26: "18-29", 27: "18-29", 28: "18-29", 29: "18-29",
	30: "30-39", 31: "30-39", 32: "30-39", 33: "30-39", 34: "30-39", 35: "30-39", 36: "30-39", 37: "30-39", 38: "30-39", 39: "30-39",
	40: "40-49", 41: "40-49", 42: "40-49", 43: "40-49", 44: "40-49", 45: "40-49", 46: "40-49", 47: "40-49", 48: "40-49", 49: "40-49",
	50: "50-59", 51: "50-59", 52: "50-59", 53: "50-59", 54: "50-59", 55: "50-59", 56: "50-59", 57: "50-59", 58: "50-59", 59: "50-59",
}

// Weight category mappings (in kg)
var WeightCategoryMap = map[string][2]float64{
	"light":  {0, 70},    // Under 70kg
	"medium": {70, 90},   // 70-90kg
	"heavy":  {90, 1000}, // Over 90kg
}

// Exercise name mappings for elite standards lookup
var ExerciseToStandardMap = map[string]string{
	"squat":          "squat",
	"back squat":     "squat",
	"front squat":    "squat",
	"bench press":    "bench",
	"incline bench":  "bench",
	"deadlift":       "deadlift",
	"sumo deadlift":  "deadlift",
	"overhead press": "overhead",
	"military press": "overhead",
	"push press":     "overhead",
}

// Compound Metrics
type CompoundMetrics struct {
	FitnessAge                 float64 `bson:"fitnessAge" json:"fitnessAge"`                                 // Fitness Age = Chronological Age × (Population Average Score / User Score)
	OverallFitnessScore        float64 `bson:"overallFitnessScore" json:"overallFitnessScore"`               // OFS = (Strength Score × 0.3 + Volume Score × 0.2 + Consistency × 0.2 + Progress × 0.3) × 100
	TrainingEfficiencyQuotient float64 `bson:"trainingEfficiencyQuotient" json:"trainingEfficiencyQuotient"` // TEQ = (Performance Gains / (Time Invested × Fatigue Generated)) × 100
}

// Predictive Analytics Metrics
type PredictiveAnalyticsMetrics struct {
	PlateauProbability          float64            `bson:"plateauProbability" json:"plateauProbability"`                   // Plateau Probability = 1 / (1 + e^(Progress Rate - 0.02))
	PerformanceTrajectory       map[string]float64 `bson:"performanceTrajectory" json:"performanceTrajectory"`             // Exercise ID -> Future Performance = Current + (Weekly Gain Rate × Weeks × Decay Factor^Weeks)
	GoalAchievementTimeline     map[string]float64 `bson:"goalAchievementTimeline" json:"goalAchievementTimeline"`         // Exercise ID -> Weeks to Goal = (Goal - Current) / (Average Weekly Progress × Expected Decay)
	DetrainingRisk              float64            `bson:"detrainingRisk" json:"detrainingRisk"`                           // Detraining Risk = Days Since Last Workout / (CTL / 10)
	ProjectedPerformance4Weeks  map[string]float64 `bson:"projectedPerformance4Weeks" json:"projectedPerformance4Weeks"`   // Exercise ID -> 4-week performance projection
	ProjectedPerformance12Weeks map[string]float64 `bson:"projectedPerformance12Weeks" json:"projectedPerformance12Weeks"` // Exercise ID -> 12-week performance projection
	ProjectedPerformance26Weeks map[string]float64 `bson:"projectedPerformance26Weeks" json:"projectedPerformance26Weeks"` // Exercise ID -> 26-week performance projection
}

// Constants for Compound Metrics
const (
	// Fitness Age calculation constants
	DefaultPopulationAverageScore = 75.0 // Population average fitness score baseline

	// Overall Fitness Score weights
	StrengthScoreWeight    = 0.3 // 30% weight for strength component
	VolumeScoreWeight      = 0.2 // 20% weight for volume component
	ConsistencyScoreWeight = 0.2 // 20% weight for consistency component
	ProgressScoreWeight    = 0.3 // 30% weight for progress component

	// Training Efficiency Quotient constants
	TEQMultiplier = 100.0 // Multiplier for TEQ calculation
)

// Constants for Predictive Analytics
const (
	// Plateau prediction constants
	PlateauThreshold = 0.02 // 2% progress rate threshold for plateau probability

	// Performance trajectory constants
	PerformanceDecayFactor = 0.95 // Weekly decay factor for performance gains (95% retention)

	// Goal achievement constants
	MinWeeklyProgress = 0.001 // Minimum weekly progress to avoid division by zero

	// Detraining risk constants
	CTLDivisor = 10.0 // Divisor for CTL in detraining risk calculation

	// Projection time horizons (in weeks)
	ShortTermProjection  = 4  // 4 weeks
	MediumTermProjection = 12 // 12 weeks (3 months)
	LongTermProjection   = 26 // 26 weeks (6 months)
)
