package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gymlog-backend/pkg/database"
	"gymlog-backend/pkg/models"
	pb "gymlog-backend/proto/gymlog/v1"
)

// UserService implements the gRPC UserService
type UserService struct {
	pb.UnimplementedUserServiceServer
	db        *database.MongoDB
	usersColl *mongo.Collection
}

// NewUserService creates a new UserService instance
func NewUserService(db *database.MongoDB) *UserService {
	return &UserService{
		db:        db,
		usersColl: db.GetCollection("users"),
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Validate goal if provided
	if req.Goal != "" && req.Goal != models.GoalLoseFat && req.Goal != models.GoalGainMuscle && req.Goal != models.GoalMaintain {
		return nil, status.Error(codes.InvalidArgument, "invalid goal: must be lose_fat, gain_muscle, or maintain")
	}

	// Check if user already exists
	var existingUser models.User
	err := s.usersColl.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	} else if err != mongo.ErrNoDocuments {
		return nil, status.Error(codes.Internal, "failed to check existing user")
	}

	now := time.Now()
	user := models.User{
		Email:     req.Email,
		Name:      fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		Height:    req.Height,
		Weight:    req.Weight,
		Age:       req.Age,
		Goal:      req.Goal,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := s.usersColl.InsertOne(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return &pb.User{
		Id:        user.ID.Hex(),
		Email:     user.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Height:    user.Height,
		Weight:    user.Weight,
		Age:       user.Age,
		Goal:      user.Goal,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	var user models.User
	err = s.usersColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "user not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return s.modelToProto(&user), nil
}

// UpdateUser updates an existing user or creates one if it doesn't exist (upsert)
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	// Validate goal if provided
	if req.Goal != "" && req.Goal != models.GoalLoseFat && req.Goal != models.GoalGainMuscle && req.Goal != models.GoalMaintain {
		return nil, status.Error(codes.InvalidArgument, "invalid goal: must be lose_fat, gain_muscle, or maintain")
	}

	// Check if user exists
	var existingUser models.User
	err = s.usersColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingUser)
	userExists := err == nil

	now := time.Now()

	if userExists {
		// Update existing user
		update := bson.M{
			"updated_at": now,
		}

		if req.Email != "" {
			update["email"] = req.Email
		}
		if req.FirstName != "" || req.LastName != "" {
			name := fmt.Sprintf("%s %s", req.FirstName, req.LastName)
			update["name"] = name
		}
		if req.Height > 0 {
			update["height"] = req.Height
		}
		if req.Weight > 0 {
			update["weight"] = req.Weight
		}
		if req.Age > 0 {
			update["age"] = req.Age
		}
		if req.Goal != "" {
			update["goal"] = req.Goal
		}

		_, err := s.usersColl.UpdateOne(
			ctx,
			bson.M{"_id": objectID},
			bson.M{"$set": update},
		)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to update user")
		}
	} else {
		// Create new user with the provided ID
		user := models.User{
			ID:        objectID,
			Email:     req.Email,
			Name:      fmt.Sprintf("%s %s", req.FirstName, req.LastName),
			Height:    req.Height,
			Weight:    req.Weight,
			Age:       req.Age,
			Goal:      req.Goal,
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err := s.usersColl.InsertOne(ctx, user)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to create user")
		}
	}

	return s.GetUser(ctx, &pb.GetUserRequest{Id: req.Id})
}

// ListUsers lists users with pagination
func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	pageSize := int64(req.PageSize)
	if pageSize <= 0 {
		pageSize = 50
	}

	cursor, err := s.usersColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, status.Error(codes.Internal, "failed to decode users")
	}

	var protoUsers []*pb.User
	for _, user := range users {
		protoUsers = append(protoUsers, s.modelToProto(&user))
	}

	return &pb.ListUsersResponse{
		Users: protoUsers,
	}, nil
}

// CreateOrUpdateClerkUser creates or updates a user from Clerk authentication with full profile data
func (s *UserService) CreateOrUpdateClerkUser(ctx context.Context, clerkID, email, name, picture string, height float64, weight float64, age int32, goal string) (*models.User, error) {
	// Try to find existing user by clerk_id
	filter := bson.M{"clerk_id": clerkID}

	var user models.User
	err := s.usersColl.FindOne(ctx, filter).Decode(&user)

	now := time.Now()

	if err == mongo.ErrNoDocuments {
		// Create new user with all profile data
		user = models.User{
			ClerkID:   clerkID,
			Email:     email,
			Name:      name,
			Picture:   picture,
			Height:    height,
			Weight:    weight,
			Age:       age,
			Goal:      goal,
			CreatedAt: now,
			UpdatedAt: now,
		}

		result, err := s.usersColl.InsertOne(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		user.ID = result.InsertedID.(primitive.ObjectID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	} else {
		// Update existing user with all provided data
		update := bson.M{
			"name":       name,
			"email":      email,
			"picture":    picture,
			"updated_at": now,
		}

		// Only update profile fields if they are provided (non-zero values)
		if height > 0 {
			update["height"] = height
		}
		if weight > 0 {
			update["weight"] = weight
		}
		if age > 0 {
			update["age"] = age
		}
		if goal != "" {
			update["goal"] = goal
		}

		_, err = s.usersColl.UpdateOne(
			ctx,
			bson.M{"_id": user.ID},
			bson.M{"$set": update},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		// Update local user object with the new values
		user.Name = name
		user.Email = email
		user.Picture = picture
		if height > 0 {
			user.Height = height
		}
		if weight > 0 {
			user.Weight = weight
		}
		if age > 0 {
			user.Age = age
		}
		if goal != "" {
			user.Goal = goal
		}
		user.UpdatedAt = now
	}

	return &user, nil
}

// GetUserByClerkID retrieves a user by their Clerk ID
func (s *UserService) GetUserByClerkID(ctx context.Context, clerkID string) (*models.User, error) {
	var user models.User
	err := s.usersColl.FindOne(ctx, bson.M{"clerk_id": clerkID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ObjectID
func (s *UserService) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.usersColl.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// modelToProto converts a user model to protobuf user
func (s *UserService) modelToProto(user *models.User) *pb.User {
	// Parse name to get first and last name
	nameParts := strings.Split(user.Name, " ")
	firstName := ""
	lastName := ""
	if len(nameParts) > 0 {
		firstName = nameParts[0]
		if len(nameParts) > 1 {
			lastName = strings.Join(nameParts[1:], " ")
		}
	}

	return &pb.User{
		Id:        user.ID.Hex(),
		Email:     user.Email,
		FirstName: firstName,
		LastName:  lastName,
		Height:    user.Height,
		Weight:    user.Weight,
		Age:       user.Age,
		Goal:      user.Goal,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
