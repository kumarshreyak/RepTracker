package services

import (
	"context"
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

// WorkoutService implements the gRPC WorkoutService
type WorkoutService struct {
	pb.UnimplementedWorkoutServiceServer
	db          *database.MongoDB
	workoutColl *mongo.Collection
}

// NewWorkoutService creates a new WorkoutService instance
func NewWorkoutService(db *database.MongoDB) *WorkoutService {
	return &WorkoutService{
		db:          db,
		workoutColl: db.GetCollection("workouts"),
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

	// Convert user ID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	// Validate and convert exercises
	var workoutExercises []models.WorkoutExercise
	for i, exercise := range req.Exercises {
		if exercise.ExerciseId == "" {
			return nil, status.Errorf(codes.InvalidArgument, "exercise_id is required for exercise %d", i)
		}
		if len(exercise.Sets) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "at least one set is required for exercise %d", i)
		}

		exerciseObjectID, err := primitive.ObjectIDFromHex(exercise.ExerciseId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid exercise_id for exercise %d", i)
		}

		// Convert sets
		var workoutSets []models.WorkoutSet
		for j, set := range exercise.Sets {
			if set.Reps <= 0 {
				return nil, status.Errorf(codes.InvalidArgument, "reps must be positive for exercise %d, set %d", i, j)
			}
			if set.Weight < 0 {
				return nil, status.Errorf(codes.InvalidArgument, "weight cannot be negative for exercise %d, set %d", i, j)
			}

			workoutSets = append(workoutSets, models.WorkoutSet{
				Reps:            set.Reps,
				Weight:          set.Weight,
				DurationSeconds: int32(set.DurationSeconds),
				Distance:        set.Distance,
				Notes:           set.Notes,
			})
		}

		workoutExercises = append(workoutExercises, models.WorkoutExercise{
			ExerciseID:  exerciseObjectID,
			Sets:        workoutSets,
			Notes:       exercise.Notes,
			RestSeconds: exercise.RestSeconds,
		})
	}

	now := time.Now()
	var startedAt *time.Time
	if req.StartedAt != nil {
		t := req.StartedAt.AsTime()
		startedAt = &t
	}

	workout := models.Workout{
		UserID:      userObjectID,
		Name:        req.Name,
		Description: req.Description,
		Exercises:   workoutExercises,
		StartedAt:   startedAt,
		Notes:       req.Notes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result, err := s.workoutColl.InsertOne(ctx, workout)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create workout")
	}

	workout.ID = result.InsertedID.(primitive.ObjectID)
	return s.modelToProto(&workout), nil
}

// GetWorkout retrieves a workout by ID
func (s *WorkoutService) GetWorkout(ctx context.Context, req *pb.GetWorkoutRequest) (*pb.Workout, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "workout id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid workout ID")
	}

	var workout models.Workout
	err = s.workoutColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&workout)
	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "workout not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "failed to get workout")
	}

	return s.modelToProto(&workout), nil
}

// UpdateWorkout updates an existing workout
func (s *WorkoutService) UpdateWorkout(ctx context.Context, req *pb.UpdateWorkoutRequest) (*pb.Workout, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "workout id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid workout ID")
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
		// Validate and convert exercises
		var workoutExercises []models.WorkoutExercise
		for j, exercise := range req.Exercises {
			if exercise.ExerciseId == "" {
				return nil, status.Errorf(codes.InvalidArgument, "exercise_id is required for exercise %d", j)
			}
			if len(exercise.Sets) == 0 {
				return nil, status.Errorf(codes.InvalidArgument, "at least one set is required for exercise %d", j)
			}

			exerciseObjectID, err := primitive.ObjectIDFromHex(exercise.ExerciseId)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid exercise_id for exercise %d", j)
			}

			// Convert sets
			var workoutSets []models.WorkoutSet
			for k, set := range exercise.Sets {
				if set.Reps <= 0 {
					return nil, status.Errorf(codes.InvalidArgument, "reps must be positive for exercise %d, set %d", j, k)
				}
				if set.Weight < 0 {
					return nil, status.Errorf(codes.InvalidArgument, "weight cannot be negative for exercise %d, set %d", j, k)
				}

				workoutSets = append(workoutSets, models.WorkoutSet{
					Reps:            set.Reps,
					Weight:          set.Weight,
					DurationSeconds: int32(set.DurationSeconds),
					Distance:        set.Distance,
					Notes:           set.Notes,
				})
			}

			workoutExercises = append(workoutExercises, models.WorkoutExercise{
				ExerciseID:  exerciseObjectID,
				Sets:        workoutSets,
				Notes:       exercise.Notes,
				RestSeconds: exercise.RestSeconds,
			})
		}
		update["exercises"] = workoutExercises
	}
	if req.FinishedAt != nil {
		t := req.FinishedAt.AsTime()
		update["finished_at"] = &t
	}
	if req.DurationSeconds > 0 {
		update["duration_seconds"] = req.DurationSeconds
	}
	if req.Notes != "" {
		update["notes"] = req.Notes
	}

	result, err := s.workoutColl.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update workout")
	}
	if result.MatchedCount == 0 {
		return nil, status.Error(codes.NotFound, "workout not found")
	}

	return s.GetWorkout(ctx, &pb.GetWorkoutRequest{Id: req.Id})
}

// DeleteWorkout deletes a workout
func (s *WorkoutService) DeleteWorkout(ctx context.Context, req *pb.DeleteWorkoutRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "workout id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid workout ID")
	}

	result, err := s.workoutColl.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete workout")
	}
	if result.DeletedCount == 0 {
		return nil, status.Error(codes.NotFound, "workout not found")
	}

	return &emptypb.Empty{}, nil
}

// ListWorkouts lists workouts with optional filtering
func (s *WorkoutService) ListWorkouts(ctx context.Context, req *pb.ListWorkoutsRequest) (*pb.ListWorkoutsResponse, error) {
	filter := bson.M{}

	// Apply filters
	if req.UserId != "" {
		userObjectID, err := primitive.ObjectIDFromHex(req.UserId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid user_id")
		}
		filter["user_id"] = userObjectID
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

	opts := options.Find().SetLimit(pageSize).SetSkip(skip).SetSort(bson.D{{"updated_at", -1}})
	cursor, err := s.workoutColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list workouts")
	}
	defer cursor.Close(ctx)

	var workouts []models.Workout
	if err := cursor.All(ctx, &workouts); err != nil {
		return nil, status.Error(codes.Internal, "failed to decode workouts")
	}

	var protoWorkouts []*pb.Workout
	for _, workout := range workouts {
		protoWorkouts = append(protoWorkouts, s.modelToProto(&workout))
	}

	response := &pb.ListWorkoutsResponse{
		Workouts: protoWorkouts,
	}

	// Set next page token if there are more results
	if int64(len(workouts)) == pageSize {
		response.NextPageToken = "next_page" // In real app, encode the next skip position
	}

	return response, nil
}

// modelToProto converts a workout model to protobuf workout
func (s *WorkoutService) modelToProto(workout *models.Workout) *pb.Workout {
	var exercises []*pb.WorkoutExercise
	for _, exercise := range workout.Exercises {
		var sets []*pb.WorkoutSet
		for _, set := range exercise.Sets {
			sets = append(sets, &pb.WorkoutSet{
				Reps:            set.Reps,
				Weight:          set.Weight,
				DurationSeconds: float32(set.DurationSeconds),
				Distance:        set.Distance,
				Notes:           set.Notes,
			})
		}

		exercises = append(exercises, &pb.WorkoutExercise{
			ExerciseId:  exercise.ExerciseID.Hex(),
			Sets:        sets,
			Notes:       exercise.Notes,
			RestSeconds: exercise.RestSeconds,
		})
	}

	var startedAt *timestamppb.Timestamp
	if workout.StartedAt != nil {
		startedAt = timestamppb.New(*workout.StartedAt)
	}

	var finishedAt *timestamppb.Timestamp
	if workout.FinishedAt != nil {
		finishedAt = timestamppb.New(*workout.FinishedAt)
	}

	return &pb.Workout{
		Id:              workout.ID.Hex(),
		UserId:          workout.UserID.Hex(),
		Name:            workout.Name,
		Description:     workout.Description,
		Exercises:       exercises,
		StartedAt:       startedAt,
		FinishedAt:      finishedAt,
		DurationSeconds: workout.DurationSeconds,
		Notes:           workout.Notes,
		CreatedAt:       timestamppb.New(workout.CreatedAt),
		UpdatedAt:       timestamppb.New(workout.UpdatedAt),
	}
}
