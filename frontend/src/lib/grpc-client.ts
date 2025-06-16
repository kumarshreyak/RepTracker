import * as grpc from '@grpc/grpc-js'
import * as protoLoader from '@grpc/proto-loader'
import { join } from 'path'

// Define the proto file path
const PROTO_PATH = join(process.cwd(), '../../proto/gymlog/v1/gymlog.proto')

// Load the protobuf
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
})

const gymlogProto = grpc.loadPackageDefinition(packageDefinition).gymlog.v1 as any

// Create gRPC client
export function createGrpcClient() {
  const serverUrl = process.env.GRPC_SERVER_URL || 'localhost:50051'
  
  return {
    userService: new gymlogProto.UserService(serverUrl, grpc.credentials.createInsecure()),
    exerciseService: new gymlogProto.ExerciseService(serverUrl, grpc.credentials.createInsecure()),
    workoutService: new gymlogProto.WorkoutService(serverUrl, grpc.credentials.createInsecure()),
  }
}

// Utility function for Google user creation/update
export async function createOrUpdateGoogleUser(
  googleId: string,
  email: string,
  name: string,
  picture: string,
  accessToken: string,
  refreshToken?: string,
  expiresAt?: string
) {
  const client = createGrpcClient()
  
  return new Promise((resolve, reject) => {
    // This would be a custom method we'd need to add to the proto
    // For now, let's use the standard CreateUser method
    client.userService.CreateUser({
      email,
      username: email.split('@')[0],
      first_name: name.split(' ')[0] || '',
      last_name: name.split(' ').slice(1).join(' ') || '',
    }, (error: any, response: any) => {
      if (error) {
        reject(error)
      } else {
        resolve(response)
      }
    })
  })
} 