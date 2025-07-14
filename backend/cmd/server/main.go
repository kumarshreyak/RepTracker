package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"gymlog-backend/internal/services"
	"gymlog-backend/pkg/database"
	httpserver "gymlog-backend/pkg/http"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "gymlog"
	}

	// Connect to MongoDB
	db, err := database.NewMongoDB(mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Close()

	// Create service instances
	userService := services.NewUserService(db)
	exerciseService := services.NewExerciseService(db)
	workoutService := services.NewWorkoutService(db)
	metricsService := services.NewMetricsService(db, exerciseService)
	workoutSessionService, err := services.NewWorkoutSessionService(db, workoutService, exerciseService, metricsService, userService)
	if err != nil {
		log.Fatalf("Failed to create workout session service: %v", err)
	}

	// Create insights service
	insightsService, err := services.NewInsightsService(db.Database)
	if err != nil {
		log.Fatalf("Failed to create insights service: %v", err)
	}

	// Create suggestions service
	suggestionsService, err := services.NewSuggestionsService(db.Database, workoutService, userService, exerciseService, insightsService)
	if err != nil {
		log.Fatalf("Failed to create suggestions service: %v", err)
	}

	// Get port from environment (Cloud Run uses PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server with all services
	httpServer := httpserver.NewServer(userService, exerciseService, workoutService, workoutSessionService, metricsService, insightsService, suggestionsService, db)
	log.Printf("Starting HTTP server on port %s...", port)
	if err := httpServer.Start(port); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
