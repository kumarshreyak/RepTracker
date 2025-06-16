package services

import (
	"context"
	"fmt"
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
	db           *database.MongoDB
	usersColl    *mongo.Collection
	sessionsColl *mongo.Collection
}

// NewUserService creates a new UserService instance
func NewUserService(db *database.MongoDB) *UserService {
	return &UserService{
		db:           db,
		usersColl:    db.GetCollection("users"),
		sessionsColl: db.GetCollection("sessions"),
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
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

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	update := bson.M{
		"updated_at": time.Now(),
	}

	if req.Email != "" {
		update["email"] = req.Email
	}
	if req.FirstName != "" || req.LastName != "" {
		name := fmt.Sprintf("%s %s", req.FirstName, req.LastName)
		update["name"] = name
	}

	result, err := s.usersColl.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update user")
	}
	if result.MatchedCount == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
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

// CreateOrUpdateGoogleUser creates or updates a user from Google OAuth
func (s *UserService) CreateOrUpdateGoogleUser(ctx context.Context, googleID, email, name, picture string) (*models.User, error) {
	// Try to find existing user by google_id or email
	filter := bson.M{
		"$or": []bson.M{
			{"google_id": googleID},
			{"email": email},
		},
	}

	var user models.User
	err := s.usersColl.FindOne(ctx, filter).Decode(&user)

	now := time.Now()

	if err == mongo.ErrNoDocuments {
		// Create new user
		user = models.User{
			Email:     email,
			Name:      name,
			GoogleID:  googleID,
			Picture:   picture,
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
		// Update existing user
		update := bson.M{
			"name":       name,
			"picture":    picture,
			"updated_at": now,
		}

		// Update google_id if not already set
		if user.GoogleID == "" {
			update["google_id"] = googleID
		}

		_, err = s.usersColl.UpdateOne(
			ctx,
			bson.M{"_id": user.ID},
			bson.M{"$set": update},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		user.Name = name
		user.Picture = picture
		user.UpdatedAt = now
		if user.GoogleID == "" {
			user.GoogleID = googleID
		}
	}

	return &user, nil
}

// CreateSession creates a new user session
func (s *UserService) CreateSession(ctx context.Context, userID primitive.ObjectID, accessToken, refreshToken string, expiresAt time.Time) (*models.Session, error) {
	session := models.Session{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		CreatedAt:    time.Now(),
	}

	result, err := s.sessionsColl.InsertOne(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	session.ID = result.InsertedID.(primitive.ObjectID)
	return &session, nil
}

// ValidateSession validates a session token and returns the user
func (s *UserService) ValidateSession(ctx context.Context, accessToken string) (*models.User, error) {
	var session models.Session
	err := s.sessionsColl.FindOne(ctx, bson.M{
		"access_token": accessToken,
		"expires_at":   bson.M{"$gt": time.Now()},
	}).Decode(&session)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("invalid or expired session")
	} else if err != nil {
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}

	var user models.User
	err = s.usersColl.FindOne(ctx, bson.M{"_id": session.UserID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// DeleteSession deletes a session
func (s *UserService) DeleteSession(ctx context.Context, accessToken string) error {
	_, err := s.sessionsColl.DeleteOne(ctx, bson.M{"access_token": accessToken})
	return err
}

// modelToProto converts a user model to protobuf user
func (s *UserService) modelToProto(user *models.User) *pb.User {
	return &pb.User{
		Id:        user.ID.Hex(),
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
