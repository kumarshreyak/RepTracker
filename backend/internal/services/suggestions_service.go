package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gymlog-backend/pkg/models"
	pb "gymlog-backend/proto/gymlog/v1"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SuggestionsService handles AI-powered workout suggestions
type SuggestionsService struct {
	pb.UnimplementedSuggestionsServiceServer
	db              *mongo.Database
	genAIClient     *genai.Client
	workoutService  *WorkoutService
	userService     *UserService
	exerciseService *ExerciseService
	insightsService *InsightsService
}

// NewSuggestionsService creates a new suggestions service
func NewSuggestionsService(db *mongo.Database, workoutService *WorkoutService, userService *UserService, exerciseService *ExerciseService, insightsService *InsightsService) (*SuggestionsService, error) {
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

	return &SuggestionsService{
		db:              db,
		genAIClient:     client,
		workoutService:  workoutService,
		userService:     userService,
		exerciseService: exerciseService,
		insightsService: insightsService,
	}, nil
}

// GenerateWorkoutSuggestions generates AI-powered workout suggestions based on recent sessions
func (s *SuggestionsService) GenerateWorkoutSuggestions(ctx context.Context, req *pb.GenerateWorkoutSuggestionsRequest) (*pb.GenerateWorkoutSuggestionsResponse, error) {
	log.Printf("GenerateWorkoutSuggestions: Starting for user %s", req.UserId)

	// Validate user ID
	userObjectID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	// Set defaults
	daysToAnalyze := req.DaysToAnalyze
	if daysToAnalyze <= 0 {
		daysToAnalyze = 14 // Default to 2 weeks
	}
	maxSuggestions := req.MaxSuggestions
	if maxSuggestions <= 0 {
		maxSuggestions = 3 // Default to 3 suggestions
	}

	// Fetch user data and goal
	log.Printf("GenerateWorkoutSuggestions: Fetching user data for %s", req.UserId)
	userReq := &pb.GetUserRequest{Id: req.UserId}
	user, err := s.userService.GetUser(ctx, userReq)
	if err != nil {
		log.Printf("GenerateWorkoutSuggestions: Failed to fetch user: %v", err)
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}
	log.Printf("GenerateWorkoutSuggestions: User goal: %s", user.Goal)

	// Fetch recent workout sessions
	log.Printf("GenerateWorkoutSuggestions: Fetching workout sessions from last %d days", daysToAnalyze)
	sessions, err := s.fetchRecentWorkoutSessions(ctx, userObjectID, int(daysToAnalyze))
	if err != nil {
		log.Printf("GenerateWorkoutSuggestions: Failed to fetch workout sessions: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to fetch workout sessions: %v", err)
	}
	log.Printf("GenerateWorkoutSuggestions: Found %d workout sessions", len(sessions))

	if len(sessions) == 0 {
		return &pb.GenerateWorkoutSuggestionsResponse{
			Suggestions:     []*pb.SuggestedWorkout{},
			AnalysisSummary: "No recent workout sessions found to analyze.",
		}, nil
	}

	// Fetch user's routines/workouts
	log.Printf("GenerateWorkoutSuggestions: Fetching user's routines")
	routines, err := s.fetchUserRoutines(ctx, req.UserId)
	if err != nil {
		log.Printf("GenerateWorkoutSuggestions: Failed to fetch routines: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to fetch routines: %v", err)
	}
	log.Printf("GenerateWorkoutSuggestions: Found %d routines", len(routines))

	// Fetch recent insights
	log.Printf("GenerateWorkoutSuggestions: Fetching recent insights")
	insights, err := s.fetchRecentInsights(ctx, userObjectID)
	if err != nil {
		// Not critical, continue without insights
		log.Printf("GenerateWorkoutSuggestions: Failed to fetch insights: %v", err)
		insights = []WorkoutInsight{}
	}
	log.Printf("GenerateWorkoutSuggestions: Found %d insights", len(insights))

	// Generate suggestions using AI (Gemini will analyze the patterns)
	log.Printf("GenerateWorkoutSuggestions: Generating AI suggestions with pattern analysis")
	suggestions, analysisSummary, err := s.generateAISuggestions(ctx, user, sessions, routines, insights, int(maxSuggestions))
	if err != nil {
		log.Printf("GenerateWorkoutSuggestions: Failed to generate AI suggestions: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to generate suggestions: %v", err)
	}

	log.Printf("GenerateWorkoutSuggestions: Successfully generated %d suggestions", len(suggestions))
	return &pb.GenerateWorkoutSuggestionsResponse{
		Suggestions:     suggestions,
		AnalysisSummary: analysisSummary,
	}, nil
}

// fetchRecentWorkoutSessions fetches workout sessions from the past N days
func (s *SuggestionsService) fetchRecentWorkoutSessions(ctx context.Context, userID primitive.ObjectID, days int) ([]*pb.WorkoutSession, error) {
	collection := s.db.Collection("workout_sessions")

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	filter := bson.M{
		"user_id": userID,
		"finished_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
			"$ne":  nil,
		},
		"is_active": false,
	}

	opts := options.Find().SetSort(bson.D{{Key: "finished_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.WorkoutSession
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	// Convert to protobuf
	var pbSessions []*pb.WorkoutSession
	for _, session := range sessions {
		pbSession := s.workoutSessionToProto(&session)

		// Populate exercise details
		for i, exercise := range pbSession.Exercises {
			exerciseReq := &pb.GetExerciseRequest{Id: exercise.ExerciseId}
			exerciseDetail, err := s.exerciseService.GetExercise(ctx, exerciseReq)
			if err == nil {
				pbSession.Exercises[i].Exercise = exerciseDetail
			}
		}

		pbSessions = append(pbSessions, pbSession)
	}

	return pbSessions, nil
}

// fetchUserRoutines fetches the user's workout routines
func (s *SuggestionsService) fetchUserRoutines(ctx context.Context, userID string) ([]*pb.Workout, error) {
	req := &pb.ListWorkoutsRequest{
		UserId:   userID,
		PageSize: 50,
	}

	resp, err := s.workoutService.ListWorkouts(ctx, req)
	if err != nil {
		return nil, err
	}

	// Filter for routines (no start/finish times)
	var routines []*pb.Workout
	for _, workout := range resp.Workouts {
		if workout.StartedAt == nil && workout.FinishedAt == nil {
			// Populate exercise details
			for i, exercise := range workout.Exercises {
				exerciseReq := &pb.GetExerciseRequest{Id: exercise.ExerciseId}
				exerciseDetail, err := s.exerciseService.GetExercise(ctx, exerciseReq)
				if err == nil {
					workout.Exercises[i].Exercise = exerciseDetail
				}
			}
			routines = append(routines, workout)
		}
	}

	return routines, nil
}

// fetchRecentInsights fetches recent insights for the user
func (s *SuggestionsService) fetchRecentInsights(ctx context.Context, userID primitive.ObjectID) ([]WorkoutInsight, error) {
	return s.insightsService.GetRecentInsights(ctx, userID, 10)
}

// generateAISuggestions uses Gemini AI to generate workout suggestions
func (s *SuggestionsService) generateAISuggestions(ctx context.Context, user *pb.User, sessions []*pb.WorkoutSession, routines []*pb.Workout, insights []WorkoutInsight, maxSuggestions int) ([]*pb.SuggestedWorkout, string, error) {
	log.Printf("generateAISuggestions: Starting AI generation")

	// Prepare data for AI
	prompt := s.prepareAIPrompt(user, sessions, routines, insights, maxSuggestions)
	log.Printf("generateAISuggestions: Prompt length: %d characters", len(prompt))

	// Call Gemini API
	log.Printf("generateAISuggestions: Calling Gemini API")
	result, err := s.genAIClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Printf("generateAISuggestions: Gemini API error: %v", err)
		return nil, "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Parse the AI response
	log.Printf("generateAISuggestions: Parsing AI response")
	responseText := result.Text()
	suggestions, analysisSummary, err := s.parseAIResponse(responseText, routines)
	if err != nil {
		log.Printf("generateAISuggestions: Failed to parse AI response: %v", err)
		return nil, "", fmt.Errorf("failed to parse AI response: %w", err)
	}

	log.Printf("generateAISuggestions: Successfully parsed %d suggestions", len(suggestions))
	return suggestions, analysisSummary, nil
}

func (s *SuggestionsService) prepareAIPrompt(user *pb.User, sessions []*pb.WorkoutSession, routines []*pb.Workout, insights []WorkoutInsight, maxSuggestions int) string {
	var sb strings.Builder

	sb.WriteString("You are an expert personal trainer analyzing a user's workout history to suggest improved workout routines.\n\n")

	// User profile
	sb.WriteString("USER PROFILE:\n")
	sb.WriteString(fmt.Sprintf("- Goal: %s\n", user.Goal))
	sb.WriteString(fmt.Sprintf("- Height: %.1f cm, Weight: %.1f kg, Age: %d\n", user.Height, user.Weight, user.Age))
	sb.WriteString("\n")

	// Recent workout sessions raw data
	sb.WriteString("RECENT WORKOUT SESSIONS (Last 2 weeks):\n")
	sb.WriteString(fmt.Sprintf("Total sessions: %d\n\n", len(sessions)))

	for i, session := range sessions {
		if i >= 10 { // Limit to most recent 10 sessions to avoid overwhelming the prompt
			break
		}

		duration := ""
		if session.DurationSeconds > 0 {
			duration = fmt.Sprintf("%.1f minutes", float64(session.DurationSeconds)/60.0)
		}

		rpe := ""
		if session.RpeRating > 0 {
			rpe = fmt.Sprintf("RPE: %d/10", session.RpeRating)
		}

		sb.WriteString(fmt.Sprintf("Session %d (%s):\n", i+1, session.FinishedAt.AsTime().Format("2006-01-02")))
		sb.WriteString(fmt.Sprintf("  Duration: %s, %s\n", duration, rpe))
		sb.WriteString("  Exercises:\n")

		for _, exercise := range session.Exercises {
			if exercise.Exercise == nil {
				continue
			}

			completedSets := 0
			totalVolume := 0.0
			maxWeight := 0.0

			for _, set := range exercise.Sets {
				if set.Completed {
					completedSets++
					weight := float64(set.ActualWeight)
					reps := float64(set.ActualReps)
					totalVolume += weight * reps
					if weight > maxWeight {
						maxWeight = weight
					}
				}
			}

			muscles := strings.Join(exercise.Exercise.PrimaryMuscles, ", ")
			sb.WriteString(fmt.Sprintf("    - %s (%s): %d sets, max weight: %.1fkg, volume: %.1fkg\n",
				exercise.Exercise.Name, muscles, completedSets, maxWeight, totalVolume))
		}
		sb.WriteString("\n")
	}

	// Recent insights
	if len(insights) > 0 {
		sb.WriteString("RECENT INSIGHTS:\n")
		for _, insight := range insights[:min(5, len(insights))] {
			sb.WriteString(fmt.Sprintf("- %s\n", insight.Insight))
		}
		sb.WriteString("\n")
	}

	// Current routines
	sb.WriteString("CURRENT ROUTINES:\n")
	for i, routine := range routines {
		sb.WriteString(fmt.Sprintf("\nRoutine %d: %s\n", i+1, routine.Name))
		for _, exercise := range routine.Exercises {
			if exercise.Exercise != nil {
				sb.WriteString(fmt.Sprintf("  - %s: %d sets\n", exercise.Exercise.Name, len(exercise.Sets)))
			}
		}
	}
	sb.WriteString("\n")

	// Instructions
	sb.WriteString(fmt.Sprintf(`TASK: 
1. ANALYZE the workout data above to identify:
   - Workout frequency and patterns
   - Exercise progression trends (weight, reps, volume changes over time)
   - Muscle group distribution and potential imbalances
   - Recovery patterns between similar muscle groups
   - Average workout metrics (duration, exercises, sets, RPE)
   - Strengths and weaknesses in their current approach

2. Generate %d improved workout routine suggestions based on the user's goal (%s) and your analysis.

For each suggestion, provide:
1. A modified version of one of their existing routines
2. Specific changes with one-line explanations
3. Overall reasoning for the suggestion

OUTPUT FORMAT (JSON):
{
  "analysis_summary": "2-3 sentence summary of key patterns, progress trends, and recommendations based on your analysis of the workout data",
  "suggestions": [
    {
      "original_workout_id": "routine_id",
      "name": "Improved Routine Name",
      "description": "Brief description of improvements",
      "overall_reasoning": "Why these changes align with their goal based on the patterns you identified",
      "priority": 5,
      "changes": [
        {
          "type": "reps|sets|weight|rest|add|remove|exercise",
          "exercise_name": "Exercise Name",
          "old_value": "3 sets x 10 reps",
          "new_value": "4 sets x 8 reps",
          "reason": "Increased sets and lowered reps to build more strength based on observed weight progression"
        }
      ],
      "exercises": [
        {
          "exercise_id": "existing_exercise_id",
          "sets": [
            {"reps": 8, "weight": 0, "duration_seconds": 0, "distance": 0, "notes": ""},
          ],
          "notes": "",
          "rest_seconds": 90
        }
      ]
    }
  ]
}

Focus your analysis on:
- Progressive overload opportunities based on exercise performance trends
- Addressing muscle imbalances you identify in the data
- Optimizing recovery based on observed workout frequency patterns
- Aligning changes with the user's specific goal
- Being specific with numbers and reasoning based on actual data patterns
- Identifying which exercises are progressing well vs struggling
`, maxSuggestions, user.Goal))

	return sb.String()
}

func (s *SuggestionsService) parseAIResponse(responseText string, routines []*pb.Workout) ([]*pb.SuggestedWorkout, string, error) {
	// Extract JSON from response
	startIdx := strings.Index(responseText, "{")
	endIdx := strings.LastIndex(responseText, "}")
	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return nil, "", fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := responseText[startIdx : endIdx+1]

	// Parse JSON
	var aiResponse struct {
		AnalysisSummary string `json:"analysis_summary"`
		Suggestions     []struct {
			OriginalWorkoutID string `json:"original_workout_id"`
			Name              string `json:"name"`
			Description       string `json:"description"`
			OverallReasoning  string `json:"overall_reasoning"`
			Priority          int32  `json:"priority"`
			Changes           []struct {
				Type         string `json:"type"`
				ExerciseName string `json:"exercise_name"`
				OldValue     string `json:"old_value"`
				NewValue     string `json:"new_value"`
				Reason       string `json:"reason"`
			} `json:"changes"`
			Exercises []struct {
				ExerciseID string `json:"exercise_id"`
				Sets       []struct {
					Reps            int32   `json:"reps"`
					Weight          float32 `json:"weight"`
					DurationSeconds float32 `json:"duration_seconds"`
					Distance        float32 `json:"distance"`
					Notes           string  `json:"notes"`
				} `json:"sets"`
				Notes       string `json:"notes"`
				RestSeconds int32  `json:"rest_seconds"`
			} `json:"exercises"`
		} `json:"suggestions"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &aiResponse); err != nil {
		return nil, "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to protobuf
	var suggestions []*pb.SuggestedWorkout
	for _, suggestion := range aiResponse.Suggestions {
		// Find exercise IDs for changes
		exerciseIDMap := make(map[string]string)
		for _, routine := range routines {
			if routine.Id == suggestion.OriginalWorkoutID {
				for _, ex := range routine.Exercises {
					if ex.Exercise != nil {
						exerciseIDMap[ex.Exercise.Name] = ex.ExerciseId
					}
				}
			}
		}

		// Convert changes
		var changes []*pb.WorkoutChange
		for _, change := range suggestion.Changes {
			pbChange := &pb.WorkoutChange{
				Type:         change.Type,
				ExerciseName: change.ExerciseName,
				OldValue:     change.OldValue,
				NewValue:     change.NewValue,
				Reason:       change.Reason,
			}

			// Try to find exercise ID
			if exerciseID, ok := exerciseIDMap[change.ExerciseName]; ok {
				pbChange.ExerciseId = exerciseID
			}

			changes = append(changes, pbChange)
		}

		// Convert exercises
		var exercises []*pb.WorkoutExercise
		for _, ex := range suggestion.Exercises {
			var sets []*pb.WorkoutSet
			for _, set := range ex.Sets {
				sets = append(sets, &pb.WorkoutSet{
					Reps:            set.Reps,
					Weight:          set.Weight,
					DurationSeconds: set.DurationSeconds,
					Distance:        set.Distance,
					Notes:           set.Notes,
				})
			}

			exercises = append(exercises, &pb.WorkoutExercise{
				ExerciseId:  ex.ExerciseID,
				Sets:        sets,
				Notes:       ex.Notes,
				RestSeconds: ex.RestSeconds,
			})
		}

		suggestions = append(suggestions, &pb.SuggestedWorkout{
			OriginalWorkoutId: suggestion.OriginalWorkoutID,
			Name:              suggestion.Name,
			Description:       suggestion.Description,
			Exercises:         exercises,
			Changes:           changes,
			OverallReasoning:  suggestion.OverallReasoning,
			Priority:          suggestion.Priority,
		})
	}

	return suggestions, aiResponse.AnalysisSummary, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function to convert workout session model to proto
func (s *SuggestionsService) workoutSessionToProto(session *models.WorkoutSession) *pb.WorkoutSession {
	exercises := make([]*pb.WorkoutSessionExercise, len(session.Exercises))
	for i, ex := range session.Exercises {
		sets := make([]*pb.WorkoutSessionSet, len(ex.Sets))
		for j, set := range ex.Sets {
			var startedAt, finishedAt *timestamppb.Timestamp
			if set.StartedAt != nil {
				startedAt = timestamppb.New(*set.StartedAt)
			}
			if set.FinishedAt != nil {
				finishedAt = timestamppb.New(*set.FinishedAt)
			}

			sets[j] = &pb.WorkoutSessionSet{
				TargetReps:      set.TargetReps,
				TargetWeight:    set.TargetWeight,
				ActualReps:      set.ActualReps,
				ActualWeight:    set.ActualWeight,
				DurationSeconds: float32(set.DurationSeconds),
				Distance:        set.Distance,
				Notes:           set.Notes,
				Completed:       set.Completed,
				StartedAt:       startedAt,
				FinishedAt:      finishedAt,
			}
		}

		var exerciseStartedAt, exerciseFinishedAt *timestamppb.Timestamp
		if ex.StartedAt != nil {
			exerciseStartedAt = timestamppb.New(*ex.StartedAt)
		}
		if ex.FinishedAt != nil {
			exerciseFinishedAt = timestamppb.New(*ex.FinishedAt)
		}

		exercises[i] = &pb.WorkoutSessionExercise{
			ExerciseId:  ex.ExerciseID.Hex(),
			Sets:        sets,
			Notes:       ex.Notes,
			RestSeconds: ex.RestSeconds,
			Completed:   ex.Completed,
			StartedAt:   exerciseStartedAt,
			FinishedAt:  exerciseFinishedAt,
		}
	}

	var finishedAt *timestamppb.Timestamp
	if session.FinishedAt != nil {
		finishedAt = timestamppb.New(*session.FinishedAt)
	}

	return &pb.WorkoutSession{
		Id:              session.ID.Hex(),
		UserId:          session.UserID.Hex(),
		RoutineId:       session.RoutineID.Hex(),
		Name:            session.Name,
		Description:     session.Description,
		Exercises:       exercises,
		StartedAt:       timestamppb.New(session.StartedAt),
		FinishedAt:      finishedAt,
		DurationSeconds: session.DurationSeconds,
		Notes:           session.Notes,
		RpeRating:       session.RPERating,
		IsActive:        session.IsActive,
		CreatedAt:       timestamppb.New(session.CreatedAt),
		UpdatedAt:       timestamppb.New(session.UpdatedAt),
	}
}
