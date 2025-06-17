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
	exerciseService := services.NewExerciseService()
	workoutService := services.NewWorkoutService()
	workoutSessionService := services.NewWorkoutSessionService(workoutService, exerciseService)

	// Start HTTP server with all services
	httpServer := httpserver.NewServer(userService, exerciseService, workoutService, workoutSessionService, db)
	log.Println("Starting HTTP server on port 8080...")
	if err := httpServer.Start("8080"); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
