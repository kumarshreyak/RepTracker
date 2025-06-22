package services

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gymlog-backend/pkg/database"
	"gymlog-backend/pkg/models"
	pb "gymlog-backend/proto/gymlog/v1"
)

type MetricsService struct {
	pb.UnimplementedMetricsServiceServer
	db              *database.MongoDB
	metricsColl     *mongo.Collection
	userMetricsColl *mongo.Collection
	exerciseService *ExerciseService
}

func NewMetricsService(db *database.MongoDB, exerciseService *ExerciseService) *MetricsService {
	return &MetricsService{
		db:              db,
		metricsColl:     db.GetCollection("workout_metrics"),
		userMetricsColl: db.GetCollection("user_metrics"),
		exerciseService: exerciseService,
	}
}

func (m *MetricsService) CalculateWorkoutMetrics(ctx context.Context, session *models.WorkoutSession) (*models.WorkoutMetrics, error) {
	if session.FinishedAt == nil {
		return nil, fmt.Errorf("workout session is not finished")
	}

	metrics := &models.WorkoutMetrics{
		UserID:    session.UserID,
		SessionID: session.ID,
		RoutineID: session.RoutineID,
		Date:      session.StartedAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	metrics.WorkoutDurationSecs = session.DurationSeconds

	var err error
	metrics.VolumeMetrics, err = m.calculateVolumeMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate volume metrics: %w", err)
	}

	metrics.PerformanceMetrics = m.calculatePerformanceMetrics(session)

	metrics.IntensityMetrics, err = m.calculateIntensityMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate intensity metrics: %w", err)
	}

	metrics.StrengthMetrics, err = m.calculateStrengthMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate strength metrics: %w", err)
	}

	metrics.ProgressAdaptationMetrics, err = m.calculateProgressAdaptationMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate progress adaptation metrics: %w", err)
	}

	metrics.SetMetrics = m.calculateSetMetrics(session)

	metrics.ExerciseMetrics, err = m.calculateExerciseMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate exercise metrics: %w", err)
	}

	metrics.RecoveryFatigueMetrics, err = m.calculateRecoveryFatigueMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate recovery fatigue metrics: %w", err)
	}

	metrics.BodyCompositionMetrics, err = m.calculateBodyCompositionMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate body composition metrics: %w", err)
	}

	metrics.MuscleSpecificMetrics, err = m.calculateMuscleSpecificMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate muscle specific metrics: %w", err)
	}

	metrics.WorkCapacityMetrics, err = m.calculateWorkCapacityMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate work capacity metrics: %w", err)
	}

	metrics.TrainingPatternMetrics, err = m.calculateTrainingPatternMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate training pattern metrics: %w", err)
	}

	metrics.PeriodizationMetrics, err = m.calculatePeriodizationMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate periodization metrics: %w", err)
	}

	metrics.InjuryRiskPreventionMetrics, err = m.calculateInjuryRiskPreventionMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate injury risk prevention metrics: %w", err)
	}

	metrics.EfficiencyTechniqueMetrics, err = m.calculateEfficiencyTechniqueMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate efficiency technique metrics: %w", err)
	}

	return metrics, nil
}

func (m *MetricsService) calculateVolumeMetrics(ctx context.Context, session *models.WorkoutSession) (models.VolumeMetrics, error) {
	var totalVolumeLoad, tonnage float64
	var effectiveReps, hardSets int32
	muscleGroupVolume := make(map[string]float64)

	for _, exercise := range session.Exercises {
		exerciseDetails, err := m.exerciseService.GetExercise(ctx, &pb.GetExerciseRequest{Id: exercise.ExerciseID.Hex()})
		if err != nil {
			continue
		}

		for _, set := range exercise.Sets {
			if !set.Completed {
				continue
			}

			setVolumeLoad := float64(set.ActualReps) * float64(set.ActualWeight)
			setTonnage := float64(set.ActualReps) * float64(set.ActualWeight)

			totalVolumeLoad += setVolumeLoad
			tonnage += setTonnage

			if session.RPERating >= 7 {
				effectiveReps += set.ActualReps
				hardSets++
			}

			for _, muscle := range exerciseDetails.PrimaryMuscles {
				muscleGroupVolume[muscle] += setVolumeLoad * models.DefaultMuscleGroupWeights.PrimaryMultiplier
			}

			for _, muscle := range exerciseDetails.SecondaryMuscles {
				muscleGroupVolume[muscle] += setVolumeLoad * models.DefaultMuscleGroupWeights.SecondaryMultiplier
			}
		}
	}

	return models.VolumeMetrics{
		TotalVolumeLoad:   totalVolumeLoad,
		Tonnage:           tonnage,
		RelativeVolume:    totalVolumeLoad, // TODO: Divide by body weight when available
		EffectiveReps:     effectiveReps,
		HardSets:          hardSets,
		MuscleGroupVolume: muscleGroupVolume,
	}, nil
}

func (m *MetricsService) calculatePerformanceMetrics(session *models.WorkoutSession) models.PerformanceMetrics {
	rpeDistribution := make(map[int32]int32)
	var totalRPE int32
	var setCount int32
	var maxRPE int32

	for _, exercise := range session.Exercises {
		for _, set := range exercise.Sets {
			if !set.Completed {
				continue
			}

			rpe := session.RPERating
			if rpe > 0 {
				rpeDistribution[rpe]++
				totalRPE += rpe
				setCount++
				if rpe > maxRPE {
					maxRPE = rpe
				}
			}
		}
	}

	var averageRPE float64
	if setCount > 0 {
		averageRPE = float64(totalRPE) / float64(setCount)
	}

	return models.PerformanceMetrics{
		AverageRPE:      averageRPE,
		MaxRPE:          maxRPE,
		RPEDistribution: rpeDistribution,
	}
}

func (m *MetricsService) calculateIntensityMetrics(ctx context.Context, session *models.WorkoutSession) (models.IntensityMetrics, error) {
	var totalWeightReps, totalOneRMReps float64
	var totalRPEAdjustedLoad float64
	var totalSets int32
	var lightSets, moderateSets, heavySets int32

	// Default body weight for relative intensity calculation (TODO: get from user profile)
	bodyWeight := 75.0 // kg - this should be retrieved from user profile when available

	for _, exercise := range session.Exercises {
		for _, set := range exercise.Sets {
			if !set.Completed {
				continue
			}

			weight := float64(set.ActualWeight)
			reps := float64(set.ActualReps)

			// Calculate 1RM using Epley formula: 1RM = weight × (1 + reps/30)
			oneRM := m.calculateOneRM(weight, int32(reps))

			// For average intensity calculation
			totalWeightReps += weight * reps
			totalOneRMReps += oneRM * reps

			// RPE-adjusted load calculation
			rpe := float64(session.RPERating)
			if rpe > 0 {
				rpeAdjustedLoad := weight * math.Pow(rpe/10.0, 2)
				totalRPEAdjustedLoad += rpeAdjustedLoad
			}

			// Intensity distribution - categorize based on % of 1RM
			if oneRM > 0 {
				intensityPercent := (weight / oneRM) * 100
				if intensityPercent < 60 {
					lightSets++
				} else if intensityPercent <= 80 {
					moderateSets++
				} else {
					heavySets++
				}
			}

			totalSets++
		}
	}

	// Calculate average intensity
	var averageIntensity float64
	if totalOneRMReps > 0 {
		averageIntensity = (totalWeightReps / totalOneRMReps) * 100
	}

	// Calculate relative intensity (using max weight lifted in session)
	var maxWeight float64
	for _, exercise := range session.Exercises {
		for _, set := range exercise.Sets {
			if set.Completed && float64(set.ActualWeight) > maxWeight {
				maxWeight = float64(set.ActualWeight)
			}
		}
	}
	relativeIntensity := maxWeight / bodyWeight

	// Calculate intensity distribution percentages
	intensityDistribution := make(map[string]float64)
	if totalSets > 0 {
		intensityDistribution["light"] = (float64(lightSets) / float64(totalSets)) * 100
		intensityDistribution["moderate"] = (float64(moderateSets) / float64(totalSets)) * 100
		intensityDistribution["heavy"] = (float64(heavySets) / float64(totalSets)) * 100
	}

	// Calculate load density (Total Volume / Total Workout Time in minutes)
	var loadDensity float64
	if session.DurationSeconds > 0 {
		workoutTimeMinutes := float64(session.DurationSeconds) / 60.0
		// Use total volume load from volume metrics calculation
		volumeMetrics, err := m.calculateVolumeMetrics(ctx, session)
		if err == nil {
			loadDensity = volumeMetrics.TotalVolumeLoad / workoutTimeMinutes
		}
	}

	return models.IntensityMetrics{
		AverageIntensity:      averageIntensity,
		RelativeIntensity:     relativeIntensity,
		RPEAdjustedLoad:       totalRPEAdjustedLoad,
		IntensityDistribution: intensityDistribution,
		LoadDensity:           loadDensity,
	}, nil
}

// calculateOneRM estimates 1RM using the Epley formula
func (m *MetricsService) calculateOneRM(weight float64, reps int32) float64 {
	if reps <= 0 {
		return 0
	}
	if reps == 1 {
		return weight
	}
	// Epley formula: 1RM = weight × (1 + reps/30)
	return weight * (1 + float64(reps)/30.0)
}

// calculateOneRMBrzycki estimates 1RM using the Brzycki formula
func (m *MetricsService) calculateOneRMBrzycki(weight float64, reps int32) float64 {
	if reps <= 0 {
		return 0
	}
	if reps == 1 {
		return weight
	}
	if reps >= 37 {
		// Formula becomes invalid for reps >= 37, use Epley instead
		return m.calculateOneRM(weight, reps)
	}
	// Brzycki formula: 1RM = weight × (36 / (37 - reps))
	return weight * (36.0 / (37.0 - float64(reps)))
}

func (m *MetricsService) calculateStrengthMetrics(ctx context.Context, session *models.WorkoutSession) (models.StrengthMetrics, error) {
	estimatedOneRMEpley := make(map[string]float64)
	estimatedOneRMBrzycki := make(map[string]float64)

	var pushVolume, pullVolume float64
	var totalPowerOutput float64
	var maxWeightLifted float64

	// Default body weight for Wilks calculation (TODO: get from user profile)
	bodyWeight := 75.0 // kg - this should be retrieved from user profile when available

	// Default gender for Wilks calculation (TODO: get from user profile)
	isMale := true // this should be retrieved from user profile when available

	for _, exercise := range session.Exercises {
		exerciseDetails, err := m.exerciseService.GetExercise(ctx, &pb.GetExerciseRequest{Id: exercise.ExerciseID.Hex()})
		if err != nil {
			continue
		}

		var maxWeightForExercise float64
		var maxRepsAtMaxWeight int32

		for _, set := range exercise.Sets {
			if !set.Completed {
				continue
			}

			weight := float64(set.ActualWeight)
			reps := set.ActualReps

			// Track max weight for this exercise
			if weight > maxWeightForExercise {
				maxWeightForExercise = weight
				maxRepsAtMaxWeight = reps
			}

			// Track overall max weight for Wilks calculation
			if weight > maxWeightLifted {
				maxWeightLifted = weight
			}

			// Calculate push/pull volume
			setVolume := weight * float64(reps)
			isPush := false
			isPull := false

			for _, muscle := range exerciseDetails.PrimaryMuscles {
				if models.PushMuscles[muscle] {
					isPush = true
				}
				if models.PullMuscles[muscle] {
					isPull = true
				}
			}

			if isPush {
				pushVolume += setVolume
			}
			if isPull {
				pullVolume += setVolume
			}

			// Calculate power output (simplified - assuming 1 meter distance for most exercises)
			// Power = (Weight × Distance × Reps) / Time
			// For now, we'll use a simplified calculation without precise distance and time per set
			distance := 1.0 // meters - this is a simplification
			if session.DurationSeconds > 0 {
				timePerSet := float64(session.DurationSeconds) / float64(m.getTotalSets(session)) // approximate time per set
				if timePerSet > 0 {
					setPower := (weight * distance * float64(reps)) / timePerSet
					totalPowerOutput += setPower
				}
			}
		}

		// Calculate 1RM estimates for this exercise
		if maxWeightForExercise > 0 && maxRepsAtMaxWeight > 0 {
			exerciseID := exercise.ExerciseID.Hex()
			estimatedOneRMEpley[exerciseID] = m.calculateOneRM(maxWeightForExercise, maxRepsAtMaxWeight)
			estimatedOneRMBrzycki[exerciseID] = m.calculateOneRMBrzycki(maxWeightForExercise, maxRepsAtMaxWeight)
		}
	}

	// Calculate push:pull ratio
	var pushPullRatio float64
	if pullVolume > 0 {
		pushPullRatio = pushVolume / pullVolume
	}

	// Calculate Wilks score using the heaviest lift
	wilksScore := m.calculateWilksScore(maxWeightLifted, bodyWeight, isMale)

	return models.StrengthMetrics{
		EstimatedOneRMEpley:   estimatedOneRMEpley,
		EstimatedOneRMBrzycki: estimatedOneRMBrzycki,
		WilksScore:            wilksScore,
		PushPullRatio:         pushPullRatio,
		PowerOutput:           totalPowerOutput,
	}, nil
}

// calculateWilksScore calculates the Wilks coefficient for powerlifting comparison
func (m *MetricsService) calculateWilksScore(weightLifted, bodyWeight float64, isMale bool) float64 {
	if weightLifted <= 0 || bodyWeight <= 0 {
		return 0
	}

	var coeffs models.WilksCoefficients
	if isMale {
		coeffs = models.WilksMaleCoefficients
	} else {
		coeffs = models.WilksFemaleCoefficients
	}

	// Wilks formula: Wilks = Weight Lifted × 500 / (a + b×BW + c×BW² + d×BW³ + e×BW⁴ + f×BW⁵)
	bw := bodyWeight
	denominator := coeffs.A +
		coeffs.B*bw +
		coeffs.C*bw*bw +
		coeffs.D*bw*bw*bw +
		coeffs.E*bw*bw*bw*bw +
		coeffs.F*bw*bw*bw*bw*bw

	if denominator == 0 {
		return 0
	}

	return weightLifted * 500.0 / denominator
}

// getTotalSets returns the total number of sets in a workout session
func (m *MetricsService) getTotalSets(session *models.WorkoutSession) int {
	totalSets := 0
	for _, exercise := range session.Exercises {
		totalSets += len(exercise.Sets)
	}
	return totalSets
}

func (m *MetricsService) calculateProgressAdaptationMetrics(ctx context.Context, session *models.WorkoutSession) (models.ProgressAdaptationMetrics, error) {
	progressAdaptationMetrics := models.ProgressAdaptationMetrics{
		WeekOverWeekProgressRate: make(map[string]float64),
		PlateauDetection:         make(map[string]models.PlateauStatus),
		StrengthGainVelocity:     make(map[string]float64),
		AdaptationRate:           make(map[string]float64),
	}

	// Calculate Progressive Overload Index
	poi, err := m.calculateProgressiveOverloadIndex(ctx, session)
	if err != nil {
		// Log error but don't fail - set to 0 if calculation fails
		poi = 0
	}
	progressAdaptationMetrics.ProgressiveOverloadIndex = poi

	// Calculate metrics for each exercise
	for _, exercise := range session.Exercises {
		exerciseID := exercise.ExerciseID.Hex()

		// Calculate Week-over-Week Progress Rate
		progressRate, err := m.calculateWeekOverWeekProgressRate(ctx, session.UserID, exerciseID, session.StartedAt)
		if err == nil {
			progressAdaptationMetrics.WeekOverWeekProgressRate[exerciseID] = progressRate
		}

		// Calculate Plateau Detection
		plateauStatus, err := m.calculatePlateauDetection(ctx, session.UserID, exerciseID, session.StartedAt)
		if err == nil {
			progressAdaptationMetrics.PlateauDetection[exerciseID] = plateauStatus
		}

		// Calculate Strength Gain Velocity
		sgv, err := m.calculateStrengthGainVelocity(ctx, session.UserID, exerciseID, session.StartedAt)
		if err == nil {
			progressAdaptationMetrics.StrengthGainVelocity[exerciseID] = sgv
		}

		// Calculate Adaptation Rate
		adaptationRate, err := m.calculateAdaptationRate(ctx, session.UserID, exerciseID, session.StartedAt)
		if err == nil {
			progressAdaptationMetrics.AdaptationRate[exerciseID] = adaptationRate
		}
	}

	return progressAdaptationMetrics, nil
}

// calculateProgressiveOverloadIndex calculates POI = (Current Week Volume × Avg Intensity) / (Previous Week Volume × Avg Intensity)
func (m *MetricsService) calculateProgressiveOverloadIndex(ctx context.Context, session *models.WorkoutSession) (float64, error) {
	now := session.StartedAt
	currentWeekStart := now.AddDate(0, 0, -int(now.Weekday()))
	previousWeekStart := currentWeekStart.AddDate(0, 0, -7)
	previousWeekEnd := currentWeekStart.AddDate(0, 0, -1)

	// Get current week metrics
	currentWeekMetrics, err := m.calculatePeriodMetrics(ctx, session.UserID, currentWeekStart, now)
	if err != nil {
		return 0, err
	}

	// Get previous week metrics
	previousWeekMetrics, err := m.calculatePeriodMetrics(ctx, session.UserID, previousWeekStart, previousWeekEnd)
	if err != nil {
		return 0, err
	}

	// Calculate average intensity for both weeks (simplified - using total volume as proxy)
	currentWeekVolumeIntensity := currentWeekMetrics.TotalVolumeLoad
	previousWeekVolumeIntensity := previousWeekMetrics.TotalVolumeLoad

	if previousWeekVolumeIntensity == 0 {
		return 0, nil // No previous week data
	}

	poi := currentWeekVolumeIntensity / previousWeekVolumeIntensity
	return poi, nil
}

// calculateWeekOverWeekProgressRate calculates Progress Rate = (Current 1RM - Previous 1RM) / Previous 1RM × 100
func (m *MetricsService) calculateWeekOverWeekProgressRate(ctx context.Context, userID primitive.ObjectID, exerciseID string, currentDate time.Time) (float64, error) {
	// Get current week's best 1RM for this exercise
	currentWeekStart := currentDate.AddDate(0, 0, -int(currentDate.Weekday()))
	currentOneRM, err := m.getBestOneRMForPeriod(ctx, userID, exerciseID, currentWeekStart, currentDate)
	if err != nil {
		return 0, err
	}

	// Get previous week's best 1RM for this exercise
	previousWeekStart := currentWeekStart.AddDate(0, 0, -7)
	previousWeekEnd := currentWeekStart.AddDate(0, 0, -1)
	previousOneRM, err := m.getBestOneRMForPeriod(ctx, userID, exerciseID, previousWeekStart, previousWeekEnd)
	if err != nil {
		return 0, err
	}

	if previousOneRM == 0 {
		return 0, nil // No previous week data
	}

	progressRate := ((currentOneRM - previousOneRM) / previousOneRM) * 100
	return progressRate, nil
}

// calculatePlateauDetection detects if Progress Rate < 1% for 3+ consecutive weeks
func (m *MetricsService) calculatePlateauDetection(ctx context.Context, userID primitive.ObjectID, exerciseID string, currentDate time.Time) (models.PlateauStatus, error) {
	plateauStatus := models.PlateauStatus{}

	// Get progress rates for the last 4 weeks
	var progressRates []float64
	for i := 0; i < 4; i++ {
		weekStart := currentDate.AddDate(0, 0, -int(currentDate.Weekday())-(i*7))
		weekEnd := weekStart.AddDate(0, 0, 6)

		if i == 0 {
			// Current week
			progressRate, err := m.calculateWeekOverWeekProgressRate(ctx, userID, exerciseID, currentDate)
			if err == nil {
				progressRates = append(progressRates, progressRate)
			}
		} else {
			// Previous weeks
			progressRate, err := m.calculateWeekOverWeekProgressRate(ctx, userID, exerciseID, weekEnd)
			if err == nil {
				progressRates = append(progressRates, progressRate)
			}
		}
	}

	if len(progressRates) == 0 {
		return plateauStatus, nil
	}

	plateauStatus.LastProgressRate = progressRates[0]

	// Count consecutive weeks with <1% progress
	consecutiveWeeks := int32(0)
	for _, rate := range progressRates {
		if rate < 1.0 {
			consecutiveWeeks++
		} else {
			break
		}
	}

	plateauStatus.ConsecutiveWeeks = consecutiveWeeks
	plateauStatus.IsPlateaued = consecutiveWeeks >= 3

	// Count weeks since last significant progress (>=1%)
	weeksSinceProgress := int32(0)
	for _, rate := range progressRates {
		if rate >= 1.0 {
			break
		}
		weeksSinceProgress++
	}
	plateauStatus.WeeksSinceProgress = weeksSinceProgress

	return plateauStatus, nil
}

// calculateStrengthGainVelocity calculates SGV = (Current 1RM - Initial 1RM) / Training Weeks
func (m *MetricsService) calculateStrengthGainVelocity(ctx context.Context, userID primitive.ObjectID, exerciseID string, currentDate time.Time) (float64, error) {
	// Get the earliest workout with this exercise (initial 1RM)
	initialOneRM, initialDate, err := m.getInitialOneRM(ctx, userID, exerciseID)
	if err != nil {
		return 0, err
	}

	// Get current best 1RM
	currentOneRM, err := m.getBestOneRMForPeriod(ctx, userID, exerciseID, time.Time{}, currentDate)
	if err != nil {
		return 0, err
	}

	if initialOneRM == 0 {
		return 0, nil // No initial data
	}

	// Calculate training weeks
	trainingWeeks := currentDate.Sub(initialDate).Hours() / (24 * 7) // Convert to weeks
	if trainingWeeks <= 0 {
		return 0, nil
	}

	sgv := (currentOneRM - initialOneRM) / trainingWeeks
	return sgv, nil
}

// calculateAdaptationRate calculates Adaptation Rate = Δ Performance / Δ Volume
func (m *MetricsService) calculateAdaptationRate(ctx context.Context, userID primitive.ObjectID, exerciseID string, currentDate time.Time) (float64, error) {
	// Get current week's volume and 1RM for this exercise
	currentWeekStart := currentDate.AddDate(0, 0, -int(currentDate.Weekday()))
	currentVolume, err := m.getExerciseVolumeForPeriod(ctx, userID, exerciseID, currentWeekStart, currentDate)
	if err != nil {
		return 0, err
	}
	currentOneRM, err := m.getBestOneRMForPeriod(ctx, userID, exerciseID, currentWeekStart, currentDate)
	if err != nil {
		return 0, err
	}

	// Get previous week's volume and 1RM for this exercise
	previousWeekStart := currentWeekStart.AddDate(0, 0, -7)
	previousWeekEnd := currentWeekStart.AddDate(0, 0, -1)
	previousVolume, err := m.getExerciseVolumeForPeriod(ctx, userID, exerciseID, previousWeekStart, previousWeekEnd)
	if err != nil {
		return 0, err
	}
	previousOneRM, err := m.getBestOneRMForPeriod(ctx, userID, exerciseID, previousWeekStart, previousWeekEnd)
	if err != nil {
		return 0, err
	}

	// Calculate deltas
	deltaPerformance := currentOneRM - previousOneRM
	deltaVolume := currentVolume - previousVolume

	if deltaVolume == 0 {
		return 0, nil // No volume change
	}

	adaptationRate := deltaPerformance / deltaVolume
	return adaptationRate, nil
}

// Helper methods for Progress & Adaptation metrics

// getBestOneRMForPeriod gets the best estimated 1RM for an exercise in a given period
func (m *MetricsService) getBestOneRMForPeriod(ctx context.Context, userID primitive.ObjectID, exerciseID string, startDate, endDate time.Time) (float64, error) {
	filter := bson.M{
		"userId": userID,
	}

	if !startDate.IsZero() {
		filter["date"] = bson.M{"$gte": startDate, "$lte": endDate}
	} else {
		filter["date"] = bson.M{"$lte": endDate}
	}

	cursor, err := m.metricsColl.Find(ctx, filter, options.Find().SetSort(bson.D{{"date", -1}}))
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var bestOneRM float64
	for cursor.Next(ctx) {
		var metrics models.WorkoutMetrics
		if err := cursor.Decode(&metrics); err != nil {
			continue
		}

		if oneRM, exists := metrics.StrengthMetrics.EstimatedOneRMEpley[exerciseID]; exists {
			if oneRM > bestOneRM {
				bestOneRM = oneRM
			}
		}
	}

	return bestOneRM, nil
}

// getInitialOneRM gets the first recorded 1RM for an exercise
func (m *MetricsService) getInitialOneRM(ctx context.Context, userID primitive.ObjectID, exerciseID string) (float64, time.Time, error) {
	filter := bson.M{
		"userId": userID,
	}

	cursor, err := m.metricsColl.Find(ctx, filter, options.Find().SetSort(bson.D{{"date", 1}}))
	if err != nil {
		return 0, time.Time{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var metrics models.WorkoutMetrics
		if err := cursor.Decode(&metrics); err != nil {
			continue
		}

		if oneRM, exists := metrics.StrengthMetrics.EstimatedOneRMEpley[exerciseID]; exists && oneRM > 0 {
			return oneRM, metrics.Date, nil
		}
	}

	return 0, time.Time{}, fmt.Errorf("no initial 1RM found for exercise %s", exerciseID)
}

// getExerciseVolumeForPeriod gets the total volume for a specific exercise in a given period
func (m *MetricsService) getExerciseVolumeForPeriod(ctx context.Context, userID primitive.ObjectID, exerciseID string, startDate, endDate time.Time) (float64, error) {
	filter := bson.M{
		"userId": userID,
		"date":   bson.M{"$gte": startDate, "$lte": endDate},
	}

	cursor, err := m.metricsColl.Find(ctx, filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var totalVolume float64
	for cursor.Next(ctx) {
		var metrics models.WorkoutMetrics
		if err := cursor.Decode(&metrics); err != nil {
			continue
		}

		// Find the specific exercise in the workout metrics
		for _, exerciseMetric := range metrics.ExerciseMetrics {
			if exerciseMetric.ExerciseID.Hex() == exerciseID {
				totalVolume += exerciseMetric.VolumeLoad
				break
			}
		}
	}

	return totalVolume, nil
}

func (m *MetricsService) calculateSetMetrics(session *models.WorkoutSession) models.SetMetrics {
	var totalSets, completedSets int32

	for _, exercise := range session.Exercises {
		for _, set := range exercise.Sets {
			totalSets++
			if set.Completed {
				completedSets++
			}
		}
	}

	var completionRate float64
	if totalSets > 0 {
		completionRate = float64(completedSets) / float64(totalSets)
	}

	return models.SetMetrics{
		TotalSets:      totalSets,
		CompletedSets:  completedSets,
		CompletionRate: completionRate,
	}
}

func (m *MetricsService) calculateExerciseMetrics(ctx context.Context, session *models.WorkoutSession) ([]models.ExerciseMetrics, error) {
	var exerciseMetrics []models.ExerciseMetrics

	for _, exercise := range session.Exercises {
		exerciseDetails, err := m.exerciseService.GetExercise(ctx, &pb.GetExerciseRequest{Id: exercise.ExerciseID.Hex()})
		if err != nil {
			continue
		}

		var volumeLoad, tonnage float64
		var setCount, totalReps int32
		var weightSum, maxWeight float64

		for _, set := range exercise.Sets {
			if !set.Completed {
				continue
			}

			setVolumeLoad := float64(set.ActualReps) * float64(set.ActualWeight)
			setTonnage := float64(set.ActualReps) * float64(set.ActualWeight)

			volumeLoad += setVolumeLoad
			tonnage += setTonnage
			setCount++
			totalReps += set.ActualReps
			weightSum += float64(set.ActualWeight)

			if float64(set.ActualWeight) > maxWeight {
				maxWeight = float64(set.ActualWeight)
			}
		}

		var averageWeight float64
		if setCount > 0 {
			averageWeight = weightSum / float64(setCount)
		}

		exerciseMetric := models.ExerciseMetrics{
			ExerciseID:       exercise.ExerciseID,
			ExerciseName:     exerciseDetails.Name,
			VolumeLoad:       volumeLoad,
			Tonnage:          tonnage,
			SetCount:         setCount,
			TotalReps:        totalReps,
			AverageWeight:    averageWeight,
			MaxWeight:        maxWeight,
			PrimaryMuscles:   exerciseDetails.PrimaryMuscles,
			SecondaryMuscles: exerciseDetails.SecondaryMuscles,
		}

		exerciseMetrics = append(exerciseMetrics, exerciseMetric)
	}

	return exerciseMetrics, nil
}

func (m *MetricsService) calculateRecoveryFatigueMetrics(ctx context.Context, session *models.WorkoutSession) (models.RecoveryFatigueMetrics, error) {
	// Calculate acute and chronic workloads
	acuteWorkload, err := m.calculateAcuteWorkload(ctx, session.UserID, session.StartedAt)
	if err != nil {
		return models.RecoveryFatigueMetrics{}, err
	}

	chronicWorkload, err := m.calculateChronicWorkload(ctx, session.UserID, session.StartedAt)
	if err != nil {
		return models.RecoveryFatigueMetrics{}, err
	}

	// Calculate ACWR (Acute:Chronic Workload Ratio)
	var acwr float64
	if chronicWorkload > 0 {
		acwr = acuteWorkload / chronicWorkload
	}

	// Calculate current session volume and intensity for TSS
	volumeMetrics, err := m.calculateVolumeMetrics(ctx, session)
	if err != nil {
		return models.RecoveryFatigueMetrics{}, err
	}

	intensityMetrics, err := m.calculateIntensityMetrics(ctx, session)
	if err != nil {
		return models.RecoveryFatigueMetrics{}, err
	}

	// Calculate Training Strain Score (TSS)
	tss := (volumeMetrics.TotalVolumeLoad * intensityMetrics.AverageIntensity * float64(session.RPERating)) / 1000.0

	// Calculate Recovery Need Index (RNI)
	hoursSinceLastSession, err := m.getHoursSinceLastSession(ctx, session.UserID, session.StartedAt)
	if err != nil {
		hoursSinceLastSession = 24.0 // Default to 24 hours if can't determine
	}

	var rni float64
	if hoursSinceLastSession > 0 {
		rni = (volumeMetrics.TotalVolumeLoad * intensityMetrics.AverageIntensity * float64(session.RPERating)) / (hoursSinceLastSession * 10.0)
	}

	// Calculate Recovery Score (simplified - using defaults from user profile)
	profile := models.DefaultProfile
	recoveryScore := hoursSinceLastSession * profile.SleepQuality * profile.NutritionFactor

	// Calculate Fatigue Accumulation Index (FAI)
	dailyTSS, err := m.calculateDailyTSS(ctx, session.UserID, session.StartedAt)
	if err != nil {
		dailyTSS = tss // Use current session TSS if historical data unavailable
	}

	dailyRecoveryScore, err := m.calculateDailyRecoveryScore(ctx, session.UserID, session.StartedAt)
	if err != nil {
		dailyRecoveryScore = recoveryScore // Use current recovery score if historical data unavailable
	}

	fai := dailyTSS - dailyRecoveryScore

	// Calculate Overtraining Risk Score (OTS)
	rpeTrend, err := m.calculateRPETrend(ctx, session.UserID, session.StartedAt)
	if err != nil {
		rpeTrend = 1.0 // Neutral trend if can't calculate
	}

	performanceDecline, err := m.calculatePerformanceDecline(ctx, session.UserID, session.StartedAt)
	if err != nil {
		performanceDecline = 1.0 // No decline if can't calculate
	}

	recoveryQuality := (profile.SleepQuality + profile.NutritionFactor) / 2.0 // Average of sleep and nutrition
	var ots float64
	if recoveryQuality > 0 {
		ots = (acwr * rpeTrend * performanceDecline) / recoveryQuality
	}

	return models.RecoveryFatigueMetrics{
		AcuteChronicWorkloadRatio: acwr,
		TrainingStrainScore:       tss,
		RecoveryNeedIndex:         rni,
		FatigueAccumulationIndex:  fai,
		OvertrainingRiskScore:     ots,
		RecoveryScore:             recoveryScore,
	}, nil
}

func (m *MetricsService) calculateBodyCompositionMetrics(ctx context.Context, session *models.WorkoutSession) (models.BodyCompositionMetrics, error) {
	// Use default profile values (TODO: get from user profile when available)
	profile := models.DefaultProfile

	// Calculate BMI
	bmi := profile.BodyWeight / (profile.Height * profile.Height)

	// Calculate Strength-to-Weight Ratio (Big 3 lifts)
	big3Total, err := m.calculateBig3Total(ctx, session)
	if err != nil {
		return models.BodyCompositionMetrics{}, err
	}

	strengthToWeightRatio := big3Total / profile.BodyWeight

	// Calculate Allometric Scaling Score (using heaviest lift from session)
	maxWeight := m.getMaxWeightFromSession(session)
	allometricScalingScore := maxWeight / math.Pow(profile.BodyWeight, 2.0/3.0)

	// Calculate Ponderal Index
	ponderalIndex := profile.BodyWeight / (profile.Height * profile.Height * profile.Height)

	return models.BodyCompositionMetrics{
		BMI:                    bmi,
		StrengthToWeightRatio:  strengthToWeightRatio,
		AllometricScalingScore: allometricScalingScore,
		PonderalIndex:          ponderalIndex,
	}, nil
}

// Helper methods for Recovery & Fatigue Metrics

func (m *MetricsService) calculateAcuteWorkload(ctx context.Context, userID primitive.ObjectID, currentDate time.Time) (float64, error) {
	// Get workouts from last 7 days
	sevenDaysAgo := currentDate.AddDate(0, 0, -7)
	return m.getWorkloadForPeriod(ctx, userID, sevenDaysAgo, currentDate)
}

func (m *MetricsService) calculateChronicWorkload(ctx context.Context, userID primitive.ObjectID, currentDate time.Time) (float64, error) {
	// Get workouts from last 28 days
	twentyEightDaysAgo := currentDate.AddDate(0, 0, -28)
	return m.getWorkloadForPeriod(ctx, userID, twentyEightDaysAgo, currentDate)
}

func (m *MetricsService) getWorkloadForPeriod(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (float64, error) {
	filter := bson.M{
		"userId": userID,
		"date":   bson.M{"$gte": startDate, "$lte": endDate},
	}

	cursor, err := m.metricsColl.Find(ctx, filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var totalWorkload float64
	for cursor.Next(ctx) {
		var metrics models.WorkoutMetrics
		if err := cursor.Decode(&metrics); err != nil {
			continue
		}
		// Use TSS as workload measure
		totalWorkload += metrics.RecoveryFatigueMetrics.TrainingStrainScore
	}

	return totalWorkload, nil
}

func (m *MetricsService) getHoursSinceLastSession(ctx context.Context, userID primitive.ObjectID, currentDate time.Time) (float64, error) {
	filter := bson.M{
		"userId": userID,
		"date":   bson.M{"$lt": currentDate},
	}

	opts := options.FindOne().SetSort(bson.D{{"date", -1}})
	var lastMetrics models.WorkoutMetrics
	err := m.metricsColl.FindOne(ctx, filter, opts).Decode(&lastMetrics)
	if err == mongo.ErrNoDocuments {
		return 24.0, nil // Default to 24 hours if no previous session
	} else if err != nil {
		return 0, err
	}

	duration := currentDate.Sub(lastMetrics.Date)
	return duration.Hours(), nil
}

func (m *MetricsService) calculateDailyTSS(ctx context.Context, userID primitive.ObjectID, currentDate time.Time) (float64, error) {
	// Get TSS for the current day
	startOfDay := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	filter := bson.M{
		"userId": userID,
		"date":   bson.M{"$gte": startOfDay, "$lt": endOfDay},
	}

	cursor, err := m.metricsColl.Find(ctx, filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var totalTSS float64
	for cursor.Next(ctx) {
		var metrics models.WorkoutMetrics
		if err := cursor.Decode(&metrics); err != nil {
			continue
		}
		totalTSS += metrics.RecoveryFatigueMetrics.TrainingStrainScore
	}

	return totalTSS, nil
}

func (m *MetricsService) calculateDailyRecoveryScore(ctx context.Context, userID primitive.ObjectID, currentDate time.Time) (float64, error) {
	// Simplified: return a default recovery score based on time since last session
	hoursSinceLastSession, err := m.getHoursSinceLastSession(ctx, userID, currentDate)
	if err != nil {
		return 0, err
	}

	profile := models.DefaultProfile
	return hoursSinceLastSession * profile.SleepQuality * profile.NutritionFactor, nil
}

func (m *MetricsService) calculateRPETrend(ctx context.Context, userID primitive.ObjectID, currentDate time.Time) (float64, error) {
	// Get RPE values from last 7 days
	sevenDaysAgo := currentDate.AddDate(0, 0, -7)
	filter := bson.M{
		"userId": userID,
		"date":   bson.M{"$gte": sevenDaysAgo, "$lte": currentDate},
	}

	cursor, err := m.metricsColl.Find(ctx, filter, options.Find().SetSort(bson.D{{"date", 1}}))
	if err != nil {
		return 1.0, err
	}
	defer cursor.Close(ctx)

	var rpeValues []float64
	for cursor.Next(ctx) {
		var metrics models.WorkoutMetrics
		if err := cursor.Decode(&metrics); err != nil {
			continue
		}
		rpeValues = append(rpeValues, metrics.PerformanceMetrics.AverageRPE)
	}

	if len(rpeValues) < 2 {
		return 1.0, nil // Neutral trend if insufficient data
	}

	// Calculate simple trend (last value / first value)
	return rpeValues[len(rpeValues)-1] / rpeValues[0], nil
}

func (m *MetricsService) calculatePerformanceDecline(ctx context.Context, userID primitive.ObjectID, currentDate time.Time) (float64, error) {
	// Get volume metrics from last 14 days to assess performance decline
	fourteenDaysAgo := currentDate.AddDate(0, 0, -14)
	filter := bson.M{
		"userId": userID,
		"date":   bson.M{"$gte": fourteenDaysAgo, "$lte": currentDate},
	}

	cursor, err := m.metricsColl.Find(ctx, filter, options.Find().SetSort(bson.D{{"date", 1}}))
	if err != nil {
		return 1.0, err
	}
	defer cursor.Close(ctx)

	var volumeValues []float64
	for cursor.Next(ctx) {
		var metrics models.WorkoutMetrics
		if err := cursor.Decode(&metrics); err != nil {
			continue
		}
		volumeValues = append(volumeValues, metrics.VolumeMetrics.TotalVolumeLoad)
	}

	if len(volumeValues) < 2 {
		return 1.0, nil // No decline if insufficient data
	}

	// Calculate performance decline (higher values indicate more decline)
	recentAvg := average(volumeValues[len(volumeValues)/2:])
	earlierAvg := average(volumeValues[:len(volumeValues)/2])

	if earlierAvg == 0 {
		return 1.0, nil
	}

	decline := earlierAvg / recentAvg // Values > 1.0 indicate decline
	return decline, nil
}

// Helper methods for Body Composition Metrics

func (m *MetricsService) calculateBig3Total(ctx context.Context, session *models.WorkoutSession) (float64, error) {
	var big3Total float64

	for _, exercise := range session.Exercises {
		exerciseDetails, err := m.exerciseService.GetExercise(ctx, &pb.GetExerciseRequest{Id: exercise.ExerciseID.Hex()})
		if err != nil {
			continue
		}

		// Check if this is a Big 3 exercise (case-insensitive)
		exerciseName := strings.ToLower(exerciseDetails.Name)
		isBig3 := false
		for big3Exercise := range models.Big3Exercises {
			if strings.Contains(exerciseName, big3Exercise) {
				isBig3 = true
				break
			}
		}

		if isBig3 {
			// Get the heaviest weight lifted for this exercise
			var maxWeight float64
			var maxReps int32
			for _, set := range exercise.Sets {
				if set.Completed && float64(set.ActualWeight) > maxWeight {
					maxWeight = float64(set.ActualWeight)
					maxReps = set.ActualReps
				}
			}

			// Calculate 1RM for this exercise and add to Big 3 total
			if maxWeight > 0 {
				oneRM := m.calculateOneRM(maxWeight, maxReps)
				big3Total += oneRM
			}
		}
	}

	return big3Total, nil
}

func (m *MetricsService) getMaxWeightFromSession(session *models.WorkoutSession) float64 {
	var maxWeight float64
	for _, exercise := range session.Exercises {
		for _, set := range exercise.Sets {
			if set.Completed && float64(set.ActualWeight) > maxWeight {
				maxWeight = float64(set.ActualWeight)
			}
		}
	}
	return maxWeight
}

// calculateMuscleSpecificMetrics calculates muscle-specific metrics for a workout session
func (m *MetricsService) calculateMuscleSpecificMetrics(ctx context.Context, session *models.WorkoutSession) (models.MuscleSpecificMetrics, error) {
	muscleSpecificMetrics := models.MuscleSpecificMetrics{
		MuscleGroupDistribution: make(map[string]float64),
		MuscleImbalanceIndex:    make(map[string]float64),
		AntagonistRatio:         make(map[string]float64),
		StimulusToFatigueRatio:  make(map[string]float64),
	}

	// Calculate volume metrics to get total volume
	volumeMetrics, err := m.calculateVolumeMetrics(ctx, session)
	if err != nil {
		return muscleSpecificMetrics, err
	}

	// Calculate Muscle Group Volume Distribution
	totalVolume := volumeMetrics.TotalVolumeLoad
	if totalVolume > 0 {
		for muscleGroup, volume := range volumeMetrics.MuscleGroupVolume {
			distribution := (volume / totalVolume) * 100
			muscleSpecificMetrics.MuscleGroupDistribution[muscleGroup] = distribution
		}
	}

	// Calculate Muscle Imbalance Index (Left vs Right sides)
	muscleStrengths := make(map[string]float64)
	for _, exercise := range session.Exercises {
		exerciseDetails, err := m.exerciseService.GetExercise(ctx, &pb.GetExerciseRequest{Id: exercise.ExerciseID.Hex()})
		if err != nil {
			continue
		}

		// Get the heaviest weight for this exercise
		var maxWeight float64
		var maxReps int32
		for _, set := range exercise.Sets {
			if set.Completed && float64(set.ActualWeight) > maxWeight {
				maxWeight = float64(set.ActualWeight)
				maxReps = set.ActualReps
			}
		}

		if maxWeight > 0 {
			oneRM := m.calculateOneRM(maxWeight, maxReps)
			for _, muscle := range exerciseDetails.PrimaryMuscles {
				if currentStrength, exists := muscleStrengths[muscle]; exists {
					muscleStrengths[muscle] = math.Max(currentStrength, oneRM)
				} else {
					muscleStrengths[muscle] = oneRM
				}
			}
		}
	}

	// Calculate imbalance indices for left-right pairs
	for leftMuscle, rightMuscle := range models.LeftRightPairs {
		leftStrength, leftExists := muscleStrengths[leftMuscle]
		rightStrength, rightExists := muscleStrengths[rightMuscle]

		if leftExists && rightExists && (leftStrength > 0 || rightStrength > 0) {
			averageStrength := (leftStrength + rightStrength) / 2
			if averageStrength > 0 {
				imbalance := math.Abs(leftStrength-rightStrength) / averageStrength * 100
				baseMuscle := strings.TrimPrefix(strings.TrimPrefix(leftMuscle, "left "), "right ")
				muscleSpecificMetrics.MuscleImbalanceIndex[baseMuscle] = imbalance
			}
		}
	}

	// Calculate Antagonist Ratios
	for agonist, antagonist := range models.AntagonistPairs {
		agonistStrength, agonistExists := muscleStrengths[agonist]
		antagonistStrength, antagonistExists := muscleStrengths[antagonist]

		if agonistExists && antagonistExists && agonistStrength > 0 {
			ratio := antagonistStrength / agonistStrength
			muscleSpecificMetrics.AntagonistRatio[agonist+"_"+antagonist] = ratio
		}
	}

	// Calculate Stimulus-to-Fatigue Ratio
	// Get previous session metrics for baseline comparison
	previousMetrics, err := m.getPreviousSessionMetrics(ctx, session.UserID, session.StartedAt)
	if err == nil && previousMetrics != nil {
		for muscleGroup, currentVolume := range volumeMetrics.MuscleGroupVolume {
			if previousVolume, exists := previousMetrics.VolumeMetrics.MuscleGroupVolume[muscleGroup]; exists && previousVolume > 0 {
				performanceGain := (currentVolume - previousVolume) / previousVolume
				fatigueFactor := float64(session.RPERating) * currentVolume
				if fatigueFactor > 0 {
					sfr := performanceGain / fatigueFactor
					muscleSpecificMetrics.StimulusToFatigueRatio[muscleGroup] = sfr
				}
			}
		}
	}

	return muscleSpecificMetrics, nil
}

// calculateWorkCapacityMetrics calculates work capacity metrics for a workout session
func (m *MetricsService) calculateWorkCapacityMetrics(ctx context.Context, session *models.WorkoutSession) (models.WorkCapacityMetrics, error) {
	workCapacityMetrics := models.WorkCapacityMetrics{}

	var totalWorkCapacity, totalTimeUnderTension, totalMechanicalTension float64
	workoutTimeMinutes := float64(session.DurationSeconds) / 60.0

	// Calculate volume metrics for density calculations
	volumeMetrics, err := m.calculateVolumeMetrics(ctx, session)
	if err != nil {
		return workCapacityMetrics, err
	}

	for _, exercise := range session.Exercises {
		exerciseDetails, err := m.exerciseService.GetExercise(ctx, &pb.GetExerciseRequest{Id: exercise.ExerciseID.Hex()})
		if err != nil {
			continue
		}

		// Get exercise tempo (default if not specified)
		tempo := models.DefaultExerciseTempos["default"]
		exerciseName := strings.ToLower(exerciseDetails.Name)
		for knownExercise, knownTempo := range models.DefaultExerciseTempos {
			if strings.Contains(exerciseName, knownExercise) {
				tempo = knownTempo
				break
			}
		}

		for i, set := range exercise.Sets {
			if !set.Completed {
				continue
			}

			weight := float64(set.ActualWeight)
			reps := float64(set.ActualReps)

			// Calculate rest time factor (assume default rest time between sets)
			restTime := models.DefaultRestTime
			if i == 0 {
				restTime = 0 // No rest before first set
			}
			restFactor := 1 - (restTime / 300.0) // 300 seconds = 5 minutes max rest
			if restFactor < 0 {
				restFactor = 0
			}

			// Total Work Capacity: TWC = Σ(Sets × Reps × Weight × (1 - Rest Time/300))
			setWorkCapacity := 1 * reps * weight * restFactor
			totalWorkCapacity += setWorkCapacity

			// Time Under Tension: TUT = Σ(Reps × Tempo in seconds)
			repTempo := tempo.Eccentric + tempo.Pause1 + tempo.Concentric + tempo.Pause2
			setTUT := reps * repTempo
			totalTimeUnderTension += setTUT

			// Mechanical Tension Score: MTS = Weight × TUT × (RPE/10)
			rpeFactor := float64(session.RPERating) / 10.0
			setMTS := weight * setTUT * rpeFactor
			totalMechanicalTension += setMTS
		}
	}

	workCapacityMetrics.TotalWorkCapacity = totalWorkCapacity
	workCapacityMetrics.TimeUnderTension = totalTimeUnderTension
	workCapacityMetrics.MechanicalTensionScore = totalMechanicalTension

	// Density Training Index: Density = Total Volume / Total Time
	if workoutTimeMinutes > 0 {
		workCapacityMetrics.DensityTrainingIndex = volumeMetrics.TotalVolumeLoad / workoutTimeMinutes
	}

	// Density Progress Percent: Progress% = (Current Density - Previous Density) / Previous Density × 100
	previousMetrics, err := m.getPreviousSessionMetrics(ctx, session.UserID, session.StartedAt)
	if err == nil && previousMetrics != nil {
		previousDensity := previousMetrics.WorkCapacityMetrics.DensityTrainingIndex
		currentDensity := workCapacityMetrics.DensityTrainingIndex
		if previousDensity > 0 {
			densityProgress := ((currentDensity - previousDensity) / previousDensity) * 100
			workCapacityMetrics.DensityProgressPercent = densityProgress
		}
	}

	return workCapacityMetrics, nil
}

// getPreviousSessionMetrics gets the metrics from the previous workout session for comparison
func (m *MetricsService) getPreviousSessionMetrics(ctx context.Context, userID primitive.ObjectID, currentDate time.Time) (*models.WorkoutMetrics, error) {
	filter := bson.M{
		"userId": userID,
		"date":   bson.M{"$lt": currentDate},
	}

	opts := options.FindOne().SetSort(bson.D{{"date", -1}})
	var previousMetrics models.WorkoutMetrics
	err := m.metricsColl.FindOne(ctx, filter, opts).Decode(&previousMetrics)
	if err == mongo.ErrNoDocuments {
		return nil, nil // No previous session found
	} else if err != nil {
		return nil, err
	}

	return &previousMetrics, nil
}

// calculateTrainingPatternMetrics calculates training pattern analytics metrics
func (m *MetricsService) calculateTrainingPatternMetrics(ctx context.Context, session *models.WorkoutSession) (models.TrainingPatternMetrics, error) {
	trainingPatternMetrics := models.TrainingPatternMetrics{}

	// Calculate Volume and Intensity for Recovery Time calculation
	volumeMetrics, err := m.calculateVolumeMetrics(ctx, session)
	if err != nil {
		return trainingPatternMetrics, err
	}

	intensityMetrics, err := m.calculateIntensityMetrics(ctx, session)
	if err != nil {
		return trainingPatternMetrics, err
	}

	// Calculate Recovery Time = 24 × (Volume/10) × (Intensity/70) × (RPE/7)
	volume := volumeMetrics.TotalVolumeLoad
	intensity := intensityMetrics.AverageIntensity
	rpe := float64(session.RPERating)

	if volume > 0 && intensity > 0 && rpe > 0 {
		recoveryTime := 24.0 * (volume / models.BaseVolumeThreshold) * (intensity / models.BaseIntensityPercent) * (rpe / models.BaseRPEThreshold)
		trainingPatternMetrics.RecoveryTime = recoveryTime

		// Calculate Optimal Frequency = Recovery Time + 24-48 hours
		optimalFrequency := recoveryTime + models.MinimumRecoveryHours
		if optimalFrequency > recoveryTime+models.MaximumRecoveryHours {
			optimalFrequency = recoveryTime + models.MaximumRecoveryHours
		}
		trainingPatternMetrics.OptimalFrequency = optimalFrequency
	}

	// Calculate Exercise Selection Diversity Index = Unique Exercises / Total Exercises × 100
	uniqueExercises := make(map[string]bool)
	totalExercises := len(session.Exercises)

	for _, exercise := range session.Exercises {
		uniqueExercises[exercise.ExerciseID.Hex()] = true
	}

	if totalExercises > 0 {
		diversityIndex := (float64(len(uniqueExercises)) / float64(totalExercises)) * 100.0
		trainingPatternMetrics.ExerciseSelectionIndex = diversityIndex
	}

	// Calculate Workout Completion Rate = (Completed Sets / Planned Sets) × 100
	setMetrics := m.calculateSetMetrics(session)
	trainingPatternMetrics.WorkoutCompletionRate = setMetrics.CompletionRate * 100.0

	// Calculate Consistency Score = (Actual Workouts / Planned Workouts) × 100
	// For now, we'll use the workout completion rate as a proxy since we don't have planned workout data
	// This could be enhanced when we have more detailed workout planning functionality
	trainingPatternMetrics.ConsistencyScore = trainingPatternMetrics.WorkoutCompletionRate

	return trainingPatternMetrics, nil
}

// calculatePeriodizationMetrics calculates advanced periodization metrics
func (m *MetricsService) calculatePeriodizationMetrics(ctx context.Context, session *models.WorkoutSession) (models.PeriodizationMetrics, error) {
	periodizationMetrics := models.PeriodizationMetrics{}

	// Calculate current session's TSS for the calculations
	recoveryFatigueMetrics, err := m.calculateRecoveryFatigueMetrics(ctx, session)
	if err != nil {
		return periodizationMetrics, err
	}

	currentTSS := recoveryFatigueMetrics.TrainingStrainScore

	// Calculate Chronic Training Load (CTL) = Exponentially Weighted Average of daily TSS over 42 days
	ctl, err := m.calculateExponentiallyWeightedAverage(ctx, session.UserID, session.StartedAt, models.ChronicTrainingLoadDays)
	if err != nil {
		// If calculation fails, use current TSS as fallback
		ctl = currentTSS
	}
	periodizationMetrics.ChronicTrainingLoad = ctl

	// Calculate Acute Training Load (ATL) = Exponentially Weighted Average of daily TSS over 7 days
	atl, err := m.calculateExponentiallyWeightedAverage(ctx, session.UserID, session.StartedAt, models.AcuteTrainingLoadDays)
	if err != nil {
		// If calculation fails, use current TSS as fallback
		atl = currentTSS
	}
	periodizationMetrics.AcuteTrainingLoad = atl

	// Calculate Training Stress Balance (TSB) = CTL - ATL
	tsb := ctl - atl
	periodizationMetrics.TrainingStressBalance = tsb

	// Calculate Form/Freshness Index (FFI) = TSB / CTL × 100
	var ffi float64
	if ctl > 0 {
		ffi = (tsb / ctl) * 100.0
	}
	periodizationMetrics.FormFreshnessIndex = ffi

	return periodizationMetrics, nil
}

// calculateExponentiallyWeightedAverage calculates the exponentially weighted average of TSS over a given period
func (m *MetricsService) calculateExponentiallyWeightedAverage(ctx context.Context, userID primitive.ObjectID, currentDate time.Time, days int) (float64, error) {
	// Get TSS values for the specified period
	startDate := currentDate.AddDate(0, 0, -days)

	filter := bson.M{
		"userId": userID,
		"date":   bson.M{"$gte": startDate, "$lte": currentDate},
	}

	cursor, err := m.metricsColl.Find(ctx, filter, options.Find().SetSort(bson.D{{"date", 1}}))
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var tssValues []float64
	var dates []time.Time

	for cursor.Next(ctx) {
		var metrics models.WorkoutMetrics
		if err := cursor.Decode(&metrics); err != nil {
			continue
		}
		tssValues = append(tssValues, metrics.RecoveryFatigueMetrics.TrainingStrainScore)
		dates = append(dates, metrics.Date)
	}

	if len(tssValues) == 0 {
		return 0, nil
	}

	// Calculate exponentially weighted average
	// More recent values have higher weight
	var weightedSum, totalWeight float64

	for i, tss := range tssValues {
		// Calculate days ago from current date
		daysAgo := currentDate.Sub(dates[i]).Hours() / 24.0

		// Apply exponential decay: weight = e^(-decay_factor * days_ago)
		weight := math.Exp(-models.ExponentialDecayFactor * daysAgo)

		weightedSum += tss * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0, nil
	}

	return weightedSum / totalWeight, nil
}

func (m *MetricsService) SaveWorkoutMetrics(ctx context.Context, metrics *models.WorkoutMetrics) error {
	result, err := m.metricsColl.InsertOne(ctx, metrics)
	if err != nil {
		return fmt.Errorf("failed to save workout metrics: %w", err)
	}

	metrics.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (m *MetricsService) getUserMetrics(ctx context.Context, userID primitive.ObjectID, period string) (*models.UserMetrics, error) {
	var userMetrics models.UserMetrics
	err := m.userMetricsColl.FindOne(ctx, bson.M{"userId": userID}).Decode(&userMetrics)
	if err == mongo.ErrNoDocuments {
		userMetrics = models.UserMetrics{
			UserID:    userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user metrics: %w", err)
	}

	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var err2 error
	userMetrics.WeeklyMetrics, err2 = m.calculatePeriodMetrics(ctx, userID, weekStart, now)
	if err2 != nil {
		return nil, fmt.Errorf("failed to calculate weekly metrics: %w", err2)
	}

	userMetrics.MonthlyMetrics, err2 = m.calculatePeriodMetrics(ctx, userID, monthStart, now)
	if err2 != nil {
		return nil, fmt.Errorf("failed to calculate monthly metrics: %w", err2)
	}

	userMetrics.AllTimeMetrics, err2 = m.calculatePeriodMetrics(ctx, userID, time.Time{}, now)
	if err2 != nil {
		return nil, fmt.Errorf("failed to calculate all-time metrics: %w", err2)
	}

	// Calculate volume landmarks
	volumeLandmarks, err2 := m.CalculateVolumeLandmarks(ctx, userID)
	if err2 != nil {
		return nil, fmt.Errorf("failed to calculate volume landmarks: %w", err2)
	}
	userMetrics.VolumeLandmarks = *volumeLandmarks

	return &userMetrics, nil
}

func (m *MetricsService) calculatePeriodMetrics(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (models.PeriodMetrics, error) {
	filter := bson.M{"userId": userID}
	if !startDate.IsZero() {
		filter["date"] = bson.M{"$gte": startDate, "$lte": endDate}
	} else {
		filter["date"] = bson.M{"$lte": endDate}
	}

	cursor, err := m.metricsColl.Find(ctx, filter, options.Find().SetSort(bson.D{{"date", 1}}))
	if err != nil {
		return models.PeriodMetrics{}, fmt.Errorf("failed to query workout metrics: %w", err)
	}
	defer cursor.Close(ctx)

	var workoutMetrics []models.WorkoutMetrics
	if err := cursor.All(ctx, &workoutMetrics); err != nil {
		return models.PeriodMetrics{}, fmt.Errorf("failed to decode workout metrics: %w", err)
	}

	var totalVolumeLoad, totalTonnage, totalWorkoutTime float64
	var totalWorkouts int32
	muscleGroupVolume := make(map[string]float64)
	var volumeProgression []float64

	for _, metrics := range workoutMetrics {
		totalVolumeLoad += metrics.VolumeMetrics.TotalVolumeLoad
		totalTonnage += metrics.VolumeMetrics.Tonnage
		totalWorkoutTime += float64(metrics.WorkoutDurationSecs)
		totalWorkouts++

		for muscle, volume := range metrics.VolumeMetrics.MuscleGroupVolume {
			muscleGroupVolume[muscle] += volume
		}

		volumeProgression = append(volumeProgression, metrics.VolumeMetrics.TotalVolumeLoad)
	}

	var averageWorkoutTime float64
	if totalWorkouts > 0 {
		averageWorkoutTime = totalWorkoutTime / float64(totalWorkouts)
	}

	return models.PeriodMetrics{
		StartDate:          startDate,
		EndDate:            endDate,
		TotalWorkouts:      totalWorkouts,
		TotalVolumeLoad:    totalVolumeLoad,
		TotalTonnage:       totalTonnage,
		AverageWorkoutTime: averageWorkoutTime,
		MuscleGroupVolume:  muscleGroupVolume,
		VolumeProgression:  volumeProgression,
	}, nil
}

func (m *MetricsService) CalculateVolumeTrends(ctx context.Context, userID primitive.ObjectID, period string) (*models.TrendMetrics, error) {
	var startDate time.Time
	now := time.Now()

	switch period {
	case "weekly":
		startDate = now.AddDate(0, 0, -7*12) // Last 12 weeks
	case "monthly":
		startDate = now.AddDate(0, -12, 0) // Last 12 months
	default:
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	filter := bson.M{
		"userId": userID,
		"date":   bson.M{"$gte": startDate, "$lte": now},
	}

	cursor, err := m.metricsColl.Find(ctx, filter, options.Find().SetSort(bson.D{{"date", 1}}))
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics for trends: %w", err)
	}
	defer cursor.Close(ctx)

	var workoutMetrics []models.WorkoutMetrics
	if err := cursor.All(ctx, &workoutMetrics); err != nil {
		return nil, fmt.Errorf("failed to decode metrics: %w", err)
	}

	var volumeProgression []models.VolumeProgressionPoint
	var previousVolume float64
	var volumeGrowthRate float64

	for _, metrics := range workoutMetrics {
		point := models.VolumeProgressionPoint{
			Date:       metrics.Date,
			VolumeLoad: metrics.VolumeMetrics.TotalVolumeLoad,
			Tonnage:    metrics.VolumeMetrics.Tonnage,
		}
		volumeProgression = append(volumeProgression, point)

		if previousVolume > 0 {
			growth := (metrics.VolumeMetrics.TotalVolumeLoad - previousVolume) / previousVolume * 100
			volumeGrowthRate = growth
		}
		previousVolume = metrics.VolumeMetrics.TotalVolumeLoad
	}

	muscleGroupTrends := make(map[string][]float64)
	for _, metrics := range workoutMetrics {
		for muscle, volume := range metrics.VolumeMetrics.MuscleGroupVolume {
			muscleGroupTrends[muscle] = append(muscleGroupTrends[muscle], volume)
		}
	}

	var performanceTrends []float64
	for _, metrics := range workoutMetrics {
		performanceTrends = append(performanceTrends, metrics.PerformanceMetrics.AverageRPE)
	}

	return &models.TrendMetrics{
		Period:            period,
		VolumeProgression: volumeProgression,
		VolumeGrowthRate:  volumeGrowthRate,
		MuscleGroupTrends: muscleGroupTrends,
		PerformanceTrends: performanceTrends,
	}, nil
}

func (m *MetricsService) CalculateVolumeLandmarks(ctx context.Context, userID primitive.ObjectID) (*models.VolumeLandmarks, error) {
	// Get recent workout metrics (last 8 weeks for better landmark calculation)
	now := time.Now()
	eightWeeksAgo := now.AddDate(0, 0, -56) // 8 weeks * 7 days

	filter := bson.M{
		"userId": userID,
		"date":   bson.M{"$gte": eightWeeksAgo, "$lte": now},
	}

	cursor, err := m.metricsColl.Find(ctx, filter, options.Find().SetSort(bson.D{{"date", -1}}))
	if err != nil {
		return nil, fmt.Errorf("failed to query recent metrics: %w", err)
	}
	defer cursor.Close(ctx)

	var recentMetrics []models.WorkoutMetrics
	if err := cursor.All(ctx, &recentMetrics); err != nil {
		return nil, fmt.Errorf("failed to decode recent metrics: %w", err)
	}

	mev := make(map[string]float64)
	mav := make(map[string]float64)
	mrv := make(map[string]float64)

	// Group volume data by muscle group
	muscleGroupVolumes := make(map[string][]float64)
	for _, metrics := range recentMetrics {
		for muscle, volume := range metrics.VolumeMetrics.MuscleGroupVolume {
			muscleGroupVolumes[muscle] = append(muscleGroupVolumes[muscle], volume)
		}
	}

	// Calculate landmarks for each muscle group
	for muscle, volumes := range muscleGroupVolumes {
		if len(volumes) < 3 {
			// Not enough data, use default values
			mev[muscle] = 10.0 // Default MEV (10 sets/week equivalent)
			mav[muscle] = 17.5 // Default MAV (17.5 sets/week equivalent)
			mrv[muscle] = 22.5 // Default MRV (22.5 sets/week equivalent)
			continue
		}

		avgVolume := average(volumes)
		stdDev := standardDeviation(volumes, avgVolume)

		// MEV: Minimum Effective Volume (typically 70-80% of average)
		mev[muscle] = math.Max(avgVolume*0.75, avgVolume-stdDev)

		// MAV: Maximum Adaptive Volume (typically average + 0.5 standard deviations)
		mav[muscle] = avgVolume + stdDev*0.5

		// MRV: Maximum Recoverable Volume (typically average + 1.5 standard deviations)
		mrv[muscle] = avgVolume + stdDev*1.5

		// Ensure MEV < MAV < MRV
		if mev[muscle] >= mav[muscle] {
			mav[muscle] = mev[muscle] * 1.25
		}
		if mav[muscle] >= mrv[muscle] {
			mrv[muscle] = mav[muscle] * 1.25
		}
	}

	return &models.VolumeLandmarks{
		MEV: mev,
		MAV: mav,
		MRV: mrv,
	}, nil
}

// UpdateUserMetrics recalculates and updates user metrics after a new workout
func (m *MetricsService) UpdateUserMetrics(ctx context.Context, userID primitive.ObjectID) error {
	// Get current user metrics
	userMetrics, err := m.getUserMetrics(ctx, userID, "all")
	if err != nil {
		return fmt.Errorf("failed to get current user metrics: %w", err)
	}

	// Update timestamp
	userMetrics.UpdatedAt = time.Now()

	// Upsert user metrics
	filter := bson.M{"userId": userID}
	update := bson.M{"$set": userMetrics}
	opts := options.Update().SetUpsert(true)

	_, err = m.userMetricsColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to update user metrics: %w", err)
	}

	return nil
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func standardDeviation(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	variance := 0.0
	for _, v := range values {
		variance += math.Pow(v-mean, 2)
	}
	return math.Sqrt(variance / float64(len(values)))
}

// gRPC Methods

// GetUserMetricsRPC retrieves user metrics for a specific period (gRPC method)
func (m *MetricsService) GetUserMetrics(ctx context.Context, req *pb.GetUserMetricsRequest) (*pb.UserMetrics, error) {
	if req.UserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	userObjectID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id")
	}

	userMetrics, err := m.getUserMetrics(ctx, userObjectID, req.Period)
	if err != nil {
		return nil, fmt.Errorf("failed to get user metrics: %w", err)
	}

	return m.userMetricsToProto(userMetrics), nil
}

// GetWorkoutMetrics retrieves metrics for a specific workout session
func (m *MetricsService) GetWorkoutMetrics(ctx context.Context, req *pb.GetWorkoutMetricsRequest) (*pb.WorkoutMetrics, error) {
	if req.SessionId == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	sessionObjectID, err := primitive.ObjectIDFromHex(req.SessionId)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id")
	}

	var workoutMetrics models.WorkoutMetrics
	err = m.metricsColl.FindOne(ctx, bson.M{"sessionId": sessionObjectID}).Decode(&workoutMetrics)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("workout metrics not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get workout metrics: %w", err)
	}

	return m.workoutMetricsToProto(&workoutMetrics), nil
}

// GetVolumeTrends retrieves volume trends for a user
func (m *MetricsService) GetVolumeTrends(ctx context.Context, req *pb.GetVolumeTrendsRequest) (*pb.TrendMetrics, error) {
	if req.UserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	userObjectID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id")
	}

	trendMetrics, err := m.CalculateVolumeTrends(ctx, userObjectID, req.Period)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate volume trends: %w", err)
	}

	return m.trendMetricsToProto(trendMetrics), nil
}

// Helper methods to convert models to protobuf

func (m *MetricsService) userMetricsToProto(userMetrics *models.UserMetrics) *pb.UserMetrics {
	return &pb.UserMetrics{
		Id:              userMetrics.ID.Hex(),
		UserId:          userMetrics.UserID.Hex(),
		WeeklyMetrics:   m.periodMetricsToProto(&userMetrics.WeeklyMetrics),
		MonthlyMetrics:  m.periodMetricsToProto(&userMetrics.MonthlyMetrics),
		AllTimeMetrics:  m.periodMetricsToProto(&userMetrics.AllTimeMetrics),
		VolumeLandmarks: m.volumeLandmarksToProto(&userMetrics.VolumeLandmarks),
		CreatedAt:       timestamppb.New(userMetrics.CreatedAt),
		UpdatedAt:       timestamppb.New(userMetrics.UpdatedAt),
	}
}

func (m *MetricsService) periodMetricsToProto(periodMetrics *models.PeriodMetrics) *pb.PeriodMetrics {
	var startDate, endDate *timestamppb.Timestamp
	if !periodMetrics.StartDate.IsZero() {
		startDate = timestamppb.New(periodMetrics.StartDate)
	}
	if !periodMetrics.EndDate.IsZero() {
		endDate = timestamppb.New(periodMetrics.EndDate)
	}

	return &pb.PeriodMetrics{
		StartDate:          startDate,
		EndDate:            endDate,
		TotalWorkouts:      periodMetrics.TotalWorkouts,
		TotalVolumeLoad:    periodMetrics.TotalVolumeLoad,
		TotalTonnage:       periodMetrics.TotalTonnage,
		AverageWorkoutTime: periodMetrics.AverageWorkoutTime,
		MuscleGroupVolume:  periodMetrics.MuscleGroupVolume,
		VolumeProgression:  periodMetrics.VolumeProgression,
	}
}

func (m *MetricsService) volumeLandmarksToProto(landmarks *models.VolumeLandmarks) *pb.VolumeLandmarks {
	return &pb.VolumeLandmarks{
		Mev: landmarks.MEV,
		Mav: landmarks.MAV,
		Mrv: landmarks.MRV,
	}
}

func (m *MetricsService) workoutMetricsToProto(workoutMetrics *models.WorkoutMetrics) *pb.WorkoutMetrics {
	var exerciseMetrics []*pb.ExerciseMetrics
	for _, em := range workoutMetrics.ExerciseMetrics {
		exerciseMetrics = append(exerciseMetrics, &pb.ExerciseMetrics{
			ExerciseId:       em.ExerciseID.Hex(),
			ExerciseName:     em.ExerciseName,
			VolumeLoad:       em.VolumeLoad,
			Tonnage:          em.Tonnage,
			SetCount:         em.SetCount,
			TotalReps:        em.TotalReps,
			AverageWeight:    em.AverageWeight,
			MaxWeight:        em.MaxWeight,
			PrimaryMuscles:   em.PrimaryMuscles,
			SecondaryMuscles: em.SecondaryMuscles,
		})
	}

	return &pb.WorkoutMetrics{
		Id:                          workoutMetrics.ID.Hex(),
		UserId:                      workoutMetrics.UserID.Hex(),
		SessionId:                   workoutMetrics.SessionID.Hex(),
		RoutineId:                   workoutMetrics.RoutineID.Hex(),
		Date:                        timestamppb.New(workoutMetrics.Date),
		VolumeMetrics:               m.volumeMetricsToProto(&workoutMetrics.VolumeMetrics),
		PerformanceMetrics:          m.performanceMetricsToProto(&workoutMetrics.PerformanceMetrics),
		IntensityMetrics:            m.intensityMetricsToProto(&workoutMetrics.IntensityMetrics),
		StrengthMetrics:             m.strengthMetricsToProto(&workoutMetrics.StrengthMetrics),
		ProgressAdaptationMetrics:   m.progressAdaptationMetricsToProto(&workoutMetrics.ProgressAdaptationMetrics),
		RecoveryFatigueMetrics:      m.recoveryFatigueMetricsToProto(&workoutMetrics.RecoveryFatigueMetrics),
		BodyCompositionMetrics:      m.bodyCompositionMetricsToProto(&workoutMetrics.BodyCompositionMetrics),
		MuscleSpecificMetrics:       m.muscleSpecificMetricsToProto(&workoutMetrics.MuscleSpecificMetrics),
		WorkCapacityMetrics:         m.workCapacityMetricsToProto(&workoutMetrics.WorkCapacityMetrics),
		TrainingPatternMetrics:      m.trainingPatternMetricsToProto(&workoutMetrics.TrainingPatternMetrics),
		PeriodizationMetrics:        m.periodizationMetricsToProto(&workoutMetrics.PeriodizationMetrics),
		InjuryRiskPreventionMetrics: m.injuryRiskPreventionMetricsToProto(&workoutMetrics.InjuryRiskPreventionMetrics),
		EfficiencyTechniqueMetrics:  m.efficiencyTechniqueMetricsToProto(&workoutMetrics.EfficiencyTechniqueMetrics),
		SetMetrics:                  m.setMetricsToProto(&workoutMetrics.SetMetrics),
		ExerciseMetrics:             exerciseMetrics,
		WorkoutDurationSecs:         workoutMetrics.WorkoutDurationSecs,
		CreatedAt:                   timestamppb.New(workoutMetrics.CreatedAt),
		UpdatedAt:                   timestamppb.New(workoutMetrics.UpdatedAt),
	}
}

func (m *MetricsService) volumeMetricsToProto(vm *models.VolumeMetrics) *pb.VolumeMetrics {
	return &pb.VolumeMetrics{
		TotalVolumeLoad:   vm.TotalVolumeLoad,
		Tonnage:           vm.Tonnage,
		RelativeVolume:    vm.RelativeVolume,
		EffectiveReps:     vm.EffectiveReps,
		HardSets:          vm.HardSets,
		MuscleGroupVolume: vm.MuscleGroupVolume,
	}
}

func (m *MetricsService) performanceMetricsToProto(pm *models.PerformanceMetrics) *pb.PerformanceMetrics {
	return &pb.PerformanceMetrics{
		AverageRpe:      pm.AverageRPE,
		MaxRpe:          pm.MaxRPE,
		RpeDistribution: pm.RPEDistribution,
	}
}

func (m *MetricsService) intensityMetricsToProto(im *models.IntensityMetrics) *pb.IntensityMetrics {
	return &pb.IntensityMetrics{
		AverageIntensity:      im.AverageIntensity,
		RelativeIntensity:     im.RelativeIntensity,
		RpeAdjustedLoad:       im.RPEAdjustedLoad,
		IntensityDistribution: im.IntensityDistribution,
		LoadDensity:           im.LoadDensity,
	}
}

func (m *MetricsService) setMetricsToProto(sm *models.SetMetrics) *pb.SetMetrics {
	return &pb.SetMetrics{
		TotalSets:      sm.TotalSets,
		CompletedSets:  sm.CompletedSets,
		CompletionRate: sm.CompletionRate,
	}
}

func (m *MetricsService) strengthMetricsToProto(sm *models.StrengthMetrics) *pb.StrengthMetrics {
	return &pb.StrengthMetrics{
		EstimatedOneRmEpley:   sm.EstimatedOneRMEpley,
		EstimatedOneRmBrzycki: sm.EstimatedOneRMBrzycki,
		WilksScore:            sm.WilksScore,
		PushPullRatio:         sm.PushPullRatio,
		PowerOutput:           sm.PowerOutput,
	}
}

func (m *MetricsService) trendMetricsToProto(tm *models.TrendMetrics) *pb.TrendMetrics {
	var volumeProgression []*pb.VolumeProgressionPoint
	for _, vp := range tm.VolumeProgression {
		volumeProgression = append(volumeProgression, &pb.VolumeProgressionPoint{
			Date:       timestamppb.New(vp.Date),
			VolumeLoad: vp.VolumeLoad,
			Tonnage:    vp.Tonnage,
		})
	}

	muscleGroupTrends := make(map[string]*pb.DoubleList)
	for muscle, trends := range tm.MuscleGroupTrends {
		muscleGroupTrends[muscle] = &pb.DoubleList{Values: trends}
	}

	return &pb.TrendMetrics{
		Period:            tm.Period,
		VolumeProgression: volumeProgression,
		VolumeGrowthRate:  tm.VolumeGrowthRate,
		MuscleGroupTrends: muscleGroupTrends,
		PerformanceTrends: tm.PerformanceTrends,
	}
}

func (m *MetricsService) progressAdaptationMetricsToProto(pam *models.ProgressAdaptationMetrics) *pb.ProgressAdaptationMetrics {
	plateauDetection := make(map[string]*pb.PlateauStatus)
	for exerciseID, status := range pam.PlateauDetection {
		plateauDetection[exerciseID] = &pb.PlateauStatus{
			IsPlateaued:        status.IsPlateaued,
			ConsecutiveWeeks:   status.ConsecutiveWeeks,
			LastProgressRate:   status.LastProgressRate,
			WeeksSinceProgress: status.WeeksSinceProgress,
		}
	}

	return &pb.ProgressAdaptationMetrics{
		ProgressiveOverloadIndex: pam.ProgressiveOverloadIndex,
		WeekOverWeekProgressRate: pam.WeekOverWeekProgressRate,
		PlateauDetection:         plateauDetection,
		StrengthGainVelocity:     pam.StrengthGainVelocity,
		AdaptationRate:           pam.AdaptationRate,
	}
}

func (m *MetricsService) recoveryFatigueMetricsToProto(rfm *models.RecoveryFatigueMetrics) *pb.RecoveryFatigueMetrics {
	return &pb.RecoveryFatigueMetrics{
		AcuteChronicWorkloadRatio: rfm.AcuteChronicWorkloadRatio,
		TrainingStrainScore:       rfm.TrainingStrainScore,
		RecoveryNeedIndex:         rfm.RecoveryNeedIndex,
		FatigueAccumulationIndex:  rfm.FatigueAccumulationIndex,
		OvertrainingRiskScore:     rfm.OvertrainingRiskScore,
		RecoveryScore:             rfm.RecoveryScore,
	}
}

func (m *MetricsService) bodyCompositionMetricsToProto(bcm *models.BodyCompositionMetrics) *pb.BodyCompositionMetrics {
	return &pb.BodyCompositionMetrics{
		Bmi:                    bcm.BMI,
		StrengthToWeightRatio:  bcm.StrengthToWeightRatio,
		AllometricScalingScore: bcm.AllometricScalingScore,
		PonderalIndex:          bcm.PonderalIndex,
	}
}

func (m *MetricsService) muscleSpecificMetricsToProto(msm *models.MuscleSpecificMetrics) *pb.MuscleSpecificMetrics {
	return &pb.MuscleSpecificMetrics{
		MuscleGroupDistribution: msm.MuscleGroupDistribution,
		MuscleImbalanceIndex:    msm.MuscleImbalanceIndex,
		AntagonistRatio:         msm.AntagonistRatio,
		StimulusToFatigueRatio:  msm.StimulusToFatigueRatio,
	}
}

func (m *MetricsService) workCapacityMetricsToProto(wcm *models.WorkCapacityMetrics) *pb.WorkCapacityMetrics {
	return &pb.WorkCapacityMetrics{
		TotalWorkCapacity:      wcm.TotalWorkCapacity,
		DensityTrainingIndex:   wcm.DensityTrainingIndex,
		DensityProgressPercent: wcm.DensityProgressPercent,
		TimeUnderTension:       wcm.TimeUnderTension,
		MechanicalTensionScore: wcm.MechanicalTensionScore,
	}
}

func (m *MetricsService) trainingPatternMetricsToProto(tpm *models.TrainingPatternMetrics) *pb.TrainingPatternMetrics {
	return &pb.TrainingPatternMetrics{
		OptimalFrequency:       tpm.OptimalFrequency,
		RecoveryTime:           tpm.RecoveryTime,
		ExerciseSelectionIndex: tpm.ExerciseSelectionIndex,
		ConsistencyScore:       tpm.ConsistencyScore,
		WorkoutCompletionRate:  tpm.WorkoutCompletionRate,
	}
}

func (m *MetricsService) periodizationMetricsToProto(pm *models.PeriodizationMetrics) *pb.PeriodizationMetrics {
	return &pb.PeriodizationMetrics{
		ChronicTrainingLoad:   pm.ChronicTrainingLoad,
		AcuteTrainingLoad:     pm.AcuteTrainingLoad,
		TrainingStressBalance: pm.TrainingStressBalance,
		FormFreshnessIndex:    pm.FormFreshnessIndex,
	}
}

// calculateInjuryRiskPreventionMetrics calculates injury risk and prevention metrics
func (m *MetricsService) calculateInjuryRiskPreventionMetrics(ctx context.Context, session *models.WorkoutSession) (models.InjuryRiskPreventionMetrics, error) {
	injuryRiskMetrics := models.InjuryRiskPreventionMetrics{}

	// Get recovery fatigue metrics for ACWR
	recoveryFatigueMetrics, err := m.calculateRecoveryFatigueMetrics(ctx, session)
	if err != nil {
		return injuryRiskMetrics, err
	}

	// Get muscle specific metrics for imbalance index
	muscleSpecificMetrics, err := m.calculateMuscleSpecificMetrics(ctx, session)
	if err != nil {
		return injuryRiskMetrics, err
	}

	// Calculate average imbalance index across all muscle groups
	var totalImbalance float64
	var imbalanceCount int
	for _, imbalance := range muscleSpecificMetrics.MuscleImbalanceIndex {
		totalImbalance += imbalance
		imbalanceCount++
	}

	var averageImbalance float64
	if imbalanceCount > 0 {
		averageImbalance = totalImbalance / float64(imbalanceCount)
	}

	// Calculate fatigue score (using fatigue accumulation index as proxy)
	fatigueScore := recoveryFatigueMetrics.FatigueAccumulationIndex

	// Get recovery quality (simplified using defaults)
	profile := models.DefaultProfile
	recoveryQuality := (profile.SleepQuality + profile.NutritionFactor) / 2.0

	// Calculate Injury Risk Score: IRS = (ACWR × Imbalance Index × Fatigue Score) / Recovery Quality
	acwr := recoveryFatigueMetrics.AcuteChronicWorkloadRatio
	if recoveryQuality > 0 {
		injuryRiskScore := (acwr * averageImbalance * math.Abs(fatigueScore)) / recoveryQuality
		injuryRiskMetrics.InjuryRiskScore = injuryRiskScore
	}

	// Calculate Load Spike Alert: Weekly Volume > 1.5 × Average of Last 4 Weeks
	currentWeekVolume, err := m.getCurrentWeekVolume(ctx, session.UserID, session.StartedAt)
	if err == nil {
		averageFourWeekVolume, err := m.getAverageFourWeekVolume(ctx, session.UserID, session.StartedAt)
		if err == nil && averageFourWeekVolume > 0 {
			spikeThreshold := averageFourWeekVolume * models.LoadSpikeMultiplier
			injuryRiskMetrics.LoadSpikeAlert = currentWeekVolume > spikeThreshold
		}
	}

	// Calculate Asymmetry Development: |(Left Performance - Right Performance)| / Average × 100
	asymmetryDevelopment, err := m.calculateAsymmetryDevelopment(ctx, session)
	if err == nil {
		injuryRiskMetrics.AsymmetryDevelopment = asymmetryDevelopment
	}

	return injuryRiskMetrics, nil
}

// calculateEfficiencyTechniqueMetrics calculates efficiency and technique metrics
func (m *MetricsService) calculateEfficiencyTechniqueMetrics(ctx context.Context, session *models.WorkoutSession) (models.EfficiencyTechniqueMetrics, error) {
	efficiencyMetrics := models.EfficiencyTechniqueMetrics{}

	// Calculate volume metrics for total volume
	volumeMetrics, err := m.calculateVolumeMetrics(ctx, session)
	if err != nil {
		return efficiencyMetrics, err
	}

	// Calculate Strength Efficiency: SE = (1RM Gain / Total Volume) × 1000
	strengthEfficiency, err := m.calculateStrengthEfficiency(ctx, session, volumeMetrics.TotalVolumeLoad)
	if err == nil {
		efficiencyMetrics.StrengthEfficiency = strengthEfficiency
	}

	// Calculate Volume Efficiency: VE = Performance Improvement / Total Training Volume
	volumeEfficiency, err := m.calculateVolumeEfficiency(ctx, session, volumeMetrics.TotalVolumeLoad)
	if err == nil {
		efficiencyMetrics.VolumeEfficiency = volumeEfficiency
	}

	// Calculate RPE-Performance Correlation: Pearson(Actual Reps, 10 - RPE + Expected Reps at RPE)
	rpeCorrelation, err := m.calculateRPEPerformanceCorrelation(ctx, session)
	if err == nil {
		efficiencyMetrics.RPEPerformanceCorrelation = rpeCorrelation
	}

	// Calculate Technique Consistency: 1 - (Standard Deviation of Rep Times / Average Rep Time)
	techniqueConsistency, err := m.calculateTechniqueConsistency(ctx, session)
	if err == nil {
		efficiencyMetrics.TechniqueConsistency = techniqueConsistency
	}

	return efficiencyMetrics, nil
}

// Helper methods for Injury Risk Prevention Metrics

func (m *MetricsService) getCurrentWeekVolume(ctx context.Context, userID primitive.ObjectID, currentDate time.Time) (float64, error) {
	weekStart := currentDate.AddDate(0, 0, -int(currentDate.Weekday()))
	periodMetrics, err := m.calculatePeriodMetrics(ctx, userID, weekStart, currentDate)
	if err != nil {
		return 0, err
	}
	return periodMetrics.TotalVolumeLoad, nil
}

func (m *MetricsService) getAverageFourWeekVolume(ctx context.Context, userID primitive.ObjectID, currentDate time.Time) (float64, error) {
	// Get volume for last 4 weeks (excluding current week)
	currentWeekStart := currentDate.AddDate(0, 0, -int(currentDate.Weekday()))

	var weeklyVolumes []float64

	for i := 0; i < 4; i++ {
		weekStart := currentWeekStart.AddDate(0, 0, -(i+1)*7)
		weekEnd := weekStart.AddDate(0, 0, 6)

		periodMetrics, err := m.calculatePeriodMetrics(ctx, userID, weekStart, weekEnd)
		if err == nil {
			weeklyVolumes = append(weeklyVolumes, periodMetrics.TotalVolumeLoad)
		}
	}

	if len(weeklyVolumes) == 0 {
		return 0, fmt.Errorf("no historical volume data found")
	}

	// Calculate average
	total := 0.0
	for _, volume := range weeklyVolumes {
		total += volume
	}

	return total / float64(len(weeklyVolumes)), nil
}

func (m *MetricsService) calculateAsymmetryDevelopment(ctx context.Context, session *models.WorkoutSession) (float64, error) {
	// For now, use the largest imbalance found in muscle-specific metrics as proxy
	muscleSpecificMetrics, err := m.calculateMuscleSpecificMetrics(ctx, session)
	if err != nil {
		return 0, err
	}

	var maxAsymmetry float64
	for _, imbalance := range muscleSpecificMetrics.MuscleImbalanceIndex {
		if imbalance > maxAsymmetry {
			maxAsymmetry = imbalance
		}
	}

	return maxAsymmetry, nil
}

// Helper methods for Efficiency & Technique Metrics

func (m *MetricsService) calculateStrengthEfficiency(ctx context.Context, session *models.WorkoutSession, totalVolume float64) (float64, error) {
	// Get current session's best 1RM estimates
	strengthMetrics, err := m.calculateStrengthMetrics(ctx, session)
	if err != nil {
		return 0, err
	}

	// Get previous session's 1RM estimates for comparison
	previousMetrics, err := m.getPreviousSessionMetrics(ctx, session.UserID, session.StartedAt)
	if err != nil || previousMetrics == nil {
		return 0, nil // No previous data for comparison
	}

	var totalOneRMGain float64
	var exerciseCount int

	// Calculate 1RM gains for each exercise
	for exerciseID, currentOneRM := range strengthMetrics.EstimatedOneRMEpley {
		if previousOneRM, exists := previousMetrics.StrengthMetrics.EstimatedOneRMEpley[exerciseID]; exists {
			oneRMGain := currentOneRM - previousOneRM
			if oneRMGain > 0 { // Only count positive gains
				totalOneRMGain += oneRMGain
				exerciseCount++
			}
		}
	}

	if totalVolume <= 0 || totalOneRMGain <= 0 {
		return 0, nil
	}

	// SE = (1RM Gain / Total Volume) × 1000
	strengthEfficiency := (totalOneRMGain / totalVolume) * models.StrengthEfficiencyMultiplier
	return strengthEfficiency, nil
}

func (m *MetricsService) calculateVolumeEfficiency(ctx context.Context, session *models.WorkoutSession, totalVolume float64) (float64, error) {
	// Calculate performance improvement using average 1RM improvement
	strengthMetrics, err := m.calculateStrengthMetrics(ctx, session)
	if err != nil {
		return 0, err
	}

	previousMetrics, err := m.getPreviousSessionMetrics(ctx, session.UserID, session.StartedAt)
	if err != nil || previousMetrics == nil {
		return 0, nil // No previous data for comparison
	}

	var totalPerformanceImprovement float64
	var exerciseCount int

	// Calculate performance improvement for each exercise
	for exerciseID, currentOneRM := range strengthMetrics.EstimatedOneRMEpley {
		if previousOneRM, exists := previousMetrics.StrengthMetrics.EstimatedOneRMEpley[exerciseID]; exists && previousOneRM > 0 {
			performanceImprovement := (currentOneRM - previousOneRM) / previousOneRM
			totalPerformanceImprovement += performanceImprovement
			exerciseCount++
		}
	}

	if exerciseCount == 0 || totalVolume <= 0 {
		return 0, nil
	}

	averagePerformanceImprovement := totalPerformanceImprovement / float64(exerciseCount)

	// VE = Performance Improvement / Total Training Volume
	volumeEfficiency := averagePerformanceImprovement / totalVolume
	return volumeEfficiency, nil
}

func (m *MetricsService) calculateRPEPerformanceCorrelation(ctx context.Context, session *models.WorkoutSession) (float64, error) {
	var actualReps []float64
	var expectedReps []float64

	// Collect data points from all completed sets
	for _, exercise := range session.Exercises {
		for _, set := range exercise.Sets {
			if !set.Completed {
				continue
			}

			actualRep := float64(set.ActualReps)
			rpe := float64(session.RPERating)

			// Expected reps at given RPE (simplified formula: higher RPE = fewer reps expected)
			// This is a simplified model: Expected = 10 - RPE + baseline reps
			expectedRep := 10.0 - rpe + float64(set.TargetReps)

			actualReps = append(actualReps, actualRep)
			expectedReps = append(expectedReps, expectedRep)
		}
	}

	// Need minimum sample size for meaningful correlation
	if len(actualReps) < models.MinCorrelationSampleSize {
		return 0, nil
	}

	// Calculate Pearson correlation coefficient
	correlation := calculatePearsonCorrelation(actualReps, expectedReps)
	return correlation, nil
}

func (m *MetricsService) calculateTechniqueConsistency(ctx context.Context, session *models.WorkoutSession) (float64, error) {
	var repTimes []float64

	// For now, use default rep times based on exercise tempo
	// In a real implementation, this would use actual timing data from the workout
	for _, exercise := range session.Exercises {
		exerciseDetails, err := m.exerciseService.GetExercise(ctx, &pb.GetExerciseRequest{Id: exercise.ExerciseID.Hex()})
		if err != nil {
			continue
		}

		// Get exercise tempo
		tempo := models.DefaultExerciseTempos["default"]
		exerciseName := strings.ToLower(exerciseDetails.Name)
		for knownExercise, knownTempo := range models.DefaultExerciseTempos {
			if strings.Contains(exerciseName, knownExercise) {
				tempo = knownTempo
				break
			}
		}

		repTime := tempo.Eccentric + tempo.Pause1 + tempo.Concentric + tempo.Pause2

		for _, set := range exercise.Sets {
			if set.Completed {
				// Add some simulated variance (±10%) to represent technique consistency
				// In real implementation, this would be actual measured rep times
				variance := 0.1 * repTime * (math.Sin(float64(set.ActualReps)) - 0.5) // Pseudo-random variance
				simulatedRepTime := repTime + variance
				repTimes = append(repTimes, simulatedRepTime)
			}
		}
	}

	if len(repTimes) < 2 {
		return 1.0, nil // Perfect consistency if only one rep
	}

	// Calculate average and standard deviation
	averageRepTime := average(repTimes)
	stdDev := standardDeviation(repTimes, averageRepTime)

	if averageRepTime <= 0 {
		return 0, nil
	}

	// Technique Consistency = 1 - (Standard Deviation / Average Rep Time)
	consistency := 1.0 - (stdDev / averageRepTime)

	// Ensure consistency is between 0 and 1
	if consistency < 0 {
		consistency = 0
	} else if consistency > 1 {
		consistency = 1
	}

	return consistency, nil
}

// calculatePearsonCorrelation calculates the Pearson correlation coefficient between two slices
func calculatePearsonCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}

	// Calculate means
	meanX := average(x)
	meanY := average(y)

	// Calculate correlation coefficient
	var numerator, denomX, denomY float64

	for i := range x {
		dx := x[i] - meanX
		dy := y[i] - meanY

		numerator += dx * dy
		denomX += dx * dx
		denomY += dy * dy
	}

	if denomX == 0 || denomY == 0 {
		return 0
	}

	correlation := numerator / math.Sqrt(denomX*denomY)
	return correlation
}

func (m *MetricsService) injuryRiskPreventionMetricsToProto(irpm *models.InjuryRiskPreventionMetrics) *pb.InjuryRiskPreventionMetrics {
	return &pb.InjuryRiskPreventionMetrics{
		InjuryRiskScore:      irpm.InjuryRiskScore,
		LoadSpikeAlert:       irpm.LoadSpikeAlert,
		AsymmetryDevelopment: irpm.AsymmetryDevelopment,
	}
}

func (m *MetricsService) efficiencyTechniqueMetricsToProto(etm *models.EfficiencyTechniqueMetrics) *pb.EfficiencyTechniqueMetrics {
	return &pb.EfficiencyTechniqueMetrics{
		StrengthEfficiency:        etm.StrengthEfficiency,
		VolumeEfficiency:          etm.VolumeEfficiency,
		RpePerformanceCorrelation: etm.RPEPerformanceCorrelation,
		TechniqueConsistency:      etm.TechniqueConsistency,
	}
}
