package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gymlog-backend/internal/services"
	pb "gymlog-backend/proto/gymlog/v1"
)

func main() {
	// Create a listener on TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Create service instances
	// userService := services.NewUserService()
	exerciseService := services.NewExerciseService()
	workoutService := services.NewWorkoutService()

	// Register services with the gRPC server
	// pb.RegisterUserServiceServer(s, userService)
	pb.RegisterExerciseServiceServer(s, exerciseService)
	pb.RegisterWorkoutServiceServer(s, workoutService)

	// Register reflection service on gRPC server for debugging
	reflection.Register(s)

	log.Println("gRPC server starting on port 50051...")

	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
