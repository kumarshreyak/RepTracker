package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gymlog-backend/internal/services"
	"gymlog-backend/pkg/database"
	httpserver "gymlog-backend/pkg/http"
	pb "gymlog-backend/proto/gymlog/v1"
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

	// Start HTTP server in a goroutine
	go func() {
		httpServer := httpserver.NewServer(userService, db)
		log.Println("Starting HTTP server on port 8080...")
		if err := httpServer.Start("8080"); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Start gRPC server
	// Create a listener on TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Register services with the gRPC server
	pb.RegisterUserServiceServer(s, userService)
	pb.RegisterExerciseServiceServer(s, exerciseService)
	pb.RegisterWorkoutServiceServer(s, workoutService)

	// Register reflection service on gRPC server for debugging
	reflection.Register(s)

	log.Println("gRPC server starting on port 50051...")

	// Start the gRPC server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
