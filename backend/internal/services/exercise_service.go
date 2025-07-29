package services

import (
	"context"
	"log"
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
			Name:             "Push-ups",
			Description:      "A classic upper body exercise that works the chest, shoulders, and triceps.",
			Category:         string(models.CategoryCalisthenics),
			Equipment:        []string{string(models.EquipmentNone)},
			PrimaryMuscles:   []string{string(models.MuscleChest)},
			SecondaryMuscles: []string{string(models.MuscleTriceps), string(models.MuscleShoulders)},
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
			Name:             "Squats",
			Description:      "A fundamental lower body exercise targeting the quadriceps, glutes, and hamstrings.",
			Category:         string(models.CategoryStrength),
			Equipment:        []string{string(models.EquipmentNone)},
			PrimaryMuscles:   []string{string(models.MuscleQuads)},
			SecondaryMuscles: []string{string(models.MuscleGlutes), string(models.MuscleHamstrings)},
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
			Name:             "Bench Press",
			Description:      "A compound exercise for building chest, shoulder, and tricep strength.",
			Category:         string(models.CategoryStrength),
			Equipment:        []string{string(models.EquipmentBarbell), string(models.EquipmentBench)},
			PrimaryMuscles:   []string{string(models.MuscleChest)},
			SecondaryMuscles: []string{string(models.MuscleTriceps), string(models.MuscleShoulders)},
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
			Name:             "Deadlifts",
			Description:      "A compound movement that works the entire posterior chain.",
			Category:         string(models.CategoryStrength),
			Equipment:        []string{string(models.EquipmentBarbell)},
			PrimaryMuscles:   []string{string(models.MuscleGlutes), string(models.MuscleHamstrings)},
			SecondaryMuscles: []string{string(models.MuscleLowerBack), string(models.MuscleTraps), string(models.MuscleQuads)},
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
			Name:             "Pull-ups",
			Description:      "An upper body exercise that primarily targets the latissimus dorsi.",
			Category:         string(models.CategoryCalisthenics),
			Equipment:        []string{string(models.EquipmentPullUpBar)},
			PrimaryMuscles:   []string{string(models.MuscleLats)},
			SecondaryMuscles: []string{string(models.MuscleBiceps), string(models.MuscleRhomboids)},
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
			Name:             "Shoulder Press",
			Description:      "An overhead pressing movement for shoulder and tricep development.",
			Category:         string(models.CategoryStrength),
			Equipment:        []string{string(models.EquipmentDumbbell)},
			PrimaryMuscles:   []string{string(models.MuscleShoulders)},
			SecondaryMuscles: []string{string(models.MuscleTriceps)},
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
			Name:             "Bicep Curls",
			Description:      "An isolation exercise targeting the biceps brachii.",
			Category:         string(models.CategoryStrength),
			Equipment:        []string{string(models.EquipmentDumbbell)},
			PrimaryMuscles:   []string{string(models.MuscleBiceps)},
			SecondaryMuscles: []string{},
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
			Name:             "Planks",
			Description:      "An isometric core exercise for building abdominal strength.",
			Category:         string(models.CategoryCalisthenics),
			Equipment:        []string{string(models.EquipmentNone)},
			PrimaryMuscles:   []string{string(models.MuscleAbs)},
			SecondaryMuscles: []string{string(models.MuscleObliques)},
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
			Name:             "Lunges",
			Description:      "A unilateral leg exercise that targets quads, glutes, and hamstrings.",
			Category:         string(models.CategoryStrength),
			Equipment:        []string{string(models.EquipmentNone)},
			PrimaryMuscles:   []string{string(models.MuscleQuads)},
			SecondaryMuscles: []string{string(models.MuscleGlutes), string(models.MuscleHamstrings)},
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
			Name:             "Rows",
			Description:      "A pulling exercise that targets the back muscles and biceps.",
			Category:         string(models.CategoryStrength),
			Equipment:        []string{string(models.EquipmentCable)},
			PrimaryMuscles:   []string{string(models.MuscleLats), string(models.MuscleRhomboids)},
			SecondaryMuscles: []string{string(models.MuscleBiceps), string(models.MuscleTraps)},
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
		Name:             req.Name,
		Description:      req.Description,
		Category:         req.Category,
		Equipment:        req.Equipment,
		PrimaryMuscles:   req.PrimaryMuscles,
		SecondaryMuscles: req.SecondaryMuscles,
		Instructions:     req.Instructions,
		Video:            req.Video,
		VariationsOn:     req.VariationsOn,
		VariationOn:      req.VariationOn,
		CreatedAt:        now,
		UpdatedAt:        now,
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
		"updatedAt": time.Now(),
	}

	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.Category != "" {
		update["category"] = req.Category
	}
	if len(req.Equipment) > 0 {
		update["equipment"] = req.Equipment
	}
	if len(req.PrimaryMuscles) > 0 {
		update["primaryMuscles"] = req.PrimaryMuscles
	}
	if len(req.SecondaryMuscles) > 0 {
		update["secondaryMuscles"] = req.SecondaryMuscles
	}
	if len(req.Instructions) > 0 {
		update["instructions"] = req.Instructions
	}
	if req.Video != "" {
		update["video"] = req.Video
	}
	if len(req.VariationsOn) > 0 {
		update["variationsOn"] = req.VariationsOn
	}
	if len(req.VariationOn) > 0 {
		update["variationOn"] = req.VariationOn
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
				{"category": searchRegex},
				{"primaryMuscles": searchRegex},
				{"secondaryMuscles": searchRegex},
				{"equipment": searchRegex},
			}
		}
	}

	// Apply filters
	if len(req.Categories) > 0 {
		filter["category"] = bson.M{"$in": req.Categories}
	}
	if len(req.Equipment) > 0 {
		filter["equipment"] = bson.M{"$in": req.Equipment}
	}
	if len(req.PrimaryMuscles) > 0 {
		filter["primaryMuscles"] = bson.M{"$in": req.PrimaryMuscles}
	}
	if len(req.SecondaryMuscles) > 0 {
		filter["secondaryMuscles"] = bson.M{"$in": req.SecondaryMuscles}
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

// GetQuickAddExercises returns a list of exercises for quick adding to routines
// Returns the last 5 exercises that were added by user in some routine or fallback to default exercises
func (s *ExerciseService) GetQuickAddExercises(ctx context.Context, req *pb.GetQuickAddExercisesRequest) (*pb.GetQuickAddExercisesResponse, error) {
	log.Printf("🏋️ GetQuickAddExercises: Starting - userID: %s, limit: %d", req.UserId, req.Limit)

	limit := int64(req.Limit)
	if limit <= 0 {
		limit = 5 // Default limit
	}

	var exercises []models.Exercise

	// If user_id is provided, try to get user's recent exercises from workouts
	if req.UserId != "" {
		log.Printf("🏋️ GetQuickAddExercises: Attempting to get user recent exercises")
		userObjectID, err := primitive.ObjectIDFromHex(req.UserId)
		if err == nil {
			exercises = s.getUserRecentExercises(ctx, userObjectID, limit)
			log.Printf("🏋️ GetQuickAddExercises: Found %d recent exercises for user", len(exercises))
		} else {
			log.Printf("⚠️ GetQuickAddExercises: Invalid user ID format: %v", err)
		}
	} else {
		log.Printf("🏋️ GetQuickAddExercises: No user ID provided, will use default exercises")
	}

	// If we don't have enough exercises from user history, fill with default exercises
	if int64(len(exercises)) < limit {
		log.Printf("🏋️ GetQuickAddExercises: Need %d more exercises, getting defaults", limit-int64(len(exercises)))
		defaultExercises := s.getDefaultQuickAddExercises(ctx, limit-int64(len(exercises)))
		log.Printf("🏋️ GetQuickAddExercises: Found %d default exercises", len(defaultExercises))

		// Merge without duplicates
		exerciseMap := make(map[string]bool)
		for _, ex := range exercises {
			exerciseMap[ex.ID.Hex()] = true
		}

		for _, defaultEx := range defaultExercises {
			if !exerciseMap[defaultEx.ID.Hex()] {
				exercises = append(exercises, defaultEx)
				if int64(len(exercises)) >= limit {
					break
				}
			}
		}
	}

	log.Printf("🏋️ GetQuickAddExercises: Final exercise count: %d", len(exercises))

	// Convert to protobuf
	var protoExercises []*pb.Exercise
	for _, exercise := range exercises {
		protoExercises = append(protoExercises, s.modelToProto(&exercise))
	}

	log.Printf("🏋️ GetQuickAddExercises: Returning %d exercises", len(protoExercises))
	return &pb.GetQuickAddExercisesResponse{
		Exercises: protoExercises,
	}, nil
}

// getUserRecentExercises gets the most recently used exercises by a user from their workouts
func (s *ExerciseService) getUserRecentExercises(ctx context.Context, userID primitive.ObjectID, limit int64) []models.Exercise {
	// Get workouts collection
	workoutColl := s.db.GetCollection("workouts")

	// Aggregation pipeline to get recent exercises from user's workouts
	pipeline := []bson.M{
		// Match workouts by user
		{"$match": bson.M{"user_id": userID}},
		// Sort by creation date (most recent first)
		{"$sort": bson.M{"created_at": -1}},
		// Limit to recent workouts
		{"$limit": 20},
		// Unwind exercises array
		{"$unwind": "$exercises"},
		// Group by exercise_id to get unique exercises with their latest usage
		{"$group": bson.M{
			"_id":       "$exercises.exercise_id",
			"last_used": bson.M{"$max": "$created_at"},
		}},
		// Sort by last used (most recent first)
		{"$sort": bson.M{"last_used": -1}},
		// Limit to requested number
		{"$limit": limit},
	}

	cursor, err := workoutColl.Aggregate(ctx, pipeline)
	if err != nil {
		return []models.Exercise{}
	}
	defer cursor.Close(ctx)

	var exerciseIDs []primitive.ObjectID
	for cursor.Next(ctx) {
		var result struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&result); err == nil {
			exerciseIDs = append(exerciseIDs, result.ID)
		}
	}

	if len(exerciseIDs) == 0 {
		return []models.Exercise{}
	}

	// Get the actual exercise documents
	filter := bson.M{"_id": bson.M{"$in": exerciseIDs}}
	exerciseCursor, err := s.exerciseColl.Find(ctx, filter)
	if err != nil {
		return []models.Exercise{}
	}
	defer exerciseCursor.Close(ctx)

	var exercises []models.Exercise
	if err := exerciseCursor.All(ctx, &exercises); err != nil {
		return []models.Exercise{}
	}

	// Sort exercises by the order they appeared in exerciseIDs (most recent first)
	exerciseMap := make(map[primitive.ObjectID]models.Exercise)
	for _, ex := range exercises {
		exerciseMap[ex.ID] = ex
	}

	var sortedExercises []models.Exercise
	for _, id := range exerciseIDs {
		if ex, exists := exerciseMap[id]; exists {
			sortedExercises = append(sortedExercises, ex)
		}
	}

	return sortedExercises
}

// getDefaultQuickAddExercises returns default exercises for quick add
func (s *ExerciseService) getDefaultQuickAddExercises(ctx context.Context, limit int64) []models.Exercise {
	// Default exercise names in order of preference
	defaultNames := []string{
		"Push-ups",
		"Squats",
		"Bench Press",
		"Shoulder Press",
		"Bicep Curls",
	}

	var exercises []models.Exercise

	for _, name := range defaultNames {
		if int64(len(exercises)) >= limit {
			break
		}

		var exercise models.Exercise
		err := s.exerciseColl.FindOne(ctx, bson.M{"name": bson.M{"$regex": "^" + regexp.QuoteMeta(name) + "$", "$options": "i"}}).Decode(&exercise)
		if err == nil {
			exercises = append(exercises, exercise)
		}
	}

	// If we still don't have enough, get any exercises up to the limit
	if int64(len(exercises)) < limit {
		remaining := limit - int64(len(exercises))

		// Get exercise IDs we already have
		var existingIDs []primitive.ObjectID
		for _, ex := range exercises {
			existingIDs = append(existingIDs, ex.ID)
		}

		filter := bson.M{}
		if len(existingIDs) > 0 {
			filter["_id"] = bson.M{"$nin": existingIDs}
		}

		opts := options.Find().SetLimit(remaining)
		cursor, err := s.exerciseColl.Find(ctx, filter, opts)
		if err == nil {
			defer cursor.Close(ctx)

			var additionalExercises []models.Exercise
			if cursor.All(ctx, &additionalExercises) == nil {
				exercises = append(exercises, additionalExercises...)
			}
		}
	}

	return exercises
}

// GetExerciseHistory retrieves the workout session history for a specific exercise
func (s *ExerciseService) GetExerciseHistory(ctx context.Context, req *pb.GetExerciseHistoryRequest) (*pb.GetExerciseHistoryResponse, error) {
	if req.ExerciseId == "" {
		return nil, status.Error(codes.InvalidArgument, "exercise_id is required")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Convert exercise ID to ObjectID for MongoDB queries
	exerciseObjectID, err := primitive.ObjectIDFromHex(req.ExerciseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise_id")
	}

	log.Printf("🏋️ GetExerciseHistory: Starting for exercise %s", req.ExerciseId)

	// First get the exercise name for the response
	var exercise models.Exercise
	err = s.exerciseColl.FindOne(ctx, bson.M{"_id": exerciseObjectID}).Decode(&exercise)
	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "exercise not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "failed to get exercise")
	}

	// Get workout sessions collection
	workoutSessionColl := s.db.GetCollection("workout_sessions")
	aiProgressiveOverloadColl := s.db.GetCollection("ai_progressive_overload_recommendations")

	// Convert user ID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	// Build filter for workout sessions
	filter := bson.M{
		"exercises.exercise_id": exerciseObjectID,
		"finished_at":           bson.M{"$exists": true, "$ne": nil}, // Only completed sessions
		"user_id":               userObjectID,                        // Filter by user (required)
	}

	// Query options - sort by finished_at descending, limit to last 3 sessions
	opts := options.Find().
		SetSort(bson.D{{Key: "finished_at", Value: -1}}).
		SetLimit(3)

	cursor, err := workoutSessionColl.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("GetExerciseHistory: Failed to query workout sessions: %v", err)
		return nil, status.Error(codes.Internal, "failed to query workout sessions")
	}
	defer cursor.Close(ctx)

	var sessions []models.WorkoutSession
	if err = cursor.All(ctx, &sessions); err != nil {
		log.Printf("GetExerciseHistory: Failed to decode workout sessions: %v", err)
		return nil, status.Error(codes.Internal, "failed to decode workout sessions")
	}

	log.Printf("🏋️ GetExerciseHistory: Found %d sessions with exercise %s", len(sessions), req.ExerciseId)

	var historyEntries []*pb.ExerciseHistoryEntry

	for _, session := range sessions {
		// Find the specific exercise in this session
		var sessionExercise *models.WorkoutSessionExercise
		for _, ex := range session.Exercises {
			if ex.ExerciseID == exerciseObjectID {
				sessionExercise = &ex
				break
			}
		}

		if sessionExercise == nil {
			continue // This shouldn't happen given our query, but be safe
		}

		// Convert session exercise to protobuf
		protoSessionExercise := s.workoutSessionExerciseToProto(sessionExercise, &exercise)

		// Create minimal session info for context
		protoSessionInfo := &pb.WorkoutSession{
			Id:         session.ID.Hex(),
			Name:       session.Name,
			FinishedAt: timestamppb.New(*session.FinishedAt),
		}

		// Create history entry
		historyEntry := &pb.ExerciseHistoryEntry{
			SessionExercise: protoSessionExercise,
			SessionInfo:     protoSessionInfo,
		}

		// Try to find corresponding AI progressive overload data
		aiFilter := bson.M{
			"workoutSessionId":            session.ID,
			"updatedExercises.exerciseId": req.ExerciseId, // Use string ID for AI data
		}

		var aiResponse models.AIProgressiveOverloadResponse
		err = aiProgressiveOverloadColl.FindOne(ctx, aiFilter).Decode(&aiResponse)
		if err == nil {
			// Found AI data, find the specific exercise
			for _, aiExercise := range aiResponse.UpdatedExercises {
				if aiExercise.ExerciseID == req.ExerciseId {
					protoAIExercise := s.aiExerciseToProto(&aiExercise)
					historyEntry.AiExercise = protoAIExercise
					break
				}
			}
		} else if err != mongo.ErrNoDocuments {
			log.Printf("GetExerciseHistory: Warning - failed to query AI data for session %s: %v", session.ID.Hex(), err)
		}

		historyEntries = append(historyEntries, historyEntry)
	}

	log.Printf("🏋️ GetExerciseHistory: Returning %d history entries", len(historyEntries))

	return &pb.GetExerciseHistoryResponse{
		History:      historyEntries,
		ExerciseName: exercise.Name,
	}, nil
}

// Helper function to convert workout session exercise to protobuf
func (s *ExerciseService) workoutSessionExerciseToProto(sessionExercise *models.WorkoutSessionExercise, exercise *models.Exercise) *pb.WorkoutSessionExercise {
	var protoSets []*pb.WorkoutSessionSet
	for _, set := range sessionExercise.Sets {
		protoSet := &pb.WorkoutSessionSet{
			TargetReps:      set.TargetReps,
			TargetWeight:    set.TargetWeight,
			ActualReps:      set.ActualReps,
			ActualWeight:    set.ActualWeight,
			DurationSeconds: float32(set.DurationSeconds),
			Distance:        set.Distance,
			Notes:           set.Notes,
			Completed:       set.Completed,
		}
		if set.StartedAt != nil {
			protoSet.StartedAt = timestamppb.New(*set.StartedAt)
		}
		if set.FinishedAt != nil {
			protoSet.FinishedAt = timestamppb.New(*set.FinishedAt)
		}
		protoSets = append(protoSets, protoSet)
	}

	protoExercise := &pb.WorkoutSessionExercise{
		ExerciseId:  sessionExercise.ExerciseID.Hex(),
		Exercise:    s.modelToProto(exercise),
		Sets:        protoSets,
		Notes:       sessionExercise.Notes,
		RestSeconds: sessionExercise.RestSeconds,
		Completed:   sessionExercise.Completed,
	}

	if sessionExercise.StartedAt != nil {
		protoExercise.StartedAt = timestamppb.New(*sessionExercise.StartedAt)
	}
	if sessionExercise.FinishedAt != nil {
		protoExercise.FinishedAt = timestamppb.New(*sessionExercise.FinishedAt)
	}

	return protoExercise
}

// Helper function to convert AI exercise to protobuf
func (s *ExerciseService) aiExerciseToProto(aiExercise *models.AIProgressiveOverloadExercise) *pb.AIProgressiveOverloadExercise {
	var protoSets []*pb.WorkoutSet
	for _, set := range aiExercise.Sets {
		protoSet := &pb.WorkoutSet{
			Reps:            set.Reps,
			Weight:          set.Weight,
			DurationSeconds: float32(set.DurationSeconds),
			Distance:        set.Distance,
			Notes:           set.Notes,
		}
		protoSets = append(protoSets, protoSet)
	}

	return &pb.AIProgressiveOverloadExercise{
		ExerciseId:         aiExercise.ExerciseID,
		ExerciseName:       aiExercise.ExerciseName,
		ProgressionApplied: aiExercise.ProgressionApplied,
		Reasoning:          aiExercise.Reasoning,
		ChangesMade:        aiExercise.ChangesMade,
		Sets:               protoSets,
		Notes:              aiExercise.Notes,
		RestSeconds:        aiExercise.RestSeconds,
	}
}

// modelToProto converts an exercise model to protobuf exercise
func (s *ExerciseService) modelToProto(exercise *models.Exercise) *pb.Exercise {
	return &pb.Exercise{
		Id:               exercise.ID.Hex(),
		Name:             exercise.Name,
		Description:      exercise.Description,
		Category:         exercise.Category,
		Equipment:        exercise.Equipment,
		PrimaryMuscles:   exercise.PrimaryMuscles,
		SecondaryMuscles: exercise.SecondaryMuscles,
		Instructions:     exercise.Instructions,
		Video:            exercise.Video,
		VariationsOn:     exercise.VariationsOn,
		VariationOn:      exercise.VariationOn,
		CreatedAt:        timestamppb.New(exercise.CreatedAt),
		UpdatedAt:        timestamppb.New(exercise.UpdatedAt),
	}
}
