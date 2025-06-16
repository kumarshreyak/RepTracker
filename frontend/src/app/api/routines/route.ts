import { NextRequest, NextResponse } from "next/server";
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import path from 'path';

interface WorkoutSet {
  reps: number;
  weight: number;
  duration_seconds: number;
  distance: number;
  notes: string;
}

interface WorkoutExercise {
  exercise_id: string;
  sets: WorkoutSet[];
  notes: string;
  rest_seconds: number;
}

interface CreateWorkoutRequest {
  user_id: string;
  name: string;
  description: string;
  exercises: WorkoutExercise[];
  started_at: string | null;
  notes: string;
}

// gRPC client setup
const PROTO_PATH = path.join(process.cwd(), '..', 'proto', 'gymlog', 'v1', 'gymlog.proto');

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const gymlogProto = grpc.loadPackageDefinition(packageDefinition) as any;

function getWorkoutServiceClient() {
  return new gymlogProto.gymlog.v1.WorkoutService(
    'localhost:50051',
    grpc.credentials.createInsecure()
  );
}

export async function POST(request: NextRequest) {
  try {
    const body: CreateWorkoutRequest = await request.json();

    // Validate request body
    if (!body.name || !body.user_id || !body.exercises || body.exercises.length === 0) {
      return NextResponse.json(
        { error: "Missing required fields: name, user_id, and exercises are required" },
        { status: 400 }
      );
    }

    // Validate exercises
    for (let i = 0; i < body.exercises.length; i++) {
      const exercise = body.exercises[i];
      if (!exercise.exercise_id || !exercise.sets || exercise.sets.length === 0) {
        return NextResponse.json(
          { error: `Exercise ${i + 1} is missing required fields: exercise_id and sets` },
          { status: 400 }
        );
      }

      // Validate sets
      for (let j = 0; j < exercise.sets.length; j++) {
        const set = exercise.sets[j];
        if (set.reps <= 0) {
          return NextResponse.json(
            { error: `Exercise ${i + 1}, Set ${j + 1}: reps must be greater than 0` },
            { status: 400 }
          );
        }
        if (set.weight < 0) {
          return NextResponse.json(
            { error: `Exercise ${i + 1}, Set ${j + 1}: weight cannot be negative` },
            { status: 400 }
          );
        }
      }
    }

    // Create gRPC client and call CreateWorkout
    const client = getWorkoutServiceClient();
    
    const grpcRequest = {
      user_id: body.user_id,
      name: body.name,
      description: body.description,
      exercises: body.exercises,
      started_at: body.started_at ? { seconds: Math.floor(new Date(body.started_at).getTime() / 1000) } : null,
      notes: body.notes,
    };

    // Make the gRPC call
    const workout = await new Promise((resolve, reject) => {
      client.CreateWorkout(grpcRequest, (error: any, response: any) => {
        if (error) {
          console.error("gRPC Error:", error);
          reject(error);
        } else {
          resolve(response);
        }
      });
    });

    return NextResponse.json({
      success: true,
      workout: workout,
    });

  } catch (error) {
    console.error("Error creating routine:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    );
  }
}

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const userId = searchParams.get("user_id");

    if (!userId) {
      return NextResponse.json(
        { error: "user_id parameter is required" },
        { status: 400 }
      );
    }

    // Create gRPC client and call ListWorkouts
    const client = getWorkoutServiceClient();
    
    const grpcRequest = {
      user_id: userId,
      page_size: 50, // Default page size
      page_token: "",
    };

    // Make the gRPC call
    const workouts = await new Promise((resolve, reject) => {
      client.ListWorkouts(grpcRequest, (error: any, response: any) => {
        if (error) {
          console.error("gRPC Error:", error);
          reject(error);
        } else {
          resolve(response.workouts || []);
        }
      });
    });

    return NextResponse.json({
      success: true,
      workouts: workouts,
    });

  } catch (error) {
    console.error("Error fetching routines:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    );
  }
} 