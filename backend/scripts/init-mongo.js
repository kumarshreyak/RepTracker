// MongoDB initialization script for GymLog
print('Initializing GymLog database...');

// Switch to the gymlog database
db = db.getSiblingDB('gymlog');

// Create collections
db.createCollection('users');
db.createCollection('sessions');
db.createCollection('exercises');
db.createCollection('workouts');
db.createCollection('workout_sessions');

// Create indexes for better performance
print('Creating indexes...');

// Users collection indexes
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "google_id": 1 }, { unique: true });

// Sessions collection indexes
db.sessions.createIndex({ "user_id": 1 });
db.sessions.createIndex({ "access_token": 1 }, { unique: true });
db.sessions.createIndex({ "expires_at": 1 }, { expireAfterSeconds: 0 });

// Exercises collection indexes
db.exercises.createIndex({ "name": 1 });
db.exercises.createIndex({ "muscle_group": 1 });
db.exercises.createIndex({ "equipment": 1 });

// Workouts collection indexes
db.workouts.createIndex({ "user_id": 1 });
db.workouts.createIndex({ "name": 1 });
db.workouts.createIndex({ "created_at": -1 });

// Workout sessions collection indexes
db.workout_sessions.createIndex({ "user_id": 1 });
db.workout_sessions.createIndex({ "routine_id": 1 });
db.workout_sessions.createIndex({ "is_active": 1 });
db.workout_sessions.createIndex({ "created_at": -1 });

print('GymLog database initialized successfully!'); 