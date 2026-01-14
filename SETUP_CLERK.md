# GymLog Setup Guide with Clerk Authentication

## Prerequisites

- Node.js 18+ and npm
- Go 1.23+
- MongoDB (local or cloud)
- Clerk account (free tier available)

## Step 1: Clerk Setup (5 minutes)

### 1.1 Create Clerk Account

1. Go to https://clerk.com
2. Sign up for a free account
3. Create a new application named "GymLog"
4. Select **Email** as the authentication method
5. Click **Create Application**

### 1.2 Get Your Clerk Keys

1. In the Clerk dashboard, go to **API Keys**
2. Copy your **Publishable Key** (starts with `pk_test_`)
3. Keep this key handy - you'll need it for both frontend and backend

## Step 2: Backend Setup

### 2.1 Create Environment File

Create a `.env` file in the `backend/` directory:

```bash
cd backend
cat > .env << 'EOF'
# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
DB_NAME=gymlog

# Clerk Configuration
EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_YOUR_KEY_HERE

# Google Gemini API Configuration (optional for now)
GEMINI_API_KEY=your-gemini-api-key-here
EOF
```

**Replace `pk_test_YOUR_KEY_HERE` with your actual Clerk publishable key!**

### 2.2 Install Dependencies

```bash
go mod tidy
```

### 2.3 Start MongoDB

If running MongoDB locally:

```bash
mongod
```

Or use MongoDB Atlas (cloud) and update the `MONGODB_URI` in `.env`.

### 2.4 Initialize Database

```bash
# Initialize MongoDB collections
mongosh < scripts/init-mongo.js

# (Optional) Populate exercises
go run scripts/populate_exercises.go
```

### 2.5 Start Backend Server

```bash
go run cmd/server/main.go
```

The server will start on:
- HTTP: `http://localhost:8080`
- gRPC: `localhost:50051`

## Step 3: Frontend Setup

### 3.1 Create Environment File

Create a `.env` file in the `frontend-native/` directory:

```bash
cd frontend-native
cat > .env << 'EOF'
EXPO_PUBLIC_API_BASE_URL=http://localhost:8080
EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_YOUR_KEY_HERE
EOF
```

**Use the same Clerk publishable key as the backend!**

### 3.2 Install Dependencies

```bash
npm install
```

### 3.3 Start Development Server

```bash
npm start
```

This will start the Expo development server. You can then:
- Press `i` for iOS simulator
- Press `a` for Android emulator
- Scan QR code with Expo Go app on your phone

## Step 4: Test Authentication

### 4.1 Sign Up

1. Open the app
2. Tap **Sign Up**
3. Enter your email and password
4. Check your email for verification code
5. Enter the 6-digit code
6. Complete onboarding (height, weight, age, goal)

### 4.2 Sign In

1. Open the app
2. Tap **Sign In**
3. Enter your email and password
4. You should be logged in and see the home screen

### 4.3 Test API Calls

1. Create a workout routine
2. Start a workout session
3. Check the Insights tab

All API calls should work with Clerk authentication!

## Troubleshooting

### Backend Issues

**Error: "EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY environment variable not set"**
- Make sure you created the `.env` file in the `backend/` directory
- Verify the key is correct and starts with `pk_test_`

**Error: "Failed to connect to MongoDB"**
- Make sure MongoDB is running
- Check the `MONGODB_URI` in `.env`

**Error: "Failed to fetch JWKS"**
- Check your internet connection
- Verify the Clerk publishable key is correct

### Frontend Issues

**Error: "Missing EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY"**
- Make sure you created the `.env` file in the `frontend-native/` directory
- Restart the Expo development server after creating `.env`

**Error: "Network request failed"**
- Make sure the backend server is running
- Check the `EXPO_PUBLIC_API_BASE_URL` in `.env`
- For Android emulator, use `http://10.0.2.2:8080` instead of `localhost`

**Email verification not received**
- Check spam folder
- Verify email settings in Clerk dashboard
- Try resending the verification email

### Authentication Issues

**Error: "Invalid or expired token"**
- Sign out and sign in again
- Check that both frontend and backend use the same Clerk publishable key
- Verify your internet connection (Clerk needs to validate tokens)

**Error: "User not found"**
- The user record should be created automatically on first API call
- Check backend logs for errors
- Verify MongoDB connection

## Environment Variables Reference

### Backend `.env`

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `MONGODB_URI` | Yes | MongoDB connection string | `mongodb://localhost:27017` |
| `DB_NAME` | Yes | MongoDB database name | `gymlog` |
| `EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY` | Yes | Clerk publishable key | `pk_test_...` |
| `GEMINI_API_KEY` | No | Google Gemini API key (for AI features) | `AIza...` |

### Frontend `.env`

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `EXPO_PUBLIC_API_BASE_URL` | Yes | Backend API URL | `http://localhost:8080` |
| `EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY` | Yes | Clerk publishable key | `pk_test_...` |

## Production Deployment

### Backend

1. Update `.env` with production values:
   - Use production MongoDB URI
   - Use production Clerk publishable key (starts with `pk_live_`)
   - Set production Gemini API key

2. Build and deploy:
   ```bash
   go build -o gymlog-server cmd/server/main.go
   ./gymlog-server
   ```

### Frontend

1. Update `.env` with production values:
   - Use production API URL
   - Use production Clerk publishable key (starts with `pk_live_`)

2. Build for production:
   ```bash
   # iOS
   eas build --platform ios --profile production
   
   # Android
   eas build --platform android --profile production
   ```

## Next Steps

- [ ] Enable social OAuth (Google, GitHub) in Clerk dashboard
- [ ] Add multi-factor authentication in Clerk dashboard
- [ ] Set up Clerk webhooks for user events
- [ ] Configure custom email templates in Clerk
- [ ] Add user roles and permissions

## Support

- **Clerk Documentation**: https://clerk.com/docs
- **Clerk Dashboard**: https://dashboard.clerk.com
- **GymLog Migration Guide**: See `CLERK_MIGRATION.md`

## Quick Commands Reference

### Start Everything

```bash
# Terminal 1 - MongoDB
mongod

# Terminal 2 - Backend
cd backend
go run cmd/server/main.go

# Terminal 3 - Frontend
cd frontend-native
npm start
```

### Reset Everything

```bash
# Drop database
mongosh gymlog --eval "db.dropDatabase()"

# Reinitialize
mongosh < backend/scripts/init-mongo.js

# Restart servers
```

That's it! You're now running GymLog with Clerk authentication. 🎉

