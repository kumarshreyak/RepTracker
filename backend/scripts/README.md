# Scripts

This directory contains utility scripts for the GymLog backend.

## Exercise Population Script

### Purpose
This script clears the existing exercises collection in MongoDB and populates it with exercises from the `strengthlog_exercises.json` file.

### Files
- `populate_exercises.go` - Go script that connects to MongoDB and performs the data operations
- `populate_exercises.sh` - Shell script wrapper for easy execution
- `strengthlog_exercises.json` - Source data containing 325 exercises from StrengthLog

### Usage

#### Prerequisites
1. Make sure MongoDB is running
2. Ensure you have the proper environment variables set (MONGODB_URI, DB_NAME) or the script will use defaults:
   - MONGODB_URI: `mongodb://localhost:27017`
   - DB_NAME: `gymlog`

#### Running the Script

**Option 1: Using the shell script (recommended)**
```bash
cd backend
./scripts/populate_exercises.sh
```

**Option 2: Running the Go script directly**
```bash
cd backend
go run scripts/populate_exercises.go
```

### What the Script Does

1. **Connects to MongoDB** using the configuration from environment variables or defaults
2. **Clears existing exercises** - Deletes all documents from the `exercises` collection
3. **Reads the JSON file** - Parses `strengthlog_exercises.json` containing 325 exercises
4. **Transforms the data** - Converts JSON format to MongoDB document format
5. **Inserts in batches** - Adds exercises to the database in batches of 100 for efficiency
6. **Reports progress** - Shows how many exercises were deleted and inserted

### Exercise Data Structure

Each exercise contains:
- **Name** - Exercise name (e.g., "Push-ups", "Bench Press")
- **Description** - Detailed description of the exercise
- **MuscleGroup** - Primary muscle group targeted (e.g., "Chest", "Back", "Leg")
- **Equipment** - Required equipment (e.g., "Bodyweight", "Barbell", "Dumbbell")
- **Difficulty** - Difficulty level (e.g., "Beginner", "Intermediate", "Advanced")
- **Instructions** - Step-by-step instructions as an array of strings
- **CreatedAt/UpdatedAt** - Timestamps for when the exercise was added/modified

### Sample Output
```
Starting exercise population script...
Building population script...
Running population script...
2025-01-17 10:30:15 Successfully connected to MongoDB
Found 325 exercises in JSON file
Deleted 10 existing exercises
Inserted batch 0-99 (100 exercises)
Inserted batch 100-199 (100 exercises)
Inserted batch 200-299 (100 exercises)
Inserted batch 300-324 (25 exercises)
Successfully populated exercises collection with 325 exercises
Exercise population completed!
```

### Error Handling

The script will fail and show an error message if:
- MongoDB is not accessible
- The JSON file is missing or corrupted
- There are issues with database permissions
- The Go build fails

### Database Impact

⚠️ **Warning**: This script will delete ALL existing exercises in your database before adding the new ones. Make sure you have backups if needed.

The script is designed to be idempotent - you can run it multiple times safely. 