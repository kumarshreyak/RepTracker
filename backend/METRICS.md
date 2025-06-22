# GymLog Metrics System Documentation

## Overview

The GymLog metrics system automatically calculates comprehensive workout analytics based on completed workout sessions. It provides volume-based metrics, performance tracking, and personalized training recommendations through Volume Landmarks (MEV/MAV/MRV).

## Table of Contents

- [Metrics Calculated](#metrics-calculated)
- [When Metrics Are Calculated](#when-metrics-are-calculated)
- [Data Models](#data-models)
- [API Endpoints](#api-endpoints)
- [Database Collections](#database-collections)
- [Implementation Details](#implementation-details)
- [Usage Examples](#usage-examples)

## Metrics Calculated

### Volume-Based Metrics

#### 1. Total Volume Load
**Formula**: `Σ(Sets × Reps × Weight)`
- Calculated for all completed sets in a workout
- Primary metric for tracking training volume
- Used in volume progression analysis

#### 2. Tonnage
**Formula**: `Σ(Reps × Weight)`
- Total weight moved during the workout
- Identical to Volume Load for individual sets
- Useful for comparing workouts of different set/rep schemes

#### 3. Volume per Muscle Group
**Formula**: `Σ(Sets × Reps × Weight × Muscle Group Multiplier)`
- **Primary Muscles**: 100% weight (multiplier = 1.0)
- **Secondary Muscles**: 60% weight (multiplier = 0.6)
- Tracks volume distribution across muscle groups
- Configurable via `models.DefaultMuscleGroupWeights`

#### 4. Relative Volume
**Formula**: `Total Volume Load / Body Weight`
- Currently equals Total Volume Load (body weight feature pending)
- Will normalize volume based on user's body weight when implemented

#### 5. Effective Reps
**Formula**: `Σ(Reps where RPE ≥ 7)`
- Counts reps from sets with high perceived exertion
- Uses session-level RPE rating
- Indicates quality training volume

#### 6. Hard Sets
**Formula**: `Count(Sets where RPE ≥ 7)`
- Number of challenging sets in the workout
- Uses session-level RPE rating
- Key indicator for training intensity

### Performance Metrics

#### 7. RPE Distribution
- Breakdown of sets by RPE rating
- Average RPE across all completed sets
- Maximum RPE in the session

#### 8. Set Completion Rate
**Formula**: `Completed Sets / Total Sets`
- Percentage of planned sets actually completed
- Indicates workout adherence

### Intensity Metrics

#### 9. Average Intensity
**Formula**: `Σ(Weight × Reps) / Σ(1RM × Reps) × 100`
- Percentage of 1RM lifted across all sets
- Uses Epley formula for 1RM estimation: `1RM = Weight × (1 + Reps/30)`
- Indicates training intensity relative to maximum strength

#### 10. Relative Intensity
**Formula**: `Max Weight Lifted / Body Weight`
- Strength relative to body weight
- Currently uses default body weight (75kg) - will use user profile when available
- Useful for comparing strength across different body weights

#### 11. RPE-Adjusted Load
**Formula**: `Σ(Weight × (RPE/10)²)`
- Weights adjusted for perceived exertion
- Accounts for subjective difficulty of each set
- Higher RPE increases the effective load exponentially

#### 12. Intensity Distribution
**Formula**: `(Sets in Zone / Total Sets) × 100`
- **Light Zone**: <60% of 1RM
- **Moderate Zone**: 60-80% of 1RM  
- **Heavy Zone**: >80% of 1RM
- Shows training zone distribution for periodization insights

#### 13. Load Density
**Formula**: `Total Volume Load / Workout Time (minutes)`
- Volume per unit time
- Indicates workout efficiency and pace
- Higher density suggests more intense, faster-paced training

### Strength & Performance Metrics

#### 14. Estimated 1RM (Epley Formula)
**Formula**: `1RM = Weight × (1 + Reps/30)`
- Classic formula for estimating one-rep maximum
- Calculated per exercise using heaviest set
- Accurate for rep ranges 1-10, less reliable above 10 reps

#### 15. Estimated 1RM (Brzycki Formula)
**Formula**: `1RM = Weight × (36 / (37 - Reps))`
- Alternative 1RM estimation formula
- Generally more conservative than Epley
- Invalid for reps ≥ 37 (falls back to Epley)
- Better for higher rep ranges (6-15 reps)

#### 16. Wilks Score
**Formula**: `Wilks = Weight Lifted × 500 / (a + b×BW + c×BW² + d×BW³ + e×BW⁴ + f×BW⁵)`
- Powerlifting comparison metric accounting for body weight
- Uses gender-specific coefficients (male/female)
- Allows comparison of strength across different body weights
- Based on heaviest lift in the session

**Male Coefficients:**
- a = -216.0475144, b = 16.2606339, c = -0.002388645
- d = -0.00113732, e = 7.01863e-06, f = -1.291e-08

**Female Coefficients:**
- a = 594.31747775582, b = -27.23842536447, c = 0.82112226871
- d = -0.00930733913, e = 4.731582e-05, f = -9.054e-08

#### 17. Push:Pull Ratio
**Formula**: `Σ(Push Volume) / Σ(Pull Volume)`
- Balance metric for upper body training
- **Push muscles**: chest, triceps, shoulders
- **Pull muscles**: lats, middle back, rhomboids, traps, biceps
- **Ideal range**: 1:1 to 1:1.5 (slightly favoring pull)
- Values >1.5 may indicate push dominance

#### 18. Power Output
**Formula**: `(Weight × Distance × Reps) / Time`
- Measures training power and explosiveness
- Currently uses simplified distance calculation (1m per rep)
- Time distributed across all sets in workout
- Higher values indicate more explosive, powerful training

### Volume Landmarks (MEV/MAV/MRV)

#### 19. Minimum Effective Volume (MEV)
**Formula**: `max(average × 0.75, average - std_deviation)`
- Minimum volume needed for progress
- Based on 8 weeks of historical data
- Falls back to 10 sets/week equivalent if insufficient data

#### 20. Maximum Adaptive Volume (MAV)
**Formula**: `average + (std_deviation × 0.5)`
- Volume at peak progress rate
- Sweet spot for muscle growth
- Falls back to 17.5 sets/week equivalent if insufficient data

#### 21. Maximum Recoverable Volume (MRV)
**Formula**: `average + (std_deviation × 1.5)`
- Maximum volume before performance decline
- Upper limit for training volume
- Falls back to 22.5 sets/week equivalent if insufficient data

### Muscle-Specific Metrics

#### 22. Muscle Group Volume Distribution
**Formula**: `MG Distribution% = (MG Volume / Total Volume) × 100`
- Percentage contribution of each muscle group to total workout volume
- Shows training balance across muscle groups
- Calculated from muscle group volume totals
- Helps identify muscle group emphasis and potential imbalances

#### 23. Muscle Imbalance Index
**Formula**: `Imbalance = |Left Side Strength - Right Side Strength| / Average Strength × 100`
- Percentage strength difference between left and right sides
- Calculated using estimated 1RM for bilateral exercises
- Higher values indicate greater imbalances (>15% may need attention)
- Tracks muscle groups: quadriceps, hamstrings, chest, biceps, triceps

#### 24. Antagonist Ratio
**Formula**: `Antagonist Ratio = Antagonist Strength / Agonist Strength`
- Strength ratio between opposing muscle groups
- **Key ratios**:
  - Hamstring/Quadriceps: ideal 0.6-0.8
  - Middle Back/Chest: ideal 0.8-1.2
  - Biceps/Triceps: ideal 0.8-1.0
  - Posterior/Anterior Deltoids: ideal 0.8-1.2
- Values outside ideal ranges may indicate imbalances

#### 25. Stimulus-to-Fatigue Ratio
**Formula**: `SFR = (Performance Gain / Baseline) / (RPE × Volume)`
- Efficiency of muscle group training stimulus
- Higher values indicate better stimulus with less fatigue
- Calculated per muscle group using volume changes vs previous session
- Helps optimize training efficiency and recovery

### Work Capacity Metrics

#### 26. Total Work Capacity
**Formula**: `TWC = Σ(Sets × Reps × Weight × (1 - Rest Time/300))`
- Accounts for rest periods in work capacity calculation
- Rest time capped at 300 seconds (5 minutes) for calculation
- Higher values indicate greater work capacity and conditioning
- Factors in workout density and rest efficiency

#### 27. Density Training Index
**Formula**: `Density = Total Volume / Total Time`
- Volume per unit time (volume/minute)
- Measures workout efficiency and pace
- Higher density indicates more volume completed in less time
- Useful for tracking conditioning improvements

#### 28. Density Progress Percent
**Formula**: `Progress% = (Current Density - Previous Density) / Previous Density × 100`
- Percentage change in training density compared to previous session
- Positive values indicate improved work capacity
- Tracks conditioning and efficiency improvements over time

#### 29. Time Under Tension (TUT)
**Formula**: `TUT = Σ(Reps × Tempo in seconds)`
- Total time muscles spend under load during workout
- Uses exercise-specific tempo patterns:
  - **Squat**: 2s eccentric, 1s pause, 1s concentric
  - **Bench Press**: 2s eccentric, 1s pause, 1s concentric  
  - **Deadlift**: 2s eccentric, 0s pause, 1s concentric
  - **Default**: 2s eccentric, 0.5s pause, 1s concentric
- Indicates muscle stimulation quality and training style

#### 30. Mechanical Tension Score
**Formula**: `MTS = Weight × TUT × (RPE/10)`
- Combines load, time under tension, and perceived difficulty
- Weighted by RPE to account for subjective intensity
- Higher scores indicate greater mechanical stimulus
- Useful for tracking hypertrophy-focused training effectiveness

### Training Pattern Analytics Metrics

#### 31. Optimal Training Frequency
**Formula**: `Optimal Frequency = Recovery Time + 24-48 hours`
- Recommends ideal time between training sessions
- Based on calculated recovery time and minimum rest period
- Ranges between Recovery Time + 24 hours to Recovery Time + 48 hours
- Helps optimize training frequency for individual recovery capacity

#### 32. Recovery Time
**Formula**: `Recovery Time = 24 × (Volume/10) × (Intensity/70) × (RPE/7)`
- Estimates time needed for full recovery from workout
- Scales with workout volume, intensity, and perceived exertion
- Base factors: 10 volume units, 70% intensity, RPE 7
- Used to calculate optimal training frequency

#### 33. Exercise Selection Diversity Index
**Formula**: `Diversity Index = (Unique Exercises / Total Exercises) × 100`
- Measures variety in exercise selection within a workout
- Values close to 100% indicate high exercise diversity
- Lower values suggest exercise repetition or specialization
- Useful for tracking program variety and movement patterns

#### 34. Consistency Score
**Formula**: `Consistency = (Actual Workouts / Planned Workouts) × 100`
- Measures adherence to planned training schedule
- Currently uses workout completion rate as proxy
- Will be enhanced when detailed workout planning is implemented
- Higher values indicate better training consistency

#### 35. Workout Completion Rate
**Formula**: `Completion Rate = (Completed Sets / Planned Sets) × 100`
- Percentage of planned sets actually completed
- Indicates workout adherence and execution quality
- Values close to 100% show strong workout completion
- Lower values may indicate fatigue, time constraints, or programming issues

### Advanced Periodization Metrics

#### 36. Chronic Training Load (CTL)
**Formula**: `CTL = Exponentially Weighted Average of daily TSS over 42 days`
- Measures long-term training stress and fitness level
- Based on 6 weeks (42 days) of Training Strain Score data
- Uses exponential decay to weight recent sessions more heavily
- Higher CTL indicates greater long-term training capacity

#### 37. Acute Training Load (ATL)
**Formula**: `ATL = Exponentially Weighted Average of daily TSS over 7 days`
- Measures recent training stress and fatigue level
- Based on 1 week (7 days) of Training Strain Score data
- Uses exponential decay with recent sessions weighted more heavily
- Higher ATL indicates greater recent training stress

#### 38. Training Stress Balance (TSB)
**Formula**: `TSB = CTL - ATL`
- Indicates training stress balance and freshness
- Positive values suggest freshness and readiness for hard training
- Negative values indicate fatigue and need for recovery
- **Optimal range**: -10 to +25 for most athletes

#### 39. Form/Freshness Index (FFI)
**Formula**: `FFI = (TSB / CTL) × 100`
- Percentage-based indicator of form and freshness
- Positive values indicate good form and readiness for performance
- Negative values suggest fatigue and need for recovery
- **Optimal range**: -5% to +15% for peak performance readiness

### Injury Risk & Prevention Metrics

#### 40. Injury Risk Score (IRS)
**Formula**: `IRS = (ACWR × Imbalance Index × Fatigue Score) / Recovery Quality`
- Combines acute:chronic workload ratio, muscle imbalances, and fatigue factors
- Divided by recovery quality (sleep + nutrition average)
- **High risk**: IRS > 1.5 indicates elevated injury risk
- Helps identify when to reduce training intensity or focus on recovery

#### 41. Load Spike Alert
**Formula**: `Spike = True if Weekly Volume > 1.5 × Average of Last 4 Weeks`
- Detects sudden increases in training volume
- Based on 4-week rolling average of total volume load
- **Alert threshold**: Current week volume exceeds 150% of 4-week average
- Rapid volume increases are associated with higher injury risk

#### 42. Asymmetry Development
**Formula**: `Asymmetry = |(Left Performance - Right Performance)| / Average × 100`
- Measures strength imbalances between left and right sides
- Uses maximum imbalance found across all muscle groups
- **Risk threshold**: Asymmetry > 15% may indicate injury risk
- Tracks bilateral exercise performance differences

### Efficiency & Technique Metrics

#### 43. Strength Efficiency (SE)
**Formula**: `SE = (1RM Gain / Total Volume) × 1000`
- Measures strength gains relative to training volume invested
- Higher values indicate more efficient strength development
- Compares current session 1RM estimates to previous session
- Multiplied by 1000 for easier interpretation (e.g., 2.5 vs 0.0025)

#### 44. Volume Efficiency (VE)
**Formula**: `VE = Performance Improvement / Total Training Volume`
- Assesses how effectively training volume translates to performance
- Uses average performance improvement across all exercises
- Lower volume with same gains = higher efficiency
- Helps optimize training volume for individual response

#### 45. RPE-Performance Correlation
**Formula**: `Correlation = Pearson(Actual Reps, 10 - RPE + Expected Reps at RPE)`
- Measures consistency between perceived exertion and actual performance
- Strong positive correlation indicates good RPE calibration
- Weak correlation may suggest RPE mis-calibration or variable factors
- **Minimum sample size**: 5 sets required for meaningful correlation

#### 46. Technique Consistency
**Formula**: `Consistency = 1 - (Standard Deviation of Rep Times / Average Rep Time)`
- Measures consistency of movement execution across reps
- Based on simulated rep times using exercise tempo patterns
- Values closer to 1.0 indicate more consistent technique
- Helps identify exercises where technique may be breaking down

### Trend Metrics

#### 31. Volume Progression
- Historical volume data points over time
- Tracks workout-to-workout changes
- Supports weekly and monthly views

#### 32. Volume Growth Rate
**Formula**: `(Current Period Volume - Previous Period Volume) / Previous Period Volume × 100`
- Percentage change in volume between periods
- Indicates training progression trends

### Progress & Adaptation Metrics

#### 24. Progressive Overload Index (POI)
**Formula**: `POI = (Current Week Volume × Avg Intensity) / (Previous Week Volume × Avg Intensity)`
- Measures progressive overload implementation
- Values >1.0 indicate increasing training stimulus
- Values <1.0 may indicate deload or regression
- Simplified calculation using total volume as intensity proxy

#### 25. Week-over-Week Progress Rate
**Formula**: `Progress Rate = (Current 1RM - Previous 1RM) / Previous 1RM × 100`
- Percentage change in estimated 1RM between consecutive weeks
- Calculated per exercise using Epley formula estimates
- Positive values indicate strength gains
- Values <1% may indicate plateau or stagnation

#### 26. Plateau Detection
**Criteria**: `Plateau = True if Progress Rate < 1% for 3+ consecutive weeks`
- Automated detection of training plateaus
- Tracks consecutive weeks with minimal progress
- Records weeks since last significant progress (≥1%)
- Provides last recorded progress rate for context

#### 27. Strength Gain Velocity (SGV)
**Formula**: `SGV = (Current 1RM - Initial 1RM) / Training Weeks`
- Rate of strength improvement over time
- Measured in weight units per week (kg/week, lbs/week)
- Calculated from first recorded 1RM to current best
- Indicates overall training effectiveness

#### 28. Adaptation Rate
**Formula**: `Adaptation Rate = Δ Performance / Δ Volume`
- Efficiency of strength gains relative to volume changes
- Higher values indicate better adaptation to training stimulus
- Calculated as change in 1RM divided by change in exercise volume
- Helps optimize volume-performance relationships

## When Metrics Are Calculated

### Automatic Calculation Triggers

1. **Workout Session Completion**
   - Triggered when `UpdateWorkoutSession` is called with `FinishedAt` timestamp
   - Location: `workout_session_service.go:UpdateWorkoutSession()`
   - Process:
     ```go
     if req.FinishedAt != nil {
         // Calculate workout metrics
         workoutMetrics, err := s.metricsService.CalculateWorkoutMetrics(ctx, &updatedSession)
         // Save workout metrics
         err = s.metricsService.SaveWorkoutMetrics(ctx, workoutMetrics)
         // Update user metrics
         _ = s.metricsService.UpdateUserMetrics(ctx, updatedSession.UserID)
     }
     ```

2. **User Metrics Refresh**
   - Recalculated on-demand when requested via API
   - Aggregates data from `workout_metrics` collection
   - Updates weekly, monthly, and all-time metrics

### Manual Calculation

Developers can manually trigger metrics calculation:

```go
// Calculate metrics for a completed session
workoutMetrics, err := metricsService.CalculateWorkoutMetrics(ctx, session)

// Save to database
err = metricsService.SaveWorkoutMetrics(ctx, workoutMetrics)

// Update user's aggregate metrics
err = metricsService.UpdateUserMetrics(ctx, userID)
```

## Data Models

### Core Metrics Models

```go
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
    AverageIntensity      float64            `bson:"averageIntensity" json:"averageIntensity"`
    RelativeIntensity     float64            `bson:"relativeIntensity" json:"relativeIntensity"`
    RPEAdjustedLoad       float64            `bson:"rpeAdjustedLoad" json:"rpeAdjustedLoad"`
    IntensityDistribution map[string]float64 `bson:"intensityDistribution" json:"intensityDistribution"`
    LoadDensity           float64            `bson:"loadDensity" json:"loadDensity"`
}

type StrengthMetrics struct {
    EstimatedOneRMEpley   map[string]float64 `bson:"estimatedOneRmEpley" json:"estimatedOneRmEpley"`
    EstimatedOneRMBrzycki map[string]float64 `bson:"estimatedOneRmBrzycki" json:"estimatedOneRmBrzycki"`
    WilksScore            float64            `bson:"wilksScore" json:"wilksScore"`
    PushPullRatio         float64            `bson:"pushPullRatio" json:"pushPullRatio"`
    PowerOutput           float64            `bson:"powerOutput" json:"powerOutput"`
}

type ProgressAdaptationMetrics struct {
    ProgressiveOverloadIndex float64                    `bson:"progressiveOverloadIndex" json:"progressiveOverloadIndex"`
    WeekOverWeekProgressRate map[string]float64         `bson:"weekOverWeekProgressRate" json:"weekOverWeekProgressRate"`
    PlateauDetection         map[string]PlateauStatus   `bson:"plateauDetection" json:"plateauDetection"`
    StrengthGainVelocity     map[string]float64         `bson:"strengthGainVelocity" json:"strengthGainVelocity"`
    AdaptationRate           map[string]float64         `bson:"adaptationRate" json:"adaptationRate"`
}

type PlateauStatus struct {
    IsPlateaued        bool    `bson:"isPlateaued" json:"isPlateaued"`
    ConsecutiveWeeks   int32   `bson:"consecutiveWeeks" json:"consecutiveWeeks"`
    LastProgressRate   float64 `bson:"lastProgressRate" json:"lastProgressRate"`
    WeeksSinceProgress int32   `bson:"weeksSinceProgress" json:"weeksSinceProgress"`
}

type MuscleSpecificMetrics struct {
    MuscleGroupDistribution map[string]float64 `bson:"muscleGroupDistribution" json:"muscleGroupDistribution"` // MG Distribution% = (MG Volume / Total Volume) × 100
    MuscleImbalanceIndex    map[string]float64 `bson:"muscleImbalanceIndex" json:"muscleImbalanceIndex"`       // Imbalance = |Left Side Strength - Right Side Strength| / Average Strength × 100
    AntagonistRatio         map[string]float64 `bson:"antagonistRatio" json:"antagonistRatio"`                 // Antagonist Ratio = Antagonist Strength / Agonist Strength
    StimulusToFatigueRatio  map[string]float64 `bson:"stimulusToFatigueRatio" json:"stimulusToFatigueRatio"`   // SFR = (Performance Gain / Baseline) / (RPE × Volume)
}

type WorkCapacityMetrics struct {
    TotalWorkCapacity      float64 `bson:"totalWorkCapacity" json:"totalWorkCapacity"`           // TWC = Σ(Sets × Reps × Weight × (1 - Rest Time/300))
    DensityTrainingIndex   float64 `bson:"densityTrainingIndex" json:"densityTrainingIndex"`     // Density = Total Volume / Total Time
    DensityProgressPercent float64 `bson:"densityProgressPercent" json:"densityProgressPercent"` // Progress% = (Current Density - Previous Density) / Previous Density × 100
    TimeUnderTension       float64 `bson:"timeUnderTension" json:"timeUnderTension"`             // TUT = Σ(Reps × Tempo in seconds)
    MechanicalTensionScore float64 `bson:"mechanicalTensionScore" json:"mechanicalTensionScore"` // MTS = Weight × TUT × (RPE/10)
}

type TrainingPatternMetrics struct {
    OptimalFrequency        float64 `bson:"optimalFrequency" json:"optimalFrequency"`               // Optimal Frequency = Recovery Time + 24-48 hours
    RecoveryTime            float64 `bson:"recoveryTime" json:"recoveryTime"`                       // Recovery Time = 24 × (Volume/10) × (Intensity/70) × (RPE/7)
    ExerciseSelectionIndex  float64 `bson:"exerciseSelectionIndex" json:"exerciseSelectionIndex"`   // Diversity Index = Unique Exercises / Total Exercises × 100
    ConsistencyScore        float64 `bson:"consistencyScore" json:"consistencyScore"`               // Consistency = (Actual Workouts / Planned Workouts) × 100
    WorkoutCompletionRate   float64 `bson:"workoutCompletionRate" json:"workoutCompletionRate"`     // Completion Rate = (Completed Sets / Planned Sets) × 100
}

type PeriodizationMetrics struct {
    ChronicTrainingLoad    float64 `bson:"chronicTrainingLoad" json:"chronicTrainingLoad"`       // CTL = Exponentially Weighted Average of daily TSS over 42 days
    AcuteTrainingLoad      float64 `bson:"acuteTrainingLoad" json:"acuteTrainingLoad"`           // ATL = Exponentially Weighted Average of daily TSS over 7 days
    TrainingStressBalance  float64 `bson:"trainingStressBalance" json:"trainingStressBalance"`   // TSB = CTL - ATL
    FormFreshnessIndex     float64 `bson:"formFreshnessIndex" json:"formFreshnessIndex"`         // FFI = TSB / CTL × 100
}

type InjuryRiskPreventionMetrics struct {
    InjuryRiskScore     float64 `bson:"injuryRiskScore" json:"injuryRiskScore"`         // IRS = (ACWR × Imbalance Index × Fatigue Score) / Recovery Quality, high risk if > 1.5
    LoadSpikeAlert      bool    `bson:"loadSpikeAlert" json:"loadSpikeAlert"`           // Spike = True if Weekly Volume > 1.5 × Average of Last 4 Weeks
    AsymmetryDevelopment float64 `bson:"asymmetryDevelopment" json:"asymmetryDevelopment"` // Asymmetry = |(Left Performance - Right Performance)| / Average × 100
}

type EfficiencyTechniqueMetrics struct {
	StrengthEfficiency        float64 `bson:"strengthEfficiency" json:"strengthEfficiency"`               // SE = (1RM Gain / Total Volume) × 1000
	VolumeEfficiency          float64 `bson:"volumeEfficiency" json:"volumeEfficiency"`                   // VE = Performance Improvement / Total Training Volume
	RPEPerformanceCorrelation float64 `bson:"rpePerformanceCorrelation" json:"rpePerformanceCorrelation"` // Correlation = Pearson(Actual Reps, 10 - RPE + Expected Reps at RPE)
	TechniqueConsistency      float64 `bson:"techniqueConsistency" json:"techniqueConsistency"`           // Consistency = 1 - (Standard Deviation of Rep Times / Average Rep Time)
}

type CompoundMetrics struct {
	FitnessAge                 float64 `bson:"fitnessAge" json:"fitnessAge"`                                 // Fitness Age = Chronological Age × (Population Average Score / User Score)
	OverallFitnessScore        float64 `bson:"overallFitnessScore" json:"overallFitnessScore"`               // OFS = (Strength Score × 0.3 + Volume Score × 0.2 + Consistency × 0.2 + Progress × 0.3) × 100
	TrainingEfficiencyQuotient float64 `bson:"trainingEfficiencyQuotient" json:"trainingEfficiencyQuotient"` // TEQ = (Performance Gains / (Time Invested × Fatigue Generated)) × 100
}

type PredictiveAnalyticsMetrics struct {
	PlateauProbability          float64            `bson:"plateauProbability" json:"plateauProbability"`                   // Plateau Probability = 1 / (1 + e^(Progress Rate - 0.02))
	PerformanceTrajectory       map[string]float64 `bson:"performanceTrajectory" json:"performanceTrajectory"`             // Exercise ID -> Future Performance = Current + (Weekly Gain Rate × Weeks × Decay Factor^Weeks)
	GoalAchievementTimeline     map[string]float64 `bson:"goalAchievementTimeline" json:"goalAchievementTimeline"`         // Exercise ID -> Weeks to Goal = (Goal - Current) / (Average Weekly Progress × Expected Decay)
	DetrainingRisk              float64            `bson:"detrainingRisk" json:"detrainingRisk"`                           // Detraining Risk = Days Since Last Workout / (CTL / 10)
	ProjectedPerformance4Weeks  map[string]float64 `bson:"projectedPerformance4Weeks" json:"projectedPerformance4Weeks"`   // Exercise ID -> 4-week performance projection
	ProjectedPerformance12Weeks map[string]float64 `bson:"projectedPerformance12Weeks" json:"projectedPerformance12Weeks"` // Exercise ID -> 12-week performance projection
	ProjectedPerformance26Weeks map[string]float64 `bson:"projectedPerformance26Weeks" json:"projectedPerformance26Weeks"` // Exercise ID -> 26-week performance projection
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
```

### Volume Landmarks

```go
type VolumeLandmarks struct {
    MEV map[string]float64 `bson:"mev" json:"mev"` // Minimum Effective Volume
    MAV map[string]float64 `bson:"mav" json:"mav"` // Maximum Adaptive Volume  
    MRV map[string]float64 `bson:"mrv" json:"mrv"` // Maximum Recoverable Volume
}
```

### User Metrics Aggregation

```go
type UserMetrics struct {
    ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID           primitive.ObjectID `bson:"userId" json:"userId"`
    WeeklyMetrics    PeriodMetrics      `bson:"weeklyMetrics" json:"weeklyMetrics"`
    MonthlyMetrics   PeriodMetrics      `bson:"monthlyMetrics" json:"monthlyMetrics"`
    AllTimeMetrics   PeriodMetrics      `bson:"allTimeMetrics" json:"allTimeMetrics"`
    VolumeLandmarks  VolumeLandmarks    `bson:"volumeLandmarks" json:"volumeLandmarks"`
    CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
    UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}
```

## API Endpoints

### REST API

#### Get User Metrics
```http
GET /api/users/{userId}/metrics?period={weekly|monthly|all}
```
- Returns aggregated metrics for the specified period
- Default period: "all"
- Includes volume landmarks (MEV/MAV/MRV)

#### Get Volume Trends  
```http
GET /api/users/{userId}/metrics/trends?period={weekly|monthly}
```
- Returns historical volume progression
- Default period: "weekly" 
- Includes 12 weeks or 12 months of data

#### Get Workout Metrics
```http
GET /api/workout-sessions/{sessionId}/metrics
```
- Returns detailed metrics for a specific workout session
- Includes per-exercise breakdowns

### gRPC Methods

```protobuf
service MetricsService {
  rpc GetUserMetrics(GetUserMetricsRequest) returns (UserMetrics);
  rpc GetWorkoutMetrics(GetWorkoutMetricsRequest) returns (WorkoutMetrics);
  rpc GetVolumeTrends(GetVolumeTrendsRequest) returns (TrendMetrics);
}
```

## Database Collections

### `workout_metrics`
- Stores individual workout session metrics
- Indexed by `userId`, `sessionId`, `date`
- Used for aggregating user metrics and trends

### `user_metrics`
- Stores aggregated user metrics (cached)
- Updated after each completed workout
- Includes weekly, monthly, and all-time aggregations

## Implementation Details

### Service Layer

**Location**: `backend/internal/services/metrics_service.go`

**Key Methods**:
- `CalculateWorkoutMetrics()` - Main calculation entry point
- `SaveWorkoutMetrics()` - Persists workout metrics
- `UpdateUserMetrics()` - Refreshes user aggregations
- `CalculateVolumeLandmarks()` - Computes MEV/MAV/MRV
- `calculateTrainingPatternMetrics()` - Computes training pattern analytics
- `calculatePeriodizationMetrics()` - Computes advanced periodization metrics
- `calculateExponentiallyWeightedAverage()` - Calculates exponentially weighted TSS averages

### Integration Points

1. **WorkoutSessionService**
   - Automatically triggers metrics calculation on session completion
   - Location: `workout_session_service.go:UpdateWorkoutSession()`

2. **HTTP Server**
   - Exposes REST endpoints for metrics access
   - Location: `pkg/http/server.go`

3. **Main Server**
   - Initializes MetricsService with dependencies
   - Location: `cmd/server/main.go`

### Configuration

**Muscle Group Weights**:
```go
var DefaultMuscleGroupWeights = MuscleGroupConstants{
    PrimaryMultiplier:   1.0,  // 100% weight for primary muscles
    SecondaryMultiplier: 0.6,  // 60% weight for secondary muscles
}
```

**Push/Pull Muscle Classifications**:
```go
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
```

**Wilks Coefficients**:
```go
var WilksMaleCoefficients = WilksCoefficients{
    A: -216.0475144, B: 16.2606339, C: -0.002388645,
    D: -0.00113732, E: 7.01863e-06, F: -1.291e-08,
}

var WilksFemaleCoefficients = WilksCoefficients{
    A: 594.31747775582, B: -27.23842536447, C: 0.82112226871,
    D: -0.00930733913, E: 4.731582e-05, F: -9.054e-08,
}
```

**RPE Thresholds**:
- Effective/Hard sets: RPE ≥ 7
- Configurable in calculation logic

**1RM Estimation**:
- Uses Epley formula: `1RM = Weight × (1 + Reps/30)`
- Uses Brzycki formula: `1RM = Weight × (36 / (37 - Reps))`
- Applied for intensity calculations and zone classifications
- Brzycki falls back to Epley for reps ≥ 37

**Volume Landmark Periods**:
- Based on 8 weeks of historical data
- Fallback to sensible defaults when insufficient data

**Strength Metrics Configuration**:
- Body weight defaults to 75kg (TODO: integrate with user profile)
- Gender defaults to male (TODO: integrate with user profile)
- Power output uses simplified 1m distance per rep
- Push/pull classification based on primary muscles only

**Progress & Adaptation Configuration**:
- Progressive Overload Index calculated weekly (current vs previous week)
- Plateau detection threshold: <1% progress for 3+ consecutive weeks
- Strength Gain Velocity measured from first recorded 1RM to current
- Adaptation Rate calculated weekly (performance change / volume change)
- All calculations are per-exercise and resilient to missing data

**Muscle-Specific Metrics Configuration**:
```go
// Antagonist muscle group pairs for ratio calculations
var AntagonistPairs = map[string]string{
    "quadriceps": "hamstrings",
    "hamstrings": "quadriceps",
    "chest":      "middle back",
    "middle back": "chest",
    "biceps":     "triceps",
    "triceps":    "biceps",
    "anterior deltoids": "posterior deltoids",
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
```

**Work Capacity Metrics Configuration**:
```go
// Default tempo values for exercises (in seconds)
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
```

**Training Pattern Analytics Configuration**:
```go
// Constants for Training Pattern Analytics
const (
    MinimumRecoveryHours = 24.0
    MaximumRecoveryHours = 48.0
    BaseVolumeThreshold  = 10.0
    BaseIntensityPercent = 70.0
    BaseRPEThreshold     = 7.0
)
```

**Periodization Metrics Configuration**:
```go
// Constants for Periodization Metrics 
const (
    ChronicTrainingLoadDays = 42 // 6 weeks
    AcuteTrainingLoadDays   = 7  // 1 week
    ExponentialDecayFactor  = 0.07 // Standard decay factor for training load
)
```

### Comparative & Normative Metrics

#### 47. Percentile Ranking
**Formula**: `Percentile = (Number of people below score / Total people) × 100`
- User's strength performance compared to similar demographics
- Adjusted for age, weight, and gender categories
- Based on Wilks score for overall strength comparison
- Provides context for training progress relative to population

#### 48. Training Age-Adjusted Expectations
**Formula**: `Expected Strength = Base Strength × (1 + 0.1 × sqrt(Training Years))`
- Strength expectations based on training experience
- Accounts for diminishing returns as training age increases
- Per-exercise predictions using body weight multipliers
- Helps set realistic training goals based on experience level

#### 49. Genetic Potential Estimate
**Formula**: `Potential = Current Max × (1.5 - 0.5 × (Current/Elite Standard))`
- Rough estimate of genetic strength potential
- Based on distance from elite performance standards
- Per-exercise potential calculations using body weight ratios
- Provides long-term training targets and expectations

### Psychological & Behavioral Metrics

#### 50. RPE Accuracy
**Formula**: `RPE Accuracy = 1 - |Predicted Reps at RPE - Actual Reps| / Actual Reps`
- Measures how well perceived exertion matches actual performance
- Higher accuracy indicates better self-awareness and pacing
- Helps calibrate effort perception for optimal training
- Improved accuracy leads to better autoregulation

#### 51. Motivation Index
**Formula**: `MI = (Consistency × Voluntary Extra Sets × (10 - Avg RPE)) / 100`
- Measures training motivation through behavioral indicators
- Considers workout consistency and voluntary effort
- Accounts for willingness to push through challenging training
- Values between 0-1, with higher indicating greater motivation

#### 52. Burnout Risk Score
**Formula**: `BRS = (Decreasing Performance + Increasing RPE + Decreased Frequency) / 3`
- Early warning system for training burnout
- Combines performance, effort, and frequency trends
- Values >0.7 indicate elevated burnout risk
- Helps prevent overtraining and maintain long-term progress

### Time-Based Analytics Metrics

#### 53. Optimal Training Time
**Formula**: `Performance Score by Hour = Average (1RM% × Completion Rate × (10 - RPE))`
- Identifies the hour of day with best training performance
- Based on historical performance across different times
- Helps optimize training schedule for peak performance
- Accounts for strength, completion, and subjective factors

#### 54. Optimal Rest Periods
**Formula**: `Optimal Rest = 2^(Intensity%/25) minutes`
- Exercise-specific rest period recommendations
- Based on training intensity and exercise demands
- Ranges from 1-10 minutes with intensity-based scaling
- Helps optimize recovery between sets for performance

#### 55. Session Duration Efficiency
**Formula**: `SDE = Total Effective Volume / Session Duration`
- Measures training efficiency in volume per minute
- Higher values indicate more productive time usage
- Helps identify optimal workout duration and pacing
- Useful for time-constrained training situations

### Compound Metrics

#### 56. Fitness Age
**Formula**: `Fitness Age = Chronological Age × (Population Average Score / User Score)`
- Biological fitness age compared to chronological age
- Based on overall fitness score relative to population average (75.0 baseline)
- Lower fitness age indicates better-than-average fitness for chronological age
- Higher fitness age suggests room for fitness improvement

#### 57. Overall Fitness Score (OFS)
**Formula**: `OFS = (Strength Score × 0.3 + Volume Score × 0.2 + Consistency × 0.2 + Progress × 0.3) × 100`
- Comprehensive fitness assessment combining multiple factors
- **Strength Score (30%)**: Based on Wilks score normalized to 0-100 scale
- **Volume Score (20%)**: Based on proximity to optimal volume landmarks (MAV)
- **Consistency Score (20%)**: Based on training adherence and frequency
- **Progress Score (30%)**: Based on week-over-week strength improvements
- Scale: 0-100, with higher scores indicating better overall fitness

#### 58. Training Efficiency Quotient (TEQ)
**Formula**: `TEQ = (Performance Gains / (Time Invested × Fatigue Generated)) × 100`
- Measures efficiency of training in terms of gains per unit effort
- **Performance Gains**: Average 1RM improvement across exercises
- **Time Invested**: Workout duration in hours
- **Fatigue Generated**: RPE normalized to 0-1 scale
- Higher values indicate more efficient training methods

### Predictive Analytics Metrics

#### 59. Plateau Probability
**Formula**: `Plateau Probability = 1 / (1 + e^(Progress Rate - 0.02))`
- Probability of hitting a training plateau based on current progress
- Uses sigmoid function with 2% progress rate threshold
- Values closer to 1.0 indicate higher plateau risk
- Based on average progress rate across all exercises

#### 60. Performance Trajectory
**Formula**: `Future Performance = Current + (Weekly Gain Rate × Weeks × Decay Factor^Weeks)`
- Projects future performance based on current trends
- **Decay Factor**: 0.95 (95% retention of gains over time)
- Calculated for multiple time horizons: 4, 12, and 26 weeks
- Per-exercise projections based on historical 1RM improvements

#### 61. Goal Achievement Timeline
**Formula**: `Weeks to Goal = (Goal - Current) / (Average Weekly Progress × Expected Decay)`
- Estimates time needed to reach performance goals
- Assumes 10% improvement goal by default
- Accounts for diminishing returns through decay factor
- Returns 999 weeks for exercises with no measurable progress

#### 62. Detraining Risk
**Formula**: `Detraining Risk = Days Since Last Workout / (CTL / 10)`
- Risk of fitness loss due to training breaks
- Based on time since last session and chronic training load
- Values >1.0 indicate elevated detraining risk
- Helps identify when to resume training after breaks

#### 63. Projected Performance (Multiple Horizons)
**Formulas**: 
- **4-week projection**: `Current + (Weekly Gain × 4 × 0.95^4)`
- **12-week projection**: `Current + (Weekly Gain × 12 × 0.95^12)`
- **26-week projection**: `Current + (Weekly Gain × 26 × 0.95^26)`
- Short, medium, and long-term performance forecasts
- Incorporates performance decay over time
- Per-exercise projections for goal setting and planning

## Usage Examples

### Frontend Integration

```typescript
// Get user's current metrics
const response = await fetch(`/api/users/${userId}/metrics?period=monthly`);
const metrics = await response.json();

// Access volume landmarks for training recommendations
const { mev, mav, mrv } = metrics.volumeLandmarks;
console.log(`Chest MEV: ${mev.chest}, MAV: ${mav.chest}, MRV: ${mrv.chest}`);

// Get volume trends for progress tracking
const trends = await fetch(`/api/users/${userId}/metrics/trends?period=weekly`);
const trendData = await trends.json();
```

### Backend Usage

```go
// Calculate metrics after workout completion
session := &models.WorkoutSession{
    // ... session data with FinishedAt set
}

workoutMetrics, err := metricsService.CalculateWorkoutMetrics(ctx, session)
if err != nil {
    return fmt.Errorf("failed to calculate metrics: %w", err)
}

// Access calculated metrics
fmt.Printf("Total Volume: %.2f", workoutMetrics.VolumeMetrics.TotalVolumeLoad)
fmt.Printf("Hard Sets: %d", workoutMetrics.VolumeMetrics.HardSets)

// Access intensity metrics
fmt.Printf("Average Intensity: %.1f%%", workoutMetrics.IntensityMetrics.AverageIntensity)
fmt.Printf("Load Density: %.2f", workoutMetrics.IntensityMetrics.LoadDensity)
fmt.Printf("Light Sets: %.1f%%", workoutMetrics.IntensityMetrics.IntensityDistribution["light"])

// Access strength metrics
fmt.Printf("Wilks Score: %.2f", workoutMetrics.StrengthMetrics.WilksScore)
fmt.Printf("Push:Pull Ratio: %.2f", workoutMetrics.StrengthMetrics.PushPullRatio)
fmt.Printf("Power Output: %.2f", workoutMetrics.StrengthMetrics.PowerOutput)

// Access 1RM estimates per exercise
for exerciseID, oneRM := range workoutMetrics.StrengthMetrics.EstimatedOneRMEpley {
    fmt.Printf("Exercise %s - Estimated 1RM (Epley): %.2f kg", exerciseID, oneRM)
}

for exerciseID, oneRM := range workoutMetrics.StrengthMetrics.EstimatedOneRMBrzycki {
    fmt.Printf("Exercise %s - Estimated 1RM (Brzycki): %.2f kg", exerciseID, oneRM)
}

// Access Progress & Adaptation metrics
fmt.Printf("Progressive Overload Index: %.2f", workoutMetrics.ProgressAdaptationMetrics.ProgressiveOverloadIndex)

// Access per-exercise progress metrics
for exerciseID, progressRate := range workoutMetrics.ProgressAdaptationMetrics.WeekOverWeekProgressRate {
    fmt.Printf("Exercise %s - Progress Rate: %.2f%%", exerciseID, progressRate)
}

// Access plateau detection
for exerciseID, plateau := range workoutMetrics.ProgressAdaptationMetrics.PlateauDetection {
    fmt.Printf("Exercise %s - Plateaued: %t (Consecutive weeks: %d)", 
        exerciseID, plateau.IsPlateaued, plateau.ConsecutiveWeeks)
}

// Access strength gain velocity
for exerciseID, sgv := range workoutMetrics.ProgressAdaptationMetrics.StrengthGainVelocity {
    fmt.Printf("Exercise %s - Strength Gain Velocity: %.2f kg/week", exerciseID, sgv)
}

// Access adaptation rate
for exerciseID, adaptationRate := range workoutMetrics.ProgressAdaptationMetrics.AdaptationRate {
    fmt.Printf("Exercise %s - Adaptation Rate: %.2f", exerciseID, adaptationRate)
}

// Access Muscle-Specific metrics
fmt.Printf("Muscle Group Distribution:")
for muscleGroup, percentage := range workoutMetrics.MuscleSpecificMetrics.MuscleGroupDistribution {
    fmt.Printf("  %s: %.1f%%", muscleGroup, percentage)
}

fmt.Printf("Muscle Imbalance Indices:")
for muscleGroup, imbalance := range workoutMetrics.MuscleSpecificMetrics.MuscleImbalanceIndex {
    fmt.Printf("  %s: %.1f%%", muscleGroup, imbalance)
}

fmt.Printf("Antagonist Ratios:")
for pairing, ratio := range workoutMetrics.MuscleSpecificMetrics.AntagonistRatio {
    fmt.Printf("  %s: %.2f", pairing, ratio)
}

fmt.Printf("Stimulus-to-Fatigue Ratios:")
for muscleGroup, sfr := range workoutMetrics.MuscleSpecificMetrics.StimulusToFatigueRatio {
    fmt.Printf("  %s: %.3f", muscleGroup, sfr)
}

// Access Work Capacity metrics
fmt.Printf("Total Work Capacity: %.2f", workoutMetrics.WorkCapacityMetrics.TotalWorkCapacity)
fmt.Printf("Density Training Index: %.2f volume/min", workoutMetrics.WorkCapacityMetrics.DensityTrainingIndex)
fmt.Printf("Density Progress: %.1f%%", workoutMetrics.WorkCapacityMetrics.DensityProgressPercent)
fmt.Printf("Time Under Tension: %.1f seconds", workoutMetrics.WorkCapacityMetrics.TimeUnderTension)
fmt.Printf("Mechanical Tension Score: %.2f", workoutMetrics.WorkCapacityMetrics.MechanicalTensionScore)

// Access Training Pattern Analytics metrics
fmt.Printf("Optimal Training Frequency: %.1f hours", workoutMetrics.TrainingPatternMetrics.OptimalFrequency)
fmt.Printf("Recovery Time: %.1f hours", workoutMetrics.TrainingPatternMetrics.RecoveryTime)
fmt.Printf("Exercise Selection Diversity: %.1f%%", workoutMetrics.TrainingPatternMetrics.ExerciseSelectionIndex)
fmt.Printf("Consistency Score: %.1f%%", workoutMetrics.TrainingPatternMetrics.ConsistencyScore)
fmt.Printf("Workout Completion Rate: %.1f%%", workoutMetrics.TrainingPatternMetrics.WorkoutCompletionRate)

// Access Advanced Periodization metrics
fmt.Printf("Chronic Training Load (CTL): %.2f", workoutMetrics.PeriodizationMetrics.ChronicTrainingLoad)
fmt.Printf("Acute Training Load (ATL): %.2f", workoutMetrics.PeriodizationMetrics.AcuteTrainingLoad)
fmt.Printf("Training Stress Balance (TSB): %.2f", workoutMetrics.PeriodizationMetrics.TrainingStressBalance)
fmt.Printf("Form/Freshness Index (FFI): %.1f%%", workoutMetrics.PeriodizationMetrics.FormFreshnessIndex)

// Interpret periodization metrics
if workoutMetrics.PeriodizationMetrics.TrainingStressBalance > 0 {
    fmt.Printf("Status: Fresh and ready for hard training")
} else if workoutMetrics.PeriodizationMetrics.TrainingStressBalance < -10 {
    fmt.Printf("Status: Fatigued - consider recovery or deload")
} else {
    fmt.Printf("Status: Balanced training stress")
}

if workoutMetrics.PeriodizationMetrics.FormFreshnessIndex >= -5 && workoutMetrics.PeriodizationMetrics.FormFreshnessIndex <= 15 {
    fmt.Printf("Form Status: Optimal for peak performance")
} else {
    fmt.Printf("Form Status: Outside optimal range - adjust training load")
}

// Access Injury Risk & Prevention metrics
fmt.Printf("Injury Risk Score: %.2f", workoutMetrics.InjuryRiskPreventionMetrics.InjuryRiskScore)
fmt.Printf("Load Spike Alert: %t", workoutMetrics.InjuryRiskPreventionMetrics.LoadSpikeAlert)
fmt.Printf("Asymmetry Development: %.1f%%", workoutMetrics.InjuryRiskPreventionMetrics.AsymmetryDevelopment)

// Interpret injury risk
if workoutMetrics.InjuryRiskPreventionMetrics.InjuryRiskScore > models.InjuryRiskThreshold {
    fmt.Printf("⚠️ High injury risk detected - consider reducing training load")
}

if workoutMetrics.InjuryRiskPreventionMetrics.LoadSpikeAlert {
    fmt.Printf("⚠️ Training volume spike detected - monitor recovery closely")
}

if workoutMetrics.InjuryRiskPreventionMetrics.AsymmetryDevelopment > models.AsymmetryRiskThreshold {
    fmt.Printf("⚠️ High asymmetry detected - focus on corrective exercises")
}

// Access Efficiency & Technique metrics
fmt.Printf("Strength Efficiency: %.2f", workoutMetrics.EfficiencyTechniqueMetrics.StrengthEfficiency)
fmt.Printf("Volume Efficiency: %.6f", workoutMetrics.EfficiencyTechniqueMetrics.VolumeEfficiency)
fmt.Printf("RPE-Performance Correlation: %.3f", workoutMetrics.EfficiencyTechniqueMetrics.RPEPerformanceCorrelation)
fmt.Printf("Technique Consistency: %.3f", workoutMetrics.EfficiencyTechniqueMetrics.TechniqueConsistency)

// Interpret efficiency metrics
if workoutMetrics.EfficiencyTechniqueMetrics.StrengthEfficiency > 2.0 {
    fmt.Printf("✅ High strength efficiency - optimal strength gains per volume")
} else if workoutMetrics.EfficiencyTechniqueMetrics.StrengthEfficiency < 0.5 {
    fmt.Printf("📊 Low strength efficiency - consider program adjustments")
}

if workoutMetrics.EfficiencyTechniqueMetrics.RPEPerformanceCorrelation > 0.7 {
    fmt.Printf("✅ Good RPE calibration - accurate effort perception")
} else if workoutMetrics.EfficiencyTechniqueMetrics.RPEPerformanceCorrelation < 0.3 {
    fmt.Printf("📊 Poor RPE calibration - work on effort perception")
}

if workoutMetrics.EfficiencyTechniqueMetrics.TechniqueConsistency > 0.8 {
    fmt.Printf("✅ Excellent technique consistency")
} else if workoutMetrics.EfficiencyTechniqueMetrics.TechniqueConsistency < 0.6 {
    fmt.Printf("📊 Technique inconsistency detected - focus on movement quality")
}

// Access Comparative & Normative metrics
fmt.Printf("Percentile Ranking: %.1f%%", workoutMetrics.ComparativeNormativeMetrics.PercentileRanking)
fmt.Printf("Performance Category: %s", workoutMetrics.ComparativeNormativeMetrics.PerformanceCategory)
fmt.Printf("Relative to Expectations: %.2f", workoutMetrics.ComparativeNormativeMetrics.RelativeToExpectations)

// Access training age-adjusted expectations
fmt.Printf("Training Age-Adjusted Expected Strength:")
for exerciseID, expected := range workoutMetrics.ComparativeNormativeMetrics.TrainingAgeAdjustedExpected {
    fmt.Printf("  Exercise %s: %.2f kg expected", exerciseID, expected)
}

// Access genetic potential estimates
fmt.Printf("Genetic Potential Estimates:")
for exerciseID, potential := range workoutMetrics.ComparativeNormativeMetrics.GeneticPotentialEstimate {
    fmt.Printf("  Exercise %s: %.2f kg potential", exerciseID, potential)
}

// Interpret comparative metrics
if workoutMetrics.ComparativeNormativeMetrics.PercentileRanking >= 85.0 {
    fmt.Printf("🏆 Elite performance level - top 15%% of population")
} else if workoutMetrics.ComparativeNormativeMetrics.PercentileRanking >= 60.0 {
    fmt.Printf("💪 Advanced performance level")
} else if workoutMetrics.ComparativeNormativeMetrics.PercentileRanking >= 25.0 {
    fmt.Printf("📈 Intermediate performance level")
} else {
    fmt.Printf("🚀 Beginner level - lots of room for growth")
}

if workoutMetrics.ComparativeNormativeMetrics.RelativeToExpectations > 1.2 {
    fmt.Printf("⭐ Exceeding expectations for training age")
} else if workoutMetrics.ComparativeNormativeMetrics.RelativeToExpectations < 0.8 {
    fmt.Printf("📊 Below expectations - review training approach")
}

// Access Psychological & Behavioral metrics
fmt.Printf("RPE Accuracy: %.3f", workoutMetrics.PsychologicalBehavioralMetrics.RPEAccuracy)
fmt.Printf("Motivation Index: %.3f", workoutMetrics.PsychologicalBehavioralMetrics.MotivationIndex)
fmt.Printf("Burnout Risk Score: %.3f", workoutMetrics.PsychologicalBehavioralMetrics.BurnoutRiskScore)
fmt.Printf("Training Adherence: %.1f%%", workoutMetrics.PsychologicalBehavioralMetrics.TrainingAdherence)
fmt.Printf("Consistency Trend: %.3f", workoutMetrics.PsychologicalBehavioralMetrics.ConsistencyTrend)

// Interpret psychological metrics
if workoutMetrics.PsychologicalBehavioralMetrics.RPEAccuracy > 0.8 {
    fmt.Printf("✅ Excellent RPE accuracy - great body awareness")
} else if workoutMetrics.PsychologicalBehavioralMetrics.RPEAccuracy < 0.5 {
    fmt.Printf("📊 Poor RPE accuracy - work on effort calibration")
}

if workoutMetrics.PsychologicalBehavioralMetrics.MotivationIndex > 0.7 {
    fmt.Printf("🔥 High motivation - great training drive")
} else if workoutMetrics.PsychologicalBehavioralMetrics.MotivationIndex < 0.3 {
    fmt.Printf("💡 Low motivation - consider program variety or goals review")
}

if workoutMetrics.PsychologicalBehavioralMetrics.BurnoutRiskScore > 0.7 {
    fmt.Printf("⚠️ High burnout risk - consider deload or rest")
} else if workoutMetrics.PsychologicalBehavioralMetrics.BurnoutRiskScore < 0.3 {
    fmt.Printf("✅ Low burnout risk - sustainable training pace")
}

if workoutMetrics.PsychologicalBehavioralMetrics.ConsistencyTrend > 0.1 {
    fmt.Printf("📈 Improving consistency - great progress")
} else if workoutMetrics.PsychologicalBehavioralMetrics.ConsistencyTrend < -0.1 {
    fmt.Printf("📉 Declining consistency - address barriers to training")
}

// Access Time-Based Analytics metrics
fmt.Printf("Optimal Training Time: %d:00", workoutMetrics.TimeBasedAnalyticsMetrics.OptimalTrainingTimeHour)
fmt.Printf("Time of Day Preference: %s", workoutMetrics.TimeBasedAnalyticsMetrics.TimeOfDayPreference)
fmt.Printf("Session Duration Efficiency: %.2f reps/min", workoutMetrics.TimeBasedAnalyticsMetrics.SessionDurationEfficiency)
fmt.Printf("Workout Timing Consistency: %.3f", workoutMetrics.TimeBasedAnalyticsMetrics.WorkoutTimingConsistency)

// Access performance by hour
fmt.Printf("Performance by Hour:")
for hour, performance := range workoutMetrics.TimeBasedAnalyticsMetrics.PerformanceByHour {
    fmt.Printf("  %s:00 - %.3f performance score", hour, performance)
}

// Access optimal rest periods
fmt.Printf("Optimal Rest Periods:")
for exerciseID, restTime := range workoutMetrics.TimeBasedAnalyticsMetrics.OptimalRestPeriods {
    fmt.Printf("  Exercise %s: %.1f minutes", exerciseID, restTime)
}

// Interpret time-based metrics
preferredTime := workoutMetrics.TimeBasedAnalyticsMetrics.OptimalTrainingTimeHour
if preferredTime >= 6 && preferredTime < 12 {
    fmt.Printf("🌅 Morning training shows best performance")
} else if preferredTime >= 12 && preferredTime < 18 {
    fmt.Printf("☀️ Afternoon training shows best performance")
} else {
    fmt.Printf("🌙 Evening training shows best performance")
}

if workoutMetrics.TimeBasedAnalyticsMetrics.SessionDurationEfficiency > 1.0 {
    fmt.Printf("⚡ High training efficiency - good volume per minute")
} else if workoutMetrics.TimeBasedAnalyticsMetrics.SessionDurationEfficiency < 0.5 {
    fmt.Printf("⏰ Low training efficiency - consider workout pacing")
}

if workoutMetrics.TimeBasedAnalyticsMetrics.WorkoutTimingConsistency > 0.8 {
    fmt.Printf("🕒 Consistent workout timing - excellent routine")
} else if workoutMetrics.TimeBasedAnalyticsMetrics.WorkoutTimingConsistency < 0.5 {
	fmt.Printf("📅 Inconsistent workout timing - try to establish routine")
}

// Access Compound Metrics
fmt.Printf("Fitness Age: %.1f years", workoutMetrics.CompoundMetrics.FitnessAge)
fmt.Printf("Overall Fitness Score: %.1f/100", workoutMetrics.CompoundMetrics.OverallFitnessScore)
fmt.Printf("Training Efficiency Quotient: %.1f", workoutMetrics.CompoundMetrics.TrainingEfficiencyQuotient)

// Interpret compound metrics
if workoutMetrics.CompoundMetrics.FitnessAge < 25.0 {
	fmt.Printf("🌟 Excellent fitness age - you're aging like fine wine!")
} else if workoutMetrics.CompoundMetrics.FitnessAge < 35.0 {
	fmt.Printf("💪 Good fitness age - above average fitness level")
} else {
	fmt.Printf("📈 Room for improvement - focus on consistency and progress")
}

if workoutMetrics.CompoundMetrics.OverallFitnessScore >= 80.0 {
	fmt.Printf("🏆 Excellent overall fitness - top tier performance")
} else if workoutMetrics.CompoundMetrics.OverallFitnessScore >= 60.0 {
	fmt.Printf("👍 Good overall fitness - solid foundation")
} else if workoutMetrics.CompoundMetrics.OverallFitnessScore >= 40.0 {
	fmt.Printf("📊 Average fitness level - potential for growth")
} else {
	fmt.Printf("🚀 Building fitness foundation - keep up the great work!")
}

if workoutMetrics.CompoundMetrics.TrainingEfficiencyQuotient >= 100.0 {
	fmt.Printf("⚡ Highly efficient training - great gains per effort")
} else if workoutMetrics.CompoundMetrics.TrainingEfficiencyQuotient >= 50.0 {
	fmt.Printf("✅ Efficient training approach")
} else {
	fmt.Printf("💡 Consider optimizing training efficiency")
}

// Access Predictive Analytics Metrics
fmt.Printf("Plateau Probability: %.1f%%", workoutMetrics.PredictiveAnalyticsMetrics.PlateauProbability*100)
fmt.Printf("Detraining Risk: %.2f", workoutMetrics.PredictiveAnalyticsMetrics.DetrainingRisk)

// Access performance projections
fmt.Printf("Performance Projections:")
for exerciseID, projection := range workoutMetrics.PredictiveAnalyticsMetrics.ProjectedPerformance4Weeks {
	fmt.Printf("  Exercise %s - 4 weeks: %.1f kg", exerciseID, projection)
}

for exerciseID, projection := range workoutMetrics.PredictiveAnalyticsMetrics.ProjectedPerformance12Weeks {
	fmt.Printf("  Exercise %s - 12 weeks: %.1f kg", exerciseID, projection)
}

for exerciseID, projection := range workoutMetrics.PredictiveAnalyticsMetrics.ProjectedPerformance26Weeks {
	fmt.Printf("  Exercise %s - 26 weeks: %.1f kg", exerciseID, projection)
}

// Access goal achievement timelines
fmt.Printf("Goal Achievement Timelines (10% improvement):")
for exerciseID, weeks := range workoutMetrics.PredictiveAnalyticsMetrics.GoalAchievementTimeline {
	if weeks < 999 {
		fmt.Printf("  Exercise %s: %.1f weeks to reach goal", exerciseID, weeks)
	} else {
		fmt.Printf("  Exercise %s: Goal may take very long - focus on consistency", exerciseID)
	}
}

// Interpret predictive metrics
if workoutMetrics.PredictiveAnalyticsMetrics.PlateauProbability > 0.7 {
	fmt.Printf("⚠️ High plateau risk - consider program variation or deload")
} else if workoutMetrics.PredictiveAnalyticsMetrics.PlateauProbability > 0.5 {
	fmt.Printf("📊 Moderate plateau risk - monitor progress closely")
} else {
	fmt.Printf("✅ Low plateau risk - continue current approach")
}

if workoutMetrics.PredictiveAnalyticsMetrics.DetrainingRisk > 1.0 {
	fmt.Printf("⚠️ High detraining risk - consider resuming training soon")
} else if workoutMetrics.PredictiveAnalyticsMetrics.DetrainingRisk > 0.5 {
	fmt.Printf("📅 Moderate detraining risk - maintain regular training")
} else {
	fmt.Printf("✅ Low detraining risk - good training consistency")
}
```

### Querying Historical Data

```