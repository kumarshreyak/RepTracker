package services

import (
	"context"
	"regexp"
	"strings"
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

// ExerciseService implements the gRPC ExerciseService
type ExerciseService struct {
	pb.UnimplementedExerciseServiceServer
	db           *database.MongoDB
	exerciseColl *mongo.Collection
}

// NewExerciseService creates a new ExerciseService instance
func NewExerciseService(db *database.MongoDB) *ExerciseService {
	service := &ExerciseService{
		db:           db,
		exerciseColl: db.GetCollection("exercises"),
	}
	service.initializeExercises()
	return service
}

// initializeExercises populates the service with sample exercises if collection is empty
func (s *ExerciseService) initializeExercises() {
	ctx := context.Background()

	// Check if collection is empty
	count, err := s.exerciseColl.CountDocuments(ctx, bson.M{})
	if err != nil || count > 0 {
		return // Already has data or error checking
	}

	now := time.Now()
	sampleExercises := []models.Exercise{
		{
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

	// Insert sample exercises
	var docs []interface{}
	for _, exercise := range sampleExercises {
		docs = append(docs, exercise)
	}
	s.exerciseColl.InsertMany(ctx, docs)
}

// CreateExercise creates a new exercise
func (s *ExerciseService) CreateExercise(ctx context.Context, req *pb.CreateExerciseRequest) (*pb.Exercise, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "exercise name is required")
	}

	now := time.Now()
	exercise := models.Exercise{
		Name:         req.Name,
		Description:  req.Description,
		MuscleGroup:  req.MuscleGroup,
		Equipment:    req.Equipment,
		Difficulty:   req.Difficulty,
		Instructions: req.Instructions,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result, err := s.exerciseColl.InsertOne(ctx, exercise)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create exercise")
	}

	exercise.ID = result.InsertedID.(primitive.ObjectID)
	return s.modelToProto(&exercise), nil
}

// GetExercise retrieves an exercise by ID
func (s *ExerciseService) GetExercise(ctx context.Context, req *pb.GetExerciseRequest) (*pb.Exercise, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise ID")
	}

	var exercise models.Exercise
	err = s.exerciseColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&exercise)
	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "exercise not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "failed to get exercise")
	}

	return s.modelToProto(&exercise), nil
}

// UpdateExercise updates an existing exercise
func (s *ExerciseService) UpdateExercise(ctx context.Context, req *pb.UpdateExerciseRequest) (*pb.Exercise, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise ID")
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
	if req.MuscleGroup != "" {
		update["muscle_group"] = req.MuscleGroup
	}
	if req.Equipment != "" {
		update["equipment"] = req.Equipment
	}
	if req.Difficulty != "" {
		update["difficulty"] = req.Difficulty
	}
	if len(req.Instructions) > 0 {
		update["instructions"] = req.Instructions
	}

	result, err := s.exerciseColl.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update exercise")
	}
	if result.MatchedCount == 0 {
		return nil, status.Error(codes.NotFound, "exercise not found")
	}

	return s.GetExercise(ctx, &pb.GetExerciseRequest{Id: req.Id})
}

// DeleteExercise deletes an exercise
func (s *ExerciseService) DeleteExercise(ctx context.Context, req *pb.DeleteExerciseRequest) (*emptypb.Empty, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise ID")
	}

	result, err := s.exerciseColl.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete exercise")
	}
	if result.DeletedCount == 0 {
		return nil, status.Error(codes.NotFound, "exercise not found")
	}

	return &emptypb.Empty{}, nil
}

// ListExercises lists exercises with optional filtering
func (s *ExerciseService) ListExercises(ctx context.Context, req *pb.ListExercisesRequest) (*pb.ListExercisesResponse, error) {
	filter := bson.M{}

	// Apply search filter (searches across multiple fields with flexible word matching)
	if req.Search != "" {
		// Normalize search query: remove special chars, extra spaces, make lowercase
		normalizedSearch := regexp.MustCompile(`[^\w\s]`).ReplaceAllString(req.Search, "")
		normalizedSearch = regexp.MustCompile(`\s+`).ReplaceAllString(normalizedSearch, " ")
		normalizedSearch = strings.TrimSpace(strings.ToLower(normalizedSearch))

		if normalizedSearch != "" {
			// Split into individual words for flexible matching
			searchWords := strings.Fields(normalizedSearch)

			// Create regex patterns for each word
			var wordPatterns []string
			for _, word := range searchWords {
				// Escape special regex characters
				escapedWord := regexp.QuoteMeta(word)
				wordPatterns = append(wordPatterns, escapedWord)
			}

			// Create pattern that matches all words in any order with flexible separators
			// This will match "push" and "up" in "Push-ups", "Push ups", "pushups", etc.
			searchPattern := "(?i)"
			for i, pattern := range wordPatterns {
				if i > 0 {
					// Allow any non-letter characters or no characters between words
					searchPattern += "[^a-zA-Z]*"
				}
				searchPattern += pattern
			}

			// Search across multiple fields
			searchRegex := bson.M{"$regex": searchPattern, "$options": "i"}
			filter["$or"] = []bson.M{
				{"name": searchRegex},
				{"description": searchRegex},
				{"muscle_group": searchRegex},
				{"equipment": searchRegex},
			}
		}
	}

	// Apply filters
	if req.MuscleGroup != "" {
		filter["muscle_group"] = bson.M{"$regex": req.MuscleGroup, "$options": "i"}
	}
	if req.Equipment != "" {
		filter["equipment"] = bson.M{"$regex": req.Equipment, "$options": "i"}
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

	opts := options.Find().SetLimit(pageSize).SetSkip(skip)
	cursor, err := s.exerciseColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list exercises")
	}
	defer cursor.Close(ctx)

	var exercises []models.Exercise
	if err := cursor.All(ctx, &exercises); err != nil {
		return nil, status.Error(codes.Internal, "failed to decode exercises")
	}

	var protoExercises []*pb.Exercise
	for _, exercise := range exercises {
		protoExercises = append(protoExercises, s.modelToProto(&exercise))
	}

	response := &pb.ListExercisesResponse{
		Exercises: protoExercises,
	}

	// Set next page token if there are more results
	if int64(len(exercises)) == pageSize {
		response.NextPageToken = "next_page" // In real app, encode the next skip position
	}

	return response, nil
}

// modelToProto converts an exercise model to protobuf exercise
func (s *ExerciseService) modelToProto(exercise *models.Exercise) *pb.Exercise {
	return &pb.Exercise{
		Id:           exercise.ID.Hex(),
		Name:         exercise.Name,
		Description:  exercise.Description,
		MuscleGroup:  exercise.MuscleGroup,
		Equipment:    exercise.Equipment,
		Difficulty:   exercise.Difficulty,
		Instructions: exercise.Instructions,
		CreatedAt:    timestamppb.New(exercise.CreatedAt),
		UpdatedAt:    timestamppb.New(exercise.UpdatedAt),
	}
}
