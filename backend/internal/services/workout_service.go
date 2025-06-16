package services

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gymlog-backend/proto/gymlog/v1"
)

// WorkoutService implements the gRPC WorkoutService
type WorkoutService struct {
	pb.UnimplementedWorkoutServiceServer
	// In a real app, this would have a database connection
	workouts []*pb.Workout
	nextID   int
}

// NewWorkoutService creates a new WorkoutService instance
func NewWorkoutService() *WorkoutService {
	return &WorkoutService{
		workouts: []*pb.Workout{},
		nextID:   1,
	}
}

// CreateWorkout creates a new workout
func (s *WorkoutService) CreateWorkout(ctx context.Context, req *pb.CreateWorkoutRequest) (*pb.Workout, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "workout name is required")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if len(req.Exercises) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one exercise is required")
	}

	// Validate exercises
	for i, exercise := range req.Exercises {
		if exercise.ExerciseId == "" {
			return nil, status.Errorf(codes.InvalidArgument, "exercise_id is required for exercise %d", i)
		}
		if len(exercise.Sets) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "at least one set is required for exercise %d", i)
		}

		// Validate sets
		for j, set := range exercise.Sets {
			if set.Reps <= 0 {
				return nil, status.Errorf(codes.InvalidArgument, "reps must be positive for exercise %d, set %d", i, j)
			}
			if set.Weight < 0 {
				return nil, status.Errorf(codes.InvalidArgument, "weight cannot be negative for exercise %d, set %d", i, j)
			}
		}
	}

	now := timestamppb.Now()
	id := strconv.Itoa(s.nextID)
	s.nextID++

	// Set started_at to now if not provided
	startedAt := req.StartedAt
	if startedAt == nil {
		startedAt = now
	}

	workout := &pb.Workout{
		Id:          id,
		UserId:      req.UserId,
		Name:        req.Name,
		Description: req.Description,
		Exercises:   req.Exercises,
		StartedAt:   startedAt,
		Notes:       req.Notes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.workouts = append(s.workouts, workout)
	return workout, nil
}

// GetWorkout retrieves a workout by ID
func (s *WorkoutService) GetWorkout(ctx context.Context, req *pb.GetWorkoutRequest) (*pb.Workout, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "workout id is required")
	}

	for _, workout := range s.workouts {
		if workout.Id == req.Id {
			return workout, nil
		}
	}
	return nil, status.Error(codes.NotFound, "workout not found")
}

// UpdateWorkout updates an existing workout
func (s *WorkoutService) UpdateWorkout(ctx context.Context, req *pb.UpdateWorkoutRequest) (*pb.Workout, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "workout id is required")
	}

	for i, workout := range s.workouts {
		if workout.Id == req.Id {
			// Update fields if provided
			if req.Name != "" {
				workout.Name = req.Name
			}
			if req.Description != "" {
				workout.Description = req.Description
			}
			if len(req.Exercises) > 0 {
				// Validate exercises before updating
				for j, exercise := range req.Exercises {
					if exercise.ExerciseId == "" {
						return nil, status.Errorf(codes.InvalidArgument, "exercise_id is required for exercise %d", j)
					}
					if len(exercise.Sets) == 0 {
						return nil, status.Errorf(codes.InvalidArgument, "at least one set is required for exercise %d", j)
					}

					// Validate sets
					for k, set := range exercise.Sets {
						if set.Reps <= 0 {
							return nil, status.Errorf(codes.InvalidArgument, "reps must be positive for exercise %d, set %d", j, k)
						}
						if set.Weight < 0 {
							return nil, status.Errorf(codes.InvalidArgument, "weight cannot be negative for exercise %d, set %d", j, k)
						}
					}
				}
				workout.Exercises = req.Exercises
			}
			if req.FinishedAt != nil {
				workout.FinishedAt = req.FinishedAt
			}
			if req.DurationSeconds > 0 {
				workout.DurationSeconds = req.DurationSeconds
			}
			if req.Notes != "" {
				workout.Notes = req.Notes
			}
			workout.UpdatedAt = timestamppb.Now()

			s.workouts[i] = workout
			return workout, nil
		}
	}
	return nil, status.Error(codes.NotFound, "workout not found")
}

// DeleteWorkout deletes a workout
func (s *WorkoutService) DeleteWorkout(ctx context.Context, req *pb.DeleteWorkoutRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "workout id is required")
	}

	for i, workout := range s.workouts {
		if workout.Id == req.Id {
			// Remove workout from slice
			s.workouts = append(s.workouts[:i], s.workouts[i+1:]...)
			return &emptypb.Empty{}, nil
		}
	}
	return nil, status.Error(codes.NotFound, "workout not found")
}

// ListWorkouts lists workouts with optional filtering
func (s *WorkoutService) ListWorkouts(ctx context.Context, req *pb.ListWorkoutsRequest) (*pb.ListWorkoutsResponse, error) {
	var filteredWorkouts []*pb.Workout

	// Apply filters
	for _, workout := range s.workouts {
		include := true

		// Filter by user_id
		if req.UserId != "" {
			if workout.UserId != req.UserId {
				include = false
			}
		}

		// Filter by date range
		if req.StartDate != nil && workout.StartedAt != nil {
			if workout.StartedAt.AsTime().Before(req.StartDate.AsTime()) {
				include = false
			}
		}
		if req.EndDate != nil && workout.StartedAt != nil {
			if workout.StartedAt.AsTime().After(req.EndDate.AsTime()) {
				include = false
			}
		}

		if include {
			filteredWorkouts = append(filteredWorkouts, workout)
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
	if end > len(filteredWorkouts) {
		end = len(filteredWorkouts)
	}

	result := filteredWorkouts[start:end]

	response := &pb.ListWorkoutsResponse{
		Workouts: result,
	}

	// Set next page token if there are more results
	if end < len(filteredWorkouts) {
		response.NextPageToken = "next_page" // In real app, encode the next start position
	}

	return response, nil
}
