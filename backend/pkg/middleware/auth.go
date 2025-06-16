package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"gymlog-backend/internal/services"
)

// AuthInterceptor provides authentication for gRPC services
func AuthInterceptor(userService *services.UserService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip auth for certain methods
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract token from metadata
		token, err := extractToken(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing or invalid token")
		}

		// Validate session
		user, err := userService.ValidateSession(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired session")
		}

		// Add user info to context
		ctx = context.WithValue(ctx, "user_id", user.ID.Hex())
		ctx = context.WithValue(ctx, "user_email", user.Email)

		return handler(ctx, req)
	}
}

func isPublicMethod(method string) bool {
	publicMethods := []string{
		"/gymlog.v1.UserService/CreateUser",
		// Add other public methods here
	}

	for _, publicMethod := range publicMethods {
		if method == publicMethod {
			return true
		}
	}
	return false
}

func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization header")
	}

	token := authHeader[0]
	if !strings.HasPrefix(token, "Bearer ") {
		return "", status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	return strings.TrimPrefix(token, "Bearer "), nil
}
