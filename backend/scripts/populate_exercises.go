package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"

	"gymlog-backend/pkg/database"
	"gymlog-backend/pkg/models"
)

// JSONExercise represents the structure of exercises in the JSON file
type JSONExercise struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	MuscleGroup  string   `json:"muscle_group"`
	Equipment    string   `json:"equipment"`
	Difficulty   string   `json:"difficulty"`
	Instructions []string `json:"instructions"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

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

	// Get exercises collection
	exerciseColl := db.GetCollection("exercises")

	// Read the JSON file
	jsonFile := filepath.Join("scripts", "strengthlog_exercises.json")
	jsonData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	// Parse JSON data
	var jsonExercises []JSONExercise
	if err := json.Unmarshal(jsonData, &jsonExercises); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	log.Printf("Found %d exercises in JSON file", len(jsonExercises))

	// Clear existing exercises
	ctx := context.Background()
	result, err := exerciseColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to clear exercises collection: %v", err)
	}
	log.Printf("Deleted %d existing exercises", result.DeletedCount)

	// Convert JSON exercises to MongoDB models
	var exercises []interface{}
	now := time.Now()

	for _, jsonEx := range jsonExercises {
		// Parse timestamps if they exist, otherwise use current time
		var createdAt, updatedAt time.Time
		if jsonEx.CreatedAt != "" {
			if parsed, err := time.Parse(time.RFC3339, jsonEx.CreatedAt); err == nil {
				createdAt = parsed
			} else {
				createdAt = now
			}
		} else {
			createdAt = now
		}

		if jsonEx.UpdatedAt != "" {
			if parsed, err := time.Parse(time.RFC3339, jsonEx.UpdatedAt); err == nil {
				updatedAt = parsed
			} else {
				updatedAt = now
			}
		} else {
			updatedAt = now
		}

		exercise := models.Exercise{
			Name:         jsonEx.Name,
			Description:  jsonEx.Description,
			MuscleGroup:  jsonEx.MuscleGroup,
			Equipment:    jsonEx.Equipment,
			Difficulty:   jsonEx.Difficulty,
			Instructions: jsonEx.Instructions,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		}

		exercises = append(exercises, exercise)
	}

	// Insert exercises in batches
	batchSize := 100
	for i := 0; i < len(exercises); i += batchSize {
		end := i + batchSize
		if end > len(exercises) {
			end = len(exercises)
		}

		batch := exercises[i:end]
		result, err := exerciseColl.InsertMany(ctx, batch)
		if err != nil {
			log.Fatalf("Failed to insert exercises batch %d-%d: %v", i, end-1, err)
		}

		log.Printf("Inserted batch %d-%d (%d exercises)", i, end-1, len(result.InsertedIDs))
	}

	log.Printf("Successfully populated exercises collection with %d exercises", len(exercises))
}
