package models

// ExerciseFilter represents filters that can be applied when querying exercises
type ExerciseFilter struct {
	Categories       []ExerciseCategory `json:"categories,omitempty"`
	Equipment        []EquipmentType    `json:"equipment,omitempty"`
	PrimaryMuscles   []MuscleType       `json:"primaryMuscles,omitempty"`
	SecondaryMuscles []MuscleType       `json:"secondaryMuscles,omitempty"`
	MuscleGroups     []MuscleGroup      `json:"muscleGroups,omitempty"`
	SearchTerm       string             `json:"searchTerm,omitempty"`
}

// ExerciseSort represents sorting options for exercise queries
type ExerciseSort struct {
	Field string `json:"field"`
	Order int    `json:"order"` // 1 for ascending, -1 for descending
}

// ExercisePagination represents pagination options
type ExercisePagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// ExerciseQuery combines all query options
type ExerciseQuery struct {
	Filter     ExerciseFilter     `json:"filter"`
	Sort       ExerciseSort       `json:"sort"`
	Pagination ExercisePagination `json:"pagination"`
}

// ExerciseListResponse represents a paginated response of exercises
type ExerciseListResponse struct {
	Exercises   []Exercise `json:"exercises"`
	Total       int64      `json:"total"`
	Page        int        `json:"page"`
	Limit       int        `json:"limit"`
	TotalPages  int        `json:"totalPages"`
	HasNextPage bool       `json:"hasNextPage"`
	HasPrevPage bool       `json:"hasPrevPage"`
}

// ExerciseStats represents statistics about exercises in the database
type ExerciseStats struct {
	TotalExercises      int              `json:"totalExercises"`
	ByCategory          map[string]int   `json:"byCategory"`
	ByEquipment         map[string]int   `json:"byEquipment"`
	ByMuscleGroup       map[string]int   `json:"byMuscleGroup"`
	MostCommonEquipment []EquipmentCount `json:"mostCommonEquipment"`
	MostTargetedMuscles []MuscleCount    `json:"mostTargetedMuscles"`
}

// EquipmentCount represents equipment usage count
type EquipmentCount struct {
	Equipment EquipmentType `json:"equipment"`
	Count     int           `json:"count"`
}

// MuscleCount represents muscle targeting count
type MuscleCount struct {
	Muscle MuscleType `json:"muscle"`
	Count  int        `json:"count"`
}

// NewExerciseFilter creates a new exercise filter with default values
func NewExerciseFilter() *ExerciseFilter {
	return &ExerciseFilter{
		Categories:       []ExerciseCategory{},
		Equipment:        []EquipmentType{},
		PrimaryMuscles:   []MuscleType{},
		SecondaryMuscles: []MuscleType{},
		MuscleGroups:     []MuscleGroup{},
		SearchTerm:       "",
	}
}

// NewExerciseQuery creates a new exercise query with default values
func NewExerciseQuery() *ExerciseQuery {
	return &ExerciseQuery{
		Filter: *NewExerciseFilter(),
		Sort: ExerciseSort{
			Field: "name",
			Order: 1,
		},
		Pagination: ExercisePagination{
			Page:  1,
			Limit: 20,
		},
	}
}

// HasFilters checks if any filters are applied
func (f *ExerciseFilter) HasFilters() bool {
	return len(f.Categories) > 0 ||
		len(f.Equipment) > 0 ||
		len(f.PrimaryMuscles) > 0 ||
		len(f.SecondaryMuscles) > 0 ||
		len(f.MuscleGroups) > 0 ||
		f.SearchTerm != ""
}

// ExpandMuscleGroups expands muscle groups into individual muscles
func (f *ExerciseFilter) ExpandMuscleGroups() []MuscleType {
	muscleSet := make(map[MuscleType]bool)

	// Add primary muscles
	for _, muscle := range f.PrimaryMuscles {
		muscleSet[muscle] = true
	}

	// Expand muscle groups
	for _, group := range f.MuscleGroups {
		muscles := GetMusclesInGroup(group)
		for _, muscle := range muscles {
			muscleSet[muscle] = true
		}
	}

	// Convert set to slice
	result := make([]MuscleType, 0, len(muscleSet))
	for muscle := range muscleSet {
		result = append(result, muscle)
	}

	return result
}
