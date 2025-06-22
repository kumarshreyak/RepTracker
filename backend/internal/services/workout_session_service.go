package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	db              *database.MongoDB
	sessionColl     *mongo.Collection
	workoutService  *WorkoutService
	exerciseService *ExerciseService
}

// NewWorkoutSessionService creates a new WorkoutSessionService instance
func NewWorkoutSessionService(db *database.MongoDB, workoutService *WorkoutService, exerciseService *ExerciseService) *WorkoutSessionService {
	return &WorkoutSessionService{
		db:              db,
		sessionColl:     db.GetCollection("workout_sessions"),
		workoutService:  workoutService,
		exerciseService: exerciseService,
	}
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
