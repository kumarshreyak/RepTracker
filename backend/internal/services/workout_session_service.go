package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gymlog-backend/pkg/database"
	"gymlog-backend/pkg/models"
	pb "gymlog-backend/proto/gymlog/v1"
)

// WorkoutSessionService implements the gRPC WorkoutSessionService
type WorkoutSessionService struct {
	pb.UnimplementedWorkoutSessionServiceServer
	db                        *database.MongoDB
	sessionColl               *mongo.Collection
	aiProgressiveOverloadColl *mongo.Collection
	workoutService            *WorkoutService
	exerciseService           *ExerciseService
	metricsService            *MetricsService
	genAIClient               *genai.Client
	userService               *UserService
	tempAIExerciseDetails     map[string]models.AIProgressiveOverloadExercise // Added for storing AI exercise details
}

// NewWorkoutSessionService creates a new WorkoutSessionService instance
func NewWorkoutSessionService(db *database.MongoDB, workoutService *WorkoutService, exerciseService *ExerciseService, metricsService *MetricsService, userService *UserService) (*WorkoutSessionService, error) {
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

	return &WorkoutSessionService{
		db:                        db,
		sessionColl:               db.GetCollection("workout_sessions"),
		aiProgressiveOverloadColl: db.GetCollection("ai_progressive_overload_recommendations"),
		workoutService:            workoutService,
		exerciseService:           exerciseService,
		metricsService:            metricsService,
		genAIClient:               client,
		userService:               userService,
	}, nil
}

// CreateWorkoutSession creates a new workout session from a routine/workout
func (s *WorkoutSessionService) CreateWorkoutSession(ctx context.Context, req *pb.CreateWorkoutSessionRequest) (*pb.WorkoutSession, error) {
	// Validate required fields
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.RoutineId == "" {
		return nil, status.Error(codes.InvalidArgument, "routine_id is required")
	}

	// Convert IDs to ObjectIDs
	userObjectID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	routineObjectID, err := primitive.ObjectIDFromHex(req.RoutineId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid routine_id")
	}

	// Get the routine/workout template
	workoutReq := &pb.GetWorkoutRequest{Id: req.RoutineId}
	routine, err := s.workoutService.GetWorkout(ctx, workoutReq)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "routine not found: %v", err)
	}

	now := time.Now()

	// Convert routine exercises to session exercises with target values
	var sessionExercises []models.WorkoutSessionExercise
	for _, routineExercise := range routine.Exercises {
		// Get exercise details (optional, for reference)
		exerciseReq := &pb.GetExerciseRequest{Id: routineExercise.ExerciseId}
		_, err := s.exerciseService.GetExercise(ctx, exerciseReq)
		if err != nil {
			// Log warning but continue - exercise details are optional
		}

		exerciseObjectID, err := primitive.ObjectIDFromHex(routineExercise.ExerciseId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid exercise_id in routine")
		}

		// Convert sets to session sets with target and actual values
		var sessionSets []models.WorkoutSessionSet
		for _, routineSet := range routineExercise.Sets {
			sessionSets = append(sessionSets, models.WorkoutSessionSet{
				TargetReps:      routineSet.Reps,
				TargetWeight:    routineSet.Weight,
				ActualReps:      0, // To be filled during workout
				ActualWeight:    0, // To be filled during workout
				DurationSeconds: int32(routineSet.DurationSeconds),
				Distance:        routineSet.Distance,
				Notes:           routineSet.Notes,
				Completed:       false,
				StartedAt:       nil,
				FinishedAt:      nil,
			})
		}

		sessionExercises = append(sessionExercises, models.WorkoutSessionExercise{
			ExerciseID:  exerciseObjectID,
			Sets:        sessionSets,
			Notes:       routineExercise.Notes,
			RestSeconds: routineExercise.RestSeconds,
			Completed:   false,
			StartedAt:   nil,
			FinishedAt:  nil,
		})
	}

	// Use provided name or generate one
	sessionName := req.Name
	if sessionName == "" {
		sessionName = fmt.Sprintf("%s Session - %s", routine.Name, now.Format("Jan 2, 2006"))
	}

	session := models.WorkoutSession{
		UserID:          userObjectID,
		RoutineID:       routineObjectID,
		Name:            sessionName,
		Description:     req.Description,
		Exercises:       sessionExercises,
		StartedAt:       now,
		FinishedAt:      nil,
		DurationSeconds: 0,
		Notes:           req.Notes,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	result, err := s.sessionColl.InsertOne(ctx, session)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create workout session")
	}

	session.ID = result.InsertedID.(primitive.ObjectID)
	return s.modelToProto(&session), nil
}

// GetWorkoutSession retrieves a workout session by ID
func (s *WorkoutSessionService) GetWorkoutSession(ctx context.Context, req *pb.GetWorkoutSessionRequest) (*pb.WorkoutSession, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "session id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid session ID")
	}

	var session models.WorkoutSession
	err = s.sessionColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&session)
	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "workout session not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "failed to get workout session")
	}

	return s.modelToProto(&session), nil
}

// UpdateWorkoutSession updates an existing workout session
func (s *WorkoutSessionService) UpdateWorkoutSession(ctx context.Context, req *pb.UpdateWorkoutSessionRequest) (*pb.WorkoutSession, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "session id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid session ID")
	}

	update := bson.M{
		"updated_at": time.Now(),
	}

	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if len(req.Exercises) > 0 {
		// Convert protobuf exercises to model exercises
		var sessionExercises []models.WorkoutSessionExercise
		for _, exercise := range req.Exercises {
			exerciseObjectID, err := primitive.ObjectIDFromHex(exercise.ExerciseId)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "invalid exercise_id")
			}

			var sessionSets []models.WorkoutSessionSet
			for _, set := range exercise.Sets {
				var startedAt, finishedAt *time.Time
				if set.StartedAt != nil {
					t := set.StartedAt.AsTime()
					startedAt = &t
				}
				if set.FinishedAt != nil {
					t := set.FinishedAt.AsTime()
					finishedAt = &t
				}

				sessionSets = append(sessionSets, models.WorkoutSessionSet{
					TargetReps:      set.TargetReps,
					TargetWeight:    set.TargetWeight,
					ActualReps:      set.ActualReps,
					ActualWeight:    set.ActualWeight,
					DurationSeconds: int32(set.DurationSeconds),
					Distance:        set.Distance,
					Notes:           set.Notes,
					Completed:       set.Completed,
					StartedAt:       startedAt,
					FinishedAt:      finishedAt,
				})
			}

			var exerciseStartedAt, exerciseFinishedAt *time.Time
			if exercise.StartedAt != nil {
				t := exercise.StartedAt.AsTime()
				exerciseStartedAt = &t
			}
			if exercise.FinishedAt != nil {
				t := exercise.FinishedAt.AsTime()
				exerciseFinishedAt = &t
			}

			sessionExercises = append(sessionExercises, models.WorkoutSessionExercise{
				ExerciseID:  exerciseObjectID,
				Sets:        sessionSets,
				Notes:       exercise.Notes,
				RestSeconds: exercise.RestSeconds,
				Completed:   exercise.Completed,
				StartedAt:   exerciseStartedAt,
				FinishedAt:  exerciseFinishedAt,
			})
		}
		update["exercises"] = sessionExercises
	}
	if req.FinishedAt != nil {
		t := req.FinishedAt.AsTime()
		update["finished_at"] = &t
		update["is_active"] = false

		// Calculate duration if session is finished
		var session models.WorkoutSession
		err = s.sessionColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&session)
		if err == nil {
			duration := t.Sub(session.StartedAt)
			update["duration_seconds"] = int32(duration.Seconds())
		}
	}
	if req.DurationSeconds > 0 {
		update["duration_seconds"] = req.DurationSeconds
	}
	if req.Notes != "" {
		update["notes"] = req.Notes
	}
	if req.RpeRating > 0 {
		update["rpe_rating"] = req.RpeRating
	}
	if req.IsActive != update["is_active"] {
		update["is_active"] = req.IsActive
	}

	result, err := s.sessionColl.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update workout session")
	}
	if result.MatchedCount == 0 {
		return nil, status.Error(codes.NotFound, "workout session not found")
	}

	// If session was finished, calculate and save metrics
	if req.FinishedAt != nil {
		// Get the updated session
		var updatedSession models.WorkoutSession
		err = s.sessionColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedSession)
		if err == nil && s.metricsService != nil {
			// Calculate workout metrics
			workoutMetrics, err := s.metricsService.CalculateWorkoutMetrics(ctx, &updatedSession)
			if err == nil {
				// Save workout metrics
				err = s.metricsService.SaveWorkoutMetrics(ctx, workoutMetrics)
				if err == nil {
					// Update user metrics
					_ = s.metricsService.UpdateUserMetrics(ctx, updatedSession.UserID)
				}
			}
			// Note: We don't return errors here to avoid failing the session update
			// Metrics calculation is secondary to the main workout functionality
		}
	}

	return s.GetWorkoutSession(ctx, &pb.GetWorkoutSessionRequest{Id: req.Id})
}

// DeleteWorkoutSession deletes a workout session
func (s *WorkoutSessionService) DeleteWorkoutSession(ctx context.Context, req *pb.DeleteWorkoutSessionRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "session id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid session ID")
	}

	result, err := s.sessionColl.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete workout session")
	}
	if result.DeletedCount == 0 {
		return nil, status.Error(codes.NotFound, "workout session not found")
	}

	return &emptypb.Empty{}, nil
}

// ListWorkoutSessions lists workout sessions with optional filtering
func (s *WorkoutSessionService) ListWorkoutSessions(ctx context.Context, req *pb.ListWorkoutSessionsRequest) (*pb.ListWorkoutSessionsResponse, error) {
	filter := bson.M{}

	// Apply filters
	if req.UserId != "" {
		userObjectID, err := primitive.ObjectIDFromHex(req.UserId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid user_id")
		}
		filter["user_id"] = userObjectID
	}

	// Filter by active status
	if req.ActiveOnly {
		filter["is_active"] = true
	}

	// Filter by date range
	if req.StartDate != nil || req.EndDate != nil {
		dateFilter := bson.M{}
		if req.StartDate != nil {
			dateFilter["$gte"] = req.StartDate.AsTime()
		}
		if req.EndDate != nil {
			dateFilter["$lte"] = req.EndDate.AsTime()
		}
		filter["started_at"] = dateFilter
	}

	// Apply pagination
	pageSize := int64(req.PageSize)
	if pageSize <= 0 {
		pageSize = 50 // Default page size
	}

	skip := int64(0)
	if req.PageToken != "" {
		// In a real app, decode the page token to get the skip value
		// For simplicity, we're not implementing pagination here
	}

	opts := options.Find().SetLimit(pageSize).SetSkip(skip).SetSort(bson.D{{"created_at", -1}})
	cursor, err := s.sessionColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list workout sessions")
	}
	defer cursor.Close(ctx)

	var sessions []models.WorkoutSession
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, status.Error(codes.Internal, "failed to decode workout sessions")
	}

	var protoSessions []*pb.WorkoutSession
	for _, session := range sessions {
		protoSessions = append(protoSessions, s.modelToProto(&session))
	}

	response := &pb.ListWorkoutSessionsResponse{
		Sessions: protoSessions,
	}

	// Set next page token if there are more results
	if int64(len(sessions)) == pageSize {
		response.NextPageToken = "next_page" // In real app, encode the next skip position
	}

	return response, nil
}

// StartExercise marks an exercise as started in a workout session
func (s *WorkoutSessionService) StartExercise(ctx context.Context, req *pb.StartExerciseRequest) (*pb.WorkoutSession, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	session, err := s.GetWorkoutSession(ctx, &pb.GetWorkoutSessionRequest{Id: req.SessionId})
	if err != nil {
		return nil, err
	}

	if req.ExerciseIndex < 0 || int(req.ExerciseIndex) >= len(session.Exercises) {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise_index")
	}

	// Mark exercise as started
	now := timestamppb.Now()
	session.Exercises[req.ExerciseIndex].StartedAt = now

	// Update the session
	updateReq := &pb.UpdateWorkoutSessionRequest{
		Id:        req.SessionId,
		Exercises: session.Exercises,
	}

	return s.UpdateWorkoutSession(ctx, updateReq)
}

// FinishExercise marks an exercise as completed in a workout session
func (s *WorkoutSessionService) FinishExercise(ctx context.Context, req *pb.FinishExerciseRequest) (*pb.WorkoutSession, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	session, err := s.GetWorkoutSession(ctx, &pb.GetWorkoutSessionRequest{Id: req.SessionId})
	if err != nil {
		return nil, err
	}

	if req.ExerciseIndex < 0 || int(req.ExerciseIndex) >= len(session.Exercises) {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise_index")
	}

	// Mark exercise as completed
	now := timestamppb.Now()
	exercise := session.Exercises[req.ExerciseIndex]
	exercise.FinishedAt = now
	exercise.Completed = true
	if req.Notes != "" {
		exercise.Notes = req.Notes
	}

	// Update the session
	updateReq := &pb.UpdateWorkoutSessionRequest{
		Id:        req.SessionId,
		Exercises: session.Exercises,
	}

	return s.UpdateWorkoutSession(ctx, updateReq)
}

// UpdateSet updates a specific set in a workout session
func (s *WorkoutSessionService) UpdateSet(ctx context.Context, req *pb.UpdateSetRequest) (*pb.WorkoutSession, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	session, err := s.GetWorkoutSession(ctx, &pb.GetWorkoutSessionRequest{Id: req.SessionId})
	if err != nil {
		return nil, err
	}

	if req.ExerciseIndex < 0 || int(req.ExerciseIndex) >= len(session.Exercises) {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise_index")
	}

	exercise := session.Exercises[req.ExerciseIndex]
	if req.SetIndex < 0 || int(req.SetIndex) >= len(exercise.Sets) {
		return nil, status.Error(codes.InvalidArgument, "invalid set_index")
	}

	// Update set data
	set := exercise.Sets[req.SetIndex]
	now := timestamppb.Now()

	if req.ActualReps > 0 {
		set.ActualReps = req.ActualReps
	}
	if req.ActualWeight >= 0 {
		set.ActualWeight = req.ActualWeight
	}
	if req.DurationSeconds >= 0 {
		set.DurationSeconds = req.DurationSeconds
	}
	if req.Distance >= 0 {
		set.Distance = req.Distance
	}
	if req.Notes != "" {
		set.Notes = req.Notes
	}

	// Handle completion state
	if req.Completed != set.Completed {
		set.Completed = req.Completed
		if req.Completed {
			set.FinishedAt = now
			if set.StartedAt == nil {
				set.StartedAt = now
			}
		} else {
			set.FinishedAt = nil
		}
	}

	// Check if all sets in exercise are completed to auto-complete exercise
	allSetsCompleted := true
	for _, s := range exercise.Sets {
		if !s.Completed {
			allSetsCompleted = false
			break
		}
	}
	exercise.Completed = allSetsCompleted
	if allSetsCompleted && exercise.FinishedAt == nil {
		exercise.FinishedAt = now
	}

	// Update the session
	updateReq := &pb.UpdateWorkoutSessionRequest{
		Id:        req.SessionId,
		Exercises: session.Exercises,
	}

	return s.UpdateWorkoutSession(ctx, updateReq)
}

// ApplyProgressiveOverload applies progressive overload to a workout based on session performance
func (s *WorkoutSessionService) ApplyProgressiveOverload(ctx context.Context, req *pb.ApplyProgressiveOverloadRequest) (*pb.ApplyProgressiveOverloadResponse, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	if req.WorkoutId == "" {
		return nil, status.Error(codes.InvalidArgument, "workout_id is required")
	}

	// Get the workout session
	sessionReq := &pb.GetWorkoutSessionRequest{Id: req.SessionId}
	session, err := s.GetWorkoutSession(ctx, sessionReq)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "workout session not found: %v", err)
	}

	// Get the workout/routine
	workoutReq := &pb.GetWorkoutRequest{Id: req.WorkoutId}
	workout, err := s.workoutService.GetWorkout(ctx, workoutReq)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "workout not found: %v", err)
	}

	// Create a map of exercise IDs to session exercises for quick lookup
	sessionExerciseMap := make(map[string]*pb.WorkoutSessionExercise)
	for _, sessionExercise := range session.Exercises {
		sessionExerciseMap[sessionExercise.ExerciseId] = sessionExercise
	}

	// Track if any changes were made
	modified := false

	// Process each exercise in the workout
	for i, workoutExercise := range workout.Exercises {
		sessionExercise, exists := sessionExerciseMap[workoutExercise.ExerciseId]
		if !exists {
			// Exercise doesn't exist in session, skip
			continue
		}

		// Process each set in the workout exercise
		for j := range workoutExercise.Sets {
			// Find the corresponding session set (by index)
			if j >= len(sessionExercise.Sets) {
				// No corresponding session set, skip
				continue
			}

			sessionSet := sessionExercise.Sets[j]

			// Only process completed sets
			if !sessionSet.Completed {
				continue
			}

			// Apply progressive overload logic
			// Step 1: Increase rep count by 1
			newReps := sessionSet.ActualReps + 1
			newWeight := sessionSet.ActualWeight

			// Step 2: If rep count > 15, increase weight by 2.5 and set reps to 8 (only if weight is non-zero)
			if newReps > 15 && sessionSet.ActualWeight > 0 {
				newWeight = sessionSet.ActualWeight + 2.5
				newReps = 8
			}

			// Update the workout set with new values
			workout.Exercises[i].Sets[j].Reps = newReps
			workout.Exercises[i].Sets[j].Weight = newWeight
			modified = true
		}
	}

	if !modified {
		return &pb.ApplyProgressiveOverloadResponse{
			Success:        true,
			Message:        "No changes made - no completed sets found or no matching exercises",
			UpdatedWorkout: workout,
		}, nil
	}

	// Update the workout with the new values
	updateWorkoutReq := &pb.UpdateWorkoutRequest{
		Id:          req.WorkoutId,
		Name:        workout.Name,
		Description: workout.Description,
		Exercises:   workout.Exercises,
		Notes:       workout.Notes,
	}

	updatedWorkout, err := s.workoutService.UpdateWorkout(ctx, updateWorkoutReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update workout: %v", err)
	}

	return &pb.ApplyProgressiveOverloadResponse{
		Success:        true,
		Message:        "Progressive overload applied successfully",
		UpdatedWorkout: updatedWorkout,
	}, nil
}

// ApplyAIProgressiveOverload applies progressive overload using AI analysis of the workout session
func (s *WorkoutSessionService) ApplyAIProgressiveOverload(ctx context.Context, req *pb.ApplyProgressiveOverloadRequest) (*pb.ApplyProgressiveOverloadResponse, error) {
	startTime := time.Now()

	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	if req.WorkoutId == "" {
		return nil, status.Error(codes.InvalidArgument, "workout_id is required")
	}

	log.Printf("ApplyAIProgressiveOverload: Starting for session %s, workout %s", req.SessionId, req.WorkoutId)

	// Get the workout session
	sessionReq := &pb.GetWorkoutSessionRequest{Id: req.SessionId}
	session, err := s.GetWorkoutSession(ctx, sessionReq)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "workout session not found: %v", err)
	}

	// Get the workout/routine
	workoutReq := &pb.GetWorkoutRequest{Id: req.WorkoutId}
	workout, err := s.workoutService.GetWorkout(ctx, workoutReq)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "workout not found: %v", err)
	}

	// Get user profile
	userReq := &pb.GetUserRequest{Id: session.UserId}
	user, err := s.userService.GetUser(ctx, userReq)
	if err != nil {
		log.Printf("ApplyAIProgressiveOverload: Failed to fetch user profile: %v", err)
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	log.Printf("ApplyAIProgressiveOverload: User goal: %s", user.Goal)

	// Generate AI recommendations
	log.Printf("ApplyAIProgressiveOverload: Generating AI recommendations")
	updatedExercises, analysisSummary, recentSessionsCount, err := s.generateAIProgressiveOverload(ctx, user, session, workout)
	if err != nil {
		log.Printf("ApplyAIProgressiveOverload: Failed to generate AI recommendations: %v", err)

		// Store failed response
		processingTime := time.Since(startTime).Milliseconds()
		_ = s.storeAIProgressiveOverloadResponse(ctx, session, workout, "", nil, false, fmt.Sprintf("Failed to generate AI recommendations: %v", err), 0, processingTime)

		return nil, status.Errorf(codes.Internal, "failed to generate AI recommendations: %v", err)
	}

	log.Printf("ApplyAIProgressiveOverload: AI analysis: %s", analysisSummary)
	log.Printf("ApplyAIProgressiveOverload: Generated %d updated exercises", len(updatedExercises))

	successMessage := "AI progressive overload applied successfully - " + analysisSummary

	if len(updatedExercises) == 0 {
		successMessage = "No changes recommended by AI - " + analysisSummary

		// Store no-changes response
		processingTime := time.Since(startTime).Milliseconds()
		_ = s.storeAIProgressiveOverloadResponse(ctx, session, workout, analysisSummary, nil, true, successMessage, recentSessionsCount, processingTime)

		return &pb.ApplyProgressiveOverloadResponse{
			Success:        true,
			Message:        successMessage,
			UpdatedWorkout: workout,
		}, nil
	}

	// Apply the AI recommendations to the workout
	workout.Exercises = updatedExercises

	// Update the workout with the new values
	updateWorkoutReq := &pb.UpdateWorkoutRequest{
		Id:          req.WorkoutId,
		Name:        workout.Name,
		Description: workout.Description,
		Exercises:   workout.Exercises,
		Notes:       workout.Notes,
	}

	updatedWorkout, err := s.workoutService.UpdateWorkout(ctx, updateWorkoutReq)
	if err != nil {
		// Store failed response
		processingTime := time.Since(startTime).Milliseconds()
		_ = s.storeAIProgressiveOverloadResponse(ctx, session, workout, analysisSummary, updatedExercises, false, fmt.Sprintf("Failed to update workout: %v", err), recentSessionsCount, processingTime)

		return nil, status.Errorf(codes.Internal, "failed to update workout: %v", err)
	}

	// Store successful response
	processingTime := time.Since(startTime).Milliseconds()
	err = s.storeAIProgressiveOverloadResponse(ctx, session, updatedWorkout, analysisSummary, updatedExercises, true, successMessage, recentSessionsCount, processingTime)
	if err != nil {
		log.Printf("ApplyAIProgressiveOverload: Failed to store AI response: %v", err)
		// Don't fail the request if storage fails, just log the error
	}

	log.Printf("ApplyAIProgressiveOverload: Successfully applied AI recommendations")

	return &pb.ApplyProgressiveOverloadResponse{
		Success:        true,
		Message:        successMessage,
		UpdatedWorkout: updatedWorkout,
	}, nil
}

// generateAIProgressiveOverload uses Gemini AI to generate progressive overload recommendations
func (s *WorkoutSessionService) generateAIProgressiveOverload(ctx context.Context, user *pb.User, session *pb.WorkoutSession, workout *pb.Workout) ([]*pb.WorkoutExercise, string, int, error) {
	log.Printf("generateAIProgressiveOverload: Starting AI generation")

	// Fetch recent workout sessions (last 2 weeks) for better context
	log.Printf("generateAIProgressiveOverload: Fetching recent workout sessions for user %s", session.UserId)
	recentSessions, err := s.fetchRecentWorkoutSessions(ctx, session.UserId, 14) // Last 14 days
	if err != nil {
		log.Printf("generateAIProgressiveOverload: Failed to fetch recent sessions: %v", err)
		// Continue with just the current session if recent sessions fail
		recentSessions = []*pb.WorkoutSession{session}
	} else {
		log.Printf("generateAIProgressiveOverload: Found %d recent workout sessions", len(recentSessions))
		// Ensure current session is included if not already in recent sessions
		sessionFound := false
		for _, rs := range recentSessions {
			if rs.Id == session.Id {
				sessionFound = true
				break
			}
		}
		if !sessionFound {
			recentSessions = append([]*pb.WorkoutSession{session}, recentSessions...)
		}
	}

	// Prepare data for AI with recent sessions context
	prompt := s.prepareAIProgressiveOverloadPrompt(user, session, workout, recentSessions)
	log.Printf("generateAIProgressiveOverload: Prompt length: %d characters", len(prompt))

	// Log the full prompt for debugging
	log.Printf("generateAIProgressiveOverload: GEMINI PROMPT:\n%s", prompt)

	// Call Gemini API
	log.Printf("generateAIProgressiveOverload: Calling Gemini API")
	result, err := s.genAIClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Printf("generateAIProgressiveOverload: Gemini API error: %v", err)
		return nil, "", 0, fmt.Errorf("failed to generate content: %w", err)
	}

	// Parse the AI response
	log.Printf("generateAIProgressiveOverload: Parsing AI response")
	responseText := result.Text()

	// Log the full response for debugging
	log.Printf("generateAIProgressiveOverload: GEMINI RESPONSE:\n%s", responseText)

	updatedExercises, analysisSummary, aiExerciseDetails, err := s.parseAIProgressiveOverloadResponse(responseText, workout)
	if err != nil {
		log.Printf("generateAIProgressiveOverload: Failed to parse AI response: %v", err)
		return nil, "", 0, fmt.Errorf("failed to parse AI response: %w", err)
	}

	log.Printf("generateAIProgressiveOverload: Successfully parsed %d updated exercises", len(updatedExercises))

	// Store AI exercise details for later use in storage
	s.tempAIExerciseDetails = aiExerciseDetails

	return updatedExercises, analysisSummary, len(recentSessions), nil
}

// fetchRecentWorkoutSessions fetches workout sessions from the past N days for the given user
func (s *WorkoutSessionService) fetchRecentWorkoutSessions(ctx context.Context, userID string, days int) ([]*pb.WorkoutSession, error) {
	log.Printf("fetchRecentWorkoutSessions: Starting for user %s, last %d days", userID, days)

	// Convert user ID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Calculate date range
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)

	log.Printf("fetchRecentWorkoutSessions: Fetching sessions from %s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))

	// Build MongoDB filter
	filter := bson.M{
		"user_id": userObjectID,
		"finished_at": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
		"is_active": false, // Only completed sessions
	}

	// Query options - sort by finished time descending, limit to reasonable number
	opts := options.Find().
		SetSort(bson.D{{Key: "finishedAt", Value: -1}}).
		SetLimit(50) // Reasonable limit to prevent excessive data

	cursor, err := s.sessionColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query workout sessions: %w", err)
	}
	defer cursor.Close(ctx)

	var sessionModels []models.WorkoutSession
	if err = cursor.All(ctx, &sessionModels); err != nil {
		return nil, fmt.Errorf("failed to decode workout sessions: %w", err)
	}

	log.Printf("fetchRecentWorkoutSessions: Found %d sessions", len(sessionModels))

	// Convert to protobuf
	var sessions []*pb.WorkoutSession
	for _, sessionModel := range sessionModels {
		pbSession := s.modelToProto(&sessionModel)
		sessions = append(sessions, pbSession)
	}

	log.Printf("fetchRecentWorkoutSessions: Successfully converted %d sessions to protobuf", len(sessions))
	return sessions, nil
}

// loadProgressiveOverloadPromptTemplate loads the AI progressive overload prompt template from file
func (s *WorkoutSessionService) loadProgressiveOverloadPromptTemplate() (string, error) {
	// Get the path relative to the project root
	promptPath := filepath.Join("prompts", "progressive_overload_prompt.txt")

	file, err := os.Open(promptPath)
	if err != nil {
		return "", fmt.Errorf("failed to open prompt template file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt template file: %w", err)
	}

	return string(content), nil
}

// prepareAIProgressiveOverloadPrompt prepares the prompt for AI progressive overload analysis
func (s *WorkoutSessionService) prepareAIProgressiveOverloadPrompt(user *pb.User, session *pb.WorkoutSession, workout *pb.Workout, recentSessions []*pb.WorkoutSession) string {
	// Load the prompt template
	template, err := s.loadProgressiveOverloadPromptTemplate()
	if err != nil {
		log.Printf("Failed to load progressive overload prompt template, using fallback: %v", err)
		template = s.getFallbackProgressiveOverloadPromptTemplate()
	}

	// Build session exercises data
	sessionExercisesData := s.buildSessionExercisesData(session)

	// Calculate workout duration
	workoutDuration := float64(0)
	if session.FinishedAt != nil && session.StartedAt != nil {
		workoutDuration = session.FinishedAt.AsTime().Sub(session.StartedAt.AsTime()).Minutes()
	}

	// Replace template variables
	prompt := strings.ReplaceAll(template, "{{USER_GOAL}}", user.Goal)
	prompt = strings.ReplaceAll(prompt, "{{USER_HEIGHT}}", fmt.Sprintf("%.1f", user.Height))
	prompt = strings.ReplaceAll(prompt, "{{USER_WEIGHT}}", fmt.Sprintf("%.1f", user.Weight))
	prompt = strings.ReplaceAll(prompt, "{{USER_AGE}}", fmt.Sprintf("%d", user.Age))
	prompt = strings.ReplaceAll(prompt, "{{SESSION_ID}}", session.Id)
	prompt = strings.ReplaceAll(prompt, "{{WORKOUT_NAME}}", session.Name)
	prompt = strings.ReplaceAll(prompt, "{{WORKOUT_DURATION}}", fmt.Sprintf("%.1f", workoutDuration))
	prompt = strings.ReplaceAll(prompt, "{{SESSION_DATE}}", session.StartedAt.AsTime().Format("2006-01-02"))
	prompt = strings.ReplaceAll(prompt, "{{SESSION_EXERCISES_DATA}}", sessionExercisesData)

	// Add context about recent workout sessions
	recentSessionsData := s.buildRecentSessionsData(recentSessions)
	prompt += "\n\n" + recentSessionsData

	return prompt
}

// buildSessionExercisesData builds the session exercises data section
func (s *WorkoutSessionService) buildSessionExercisesData(session *pb.WorkoutSession) string {
	var sb strings.Builder

	sb.WriteString("EXERCISES PERFORMED:\n\n")

	for _, exercise := range session.Exercises {
		exerciseName := "Unknown Exercise"
		if exercise.Exercise != nil {
			exerciseName = exercise.Exercise.Name
		}

		sb.WriteString(fmt.Sprintf("Exercise: %s (ID: %s)\n", exerciseName, exercise.ExerciseId))
		sb.WriteString(fmt.Sprintf("Completed: %t\n", exercise.Completed))

		if exercise.Notes != "" {
			sb.WriteString(fmt.Sprintf("Notes: %s\n", exercise.Notes))
		}

		sb.WriteString("Sets:\n")
		for i, set := range exercise.Sets {
			completedStatus := "✗"
			if set.Completed {
				completedStatus = "✓"
			}

			sb.WriteString(fmt.Sprintf("  Set %d: %s ", i+1, completedStatus))

			if set.TargetReps > 0 && set.TargetWeight > 0 {
				sb.WriteString(fmt.Sprintf("Target: %d reps @ %.1fkg", set.TargetReps, set.TargetWeight))
			} else if set.TargetReps > 0 {
				sb.WriteString(fmt.Sprintf("Target: %d reps", set.TargetReps))
			}

			if set.Completed {
				if set.ActualReps > 0 && set.ActualWeight > 0 {
					sb.WriteString(fmt.Sprintf(" | Actual: %d reps @ %.1fkg", set.ActualReps, set.ActualWeight))
				} else if set.ActualReps > 0 {
					sb.WriteString(fmt.Sprintf(" | Actual: %d reps", set.ActualReps))
				}
			}

			if set.DurationSeconds > 0 {
				sb.WriteString(fmt.Sprintf(" | Duration: %.0fs", set.DurationSeconds))
			}

			if set.Notes != "" {
				sb.WriteString(fmt.Sprintf(" | Notes: %s", set.Notes))
			}

			sb.WriteString("\n")
		}

		// Add RPE if available (note: RPE is per session, not per exercise in current model)
		if session.RpeRating > 0 {
			sb.WriteString(fmt.Sprintf("Session RPE: %d/10\n", session.RpeRating))
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// buildRecentSessionsData builds the recent workout sessions data section
func (s *WorkoutSessionService) buildRecentSessionsData(recentSessions []*pb.WorkoutSession) string {
	var sb strings.Builder

	sb.WriteString("RECENT WORKOUT SESSIONS (for context):\n\n")

	for i, session := range recentSessions {
		sb.WriteString(fmt.Sprintf("Session %d (ID: %s)\n", i+1, session.Id))
		sb.WriteString(fmt.Sprintf("Date: %s, Duration: %.1f minutes\n", session.StartedAt.AsTime().Format("2006-01-02"), float64(session.DurationSeconds)/60))
		sb.WriteString(fmt.Sprintf("RPE: %d/10\n", session.RpeRating))
		if session.Notes != "" {
			sb.WriteString(fmt.Sprintf("Notes: %s\n", session.Notes))
		}
		sb.WriteString("Exercises:\n")
		for _, exercise := range session.Exercises {
			exerciseName := "Unknown Exercise"
			if exercise.Exercise != nil {
				exerciseName = exercise.Exercise.Name
			}
			sb.WriteString(fmt.Sprintf("  - %s (ID: %s)\n", exerciseName, exercise.ExerciseId))
			sb.WriteString(fmt.Sprintf("    Completed: %t", exercise.Completed))
			if exercise.Notes != "" {
				sb.WriteString(fmt.Sprintf(", Notes: %s", exercise.Notes))
			}
			sb.WriteString("\n    Sets:\n")
			for _, set := range exercise.Sets {
				if set.Completed {
					sb.WriteString(fmt.Sprintf("      - Target: %d reps @ %.1fkg, Actual: %d reps @ %.1fkg (✓)\n",
						set.TargetReps, set.TargetWeight, set.ActualReps, set.ActualWeight))
				} else {
					sb.WriteString(fmt.Sprintf("      - Target: %d reps @ %.1fkg (✗ incomplete)\n",
						set.TargetReps, set.TargetWeight))
				}
			}
		}
		sb.WriteString("----------------------------------------\n")
	}

	return sb.String()
}

// parseAIProgressiveOverloadResponse parses the AI response and extracts updated exercises
func (s *WorkoutSessionService) parseAIProgressiveOverloadResponse(responseText string, originalWorkout *pb.Workout) ([]*pb.WorkoutExercise, string, map[string]models.AIProgressiveOverloadExercise, error) {
	// Extract JSON from response
	startIdx := strings.Index(responseText, "{")
	endIdx := strings.LastIndex(responseText, "}")
	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return nil, "", nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := responseText[startIdx : endIdx+1]

	// Parse JSON
	var aiResponse struct {
		AnalysisSummary  string `json:"analysis_summary"`
		UpdatedExercises []struct {
			ExerciseID         string `json:"exercise_id"`
			ExerciseName       string `json:"exercise_name"`
			ProgressionApplied bool   `json:"progression_applied"`
			Reasoning          string `json:"reasoning"`
			ChangesMade        string `json:"changes_made"`
			Sets               []struct {
				Reps            int32   `json:"reps"`
				Weight          float32 `json:"weight"`
				DurationSeconds float32 `json:"duration_seconds"`
				Distance        float32 `json:"distance"`
				Notes           string  `json:"notes"`
			} `json:"sets"`
			Notes       string `json:"notes"`
			RestSeconds int32  `json:"rest_seconds"`
		} `json:"updated_exercises"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &aiResponse); err != nil {
		return nil, "", nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Create a map of exercise IDs to original exercises for validation
	originalExerciseMap := make(map[string]*pb.WorkoutExercise)
	for _, exercise := range originalWorkout.Exercises {
		originalExerciseMap[exercise.ExerciseId] = exercise
	}

	// Start with original exercises and update only those recommended by AI
	var updatedExercises []*pb.WorkoutExercise

	// Copy all original exercises first
	for _, originalExercise := range originalWorkout.Exercises {
		// Copy the original exercise
		updatedExercise := &pb.WorkoutExercise{
			ExerciseId:  originalExercise.ExerciseId,
			Sets:        make([]*pb.WorkoutSet, len(originalExercise.Sets)),
			Notes:       originalExercise.Notes,
			RestSeconds: originalExercise.RestSeconds,
		}

		// Copy original sets
		for i, originalSet := range originalExercise.Sets {
			updatedExercise.Sets[i] = &pb.WorkoutSet{
				Reps:            originalSet.Reps,
				Weight:          originalSet.Weight,
				DurationSeconds: originalSet.DurationSeconds,
				Distance:        originalSet.Distance,
				Notes:           originalSet.Notes,
			}
		}

		updatedExercises = append(updatedExercises, updatedExercise)
	}

	// Create map for AI exercise details to store separately
	aiExerciseDetails := make(map[string]models.AIProgressiveOverloadExercise)

	// Apply AI recommendations
	for _, aiExercise := range aiResponse.UpdatedExercises {
		// Validate exercise ID exists in original workout
		if _, exists := originalExerciseMap[aiExercise.ExerciseID]; !exists {
			log.Printf("parseAIProgressiveOverloadResponse: Skipping invalid exercise ID '%s'", aiExercise.ExerciseID)
			continue
		}

		// Store AI exercise details for later storage
		var aiSets []models.WorkoutSet
		for _, aiSet := range aiExercise.Sets {
			aiSets = append(aiSets, models.WorkoutSet{
				Reps:            aiSet.Reps,
				Weight:          aiSet.Weight,
				DurationSeconds: int32(aiSet.DurationSeconds),
				Distance:        aiSet.Distance,
				Notes:           aiSet.Notes,
			})
		}

		aiExerciseDetails[aiExercise.ExerciseID] = models.AIProgressiveOverloadExercise{
			ExerciseID:         aiExercise.ExerciseID,
			ExerciseName:       aiExercise.ExerciseName,
			ProgressionApplied: aiExercise.ProgressionApplied,
			Reasoning:          aiExercise.Reasoning,
			ChangesMade:        aiExercise.ChangesMade,
			Sets:               aiSets,
			Notes:              aiExercise.Notes,
			RestSeconds:        aiExercise.RestSeconds,
		}

		// Find and update the corresponding exercise
		for i, exercise := range updatedExercises {
			if exercise.ExerciseId == aiExercise.ExerciseID {
				// Apply AI recommendations
				if aiExercise.ProgressionApplied {
					// Convert AI sets to protobuf sets
					var newSets []*pb.WorkoutSet
					for _, aiSet := range aiExercise.Sets {
						newSets = append(newSets, &pb.WorkoutSet{
							Reps:            aiSet.Reps,
							Weight:          aiSet.Weight,
							DurationSeconds: aiSet.DurationSeconds,
							Distance:        aiSet.Distance,
							Notes:           aiSet.Notes,
						})
					}

					updatedExercises[i].Sets = newSets
					if aiExercise.Notes != "" {
						updatedExercises[i].Notes = aiExercise.Notes
					}
					if aiExercise.RestSeconds > 0 {
						updatedExercises[i].RestSeconds = aiExercise.RestSeconds
					}

					log.Printf("parseAIProgressiveOverloadResponse: Applied progression to exercise %s: %s",
						aiExercise.ExerciseName, aiExercise.ChangesMade)
				}
				break
			}
		}
	}

	return updatedExercises, aiResponse.AnalysisSummary, aiExerciseDetails, nil
}

// storeAIProgressiveOverloadResponse stores the AI progressive overload response in MongoDB
func (s *WorkoutSessionService) storeAIProgressiveOverloadResponse(ctx context.Context, session *pb.WorkoutSession, workout *pb.Workout, analysisSummary string, updatedExercises []*pb.WorkoutExercise, success bool, message string, recentSessionsCount int, processingTimeMs int64) error {
	log.Printf("storeAIProgressiveOverloadResponse: Starting storage for session %s", session.Id)

	// Convert user, session, and workout IDs
	userObjectID, err := primitive.ObjectIDFromHex(session.UserId)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	sessionObjectID, err := primitive.ObjectIDFromHex(session.Id)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	workoutObjectID, err := primitive.ObjectIDFromHex(workout.Id)
	if err != nil {
		return fmt.Errorf("invalid workout ID: %w", err)
	}

	// Convert updated exercises to storage format
	var aiExercises []models.AIProgressiveOverloadExercise
	if updatedExercises != nil {
		for _, exercise := range updatedExercises {
			// Check if we have AI details for this exercise
			if aiDetail, exists := s.tempAIExerciseDetails[exercise.ExerciseId]; exists {
				// Use the detailed AI information
				aiExercises = append(aiExercises, aiDetail)
			} else {
				// Fallback to basic information
				var sets []models.WorkoutSet
				for _, set := range exercise.Sets {
					sets = append(sets, models.WorkoutSet{
						Reps:            set.Reps,
						Weight:          set.Weight,
						DurationSeconds: int32(set.DurationSeconds),
						Distance:        set.Distance,
						Notes:           set.Notes,
					})
				}

				// Try to find the exercise name from the workout
				exerciseName := "Unknown Exercise"
				for _, originalExercise := range workout.Exercises {
					if originalExercise.ExerciseId == exercise.ExerciseId {
						// Get exercise details if available
						if originalExercise.Exercise != nil {
							exerciseName = originalExercise.Exercise.Name
						}
						break
					}
				}

				aiExercises = append(aiExercises, models.AIProgressiveOverloadExercise{
					ExerciseID:         exercise.ExerciseId,
					ExerciseName:       exerciseName,
					ProgressionApplied: true,
					Reasoning:          "AI-generated progression based on recent performance analysis",
					ChangesMade:        "Updated based on AI analysis",
					Sets:               sets,
					Notes:              exercise.Notes,
					RestSeconds:        exercise.RestSeconds,
				})
			}
		}
	}

	// Clear temporary storage
	s.tempAIExerciseDetails = nil

	// Create the AI response document
	now := time.Now()
	aiResponse := models.AIProgressiveOverloadResponse{
		UserID:              userObjectID,
		WorkoutSessionID:    sessionObjectID,
		WorkoutID:           workoutObjectID,
		AnalysisSummary:     analysisSummary,
		UpdatedExercises:    aiExercises,
		Success:             success,
		Message:             message,
		RecentSessionsCount: int32(recentSessionsCount),
		AIModel:             "gemini-2.5-flash",
		ProcessingTimeMs:    processingTimeMs,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	// Insert into MongoDB
	result, err := s.aiProgressiveOverloadColl.InsertOne(ctx, aiResponse)
	if err != nil {
		return fmt.Errorf("failed to insert AI response: %w", err)
	}

	log.Printf("storeAIProgressiveOverloadResponse: Successfully stored AI response with ID: %v", result.InsertedID)
	return nil
}

// getFallbackProgressiveOverloadPromptTemplate returns a fallback prompt template
func (s *WorkoutSessionService) getFallbackProgressiveOverloadPromptTemplate() string {
	return `You are an expert strength training coach specializing in progressive overload optimization.

USER PROFILE:
- Goal: {{USER_GOAL}}
- Height: {{USER_HEIGHT}} cm, Weight: {{USER_WEIGHT}} kg, Age: {{USER_AGE}}

LAST WORKOUT SESSION DATA:
Session ID: {{SESSION_ID}}
Workout: {{WORKOUT_NAME}}
Duration: {{WORKOUT_DURATION}} minutes
Completed on: {{SESSION_DATE}}

{{SESSION_EXERCISES_DATA}}

TASK:
Analyze the completed workout session above and recent workout sessions data below to generate optimized progressive overload recommendations for the corresponding routine.

Apply intelligent progressive overload based on actual performance trends from recent sessions, RPE ratings, and completion patterns.
Be conservative with weight increases for safety and focus on sustainable long-term progression.
Use recent sessions to identify patterns - is the user consistently completing sets? Are they ready for progression or need consolidation?

OUTPUT FORMAT (JSON):
{
  "analysis_summary": "2-3 sentence summary of overall session performance, recent trends, and readiness for progression",
  "updated_exercises": [
    {
      "exercise_id": "ACTUAL_EXERCISE_ID_FROM_SESSION",
      "exercise_name": "Exercise Name",
      "progression_applied": true,
      "reasoning": "Specific reason for the progression/maintenance decision based on recent performance trends",
      "changes_made": "Description of what changed (e.g., 'Increased weight from 50kg to 52.5kg based on consistent completion')",
      "sets": [
        {
          "reps": 8,
          "weight": 52.5,
          "duration_seconds": 0,
          "distance": 0,
          "notes": ""
        }
      ],
      "notes": "Any additional notes about the progression based on recent trends",
      "rest_seconds": 90
    }
  ]
}

Only modify exercises that were actually performed in the session and use exact exercise IDs from the session data.
Consider recent performance patterns, RPE trends, and consistency when making progression decisions.`
}

// modelToProto converts a workout session model to protobuf workout session
func (s *WorkoutSessionService) modelToProto(session *models.WorkoutSession) *pb.WorkoutSession {
	var exercises []*pb.WorkoutSessionExercise
	for _, exercise := range session.Exercises {
		var sets []*pb.WorkoutSessionSet
		for _, set := range exercise.Sets {
			var startedAt, finishedAt *timestamppb.Timestamp
			if set.StartedAt != nil {
				startedAt = timestamppb.New(*set.StartedAt)
			}
			if set.FinishedAt != nil {
				finishedAt = timestamppb.New(*set.FinishedAt)
			}

			sets = append(sets, &pb.WorkoutSessionSet{
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
			})
		}

		var exerciseStartedAt, exerciseFinishedAt *timestamppb.Timestamp
		if exercise.StartedAt != nil {
			exerciseStartedAt = timestamppb.New(*exercise.StartedAt)
		}
		if exercise.FinishedAt != nil {
			exerciseFinishedAt = timestamppb.New(*exercise.FinishedAt)
		}

		// Get exercise details
		exerciseDetails, _ := s.exerciseService.GetExercise(context.Background(), &pb.GetExerciseRequest{Id: exercise.ExerciseID.Hex()})

		exercises = append(exercises, &pb.WorkoutSessionExercise{
			ExerciseId:  exercise.ExerciseID.Hex(),
			Exercise:    exerciseDetails,
			Sets:        sets,
			Notes:       exercise.Notes,
			RestSeconds: exercise.RestSeconds,
			Completed:   exercise.Completed,
			StartedAt:   exerciseStartedAt,
			FinishedAt:  exerciseFinishedAt,
		})
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
