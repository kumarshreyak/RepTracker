import { NextRequest, NextResponse } from 'next/server';

// Mock data for now - this will be replaced with actual gRPC calls
const mockExercises = [
  {
    id: "1",
    name: "Push-ups",
    description: "A classic upper body exercise that works the chest, shoulders, and triceps.",
    muscle_group: "Chest",
    equipment: "Bodyweight",
    difficulty: "Beginner",
    instructions: [
      "Start in plank position with hands under shoulders",
      "Lower your body until chest nearly touches floor",
      "Push back up to starting position",
      "Repeat for desired reps"
    ]
  },
  {
    id: "2",
    name: "Squats",
    description: "A fundamental lower body exercise targeting the quadriceps, glutes, and hamstrings.",
    muscle_group: "Legs",
    equipment: "Bodyweight",
    difficulty: "Beginner",
    instructions: [
      "Stand with feet shoulder-width apart",
      "Lower your body as if sitting back into a chair",
      "Keep knees behind toes and chest up",
      "Return to starting position"
    ]
  },
  {
    id: "3",
    name: "Bench Press",
    description: "A compound exercise for building chest, shoulder, and tricep strength.",
    muscle_group: "Chest",
    equipment: "Barbell",
    difficulty: "Intermediate",
    instructions: [
      "Lie on bench with feet flat on floor",
      "Grip barbell slightly wider than shoulder width",
      "Lower bar to chest with control",
      "Press bar back up to starting position"
    ]
  },
  {
    id: "4",
    name: "Deadlifts",
    description: "A compound movement that works the entire posterior chain.",
    muscle_group: "Back",
    equipment: "Barbell",
    difficulty: "Advanced",
    instructions: [
      "Stand with feet hip-width apart, bar over mid-foot",
      "Hinge at hips and knees to grab bar",
      "Keep back straight and chest up",
      "Drive through heels to stand up tall"
    ]
  },
  {
    id: "5",
    name: "Pull-ups",
    description: "An upper body exercise that primarily targets the latissimus dorsi.",
    muscle_group: "Back",
    equipment: "Pull-up Bar",
    difficulty: "Intermediate",
    instructions: [
      "Hang from bar with palms facing away",
      "Pull your body up until chin clears bar",
      "Lower yourself with control",
      "Repeat for desired reps"
    ]
  },
  {
    id: "6",
    name: "Shoulder Press",
    description: "An overhead pressing movement for shoulder and tricep development.",
    muscle_group: "Shoulders",
    equipment: "Dumbbells",
    difficulty: "Intermediate",
    instructions: [
      "Stand with feet shoulder-width apart",
      "Hold dumbbells at shoulder height",
      "Press weights overhead until arms are extended",
      "Lower back to starting position"
    ]
  },
  {
    id: "7",
    name: "Bicep Curls",
    description: "An isolation exercise targeting the biceps brachii.",
    muscle_group: "Arms",
    equipment: "Dumbbells",
    difficulty: "Beginner",
    instructions: [
      "Stand with dumbbells at your sides",
      "Keep elbows close to your body",
      "Curl weights up by contracting biceps",
      "Lower weights back to starting position"
    ]
  },
  {
    id: "8",
    name: "Planks",
    description: "An isometric core exercise for building abdominal strength.",
    muscle_group: "Core",
    equipment: "Bodyweight",
    difficulty: "Beginner",
    instructions: [
      "Start in push-up position",
      "Lower to forearms, keeping body straight",
      "Hold position while engaging core",
      "Breathe normally throughout"
    ]
  }
];

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const search = searchParams.get('search')?.toLowerCase();
    const muscleGroup = searchParams.get('muscle_group')?.toLowerCase();
    const equipment = searchParams.get('equipment')?.toLowerCase();

    let filteredExercises = mockExercises;

    // Apply search filter
    if (search) {
      filteredExercises = filteredExercises.filter(exercise =>
        exercise.name.toLowerCase().includes(search) ||
        exercise.description.toLowerCase().includes(search) ||
        exercise.muscle_group.toLowerCase().includes(search) ||
        exercise.equipment.toLowerCase().includes(search)
      );
    }

    // Apply muscle group filter
    if (muscleGroup) {
      filteredExercises = filteredExercises.filter(exercise =>
        exercise.muscle_group.toLowerCase().includes(muscleGroup)
      );
    }

    // Apply equipment filter
    if (equipment) {
      filteredExercises = filteredExercises.filter(exercise =>
        exercise.equipment.toLowerCase().includes(equipment)
      );
    }

    return NextResponse.json({
      exercises: filteredExercises,
      total: filteredExercises.length
    });

  } catch (error) {
    console.error('Error fetching exercises:', error);
    return NextResponse.json(
      { error: 'Failed to fetch exercises' },
      { status: 500 }
    );
  }
}

// TODO: Replace mock data with actual gRPC calls
// Example implementation for when gRPC backend is ready:
/*
import * as grpc from '@grpc/grpc-js';
import { ExerciseServiceClient } from '../../../proto/gymlog/v1/gymlog_grpc_pb';
import { ListExercisesRequest } from '../../../proto/gymlog/v1/gymlog_pb';

const client = new ExerciseServiceClient(
  'localhost:50051',
  grpc.credentials.createInsecure()
);

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const search = searchParams.get('search');
    const muscleGroup = searchParams.get('muscle_group');
    const equipment = searchParams.get('equipment');

    const grpcRequest = new ListExercisesRequest();
    if (muscleGroup) grpcRequest.setMuscleGroup(muscleGroup);
    if (equipment) grpcRequest.setEquipment(equipment);
    grpcRequest.setPageSize(50);

    const response = await new Promise((resolve, reject) => {
      client.listExercises(grpcRequest, (error, response) => {
        if (error) reject(error);
        else resolve(response);
      });
    });

    let exercises = response.getExercisesList().map(exercise => ({
      id: exercise.getId(),
      name: exercise.getName(),
      description: exercise.getDescription(),
      muscle_group: exercise.getMuscleGroup(),
      equipment: exercise.getEquipment(),
      difficulty: exercise.getDifficulty(),
      instructions: exercise.getInstructionsList()
    }));

    // Apply client-side search filter if needed
    if (search) {
      exercises = exercises.filter(exercise =>
        exercise.name.toLowerCase().includes(search.toLowerCase()) ||
        exercise.description.toLowerCase().includes(search.toLowerCase())
      );
    }

    return NextResponse.json({
      exercises,
      total: exercises.length
    });

  } catch (error) {
    console.error('Error fetching exercises:', error);
    return NextResponse.json(
      { error: 'Failed to fetch exercises' },
      { status: 500 }
    );
  }
}
*/ 