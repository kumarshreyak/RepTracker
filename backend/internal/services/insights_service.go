package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"gymlog-backend/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"
)

// InsightType represents different types of insights
type InsightType string

const (
	InsightTypeProgress   InsightType = "progress"
	InsightTypeVolume     InsightType = "volume"
	InsightTypeStrength   InsightType = "strength"
	InsightTypeRecovery   InsightType = "recovery"
	InsightTypeBalance    InsightType = "balance"
	InsightTypeEfficiency InsightType = "efficiency"
	InsightTypeMotivation InsightType = "motivation"
	InsightTypeRisk       InsightType = "risk"
)

// WorkoutInsight represents an AI-generated insight
type WorkoutInsight struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Type      InsightType        `bson:"type" json:"type"`
	Insight   string             `bson:"insight" json:"insight"`
	BasedOn   string             `bson:"basedOn" json:"basedOn"`   // e.g., "Last 7 days", "Previous workout"
	Priority  int                `bson:"priority" json:"priority"` // 1-5, 5 being most important
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// InsightsService handles AI-powered workout insights
type InsightsService struct {
	db          *mongo.Database
	genAIClient *genai.Client
}

// NewInsightsService creates a new insights service
func NewInsightsService(db *mongo.Database) (*InsightsService, error) {
	ctx := context.Background()

	// Get API key from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// Create Gemini client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &InsightsService{
		db:          db,
		genAIClient: client,
	}, nil
}

// GenerateWorkoutInsights generates insights for a user based on their recent metrics
func (s *InsightsService) GenerateWorkoutInsights(ctx context.Context, userID primitive.ObjectID) ([]WorkoutInsight, error) {
	// Fetch recent workout metrics
	workoutMetrics, err := s.fetchRecentWorkoutMetrics(ctx, userID, 7) // Last 7 days
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workout metrics: %w", err)
	}

	// Fetch user metrics
	userMetrics, err := s.fetchUserMetrics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user metrics: %w", err)
	}

	// Generate insights based on different aspects
	var insights []WorkoutInsight

	// Progress insight
	if insight, err := s.generateProgressInsight(ctx, workoutMetrics, userMetrics); err == nil && insight != nil {
		insights = append(insights, *insight)
	}

	// Volume insight
	if insight, err := s.generateVolumeInsight(ctx, workoutMetrics, userMetrics); err == nil && insight != nil {
		insights = append(insights, *insight)
	}

	// Recovery insight
	if insight, err := s.generateRecoveryInsight(ctx, workoutMetrics, userMetrics); err == nil && insight != nil {
		insights = append(insights, *insight)
	}

	// Balance insight
	if insight, err := s.generateBalanceInsight(ctx, workoutMetrics, userMetrics); err == nil && insight != nil {
		insights = append(insights, *insight)
	}

	// Risk insight
	if insight, err := s.generateRiskInsight(ctx, workoutMetrics, userMetrics); err == nil && insight != nil {
		insights = append(insights, *insight)
	}

	// Save insights to database
	if len(insights) > 0 {
		err = s.saveInsights(ctx, insights)
		if err != nil {
			// Log error but don't fail - insights are already generated
			fmt.Printf("Failed to save insights: %v\n", err)
		}
	}

	return insights, nil
}

// generateProgressInsight generates an insight about training progress
func (s *InsightsService) generateProgressInsight(ctx context.Context, workoutMetrics []models.WorkoutMetrics, userMetrics *models.UserMetrics) (*WorkoutInsight, error) {
	if len(workoutMetrics) == 0 {
		return nil, nil
	}

	// Prepare metrics summary for AI
	metricsData := s.prepareProgressMetricsData(workoutMetrics)

	prompt := fmt.Sprintf(`Based on these workout progress metrics from the last 7 days:
%s

Generate ONE important insight about the user's training progress in exactly one sentence. 
Focus on strength gains, progressive overload, or plateau detection.
Be specific with numbers when relevant.
Make it actionable or motivational.`, metricsData)

	insight, err := s.callGeminiAPI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return &WorkoutInsight{
		ID:        primitive.NewObjectID(),
		UserID:    workoutMetrics[0].UserID,
		Type:      InsightTypeProgress,
		Insight:   insight,
		BasedOn:   "Last 7 days",
		Priority:  s.calculatePriority(workoutMetrics, InsightTypeProgress),
		CreatedAt: time.Now(),
	}, nil
}

// generateVolumeInsight generates an insight about training volume
func (s *InsightsService) generateVolumeInsight(ctx context.Context, workoutMetrics []models.WorkoutMetrics, userMetrics *models.UserMetrics) (*WorkoutInsight, error) {
	if len(workoutMetrics) == 0 || userMetrics == nil {
		return nil, nil
	}

	// Prepare volume data with landmarks
	volumeData := s.prepareVolumeMetricsData(workoutMetrics, userMetrics)

	prompt := fmt.Sprintf(`Based on these training volume metrics and landmarks:
%s

Generate ONE important insight about the user's training volume in exactly one sentence.
Consider MEV/MAV/MRV landmarks and muscle group distribution.
Be specific and actionable.`, volumeData)

	insight, err := s.callGeminiAPI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return &WorkoutInsight{
		ID:        primitive.NewObjectID(),
		UserID:    workoutMetrics[0].UserID,
		Type:      InsightTypeVolume,
		Insight:   insight,
		BasedOn:   "Volume analysis",
		Priority:  s.calculatePriority(workoutMetrics, InsightTypeVolume),
		CreatedAt: time.Now(),
	}, nil
}

// generateRecoveryInsight generates an insight about recovery needs
func (s *InsightsService) generateRecoveryInsight(ctx context.Context, workoutMetrics []models.WorkoutMetrics, userMetrics *models.UserMetrics) (*WorkoutInsight, error) {
	if len(workoutMetrics) == 0 {
		return nil, nil
	}

	// Prepare recovery-related data
	recoveryData := s.prepareRecoveryMetricsData(workoutMetrics)

	prompt := fmt.Sprintf(`Based on these recovery and fatigue metrics:
%s

Generate ONE important insight about the user's recovery needs in exactly one sentence.
Consider training stress balance, fatigue indicators, and optimal frequency.
Be specific and actionable.`, recoveryData)

	insight, err := s.callGeminiAPI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return &WorkoutInsight{
		ID:        primitive.NewObjectID(),
		UserID:    workoutMetrics[0].UserID,
		Type:      InsightTypeRecovery,
		Insight:   insight,
		BasedOn:   "Recovery metrics",
		Priority:  s.calculatePriority(workoutMetrics, InsightTypeRecovery),
		CreatedAt: time.Now(),
	}, nil
}

// generateBalanceInsight generates an insight about muscle balance
func (s *InsightsService) generateBalanceInsight(ctx context.Context, workoutMetrics []models.WorkoutMetrics, userMetrics *models.UserMetrics) (*WorkoutInsight, error) {
	if len(workoutMetrics) == 0 {
		return nil, nil
	}

	// Prepare balance-related data
	balanceData := s.prepareBalanceMetricsData(workoutMetrics)

	prompt := fmt.Sprintf(`Based on these muscle balance and training distribution metrics:
%s

Generate ONE important insight about the user's training balance in exactly one sentence.
Consider push/pull ratio, muscle imbalances, and muscle group distribution.
Be specific and actionable.`, balanceData)

	insight, err := s.callGeminiAPI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return &WorkoutInsight{
		ID:        primitive.NewObjectID(),
		UserID:    workoutMetrics[0].UserID,
		Type:      InsightTypeBalance,
		Insight:   insight,
		BasedOn:   "Balance analysis",
		Priority:  s.calculatePriority(workoutMetrics, InsightTypeBalance),
		CreatedAt: time.Now(),
	}, nil
}

// generateRiskInsight generates an insight about injury risk
func (s *InsightsService) generateRiskInsight(ctx context.Context, workoutMetrics []models.WorkoutMetrics, userMetrics *models.UserMetrics) (*WorkoutInsight, error) {
	if len(workoutMetrics) == 0 {
		return nil, nil
	}

	// Prepare risk-related data
	riskData := s.prepareRiskMetricsData(workoutMetrics)

	prompt := fmt.Sprintf(`Based on these injury risk and prevention metrics:
%s

Generate ONE important insight about potential injury risks or prevention strategies in exactly one sentence.
Focus on load spikes, asymmetries, or injury risk scores.
Be cautious but constructive.`, riskData)

	insight, err := s.callGeminiAPI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return &WorkoutInsight{
		ID:        primitive.NewObjectID(),
		UserID:    workoutMetrics[0].UserID,
		Type:      InsightTypeRisk,
		Insight:   insight,
		BasedOn:   "Risk assessment",
		Priority:  s.calculatePriority(workoutMetrics, InsightTypeRisk),
		CreatedAt: time.Now(),
	}, nil
}

// callGeminiAPI makes a call to the Gemini API
func (s *InsightsService) callGeminiAPI(ctx context.Context, prompt string) (string, error) {
	// Add system context to ensure one-sentence responses
	fullPrompt := fmt.Sprintf(`You are a professional fitness coach analyzing workout data. 
	Talking to an absolute beginner in fitness, so keep it simple and easy to understand.
%s

Remember: Generate exactly ONE sentence only. Make it insightful, specific, and actionable.`, prompt)

	result, err := s.genAIClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(fullPrompt),
		nil, // Use default config
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract the text from the response
	insight := strings.TrimSpace(result.Text())

	// Ensure it's a single sentence
	sentences := strings.Split(insight, ". ")
	if len(sentences) > 1 {
		insight = sentences[0] + "."
	}

	return insight, nil
}

// fetchRecentWorkoutMetrics fetches recent workout metrics for a user
func (s *InsightsService) fetchRecentWorkoutMetrics(ctx context.Context, userID primitive.ObjectID, days int) ([]models.WorkoutMetrics, error) {
	collection := s.db.Collection("workout_metrics")

	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	filter := bson.M{
		"userId": userID,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var metrics []models.WorkoutMetrics
	if err = cursor.All(ctx, &metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}

// fetchUserMetrics fetches aggregated user metrics
func (s *InsightsService) fetchUserMetrics(ctx context.Context, userID primitive.ObjectID) (*models.UserMetrics, error) {
	collection := s.db.Collection("user_metrics")

	var userMetrics models.UserMetrics
	err := collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&userMetrics)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &userMetrics, nil
}

// Helper methods to prepare data for AI prompts

func (s *InsightsService) prepareProgressMetricsData(metrics []models.WorkoutMetrics) string {
	if len(metrics) == 0 {
		return "No metrics available"
	}

	latest := metrics[0]
	data := fmt.Sprintf(`Progressive Overload Index: %.2f
Week-over-Week Progress Rates: %s
Plateau Detection: %s
Average RPE: %.1f
Recent 1RM Estimates: %s`,
		latest.ProgressAdaptationMetrics.ProgressiveOverloadIndex,
		s.formatProgressRates(latest.ProgressAdaptationMetrics.WeekOverWeekProgressRate),
		s.formatPlateauStatus(latest.ProgressAdaptationMetrics.PlateauDetection),
		latest.PerformanceMetrics.AverageRPE,
		s.format1RMEstimates(latest.StrengthMetrics.EstimatedOneRMEpley))

	return data
}

func (s *InsightsService) prepareVolumeMetricsData(metrics []models.WorkoutMetrics, userMetrics *models.UserMetrics) string {
	if len(metrics) == 0 {
		return "No metrics available"
	}

	latest := metrics[0]
	data := fmt.Sprintf(`Total Volume Load: %.0f
Volume per Muscle Group: %s
Volume Landmarks (MEV/MAV/MRV): %s
Hard Sets: %d
Effective Reps: %d`,
		latest.VolumeMetrics.TotalVolumeLoad,
		s.formatMuscleGroupVolume(latest.VolumeMetrics.MuscleGroupVolume),
		s.formatVolumeLandmarks(userMetrics),
		latest.VolumeMetrics.HardSets,
		latest.VolumeMetrics.EffectiveReps)

	return data
}

func (s *InsightsService) prepareRecoveryMetricsData(metrics []models.WorkoutMetrics) string {
	if len(metrics) == 0 {
		return "No metrics available"
	}

	latest := metrics[0]
	data := fmt.Sprintf(`Training Stress Balance: %.1f
Form/Freshness Index: %.1f%%
Recovery Time Needed: %.1f hours
Optimal Training Frequency: %.1f hours
Burnout Risk Score: %.2f`,
		latest.PeriodizationMetrics.TrainingStressBalance,
		latest.PeriodizationMetrics.FormFreshnessIndex,
		latest.TrainingPatternMetrics.RecoveryTime,
		latest.TrainingPatternMetrics.OptimalFrequency,
		0.0) // TODO: Add burnout risk when available

	return data
}

func (s *InsightsService) prepareBalanceMetricsData(metrics []models.WorkoutMetrics) string {
	if len(metrics) == 0 {
		return "No metrics available"
	}

	latest := metrics[0]
	data := fmt.Sprintf(`Push:Pull Ratio: %.2f
Muscle Group Distribution: %s
Muscle Imbalances: %s
Antagonist Ratios: %s`,
		latest.StrengthMetrics.PushPullRatio,
		s.formatMuscleDistribution(latest.MuscleSpecificMetrics.MuscleGroupDistribution),
		s.formatImbalances(latest.MuscleSpecificMetrics.MuscleImbalanceIndex),
		s.formatAntagonistRatios(latest.MuscleSpecificMetrics.AntagonistRatio))

	return data
}

func (s *InsightsService) prepareRiskMetricsData(metrics []models.WorkoutMetrics) string {
	if len(metrics) == 0 {
		return "No metrics available"
	}

	latest := metrics[0]
	data := fmt.Sprintf(`Injury Risk Score: %.2f (threshold: 1.5)
Load Spike Alert: %v
Asymmetry Development: %.1f%% (threshold: 15%%)
Training Consistency: %.1f%%`,
		latest.InjuryRiskPreventionMetrics.InjuryRiskScore,
		latest.InjuryRiskPreventionMetrics.LoadSpikeAlert,
		latest.InjuryRiskPreventionMetrics.AsymmetryDevelopment,
		latest.TrainingPatternMetrics.ConsistencyScore)

	return data
}

// Formatting helper methods

func (s *InsightsService) formatProgressRates(rates map[string]float64) string {
	if len(rates) == 0 {
		return "No data"
	}

	parts := []string{}
	for exercise, rate := range rates {
		if len(parts) < 3 { // Limit to top 3
			parts = append(parts, fmt.Sprintf("%s: %.1f%%", exercise, rate))
		}
	}
	return strings.Join(parts, ", ")
}

func (s *InsightsService) formatPlateauStatus(plateaus map[string]models.PlateauStatus) string {
	plateauCount := 0
	for _, status := range plateaus {
		if status.IsPlateaued {
			plateauCount++
		}
	}

	if plateauCount == 0 {
		return "No plateaus detected"
	}
	return fmt.Sprintf("%d exercises plateaued", plateauCount)
}

func (s *InsightsService) format1RMEstimates(estimates map[string]float64) string {
	if len(estimates) == 0 {
		return "No data"
	}

	parts := []string{}
	for exercise, estimate := range estimates {
		if len(parts) < 3 {
			parts = append(parts, fmt.Sprintf("%s: %.1fkg", exercise, estimate))
		}
	}
	return strings.Join(parts, ", ")
}

func (s *InsightsService) formatMuscleGroupVolume(volume map[string]float64) string {
	if len(volume) == 0 {
		return "No data"
	}

	parts := []string{}
	for muscle, vol := range volume {
		if len(parts) < 3 && vol > 0 {
			parts = append(parts, fmt.Sprintf("%s: %.0f", muscle, vol))
		}
	}
	return strings.Join(parts, ", ")
}

func (s *InsightsService) formatVolumeLandmarks(userMetrics *models.UserMetrics) string {
	if userMetrics == nil {
		return "No landmarks available"
	}

	parts := []string{}
	for muscle, mev := range userMetrics.VolumeLandmarks.MEV {
		if mav, ok := userMetrics.VolumeLandmarks.MAV[muscle]; ok {
			if mrv, ok := userMetrics.VolumeLandmarks.MRV[muscle]; ok {
				parts = append(parts, fmt.Sprintf("%s: MEV=%.0f, MAV=%.0f, MRV=%.0f", muscle, mev, mav, mrv))
				if len(parts) >= 2 {
					break
				}
			}
		}
	}
	return strings.Join(parts, "; ")
}

func (s *InsightsService) formatMuscleDistribution(dist map[string]float64) string {
	if len(dist) == 0 {
		return "No data"
	}

	parts := []string{}
	for muscle, percentage := range dist {
		if len(parts) < 3 && percentage > 5 { // Only show significant muscle groups
			parts = append(parts, fmt.Sprintf("%s: %.0f%%", muscle, percentage))
		}
	}
	return strings.Join(parts, ", ")
}

func (s *InsightsService) formatImbalances(imbalances map[string]float64) string {
	if len(imbalances) == 0 {
		return "No imbalances"
	}

	significant := []string{}
	for muscle, imbalance := range imbalances {
		if imbalance > 10 { // Only show significant imbalances
			significant = append(significant, fmt.Sprintf("%s: %.0f%%", muscle, imbalance))
		}
	}

	if len(significant) == 0 {
		return "All balanced"
	}
	return strings.Join(significant, ", ")
}

func (s *InsightsService) formatAntagonistRatios(ratios map[string]float64) string {
	if len(ratios) == 0 {
		return "No data"
	}

	parts := []string{}
	for pairing, ratio := range ratios {
		if len(parts) < 2 {
			parts = append(parts, fmt.Sprintf("%s: %.2f", pairing, ratio))
		}
	}
	return strings.Join(parts, ", ")
}

// calculatePriority calculates the priority of an insight based on metrics
func (s *InsightsService) calculatePriority(metrics []models.WorkoutMetrics, insightType InsightType) int {
	if len(metrics) == 0 {
		return 3 // Default medium priority
	}

	latest := metrics[0]
	priority := 3 // Default

	switch insightType {
	case InsightTypeRisk:
		if latest.InjuryRiskPreventionMetrics.InjuryRiskScore > 1.5 ||
			latest.InjuryRiskPreventionMetrics.LoadSpikeAlert ||
			latest.InjuryRiskPreventionMetrics.AsymmetryDevelopment > 15 {
			priority = 5 // High priority for risks
		}
	case InsightTypeProgress:
		if latest.ProgressAdaptationMetrics.ProgressiveOverloadIndex < 0.95 {
			priority = 4 // Higher priority if not progressing
		}
	case InsightTypeRecovery:
		if latest.PeriodizationMetrics.TrainingStressBalance < -10 {
			priority = 4 // Higher priority if fatigued
		}
	case InsightTypeBalance:
		if latest.StrengthMetrics.PushPullRatio > 1.5 || latest.StrengthMetrics.PushPullRatio < 0.67 {
			priority = 4 // Higher priority for imbalances
		}
	}

	return priority
}

// saveInsights saves insights to the database
func (s *InsightsService) saveInsights(ctx context.Context, insights []WorkoutInsight) error {
	collection := s.db.Collection("workout_insights")

	// Convert to interface slice for InsertMany
	docs := make([]interface{}, len(insights))
	for i, insight := range insights {
		docs[i] = insight
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}

// GetRecentInsights retrieves recent insights for a user
func (s *InsightsService) GetRecentInsights(ctx context.Context, userID primitive.ObjectID, limit int) ([]WorkoutInsight, error) {
	collection := s.db.Collection("workout_insights")

	filter := bson.M{"userId": userID}
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var insights []WorkoutInsight
	if err = cursor.All(ctx, &insights); err != nil {
		return nil, err
	}

	return insights, nil
}
