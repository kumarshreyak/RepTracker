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
```

### Querying Historical Data

```go
// Get user's volume trends
trends, err := metricsService.CalculateVolumeTrends(ctx, userID, "weekly")
if err != nil {
    return err
}

// Access progression data
for _, point := range trends.VolumeProgression {
    fmt.Printf("Date: %s, Volume: %.2f\n", 
        point.Date.Format("2006-01-02"), 
        point.VolumeLoad)
}
```

## Error Handling

The metrics system is designed to be resilient:

1. **Non-blocking**: Metrics calculation failures don't prevent workout session updates
2. **Graceful degradation**: Missing exercise data is skipped rather than failing entirely
3. **Default values**: Volume landmarks fall back to sensible defaults when insufficient data
4. **Logging**: All errors are logged for debugging without exposing to users

## Performance Considerations

1. **Async calculation**: Metrics are calculated after successful session updates
2. **Cached aggregations**: User metrics are cached in `user_metrics` collection
3. **Efficient queries**: Database queries use appropriate indexes
4. **Batch operations**: Multiple metrics calculated in single database operations

## Future Enhancements

1. **Body Weight Integration**: Relative volume calculations when body weight tracking is added
2. **Per-set RPE**: More granular RPE tracking for improved effective reps/hard sets
3. **Exercise-specific landmarks**: MEV/MAV/MRV per exercise rather than just muscle groups
4. **Advanced analytics**: Strength curves, volume-load relationships, deload recommendations
5. **Real-time updates**: WebSocket integration for live metrics updates during workouts
6. **Enhanced Tempo Tracking**: User-defined tempo patterns per exercise rather than defaults
7. **Rest Period Tracking**: Actual rest time measurement for more accurate work capacity calculations
8. **Bilateral Exercise Detection**: Automatic identification of unilateral vs bilateral exercises for imbalance calculations
9. **Advanced Muscle Mapping**: More detailed muscle group classifications and exercise-muscle relationships
10. **Fatigue Prediction**: ML-based fatigue prediction using stimulus-to-fatigue ratios and historical data

## Contributing

When modifying the metrics system:

1. **Maintain backwards compatibility** in data models
2. **Update this documentation** for any new metrics or changes
3. **Add tests** for new calculation logic
4. **Follow camelCase naming** for all new fields and properties
5. **Preserve existing API contracts** unless versioning appropriately 