package services

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gymlog-backend/proto/gymlog/v1"
)

// ExerciseService implements the gRPC ExerciseService
type ExerciseService struct {
	pb.UnimplementedExerciseServiceServer
	// In a real app, this would have a database connection
	exercises []*pb.Exercise
}

// NewExerciseService creates a new ExerciseService instance
func NewExerciseService() *ExerciseService {
	service := &ExerciseService{}
	service.initializeExercises()
	return service
}

// initializeExercises populates the service with sample exercises
// In a real app, this would load from a database
func (s *ExerciseService) initializeExercises() {
	now := timestamppb.Now()

	s.exercises = []*pb.Exercise{
		{
			Id:          "1",
			Name:        "Push-ups",
			Description: "A classic upper body exercise that works the chest, shoulders, and triceps.",
			MuscleGroup: "Chest",
			Equipment:   "Bodyweight",
			Difficulty:  "Beginner",
			Instructions: []string{
				"Start in plank position with hands under shoulders",
				"Lower your body until chest nearly touches floor",
				"Push back up to starting position",
				"Repeat for desired reps",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Id:          "2",
			Name:        "Squats",
			Description: "A fundamental lower body exercise targeting the quadriceps, glutes, and hamstrings.",
			MuscleGroup: "Legs",
			Equipment:   "Bodyweight",
			Difficulty:  "Beginner",
			Instructions: []string{
				"Stand with feet shoulder-width apart",
				"Lower your body as if sitting back into a chair",
				"Keep knees behind toes and chest up",
				"Return to starting position",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Id:          "3",
			Name:        "Bench Press",
			Description: "A compound exercise for building chest, shoulder, and tricep strength.",
			MuscleGroup: "Chest",
			Equipment:   "Barbell",
			Difficulty:  "Intermediate",
			Instructions: []string{
				"Lie on bench with feet flat on floor",
				"Grip barbell slightly wider than shoulder width",
				"Lower bar to chest with control",
				"Press bar back up to starting position",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Id:          "4",
			Name:        "Deadlifts",
			Description: "A compound movement that works the entire posterior chain.",
			MuscleGroup: "Back",
			Equipment:   "Barbell",
			Difficulty:  "Advanced",
			Instructions: []string{
				"Stand with feet hip-width apart, bar over mid-foot",
				"Hinge at hips and knees to grab bar",
				"Keep back straight and chest up",
				"Drive through heels to stand up tall",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Id:          "5",
			Name:        "Pull-ups",
			Description: "An upper body exercise that primarily targets the latissimus dorsi.",
			MuscleGroup: "Back",
			Equipment:   "Pull-up Bar",
			Difficulty:  "Intermediate",
			Instructions: []string{
				"Hang from bar with palms facing away",
				"Pull your body up until chin clears bar",
				"Lower yourself with control",
				"Repeat for desired reps",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Id:          "6",
			Name:        "Shoulder Press",
			Description: "An overhead pressing movement for shoulder and tricep development.",
			MuscleGroup: "Shoulders",
			Equipment:   "Dumbbells",
			Difficulty:  "Intermediate",
			Instructions: []string{
				"Stand with feet shoulder-width apart",
				"Hold dumbbells at shoulder height",
				"Press weights overhead until arms are extended",
				"Lower back to starting position",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Id:          "7",
			Name:        "Bicep Curls",
			Description: "An isolation exercise targeting the biceps brachii.",
			MuscleGroup: "Arms",
			Equipment:   "Dumbbells",
			Difficulty:  "Beginner",
			Instructions: []string{
				"Stand with dumbbells at your sides",
				"Keep elbows close to your body",
				"Curl weights up by contracting biceps",
				"Lower weights back to starting position",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Id:          "8",
			Name:        "Planks",
			Description: "An isometric core exercise for building abdominal strength.",
			MuscleGroup: "Core",
			Equipment:   "Bodyweight",
			Difficulty:  "Beginner",
			Instructions: []string{
				"Start in push-up position",
				"Lower to forearms, keeping body straight",
				"Hold position while engaging core",
				"Breathe normally throughout",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Id:          "9",
			Name:        "Lunges",
			Description: "A unilateral leg exercise that targets quads, glutes, and hamstrings.",
			MuscleGroup: "Legs",
			Equipment:   "Bodyweight",
			Difficulty:  "Beginner",
			Instructions: []string{
				"Stand with feet hip-width apart",
				"Step forward with one leg, lowering hips",
				"Keep front knee over ankle",
				"Push back to starting position",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Id:          "10",
			Name:        "Rows",
			Description: "A pulling exercise that targets the back muscles and biceps.",
			MuscleGroup: "Back",
			Equipment:   "Cable Machine",
			Difficulty:  "Intermediate",
			Instructions: []string{
				"Sit with feet on platform, knees slightly bent",
				"Grip handle with both hands",
				"Pull handle to lower chest, squeezing shoulder blades",
				"Slowly return to starting position",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// CreateExercise creates a new exercise
func (s *ExerciseService) CreateExercise(ctx context.Context, req *pb.CreateExerciseRequest) (*pb.Exercise, error) {
	// TODO: Implement validation
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "exercise name is required")
	}

	// Generate a simple ID (in real app, use UUID or database auto-increment)
	id := len(s.exercises) + 1
	now := timestamppb.Now()

	exercise := &pb.Exercise{
		Id:           string(rune(id)),
		Name:         req.Name,
		Description:  req.Description,
		MuscleGroup:  req.MuscleGroup,
		Equipment:    req.Equipment,
		Difficulty:   req.Difficulty,
		Instructions: req.Instructions,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.exercises = append(s.exercises, exercise)
	return exercise, nil
}

// GetExercise retrieves an exercise by ID
func (s *ExerciseService) GetExercise(ctx context.Context, req *pb.GetExerciseRequest) (*pb.Exercise, error) {
	for _, exercise := range s.exercises {
		if exercise.Id == req.Id {
			return exercise, nil
		}
	}
	return nil, status.Error(codes.NotFound, "exercise not found")
}

// UpdateExercise updates an existing exercise
func (s *ExerciseService) UpdateExercise(ctx context.Context, req *pb.UpdateExerciseRequest) (*pb.Exercise, error) {
	for i, exercise := range s.exercises {
		if exercise.Id == req.Id {
			// Update fields if provided
			if req.Name != "" {
				exercise.Name = req.Name
			}
			if req.Description != "" {
				exercise.Description = req.Description
			}
			if req.MuscleGroup != "" {
				exercise.MuscleGroup = req.MuscleGroup
			}
			if req.Equipment != "" {
				exercise.Equipment = req.Equipment
			}
			if req.Difficulty != "" {
				exercise.Difficulty = req.Difficulty
			}
			if len(req.Instructions) > 0 {
				exercise.Instructions = req.Instructions
			}
			exercise.UpdatedAt = timestamppb.Now()

			s.exercises[i] = exercise
			return exercise, nil
		}
	}
	return nil, status.Error(codes.NotFound, "exercise not found")
}

// DeleteExercise deletes an exercise
func (s *ExerciseService) DeleteExercise(ctx context.Context, req *pb.DeleteExerciseRequest) (*emptypb.Empty, error) {
	for i, exercise := range s.exercises {
		if exercise.Id == req.Id {
			// Remove exercise from slice
			s.exercises = append(s.exercises[:i], s.exercises[i+1:]...)
			return &emptypb.Empty{}, nil
		}
	}
	return nil, status.Error(codes.NotFound, "exercise not found")
}

// ListExercises lists exercises with optional filtering
func (s *ExerciseService) ListExercises(ctx context.Context, req *pb.ListExercisesRequest) (*pb.ListExercisesResponse, error) {
	var filteredExercises []*pb.Exercise

	// Apply filters
	for _, exercise := range s.exercises {
		include := true

		// Filter by muscle group
		if req.MuscleGroup != "" {
			if !strings.EqualFold(exercise.MuscleGroup, req.MuscleGroup) {
				include = false
			}
		}

		// Filter by equipment
		if req.Equipment != "" {
			if !strings.EqualFold(exercise.Equipment, req.Equipment) {
				include = false
			}
		}

		if include {
			filteredExercises = append(filteredExercises, exercise)
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
	if end > len(filteredExercises) {
		end = len(filteredExercises)
	}

	result := filteredExercises[start:end]

	response := &pb.ListExercisesResponse{
		Exercises: result,
	}

	// Set next page token if there are more results
	if end < len(filteredExercises) {
		response.NextPageToken = "next_page" // In real app, encode the next start position
	}

	return response, nil
}
