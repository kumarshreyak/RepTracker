package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"gymlog-backend/internal/services"
	"gymlog-backend/pkg/database"
	pb "gymlog-backend/proto/gymlog/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	userService     *services.UserService
	exerciseService *services.ExerciseService
	workoutService  *services.WorkoutService
	db              *database.MongoDB
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
	GoogleID  string    `json:"googleId"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Exercise HTTP types
type ExerciseResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	MuscleGroup  string    `json:"muscleGroup"`
	Equipment    string    `json:"equipment"`
	Difficulty   string    `json:"difficulty"`
	Instructions []string  `json:"instructions"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type CreateExerciseRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	MuscleGroup  string   `json:"muscleGroup"`
	Equipment    string   `json:"equipment"`
	Difficulty   string   `json:"difficulty"`
	Instructions []string `json:"instructions"`
}

type UpdateExerciseRequest struct {
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	MuscleGroup  string   `json:"muscleGroup,omitempty"`
	Equipment    string   `json:"equipment,omitempty"`
	Difficulty   string   `json:"difficulty,omitempty"`
	Instructions []string `json:"instructions,omitempty"`
}

type ListExercisesResponse struct {
	Exercises     []ExerciseResponse `json:"exercises"`
	NextPageToken string             `json:"nextPageToken,omitempty"`
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

func NewServer(userService *services.UserService, exerciseService *services.ExerciseService, workoutService *services.WorkoutService, db *database.MongoDB) *Server {
	return &Server{
		userService:     userService,
		exerciseService: exerciseService,
		workoutService:  workoutService,
		db:              db,
	}
}

func (s *Server) Start(port string) error {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/google", s.handleGoogleAuth).Methods("POST")
	api.HandleFunc("/auth/validate", s.handleValidateSession).Methods("GET")
	api.HandleFunc("/auth/logout", s.handleLogout).Methods("POST")

	// Exercise routes
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
	response := UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		GoogleID:  user.GoogleID,
		Picture:   user.Picture,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Set session token in response header
	w.Header().Set("X-Session-Token", session.AccessToken)
	log.Printf("📤 Setting X-Session-Token header: %s", session.AccessToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("✅ Google auth response sent successfully")
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

	response := UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		GoogleID:  user.GoogleID,
		Picture:   user.Picture,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

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
		ID:           pb.Id,
		Name:         pb.Name,
		Description:  pb.Description,
		MuscleGroup:  pb.MuscleGroup,
		Equipment:    pb.Equipment,
		Difficulty:   pb.Difficulty,
		Instructions: pb.Instructions,
		CreatedAt:    pb.CreatedAt.AsTime(),
		UpdatedAt:    pb.UpdatedAt.AsTime(),
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
	muscleGroup := r.URL.Query().Get("muscleGroup")
	equipment := r.URL.Query().Get("equipment")

	req := &pb.ListExercisesRequest{
		PageSize:    pageSize,
		PageToken:   pageToken,
		MuscleGroup: muscleGroup,
		Equipment:   equipment,
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
		Name:         req.Name,
		Description:  req.Description,
		MuscleGroup:  req.MuscleGroup,
		Equipment:    req.Equipment,
		Difficulty:   req.Difficulty,
		Instructions: req.Instructions,
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
		Id:           id,
		Name:         req.Name,
		Description:  req.Description,
		MuscleGroup:  req.MuscleGroup,
		Equipment:    req.Equipment,
		Difficulty:   req.Difficulty,
		Instructions: req.Instructions,
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
