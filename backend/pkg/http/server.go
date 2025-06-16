package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"gymlog-backend/internal/services"
	"gymlog-backend/pkg/database"
)

type Server struct {
	userService *services.UserService
	db          *database.MongoDB
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

func NewServer(userService *services.UserService, db *database.MongoDB) *Server {
	return &Server{
		userService: userService,
		db:          db,
	}
}

func (s *Server) Start(port string) error {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/google", s.handleGoogleAuth).Methods("POST")
	api.HandleFunc("/auth/validate", s.handleValidateSession).Methods("GET")
	api.HandleFunc("/auth/logout", s.handleLogout).Methods("POST")

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"}, // Add your frontend URLs
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
		http.Error(w, fmt.Sprintf("Failed to create/update user: %v", err), http.StatusInternalServerError)
		return
	}

	// Create session
	expiresAt := time.Now().Add(24 * time.Hour) // 24 hours from now
	if req.ExpiresAt != "" {
		if parsedTime, err := time.Parse(time.RFC3339, req.ExpiresAt); err == nil {
			expiresAt = parsedTime
		}
	}

	session, err := s.userService.CreateSession(ctx, user.ID, req.AccessToken, req.RefreshToken, expiresAt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create session: %v", err), http.StatusInternalServerError)
		return
	}

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
