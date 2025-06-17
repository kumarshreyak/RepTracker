package services

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gymlog-backend/proto/gymlog/v1"
)

// WorkoutSessionService implements the gRPC WorkoutSessionService
type WorkoutSessionService struct {
	pb.UnimplementedWorkoutSessionServiceServer
	// In a real app, this would have a database connection
	sessions        []*pb.WorkoutSession
	nextID          int
	workoutService  *WorkoutService
	exerciseService *ExerciseService
}

// NewWorkoutSessionService creates a new WorkoutSessionService instance
func NewWorkoutSessionService(workoutService *WorkoutService, exerciseService *ExerciseService) *WorkoutSessionService {
	return &WorkoutSessionService{
		sessions:        []*pb.WorkoutSession{},
		nextID:          1,
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

	// Get the routine/workout template
	workoutReq := &pb.GetWorkoutRequest{Id: req.RoutineId}
	routine, err := s.workoutService.GetWorkout(ctx, workoutReq)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "routine not found: %v", err)
	}

	now := timestamppb.Now()
	id := strconv.Itoa(s.nextID)
	s.nextID++

	// Convert routine exercises to session exercises with target values
	sessionExercises := make([]*pb.WorkoutSessionExercise, len(routine.Exercises))
	for i, routineExercise := range routine.Exercises {
		// Get exercise details
		exerciseReq := &pb.GetExerciseRequest{Id: routineExercise.ExerciseId}
		exercise, err := s.exerciseService.GetExercise(ctx, exerciseReq)
		if err != nil {
			// Log warning but continue - exercise details are optional
			exercise = nil
		}

		// Convert sets to session sets with target and actual values
		sessionSets := make([]*pb.WorkoutSessionSet, len(routineExercise.Sets))
		for j, routineSet := range routineExercise.Sets {
			sessionSets[j] = &pb.WorkoutSessionSet{
				TargetReps:      routineSet.Reps,
				TargetWeight:    routineSet.Weight,
				ActualReps:      0, // To be filled during workout
				ActualWeight:    0, // To be filled during workout
				DurationSeconds: routineSet.DurationSeconds,
				Distance:        routineSet.Distance,
				Notes:           routineSet.Notes,
				Completed:       false,
				StartedAt:       nil,
				FinishedAt:      nil,
			}
		}

		sessionExercises[i] = &pb.WorkoutSessionExercise{
			ExerciseId:  routineExercise.ExerciseId,
			Exercise:    exercise,
			Sets:        sessionSets,
			Notes:       routineExercise.Notes,
			RestSeconds: routineExercise.RestSeconds,
			Completed:   false,
			StartedAt:   nil,
			FinishedAt:  nil,
		}
	}

	// Use provided name or generate one
	sessionName := req.Name
	if sessionName == "" {
		sessionName = fmt.Sprintf("%s Session - %s", routine.Name, now.AsTime().Format("Jan 2, 2006"))
	}

	session := &pb.WorkoutSession{
		Id:              id,
		UserId:          req.UserId,
		RoutineId:       req.RoutineId,
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

	s.sessions = append(s.sessions, session)
	return session, nil
}

// GetWorkoutSession retrieves a workout session by ID
func (s *WorkoutSessionService) GetWorkoutSession(ctx context.Context, req *pb.GetWorkoutSessionRequest) (*pb.WorkoutSession, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "session id is required")
	}

	for _, session := range s.sessions {
		if session.Id == req.Id {
			return session, nil
		}
	}
	return nil, status.Error(codes.NotFound, "workout session not found")
}

// UpdateWorkoutSession updates an existing workout session
func (s *WorkoutSessionService) UpdateWorkoutSession(ctx context.Context, req *pb.UpdateWorkoutSessionRequest) (*pb.WorkoutSession, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "session id is required")
	}

	for i, session := range s.sessions {
		if session.Id == req.Id {
			// Update fields if provided
			if req.Name != "" {
				session.Name = req.Name
			}
			if req.Description != "" {
				session.Description = req.Description
			}
			if len(req.Exercises) > 0 {
				session.Exercises = req.Exercises
			}
			if req.FinishedAt != nil {
				session.FinishedAt = req.FinishedAt
				session.IsActive = false

				// Calculate duration if session is finished
				if session.StartedAt != nil {
					duration := req.FinishedAt.AsTime().Sub(session.StartedAt.AsTime())
					session.DurationSeconds = int32(duration.Seconds())
				}
			}
			if req.DurationSeconds > 0 {
				session.DurationSeconds = req.DurationSeconds
			}
			if req.Notes != "" {
				session.Notes = req.Notes
			}
			if req.IsActive != session.IsActive {
				session.IsActive = req.IsActive
			}
			session.UpdatedAt = timestamppb.Now()

			s.sessions[i] = session
			return session, nil
		}
	}
	return nil, status.Error(codes.NotFound, "workout session not found")
}

// DeleteWorkoutSession deletes a workout session
func (s *WorkoutSessionService) DeleteWorkoutSession(ctx context.Context, req *pb.DeleteWorkoutSessionRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "session id is required")
	}

	for i, session := range s.sessions {
		if session.Id == req.Id {
			// Remove session from slice
			s.sessions = append(s.sessions[:i], s.sessions[i+1:]...)
			return &emptypb.Empty{}, nil
		}
	}
	return nil, status.Error(codes.NotFound, "workout session not found")
}

// ListWorkoutSessions lists workout sessions with optional filtering
func (s *WorkoutSessionService) ListWorkoutSessions(ctx context.Context, req *pb.ListWorkoutSessionsRequest) (*pb.ListWorkoutSessionsResponse, error) {
	var filteredSessions []*pb.WorkoutSession

	// Apply filters
	for _, session := range s.sessions {
		include := true

		// Filter by user_id
		if req.UserId != "" {
			if session.UserId != req.UserId {
				include = false
			}
		}

		// Filter by active status
		if req.ActiveOnly && !session.IsActive {
			include = false
		}

		// Filter by date range
		if req.StartDate != nil && session.StartedAt != nil {
			if session.StartedAt.AsTime().Before(req.StartDate.AsTime()) {
				include = false
			}
		}
		if req.EndDate != nil && session.StartedAt != nil {
			if session.StartedAt.AsTime().After(req.EndDate.AsTime()) {
				include = false
			}
		}

		if include {
			filteredSessions = append(filteredSessions, session)
		}
	}

	// Apply pagination
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 50 // Default page size
	}

	start := 0
	if req.PageToken != "" {
		// In a real app, decode the page token to get the start index
		// For simplicity, we're not implementing pagination here
	}

	end := start + pageSize
	if end > len(filteredSessions) {
		end = len(filteredSessions)
	}

	result := filteredSessions[start:end]

	response := &pb.ListWorkoutSessionsResponse{
		Sessions: result,
	}

	// Set next page token if there are more results
	if end < len(filteredSessions) {
		response.NextPageToken = "next_page" // In real app, encode the next start position
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
