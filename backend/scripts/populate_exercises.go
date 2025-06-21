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

// JSONFile represents the structure of the exercises.json file
type JSONFile struct {
	Categories []string       `json:"categories"`
	Equipment  []string       `json:"equipment"`
	Exercises  []JSONExercise `json:"exercises"`
}

// JSONExercise represents the structure of exercises in the JSON file
type JSONExercise struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Category         string   `json:"category"`
	Equipment        []string `json:"equipment"`
	PrimaryMuscles   []string `json:"primary_muscles"`
	SecondaryMuscles []string `json:"secondary_muscles"`
	Instructions     []string `json:"instructions"`
	Video            string   `json:"video,omitempty"`
	VariationsOn     []string `json:"variations_on,omitempty"`
	VariationOn      []string `json:"variation_on,omitempty"`
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
	jsonFile := filepath.Join("scripts", "exercises.json")
	jsonData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	// Parse JSON data
	var jsonFileData JSONFile
	if err := json.Unmarshal(jsonData, &jsonFileData); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	log.Printf("Found %d exercises in JSON file", len(jsonFileData.Exercises))

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

	for _, jsonEx := range jsonFileData.Exercises {
		exercise := models.Exercise{
			Name:             jsonEx.Name,
			Description:      jsonEx.Description,
			Category:         jsonEx.Category,
			Equipment:        jsonEx.Equipment,
			PrimaryMuscles:   jsonEx.PrimaryMuscles,
			SecondaryMuscles: jsonEx.SecondaryMuscles,
			Instructions:     jsonEx.Instructions,
			Video:            jsonEx.Video,
			VariationsOn:     jsonEx.VariationsOn,
			VariationOn:      jsonEx.VariationOn,
			CreatedAt:        now,
			UpdatedAt:        now,
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
