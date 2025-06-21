package models

// ExerciseCategory represents the category of an exercise
type ExerciseCategory string

const (
	CategoryStrength             ExerciseCategory = "strength"
	CategoryStretching           ExerciseCategory = "stretching"
	CategoryPlyometrics          ExerciseCategory = "plyometrics"
	CategoryStrongman            ExerciseCategory = "strongman"
	CategoryCardio               ExerciseCategory = "cardio"
	CategoryOlympicWeightlifting ExerciseCategory = "olympic weightlifting"
	CategoryCrossfit             ExerciseCategory = "crossfit"
	CategoryCalisthenics         ExerciseCategory = "calisthenics"
)

// EquipmentType represents the type of equipment needed for an exercise
type EquipmentType string

const (
	EquipmentNone         EquipmentType = "none"
	EquipmentEZCurlBar    EquipmentType = "ez curl bar"
	EquipmentBarbell      EquipmentType = "barbell"
	EquipmentDumbbell     EquipmentType = "dumbbell"
	EquipmentGymMat       EquipmentType = "gym mat"
	EquipmentExerciseBall EquipmentType = "exercise ball"
	EquipmentMedicineBall EquipmentType = "medicine ball"
	EquipmentPullUpBar    EquipmentType = "pull-up bar"
	EquipmentBench        EquipmentType = "bench"
	EquipmentInclineBench EquipmentType = "incline bench"
	EquipmentKettlebell   EquipmentType = "kettlebell"
	EquipmentMachine      EquipmentType = "machine"
	EquipmentCable        EquipmentType = "cable"
	EquipmentBands        EquipmentType = "bands"
	EquipmentFoamRoll     EquipmentType = "foam roll"
	EquipmentOther        EquipmentType = "other"
)

// MuscleType represents a specific muscle that can be targeted
type MuscleType string

const (
	MuscleForearms         MuscleType = "forearms"
	MuscleAbductors        MuscleType = "abductors"
	MuscleAdductors        MuscleType = "adductors"
	MuscleMiddleBack       MuscleType = "middle back"
	MuscleNeck             MuscleType = "neck"
	MuscleBiceps           MuscleType = "biceps"
	MuscleShoulders        MuscleType = "shoulders"
	MuscleSerratusAnterior MuscleType = "serratus anterior"
	MuscleChest            MuscleType = "chest"
	MuscleTriceps          MuscleType = "triceps"
	MuscleAbs              MuscleType = "abs"
	MuscleCalves           MuscleType = "calves"
	MuscleGlutes           MuscleType = "glutes"
	MuscleTraps            MuscleType = "traps"
	MuscleQuads            MuscleType = "quads"
	MuscleHamstrings       MuscleType = "hamstrings"
	MuscleLats             MuscleType = "lats"
	MuscleBrachialis       MuscleType = "brachialis"
	MuscleObliques         MuscleType = "obliques"
	MuscleSoleus           MuscleType = "soleus"
	MuscleLowerBack        MuscleType = "lower back"
	MuscleRhomboids        MuscleType = "rhomboids"
	MuscleHipFlexors       MuscleType = "hip flexors"
)

// MuscleGroup represents a group of related muscles
type MuscleGroup string

const (
	MuscleGroupArms      MuscleGroup = "arms"
	MuscleGroupBack      MuscleGroup = "back"
	MuscleGroupCalves    MuscleGroup = "calves"
	MuscleGroupChest     MuscleGroup = "chest"
	MuscleGroupCore      MuscleGroup = "core"
	MuscleGroupLegs      MuscleGroup = "legs"
	MuscleGroupShoulders MuscleGroup = "shoulders"
)

// MuscleGroupMapping maps muscle groups to their constituent muscles
var MuscleGroupMapping = map[MuscleGroup][]MuscleType{
	MuscleGroupArms: {
		MuscleForearms,
		MuscleBiceps,
		MuscleTriceps,
		MuscleBrachialis,
	},
	MuscleGroupBack: {
		MuscleNeck,
		MuscleTraps,
		MuscleLats,
		MuscleLowerBack,
		MuscleMiddleBack,
	},
	MuscleGroupCalves: {
		MuscleCalves,
		MuscleSoleus,
	},
	MuscleGroupChest: {
		MuscleChest,
		MuscleSerratusAnterior,
	},
	MuscleGroupCore: {
		MuscleAbs,
		MuscleObliques,
	},
	MuscleGroupLegs: {
		MuscleAbductors,
		MuscleAdductors,
		MuscleQuads,
		MuscleHamstrings,
		MuscleGlutes,
	},
	MuscleGroupShoulders: {
		MuscleShoulders,
	},
}

// GetAllCategories returns all available exercise categories
func GetAllCategories() []ExerciseCategory {
	return []ExerciseCategory{
		CategoryStrength,
		CategoryStretching,
		CategoryPlyometrics,
		CategoryStrongman,
		CategoryCardio,
		CategoryOlympicWeightlifting,
		CategoryCrossfit,
		CategoryCalisthenics,
	}
}

// GetAllEquipment returns all available equipment types
func GetAllEquipment() []EquipmentType {
	return []EquipmentType{
		EquipmentNone,
		EquipmentEZCurlBar,
		EquipmentBarbell,
		EquipmentDumbbell,
		EquipmentGymMat,
		EquipmentExerciseBall,
		EquipmentMedicineBall,
		EquipmentPullUpBar,
		EquipmentBench,
		EquipmentInclineBench,
		EquipmentKettlebell,
		EquipmentMachine,
		EquipmentCable,
		EquipmentBands,
		EquipmentFoamRoll,
		EquipmentOther,
	}
}

// GetAllMuscles returns all available muscle types
func GetAllMuscles() []MuscleType {
	return []MuscleType{
		MuscleForearms,
		MuscleAbductors,
		MuscleAdductors,
		MuscleMiddleBack,
		MuscleNeck,
		MuscleBiceps,
		MuscleShoulders,
		MuscleSerratusAnterior,
		MuscleChest,
		MuscleTriceps,
		MuscleAbs,
		MuscleCalves,
		MuscleGlutes,
		MuscleTraps,
		MuscleQuads,
		MuscleHamstrings,
		MuscleLats,
		MuscleBrachialis,
		MuscleObliques,
		MuscleSoleus,
		MuscleLowerBack,
		MuscleRhomboids,
		MuscleHipFlexors,
	}
}

// GetAllMuscleGroups returns all available muscle groups
func GetAllMuscleGroups() []MuscleGroup {
	return []MuscleGroup{
		MuscleGroupArms,
		MuscleGroupBack,
		MuscleGroupCalves,
		MuscleGroupChest,
		MuscleGroupCore,
		MuscleGroupLegs,
		MuscleGroupShoulders,
	}
}

// IsValidCategory checks if the given category is valid
func IsValidCategory(category string) bool {
	for _, c := range GetAllCategories() {
		if string(c) == category {
			return true
		}
	}
	return false
}

// IsValidEquipment checks if the given equipment type is valid
func IsValidEquipment(equipment string) bool {
	for _, e := range GetAllEquipment() {
		if string(e) == equipment {
			return true
		}
	}
	return false
}

// IsValidMuscle checks if the given muscle type is valid
func IsValidMuscle(muscle string) bool {
	for _, m := range GetAllMuscles() {
		if string(m) == muscle {
			return true
		}
	}
	return false
}

// IsValidMuscleGroup checks if the given muscle group is valid
func IsValidMuscleGroup(group string) bool {
	for _, g := range GetAllMuscleGroups() {
		if string(g) == group {
			return true
		}
	}
	return false
}

// GetMusclesInGroup returns all muscles in a given muscle group
func GetMusclesInGroup(group MuscleGroup) []MuscleType {
	return MuscleGroupMapping[group]
}

// GetMuscleGroup returns the muscle group for a given muscle
func GetMuscleGroup(muscle MuscleType) MuscleGroup {
	for group, muscles := range MuscleGroupMapping {
		for _, m := range muscles {
			if m == muscle {
				return group
			}
		}
	}
	return ""
}
