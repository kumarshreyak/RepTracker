package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"gymlog-backend/internal/services"
	"gymlog-backend/pkg/database"
	"gymlog-backend/pkg/models"
	pb "gymlog-backend/proto/gymlog/v1"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	userService           *services.UserService
	exerciseService       *services.ExerciseService
	workoutService        *services.WorkoutService
	workoutSessionService *services.WorkoutSessionService
	metricsService        *services.MetricsService
	insightsService       *services.InsightsService
	suggestionsService    *services.SuggestionsService
	db                    *database.MongoDB
}

type GoogleUserRequest struct {
	GoogleID     string `json:"googleId"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Picture      string `json:"picture"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresAt    string `json:"expiresAt,omitempty"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Height    float64   `json:"height"` // Height in cm
	Weight    float64   `json:"weight"` // Weight in kg
	Age       int32     `json:"age"`    // Age in years
	Goal      string    `json:"goal"`   // Goal: lose_fat, gain_muscle, maintain
	GoogleID  string    `json:"googleId"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateUserRequest struct {
	FirstName string  `json:"firstName,omitempty"`
	LastName  string  `json:"lastName,omitempty"`
	Height    float64 `json:"height,omitempty"` // Height in cm
	Weight    float64 `json:"weight,omitempty"` // Weight in kg
	Age       int32   `json:"age,omitempty"`    // Age in years
	Goal      string  `json:"goal,omitempty"`   // Goal: lose_fat, gain_muscle, maintain
}

// Exercise HTTP types
type ExerciseResponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Category         string    `json:"category"`
	Equipment        []string  `json:"equipment"`
	PrimaryMuscles   []string  `json:"primaryMuscles"`
	SecondaryMuscles []string  `json:"secondaryMuscles"`
	Instructions     []string  `json:"instructions"`
	Video            string    `json:"video,omitempty"`
	VariationsOn     []string  `json:"variationsOn,omitempty"`
	VariationOn      []string  `json:"variationOn,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CreateExerciseRequest struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Category         string   `json:"category"`
	Equipment        []string `json:"equipment"`
	PrimaryMuscles   []string `json:"primaryMuscles"`
	SecondaryMuscles []string `json:"secondaryMuscles"`
	Instructions     []string `json:"instructions"`
	Video            string   `json:"video,omitempty"`
	VariationsOn     []string `json:"variationsOn,omitempty"`
	VariationOn      []string `json:"variationOn,omitempty"`
}

type UpdateExerciseRequest struct {
	Name             string   `json:"name,omitempty"`
	Description      string   `json:"description,omitempty"`
	Category         string   `json:"category,omitempty"`
	Equipment        []string `json:"equipment,omitempty"`
	PrimaryMuscles   []string `json:"primaryMuscles,omitempty"`
	SecondaryMuscles []string `json:"secondaryMuscles,omitempty"`
	Instructions     []string `json:"instructions,omitempty"`
	Video            string   `json:"video,omitempty"`
	VariationsOn     []string `json:"variationsOn,omitempty"`
	VariationOn      []string `json:"variationOn,omitempty"`
}

type ListExercisesResponse struct {
	Exercises     []ExerciseResponse `json:"exercises"`
	NextPageToken string             `json:"nextPageToken,omitempty"`
}

type QuickAddExercisesResponse struct {
	Exercises []ExerciseResponse `json:"exercises"`
}

// Workout HTTP types
type WorkoutSetResponse struct {
	Reps            int32   `json:"reps"`
	Weight          float32 `json:"weight"`
	DurationSeconds float32 `json:"durationSeconds"`
	Distance        float32 `json:"distance"`
	Notes           string  `json:"notes"`
}

type WorkoutExerciseResponse struct {
	ExerciseID  string               `json:"exerciseId"`
	Exercise    *ExerciseResponse    `json:"exercise,omitempty"`
	Sets        []WorkoutSetResponse `json:"sets"`
	Notes       string               `json:"notes"`
	RestSeconds int32                `json:"restSeconds"`
}

type WorkoutResponse struct {
	ID              string                    `json:"id"`
	UserID          string                    `json:"userId"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	Exercises       []WorkoutExerciseResponse `json:"exercises"`
	StartedAt       *time.Time                `json:"startedAt,omitempty"`
	FinishedAt      *time.Time                `json:"finishedAt,omitempty"`
	DurationSeconds int32                     `json:"durationSeconds"`
	Notes           string                    `json:"notes"`
	CreatedAt       time.Time                 `json:"createdAt"`
	UpdatedAt       time.Time                 `json:"updatedAt"`
}

type CreateWorkoutRequest struct {
	UserID      string            `json:"userId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exercises   []WorkoutExercise `json:"exercises"`
	StartedAt   *time.Time        `json:"startedAt,omitempty"`
	Notes       string            `json:"notes"`
}

type WorkoutSet struct {
	Reps            int32   `json:"reps"`
	Weight          float32 `json:"weight"`
	DurationSeconds float32 `json:"durationSeconds"`
	Distance        float32 `json:"distance"`
	Notes           string  `json:"notes"`
}

type WorkoutExercise struct {
	ExerciseID  string       `json:"exerciseId"`
	Sets        []WorkoutSet `json:"sets"`
	Notes       string       `json:"notes"`
	RestSeconds int32        `json:"restSeconds"`
}

type UpdateWorkoutRequest struct {
	Name            string            `json:"name,omitempty"`
	Description     string            `json:"description,omitempty"`
	Exercises       []WorkoutExercise `json:"exercises,omitempty"`
	FinishedAt      *time.Time        `json:"finishedAt,omitempty"`
	DurationSeconds int32             `json:"durationSeconds,omitempty"`
	Notes           string            `json:"notes,omitempty"`
}

type ListWorkoutsResponse struct {
	Workouts      []WorkoutResponse `json:"workouts"`
	NextPageToken string            `json:"nextPageToken,omitempty"`
}

// Workout Session HTTP types
type WorkoutSessionSetResponse struct {
	TargetReps      int32      `json:"targetReps"`
	TargetWeight    float32    `json:"targetWeight"`
	ActualReps      int32      `json:"actualReps"`
	ActualWeight    float32    `json:"actualWeight"`
	DurationSeconds float32    `json:"durationSeconds"`
	Distance        float32    `json:"distance"`
	Notes           string     `json:"notes"`
	Completed       bool       `json:"completed"`
	StartedAt       *time.Time `json:"startedAt,omitempty"`
	FinishedAt      *time.Time `json:"finishedAt,omitempty"`
}

type WorkoutSessionExerciseResponse struct {
	ExerciseID  string                      `json:"exerciseId"`
	Exercise    *ExerciseResponse           `json:"exercise,omitempty"`
	Sets        []WorkoutSessionSetResponse `json:"sets"`
	Notes       string                      `json:"notes"`
	RestSeconds int32                       `json:"restSeconds"`
	Completed   bool                        `json:"completed"`
	StartedAt   *time.Time                  `json:"startedAt,omitempty"`
	FinishedAt  *time.Time                  `json:"finishedAt,omitempty"`
}

type WorkoutSessionResponse struct {
	ID              string                           `json:"id"`
	UserID          string                           `json:"userId"`
	RoutineID       string                           `json:"routineId"`
	Name            string                           `json:"name"`
	Description     string                           `json:"description"`
	Exercises       []WorkoutSessionExerciseResponse `json:"exercises"`
	StartedAt       *time.Time                       `json:"startedAt,omitempty"`
	FinishedAt      *time.Time                       `json:"finishedAt,omitempty"`
	DurationSeconds int32                            `json:"durationSeconds"`
	Notes           string                           `json:"notes"`
	RPERating       int32                            `json:"rpeRating,omitempty"`
	IsActive        bool                             `json:"isActive"`
	CreatedAt       time.Time                        `json:"createdAt"`
	UpdatedAt       time.Time                        `json:"updatedAt"`
}

type CreateWorkoutSessionRequest struct {
	UserID      string `json:"userId"`
	RoutineID   string `json:"routineId"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Notes       string `json:"notes,omitempty"`
}

type UpdateWorkoutSessionRequest struct {
	Name            string                           `json:"name,omitempty"`
	Description     string                           `json:"description,omitempty"`
	Exercises       []WorkoutSessionExerciseResponse `json:"exercises,omitempty"`
	FinishedAt      *time.Time                       `json:"finishedAt,omitempty"`
	DurationSeconds int32                            `json:"durationSeconds,omitempty"`
	Notes           string                           `json:"notes,omitempty"`
	RPERating       int32                            `json:"rpeRating,omitempty"`
	IsActive        *bool                            `json:"isActive,omitempty"`
}

type ListWorkoutSessionsResponse struct {
	Sessions      []WorkoutSessionResponse `json:"sessions"`
	NextPageToken string                   `json:"nextPageToken,omitempty"`
}

type UpdateSetRequest struct {
	ActualReps      int32   `json:"actualReps"`
	ActualWeight    float32 `json:"actualWeight"`
	DurationSeconds float32 `json:"durationSeconds,omitempty"`
	Distance        float32 `json:"distance,omitempty"`
	Notes           string  `json:"notes,omitempty"`
	Completed       bool    `json:"completed"`
}

// Insights types
type WorkoutInsightResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Type      string    `json:"type"`
	Insight   string    `json:"insight"`
	BasedOn   string    `json:"basedOn"`
	Priority  int32     `json:"priority"`
	CreatedAt time.Time `json:"createdAt"`
}

type GenerateInsightsResponse struct {
	Insights []WorkoutInsightResponse `json:"insights"`
}

type GetRecentInsightsResponse struct {
	Insights []WorkoutInsightResponse `json:"insights"`
}

// Suggestions types
type WorkoutChangeResponse struct {
	Type         string `json:"type"`
	ExerciseID   string `json:"exerciseId,omitempty"`
	ExerciseName string `json:"exerciseName"`
	OldValue     string `json:"oldValue"`
	NewValue     string `json:"newValue"`
	Reason       string `json:"reason"`
}

type SuggestedWorkoutResponse struct {
	OriginalWorkoutID string                    `json:"originalWorkoutId"`
	Name              string                    `json:"name"`
	Description       string                    `json:"description"`
	Exercises         []WorkoutExerciseResponse `json:"exercises"`
	Changes           []WorkoutChangeResponse   `json:"changes"`
	OverallReasoning  string                    `json:"overallReasoning"`
	Priority          int32                     `json:"priority"`
}

type GenerateWorkoutSuggestionsRequest struct {
	DaysToAnalyze  int32 `json:"daysToAnalyze,omitempty"`
	MaxSuggestions int32 `json:"maxSuggestions,omitempty"`
}

type GenerateWorkoutSuggestionsResponse struct {
	Suggestions     []SuggestedWorkoutResponse `json:"suggestions"`
	AnalysisSummary string                     `json:"analysisSummary"`
}

type StoredWorkoutSuggestionResponse struct {
	ID              string                     `json:"id"`
	UserID          string                     `json:"userId"`
	Suggestions     []SuggestedWorkoutResponse `json:"suggestions"`
	AnalysisSummary string                     `json:"analysisSummary"`
	DaysAnalyzed    int32                      `json:"daysAnalyzed"`
	CreatedAt       time.Time                  `json:"createdAt"`
	UpdatedAt       time.Time                  `json:"updatedAt"`
}

type GetStoredSuggestionsResponse struct {
	StoredSuggestions []StoredWorkoutSuggestionResponse `json:"storedSuggestions"`
	NextPageToken     string                            `json:"nextPageToken,omitempty"`
}

func NewServer(userService *services.UserService, exerciseService *services.ExerciseService, workoutService *services.WorkoutService, workoutSessionService *services.WorkoutSessionService, metricsService *services.MetricsService, insightsService *services.InsightsService, suggestionsService *services.SuggestionsService, db *database.MongoDB) *Server {
	return &Server{
		userService:           userService,
		exerciseService:       exerciseService,
		workoutService:        workoutService,
		workoutSessionService: workoutSessionService,
		metricsService:        metricsService,
		insightsService:       insightsService,
		suggestionsService:    suggestionsService,
		db:                    db,
	}
}

func (s *Server) Start(port string) error {
	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", s.handleHealthCheck).Methods("GET")
	r.HandleFunc("/", s.handleHealthCheck).Methods("GET") // Root endpoint for basic checks

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/google", s.handleGoogleAuth).Methods("POST")
	api.HandleFunc("/auth/validate", s.handleValidateSession).Methods("GET")
	api.HandleFunc("/auth/logout", s.handleLogout).Methods("POST")

	// User routes
	api.HandleFunc("/users/{userId}", s.handleUpdateUser).Methods("PUT")

	// Exercise routes - specific routes must come before parameterized routes
	api.HandleFunc("/exercises/quick-add", s.handleGetQuickAddExercises).Methods("GET")
	api.HandleFunc("/exercises", s.handleListExercises).Methods("GET")
	api.HandleFunc("/exercises", s.handleCreateExercise).Methods("POST")
	api.HandleFunc("/exercises/{id}", s.handleGetExercise).Methods("GET")
	api.HandleFunc("/exercises/{id}", s.handleUpdateExercise).Methods("PUT")
	api.HandleFunc("/exercises/{id}", s.handleDeleteExercise).Methods("DELETE")

	// Workout routes
	api.HandleFunc("/workouts", s.handleListWorkouts).Methods("GET")
	api.HandleFunc("/workouts", s.handleCreateWorkout).Methods("POST")
	api.HandleFunc("/workouts/{id}", s.handleGetWorkout).Methods("GET")
	api.HandleFunc("/workouts/{id}/start", s.handleStartWorkout).Methods("GET")
	api.HandleFunc("/workouts/{id}", s.handleUpdateWorkout).Methods("PUT")
	api.HandleFunc("/workouts/{id}", s.handleDeleteWorkout).Methods("DELETE")

	// Workout Session routes
	api.HandleFunc("/workout-sessions", s.handleListWorkoutSessions).Methods("GET")
	api.HandleFunc("/workout-sessions", s.handleCreateWorkoutSession).Methods("POST")
	api.HandleFunc("/workout-sessions/{id}", s.handleGetWorkoutSession).Methods("GET")
	api.HandleFunc("/workout-sessions/{id}", s.handleUpdateWorkoutSession).Methods("PUT")
	api.HandleFunc("/workout-sessions/{id}", s.handleDeleteWorkoutSession).Methods("DELETE")
	api.HandleFunc("/workout-sessions/{id}/exercises/{exerciseIndex}/start", s.handleStartExercise).Methods("POST")
	api.HandleFunc("/workout-sessions/{id}/exercises/{exerciseIndex}/finish", s.handleFinishExercise).Methods("POST")
	api.HandleFunc("/workout-sessions/{id}/exercises/{exerciseIndex}/sets/{setIndex}", s.handleUpdateSet).Methods("PUT")

	// Metrics routes
	api.HandleFunc("/users/{userId}/metrics", s.handleGetUserMetrics).Methods("GET")
	api.HandleFunc("/users/{userId}/metrics/trends", s.handleGetVolumeTrends).Methods("GET")
	api.HandleFunc("/workout-sessions/{sessionId}/metrics", s.handleGetWorkoutMetrics).Methods("GET")

	// Insights routes
	api.HandleFunc("/users/{userId}/insights", s.handleGenerateInsights).Methods("POST")
	api.HandleFunc("/users/{userId}/insights", s.handleGetRecentInsights).Methods("GET")

	// Suggestions routes
	api.HandleFunc("/users/{userId}/suggestions/workouts", s.handleGenerateWorkoutSuggestions).Methods("POST")
	api.HandleFunc("/users/{userId}/suggestions/stored", s.handleGetStoredSuggestions).Methods("GET")

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:8081", "*"}, // Add your frontend URLs
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Printf("HTTP server starting on port %s", port)
	return http.ListenAndServe(":"+port, handler)
}

func (s *Server) handleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	var req GoogleUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("🔐 Google auth request: %+v", req)

	ctx := context.Background()

	// Create or update user
	user, err := s.userService.CreateOrUpdateGoogleUser(
		ctx,
		req.GoogleID,
		req.Email,
		req.Name,
		req.Picture,
	)
	if err != nil {
		log.Printf("❌ Failed to create/update user: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create/update user: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ User created/updated: ID=%s, Email=%s", user.ID.Hex(), user.Email)

	// Create session
	expiresAt := time.Now().Add(24 * time.Hour) // 24 hours from now
	if req.ExpiresAt != "" {
		if parsedTime, err := time.Parse(time.RFC3339, req.ExpiresAt); err == nil {
			expiresAt = parsedTime
		}
	}

	// Generate access token if not provided (for React Native clients)
	accessToken := req.AccessToken
	if accessToken == "" {
		// Generate a unique session token for React Native clients
		accessToken = fmt.Sprintf("session_%s_%d", user.ID.Hex(), time.Now().UnixNano())
		log.Printf("🔑 Generated session token for React Native client: %s", accessToken)
	} else {
		log.Printf("🔑 Using provided access token: %s", accessToken)
	}

	log.Printf("📅 Session expires at: %v", expiresAt)

	session, err := s.userService.CreateSession(ctx, user.ID, accessToken, req.RefreshToken, expiresAt)
	if err != nil {
		log.Printf("❌ Failed to create session: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create session: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Session created successfully: ID=%s, Token=%s", session.ID.Hex(), session.AccessToken)

	// Return user response
	response := s.userToResponse(user)

	// Set session token in response header
	w.Header().Set("X-Session-Token", session.AccessToken)
	log.Printf("📤 Setting X-Session-Token header: %s", session.AccessToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("✅ Google auth response sent successfully")
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Basic health check - verify database connection
	ctx := context.Background()
	if err := s.db.Client.Ping(ctx, nil); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "unhealthy",
			"error":   "database connection failed",
			"service": "gymlog-backend",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"service":   "gymlog-backend",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

func (s *Server) handleValidateSession(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing authorization header", http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	ctx := context.Background()
	user, err := s.userService.ValidateSession(ctx, token)
	if err != nil {
		http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
		return
	}

	response := s.userToResponse(user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing authorization header", http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	ctx := context.Background()
	if err := s.userService.DeleteSession(ctx, token); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// Helper functions for protobuf conversion
func exerciseToResponse(pb *pb.Exercise) ExerciseResponse {
	return ExerciseResponse{
		ID:               pb.Id,
		Name:             pb.Name,
		Description:      pb.Description,
		Category:         pb.Category,
		Equipment:        pb.Equipment,
		PrimaryMuscles:   pb.PrimaryMuscles,
		SecondaryMuscles: pb.SecondaryMuscles,
		Instructions:     pb.Instructions,
		Video:            pb.Video,
		VariationsOn:     pb.VariationsOn,
		VariationOn:      pb.VariationOn,
		CreatedAt:        pb.CreatedAt.AsTime(),
		UpdatedAt:        pb.UpdatedAt.AsTime(),
	}
}

func workoutToResponse(pb *pb.Workout) WorkoutResponse {
	exercises := make([]WorkoutExerciseResponse, len(pb.Exercises))
	for i, ex := range pb.Exercises {
		sets := make([]WorkoutSetResponse, len(ex.Sets))
		for j, set := range ex.Sets {
			sets[j] = WorkoutSetResponse{
				Reps:            set.Reps,
				Weight:          set.Weight,
				DurationSeconds: set.DurationSeconds,
				Distance:        set.Distance,
				Notes:           set.Notes,
			}
		}

		exercises[i] = WorkoutExerciseResponse{
			ExerciseID:  ex.ExerciseId,
			Sets:        sets,
			Notes:       ex.Notes,
			RestSeconds: ex.RestSeconds,
		}

		if ex.Exercise != nil {
			exerciseResp := exerciseToResponse(ex.Exercise)
			exercises[i].Exercise = &exerciseResp
		}
	}

	var startedAt, finishedAt *time.Time
	if pb.StartedAt != nil {
		t := pb.StartedAt.AsTime()
		startedAt = &t
	}
	if pb.FinishedAt != nil {
		t := pb.FinishedAt.AsTime()
		finishedAt = &t
	}

	return WorkoutResponse{
		ID:              pb.Id,
		UserID:          pb.UserId,
		Name:            pb.Name,
		Description:     pb.Description,
		Exercises:       exercises,
		StartedAt:       startedAt,
		FinishedAt:      finishedAt,
		DurationSeconds: pb.DurationSeconds,
		Notes:           pb.Notes,
		CreatedAt:       pb.CreatedAt.AsTime(),
		UpdatedAt:       pb.UpdatedAt.AsTime(),
	}
}

// Exercise HTTP handlers
func (s *Server) handleListExercises(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Parse query parameters
	pageSize := int32(50) // default
	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		if parsed, err := strconv.ParseInt(ps, 10, 32); err == nil {
			pageSize = int32(parsed)
		}
	}

	pageToken := r.URL.Query().Get("pageToken")

	// Parse filter parameters
	categories := r.URL.Query()["categories"]
	equipment := r.URL.Query()["equipment"]
	primaryMuscles := r.URL.Query()["primaryMuscles"]
	secondaryMuscles := r.URL.Query()["secondaryMuscles"]
	search := r.URL.Query().Get("search")

	req := &pb.ListExercisesRequest{
		PageSize:         pageSize,
		PageToken:        pageToken,
		Categories:       categories,
		Equipment:        equipment,
		PrimaryMuscles:   primaryMuscles,
		SecondaryMuscles: secondaryMuscles,
		Search:           search,
	}

	resp, err := s.exerciseService.ListExercises(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list exercises: %v", err), http.StatusInternalServerError)
		return
	}

	exercises := make([]ExerciseResponse, len(resp.Exercises))
	for i, ex := range resp.Exercises {
		exercises[i] = exerciseToResponse(ex)
	}

	response := ListExercisesResponse{
		Exercises:     exercises,
		NextPageToken: resp.NextPageToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleCreateExercise(w http.ResponseWriter, r *http.Request) {
	var req CreateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	pbReq := &pb.CreateExerciseRequest{
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
	}

	exercise, err := s.exerciseService.CreateExercise(ctx, pbReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create exercise: %v", err), http.StatusInternalServerError)
		return
	}

	response := exerciseToResponse(exercise)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleGetExercise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()
	req := &pb.GetExerciseRequest{Id: id}

	exercise, err := s.exerciseService.GetExercise(ctx, req)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	response := exerciseToResponse(exercise)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleUpdateExercise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	pbReq := &pb.UpdateExerciseRequest{
		Id:               id,
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
	}

	exercise, err := s.exerciseService.UpdateExercise(ctx, pbReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update exercise: %v", err), http.StatusInternalServerError)
		return
	}

	response := exerciseToResponse(exercise)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleDeleteExercise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()
	req := &pb.DeleteExerciseRequest{Id: id}

	_, err := s.exerciseService.DeleteExercise(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete exercise: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetQuickAddExercises(w http.ResponseWriter, r *http.Request) {
	log.Printf("🏋️ handleGetQuickAddExercises: Request received - Method: %s, URL: %s", r.Method, r.URL.String())

	ctx := context.Background()

	// Parse query parameters
	userID := r.URL.Query().Get("userId")
	limit := int32(5) // default
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 32); err == nil {
			limit = int32(parsed)
		}
	}

	log.Printf("🏋️ handleGetQuickAddExercises: Parsed parameters - userID: %s, limit: %d", userID, limit)

	req := &pb.GetQuickAddExercisesRequest{
		UserId: userID,
		Limit:  limit,
	}

	log.Printf("🏋️ handleGetQuickAddExercises: Calling gRPC service")
	resp, err := s.exerciseService.GetQuickAddExercises(ctx, req)
	if err != nil {
		log.Printf("❌ handleGetQuickAddExercises: gRPC service error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get quick add exercises: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("🏋️ handleGetQuickAddExercises: gRPC service returned %d exercises", len(resp.Exercises))

	exercises := make([]ExerciseResponse, len(resp.Exercises))
	for i, ex := range resp.Exercises {
		exercises[i] = exerciseToResponse(ex)
	}

	response := QuickAddExercisesResponse{
		Exercises: exercises,
	}

	log.Printf("🏋️ handleGetQuickAddExercises: Sending response with %d exercises", len(exercises))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function for workout session protobuf conversion
func workoutSessionToResponse(pb *pb.WorkoutSession) WorkoutSessionResponse {
	exercises := make([]WorkoutSessionExerciseResponse, len(pb.Exercises))
	for i, ex := range pb.Exercises {
		sets := make([]WorkoutSessionSetResponse, len(ex.Sets))
		for j, set := range ex.Sets {
			var startedAt, finishedAt *time.Time
			if set.StartedAt != nil {
				t := set.StartedAt.AsTime()
				startedAt = &t
			}
			if set.FinishedAt != nil {
				t := set.FinishedAt.AsTime()
				finishedAt = &t
			}

			sets[j] = WorkoutSessionSetResponse{
				TargetReps:      set.TargetReps,
				TargetWeight:    set.TargetWeight,
				ActualReps:      set.ActualReps,
				ActualWeight:    set.ActualWeight,
				DurationSeconds: set.DurationSeconds,
				Distance:        set.Distance,
				Notes:           set.Notes,
				Completed:       set.Completed,
				StartedAt:       startedAt,
				FinishedAt:      finishedAt,
			}
		}

		var exerciseStartedAt, exerciseFinishedAt *time.Time
		if ex.StartedAt != nil {
			t := ex.StartedAt.AsTime()
			exerciseStartedAt = &t
		}
		if ex.FinishedAt != nil {
			t := ex.FinishedAt.AsTime()
			exerciseFinishedAt = &t
		}

		exercises[i] = WorkoutSessionExerciseResponse{
			ExerciseID:  ex.ExerciseId,
			Sets:        sets,
			Notes:       ex.Notes,
			RestSeconds: ex.RestSeconds,
			Completed:   ex.Completed,
			StartedAt:   exerciseStartedAt,
			FinishedAt:  exerciseFinishedAt,
		}

		if ex.Exercise != nil {
			exerciseResp := exerciseToResponse(ex.Exercise)
			exercises[i].Exercise = &exerciseResp
		}
	}

	var startedAt, finishedAt *time.Time
	if pb.StartedAt != nil {
		t := pb.StartedAt.AsTime()
		startedAt = &t
	}
	if pb.FinishedAt != nil {
		t := pb.FinishedAt.AsTime()
		finishedAt = &t
	}

	return WorkoutSessionResponse{
		ID:              pb.Id,
		UserID:          pb.UserId,
		RoutineID:       pb.RoutineId,
		Name:            pb.Name,
		Description:     pb.Description,
		Exercises:       exercises,
		StartedAt:       startedAt,
		FinishedAt:      finishedAt,
		DurationSeconds: pb.DurationSeconds,
		Notes:           pb.Notes,
		RPERating:       pb.RpeRating,
		IsActive:        pb.IsActive,
		CreatedAt:       pb.CreatedAt.AsTime(),
		UpdatedAt:       pb.UpdatedAt.AsTime(),
	}
}

// Workout Session HTTP handlers
func (s *Server) handleListWorkoutSessions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Parse query parameters
	pageSize := int32(50) // default
	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		if parsed, err := strconv.ParseInt(ps, 10, 32); err == nil {
			pageSize = int32(parsed)
		}
	}

	pageToken := r.URL.Query().Get("pageToken")
	userID := r.URL.Query().Get("userId")
	activeOnly := r.URL.Query().Get("activeOnly") == "true"

	req := &pb.ListWorkoutSessionsRequest{
		UserId:     userID,
		PageSize:   pageSize,
		PageToken:  pageToken,
		ActiveOnly: activeOnly,
	}

	// Parse date filters if provided
	if startDate := r.URL.Query().Get("startDate"); startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			req.StartDate = timestamppb.New(t)
		}
	}
	if endDate := r.URL.Query().Get("endDate"); endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			req.EndDate = timestamppb.New(t)
		}
	}

	resp, err := s.workoutSessionService.ListWorkoutSessions(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list workout sessions: %v", err), http.StatusInternalServerError)
		return
	}

	sessions := make([]WorkoutSessionResponse, len(resp.Sessions))
	for i, session := range resp.Sessions {
		sessions[i] = workoutSessionToResponse(session)
	}

	response := ListWorkoutSessionsResponse{
		Sessions:      sessions,
		NextPageToken: resp.NextPageToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleCreateWorkoutSession(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Failed to decode workout session request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("🏃 Create workout session request received: %+v", req)

	ctx := context.Background()

	pbReq := &pb.CreateWorkoutSessionRequest{
		UserId:      req.UserID,
		RoutineId:   req.RoutineID,
		Name:        req.Name,
		Description: req.Description,
		Notes:       req.Notes,
	}

	session, err := s.workoutSessionService.CreateWorkoutSession(ctx, pbReq)
	if err != nil {
		log.Printf("❌ Failed to create workout session: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create workout session: %v", err), http.StatusInternalServerError)
		return
	}

	response := workoutSessionToResponse(session)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	log.Printf("✅ Workout session created successfully: ID=%s", session.Id)
}

func (s *Server) handleGetWorkoutSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()
	req := &pb.GetWorkoutSessionRequest{Id: id}

	session, err := s.workoutSessionService.GetWorkoutSession(ctx, req)
	if err != nil {
		http.Error(w, "Workout session not found", http.StatusNotFound)
		return
	}

	response := workoutSessionToResponse(session)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleUpdateWorkoutSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateWorkoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	pbReq := &pb.UpdateWorkoutSessionRequest{
		Id:              id,
		Name:            req.Name,
		Description:     req.Description,
		DurationSeconds: req.DurationSeconds,
		Notes:           req.Notes,
		RpeRating:       req.RPERating,
	}

	if req.FinishedAt != nil {
		pbReq.FinishedAt = timestamppb.New(*req.FinishedAt)
	}

	if req.IsActive != nil {
		pbReq.IsActive = *req.IsActive
	}

	// Convert exercises if provided
	if len(req.Exercises) > 0 {
		exercises := make([]*pb.WorkoutSessionExercise, len(req.Exercises))
		for i, ex := range req.Exercises {
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
					DurationSeconds: set.DurationSeconds,
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
				ExerciseId:  ex.ExerciseID,
				Sets:        sets,
				Notes:       ex.Notes,
				RestSeconds: ex.RestSeconds,
				Completed:   ex.Completed,
				StartedAt:   exerciseStartedAt,
				FinishedAt:  exerciseFinishedAt,
			}
		}
		pbReq.Exercises = exercises
	}

	session, err := s.workoutSessionService.UpdateWorkoutSession(ctx, pbReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update workout session: %v", err), http.StatusInternalServerError)
		return
	}

	response := workoutSessionToResponse(session)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleDeleteWorkoutSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()
	req := &pb.DeleteWorkoutSessionRequest{Id: id}

	_, err := s.workoutSessionService.DeleteWorkoutSession(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete workout session: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleStartExercise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]
	exerciseIndexStr := vars["exerciseIndex"]

	exerciseIndex, err := strconv.ParseInt(exerciseIndexStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid exercise index", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	req := &pb.StartExerciseRequest{
		SessionId:     sessionID,
		ExerciseIndex: int32(exerciseIndex),
	}

	session, err := s.workoutSessionService.StartExercise(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start exercise: %v", err), http.StatusInternalServerError)
		return
	}

	response := workoutSessionToResponse(session)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleFinishExercise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]
	exerciseIndexStr := vars["exerciseIndex"]

	exerciseIndex, err := strconv.ParseInt(exerciseIndexStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid exercise index", http.StatusBadRequest)
		return
	}

	// Parse request body for notes
	var reqBody struct {
		Notes string `json:"notes,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		// Notes are optional, so ignore decode errors
	}

	ctx := context.Background()
	req := &pb.FinishExerciseRequest{
		SessionId:     sessionID,
		ExerciseIndex: int32(exerciseIndex),
		Notes:         reqBody.Notes,
	}

	session, err := s.workoutSessionService.FinishExercise(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to finish exercise: %v", err), http.StatusInternalServerError)
		return
	}

	response := workoutSessionToResponse(session)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleUpdateSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]
	exerciseIndexStr := vars["exerciseIndex"]
	setIndexStr := vars["setIndex"]

	exerciseIndex, err := strconv.ParseInt(exerciseIndexStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid exercise index", http.StatusBadRequest)
		return
	}

	setIndex, err := strconv.ParseInt(setIndexStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid set index", http.StatusBadRequest)
		return
	}

	var req UpdateSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	pbReq := &pb.UpdateSetRequest{
		SessionId:       sessionID,
		ExerciseIndex:   int32(exerciseIndex),
		SetIndex:        int32(setIndex),
		ActualReps:      req.ActualReps,
		ActualWeight:    req.ActualWeight,
		DurationSeconds: req.DurationSeconds,
		Distance:        req.Distance,
		Notes:           req.Notes,
		Completed:       req.Completed,
	}

	session, err := s.workoutSessionService.UpdateSet(ctx, pbReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update set: %v", err), http.StatusInternalServerError)
		return
	}

	response := workoutSessionToResponse(session)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Workout HTTP handlers
func (s *Server) handleListWorkouts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Parse query parameters
	pageSize := int32(50) // default
	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		if parsed, err := strconv.ParseInt(ps, 10, 32); err == nil {
			pageSize = int32(parsed)
		}
	}

	pageToken := r.URL.Query().Get("pageToken")
	userID := r.URL.Query().Get("userId")

	req := &pb.ListWorkoutsRequest{
		UserId:    userID,
		PageSize:  pageSize,
		PageToken: pageToken,
	}

	// Parse date filters if provided
	if startDate := r.URL.Query().Get("startDate"); startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			req.StartDate = timestamppb.New(t)
		}
	}
	if endDate := r.URL.Query().Get("endDate"); endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			req.EndDate = timestamppb.New(t)
		}
	}

	resp, err := s.workoutService.ListWorkouts(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list workouts: %v", err), http.StatusInternalServerError)
		return
	}

	workouts := make([]WorkoutResponse, len(resp.Workouts))
	for i, workout := range resp.Workouts {
		workouts[i] = workoutToResponse(workout)
	}

	response := ListWorkoutsResponse{
		Workouts:      workouts,
		NextPageToken: resp.NextPageToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Failed to decode workout request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("🏋️ Create workout request received: %+v", req)
	log.Printf("👤 User ID from request: '%s' (length: %d)", req.UserID, len(req.UserID))

	ctx := context.Background()

	// Convert exercises
	exercises := make([]*pb.WorkoutExercise, len(req.Exercises))
	for i, ex := range req.Exercises {
		sets := make([]*pb.WorkoutSet, len(ex.Sets))
		for j, set := range ex.Sets {
			sets[j] = &pb.WorkoutSet{
				Reps:            set.Reps,
				Weight:          set.Weight,
				DurationSeconds: set.DurationSeconds,
				Distance:        set.Distance,
				Notes:           set.Notes,
			}
		}

		exercises[i] = &pb.WorkoutExercise{
			ExerciseId:  ex.ExerciseID,
			Sets:        sets,
			Notes:       ex.Notes,
			RestSeconds: ex.RestSeconds,
		}
	}

	pbReq := &pb.CreateWorkoutRequest{
		UserId:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
		Exercises:   exercises,
		Notes:       req.Notes,
	}

	if req.StartedAt != nil {
		pbReq.StartedAt = timestamppb.New(*req.StartedAt)
	}

	log.Printf("📋 Calling gRPC CreateWorkout with UserId: '%s'", pbReq.UserId)

	workout, err := s.workoutService.CreateWorkout(ctx, pbReq)
	if err != nil {
		log.Printf("❌ gRPC CreateWorkout failed: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create workout: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Workout created successfully: ID=%s", workout.Id)

	response := workoutToResponse(workout)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleGetWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()
	req := &pb.GetWorkoutRequest{Id: id}

	workout, err := s.workoutService.GetWorkout(ctx, req)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	response := workoutToResponse(workout)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleStartWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()
	req := &pb.GetWorkoutRequest{Id: id}

	workout, err := s.workoutService.GetWorkout(ctx, req)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	// Populate exercise details for each exercise in the workout
	for i, workoutExercise := range workout.Exercises {
		exerciseReq := &pb.GetExerciseRequest{Id: workoutExercise.ExerciseId}
		exercise, err := s.exerciseService.GetExercise(ctx, exerciseReq)
		if err != nil {
			log.Printf("Warning: Could not fetch exercise details for ID %s: %v", workoutExercise.ExerciseId, err)
			// Continue without exercise details if not found
			continue
		}
		workout.Exercises[i].Exercise = exercise
	}

	response := workoutToResponse(workout)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	pbReq := &pb.UpdateWorkoutRequest{
		Id:              id,
		Name:            req.Name,
		Description:     req.Description,
		DurationSeconds: req.DurationSeconds,
		Notes:           req.Notes,
	}

	// Convert exercises if provided
	if len(req.Exercises) > 0 {
		exercises := make([]*pb.WorkoutExercise, len(req.Exercises))
		for i, ex := range req.Exercises {
			sets := make([]*pb.WorkoutSet, len(ex.Sets))
			for j, set := range ex.Sets {
				sets[j] = &pb.WorkoutSet{
					Reps:            set.Reps,
					Weight:          set.Weight,
					DurationSeconds: set.DurationSeconds,
					Distance:        set.Distance,
					Notes:           set.Notes,
				}
			}

			exercises[i] = &pb.WorkoutExercise{
				ExerciseId:  ex.ExerciseID,
				Sets:        sets,
				Notes:       ex.Notes,
				RestSeconds: ex.RestSeconds,
			}
		}
		pbReq.Exercises = exercises
	}

	if req.FinishedAt != nil {
		pbReq.FinishedAt = timestamppb.New(*req.FinishedAt)
	}

	workout, err := s.workoutService.UpdateWorkout(ctx, pbReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update workout: %v", err), http.StatusInternalServerError)
		return
	}

	response := workoutToResponse(workout)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()
	req := &pb.DeleteWorkoutRequest{Id: id}

	_, err := s.workoutService.DeleteWorkout(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete workout: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Metrics HTTP handlers

func (s *Server) handleGetUserMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "all"
	}

	ctx := context.Background()
	req := &pb.GetUserMetricsRequest{
		UserId: userID,
		Period: period,
	}

	metrics, err := s.metricsService.GetUserMetrics(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user metrics: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (s *Server) handleGetVolumeTrends(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "weekly"
	}

	ctx := context.Background()
	req := &pb.GetVolumeTrendsRequest{
		UserId: userID,
		Period: period,
	}

	trends, err := s.metricsService.GetVolumeTrends(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get volume trends: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trends)
}

func (s *Server) handleGetWorkoutMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	ctx := context.Background()
	req := &pb.GetWorkoutMetricsRequest{
		SessionId: sessionID,
	}

	metrics, err := s.metricsService.GetWorkoutMetrics(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get workout metrics: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// Insights HTTP handlers

func (s *Server) handleGenerateInsights(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	ctx := context.Background()

	// Convert userID string to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	insights, err := s.insightsService.GenerateWorkoutInsights(ctx, userObjID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate insights: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	insightResponses := make([]WorkoutInsightResponse, len(insights))
	for i, insight := range insights {
		insightResponses[i] = WorkoutInsightResponse{
			ID:        insight.ID.Hex(),
			UserID:    insight.UserID.Hex(),
			Type:      string(insight.Type),
			Insight:   insight.Insight,
			BasedOn:   insight.BasedOn,
			Priority:  int32(insight.Priority),
			CreatedAt: insight.CreatedAt,
		}
	}

	response := GenerateInsightsResponse{
		Insights: insightResponses,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleGetRecentInsights(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	// Parse limit parameter
	limit := 5 // default
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 32); err == nil {
			limit = int(parsed)
		}
	}

	ctx := context.Background()

	// Convert userID string to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	insights, err := s.insightsService.GetRecentInsights(ctx, userObjID, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get recent insights: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	insightResponses := make([]WorkoutInsightResponse, len(insights))
	for i, insight := range insights {
		insightResponses[i] = WorkoutInsightResponse{
			ID:        insight.ID.Hex(),
			UserID:    insight.UserID.Hex(),
			Type:      string(insight.Type),
			Insight:   insight.Insight,
			BasedOn:   insight.BasedOn,
			Priority:  int32(insight.Priority),
			CreatedAt: insight.CreatedAt,
		}
	}

	response := GetRecentInsightsResponse{
		Insights: insightResponses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to convert user model to response
func (s *Server) userToResponse(user *models.User) UserResponse {
	// Parse name to get first and last name
	nameParts := strings.Split(user.Name, " ")
	firstName := ""
	lastName := ""
	if len(nameParts) > 0 {
		firstName = nameParts[0]
		if len(nameParts) > 1 {
			lastName = strings.Join(nameParts[1:], " ")
		}
	}

	return UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		FirstName: firstName,
		LastName:  lastName,
		Height:    user.Height,
		Weight:    user.Weight,
		Age:       user.Age,
		Goal:      user.Goal,
		GoogleID:  user.GoogleID,
		Picture:   user.Picture,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// User update handler
func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	pbReq := &pb.UpdateUserRequest{
		Id:        userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Height:    req.Height,
		Weight:    req.Weight,
		Age:       req.Age,
		Goal:      req.Goal,
	}

	user, err := s.userService.UpdateUser(ctx, pbReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert protobuf user to models.User for response
	modelUser := &models.User{
		ID:        primitive.ObjectID{}, // Will be set from user.Id
		Email:     user.Email,
		Name:      fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Height:    user.Height,
		Weight:    user.Weight,
		Age:       user.Age,
		Goal:      user.Goal,
		CreatedAt: user.CreatedAt.AsTime(),
		UpdatedAt: user.UpdatedAt.AsTime(),
	}

	// Convert user ID
	if userObjID, err := primitive.ObjectIDFromHex(user.Id); err == nil {
		modelUser.ID = userObjID
	}

	response := s.userToResponse(modelUser)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Suggestions HTTP handlers
func (s *Server) handleGenerateWorkoutSuggestions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var req GenerateWorkoutSuggestionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Defaults will be used
	}

	ctx := context.Background()
	pbReq := &pb.GenerateWorkoutSuggestionsRequest{
		UserId:         userID,
		DaysToAnalyze:  req.DaysToAnalyze,
		MaxSuggestions: req.MaxSuggestions,
	}

	log.Printf("🤖 Generating workout suggestions for user %s", userID)

	resp, err := s.suggestionsService.GenerateWorkoutSuggestions(ctx, pbReq)
	if err != nil {
		log.Printf("❌ Failed to generate workout suggestions: %v", err)
		http.Error(w, fmt.Sprintf("Failed to generate workout suggestions: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert protobuf response to HTTP response
	suggestions := make([]SuggestedWorkoutResponse, len(resp.Suggestions))
	for i, suggestion := range resp.Suggestions {
		// Convert exercises
		exercises := make([]WorkoutExerciseResponse, len(suggestion.Exercises))
		for j, ex := range suggestion.Exercises {
			sets := make([]WorkoutSetResponse, len(ex.Sets))
			for k, set := range ex.Sets {
				sets[k] = WorkoutSetResponse{
					Reps:            set.Reps,
					Weight:          set.Weight,
					DurationSeconds: set.DurationSeconds,
					Distance:        set.Distance,
					Notes:           set.Notes,
				}
			}

			exercises[j] = WorkoutExerciseResponse{
				ExerciseID:  ex.ExerciseId,
				Sets:        sets,
				Notes:       ex.Notes,
				RestSeconds: ex.RestSeconds,
			}

			// Populate exercise details
			if ex.Exercise != nil {
				exerciseResp := exerciseToResponse(ex.Exercise)
				exercises[j].Exercise = &exerciseResp
			}
		}

		// Convert changes
		changes := make([]WorkoutChangeResponse, len(suggestion.Changes))
		for j, change := range suggestion.Changes {
			changes[j] = WorkoutChangeResponse{
				Type:         change.Type,
				ExerciseID:   change.ExerciseId,
				ExerciseName: change.ExerciseName,
				OldValue:     change.OldValue,
				NewValue:     change.NewValue,
				Reason:       change.Reason,
			}
		}

		suggestions[i] = SuggestedWorkoutResponse{
			OriginalWorkoutID: suggestion.OriginalWorkoutId,
			Name:              suggestion.Name,
			Description:       suggestion.Description,
			Exercises:         exercises,
			Changes:           changes,
			OverallReasoning:  suggestion.OverallReasoning,
			Priority:          suggestion.Priority,
		}
	}

	response := GenerateWorkoutSuggestionsResponse{
		Suggestions:     suggestions,
		AnalysisSummary: resp.AnalysisSummary,
	}

	log.Printf("✅ Successfully generated %d workout suggestions for user %s", len(suggestions), userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleGetStoredSuggestions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	// Parse query parameters
	pageSize := int32(10) // default
	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		if parsed, err := strconv.ParseInt(ps, 10, 32); err == nil && parsed > 0 {
			pageSize = int32(parsed)
		}
	}

	pageToken := r.URL.Query().Get("pageToken")

	ctx := context.Background()
	pbReq := &pb.GetStoredSuggestionsRequest{
		UserId:    userID,
		PageSize:  pageSize,
		PageToken: pageToken,
	}

	log.Printf("📊 Getting stored workout suggestions for user %s", userID)

	resp, err := s.suggestionsService.GetStoredSuggestions(ctx, pbReq)
	if err != nil {
		log.Printf("❌ Failed to get stored workout suggestions: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get stored workout suggestions: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert protobuf response to HTTP response
	storedSuggestions := make([]StoredWorkoutSuggestionResponse, len(resp.StoredSuggestions))
	for i, stored := range resp.StoredSuggestions {
		// Convert suggestions
		suggestions := make([]SuggestedWorkoutResponse, len(stored.Suggestions))
		for j, suggestion := range stored.Suggestions {
			// Convert exercises
			exercises := make([]WorkoutExerciseResponse, len(suggestion.Exercises))
			for k, ex := range suggestion.Exercises {
				sets := make([]WorkoutSetResponse, len(ex.Sets))
				for l, set := range ex.Sets {
					sets[l] = WorkoutSetResponse{
						Reps:            set.Reps,
						Weight:          set.Weight,
						DurationSeconds: set.DurationSeconds,
						Distance:        set.Distance,
						Notes:           set.Notes,
					}
				}

				exercises[k] = WorkoutExerciseResponse{
					ExerciseID:  ex.ExerciseId,
					Sets:        sets,
					Notes:       ex.Notes,
					RestSeconds: ex.RestSeconds,
				}

				// Populate exercise details if available
				if ex.Exercise != nil {
					exerciseResp := exerciseToResponse(ex.Exercise)
					exercises[k].Exercise = &exerciseResp
				}
			}

			// Convert changes
			changes := make([]WorkoutChangeResponse, len(suggestion.Changes))
			for k, change := range suggestion.Changes {
				changes[k] = WorkoutChangeResponse{
					Type:         change.Type,
					ExerciseID:   change.ExerciseId,
					ExerciseName: change.ExerciseName,
					OldValue:     change.OldValue,
					NewValue:     change.NewValue,
					Reason:       change.Reason,
				}
			}

			suggestions[j] = SuggestedWorkoutResponse{
				OriginalWorkoutID: suggestion.OriginalWorkoutId,
				Name:              suggestion.Name,
				Description:       suggestion.Description,
				Exercises:         exercises,
				Changes:           changes,
				OverallReasoning:  suggestion.OverallReasoning,
				Priority:          suggestion.Priority,
			}
		}

		storedSuggestions[i] = StoredWorkoutSuggestionResponse{
			ID:              stored.Id,
			UserID:          stored.UserId,
			Suggestions:     suggestions,
			AnalysisSummary: stored.AnalysisSummary,
			DaysAnalyzed:    stored.DaysAnalyzed,
			CreatedAt:       stored.CreatedAt.AsTime(),
			UpdatedAt:       stored.UpdatedAt.AsTime(),
		}
	}

	response := GetStoredSuggestionsResponse{
		StoredSuggestions: storedSuggestions,
		NextPageToken:     resp.NextPageToken,
	}

	log.Printf("✅ Successfully retrieved %d stored workout suggestions for user %s", len(storedSuggestions), userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
