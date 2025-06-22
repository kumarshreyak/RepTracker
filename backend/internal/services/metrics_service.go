package services

import (
	"context"
	"fmt"
	"math"
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

	metrics.SetMetrics = m.calculateSetMetrics(session)

	metrics.ExerciseMetrics, err = m.calculateExerciseMetrics(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate exercise metrics: %w", err)
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
		Id:                  workoutMetrics.ID.Hex(),
		UserId:              workoutMetrics.UserID.Hex(),
		SessionId:           workoutMetrics.SessionID.Hex(),
		RoutineId:           workoutMetrics.RoutineID.Hex(),
		Date:                timestamppb.New(workoutMetrics.Date),
		VolumeMetrics:       m.volumeMetricsToProto(&workoutMetrics.VolumeMetrics),
		PerformanceMetrics:  m.performanceMetricsToProto(&workoutMetrics.PerformanceMetrics),
		IntensityMetrics:    m.intensityMetricsToProto(&workoutMetrics.IntensityMetrics),
		StrengthMetrics:     m.strengthMetricsToProto(&workoutMetrics.StrengthMetrics),
		SetMetrics:          m.setMetricsToProto(&workoutMetrics.SetMetrics),
		ExerciseMetrics:     exerciseMetrics,
		WorkoutDurationSecs: workoutMetrics.WorkoutDurationSecs,
		CreatedAt:           timestamppb.New(workoutMetrics.CreatedAt),
		UpdatedAt:           timestamppb.New(workoutMetrics.UpdatedAt),
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
