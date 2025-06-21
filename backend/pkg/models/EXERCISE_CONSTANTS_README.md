# Exercise Constants and Models Usage Guide

This guide explains how to use the exercise constants, enums, and models in the GymLog service layer.

## Available Constants

### Exercise Categories
```go
import "gymlog-backend/pkg/models"

// Available categories:
models.CategoryStrength              // "strength"
models.CategoryStretching            // "stretching"
models.CategoryPlyometrics           // "plyometrics"
models.CategoryStrongman             // "strongman"
models.CategoryCardio                // "cardio"
models.CategoryOlympicWeightlifting  // "olympic weightlifting"
models.CategoryCrossfit              // "crossfit"
models.CategoryCalisthenics          // "calisthenics"
```

### Equipment Types
```go
models.EquipmentNone         // "none"
models.EquipmentEZCurlBar    // "ez curl bar"
models.EquipmentBarbell      // "barbell"
models.EquipmentDumbbell     // "dumbbell"
models.EquipmentGymMat       // "gym mat"
models.EquipmentExerciseBall // "exercise ball"
models.EquipmentMedicineBall // "medicine ball"
models.EquipmentPullUpBar    // "pull-up bar"
models.EquipmentBench        // "bench"
models.EquipmentInclineBench // "incline bench"
models.EquipmentKettlebell   // "kettlebell"
models.EquipmentMachine      // "machine"
models.EquipmentCable        // "cable"
models.EquipmentBands        // "bands"
models.EquipmentFoamRoll     // "foam roll"
models.EquipmentOther        // "other"
```

### Muscle Types
```go
// Examples:
models.MuscleBiceps      // "biceps"
models.MuscleChest       // "chest"
models.MuscleQuads       // "quads"
models.MuscleGlutes      // "glutes"
models.MuscleAbs         // "abs"
// ... and many more
```

### Muscle Groups
```go
models.MuscleGroupArms       // "arms"
models.MuscleGroupBack       // "back"
models.MuscleGroupCalves     // "calves"
models.MuscleGroupChest      // "chest"
models.MuscleGroupCore       // "core"
models.MuscleGroupLegs       // "legs"
models.MuscleGroupShoulders  // "shoulders"
```

## Usage Examples

### 1. Filtering Exercises by Category
```go
func GetStrengthExercises(ctx context.Context) ([]*models.Exercise, error) {
    filter := bson.M{
        "category": models.CategoryStrength,
    }
    return exerciseCollection.Find(ctx, filter)
}
```

### 2. Filtering by Equipment
```go
func GetBodyweightExercises(ctx context.Context) ([]*models.Exercise, error) {
    filter := bson.M{
        "equipment": models.EquipmentNone,
    }
    return exerciseCollection.Find(ctx, filter)
}
```

### 3. Filtering by Muscle Group
```go
func GetLegExercises(ctx context.Context) ([]*models.Exercise, error) {
    // Get all muscles in the legs group
    legMuscles := models.GetMusclesInGroup(models.MuscleGroupLegs)
    
    filter := bson.M{
        "primaryMuscles": bson.M{
            "$in": legMuscles,
        },
    }
    return exerciseCollection.Find(ctx, filter)
}
```

### 4. Using the Exercise Filter Model
```go
func SearchExercises(ctx context.Context, searchParams ExerciseSearchParams) ([]*models.Exercise, error) {
    // Create a filter
    filter := models.NewExerciseFilter()
    filter.Categories = []models.ExerciseCategory{
        models.CategoryStrength,
        models.CategoryCalisthenics,
    }
    filter.Equipment = []models.EquipmentType{
        models.EquipmentDumbbell,
        models.EquipmentBarbell,
    }
    filter.MuscleGroups = []models.MuscleGroup{
        models.MuscleGroupChest,
    }
    
    // Build MongoDB query
    mongoFilter := buildMongoFilter(filter)
    return exerciseCollection.Find(ctx, mongoFilter)
}
```

### 5. Validation
```go
func ValidateExercise(exercise *models.Exercise) error {
    // Validate category
    if !models.IsValidCategory(exercise.Category) {
        return fmt.Errorf("invalid category: %s", exercise.Category)
    }
    
    // Validate equipment
    for _, equip := range exercise.Equipment {
        if !models.IsValidEquipment(equip) {
            return fmt.Errorf("invalid equipment: %s", equip)
        }
    }
    
    // Validate muscles
    for _, muscle := range exercise.PrimaryMuscles {
        if !models.IsValidMuscle(muscle) {
            return fmt.Errorf("invalid muscle: %s", muscle)
        }
    }
    
    return nil
}
```

### 6. Using Exercise Query for Advanced Search
```go
func AdvancedExerciseSearch(ctx context.Context) (*models.ExerciseListResponse, error) {
    query := models.NewExerciseQuery()
    
    // Set filters
    query.Filter.Categories = []models.ExerciseCategory{models.CategoryStrength}
    query.Filter.Equipment = []models.EquipmentType{models.EquipmentBarbell}
    query.Filter.SearchTerm = "bench press"
    
    // Set pagination
    query.Pagination.Page = 1
    query.Pagination.Limit = 20
    
    // Set sorting
    query.Sort.Field = "name"
    query.Sort.Order = 1 // ascending
    
    // Execute search
    return exerciseService.SearchWithFilters(ctx, query)
}
```

## Helper Functions

### Get All Values
```go
// Get all categories
categories := models.GetAllCategories()

// Get all equipment types
equipment := models.GetAllEquipment()

// Get all muscles
muscles := models.GetAllMuscles()

// Get all muscle groups
muscleGroups := models.GetAllMuscleGroups()
```

### Muscle Group Mapping
```go
// Get muscles in a group
chestMuscles := models.GetMusclesInGroup(models.MuscleGroupChest)
// Returns: []MuscleType{MuscleChe  st, MuscleSerratusAnterior}

// Get muscle group for a muscle
group := models.GetMuscleGroup(models.MuscleBiceps)
// Returns: MuscleGroupArms
```

### Filter Expansion
```go
filter := models.NewExerciseFilter()
filter.MuscleGroups = []models.MuscleGroup{models.MuscleGroupArms}
filter.PrimaryMuscles = []models.MuscleType{models.MuscleChest}

// Expand muscle groups to individual muscles
expandedMuscles := filter.ExpandMuscleGroups()
// Returns all arm muscles plus chest
```

## Best Practices

1. **Always use the constants** instead of string literals:
   ```go
   // Good
   filter["category"] = models.CategoryStrength
   
   // Bad
   filter["category"] = "strength"
   ```

2. **Validate user input** using the validation functions:
   ```go
   if !models.IsValidCategory(userInput.Category) {
       return errors.New("invalid category")
   }
   ```

3. **Use type safety** in function signatures:
   ```go
   func GetExercisesByCategory(category models.ExerciseCategory) ([]*models.Exercise, error)
   ```

4. **Leverage muscle group mapping** for broader searches:
   ```go
   // Instead of searching for each arm muscle individually
   armMuscles := models.GetMusclesInGroup(models.MuscleGroupArms)
   ```

5. **Use the filter models** for complex queries to maintain clean code:
   ```go
   query := models.NewExerciseQuery()
   // Configure query...
   ``` 